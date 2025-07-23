from model_registry.core import ModelRegistryAPIClient
from model_registry.types import ModelArtifact, ModelVersion, RegisteredModel
from robot.libraries.BuiltIn import BuiltIn
import model_registry.utils


def write_to_console(s):
    print(s)
    BuiltIn().log_to_console(s)


class ModelRegistry:
    def __init__(self, host: str = "http://localhost", port: int = 8080):
        self.api = ModelRegistryAPIClient.insecure_connection(host, port)

    async def upsert_registered_model(self, registered_model: dict) -> str:
        return (
            await self.api.upsert_registered_model(RegisteredModel(**registered_model))
        ).id

    async def upsert_model_version(
        self, model_version: dict, registered_model_id: str
    ) -> str:
        write_to_console(model_version)
        p = ModelVersion(**model_version)
        write_to_console(p)
        return (await self.api.upsert_model_version(p, registered_model_id)).id

    async def upsert_model_artifact(
        self, model_artifact: dict, model_version_id: str
    ) -> str:
        write_to_console(model_artifact)
        p = ModelArtifact(**model_artifact)
        write_to_console(p)
        return (await self.api.upsert_model_artifact(p, model_version_id)).id

    def s3_uri_from(self, path, bucket, endpoint, region) -> str:
        """
        Expose util to RobotFramework
        """
        return model_registry.utils.s3_uri_from(path=path, bucket=bucket, endpoint=endpoint, region=region)


async def test():
    demo_instance = ModelRegistry()
    await demo_instance.upsert_registered_model({"name": "testing123"})
    await demo_instance.upsert_model_version({"name": "v1"}, None)


# Used only for quick smoke tests
if __name__ == "__main__":
    import asyncio

    asyncio.run(test())
