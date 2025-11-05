"""Tests for the Economic Indicators API."""

import os
import sqlite3
import tempfile
from pathlib import Path

import pytest
from fastapi.testclient import TestClient

from py_api.main import app
from py_api.db import EconDatabase


@pytest.fixture
def test_db():
    """Create a test database with sample data."""
    # Create temporary database
    fd, db_path = tempfile.mkstemp(suffix=".db")
    os.close(fd)
    
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    # Create schema
    cursor.execute("""
        CREATE TABLE series (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            unit TEXT,
            source TEXT DEFAULT 'FRED'
        )
    """)
    
    cursor.execute("""
        CREATE TABLE observations (
            series_id TEXT NOT NULL,
            date TEXT NOT NULL,
            value REAL NOT NULL,
            PRIMARY KEY (series_id, date)
        )
    """)
    
    cursor.execute("""
        CREATE TABLE refresh_log (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            source TEXT NOT NULL,
            started_at TEXT NOT NULL,
            finished_at TEXT,
            ok INTEGER DEFAULT 0,
            message TEXT
        )
    """)
    
    # Insert test data
    cursor.execute(
        "INSERT INTO series (id, name, unit) VALUES (?, ?, ?)",
        ("CPIAUCSL", "Consumer Price Index", "Index 1982-1984=100")
    )
    
    cursor.execute(
        "INSERT INTO series (id, name, unit) VALUES (?, ?, ?)",
        ("UNRATE", "Unemployment Rate", "Percent")
    )
    
    # Insert observations (2 years of monthly data for testing)
    test_observations = [
        ("CPIAUCSL", "2022-01-01", 280.0),
        ("CPIAUCSL", "2022-06-01", 290.0),
        ("CPIAUCSL", "2023-01-01", 295.0),
        ("CPIAUCSL", "2023-06-01", 300.0),
        ("UNRATE", "2022-01-01", 4.0),
        ("UNRATE", "2022-06-01", 3.8),
        ("UNRATE", "2023-01-01", 3.6),
        ("UNRATE", "2023-06-01", 3.5),
    ]
    
    for series_id, date, value in test_observations:
        cursor.execute(
            "INSERT INTO observations (series_id, date, value) VALUES (?, ?, ?)",
            (series_id, date, value)
        )
    
    # Insert refresh log
    cursor.execute(
        "INSERT INTO refresh_log (source, started_at, finished_at, ok, message) VALUES (?, ?, ?, ?, ?)",
        ("FRED", "2023-06-01T12:00:00", "2023-06-01T12:05:00", 1, "Success")
    )
    
    conn.commit()
    conn.close()
    
    yield db_path
    
    # Cleanup
    Path(db_path).unlink(missing_ok=True)


def test_database_get_all_series(test_db):
    """Test retrieving all series from database."""
    db = EconDatabase(test_db)
    series = db.get_all_series()
    
    assert len(series) == 2
    assert series[0]["id"] == "CPIAUCSL"
    assert series[1]["id"] == "UNRATE"


def test_database_get_latest_observation(test_db):
    """Test retrieving latest observation."""
    db = EconDatabase(test_db)
    date, value = db.get_latest_observation("CPIAUCSL")
    
    assert date == "2023-06-01"
    assert value == 300.0


def test_database_calculate_delta_mom(test_db):
    """Test month-over-month delta calculation."""
    db = EconDatabase(test_db)
    delta = db.calculate_delta_mom("CPIAUCSL")
    
    # (300 - 295) / 295 * 100 ≈ 1.69
    assert delta is not None
    assert abs(delta - 1.69) < 0.1


def test_database_calculate_delta_yoy(test_db):
    """Test year-over-year delta calculation."""
    db = EconDatabase(test_db)
    delta = db.calculate_delta_yoy("CPIAUCSL")
    
    # (300 - 290) / 290 * 100 ≈ 3.45
    assert delta is not None
    assert abs(delta - 3.45) < 0.1


def test_api_health_check(test_db, monkeypatch):
    """Test health check endpoint."""
    monkeypatch.setattr("py_api.main.db", EconDatabase(test_db))
    
    client = TestClient(app)
    response = client.get("/api/health")
    
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["series_count"] == 2


def test_api_get_summary(test_db, monkeypatch):
    """Test summary endpoint."""
    monkeypatch.setattr("py_api.main.db", EconDatabase(test_db))
    
    client = TestClient(app)
    response = client.get("/api/econ/summary")
    
    assert response.status_code == 200
    data = response.json()
    assert data["count"] == 2
    assert len(data["indicators"]) == 2
    
    # Check first indicator
    cpi = next(i for i in data["indicators"] if i["code"] == "CPIAUCSL")
    assert cpi["name"] == "Consumer Price Index"
    assert cpi["value"] == 300.0
    assert cpi["delta_mom"] is not None
    assert cpi["delta_yoy"] is not None


def test_api_get_series(test_db, monkeypatch):
    """Test series endpoint."""
    monkeypatch.setattr("py_api.main.db", EconDatabase(test_db))
    
    client = TestClient(app)
    response = client.get("/api/econ/series?code=UNRATE&range=max")
    
    assert response.status_code == 200
    data = response.json()
    assert data["code"] == "UNRATE"
    assert data["name"] == "Unemployment Rate"
    assert len(data["observations"]) == 4
    assert data["latest_value"] == 3.5


def test_api_get_series_not_found(test_db, monkeypatch):
    """Test series endpoint with invalid code."""
    monkeypatch.setattr("py_api.main.db", EconDatabase(test_db))
    
    client = TestClient(app)
    response = client.get("/api/econ/series?code=INVALID")
    
    assert response.status_code == 404


def test_api_get_meta(test_db, monkeypatch):
    """Test metadata endpoint."""
    monkeypatch.setattr("py_api.main.db", EconDatabase(test_db))
    
    client = TestClient(app)
    response = client.get("/api/econ/meta")
    
    assert response.status_code == 200
    data = response.json()
    assert data["count"] == 2
    assert len(data["series"]) == 2
