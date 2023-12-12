import model_registry as mr
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel
from robot.libraries.BuiltIn import BuiltIn


def write_to_console(s):
    print(s)
    BuiltIn().log_to_console(s)


class ModelRegistry(mr.core.ModelRegistryAPIClient):
    def __init__(self, host: str = "localhost", port: int = 9090):
        super().__init__(host, port)

    def upsert_registered_model(self, registered_model) -> str:
        p = RegisteredModel(None)
        for key, value in registered_model.items():
            setattr(p, key, value)
        return super().upsert_registered_model(p)

    def upsert_model_version(self, model_version, registered_model_id: str) -> str:
        write_to_console(model_version)
        p = ModelVersion(ModelArtifact("", ""), "", "")
        for key, value in model_version.items():
            setattr(p, key, value)
        write_to_console(p)
        return super().upsert_model_version(p, registered_model_id)

    def upsert_model_artifact(self, model_artifact, model_version_id: str) -> str:
        write_to_console(model_artifact)
        p = ModelArtifact(None, None)
        for key, value in model_artifact.items():
            setattr(p, key, value)
        write_to_console(p)
        return super().upsert_model_artifact(p, model_version_id)


# Used only for quick smoke tests
if __name__ == "__main__":
    demo_instance = ModelRegistry()
    demo_instance.upsert_registered_model({"name": "testing123"})
    demo_instance.upsert_model_version({"name": "v1"}, None)
