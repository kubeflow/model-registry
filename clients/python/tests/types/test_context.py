from model_registry.types.contexts import ModelVersion, RegisteredModel


def test_model_version_repr():
    """Test the repr of a ModelVersion, with name to the front, ensuring other fields are also present"""
    mv = ModelVersion(name="Version 1", id="123", author="test_author", description="test_description")
    print(mv)
    assert str(mv).startswith("name='Version 1'")
    assert "id='123'" in str(mv)
    assert "author='test_author'" in str(mv)
    assert "description='test_description'" in str(mv)
    print(repr(mv))
    assert repr(mv).startswith("ModelVersion(name='Version 1'")
    assert "id='123'" in str(mv)
    assert "author='test_author'" in str(mv)
    assert "description='test_description'" in str(mv)

    # Test with empty name, not really used in practice but to increase coverage
    mv = ModelVersion(name="", id="123", author="test_author", description="test_description")
    print(mv)
    assert str(mv).startswith("name=''")
    assert "id='123'" in str(mv)
    assert "author='test_author'" in str(mv)
    assert "description='test_description'" in str(mv)
    print(repr(mv))
    assert repr(mv).startswith("ModelVersion(name=''")
    assert "id='123'" in str(mv)
    assert "author='test_author'" in str(mv)
    assert "description='test_description'" in str(mv)


def test_registered_model_repr():
    """Test the repr of a RegisteredModel, with name to the front, ensuring other fields are also present"""
    rm = RegisteredModel(name="Model 1", id="123", owner="test_owner", description="test_description")
    print(rm)
    assert str(rm).startswith("name='Model 1'")
    assert "id='123'" in str(rm)
    assert "owner='test_owner'" in str(rm)
    assert "description='test_description'" in str(rm)
    print(repr(rm))
    assert repr(rm).startswith("RegisteredModel(name='Model 1'")
    assert "id='123'" in str(rm)
    assert "owner='test_owner'" in str(rm)
    assert "description='test_description'" in str(rm)

    # Test with empty name, not really used in practice but to increase coverage
    rm = RegisteredModel(name="", id="123", owner="test_owner", description="test_description")
    print(rm)
    assert str(rm).startswith("name=''")
    assert "id='123'" in str(rm)
    assert "owner='test_owner'" in str(rm)
    assert "description='test_description'" in str(rm)
    print(repr(rm))
    assert repr(rm).startswith("RegisteredModel(name=''")
    assert "id='123'" in str(rm)
    assert "owner='test_owner'" in str(rm)
    assert "description='test_description'" in str(rm)
