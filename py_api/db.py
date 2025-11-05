"""Database module for reading economic data from SQLite.

This module provides read-only access to the economic indicators database
that is populated by the Go fetch service.
"""

import sqlite3
from datetime import datetime, timedelta
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple


class EconDatabase:
    """Read-only interface to the economic indicators SQLite database."""

    def __init__(self, db_path: str = "./data/econ.db"):
        """Initialize database connection.
        
        Args:
            db_path: Path to the SQLite database file
        """
        self.db_path = db_path
        self._ensure_db_exists()

    def _ensure_db_exists(self):
        """Ensure the database file exists."""
        db_file = Path(self.db_path)
        if not db_file.exists():
            raise FileNotFoundError(
                f"Database not found at {self.db_path}. "
                "Please run 'bazel run //go_fetch:refresh' first."
            )

    def _get_connection(self) -> sqlite3.Connection:
        """Get a database connection."""
        conn = sqlite3.connect(self.db_path)
        conn.row_factory = sqlite3.Row
        return conn

    def get_all_series(self) -> List[Dict[str, str]]:
        """Get metadata for all available series.
        
        Returns:
            List of dictionaries containing series metadata
        """
        with self._get_connection() as conn:
            cursor = conn.execute(
                "SELECT id, name, unit, source FROM series ORDER BY id"
            )
            return [dict(row) for row in cursor.fetchall()]

    def get_series_metadata(self, series_id: str) -> Optional[Dict[str, str]]:
        """Get metadata for a specific series.
        
        Args:
            series_id: The series identifier (e.g., 'CPIAUCSL')
            
        Returns:
            Dictionary with series metadata or None if not found
        """
        with self._get_connection() as conn:
            cursor = conn.execute(
                "SELECT id, name, unit, source FROM series WHERE id = ?",
                (series_id,)
            )
            row = cursor.fetchone()
            return dict(row) if row else None

    def get_latest_observation(self, series_id: str) -> Optional[Tuple[str, float]]:
        """Get the most recent observation for a series.
        
        Args:
            series_id: The series identifier
            
        Returns:
            Tuple of (date, value) or None if not found
        """
        with self._get_connection() as conn:
            cursor = conn.execute(
                """
                SELECT date, value 
                FROM observations 
                WHERE series_id = ?
                ORDER BY date DESC
                LIMIT 1
                """,
                (series_id,)
            )
            row = cursor.fetchone()
            return (row["date"], row["value"]) if row else None

    def get_observation_at_date(self, series_id: str, target_date: str) -> Optional[float]:
        """Get the observation value closest to a target date.
        
        Args:
            series_id: The series identifier
            target_date: Target date in YYYY-MM-DD format
            
        Returns:
            The value at or before the target date, or None
        """
        with self._get_connection() as conn:
            cursor = conn.execute(
                """
                SELECT value
                FROM observations
                WHERE series_id = ? AND date <= ?
                ORDER BY date DESC
                LIMIT 1
                """,
                (series_id, target_date)
            )
            row = cursor.fetchone()
            return row["value"] if row else None

    def get_observations(
        self, 
        series_id: str, 
        start_date: Optional[str] = None,
        end_date: Optional[str] = None
    ) -> List[Dict[str, Any]]:
        """Get observations for a series within a date range.
        
        Args:
            series_id: The series identifier
            start_date: Start date in YYYY-MM-DD format (optional)
            end_date: End date in YYYY-MM-DD format (optional)
            
        Returns:
            List of observations with date and value
        """
        query = "SELECT date, value FROM observations WHERE series_id = ?"
        params = [series_id]

        if start_date:
            query += " AND date >= ?"
            params.append(start_date)
        
        if end_date:
            query += " AND date <= ?"
            params.append(end_date)
        
        query += " ORDER BY date"

        with self._get_connection() as conn:
            cursor = conn.execute(query, params)
            return [{"date": row["date"], "value": row["value"]} for row in cursor.fetchall()]

    def calculate_delta_mom(self, series_id: str) -> Optional[float]:
        """Calculate month-over-month percentage change.
        
        Args:
            series_id: The series identifier
            
        Returns:
            Percentage change or None if insufficient data
        """
        with self._get_connection() as conn:
            cursor = conn.execute(
                """
                SELECT value 
                FROM observations 
                WHERE series_id = ?
                ORDER BY date DESC
                LIMIT 2
                """,
                (series_id,)
            )
            rows = cursor.fetchall()
            
            if len(rows) < 2:
                return None
            
            current = rows[0]["value"]
            previous = rows[1]["value"]
            
            if previous == 0:
                return None
            
            return ((current - previous) / previous) * 100

    def calculate_delta_yoy(self, series_id: str) -> Optional[float]:
        """Calculate year-over-year percentage change.
        
        Args:
            series_id: The series identifier
            
        Returns:
            Percentage change or None if insufficient data
        """
        latest = self.get_latest_observation(series_id)
        if not latest:
            return None
        
        latest_date, latest_value = latest
        
        # Calculate date one year ago
        try:
            date_obj = datetime.strptime(latest_date, "%Y-%m-%d")
            year_ago = (date_obj - timedelta(days=365)).strftime("%Y-%m-%d")
        except ValueError:
            return None
        
        year_ago_value = self.get_observation_at_date(series_id, year_ago)
        
        if year_ago_value is None or year_ago_value == 0:
            return None
        
        return ((latest_value - year_ago_value) / year_ago_value) * 100

    def get_last_refresh(self) -> Optional[Dict[str, Any]]:
        """Get information about the last data refresh.
        
        Returns:
            Dictionary with refresh metadata or None
        """
        with self._get_connection() as conn:
            cursor = conn.execute(
                """
                SELECT source, started_at, finished_at, ok, message
                FROM refresh_log
                ORDER BY started_at DESC
                LIMIT 1
                """
            )
            row = cursor.fetchone()
            return dict(row) if row else None
