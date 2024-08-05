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
    ${rm_err}   Create Dictionary  code=  message=json: unknown field "ext_id"
          And Should be equal    ${rm_err}    ${err.json()}
