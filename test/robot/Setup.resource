*** Settings ***
Library    yaml
Library    OperatingSystem
Library    Collections
Library    String

*** Keywords ***
Test Setup with dummy data
    Log To Console  TEST SETUP
    ${name}    Generate Random String    8    [LETTERS]
    Set Test Variable    $name    ${name}
    ${YAML}    Get File  ${CURDIR}${/}data.yaml
    ${YAML}    yaml.Safe Load  ${YAML}
    &{registered_model}    Get From Dictionary   ${YAML}    registered_model
    Set Test Variable    &registered_model    &{registered_model}
    &{model_version}    Get From Dictionary   ${YAML}    model_version
    Set Test Variable    &model_version    &{model_version}
    &{model_artifact}    Get From Dictionary   ${YAML}    model_artifact
    Set Test Variable    &model_artifact    &{model_artifact}


*** Variables ***
&{registered_model}    
&{model_version}    
&{model_artifact}
${name}