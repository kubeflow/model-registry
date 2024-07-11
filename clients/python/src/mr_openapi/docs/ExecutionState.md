# ExecutionState

The state of the Execution. The state transitions are   NEW -> RUNNING -> COMPLETE | CACHED | FAILED | CANCELED CACHED means the execution is skipped due to cached results. CANCELED means the execution is skipped due to precondition not met. It is different from CACHED in that a CANCELED execution will not have any event associated with it. It is different from FAILED in that there is no unexpected error happened and it is regarded as a normal state.  See also: ml-metadata Execution.State

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


