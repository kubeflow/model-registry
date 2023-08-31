package server

import (
	"context"
	"errors"
	"github.com/dhirajsb/ml-metadata-go-server/ml_metadata/proto"
	"github.com/dhirajsb/ml-metadata-go-server/model/db"
	"github.com/golang/glog"
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

func (g grpcServer) PutArtifactType(ctx context.Context, request *proto.PutArtifactTypeRequest) (resp *proto.PutArtifactTypeResponse, err error) {
	ctx, dbConn := Begin(ctx, g.dbConnection)
	defer closeDbConnection(ctx, &err)

	artifactType := request.GetArtifactType()
	name := artifactType.Name
	if name == nil {
		return nil, errors.New("missing required field name")
	}
	value := &db.Type{
		Name:     *name,
		TypeKind: int32(ARTIFACT_TYPE),
	}
	if artifactType.Version != nil {
		value.Version = *artifactType.Version
	}
	if artifactType.Description != nil {
		value.Description = *artifactType.Description
	}
	if artifactType.ExternalId != nil {
		value.ExternalID = *artifactType.ExternalId
	}
	if err := dbConn.Create(value).Error; err != nil {
		glog.Errorf("error creating artifact type %s: %v", name, err)
		return nil, err
	}
	err = g.createProperties(ctx, request.ArtifactType.Properties, value)
	if err != nil {
		return nil, err
	}
	var typeId = int64(value.ID)
	return &proto.PutArtifactTypeResponse{
		TypeId: &typeId,
	}, nil
}

func (g grpcServer) createProperties(ctx context.Context, properties map[string]proto.PropertyType, value *db.Type) error {
	for propName, prop := range properties {
		number := int32(prop.Number())
		property := db.TypeProperty{
			TypeID:   value.ID,
			Name:     propName,
			DataType: &number,
		}
		dbConn, _ := FromContext(ctx)
		if err := dbConn.Create(property).Error; err != nil {
			glog.Errorf("error creating type property %s: %v", propName, err)
			return err
		}
	}
	return nil
}

func (g grpcServer) PutExecutionType(ctx context.Context, request *proto.PutExecutionTypeRequest) (resp *proto.PutExecutionTypeResponse, err error) {
	ctx, dbConn := Begin(ctx, g.dbConnection)
	defer closeDbConnection(ctx, &err)

	executionType := request.GetExecutionType()
	value := &db.Type{
		Name:        *executionType.Name,
		Version:     *executionType.Version,
		TypeKind:    int32(EXECUTION_TYPE),
		Description: *(executionType.Description),
		InputType:   executionType.InputType.String(),
		OutputType:  executionType.OutputType.String(),
		ExternalID:  *(executionType.ExternalId),
	}
	if err = dbConn.Create(value).Error; err != nil {
		glog.Errorf("error creating execution type %s: %v", executionType.Name, err)
		return nil, err
	}
	err = g.createProperties(ctx, request.ExecutionType.Properties, value)
	if err != nil {
		return nil, err
	}
	var typeId = int64(value.ID)
	return &proto.PutExecutionTypeResponse{
		TypeId: &typeId,
	}, nil
}

func (g grpcServer) PutContextType(ctx context.Context, request *proto.PutContextTypeRequest) (resp *proto.PutContextTypeResponse, err error) {
	ctx, dbConn := Begin(ctx, g.dbConnection)
	defer closeDbConnection(ctx, &err)

	contextType := request.ContextType
	value := &db.Type{
		Name:        *contextType.Name,
		Version:     *contextType.Version,
		TypeKind:    int32(CONTEXT_TYPE),
		Description: *(contextType.Description),
		ExternalID:  *(contextType.ExternalId),
	}
	if err := dbConn.Create(value).Error; err != nil {
		glog.Errorf("error creating type %s: %v", contextType.Name, err)
		return nil, err
	}
	err = g.createProperties(ctx, request.ContextType.Properties, value)
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
	defer closeDbConnection(ctx, &err)

	response := &proto.PutTypesResponse{}

	for _, ar := range request.ArtifactTypes {
		var at *proto.PutArtifactTypeResponse
		at, err = g.PutArtifactType(ctx, &proto.PutArtifactTypeRequest{
			ArtifactType:       ar,
			CanAddFields:       request.CanAddFields,
			CanOmitFields:      request.CanOmitFields,
			CanDeleteFields:    request.CanDeleteFields,
			AllFieldsMatch:     request.AllFieldsMatch,
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
			CanDeleteFields:    request.CanDeleteFields,
			AllFieldsMatch:     request.AllFieldsMatch,
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
			CanDeleteFields:    request.CanDeleteFields,
			AllFieldsMatch:     request.AllFieldsMatch,
			TransactionOptions: request.TransactionOptions,
		})
		if err != nil {
			return response, err
		}
		response.ContextTypeIds = append(response.ContextTypeIds, *cr.TypeId)
	}
	return response, nil
}

func closeDbConnection(ctx context.Context, err *error) {
	if err == nil || *err == nil {
		*err = Commit(ctx)
	} else {
		_ = Rollback(ctx)
	}
}

func (g grpcServer) PutArtifacts(ctx context.Context, request *proto.PutArtifactsRequest) (*proto.PutArtifactsResponse, error) {
	//TODO implement me
	panic("implement me")
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

func (g grpcServer) GetArtifactType(ctx context.Context, request *proto.GetArtifactTypeRequest) (*proto.GetArtifactTypeResponse, error) {
	//TODO implement me
	panic("implement me")
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

func (g grpcServer) mustEmbedUnimplementedMetadataStoreServiceServer() {
	//TODO implement me
	panic("implement me")
}
