*** Settings ***
Library    RequestsLibrary
Library    Collections
Library    String
Library    ModelRegistry.py

*** Variables ***
${MR_HOST}     %{TEST_MR_HOST=localhost}
${MODE}        %{TEST_MODE=REST}

*** Keywords ***
I create a RegisteredModel having
    [Arguments]  ${name}
    IF  $MODE == "REST"
        ${data}    Create Dictionary    name=${name}
        ${resp}    POST    url=http://${MR_HOST}:8080/api/model_registry/v1alpha1/registered_models    json=${data}    expected_status=201
        Log to console    ${resp.json()}
        ${result}    Set Variable    ${resp.json()['id']}
    ELSE
        ${data}    Create Dictionary    name=${name}
        ${result}    Upsert Registered Model    ${data}
        Log to console    ${result}
    END
    RETURN    ${result}


I create a RegisteredModel
    [Arguments]  ${payload}
    IF  $MODE == "REST"
        ${resp}    POST    url=http://${MR_HOST}:8080/api/model_registry/v1alpha1/registered_models    json=&{payload}    expected_status=201
        Log to console    ${resp.json()}
        ${result}    Set Variable    ${resp.json()['id']}
    ELSE
        ${result}    Upsert Registered Model    ${payload}
        Log to console    ${result}
    END
    RETURN    ${result}


I create a child ModelVersion having
    [Arguments]  ${registeredModelID}  ${name}
    IF  $MODE == "REST"
        ${data}=    Create Dictionary    name=${name}    registeredModelID=${registeredModelID}
        ${resp}=    POST    url=http://${MR_HOST}:8080/api/model_registry/v1alpha1/model_versions    json=${data}    expected_status=201
        Log to console    ${resp.json()}
        ${result}    Set Variable    ${resp.json()['id']}
    ELSE
        ${data}    Create Dictionary    name=${name}
        ${result}    Upsert Model Version    ${data}     ${registeredModelID}
        Log to console    ${result}
    END
    RETURN    ${result}


I create a child ModelVersion
    [Arguments]  ${registeredModelID}  ${payload}
    IF  $MODE == "REST"
        Set To Dictionary    ${payload}    registeredModelID=${registeredModelID}
        ${resp}=    POST    url=http://${MR_HOST}:8080/api/model_registry/v1alpha1/model_versions    json=&{payload}    expected_status=201
        Log to console    ${resp.json()}
        ${result}    Set Variable    ${resp.json()['id']}
    ELSE
        ${result}    Upsert Model Version    ${payload}     ${registeredModelID}
        Log to console    ${result}
    END
    RETURN    ${result}


I create a child ModelArtifact having
    [Arguments]  ${modelversionId}    ${uri}
    IF  $MODE == "REST"
        ${data}=    Create Dictionary    uri=${uri}  artifactType=model-artifact
        Log to console    ${data}
        ${resp}=    POST    url=http://${MR_HOST}:8080/api/model_registry/v1alpha1/model_versions/${modelversionId}/artifacts    json=${data}    expected_status=201
        Log to console    ${resp.json()}
        ${result}    Set Variable    ${resp.json()['id']}
    ELSE
        ${data}    Create Dictionary    uri=${uri}
        ${result}    Upsert Model Artifact    ${data}    ${modelversionId}
        Log to console    ${result}
    END
    RETURN    ${result}


I create a child ModelArtifact
    [Arguments]  ${modelversionId}    ${payload}
    IF  $MODE == "REST"
        ${resp}=    POST    url=http://${MR_HOST}:8080/api/model_registry/v1alpha1/model_versions/${modelversionId}/artifacts    json=&{payload}    expected_status=201
        Log to console    ${resp.json()}
        ${result}    Set Variable    ${resp.json()['id']}
    ELSE
        ${result}    Upsert Model Artifact    ${payload}    ${modelversionId}
        Log to console    ${result}
    END
    RETURN    ${result}


I get RegisteredModelByID
    [Arguments]    ${id}
    IF  $MODE == "REST"
        ${resp}=    GET    url=http://${MR_HOST}:8080/api/model_registry/v1alpha1/registered_models/${id}    expected_status=200
        ${result}    Set Variable    ${resp.json()}
        Log to console    ${resp.json()}
    ELSE
        Log to console    ${MODE}
        Fail    Not Implemented
    END
    RETURN    ${result}


I get ModelVersionByID
    [Arguments]    ${id}
    IF  $MODE == "REST"
        ${resp}=    GET    url=http://${MR_HOST}:8080/api/model_registry/v1alpha1/model_versions/${id}    expected_status=200
        ${result}    Set Variable    ${resp.json()}
        Log to console    ${resp.json()}
    ELSE
        Log to console    ${MODE}
        Fail    Not Implemented
    END
    RETURN    ${result}


I get ModelArtifactByID
    [Arguments]    ${id}
    IF  $MODE == "REST"
        ${resp}=    GET    url=http://${MR_HOST}:8080/api/model_registry/v1alpha1/model_artifacts/${id}    expected_status=200
        ${result}    Set Variable    ${resp.json()}
        Log to console    ${resp.json()}
    ELSE
        Log to console    ${MODE}
        Fail    Not Implemented
    END
    RETURN    ${result}
