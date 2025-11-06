"""Unit tests for the Python service."""

from fastapi.testclient import TestClient
from py_service.main import app

client = TestClient(app)


def test_read_root():
    """Test the root endpoint returns the correct greeting."""
    response = client.get("/")
    assert response.status_code == 200
    assert response.json() == {"message": "Hello from Python! 🐍"}


def test_health_check():
    """Test the health check endpoint."""
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "OK"}
