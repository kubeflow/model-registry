*** Settings ***
Resource    Setup.resource
Resource    MRkeywords.resource
Test Setup    Test Setup with dummy data


*** Comments ***
These User Story(-ies) are defined in the PM document


*** Test Cases ***
As a MLOps engineer I would like to store Model name
    ${rId}  Given I create a RegisteredModel having    name=${name}
    ${vId}  And I create a child ModelVersion having    registeredModelID=${rId}    name=v1
    ${aId}  And I create a child ModelArtifact having    modelversionId=${vId}  uri=s3://12345
    ${r}  Then I get RegisteredModelByID    id=${rId}
          And Should be equal    ${r["name"]}    ${name}
    ${r}  Then I get ModelVersionByID    id=${vId}
          And Should be equal    ${r["name"]}    v1
    ${r}  Then I get ModelArtifactByID    id=${aId}
          And Should be equal    ${r["uri"]}    s3://12345

As a MLOps engineer I would like to store a description of the model
    Set To Dictionary    ${registered_model}    description=Lorem ipsum dolor sit amet  name=${name}
    Set To Dictionary    ${model_version}    description=consectetur adipiscing elit
    Set To Dictionary    ${model_artifact}    description=sed do eiusmod tempor incididunt
    ${rId}  Given I create a RegisteredModel    payload=${registered_model}
    ${vId}  And I create a child ModelVersion    registeredModelID=${rId}  payload=&{model_version}
    ${aId}  And I create a child ModelArtifact    modelversionId=${vId}  payload=&{model_artifact}
    ${r}  Then I get RegisteredModelByID    id=${rId}
          And Should be equal    ${r["description"]}    Lorem ipsum dolor sit amet
    ${r}  Then I findRegisteredModel by name    name=${name}
          And Should be equal    ${r["description"]}    Lorem ipsum dolor sit amet
    ${r}  Then I get ModelVersionByID    id=${vId}
          And Should be equal    ${r["description"]}    consectetur adipiscing elit
    ${r}  Then I findModelVersion by name and parentResourceId    name=v1.2.3  parentResourceId=${rId}
          And Should be equal    ${r["description"]}    consectetur adipiscing elit
    ${r}  Then I get ModelArtifactByID    id=${aId}
          And Should be equal    ${r["description"]}    sed do eiusmod tempor incididunt
    ${r}  Then I findModelArtifact by name and parentResourceId    name=ModelArtifactName  parentResourceId=${vId}
          And Should be equal    ${r["description"]}    sed do eiusmod tempor incididunt

As a MLOps engineer I would like to store a longer documentation for the model
    Set To Dictionary    ${registered_model}    description=Lorem ipsum dolor sit amet  name=${name}
    Set To Dictionary    ${model_version}    description=consectetur adipiscing elit
    Set To Dictionary    ${model_artifact}    description=sed do eiusmod tempor incididunt
    Set To Dictionary    ${doc_artifact}    uri="https://README.md"
    ${rId}  Given I create a RegisteredModel    payload=${registered_model}
    ${vId}  And I create a child ModelVersion    registeredModelID=${rId}  payload=&{model_version}
    ${aId}  And I create a child ModelArtifact    modelversionId=${vId}  payload=&{model_artifact}
    ${docId}  And I create a child Artifact    modelversionId=${vId}  payload=&{doc_artifact}
    ${r}  Then I get RegisteredModelByID    id=${rId}
          And Should be equal    ${r["description"]}    Lorem ipsum dolor sit amet
    ${r}  Then I get ModelVersionByID    id=${vId}
          And Should be equal    ${r["description"]}    consectetur adipiscing elit
    ${r}  Then I get ModelArtifactByID    id=${aId}
          And Should be equal    ${r["description"]}    sed do eiusmod tempor incididunt
    ${r}    Then I get ArtifactsByModelVersionID    id=${vId}
    ${cnt}  Then Get length    ${r["items"]}
            And Should Be Equal As Integers    ${cnt}    2
