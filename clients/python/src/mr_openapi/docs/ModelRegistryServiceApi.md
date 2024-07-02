# mr_openapi.ModelRegistryServiceApi

All URIs are relative to *https://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**create_environment_inference_service**](ModelRegistryServiceApi.md#create_environment_inference_service) | **POST** /api/model_registry/v1alpha3/serving_environments/{servingenvironmentId}/inference_services | Create a InferenceService in ServingEnvironment
[**create_inference_service**](ModelRegistryServiceApi.md#create_inference_service) | **POST** /api/model_registry/v1alpha3/inference_services | Create a InferenceService
[**create_inference_service_serve**](ModelRegistryServiceApi.md#create_inference_service_serve) | **POST** /api/model_registry/v1alpha3/inference_services/{inferenceserviceId}/serves | Create a ServeModel action in a InferenceService
[**create_model_artifact**](ModelRegistryServiceApi.md#create_model_artifact) | **POST** /api/model_registry/v1alpha3/model_artifacts | Create a ModelArtifact
[**create_model_version**](ModelRegistryServiceApi.md#create_model_version) | **POST** /api/model_registry/v1alpha3/model_versions | Create a ModelVersion
[**create_model_version_artifact**](ModelRegistryServiceApi.md#create_model_version_artifact) | **POST** /api/model_registry/v1alpha3/model_versions/{modelversionId}/artifacts | Create an Artifact in a ModelVersion
[**create_registered_model**](ModelRegistryServiceApi.md#create_registered_model) | **POST** /api/model_registry/v1alpha3/registered_models | Create a RegisteredModel
[**create_registered_model_version**](ModelRegistryServiceApi.md#create_registered_model_version) | **POST** /api/model_registry/v1alpha3/registered_models/{registeredmodelId}/versions | Create a ModelVersion in RegisteredModel
[**create_serving_environment**](ModelRegistryServiceApi.md#create_serving_environment) | **POST** /api/model_registry/v1alpha3/serving_environments | Create a ServingEnvironment
[**find_inference_service**](ModelRegistryServiceApi.md#find_inference_service) | **GET** /api/model_registry/v1alpha3/inference_service | Get an InferenceServices that matches search parameters.
[**find_model_artifact**](ModelRegistryServiceApi.md#find_model_artifact) | **GET** /api/model_registry/v1alpha3/model_artifact | Get a ModelArtifact that matches search parameters.
[**find_model_version**](ModelRegistryServiceApi.md#find_model_version) | **GET** /api/model_registry/v1alpha3/model_version | Get a ModelVersion that matches search parameters.
[**find_registered_model**](ModelRegistryServiceApi.md#find_registered_model) | **GET** /api/model_registry/v1alpha3/registered_model | Get a RegisteredModel that matches search parameters.
[**find_serving_environment**](ModelRegistryServiceApi.md#find_serving_environment) | **GET** /api/model_registry/v1alpha3/serving_environment | Find ServingEnvironment
[**get_environment_inference_services**](ModelRegistryServiceApi.md#get_environment_inference_services) | **GET** /api/model_registry/v1alpha3/serving_environments/{servingenvironmentId}/inference_services | List All ServingEnvironment&#39;s InferenceServices
[**get_inference_service**](ModelRegistryServiceApi.md#get_inference_service) | **GET** /api/model_registry/v1alpha3/inference_services/{inferenceserviceId} | Get a InferenceService
[**get_inference_service_model**](ModelRegistryServiceApi.md#get_inference_service_model) | **GET** /api/model_registry/v1alpha3/inference_services/{inferenceserviceId}/model | Get InferenceService&#39;s RegisteredModel
[**get_inference_service_serves**](ModelRegistryServiceApi.md#get_inference_service_serves) | **GET** /api/model_registry/v1alpha3/inference_services/{inferenceserviceId}/serves | List All InferenceService&#39;s ServeModel actions
[**get_inference_service_version**](ModelRegistryServiceApi.md#get_inference_service_version) | **GET** /api/model_registry/v1alpha3/inference_services/{inferenceserviceId}/version | Get InferenceService&#39;s ModelVersion
[**get_inference_services**](ModelRegistryServiceApi.md#get_inference_services) | **GET** /api/model_registry/v1alpha3/inference_services | List All InferenceServices
[**get_model_artifact**](ModelRegistryServiceApi.md#get_model_artifact) | **GET** /api/model_registry/v1alpha3/model_artifacts/{modelartifactId} | Get a ModelArtifact
[**get_model_artifacts**](ModelRegistryServiceApi.md#get_model_artifacts) | **GET** /api/model_registry/v1alpha3/model_artifacts | List All ModelArtifacts
[**get_model_version**](ModelRegistryServiceApi.md#get_model_version) | **GET** /api/model_registry/v1alpha3/model_versions/{modelversionId} | Get a ModelVersion
[**get_model_version_artifacts**](ModelRegistryServiceApi.md#get_model_version_artifacts) | **GET** /api/model_registry/v1alpha3/model_versions/{modelversionId}/artifacts | List all artifacts associated with the &#x60;ModelVersion&#x60;
[**get_model_versions**](ModelRegistryServiceApi.md#get_model_versions) | **GET** /api/model_registry/v1alpha3/model_versions | List All ModelVersions
[**get_registered_model**](ModelRegistryServiceApi.md#get_registered_model) | **GET** /api/model_registry/v1alpha3/registered_models/{registeredmodelId} | Get a RegisteredModel
[**get_registered_model_versions**](ModelRegistryServiceApi.md#get_registered_model_versions) | **GET** /api/model_registry/v1alpha3/registered_models/{registeredmodelId}/versions | List All RegisteredModel&#39;s ModelVersions
[**get_registered_models**](ModelRegistryServiceApi.md#get_registered_models) | **GET** /api/model_registry/v1alpha3/registered_models | List All RegisteredModels
[**get_serving_environment**](ModelRegistryServiceApi.md#get_serving_environment) | **GET** /api/model_registry/v1alpha3/serving_environments/{servingenvironmentId} | Get a ServingEnvironment
[**get_serving_environments**](ModelRegistryServiceApi.md#get_serving_environments) | **GET** /api/model_registry/v1alpha3/serving_environments | List All ServingEnvironments
[**update_inference_service**](ModelRegistryServiceApi.md#update_inference_service) | **PATCH** /api/model_registry/v1alpha3/inference_services/{inferenceserviceId} | Update a InferenceService
[**update_model_artifact**](ModelRegistryServiceApi.md#update_model_artifact) | **PATCH** /api/model_registry/v1alpha3/model_artifacts/{modelartifactId} | Update a ModelArtifact
[**update_model_version**](ModelRegistryServiceApi.md#update_model_version) | **PATCH** /api/model_registry/v1alpha3/model_versions/{modelversionId} | Update a ModelVersion
[**update_registered_model**](ModelRegistryServiceApi.md#update_registered_model) | **PATCH** /api/model_registry/v1alpha3/registered_models/{registeredmodelId} | Update a RegisteredModel
[**update_serving_environment**](ModelRegistryServiceApi.md#update_serving_environment) | **PATCH** /api/model_registry/v1alpha3/serving_environments/{servingenvironmentId} | Update a ServingEnvironment


# **create_environment_inference_service**
> InferenceService create_environment_inference_service(servingenvironment_id, inference_service_create)

Create a InferenceService in ServingEnvironment

Creates a new instance of a `InferenceService`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.inference_service import InferenceService
from mr_openapi.models.inference_service_create import InferenceServiceCreate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    servingenvironment_id = 'servingenvironment_id_example' # str | A unique identifier for a `ServingEnvironment`.
    inference_service_create = mr_openapi.InferenceServiceCreate() # InferenceServiceCreate | A new `InferenceService` to be created.

    try:
        # Create a InferenceService in ServingEnvironment
        api_response = await api_instance.create_environment_inference_service(servingenvironment_id, inference_service_create)
        print("The response of ModelRegistryServiceApi->create_environment_inference_service:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->create_environment_inference_service: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **servingenvironment_id** | **str**| A unique identifier for a &#x60;ServingEnvironment&#x60;. | 
 **inference_service_create** | [**InferenceServiceCreate**](InferenceServiceCreate.md)| A new &#x60;InferenceService&#x60; to be created. | 

### Return type

[**InferenceService**](InferenceService.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | A response containing a &#x60;InferenceService&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **create_inference_service**
> InferenceService create_inference_service(inference_service_create)

Create a InferenceService

Creates a new instance of a `InferenceService`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.inference_service import InferenceService
from mr_openapi.models.inference_service_create import InferenceServiceCreate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    inference_service_create = mr_openapi.InferenceServiceCreate() # InferenceServiceCreate | A new `InferenceService` to be created.

    try:
        # Create a InferenceService
        api_response = await api_instance.create_inference_service(inference_service_create)
        print("The response of ModelRegistryServiceApi->create_inference_service:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->create_inference_service: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **inference_service_create** | [**InferenceServiceCreate**](InferenceServiceCreate.md)| A new &#x60;InferenceService&#x60; to be created. | 

### Return type

[**InferenceService**](InferenceService.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;InferenceService&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **create_inference_service_serve**
> ServeModel create_inference_service_serve(inferenceservice_id, serve_model_create)

Create a ServeModel action in a InferenceService

Creates a new instance of a `ServeModel` associated with `InferenceService`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.serve_model import ServeModel
from mr_openapi.models.serve_model_create import ServeModelCreate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    inferenceservice_id = 'inferenceservice_id_example' # str | A unique identifier for a `InferenceService`.
    serve_model_create = mr_openapi.ServeModelCreate() # ServeModelCreate | A new `ServeModel` to be associated with the `InferenceService`.

    try:
        # Create a ServeModel action in a InferenceService
        api_response = await api_instance.create_inference_service_serve(inferenceservice_id, serve_model_create)
        print("The response of ModelRegistryServiceApi->create_inference_service_serve:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->create_inference_service_serve: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **inferenceservice_id** | **str**| A unique identifier for a &#x60;InferenceService&#x60;. | 
 **serve_model_create** | [**ServeModelCreate**](ServeModelCreate.md)| A new &#x60;ServeModel&#x60; to be associated with the &#x60;InferenceService&#x60;. | 

### Return type

[**ServeModel**](ServeModel.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | A response containing a &#x60;ServeModel&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **create_model_artifact**
> ModelArtifact create_model_artifact(model_artifact_create)

Create a ModelArtifact

Creates a new instance of a `ModelArtifact`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_artifact import ModelArtifact
from mr_openapi.models.model_artifact_create import ModelArtifactCreate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    model_artifact_create = mr_openapi.ModelArtifactCreate() # ModelArtifactCreate | A new `ModelArtifact` to be created.

    try:
        # Create a ModelArtifact
        api_response = await api_instance.create_model_artifact(model_artifact_create)
        print("The response of ModelRegistryServiceApi->create_model_artifact:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->create_model_artifact: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **model_artifact_create** | [**ModelArtifactCreate**](ModelArtifactCreate.md)| A new &#x60;ModelArtifact&#x60; to be created. | 

### Return type

[**ModelArtifact**](ModelArtifact.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | A response containing a &#x60;ModelArtifact&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **create_model_version**
> ModelVersion create_model_version(model_version_create)

Create a ModelVersion

Creates a new instance of a `ModelVersion`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_version import ModelVersion
from mr_openapi.models.model_version_create import ModelVersionCreate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    model_version_create = mr_openapi.ModelVersionCreate() # ModelVersionCreate | A new `ModelVersion` to be created.

    try:
        # Create a ModelVersion
        api_response = await api_instance.create_model_version(model_version_create)
        print("The response of ModelRegistryServiceApi->create_model_version:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->create_model_version: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **model_version_create** | [**ModelVersionCreate**](ModelVersionCreate.md)| A new &#x60;ModelVersion&#x60; to be created. | 

### Return type

[**ModelVersion**](ModelVersion.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | A response containing a &#x60;ModelVersion&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **create_model_version_artifact**
> Artifact create_model_version_artifact(modelversion_id, artifact)

Create an Artifact in a ModelVersion

Creates a new instance of an Artifact if needed and associates it with `ModelVersion`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.artifact import Artifact
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    modelversion_id = 'modelversion_id_example' # str | A unique identifier for a `ModelVersion`.
    artifact = mr_openapi.Artifact() # Artifact | A new or existing `Artifact` to be associated with the `ModelVersion`.

    try:
        # Create an Artifact in a ModelVersion
        api_response = await api_instance.create_model_version_artifact(modelversion_id, artifact)
        print("The response of ModelRegistryServiceApi->create_model_version_artifact:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->create_model_version_artifact: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **modelversion_id** | **str**| A unique identifier for a &#x60;ModelVersion&#x60;. | 
 **artifact** | [**Artifact**](Artifact.md)| A new or existing &#x60;Artifact&#x60; to be associated with the &#x60;ModelVersion&#x60;. | 

### Return type

[**Artifact**](Artifact.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing an &#x60;Artifact&#x60; entity. |  -  |
**201** | A response containing an &#x60;Artifact&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **create_registered_model**
> RegisteredModel create_registered_model(registered_model_create)

Create a RegisteredModel

Creates a new instance of a `RegisteredModel`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.registered_model import RegisteredModel
from mr_openapi.models.registered_model_create import RegisteredModelCreate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    registered_model_create = mr_openapi.RegisteredModelCreate() # RegisteredModelCreate | A new `RegisteredModel` to be created.

    try:
        # Create a RegisteredModel
        api_response = await api_instance.create_registered_model(registered_model_create)
        print("The response of ModelRegistryServiceApi->create_registered_model:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->create_registered_model: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registered_model_create** | [**RegisteredModelCreate**](RegisteredModelCreate.md)| A new &#x60;RegisteredModel&#x60; to be created. | 

### Return type

[**RegisteredModel**](RegisteredModel.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | A response containing a &#x60;RegisteredModel&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **create_registered_model_version**
> ModelVersion create_registered_model_version(registeredmodel_id, model_version)

Create a ModelVersion in RegisteredModel

Creates a new instance of a `ModelVersion`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_version import ModelVersion
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    registeredmodel_id = 'registeredmodel_id_example' # str | A unique identifier for a `RegisteredModel`.
    model_version = mr_openapi.ModelVersion() # ModelVersion | A new `ModelVersion` to be created.

    try:
        # Create a ModelVersion in RegisteredModel
        api_response = await api_instance.create_registered_model_version(registeredmodel_id, model_version)
        print("The response of ModelRegistryServiceApi->create_registered_model_version:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->create_registered_model_version: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registeredmodel_id** | **str**| A unique identifier for a &#x60;RegisteredModel&#x60;. | 
 **model_version** | [**ModelVersion**](ModelVersion.md)| A new &#x60;ModelVersion&#x60; to be created. | 

### Return type

[**ModelVersion**](ModelVersion.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | A response containing a &#x60;ModelVersion&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **create_serving_environment**
> ServingEnvironment create_serving_environment(serving_environment_create)

Create a ServingEnvironment

Creates a new instance of a `ServingEnvironment`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.serving_environment import ServingEnvironment
from mr_openapi.models.serving_environment_create import ServingEnvironmentCreate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    serving_environment_create = mr_openapi.ServingEnvironmentCreate() # ServingEnvironmentCreate | A new `ServingEnvironment` to be created.

    try:
        # Create a ServingEnvironment
        api_response = await api_instance.create_serving_environment(serving_environment_create)
        print("The response of ModelRegistryServiceApi->create_serving_environment:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->create_serving_environment: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **serving_environment_create** | [**ServingEnvironmentCreate**](ServingEnvironmentCreate.md)| A new &#x60;ServingEnvironment&#x60; to be created. | 

### Return type

[**ServingEnvironment**](ServingEnvironment.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | A response containing a &#x60;ServingEnvironment&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **find_inference_service**
> InferenceService find_inference_service(name=name, external_id=external_id, parent_resource_id=parent_resource_id)

Get an InferenceServices that matches search parameters.

Gets the details of a single instance of `InferenceService` that matches search parameters.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.inference_service import InferenceService
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    name = 'entity-name' # str | Name of entity to search. (optional)
    external_id = '10' # str | External ID of entity to search. (optional)
    parent_resource_id = '10' # str | ID of the parent resource to use for search. (optional)

    try:
        # Get an InferenceServices that matches search parameters.
        api_response = await api_instance.find_inference_service(name=name, external_id=external_id, parent_resource_id=parent_resource_id)
        print("The response of ModelRegistryServiceApi->find_inference_service:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->find_inference_service: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **str**| Name of entity to search. | [optional] 
 **external_id** | **str**| External ID of entity to search. | [optional] 
 **parent_resource_id** | **str**| ID of the parent resource to use for search. | [optional] 

### Return type

[**InferenceService**](InferenceService.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;InferenceService&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **find_model_artifact**
> ModelArtifact find_model_artifact(name=name, external_id=external_id, parent_resource_id=parent_resource_id)

Get a ModelArtifact that matches search parameters.

Gets the details of a single instance of a `ModelArtifact` that matches search parameters.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_artifact import ModelArtifact
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    name = 'entity-name' # str | Name of entity to search. (optional)
    external_id = '10' # str | External ID of entity to search. (optional)
    parent_resource_id = '10' # str | ID of the parent resource to use for search. (optional)

    try:
        # Get a ModelArtifact that matches search parameters.
        api_response = await api_instance.find_model_artifact(name=name, external_id=external_id, parent_resource_id=parent_resource_id)
        print("The response of ModelRegistryServiceApi->find_model_artifact:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->find_model_artifact: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **str**| Name of entity to search. | [optional] 
 **external_id** | **str**| External ID of entity to search. | [optional] 
 **parent_resource_id** | **str**| ID of the parent resource to use for search. | [optional] 

### Return type

[**ModelArtifact**](ModelArtifact.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;ModelArtifact&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **find_model_version**
> ModelVersion find_model_version(name=name, external_id=external_id, parent_resource_id=parent_resource_id)

Get a ModelVersion that matches search parameters.

Gets the details of a single instance of a `ModelVersion` that matches search parameters.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_version import ModelVersion
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    name = 'entity-name' # str | Name of entity to search. (optional)
    external_id = '10' # str | External ID of entity to search. (optional)
    parent_resource_id = '10' # str | ID of the parent resource to use for search. (optional)

    try:
        # Get a ModelVersion that matches search parameters.
        api_response = await api_instance.find_model_version(name=name, external_id=external_id, parent_resource_id=parent_resource_id)
        print("The response of ModelRegistryServiceApi->find_model_version:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->find_model_version: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **str**| Name of entity to search. | [optional] 
 **external_id** | **str**| External ID of entity to search. | [optional] 
 **parent_resource_id** | **str**| ID of the parent resource to use for search. | [optional] 

### Return type

[**ModelVersion**](ModelVersion.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;ModelVersion&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **find_registered_model**
> RegisteredModel find_registered_model(name=name, external_id=external_id)

Get a RegisteredModel that matches search parameters.

Gets the details of a single instance of a `RegisteredModel` that matches search parameters.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.registered_model import RegisteredModel
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    name = 'entity-name' # str | Name of entity to search. (optional)
    external_id = '10' # str | External ID of entity to search. (optional)

    try:
        # Get a RegisteredModel that matches search parameters.
        api_response = await api_instance.find_registered_model(name=name, external_id=external_id)
        print("The response of ModelRegistryServiceApi->find_registered_model:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->find_registered_model: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **str**| Name of entity to search. | [optional] 
 **external_id** | **str**| External ID of entity to search. | [optional] 

### Return type

[**RegisteredModel**](RegisteredModel.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;RegisteredModel&#x60; entity. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **find_serving_environment**
> ServingEnvironment find_serving_environment(name=name, external_id=external_id)

Find ServingEnvironment

Finds a `ServingEnvironment` entity that matches query parameters.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.serving_environment import ServingEnvironment
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    name = 'entity-name' # str | Name of entity to search. (optional)
    external_id = '10' # str | External ID of entity to search. (optional)

    try:
        # Find ServingEnvironment
        api_response = await api_instance.find_serving_environment(name=name, external_id=external_id)
        print("The response of ModelRegistryServiceApi->find_serving_environment:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->find_serving_environment: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **str**| Name of entity to search. | [optional] 
 **external_id** | **str**| External ID of entity to search. | [optional] 

### Return type

[**ServingEnvironment**](ServingEnvironment.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;ServingEnvironment&#x60; entity. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_environment_inference_services**
> InferenceServiceList get_environment_inference_services(servingenvironment_id, name=name, external_id=external_id, page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)

List All ServingEnvironment's InferenceServices

Gets a list of all `InferenceService` entities for the `ServingEnvironment`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.inference_service_list import InferenceServiceList
from mr_openapi.models.order_by_field import OrderByField
from mr_openapi.models.sort_order import SortOrder
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    servingenvironment_id = 'servingenvironment_id_example' # str | A unique identifier for a `ServingEnvironment`.
    name = 'entity-name' # str | Name of entity to search. (optional)
    external_id = '10' # str | External ID of entity to search. (optional)
    page_size = '100' # str | Number of entities in each page. (optional)
    order_by = mr_openapi.OrderByField() # OrderByField | Specifies the order by criteria for listing entities. (optional)
    sort_order = mr_openapi.SortOrder() # SortOrder | Specifies the sort order for listing entities, defaults to ASC. (optional)
    next_page_token = 'IkhlbGxvLCB3b3JsZC4i' # str | Token to use to retrieve next page of results. (optional)

    try:
        # List All ServingEnvironment's InferenceServices
        api_response = await api_instance.get_environment_inference_services(servingenvironment_id, name=name, external_id=external_id, page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)
        print("The response of ModelRegistryServiceApi->get_environment_inference_services:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_environment_inference_services: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **servingenvironment_id** | **str**| A unique identifier for a &#x60;ServingEnvironment&#x60;. | 
 **name** | **str**| Name of entity to search. | [optional] 
 **external_id** | **str**| External ID of entity to search. | [optional] 
 **page_size** | **str**| Number of entities in each page. | [optional] 
 **order_by** | [**OrderByField**](.md)| Specifies the order by criteria for listing entities. | [optional] 
 **sort_order** | [**SortOrder**](.md)| Specifies the sort order for listing entities, defaults to ASC. | [optional] 
 **next_page_token** | **str**| Token to use to retrieve next page of results. | [optional] 

### Return type

[**InferenceServiceList**](InferenceServiceList.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a list of &#x60;InferenceService&#x60; entities. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_inference_service**
> InferenceService get_inference_service(inferenceservice_id)

Get a InferenceService

Gets the details of a single instance of a `InferenceService`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.inference_service import InferenceService
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    inferenceservice_id = 'inferenceservice_id_example' # str | A unique identifier for a `InferenceService`.

    try:
        # Get a InferenceService
        api_response = await api_instance.get_inference_service(inferenceservice_id)
        print("The response of ModelRegistryServiceApi->get_inference_service:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_inference_service: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **inferenceservice_id** | **str**| A unique identifier for a &#x60;InferenceService&#x60;. | 

### Return type

[**InferenceService**](InferenceService.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;InferenceService&#x60; entity. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_inference_service_model**
> RegisteredModel get_inference_service_model(inferenceservice_id)

Get InferenceService's RegisteredModel

Gets the `RegisteredModel` entity for the `InferenceService`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.registered_model import RegisteredModel
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    inferenceservice_id = 'inferenceservice_id_example' # str | A unique identifier for a `InferenceService`.

    try:
        # Get InferenceService's RegisteredModel
        api_response = await api_instance.get_inference_service_model(inferenceservice_id)
        print("The response of ModelRegistryServiceApi->get_inference_service_model:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_inference_service_model: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **inferenceservice_id** | **str**| A unique identifier for a &#x60;InferenceService&#x60;. | 

### Return type

[**RegisteredModel**](RegisteredModel.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;RegisteredModel&#x60; entity. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_inference_service_serves**
> ServeModelList get_inference_service_serves(inferenceservice_id, name=name, external_id=external_id, page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)

List All InferenceService's ServeModel actions

Gets a list of all `ServeModel` entities for the `InferenceService`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.order_by_field import OrderByField
from mr_openapi.models.serve_model_list import ServeModelList
from mr_openapi.models.sort_order import SortOrder
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    inferenceservice_id = 'inferenceservice_id_example' # str | A unique identifier for a `InferenceService`.
    name = 'entity-name' # str | Name of entity to search. (optional)
    external_id = '10' # str | External ID of entity to search. (optional)
    page_size = '100' # str | Number of entities in each page. (optional)
    order_by = mr_openapi.OrderByField() # OrderByField | Specifies the order by criteria for listing entities. (optional)
    sort_order = mr_openapi.SortOrder() # SortOrder | Specifies the sort order for listing entities, defaults to ASC. (optional)
    next_page_token = 'IkhlbGxvLCB3b3JsZC4i' # str | Token to use to retrieve next page of results. (optional)

    try:
        # List All InferenceService's ServeModel actions
        api_response = await api_instance.get_inference_service_serves(inferenceservice_id, name=name, external_id=external_id, page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)
        print("The response of ModelRegistryServiceApi->get_inference_service_serves:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_inference_service_serves: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **inferenceservice_id** | **str**| A unique identifier for a &#x60;InferenceService&#x60;. | 
 **name** | **str**| Name of entity to search. | [optional] 
 **external_id** | **str**| External ID of entity to search. | [optional] 
 **page_size** | **str**| Number of entities in each page. | [optional] 
 **order_by** | [**OrderByField**](.md)| Specifies the order by criteria for listing entities. | [optional] 
 **sort_order** | [**SortOrder**](.md)| Specifies the sort order for listing entities, defaults to ASC. | [optional] 
 **next_page_token** | **str**| Token to use to retrieve next page of results. | [optional] 

### Return type

[**ServeModelList**](ServeModelList.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a list of &#x60;ServeModel&#x60; entities. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_inference_service_version**
> ModelVersion get_inference_service_version(inferenceservice_id)

Get InferenceService's ModelVersion

Gets the `ModelVersion` entity for the `InferenceService`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_version import ModelVersion
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    inferenceservice_id = 'inferenceservice_id_example' # str | A unique identifier for a `InferenceService`.

    try:
        # Get InferenceService's ModelVersion
        api_response = await api_instance.get_inference_service_version(inferenceservice_id)
        print("The response of ModelRegistryServiceApi->get_inference_service_version:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_inference_service_version: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **inferenceservice_id** | **str**| A unique identifier for a &#x60;InferenceService&#x60;. | 

### Return type

[**ModelVersion**](ModelVersion.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;ModelVersion&#x60; entity. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_inference_services**
> InferenceServiceList get_inference_services(page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)

List All InferenceServices

Gets a list of all `InferenceService` entities.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.inference_service_list import InferenceServiceList
from mr_openapi.models.order_by_field import OrderByField
from mr_openapi.models.sort_order import SortOrder
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    page_size = '100' # str | Number of entities in each page. (optional)
    order_by = mr_openapi.OrderByField() # OrderByField | Specifies the order by criteria for listing entities. (optional)
    sort_order = mr_openapi.SortOrder() # SortOrder | Specifies the sort order for listing entities, defaults to ASC. (optional)
    next_page_token = 'IkhlbGxvLCB3b3JsZC4i' # str | Token to use to retrieve next page of results. (optional)

    try:
        # List All InferenceServices
        api_response = await api_instance.get_inference_services(page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)
        print("The response of ModelRegistryServiceApi->get_inference_services:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_inference_services: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page_size** | **str**| Number of entities in each page. | [optional] 
 **order_by** | [**OrderByField**](.md)| Specifies the order by criteria for listing entities. | [optional] 
 **sort_order** | [**SortOrder**](.md)| Specifies the sort order for listing entities, defaults to ASC. | [optional] 
 **next_page_token** | **str**| Token to use to retrieve next page of results. | [optional] 

### Return type

[**InferenceServiceList**](InferenceServiceList.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a list of &#x60;InferenceService&#x60; entities. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_model_artifact**
> ModelArtifact get_model_artifact(modelartifact_id)

Get a ModelArtifact

Gets the details of a single instance of a `ModelArtifact`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_artifact import ModelArtifact
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    modelartifact_id = 'modelartifact_id_example' # str | A unique identifier for a `ModelArtifact`.

    try:
        # Get a ModelArtifact
        api_response = await api_instance.get_model_artifact(modelartifact_id)
        print("The response of ModelRegistryServiceApi->get_model_artifact:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_model_artifact: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **modelartifact_id** | **str**| A unique identifier for a &#x60;ModelArtifact&#x60;. | 

### Return type

[**ModelArtifact**](ModelArtifact.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;ModelArtifact&#x60; entity. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_model_artifacts**
> ModelArtifactList get_model_artifacts(page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)

List All ModelArtifacts

Gets a list of all `ModelArtifact` entities.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_artifact_list import ModelArtifactList
from mr_openapi.models.order_by_field import OrderByField
from mr_openapi.models.sort_order import SortOrder
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    page_size = '100' # str | Number of entities in each page. (optional)
    order_by = mr_openapi.OrderByField() # OrderByField | Specifies the order by criteria for listing entities. (optional)
    sort_order = mr_openapi.SortOrder() # SortOrder | Specifies the sort order for listing entities, defaults to ASC. (optional)
    next_page_token = 'IkhlbGxvLCB3b3JsZC4i' # str | Token to use to retrieve next page of results. (optional)

    try:
        # List All ModelArtifacts
        api_response = await api_instance.get_model_artifacts(page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)
        print("The response of ModelRegistryServiceApi->get_model_artifacts:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_model_artifacts: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page_size** | **str**| Number of entities in each page. | [optional] 
 **order_by** | [**OrderByField**](.md)| Specifies the order by criteria for listing entities. | [optional] 
 **sort_order** | [**SortOrder**](.md)| Specifies the sort order for listing entities, defaults to ASC. | [optional] 
 **next_page_token** | **str**| Token to use to retrieve next page of results. | [optional] 

### Return type

[**ModelArtifactList**](ModelArtifactList.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a list of ModelArtifact entities. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_model_version**
> ModelVersion get_model_version(modelversion_id)

Get a ModelVersion

Gets the details of a single instance of a `ModelVersion`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_version import ModelVersion
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    modelversion_id = 'modelversion_id_example' # str | A unique identifier for a `ModelVersion`.

    try:
        # Get a ModelVersion
        api_response = await api_instance.get_model_version(modelversion_id)
        print("The response of ModelRegistryServiceApi->get_model_version:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_model_version: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **modelversion_id** | **str**| A unique identifier for a &#x60;ModelVersion&#x60;. | 

### Return type

[**ModelVersion**](ModelVersion.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;ModelVersion&#x60; entity. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_model_version_artifacts**
> ArtifactList get_model_version_artifacts(modelversion_id, name=name, external_id=external_id, page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)

List all artifacts associated with the `ModelVersion`

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.artifact_list import ArtifactList
from mr_openapi.models.order_by_field import OrderByField
from mr_openapi.models.sort_order import SortOrder
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    modelversion_id = 'modelversion_id_example' # str | A unique identifier for a `ModelVersion`.
    name = 'entity-name' # str | Name of entity to search. (optional)
    external_id = '10' # str | External ID of entity to search. (optional)
    page_size = '100' # str | Number of entities in each page. (optional)
    order_by = mr_openapi.OrderByField() # OrderByField | Specifies the order by criteria for listing entities. (optional)
    sort_order = mr_openapi.SortOrder() # SortOrder | Specifies the sort order for listing entities, defaults to ASC. (optional)
    next_page_token = 'IkhlbGxvLCB3b3JsZC4i' # str | Token to use to retrieve next page of results. (optional)

    try:
        # List all artifacts associated with the `ModelVersion`
        api_response = await api_instance.get_model_version_artifacts(modelversion_id, name=name, external_id=external_id, page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)
        print("The response of ModelRegistryServiceApi->get_model_version_artifacts:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_model_version_artifacts: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **modelversion_id** | **str**| A unique identifier for a &#x60;ModelVersion&#x60;. | 
 **name** | **str**| Name of entity to search. | [optional] 
 **external_id** | **str**| External ID of entity to search. | [optional] 
 **page_size** | **str**| Number of entities in each page. | [optional] 
 **order_by** | [**OrderByField**](.md)| Specifies the order by criteria for listing entities. | [optional] 
 **sort_order** | [**SortOrder**](.md)| Specifies the sort order for listing entities, defaults to ASC. | [optional] 
 **next_page_token** | **str**| Token to use to retrieve next page of results. | [optional] 

### Return type

[**ArtifactList**](ArtifactList.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a list of &#x60;Artifact&#x60; entities. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_model_versions**
> ModelVersionList get_model_versions(page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)

List All ModelVersions

Gets a list of all `ModelVersion` entities.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_version_list import ModelVersionList
from mr_openapi.models.order_by_field import OrderByField
from mr_openapi.models.sort_order import SortOrder
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    page_size = '100' # str | Number of entities in each page. (optional)
    order_by = mr_openapi.OrderByField() # OrderByField | Specifies the order by criteria for listing entities. (optional)
    sort_order = mr_openapi.SortOrder() # SortOrder | Specifies the sort order for listing entities, defaults to ASC. (optional)
    next_page_token = 'IkhlbGxvLCB3b3JsZC4i' # str | Token to use to retrieve next page of results. (optional)

    try:
        # List All ModelVersions
        api_response = await api_instance.get_model_versions(page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)
        print("The response of ModelRegistryServiceApi->get_model_versions:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_model_versions: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page_size** | **str**| Number of entities in each page. | [optional] 
 **order_by** | [**OrderByField**](.md)| Specifies the order by criteria for listing entities. | [optional] 
 **sort_order** | [**SortOrder**](.md)| Specifies the sort order for listing entities, defaults to ASC. | [optional] 
 **next_page_token** | **str**| Token to use to retrieve next page of results. | [optional] 

### Return type

[**ModelVersionList**](ModelVersionList.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a list of &#x60;ModelVersion&#x60; entities. |  -  |
**401** | Unauthorized |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_registered_model**
> RegisteredModel get_registered_model(registeredmodel_id)

Get a RegisteredModel

Gets the details of a single instance of a `RegisteredModel`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.registered_model import RegisteredModel
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    registeredmodel_id = 'registeredmodel_id_example' # str | A unique identifier for a `RegisteredModel`.

    try:
        # Get a RegisteredModel
        api_response = await api_instance.get_registered_model(registeredmodel_id)
        print("The response of ModelRegistryServiceApi->get_registered_model:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_registered_model: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registeredmodel_id** | **str**| A unique identifier for a &#x60;RegisteredModel&#x60;. | 

### Return type

[**RegisteredModel**](RegisteredModel.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;RegisteredModel&#x60; entity. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_registered_model_versions**
> ModelVersionList get_registered_model_versions(registeredmodel_id, name=name, external_id=external_id, page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)

List All RegisteredModel's ModelVersions

Gets a list of all `ModelVersion` entities for the `RegisteredModel`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_version_list import ModelVersionList
from mr_openapi.models.order_by_field import OrderByField
from mr_openapi.models.sort_order import SortOrder
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    registeredmodel_id = 'registeredmodel_id_example' # str | A unique identifier for a `RegisteredModel`.
    name = 'entity-name' # str | Name of entity to search. (optional)
    external_id = '10' # str | External ID of entity to search. (optional)
    page_size = '100' # str | Number of entities in each page. (optional)
    order_by = mr_openapi.OrderByField() # OrderByField | Specifies the order by criteria for listing entities. (optional)
    sort_order = mr_openapi.SortOrder() # SortOrder | Specifies the sort order for listing entities, defaults to ASC. (optional)
    next_page_token = 'IkhlbGxvLCB3b3JsZC4i' # str | Token to use to retrieve next page of results. (optional)

    try:
        # List All RegisteredModel's ModelVersions
        api_response = await api_instance.get_registered_model_versions(registeredmodel_id, name=name, external_id=external_id, page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)
        print("The response of ModelRegistryServiceApi->get_registered_model_versions:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_registered_model_versions: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registeredmodel_id** | **str**| A unique identifier for a &#x60;RegisteredModel&#x60;. | 
 **name** | **str**| Name of entity to search. | [optional] 
 **external_id** | **str**| External ID of entity to search. | [optional] 
 **page_size** | **str**| Number of entities in each page. | [optional] 
 **order_by** | [**OrderByField**](.md)| Specifies the order by criteria for listing entities. | [optional] 
 **sort_order** | [**SortOrder**](.md)| Specifies the sort order for listing entities, defaults to ASC. | [optional] 
 **next_page_token** | **str**| Token to use to retrieve next page of results. | [optional] 

### Return type

[**ModelVersionList**](ModelVersionList.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a list of &#x60;ModelVersion&#x60; entities. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_registered_models**
> RegisteredModelList get_registered_models(page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)

List All RegisteredModels

Gets a list of all `RegisteredModel` entities.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.order_by_field import OrderByField
from mr_openapi.models.registered_model_list import RegisteredModelList
from mr_openapi.models.sort_order import SortOrder
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    page_size = '100' # str | Number of entities in each page. (optional)
    order_by = mr_openapi.OrderByField() # OrderByField | Specifies the order by criteria for listing entities. (optional)
    sort_order = mr_openapi.SortOrder() # SortOrder | Specifies the sort order for listing entities, defaults to ASC. (optional)
    next_page_token = 'IkhlbGxvLCB3b3JsZC4i' # str | Token to use to retrieve next page of results. (optional)

    try:
        # List All RegisteredModels
        api_response = await api_instance.get_registered_models(page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)
        print("The response of ModelRegistryServiceApi->get_registered_models:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_registered_models: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page_size** | **str**| Number of entities in each page. | [optional] 
 **order_by** | [**OrderByField**](.md)| Specifies the order by criteria for listing entities. | [optional] 
 **sort_order** | [**SortOrder**](.md)| Specifies the sort order for listing entities, defaults to ASC. | [optional] 
 **next_page_token** | **str**| Token to use to retrieve next page of results. | [optional] 

### Return type

[**RegisteredModelList**](RegisteredModelList.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a list of &#x60;RegisteredModel&#x60; entities. |  -  |
**401** | Unauthorized |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_serving_environment**
> ServingEnvironment get_serving_environment(servingenvironment_id)

Get a ServingEnvironment

Gets the details of a single instance of a `ServingEnvironment`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.serving_environment import ServingEnvironment
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    servingenvironment_id = 'servingenvironment_id_example' # str | A unique identifier for a `ServingEnvironment`.

    try:
        # Get a ServingEnvironment
        api_response = await api_instance.get_serving_environment(servingenvironment_id)
        print("The response of ModelRegistryServiceApi->get_serving_environment:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_serving_environment: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **servingenvironment_id** | **str**| A unique identifier for a &#x60;ServingEnvironment&#x60;. | 

### Return type

[**ServingEnvironment**](ServingEnvironment.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;ServingEnvironment&#x60; entity. |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_serving_environments**
> ServingEnvironmentList get_serving_environments(page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)

List All ServingEnvironments

Gets a list of all `ServingEnvironment` entities.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.order_by_field import OrderByField
from mr_openapi.models.serving_environment_list import ServingEnvironmentList
from mr_openapi.models.sort_order import SortOrder
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    page_size = '100' # str | Number of entities in each page. (optional)
    order_by = mr_openapi.OrderByField() # OrderByField | Specifies the order by criteria for listing entities. (optional)
    sort_order = mr_openapi.SortOrder() # SortOrder | Specifies the sort order for listing entities, defaults to ASC. (optional)
    next_page_token = 'IkhlbGxvLCB3b3JsZC4i' # str | Token to use to retrieve next page of results. (optional)

    try:
        # List All ServingEnvironments
        api_response = await api_instance.get_serving_environments(page_size=page_size, order_by=order_by, sort_order=sort_order, next_page_token=next_page_token)
        print("The response of ModelRegistryServiceApi->get_serving_environments:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->get_serving_environments: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page_size** | **str**| Number of entities in each page. | [optional] 
 **order_by** | [**OrderByField**](.md)| Specifies the order by criteria for listing entities. | [optional] 
 **sort_order** | [**SortOrder**](.md)| Specifies the sort order for listing entities, defaults to ASC. | [optional] 
 **next_page_token** | **str**| Token to use to retrieve next page of results. | [optional] 

### Return type

[**ServingEnvironmentList**](ServingEnvironmentList.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a list of &#x60;ServingEnvironment&#x60; entities. |  -  |
**401** | Unauthorized |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **update_inference_service**
> InferenceService update_inference_service(inferenceservice_id, inference_service_update)

Update a InferenceService

Updates an existing `InferenceService`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.inference_service import InferenceService
from mr_openapi.models.inference_service_update import InferenceServiceUpdate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    inferenceservice_id = 'inferenceservice_id_example' # str | A unique identifier for a `InferenceService`.
    inference_service_update = mr_openapi.InferenceServiceUpdate() # InferenceServiceUpdate | Updated `InferenceService` information.

    try:
        # Update a InferenceService
        api_response = await api_instance.update_inference_service(inferenceservice_id, inference_service_update)
        print("The response of ModelRegistryServiceApi->update_inference_service:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->update_inference_service: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **inferenceservice_id** | **str**| A unique identifier for a &#x60;InferenceService&#x60;. | 
 **inference_service_update** | [**InferenceServiceUpdate**](InferenceServiceUpdate.md)| Updated &#x60;InferenceService&#x60; information. | 

### Return type

[**InferenceService**](InferenceService.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;InferenceService&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **update_model_artifact**
> ModelArtifact update_model_artifact(modelartifact_id, model_artifact_update)

Update a ModelArtifact

Updates an existing `ModelArtifact`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_artifact import ModelArtifact
from mr_openapi.models.model_artifact_update import ModelArtifactUpdate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    modelartifact_id = 'modelartifact_id_example' # str | A unique identifier for a `ModelArtifact`.
    model_artifact_update = mr_openapi.ModelArtifactUpdate() # ModelArtifactUpdate | Updated `ModelArtifact` information.

    try:
        # Update a ModelArtifact
        api_response = await api_instance.update_model_artifact(modelartifact_id, model_artifact_update)
        print("The response of ModelRegistryServiceApi->update_model_artifact:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->update_model_artifact: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **modelartifact_id** | **str**| A unique identifier for a &#x60;ModelArtifact&#x60;. | 
 **model_artifact_update** | [**ModelArtifactUpdate**](ModelArtifactUpdate.md)| Updated &#x60;ModelArtifact&#x60; information. | 

### Return type

[**ModelArtifact**](ModelArtifact.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;ModelArtifact&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **update_model_version**
> ModelVersion update_model_version(modelversion_id, model_version_update)

Update a ModelVersion

Updates an existing `ModelVersion`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.model_version import ModelVersion
from mr_openapi.models.model_version_update import ModelVersionUpdate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    modelversion_id = 'modelversion_id_example' # str | A unique identifier for a `ModelVersion`.
    model_version_update = mr_openapi.ModelVersionUpdate() # ModelVersionUpdate | Updated `ModelVersion` information.

    try:
        # Update a ModelVersion
        api_response = await api_instance.update_model_version(modelversion_id, model_version_update)
        print("The response of ModelRegistryServiceApi->update_model_version:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->update_model_version: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **modelversion_id** | **str**| A unique identifier for a &#x60;ModelVersion&#x60;. | 
 **model_version_update** | [**ModelVersionUpdate**](ModelVersionUpdate.md)| Updated &#x60;ModelVersion&#x60; information. | 

### Return type

[**ModelVersion**](ModelVersion.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;ModelVersion&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **update_registered_model**
> RegisteredModel update_registered_model(registeredmodel_id, registered_model_update)

Update a RegisteredModel

Updates an existing `RegisteredModel`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.registered_model import RegisteredModel
from mr_openapi.models.registered_model_update import RegisteredModelUpdate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    registeredmodel_id = 'registeredmodel_id_example' # str | A unique identifier for a `RegisteredModel`.
    registered_model_update = mr_openapi.RegisteredModelUpdate() # RegisteredModelUpdate | Updated `RegisteredModel` information.

    try:
        # Update a RegisteredModel
        api_response = await api_instance.update_registered_model(registeredmodel_id, registered_model_update)
        print("The response of ModelRegistryServiceApi->update_registered_model:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->update_registered_model: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **registeredmodel_id** | **str**| A unique identifier for a &#x60;RegisteredModel&#x60;. | 
 **registered_model_update** | [**RegisteredModelUpdate**](RegisteredModelUpdate.md)| Updated &#x60;RegisteredModel&#x60; information. | 

### Return type

[**RegisteredModel**](RegisteredModel.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;RegisteredModel&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **update_serving_environment**
> ServingEnvironment update_serving_environment(servingenvironment_id, serving_environment_update)

Update a ServingEnvironment

Updates an existing `ServingEnvironment`.

### Example

* Bearer (JWT) Authentication (Bearer):

```python
import mr_openapi
from mr_openapi.models.serving_environment import ServingEnvironment
from mr_openapi.models.serving_environment_update import ServingEnvironmentUpdate
from mr_openapi.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to https://localhost:8080
# See configuration.py for a list of all supported configuration parameters.
configuration = mr_openapi.Configuration(
    host = "https://localhost:8080"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure Bearer authorization (JWT): Bearer
configuration = mr_openapi.Configuration(
    access_token = os.environ["BEARER_TOKEN"]
)

# Enter a context with an instance of the API client
async with mr_openapi.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = mr_openapi.ModelRegistryServiceApi(api_client)
    servingenvironment_id = 'servingenvironment_id_example' # str | A unique identifier for a `ServingEnvironment`.
    serving_environment_update = mr_openapi.ServingEnvironmentUpdate() # ServingEnvironmentUpdate | Updated `ServingEnvironment` information.

    try:
        # Update a ServingEnvironment
        api_response = await api_instance.update_serving_environment(servingenvironment_id, serving_environment_update)
        print("The response of ModelRegistryServiceApi->update_serving_environment:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ModelRegistryServiceApi->update_serving_environment: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **servingenvironment_id** | **str**| A unique identifier for a &#x60;ServingEnvironment&#x60;. | 
 **serving_environment_update** | [**ServingEnvironmentUpdate**](ServingEnvironmentUpdate.md)| Updated &#x60;ServingEnvironment&#x60; information. | 

### Return type

[**ServingEnvironment**](ServingEnvironment.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | A response containing a &#x60;ServingEnvironment&#x60; entity. |  -  |
**400** | Bad Request parameters |  -  |
**401** | Unauthorized |  -  |
**404** | The specified resource was not found |  -  |
**500** | Unexpected internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

