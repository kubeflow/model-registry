from model_registry.core import ModelRegistryAPIClient
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel
from robot.libraries.BuiltIn import BuiltIn


def write_to_console(s):
    print(s)
    BuiltIn().log_to_console(s)


class ModelRegistry:
    def __init__(self, host: str = "http://localhost", port: int = 9090):
        self.api = ModelRegistryAPIClient.insecure_connection(host, port)

    def upsert_registered_model(self, registered_model: dict) -> str:
        return self.api.upsert_registered_model(RegisteredModel(**registered_model))

    def upsert_model_version(
        self, model_version: dict, registered_model_id: str
    ) -> str:
        write_to_console(model_version)
        p = ModelVersion(**model_version)
        write_to_console(p)
        return self.api.upsert_model_version(p, registered_model_id)

    def upsert_model_artifact(self, model_artifact: dict, model_version_id: str) -> str:
        write_to_console(model_artifact)
        p = ModelArtifact(**model_artifact)
        write_to_console(p)
        return self.api.upsert_model_artifact(p, model_version_id)


# Used only for quick smoke tests
if __name__ == "__main__":
    demo_instance = ModelRegistry()
    demo_instance.upsert_registered_model({"name": "testing123"})
    demo_instance.upsert_model_version({"name": "v1"}, None)
