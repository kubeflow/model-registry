package server

import (
	"context"
	"fmt"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type TypeKind int32

// artifact type values from ml-metadata table values
const (
	EXECUTION_TYPE TypeKind = iota
	ARTIFACT_TYPE
	CONTEXT_TYPE
)

func (tk TypeKind) String() string {
	switch tk {
	case EXECUTION_TYPE:
		return "Execution"
	case ARTIFACT_TYPE:
		return "Artifact"
	case CONTEXT_TYPE:
		return "Context"
	}
	return "unknown"
}

type grpcServer struct {
	proto.UnimplementedMetadataStoreServiceServer
	dbConnection *gorm.DB
}

var _ proto.MetadataStoreServiceServer = grpcServer{}
var _ proto.MetadataStoreServiceServer = (*grpcServer)(nil)

func NewGrpcServer(dbConnection *gorm.DB) proto.MetadataStoreServiceServer {
	return &grpcServer{dbConnection: dbConnection}
}

var REQUIRED_TYPE_FIELDS = []string{"name"}

func (g grpcServer) PutArtifactType(ctx context.Context, request *proto.PutArtifactTypeRequest) (resp *proto.PutArtifactTypeResponse, err error) {
	ctx, _ = Begin(ctx, g.dbConnection)
	defer handleTransaction(ctx, &err)

	artifactType := request.GetArtifactType()
	properties := request.ArtifactType.Properties
	err = requiredFields(REQUIRED_TYPE_FIELDS, artifactType.Name)
	if err != nil {
		return nil, err
	}
	value := &db.Type{
		Name:        *artifactType.Name,
		Version:     artifactType.Version,
		TypeKind:    int32(ARTIFACT_TYPE),
		Description: artifactType.Description,
		ExternalID:  artifactType.ExternalId,
	}
	err = g.createOrUpdateType(ctx, value, properties)
	if err != nil {
		return nil, err
	}
	var typeId = int64(value.ID)
	return &proto.PutArtifactTypeResponse{
		TypeId: &typeId,
	}, nil
}

func (g grpcServer) createOrUpdateType(ctx context.Context, value *db.Type,
	properties map[string]proto.PropertyType) error {
	// TODO handle CanAdd, CanOmit properties from type request
	dbConn, _ := FromContext(ctx)

	if err := dbConn.Where("name = ?", value.Name).Assign(value).FirstOrCreate(value).Error; err != nil {
		err = fmt.Errorf("error creating type %s: %v", value.Name, err)
		return err
	}
	err := g.createTypeProperties(ctx, properties, value.ID)
	if err != nil {
		return err
	}
	return nil
}

func (g grpcServer) PutExecutionType(ctx context.Context, request *proto.PutExecutionTypeRequest) (resp *proto.PutExecutionTypeResponse, err error) {
	ctx, _ = Begin(ctx, g.dbConnection)
	defer handleTransaction(ctx, &err)

	executionType := request.GetExecutionType()
	err = requiredFields(REQUIRED_TYPE_FIELDS, executionType.Name)
	if err != nil {
		return nil, err
	}
	value := &db.Type{
		Name:        *executionType.Name,
		Version:     executionType.Version,
		TypeKind:    int32(EXECUTION_TYPE),
		Description: executionType.Description,
		ExternalID:  executionType.ExternalId,
	}
	err = g.createOrUpdateType(ctx, value, executionType.Properties)
	if err != nil {
		return nil, err
	}
	var typeId = int64(value.ID)
	return &proto.PutExecutionTypeResponse{
		TypeId: &typeId,
	}, nil
}

func (g grpcServer) PutContextType(ctx context.Context, request *proto.PutContextTypeRequest) (resp *proto.PutContextTypeResponse, err error) {
	ctx, _ = Begin(ctx, g.dbConnection)
	defer handleTransaction(ctx, &err)

	contextType := request.GetContextType()
	err = requiredFields(REQUIRED_TYPE_FIELDS, contextType.Name)
	if err != nil {
		return nil, err
	}
	value := &db.Type{
		Name:        *contextType.Name,
		Version:     contextType.Version,
		TypeKind:    int32(CONTEXT_TYPE),
		Description: contextType.Description,
		ExternalID:  contextType.ExternalId,
	}
	err = g.createOrUpdateType(ctx, value, contextType.Properties)
	if err != nil {
		return nil, err
	}
	var typeId = int64(value.ID)
	return &proto.PutContextTypeResponse{
		TypeId: &typeId,
	}, nil
}

func (g grpcServer) PutTypes(ctx context.Context, request *proto.PutTypesRequest) (resp *proto.PutTypesResponse, err error) {
	ctx, _ = Begin(ctx, g.dbConnection)
	defer handleTransaction(ctx, &err)

	response := &proto.PutTypesResponse{}

	for _, ar := range request.ArtifactTypes {
		var at *proto.PutArtifactTypeResponse
		at, err = g.PutArtifactType(ctx, &proto.PutArtifactTypeRequest{
			ArtifactType:       ar,
			CanAddFields:       request.CanAddFields,
			CanOmitFields:      request.CanOmitFields,
			TransactionOptions: request.TransactionOptions,
		})
		if err != nil {
			return response, err
		}
		response.ArtifactTypeIds = append(response.ArtifactTypeIds, *at.TypeId)
	}
	for _, ex := range request.ExecutionTypes {
		var er *proto.PutExecutionTypeResponse
		er, err = g.PutExecutionType(ctx, &proto.PutExecutionTypeRequest{
			ExecutionType:      ex,
			CanAddFields:       request.CanAddFields,
			CanOmitFields:      request.CanOmitFields,
			TransactionOptions: request.TransactionOptions,
		})
		if err != nil {
			return response, err
		}
		response.ExecutionTypeIds = append(response.ExecutionTypeIds, *er.TypeId)
	}
	for _, ct := range request.ContextTypes {
		var cr *proto.PutContextTypeResponse
		cr, err = g.PutContextType(ctx, &proto.PutContextTypeRequest{
			ContextType:        ct,
			CanAddFields:       request.CanAddFields,
			CanOmitFields:      request.CanOmitFields,
			TransactionOptions: request.TransactionOptions,
		})
		if err != nil {
			return response, err
		}
		response.ContextTypeIds = append(response.ContextTypeIds, *cr.TypeId)
	}
	return response, nil
}

var REQUIRED_ARTIFACT_FIELDS = []string{"type_id", "uri"}

func (g grpcServer) PutArtifacts(ctx context.Context, request *proto.PutArtifactsRequest) (resp *proto.PutArtifactsResponse, err error) {
	ctx, dbConn := Begin(ctx, g.dbConnection)
	defer handleTransaction(ctx, &err)

	var artifactIds []int64
	for _, artifact := range request.Artifacts {
		err = requiredFields(REQUIRED_ARTIFACT_FIELDS, artifact.TypeId, artifact.Uri)
		if err != nil {
			return nil, err
		}
		value := &db.Artifact{
			TypeID:     *artifact.TypeId,
			URI:        artifact.Uri,
			Name:       artifact.Name,
			ExternalID: artifact.ExternalId,
		}
		nilSafeCopy(&value.ID, artifact.Id, identity[int64])
		nilSafeCopy(&value.State, artifact.State, artifactStateToInt64)
		// create in DB
		if err = dbConn.Create(value).Error; err != nil {
			err = fmt.Errorf("error creating artifact with type_id[%d], name[%s]: %w", value.TypeID, *value.Name, err)
			return nil, err
		}
		// create properties in DB
		err = g.createArtifactProperties(ctx, value.ID, artifact.GetProperties(), false)
		if err != nil {
			return nil, err
		}
		err = g.createArtifactProperties(ctx, value.ID, artifact.GetCustomProperties(), true)
		if err != nil {
			return nil, err
		}
		artifactIds = append(artifactIds, int64(value.ID))
	}
	resp = &proto.PutArtifactsResponse{
		ArtifactIds: artifactIds,
	}
	return resp, nil
}

func (g grpcServer) PutExecutions(ctx context.Context, request *proto.PutExecutionsRequest) (*proto.PutExecutionsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) PutEvents(ctx context.Context, request *proto.PutEventsRequest) (*proto.PutEventsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) PutExecution(ctx context.Context, request *proto.PutExecutionRequest) (*proto.PutExecutionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) PutLineageSubgraph(ctx context.Context, request *proto.PutLineageSubgraphRequest) (*proto.PutLineageSubgraphResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) PutContexts(ctx context.Context, request *proto.PutContextsRequest) (*proto.PutContextsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) PutAttributionsAndAssociations(ctx context.Context, request *proto.PutAttributionsAndAssociationsRequest) (*proto.PutAttributionsAndAssociationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) PutParentContexts(ctx context.Context, request *proto.PutParentContextsRequest) (*proto.PutParentContextsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetArtifactType(ctx context.Context, request *proto.GetArtifactTypeRequest) (resp *proto.GetArtifactTypeResponse, err error) {
	ctx, dbConn := Begin(ctx, g.dbConnection)
	defer handleTransaction(ctx, &err)

	err = requiredFields(REQUIRED_TYPE_FIELDS, request.TypeName)
	response := &proto.GetArtifactTypeResponse{}

	var results []db.Type
	rx := dbConn.Find(&results, db.Type{Name: *request.TypeName, TypeKind: int32(ARTIFACT_TYPE), Version: request.TypeVersion})
	if rx.Error != nil {
		return nil, rx.Error
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("more than one type found: %v", len(results))
	}
	if len(results) == 0 {
		return response, nil
	}

	r0 := results[0]
	artifactType := proto.ArtifactType{
		Id:          &r0.ID,
		Name:        &r0.Name,
		Version:     r0.Version,
		Description: r0.Description,
		ExternalId:  r0.ExternalID,
	}
	for _, v := range r0.Properties {
		artifactType.Properties[v.Name] = proto.PropertyType(v.DataType)
	}
	response.ArtifactType = &artifactType
	return response, nil
}

func (g grpcServer) GetArtifactTypesByID(ctx context.Context, request *proto.GetArtifactTypesByIDRequest) (*proto.GetArtifactTypesByIDResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetArtifactTypes(ctx context.Context, request *proto.GetArtifactTypesRequest) (*proto.GetArtifactTypesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetExecutionType(ctx context.Context, request *proto.GetExecutionTypeRequest) (*proto.GetExecutionTypeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetExecutionTypesByID(ctx context.Context, request *proto.GetExecutionTypesByIDRequest) (*proto.GetExecutionTypesByIDResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetExecutionTypes(ctx context.Context, request *proto.GetExecutionTypesRequest) (*proto.GetExecutionTypesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContextType(ctx context.Context, request *proto.GetContextTypeRequest) (*proto.GetContextTypeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContextTypesByID(ctx context.Context, request *proto.GetContextTypesByIDRequest) (*proto.GetContextTypesByIDResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContextTypes(ctx context.Context, request *proto.GetContextTypesRequest) (*proto.GetContextTypesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetArtifacts(ctx context.Context, request *proto.GetArtifactsRequest) (*proto.GetArtifactsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetExecutions(ctx context.Context, request *proto.GetExecutionsRequest) (*proto.GetExecutionsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContexts(ctx context.Context, request *proto.GetContextsRequest) (*proto.GetContextsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetArtifactsByID(ctx context.Context, request *proto.GetArtifactsByIDRequest) (*proto.GetArtifactsByIDResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetExecutionsByID(ctx context.Context, request *proto.GetExecutionsByIDRequest) (*proto.GetExecutionsByIDResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContextsByID(ctx context.Context, request *proto.GetContextsByIDRequest) (*proto.GetContextsByIDResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetArtifactsByType(ctx context.Context, request *proto.GetArtifactsByTypeRequest) (*proto.GetArtifactsByTypeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetExecutionsByType(ctx context.Context, request *proto.GetExecutionsByTypeRequest) (*proto.GetExecutionsByTypeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContextsByType(ctx context.Context, request *proto.GetContextsByTypeRequest) (*proto.GetContextsByTypeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetArtifactByTypeAndName(ctx context.Context, request *proto.GetArtifactByTypeAndNameRequest) (*proto.GetArtifactByTypeAndNameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetExecutionByTypeAndName(ctx context.Context, request *proto.GetExecutionByTypeAndNameRequest) (*proto.GetExecutionByTypeAndNameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContextByTypeAndName(ctx context.Context, request *proto.GetContextByTypeAndNameRequest) (*proto.GetContextByTypeAndNameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetArtifactsByURI(ctx context.Context, request *proto.GetArtifactsByURIRequest) (*proto.GetArtifactsByURIResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetEventsByExecutionIDs(ctx context.Context, request *proto.GetEventsByExecutionIDsRequest) (*proto.GetEventsByExecutionIDsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetEventsByArtifactIDs(ctx context.Context, request *proto.GetEventsByArtifactIDsRequest) (*proto.GetEventsByArtifactIDsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetArtifactsByExternalIds(ctx context.Context, request *proto.GetArtifactsByExternalIdsRequest) (*proto.GetArtifactsByExternalIdsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetExecutionsByExternalIds(ctx context.Context, request *proto.GetExecutionsByExternalIdsRequest) (*proto.GetExecutionsByExternalIdsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContextsByExternalIds(ctx context.Context, request *proto.GetContextsByExternalIdsRequest) (*proto.GetContextsByExternalIdsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetArtifactTypesByExternalIds(ctx context.Context, request *proto.GetArtifactTypesByExternalIdsRequest) (*proto.GetArtifactTypesByExternalIdsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetExecutionTypesByExternalIds(ctx context.Context, request *proto.GetExecutionTypesByExternalIdsRequest) (*proto.GetExecutionTypesByExternalIdsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContextTypesByExternalIds(ctx context.Context, request *proto.GetContextTypesByExternalIdsRequest) (*proto.GetContextTypesByExternalIdsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContextsByArtifact(ctx context.Context, request *proto.GetContextsByArtifactRequest) (*proto.GetContextsByArtifactResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetContextsByExecution(ctx context.Context, request *proto.GetContextsByExecutionRequest) (*proto.GetContextsByExecutionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetParentContextsByContext(ctx context.Context, request *proto.GetParentContextsByContextRequest) (*proto.GetParentContextsByContextResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetChildrenContextsByContext(ctx context.Context, request *proto.GetChildrenContextsByContextRequest) (*proto.GetChildrenContextsByContextResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetParentContextsByContexts(ctx context.Context, request *proto.GetParentContextsByContextsRequest) (*proto.GetParentContextsByContextsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetChildrenContextsByContexts(ctx context.Context, request *proto.GetChildrenContextsByContextsRequest) (*proto.GetChildrenContextsByContextsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetArtifactsByContext(ctx context.Context, request *proto.GetArtifactsByContextRequest) (*proto.GetArtifactsByContextResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetExecutionsByContext(ctx context.Context, request *proto.GetExecutionsByContextRequest) (*proto.GetExecutionsByContextResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetLineageGraph(ctx context.Context, request *proto.GetLineageGraphRequest) (*proto.GetLineageGraphResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g grpcServer) GetLineageSubgraph(ctx context.Context, request *proto.GetLineageSubgraphRequest) (*proto.GetLineageSubgraphResponse, error) {
	//TODO implement me
	panic("implement me")
}

//nolint:golint,unused
func (g grpcServer) mustEmbedUnimplementedMetadataStoreServiceServer() {
	// implemented to signal that server is extendable
}

func (g grpcServer) createTypeProperties(ctx context.Context, properties map[string]proto.PropertyType, typeId int64) (err error) {
	ctx, dbConn := Begin(ctx, g.dbConnection)
	defer handleTransaction(ctx, &err)

	for propName, prop := range properties {
		number := int32(prop.Number())
		property := &db.TypeProperty{
			TypeID:   typeId,
			Name:     propName,
			DataType: number,
		}
		if err = dbConn.Where("type_id = ? AND name = ?", typeId, propName).
			Assign(property).FirstOrCreate(property).Error; err != nil {
			err = fmt.Errorf("error creating type property for type_id[%d] with name[%s]: %v", typeId, propName, err)
			return err
		}
	}
	return nil
}

func (g grpcServer) createArtifactProperties(ctx context.Context, artifactId int64, properties map[string]*proto.Value, isCustomProperty bool) (err error) {
	ctx, dbConn := Begin(ctx, g.dbConnection)
	defer handleTransaction(ctx, &err)

	for propName, prop := range properties {
		property := &db.ArtifactProperty{
			ArtifactID: artifactId,
			Name:       propName,
		}
		if isCustomProperty {
			property.IsCustomProperty = true
		}
		// TODO handle polymorphic value with null columns
		intValue, ok := prop.GetValue().(*proto.Value_IntValue)
		if ok {
			property.IntValue = &intValue.IntValue
		}
		doubleValue, ok := prop.GetValue().(*proto.Value_DoubleValue)
		if ok {
			property.DoubleValue = &doubleValue.DoubleValue
		}
		stringValue, ok := prop.GetValue().(*proto.Value_StringValue)
		if ok {
			property.StringValue = &stringValue.StringValue
		}
		structValue, ok := prop.GetValue().(*proto.Value_StructValue)
		if ok {
			json, err2 := structValue.StructValue.MarshalJSON()
			if err2 != nil {
				err = fmt.Errorf("error marshaling struct %s value: %w", propName, err2)
				return err
			}
			property.ByteValue = &json
		}
		protoValue, ok := prop.GetValue().(*proto.Value_ProtoValue)
		if ok {
			property.ProtoValue = &protoValue.ProtoValue.Value
		}
		boolValue, ok := prop.GetValue().(*proto.Value_BoolValue)
		if ok {
			property.BoolValue = &boolValue.BoolValue
		}
		if err = dbConn.Create(property).Error; err != nil {
			err = fmt.Errorf("error creating artifact property for type_id[%d] with name %s: %v", artifactId, propName, err)
			return err
		}
	}
	return nil
}

func identity[T int64 | string](i T) T { return i }
func artifactStateToInt64(i proto.Artifact_State) *int64 {
	var result = int64(i)
	return &result
}

func requiredFields(names []string, args ...interface{}) error {
	var missing []string
	for i, a := range args {
		if a == nil {
			missing = append(missing, names[i])
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required properties: %s", missing)
	}
	return nil
}

func nilSafeCopy[D int32 | int64 | *int64 | string, S int64 | proto.Artifact_State | string](dest *D, src *S, f func(i S) D) {
	if src != nil {
		*dest = f(*src)
	}
}
func handleTransaction(ctx context.Context, err *error) {
	// handle panic
	if perr := recover(); perr != nil {
		_ = Rollback(ctx)
		*err = status.Errorf(codes.Internal, "server panic: %v", perr)
		return
	}
	if err == nil || *err == nil {
		*err = Commit(ctx)
	} else {
		_ = Rollback(ctx)
		if _, ok := status.FromError(*err); !ok {
			*err = status.Errorf(codes.Internal, "internal error: %v", *err)
		}
	}
}
