"""FastAPI server for Economic Indicators Dashboard API.

This server provides REST endpoints to access economic indicators data
stored in SQLite by the Go fetch service.
"""

import os
from typing import Optional

from fastapi import FastAPI, HTTPException, Query
from fastapi.middleware.cors import CORSMiddleware

from .db import EconDatabase
from . import handlers


# Initialize FastAPI app
app = FastAPI(
    title="Economic Indicators API",
    description="REST API for accessing economic indicators from FRED",
    version="1.0.0",
)

# Configure CORS for frontend access
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # In production, restrict to specific origins
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Database instance
db_path = os.environ.get("ECON_DB_PATH", "./data/econ.db")
db = EconDatabase(db_path)


@app.get("/")
async def root():
    """Root endpoint with API information."""
    return {
        "service": "Economic Indicators API",
        "version": "1.0.0",
        "endpoints": {
            "health": "/api/health",
            "summary": "/api/econ/summary",
            "series": "/api/econ/series?code=<CODE>&range=<RANGE>",
            "meta": "/api/econ/meta",
        }
    }


@app.get("/api/health")
async def health_check():
    """Health check endpoint.
    
    Returns:
        Status information about the service
    """
    try:
        # Check if we can read from database
        series = db.get_all_series()
        last_refresh = db.get_last_refresh()
        
        return {
            "status": "ok",
            "database": "connected",
            "series_count": len(series),
            "last_refresh": last_refresh.get("finished_at") if last_refresh else None,
        }
    except Exception as e:
        raise HTTPException(status_code=503, detail=f"Service unhealthy: {str(e)}")


@app.get("/api/econ/summary")
async def get_summary():
    """Get summary of all economic indicators.
    
    Returns latest value and MoM/YoY deltas for all configured indicators.
    
    Returns:
        Dictionary with list of indicators and metadata
    """
    try:
        return handlers.get_summary(db)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error fetching summary: {str(e)}")


@app.get("/api/econ/series")
async def get_series(
    code: str = Query(..., description="Series code (e.g., CPIAUCSL)"),
    range: Optional[str] = Query("max", description="Time range: 1y, 5y, or max"),
):
    """Get time series data for a specific indicator.
    
    Args:
        code: Series identifier (e.g., 'CPIAUCSL', 'UNRATE')
        range: Time range to fetch - '1y', '5y', or 'max' (default)
        
    Returns:
        Dictionary with series metadata and observations
    """
    try:
        result = handlers.get_series_data(db, code, range)
        
        if "error" in result:
            raise HTTPException(status_code=404, detail=result["error"])
        
        return result
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error fetching series: {str(e)}")


@app.get("/api/econ/meta")
async def get_metadata():
    """Get metadata about all available series.
    
    Returns:
        Dictionary with list of available series codes and names
    """
    try:
        return handlers.get_meta(db)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error fetching metadata: {str(e)}")


if __name__ == "__main__":
    import uvicorn
    
    port = int(os.environ.get("PORT", 8000))
    uvicorn.run(app, host="0.0.0.0", port=port)
