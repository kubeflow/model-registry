"""Logging utilities for signing components."""

import logging


class LogConfig:
    """Configuration for instance-level logging."""

    def __init__(self, instance_name: str, level: int | None = None):
        self.instance_name: str = instance_name
        self.level: int = level if level is not None else logging.WARNING


class InstanceLevelAdapter(logging.LoggerAdapter):
    """Checks instance log level before logging.

    Uses a LogConfig dataclass for instance-level filtering,
    while keeping extra as a standard dict for LoggerAdapter compatibility.
    """

    def __init__(self, logger: logging.Logger, log_config: LogConfig):
        super().__init__(logger, {"instance_name": log_config.instance_name})
        self.log_config = log_config

    def log(self, level, msg, *args, **kwargs):
        if level >= self.log_config.level:
            super().log(level, msg, *args, **kwargs)

    def process(self, msg, kwargs):
        return f"[{self.log_config.instance_name}] {msg}", kwargs

    def set_log_level(self, level: int) -> None:
        """Set the log level for this instance."""
        self.log_config.level = level
