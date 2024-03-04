*** Settings ***
Library    MLMetadata.py
Library    Collections
Resource   Setup.resource
Resource   MRkeywords.resource
Test Setup    Test Setup with dummy data


*** Comments ***
You can switch between REST and Python flow by environment variables,
as documented in the keyword implementation


*** Test Cases ***
Verify basic logical mapping between MR and MLMD
    Comment    This test ensures basic logical mapping bewteen MR entities and MLMD entities
    ...    based on the MR logical mapping:
    ...    RegisteredModel shall result in a MLMD Context
    ...    ModelVersion shall result in a MLMD Context and parent Context(of RegisteredModel)
    ...    ModelArtifact shall result in a MLMD Artifact and Attribution(to the parent Context of ModelVersion)

    ${rId}    Given I create a RegisteredModel having    name=${name}
    ${vId}    And I create a child ModelVersion having    registeredModelID=${rId}    name=v1
    ${aId}    And I create a child ModelArtifact having    modelversionId=${vId}  uri=s3://12345
    ${rId}    Convert To Integer    ${rId}
    ${vId}    Convert To Integer    ${vId}
    ${aId}    Convert To Integer    ${aId}

    # RegisteredModel shall result in a MLMD Context
    ${mlmdProto}    Get Context By Single Id    ${rId}
    Log To Console    ${mlmdProto}
    Should be equal    ${mlmdProto.type}    odh.RegisteredModel
    Should be equal    ${mlmdProto.name}    ${name}

    # ModelVersion shall result in a MLMD Context and parent Context(of RegisteredModel)
    ${mlmdProto}    Get Context By Single Id    ${vId}
    Log To Console    ${mlmdProto}
    Should be equal    ${mlmdProto.type}    odh.ModelVersion
    Should be equal    ${mlmdProto.name}    ${rId}:v1
    ${mlmdProto}    Get Parent Contexts By Context    ${vId}
    Should be equal    ${mlmdProto[0].id}    ${rId}

    # ModelArtifact shall result in a MLMD Artifact and Attribution(to the parent Context of ModelVersion)
    ${aNamePrefix}    Set Variable    ${vId}:
    ${mlmdProto}    Get Artifact By Single Id    ${aId}
    Log To Console    ${mlmdProto}
    Should be equal    ${mlmdProto.type}    odh.ModelArtifact
    Should Start With   ${mlmdProto.name}    ${aNamePrefix}
    Should be equal   ${mlmdProto.uri}    s3://12345
    ${mlmdProto}    Get Artifacts By Context    ${vId}
    Should be equal   ${mlmdProto[0].id}    ${aId}

Verify logical mapping of description property between MR and MLMD
    Comment    This test ensures logical mapping of the description bewteen MR entities and MLMD entities
    ...    being implemented as a custom_property

    Set To Dictionary    ${registered_model}    description=Lorem ipsum dolor sit amet  name=${name}
    Set To Dictionary    ${model_version}    description=consectetur adipiscing elit
    Set To Dictionary    ${model_artifact}    description=sed do eiusmod tempor incididunt
    ${rId}  Given I create a RegisteredModel    payload=${registered_model}
    ${vId}  And I create a child ModelVersion    registeredModelID=${rId}  payload=&{model_version}
    ${aId}  And I create a child ModelArtifact    modelversionId=${vId}  payload=&{model_artifact}
    ${rId}    Convert To Integer    ${rId}
    ${vId}    Convert To Integer    ${vId}
    ${aId}    Convert To Integer    ${aId}

    # RegisteredModel description
    ${mlmdProto}    Get Context By Single Id    ${rId}
    Should be equal    ${mlmdProto.properties['description'].string_value}    Lorem ipsum dolor sit amet

    # ModelVersion description
    ${mlmdProto}    Get Context By Single Id    ${vId}
    Should be equal    ${mlmdProto.properties['description'].string_value}    consectetur adipiscing elit

    # ModelArtifact description
    ${mlmdProto}    Get Artifact By Single Id    ${aId}
    Should be equal    ${mlmdProto.properties['description'].string_value}    sed do eiusmod tempor incididunt
