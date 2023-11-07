*** Settings ***
Library    String
Resource   MRviaREST.resource
Library    MLMetadata.py


*** Comments ***
You can replace the underlying keyword definition to go via:
- REST entrypoint
- Python entrypoint
for example by changing the imported `Resource MRviaREST.resource`
with an alternative implementation, etc.


*** Test Cases ***              
Verify basic logical mapping between MR and MLMD
    Comment    This test ensures basic logical mapping bewteen MR entities and MLMD entities
    ...    based on the MR logical mapping:
    ...    RegisteredModel shall result in a MLMD Context
    ...    ModelVersion shall result in a MLMD Context and parent Context(of RegisteredModel)
    ...    ModelArtifact shall result in a MLMD Artifact and Attribution(to the parent Context of ModelVersion)

    ${name}=  Generate Random String    8    [LETTERS] 
    ${rId}=  Given I create a RegisteredModel having    name=${name}
    ${vId}=  And I create a child ModelVersion having    registeredModelID=${rId}    name=v1  
    ${aId}=  And I create a child ModelArtifact having    modelversionId=${vId}  uri=s3://12345
    ${rId}=    Convert To Integer    ${rId}
    ${vId}=    Convert To Integer    ${vId}
    ${aId}    Convert To Integer    ${aId}

    # RegisteredModel shall result in a MLMD Context
    @{singList}    Create List    ${rId}
    ${mlmdProto}    Get Contexts By Id    ${singList}
    Log To Console    ${mlmdProto[0]}
    Should be equal    ${mlmdProto[0].type}    odh.RegisteredModel
    Should be equal    ${mlmdProto[0].name}    ${name}

    # ModelVersion shall result in a MLMD Context and parent Context(of RegisteredModel)
    @{singList}    Create List    ${vId}
    ${mlmdProto}    Get Contexts By Id    ${singList}
    Log To Console    ${mlmdProto[0]}
    Should be equal    ${mlmdProto[0].type}    odh.ModelVersion
    Should be equal    ${mlmdProto[0].name}    ${rId}:v1
    ${mlmdProto}    Get Parent Contexts By Context    ${vId}
    Should be equal    ${mlmdProto[0].id}    ${rId}

    # ModelArtifact shall result in a MLMD Artifact and Attribution(to the parent Context of ModelVersion)
    ${aNamePrefix}    Set Variable    ${vId}:
    @{singList}    Create List    ${aId}
    ${mlmdProto}    Get Artifacts By Id    ${singList}
    Log To Console    ${mlmdProto[0]}
    Should be equal    ${mlmdProto[0].type}    odh.ModelArtifact
    Should Start With   ${mlmdProto[0].name}    ${aNamePrefix}
    Should be equal   ${mlmdProto[0].uri}    s3://12345
    ${mlmdProto}    Get Artifacts By Context    ${vId}
    Should be equal   ${mlmdProto[0].id}    ${aId}
