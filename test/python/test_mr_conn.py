from model_registry import ModelRegistry


def main(server: str, port: int):
    mr = ModelRegistry(server, port, author="test", is_secure=False)

    model = mr.register_model(
        "my-model",
        "https://mybucket.uri/",
        version="2.0.0",
        model_format_name="onnx",
        model_format_version="1",
        storage_key="my-data-connection",
        storage_path="path/to/model",
        metadata={
            "day": 1,
            "split": "train",
        },
    )

    m = mr.get_registered_model("my-model")
    assert m
    assert model.id == m.id, f"{model} != {m}"


if __name__ == "__main__":
    import sys

    main(sys.argv[1], int(sys.argv[2]))
