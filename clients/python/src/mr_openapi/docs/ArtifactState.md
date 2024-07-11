# ArtifactState

 - PENDING: A state indicating that the artifact may exist.  - LIVE: A state indicating that the artifact should exist, unless something external to the system deletes it.  - MARKED_FOR_DELETION: A state indicating that the artifact should be deleted.  - DELETED: A state indicating that the artifact has been deleted.  - ABANDONED: A state indicating that the artifact has been abandoned, which may be due to a failed or cancelled execution.  - REFERENCE: A state indicating that the artifact is a reference artifact. At execution start time, the orchestrator produces an output artifact for each output key with state PENDING. However, for an intermediate artifact, this first artifact's state will be REFERENCE. Intermediate artifacts emitted during a component's execution will copy the REFERENCE artifact's attributes. At the end of an execution, the artifact state should remain REFERENCE instead of being changed to LIVE.  See also: ml-metadata Artifact.State

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


