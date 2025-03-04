package core

import (
	"context"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// SERVING ENVIRONMENT

func (suite *CoreTestSuite) TestCreateServingEnvironment() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:        &entityName,
		ExternalId:  &entityExternalId,
		Description: &entityDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// test
	createdEntity, err := service.UpsertServingEnvironment(eut)

	// checks
	suite.Nilf(err, "error creating uut: %v", err)
	suite.NotNilf(createdEntity.Id, "created uut should not have nil Id")

	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)
	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdEntityId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	suite.Equal(*createdEntity.Id, *ctxId, "returned id should match the mlmd one")
	suite.Equal(entityName, *ctx.Name, "saved name should match the provided one")
	suite.Equal(entityExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(entityDescription, ctx.Properties["description"].GetStringValue(), "saved description should match the provided one")
	suite.Equal(myCustomProp, ctx.CustomProperties["myCustomProp"].GetStringValue(), "saved myCustomProp custom property should match the provided one")

	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")
}

func (suite *CoreTestSuite) TestCreateDuplicateServingEnvironmentFailure() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:        &entityName,
		ExternalId:  &entityExternalId,
		Description: &entityDescription,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// create first serving environment
	createdEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating uut: %v", err)
	suite.NotNilf(createdEntity.Id, "created uut should not have nil Id")

	// attempt to create dupliate serving environment
	_, err = service.UpsertServingEnvironment(eut)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "cannot register a duplicate serving environment")
	suite.Equal(409, statusResp, "duplicate serving environments not allowed")
}

func (suite *CoreTestSuite) TestUpdateServingEnvironment() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalId: &entityExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// test
	createdEntity, err := service.UpsertServingEnvironment(eut)

	// checks
	suite.Nilf(err, "error creating uut: %v", err)
	suite.NotNilf(createdEntity.Id, "created uut should not have nil Id")
	createdEntityId, _ := converter.StringToInt64(createdEntity.Id)

	// checks created entity matches original one except for Id
	suite.Equal(*eut.Name, *createdEntity.Name, "returned entity should match the original one")
	suite.Equal(*eut.ExternalId, *createdEntity.ExternalId, "returned entity external id should match the original one")
	suite.Equal(*eut.CustomProperties, *createdEntity.CustomProperties, "returned entity custom props should match the original one")

	// update existing entity
	newExternalId := "newExternalId"
	newCustomProp := "newCustomProp"

	createdEntity.ExternalId = &newExternalId
	(*createdEntity.CustomProperties)["myCustomProp"] = openapi.MetadataValue{
		MetadataStringValue: converter.NewMetadataStringValue(newCustomProp),
	}

	// update the entity
	createdEntity, err = service.UpsertServingEnvironment(createdEntity)
	suite.Nilf(err, "error creating uut: %v", err)

	// still one expected MLMD type
	getAllResp, err := suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err := suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdEntityId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx := ctxById.Contexts[0]
	ctxId := converter.Int64ToString(ctx.Id)
	suite.Equal(*createdEntity.Id, *ctxId, "returned entity id should match the mlmd one")
	suite.Equal(entityName, *ctx.Name, "saved entity name should match the provided one")
	suite.Equal(newExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(newCustomProp, ctx.CustomProperties["myCustomProp"].GetStringValue(), "saved myCustomProp custom property should match the provided one")

	// update the entity under test, keeping nil name
	newExternalId = "newNewExternalId"
	createdEntity.ExternalId = &newExternalId
	createdEntity.Name = nil
	createdEntity, err = service.UpsertServingEnvironment(createdEntity)
	suite.Nilf(err, "error creating entity: %v", err)

	// still one registered entity
	getAllResp, err = suite.mlmdClient.GetContexts(context.Background(), &proto.GetContextsRequest{})
	suite.Nilf(err, "error retrieving all contexts, not related to the test itself: %v", err)
	suite.Equal(1, len(getAllResp.Contexts), "there should be just one context saved in mlmd")

	ctxById, err = suite.mlmdClient.GetContextsByID(context.Background(), &proto.GetContextsByIDRequest{
		ContextIds: []int64{*createdEntityId},
	})
	suite.Nilf(err, "error retrieving context by type and name, not related to the test itself: %v", err)

	ctx = ctxById.Contexts[0]
	ctxId = converter.Int64ToString(ctx.Id)
	suite.Equal(*createdEntity.Id, *ctxId, "returned entity id should match the mlmd one")
	suite.Equal(entityName, *ctx.Name, "saved entity name should match the provided one")
	suite.Equal(newExternalId, *ctx.ExternalId, "saved external id should match the provided one")
	suite.Equal(newCustomProp, ctx.CustomProperties["myCustomProp"].GetStringValue(), "saved myCustomProp custom property should match the provided one")
}

func (suite *CoreTestSuite) TestGetServingEnvironmentById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new entity
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalId: &entityExternalId,
		CustomProperties: &map[string]openapi.MetadataValue{
			"myCustomProp": {
				MetadataStringValue: converter.NewMetadataStringValue(myCustomProp),
			},
		},
	}

	// test
	createdEntity, err := service.UpsertServingEnvironment(eut)

	// checks
	suite.Nilf(err, "error creating eut: %v", err)

	getEntityById, err := service.GetServingEnvironmentById(*createdEntity.Id)
	suite.Nilf(err, "error getting eut by id %s: %v", *createdEntity.Id, err)

	// checks created entity matches original one except for Id
	suite.Equal(*eut.Name, *getEntityById.Name, "saved name should match the original one")
	suite.Equal(*eut.ExternalId, *getEntityById.ExternalId, "saved external id should match the original one")
	suite.Equal(*eut.CustomProperties, *getEntityById.CustomProperties, "saved custom props should match the original one")
}

func (suite *CoreTestSuite) TestGetServingEnvironmentByParamsWithNoResults() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	_, err := service.GetServingEnvironmentByParams(apiutils.Of("not-present"), nil)
	suite.NotNil(err)
	suite.Equal("no serving environments found for name=not-present, externalId=: not found", err.Error())
}

func (suite *CoreTestSuite) TestGetServingEnvironmentByParamsName() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalId: &entityExternalId,
	}

	createdEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	byName, err := service.GetServingEnvironmentByParams(&entityName, nil)
	suite.Nilf(err, "error getting ServingEnvironment by name: %v", err)

	suite.Equalf(*createdEntity.Id, *byName.Id, "the returned entity id should match the retrieved by name")
}

func (suite *CoreTestSuite) TestGetServingEnvironmentByParamsInvalid() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalId: &entityExternalId,
	}

	// must register a serving environment first, otherwise the http error will be a 404
	_, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	invalidName := "\xFF"

	_, err = service.GetServingEnvironmentByParams(&invalidName, nil)
	statusResp := api.ErrToStatus(err)
	suite.NotNilf(err, "invalid parameter used to retreive serving environemnt")
	suite.Equal(400, statusResp, "invalid parameter used to retreive serving environemnt")
}

func (suite *CoreTestSuite) TestGetServingEnvironmentByParamsExternalId() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalId: &entityExternalId,
	}

	createdEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	byName, err := service.GetServingEnvironmentByParams(nil, &entityExternalId)
	suite.Nilf(err, "error getting ServingEnvironment by external id: %v", err)

	suite.Equalf(*createdEntity.Id, *byName.Id, "the returned entity id should match the retrieved by name")
}

func (suite *CoreTestSuite) TestGetServingEnvironmentByEmptyParams() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalId: &entityExternalId,
	}

	_, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	_, err = service.GetServingEnvironmentByParams(nil, nil)
	suite.NotNil(err)
	suite.Equal("invalid parameters call, supply either name or externalId: bad request", err.Error())
}

func (suite *CoreTestSuite) TestGetServingEnvironmentsOrderedById() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	orderBy := "ID"

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalId: &entityExternalId,
	}

	_, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName := "Pricingentity2"
	newExternalId := "myExternalId2"
	eut.Name = &newName
	eut.ExternalId = &newExternalId
	_, err = service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName = "Pricingentity3"
	newExternalId = "myExternalId3"
	eut.Name = &newName
	eut.ExternalId = &newExternalId
	_, err = service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	orderedById, err := service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	suite.Nilf(err, "error getting ServingEnvironment: %v", err)

	suite.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		suite.Less(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}

	orderedById, err = service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	suite.Nilf(err, "error getting ServingEnvironments: %v", err)

	suite.Equal(3, int(orderedById.Size))
	for i := 0; i < int(orderedById.Size)-1; i++ {
		suite.Greater(*orderedById.Items[i].Id, *orderedById.Items[i+1].Id)
	}
}

func (suite *CoreTestSuite) TestGetServingEnvironmentsOrderedByLastUpdate() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	orderBy := "LAST_UPDATE_TIME"

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalId: &entityExternalId,
	}

	firstEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName := "Pricingentity2"
	newExternalId := "myExternalId2"
	eut.Name = &newName
	eut.ExternalId = &newExternalId
	secondEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName = "Pricingentity3"
	newExternalId = "myExternalId3"
	eut.Name = &newName
	eut.ExternalId = &newExternalId
	thirdEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	// update second entity
	secondEntity.ExternalId = nil
	_, err = service.UpsertServingEnvironment(secondEntity)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	orderedById, err := service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &ascOrderDirection,
	})
	suite.Nilf(err, "error getting ServingEnvironments: %v", err)

	suite.Equal(3, int(orderedById.Size))
	suite.Equal(*firstEntity.Id, *orderedById.Items[0].Id)
	suite.Equal(*thirdEntity.Id, *orderedById.Items[1].Id)
	suite.Equal(*secondEntity.Id, *orderedById.Items[2].Id)

	orderedById, err = service.GetServingEnvironments(api.ListOptions{
		OrderBy:   &orderBy,
		SortOrder: &descOrderDirection,
	})
	suite.Nilf(err, "error getting ServingEnvironments: %v", err)

	suite.Equal(3, int(orderedById.Size))
	suite.Equal(*secondEntity.Id, *orderedById.Items[0].Id)
	suite.Equal(*thirdEntity.Id, *orderedById.Items[1].Id)
	suite.Equal(*firstEntity.Id, *orderedById.Items[2].Id)
}

func (suite *CoreTestSuite) TestGetServingEnvironmentsWithPageSize() {
	// create mode registry service
	service := suite.setupModelRegistryService()

	pageSize := int32(1)
	pageSize2 := int32(2)
	entityName := "Pricingentity1"
	entityExternalId := "myExternalId1"

	// register a new ServingEnvironment
	eut := &openapi.ServingEnvironment{
		Name:       &entityName,
		ExternalId: &entityExternalId,
	}

	firstEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating registered entity: %v", err)

	newName := "Pricingentity2"
	newExternalId := "myExternalId2"
	eut.Name = &newName
	eut.ExternalId = &newExternalId
	secondEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	newName = "Pricingentity3"
	newExternalId = "myExternalId3"
	eut.Name = &newName
	eut.ExternalId = &newExternalId
	thirdEntity, err := service.UpsertServingEnvironment(eut)
	suite.Nilf(err, "error creating ServingEnvironment: %v", err)

	truncatedList, err := service.GetServingEnvironments(api.ListOptions{
		PageSize: &pageSize,
	})
	suite.Nilf(err, "error getting ServingEnvironments: %v", err)

	suite.Equal(1, int(truncatedList.Size))
	suite.NotEqual("", truncatedList.NextPageToken, "next page token should not be empty")
	suite.Equal(*firstEntity.Id, *truncatedList.Items[0].Id)

	truncatedList, err = service.GetServingEnvironments(api.ListOptions{
		PageSize:      &pageSize2,
		NextPageToken: &truncatedList.NextPageToken,
	})
	suite.Nilf(err, "error getting ServingEnvironments: %v", err)

	suite.Equal(2, int(truncatedList.Size))
	suite.Equal("", truncatedList.NextPageToken, "next page token should be empty as list item returned")
	suite.Equal(*secondEntity.Id, *truncatedList.Items[0].Id)
	suite.Equal(*thirdEntity.Id, *truncatedList.Items[1].Id)
}
