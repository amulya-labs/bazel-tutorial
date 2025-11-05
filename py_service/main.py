"""Simple Python HTTP service using FastAPI."""

from fastapi import FastAPI
import uvicorn

app = FastAPI()


@app.get("/")
def read_root():
    """Return a simple greeting message."""
    return {"message": "Hello from Python! 🐍"}


@app.get("/health")
def health_check():
    """Health check endpoint."""
    return {"status": "OK"}


if __name__ == "__main__":
    # Run the service on port 8081
    print("Python service starting on http://localhost:8081")
    uvicorn.run(app, host="0.0.0.0", port=8081)
