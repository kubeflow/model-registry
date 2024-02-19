from grpc import insecure_channel

# from ml_metadata.metadata_store import metadata_store
from ml_metadata.proto import metadata_store_pb2
from ml_metadata.proto import metadata_store_service_pb2
from ml_metadata.proto import metadata_store_service_pb2_grpc


def main():
    # connection_config = metadata_store_pb2.ConnectionConfig()
    # connection_config.sqlite.filename_uri = './metadata.sqlite'
    # connection_config.sqlite.connection_mode = 3 # READWRITE_OPENCREATE
    # store = metadata_store.MetadataStore(connection_config)

    # connection_config = metadata_store_pb2.ConnectionConfig()
    # connection_config.mysql.host = 'localhost'
    # connection_config.mysql.port = 3306
    # connection_config.mysql.database = 'mlmetadata'
    # connection_config.mysql.user = 'root'
    # connection_config.mysql.password = 'my-secret-pw'
    # store = metadata_store.MetadataStore(connection_config, enable_upgrade_migration=True)

    channel = insecure_channel("localhost:8080")
    store = metadata_store_service_pb2_grpc.MetadataStoreServiceStub(channel)

    # Create ArtifactTypes, e.g., Data and Model
    data_type = metadata_store_pb2.ArtifactType()
    data_type.name = "DataSet"
    data_type.properties["day"] = metadata_store_pb2.INT
    data_type.properties["split"] = metadata_store_pb2.STRING

    request = metadata_store_service_pb2.PutArtifactTypeRequest()
    request.all_fields_match = True
    request.artifact_type.CopyFrom(data_type)
    response = store.PutArtifactType(request)
    data_type_id = response.type_id

    model_type = metadata_store_pb2.ArtifactType()
    model_type.name = "SavedModel"
    model_type.properties["version"] = metadata_store_pb2.INT
    model_type.properties["name"] = metadata_store_pb2.STRING

    request.artifact_type.CopyFrom(model_type)
    response = store.PutArtifactType(request)
    model_type_id = response.type_id

    request = metadata_store_service_pb2.GetArtifactTypeRequest()
    request.type_name = "SavedModel"
    response = store.GetArtifactType(request)
    assert response.artifact_type.id == model_type_id
    assert response.artifact_type.name == "SavedModel"

    # Query all registered Artifact types.
    # artifact_types = store.GetArtifactTypes()

    # Create an ExecutionType, e.g., Trainer
    trainer_type = metadata_store_pb2.ExecutionType()
    trainer_type.name = "Trainer"
    trainer_type.properties["state"] = metadata_store_pb2.STRING

    request = metadata_store_service_pb2.PutExecutionTypeRequest()
    request.execution_type.CopyFrom(trainer_type)
    response = store.PutExecutionType(request)
    # trainer_type_id = response.type_id

    # # Query a registered Execution type with the returned id
    # [registered_type] = store.GetExecutionTypesByID([trainer_type_id])

    # Create an input artifact of type DataSet
    data_artifact = metadata_store_pb2.Artifact()
    data_artifact.uri = "path/to/data"
    data_artifact.properties["day"].int_value = 1
    data_artifact.properties["split"].string_value = "train"
    data_artifact.type_id = data_type_id

    request = metadata_store_service_pb2.PutArtifactsRequest()
    request.artifacts.extend([data_artifact])
    response = store.PutArtifacts(request)
    # data_artifact_id = response.artifact_ids[0]

    # # Query all registered Artifacts
    # artifacts = store.GetArtifacts()
    #
    # # Plus, there are many ways to query the same Artifact
    # [stored_data_artifact] = store.GetArtifactsByID([data_artifact_id])
    # artifacts_with_uri = store.GetArtifactsByURI(data_artifact.uri)
    # artifacts_with_conditions = store.GetArtifacts(
    #       list_options=mlmd.ListOptions(
    #           filter_query='uri LIKE "%/data" AND properties.day.int_value > 0'))
    #
    # # Register the Execution of a Trainer run
    # trainer_run = metadata_store_pb2.Execution()
    # trainer_run.type_id = trainer_type_id
    # trainer_run.properties["state"].string_value = "RUNNING"
    # [run_id] = store.PutExecutions([trainer_run])
    #
    # # Query all registered Execution
    # executions = store.GetExecutionsByID([run_id])
    # # Similarly, the same execution can be queried with conditions.
    # executions_with_conditions = store.GetExecutions(
    #     list_options = mlmd.ListOptions(
    #         filter_query='type = "Trainer" AND properties.state.string_value IS NOT NULL'))
    #
    # # Define the input event
    # input_event = metadata_store_pb2.Event()
    # input_event.artifact_id = data_artifact_id
    # input_event.execution_id = run_id
    # input_event.type = metadata_store_pb2.Event.DECLARED_INPUT
    #
    # # Record the input event in the metadata store
    # store.put_events([input_event])
    #
    # # Declare the output artifact of type SavedModel
    # model_artifact = metadata_store_pb2.Artifact()
    # model_artifact.uri = 'path/to/model/file'
    # model_artifact.properties["version"].int_value = 1
    # model_artifact.properties["name"].string_value = 'MNIST-v1'
    # model_artifact.type_id = model_type_id
    # [model_artifact_id] = store.PutArtifacts([model_artifact])
    #
    # # Declare the output event
    # output_event = metadata_store_pb2.Event()
    # output_event.artifact_id = model_artifact_id
    # output_event.execution_id = run_id
    # output_event.type = metadata_store_pb2.Event.DECLARED_OUTPUT
    #
    # # Submit output event to the Metadata Store
    # store.PutEvents([output_event])
    #
    # trainer_run.id = run_id
    # trainer_run.properties["state"].string_value = "COMPLETED"
    # store.PutExecutions([trainer_run])

    # Create a ContextType, e.g., Experiment with a note property
    experiment_type = metadata_store_pb2.ContextType()
    experiment_type.name = "Experiment"
    experiment_type.properties["note"] = metadata_store_pb2.STRING
    request = metadata_store_service_pb2.PutContextTypeRequest()
    request.context_type.CopyFrom(experiment_type)
    response = store.PutContextType(request)
    # experiment_type_id = response.type_id

    # # Group the model and the trainer run to an experiment.
    # my_experiment = metadata_store_pb2.Context()
    # my_experiment.type_id = experiment_type_id
    # # Give the experiment a name
    # my_experiment.name = "exp1"
    # my_experiment.properties["note"].string_value = "My first experiment."
    # [experiment_id] = store.PutContexts([my_experiment])
    #
    # attribution = metadata_store_pb2.Attribution()
    # attribution.artifact_id = model_artifact_id
    # attribution.context_id = experiment_id
    #
    # association = metadata_store_pb2.Association()
    # association.execution_id = run_id
    # association.context_id = experiment_id
    #
    # store.PutAttributionsAndAssociations([attribution], [association])
    #
    # # Query the Artifacts and Executions that are linked to the Context.
    # experiment_artifacts = store.GetArtifactsByContext(experiment_id)
    # experiment_executions = store.GetExecutionsByContext(experiment_id)
    #
    # # You can also use neighborhood queries to fetch these artifacts and executions
    # # with conditions.
    # experiment_artifacts_with_conditions = store.GetArtifacts(
    #     list_options = mlmd.ListOptions(
    #         filter_query=('contexts_a.type = "Experiment" AND contexts_a.name = "exp1"')))
    # experiment_executions_with_conditions = store.GetExecutions(
    #     list_options = mlmd.ListOptions(
    #         filter_query=('contexts_a.id = {}'.format(experiment_id))))


if __name__ == "__main__":
    main()
