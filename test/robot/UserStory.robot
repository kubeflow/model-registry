*** Settings ***
Resource   MRviaREST.resource


*** Comments ***
These User Story(-ies) are defined in the PM document


*** Test Cases ***              
As a MLOps engineer I would like to store Model name
    ${name}=  Generate Random String    8    [LETTERS] 
    ${rId}=  Given I create a RegisteredModel having    name=${name}
    ${vId}=  And I create a child ModelVersion having    registeredModelID=${rId}    name=v1  
    ${aId}=  And I create a child ModelArtifact having    modelversionId=${vId}  uri=s3://12345
    ${r}=  Then I get RegisteredModelByID    id=${rId}
        And Should be equal    ${r["name"]}    ${name}
    ${r}=  Then I get ModelVersionByID    id=${vId}
        And Should be equal    ${r["name"]}    v1
    ${r}=  Then I get ModelArtifactByID    id=${aId}
        And Should be equal    ${r["uri"]}    s3://12345
