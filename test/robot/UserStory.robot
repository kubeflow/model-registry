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
          And Should be equal    ${r["registeredModelId"]}    ${rId}
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

As a MLOps engineer I would like to update a description of the model
    Set To Dictionary    ${registered_model}    description=Lorem ipsum dolor sit amet  name=${name}  owner=me
    ${rId}  Given I create a RegisteredModel    payload=${registered_model}
    &{registered_model_update}  Create dictionary    description=New description
    ${rId}  And I update a RegisteredModel    id=${rId}  payload=${registered_model_update}
    ${r}  Then I get RegisteredModelByID    id=${rId}
          And Should be equal    ${r["description"]}    New description
          And Should be equal    ${r["owner"]}    me

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

As a MLOps engineer I would like to store some labels
    # MIGRATED TO test_update_logical_model_with_labels in pytest
    # A custom property of type string, with empty string value, shall be considered a Label; this is also semantically compatible for properties having empty string values in general.
    ${cp1}    Create Dictionary  my-label1=${{ {"string_value": "", "metadataType": "MetadataStringValue"} }}  my-label2=${{ {"string_value": "", "metadataType": "MetadataStringValue"} }}
    Set To Dictionary    ${registered_model}    description=Lorem ipsum dolor sit amet  name=${name}  customProperties=${cp1}
    ${cp2}    Create Dictionary  my-label3=${{ {"string_value": "", "metadataType": "MetadataStringValue"} }}  my-label4=${{ {"string_value": "", "metadataType": "MetadataStringValue"} }}
    Set To Dictionary    ${model_version}    description=consectetur adipiscing elit  customProperties=${cp2}
    ${cp3}    Create Dictionary  my-label5=${{ {"string_value": "", "metadataType": "MetadataStringValue"} }}  my-label6=${{ {"string_value": "", "metadataType": "MetadataStringValue"} }}
    Set To Dictionary    ${model_artifact}    description=sed do eiusmod tempor incididunt  customProperties=${cp3}
    ${rId}  Given I create a RegisteredModel    payload=${registered_model}
    ${vId}  And I create a child ModelVersion    registeredModelID=${rId}  payload=&{model_version}
    ${aId}  And I create a child ModelArtifact    modelversionId=${vId}  payload=&{model_artifact}
    ${r}  Then I get RegisteredModelByID    id=${rId}
          And Should be equal    ${r["description"]}    Lorem ipsum dolor sit amet
          And Dictionaries Should Be Equal   ${r["customProperties"]}  ${cp1}
    ${r}  Then I get ModelVersionByID    id=${vId}
          And Should be equal    ${r["description"]}    consectetur adipiscing elit
          And Dictionaries Should Be Equal   ${r["customProperties"]}  ${cp2}
    ${r}  Then I get ModelArtifactByID    id=${aId}
          And Should be equal    ${r["description"]}    sed do eiusmod tempor incididunt
          And Dictionaries Should Be Equal   ${r["customProperties"]}  ${cp3}

As a MLOps engineer I would like to store an owner for the RegisteredModel
    # MIGRATED TO test_register_model_with_owner in pytest
    Set To Dictionary    ${registered_model}    description=Lorem ipsum dolor sit amet  name=${name}  owner=My owner
    ${rId}  Given I create a RegisteredModel    payload=${registered_model}
    ${r}  Then I get RegisteredModelByID    id=${rId}
          And Should be equal    ${r["description"]}    Lorem ipsum dolor sit amet
          And Should be equal    ${r["owner"]}    My owner

As a MLOps engineer I want to track a Model from an S3 bucket Data Connection
    # MIGRATED TO test_register_model_with_s3_data_connection in pytest
    ${data_connection_name}  Set Variable   aws-connection-my-data-connection
    ${s3_bucket}  Set Variable  my-bucket
    ${s3_path}  Set Variable  my-path
    ${s3_endpoint}  Set Variable  https://minio-api.acme.org
    ${s3_region}  Set Variable  us-east-1

    Set To Dictionary    ${registered_model}    description=The Model  name=${name}
    Set To Dictionary    ${model_version}    description=The Version  name=v1.0  
    ${art_uri}  I use util s3_uri_from  path=${s3_path}  bucket=${s3_bucket}  endpoint=${s3_endpoint}  region=${s3_region}
    Set To Dictionary    ${model_artifact}    description=The Artifact  uri=${art_uri}  modelFormatName=onnx  modelFormatVersion=1  storageKey=${data_connection_name}  storagePath=${s3_path}
    Log to console  ${model_artifact}

    ${rId}  Given I create a RegisteredModel    payload=${registered_model}
    ${vId}  And I create a child ModelVersion    registeredModelID=${rId}  payload=&{model_version}
    ${aId}  And I create a child ModelArtifact    modelversionId=${vId}  payload=&{model_artifact}
    ${r}  Then I get RegisteredModelByID    id=${rId}
          And Should be equal    ${r["description"]}    The Model
    ${r}  Then I get ModelVersionByID    id=${vId}
          And Should be equal    ${r["description"]}    The Version
          And Should be equal    ${r["name"]}    v1.0
    ${r}  Then I get ModelArtifactByID    id=${aId}
          And Should be equal    ${r["description"]}    The Artifact
          And Should be equal    ${r["uri"]}    s3://my-bucket/my-path?endpoint=https://minio-api.acme.org&defaultRegion=us-east-1
          And Should be equal    ${r["modelFormatName"]}    onnx
          And Should be equal    ${r["modelFormatVersion"]}    1
          And Should be equal    ${r["storageKey"]}    aws-connection-my-data-connection
          And Should be equal    ${r["storagePath"]}    my-path
