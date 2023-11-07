*** Settings ***
Library    RequestsLibrary
Library    String

*** Variables ***
${mr_host}     localhost

*** Keywords ***
I create a RegisteredModel having
    [Arguments]  ${name}
    ${data}=    Create Dictionary    name=${name}
    ${resp}=    POST    url=http://${mr_host}:8080/api/model_registry/v1alpha1/registered_models    json=${data}    expected_status=201
    Log to console    ${resp.json()}
    [return]    ${resp.json()["id"]}


I create a child ModelVersion having
    [Arguments]  ${registeredModelID}  ${name}
    ${data}=    Create Dictionary    name=${name}    registeredModelID=${registeredModelID}
    ${resp}=    POST    url=http://${mr_host}:8080/api/model_registry/v1alpha1/model_versions    json=${data}    expected_status=201
    Log to console    ${resp.json()}
    [return]    ${resp.json()["id"]}

I create a child ModelArtifact having
    [Arguments]  ${modelversionId}    ${uri}
    ${data}=    Create Dictionary    uri=${uri}  artifactType=model-artifact
    Log to console    ${data}
    ${resp}=    POST    url=http://${mr_host}:8080/api/model_registry/v1alpha1/model_versions/${modelversionId}/artifacts    json=${data}    expected_status=201
    Log to console    ${resp.json()}
    [return]    ${resp.json()["id"]}


I get RegisteredModelByID
    [Arguments]    ${id}
    ${resp}=    GET    url=http://${mr_host}:8080/api/model_registry/v1alpha1/registered_models/${id}    expected_status=200
    Log to console    ${resp.json()}
    [return]    ${resp.json()}


I get ModelVersionByID
    [Arguments]    ${id}
    ${resp}=    GET    url=http://${mr_host}:8080/api/model_registry/v1alpha1/model_versions/${id}    expected_status=200
    Log to console    ${resp.json()}
    [return]    ${resp.json()}


I get ModelArtifactByID
    [Arguments]    ${id}
    ${resp}=    GET    url=http://${mr_host}:8080/api/model_registry/v1alpha1/model_artifacts/${id}    expected_status=200
    Log to console    ${resp.json()}
    [return]    ${resp.json()}

