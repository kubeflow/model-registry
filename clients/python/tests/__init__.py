"""Tests for model registry."""

from dataclasses import dataclass
from typing import Generic, TypeVar

from model_registry.store import ProtoType
from model_registry.types.base import ProtoBase

P = TypeVar("P", bound=ProtoBase)


@dataclass
class Mapped(Generic[P]):
    proto: ProtoType
    py: P
