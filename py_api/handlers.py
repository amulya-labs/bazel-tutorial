"""API handlers for economic indicators endpoints."""

from typing import Dict, Optional
from datetime import datetime, timedelta

from .db import EconDatabase


def get_summary(db: EconDatabase) -> Dict:
    """Get summary of all indicators with latest values and deltas.
    
    Args:
        db: Database instance
        
    Returns:
        Dictionary containing indicators list and metadata
    """
    series_list = db.get_all_series()
    indicators = []
    
    for series in series_list:
        series_id = series["id"]
        
        # Get latest observation
        latest = db.get_latest_observation(series_id)
        if not latest:
            continue
        
        date, value = latest
        
        # Calculate deltas
        delta_mom = db.calculate_delta_mom(series_id)
        delta_yoy = db.calculate_delta_yoy(series_id)
        
        indicators.append({
            "code": series_id,
            "name": series["name"],
            "unit": series["unit"],
            "source": series["source"],
            "last_updated": date,
            "value": value,
            "delta_mom": round(delta_mom, 2) if delta_mom is not None else None,
            "delta_yoy": round(delta_yoy, 2) if delta_yoy is not None else None,
        })
    
    # Get last refresh info
    refresh_info = db.get_last_refresh()
    
    return {
        "indicators": indicators,
        "count": len(indicators),
        "last_refresh": refresh_info.get("finished_at") if refresh_info else None,
    }


def get_series_data(
    db: EconDatabase,
    code: str,
    range_param: Optional[str] = None
) -> Dict:
    """Get time series data for a specific indicator.
    
    Args:
        db: Database instance
        code: Series code (e.g., 'CPIAUCSL')
        range_param: Time range - '1y', '5y', or 'max' (default: 'max')
        
    Returns:
        Dictionary containing series metadata and observations
    """
    # Get series metadata
    metadata = db.get_series_metadata(code)
    if not metadata:
        return {"error": f"Series {code} not found"}
    
    # Calculate date range
    end_date = None
    start_date = None
    
    if range_param and range_param != "max":
        latest = db.get_latest_observation(code)
        if latest:
            latest_date_str = latest[0]
            try:
                latest_date = datetime.strptime(latest_date_str, "%Y-%m-%d")
                
                if range_param == "1y":
                    start_date = (latest_date - timedelta(days=365)).strftime("%Y-%m-%d")
                elif range_param == "5y":
                    start_date = (latest_date - timedelta(days=365*5)).strftime("%Y-%m-%d")
            except ValueError:
                pass  # Use max range if date parsing fails
    
    # Get observations
    observations = db.get_observations(code, start_date, end_date)
    
    # Get latest value and deltas
    latest = db.get_latest_observation(code)
    delta_mom = db.calculate_delta_mom(code)
    delta_yoy = db.calculate_delta_yoy(code)
    
    return {
        "code": metadata["id"],
        "name": metadata["name"],
        "unit": metadata["unit"],
        "source": metadata["source"],
        "last_updated": latest[0] if latest else None,
        "latest_value": latest[1] if latest else None,
        "delta_mom": round(delta_mom, 2) if delta_mom is not None else None,
        "delta_yoy": round(delta_yoy, 2) if delta_yoy is not None else None,
        "observations": observations,
        "count": len(observations),
        "range": range_param or "max",
    }


def get_meta(db: EconDatabase) -> Dict:
    """Get metadata about all available series.
    
    Args:
        db: Database instance
        
    Returns:
        Dictionary containing list of available series
    """
    series_list = db.get_all_series()
    
    return {
        "series": [
            {
                "code": s["id"],
                "name": s["name"],
                "unit": s["unit"],
                "source": s["source"],
            }
            for s in series_list
        ],
        "count": len(series_list),
    }
