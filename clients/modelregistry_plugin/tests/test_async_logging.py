"""Tests for async logging functionality."""

import pytest
from unittest.mock import Mock, patch

from modelregistry_plugin.store_new import ModelRegistryStore


class TestAsyncLogging:
    """Test async logging methods."""

    def test_log_batch_async(self):
        """Test log_batch_async method."""
        store = ModelRegistryStore("modelregistry://localhost:8080", "file:///tmp")

        # Mock the async logging queue
        mock_queue = Mock()
        mock_queue.is_active.return_value = False
        mock_queue.activate.return_value = None
        mock_queue.log_batch_async.return_value = Mock()
        store._async_logging_queue = mock_queue

        # Test data
        from mlflow.entities import Metric, Param, RunTag

        metrics = [Metric(key="test_metric", value=1.0, timestamp=1234567890, step=1)]
        params = [Param(key="test_param", value="test_value")]
        tags = [RunTag(key="test_tag", value="test_value")]

        # Call the method
        result = store.log_batch_async("test_run_id", metrics, params, tags)

        # Verify the queue was activated and called
        mock_queue.activate.assert_called_once()
        mock_queue.log_batch_async.assert_called_once_with(
            run_id="test_run_id", metrics=metrics, params=params, tags=tags
        )
        assert result is not None

    def test_end_async_logging(self):
        """Test end_async_logging method."""
        store = ModelRegistryStore("modelregistry://localhost:8080", "file:///tmp")

        # Mock the async logging queue
        mock_queue = Mock()
        mock_queue.is_active.return_value = True
        mock_queue.end_async_logging.return_value = None
        store._async_logging_queue = mock_queue

        # Call the method
        store.end_async_logging()

        # Verify the queue was called
        mock_queue.end_async_logging.assert_called_once()

    def test_end_async_logging_inactive(self):
        """Test end_async_logging method when queue is inactive."""
        store = ModelRegistryStore("modelregistry://localhost:8080", "file:///tmp")

        # Mock the async logging queue
        mock_queue = Mock()
        mock_queue.is_active.return_value = False
        mock_queue.end_async_logging.return_value = None
        store._async_logging_queue = mock_queue

        # Call the method
        store.end_async_logging()

        # Verify the queue was not called
        mock_queue.end_async_logging.assert_not_called()

    def test_flush_async_logging(self):
        """Test flush_async_logging method."""
        store = ModelRegistryStore("modelregistry://localhost:8080", "file:///tmp")

        # Mock the async logging queue
        mock_queue = Mock()
        mock_queue.is_idle.return_value = False
        mock_queue.flush.return_value = None
        store._async_logging_queue = mock_queue

        # Call the method
        store.flush_async_logging()

        # Verify the queue was called
        mock_queue.flush.assert_called_once()

    def test_flush_async_logging_idle(self):
        """Test flush_async_logging method when queue is idle."""
        store = ModelRegistryStore("modelregistry://localhost:8080", "file:///tmp")

        # Mock the async logging queue
        mock_queue = Mock()
        mock_queue.is_idle.return_value = True
        mock_queue.flush.return_value = None
        store._async_logging_queue = mock_queue

        # Call the method
        store.flush_async_logging()

        # Verify the queue was not called
        mock_queue.flush.assert_not_called()

    def test_shut_down_async_logging(self):
        """Test shut_down_async_logging method."""
        store = ModelRegistryStore("modelregistry://localhost:8080", "file:///tmp")

        # Mock the async logging queue
        mock_queue = Mock()
        mock_queue.is_idle.return_value = False
        mock_queue.shut_down_async_logging.return_value = None
        store._async_logging_queue = mock_queue

        # Call the method
        store.shut_down_async_logging()

        # Verify the queue was called
        mock_queue.shut_down_async_logging.assert_called_once()

    def test_shut_down_async_logging_idle(self):
        """Test shut_down_async_logging method when queue is idle."""
        store = ModelRegistryStore("modelregistry://localhost:8080", "file:///tmp")

        # Mock the async logging queue
        mock_queue = Mock()
        mock_queue.is_idle.return_value = True
        mock_queue.shut_down_async_logging.return_value = None
        store._async_logging_queue = mock_queue

        # Call the method
        store.shut_down_async_logging()

        # Verify the queue was not called
        mock_queue.shut_down_async_logging.assert_not_called()
