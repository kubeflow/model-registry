*** Settings ***
Resource    Setup.resource
Resource    MRkeywords.resource
Test Setup    Test Setup with dummy data


*** Comments ***
Regression tests for Model Registry


*** Test Cases ***
As a MLOps engineer if I try to store a malformed RegisteredModel I get a structured error message
    ${rm}   Create Dictionary  name="model"  ext_id=123
    ${err}    POST    url=http://${MR_HOST}:${MR_PORT}/api/model_registry/v1alpha3/registered_models    json=&{rm}    expected_status=400
    ${rm_err}   Create Dictionary  code=Bad Request  message=json: unknown field "ext_id"
          And Should be equal    ${rm_err}    ${err.json()}

As a MLOps engineer if I try to store a malformed ModelVersion I get a structured error message
    ${rId}  Given I create a RegisteredModel having    name=${name}
    ${mv}   Create Dictionary   registeredModelId=${rId}
    ${err}    POST    url=http://${MR_HOST}:${MR_PORT}/api/model_registry/v1alpha3/model_versions    json=&{mv}    expected_status=422
    ${mv_err}   Create Dictionary  code=Bad Request  message=required field 'name' is zero value.
          And Should be equal    ${mv_err}    ${err.json()}
