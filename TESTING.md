# 🧪 Testing Guide - Economic Dashboard

This guide explains the testing strategy for the Economic Dashboard and how to run tests.

## Overview

All tests are **hermetic** - they don't require:
- ❌ Network access
- ❌ API keys
- ❌ External databases
- ❌ Running services

Tests use **mocks** (Go) and **fixtures** (Python) for complete isolation.

## Running Tests

### Run All Tests
```bash
bazel test //...
```

### Run Specific Tests
```bash
# Go tests
bazel test //go_fetch:fred_test

# Python tests
bazel test //py_api:api_test

# Legacy services
bazel test //go_service:go_service_test
bazel test //py_service:py_service_test
```

### Verbose Output
```bash
# Show all output
bazel test //... --test_output=all

# Show only errors
bazel test //... --test_output=errors

# Show test summary
bazel test //... --test_summary=detailed
```

## Go Tests (`go_fetch/`)

### What's Tested
- ✅ FRED API client (series metadata fetching)
- ✅ FRED API client (observations fetching)
- ✅ Error handling for invalid responses
- ✅ Data filtering (skipping "." values)

### Testing Strategy: Mock HTTP Client

The `fred_test.go` uses a **mock HTTP client** that returns predefined responses:

```go
type mockHTTPClient struct {
    responses map[string]*http.Response
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
    if resp, ok := m.responses[url]; ok {
        return resp, nil
    }
    return &http.Response{StatusCode: 404}, nil
}
```

This allows testing without making real API calls.

### Example Test
```go
func TestFREDClient_FetchSeriesMetadata(t *testing.T) {
    mockClient := &mockHTTPClient{
        responses: map[string]*http.Response{
            "https://api.stlouisfed.org/...": newMockResponse(200, `{
                "seriess": [{
                    "id": "CPIAUCSL",
                    "title": "Consumer Price Index",
                    "units": "Index 1982-1984=100"
                }]
            }`),
        },
    }
    
    client := NewFREDClientWithHTTP("test", mockClient)
    series, err := client.fetchSeriesMetadata("CPIAUCSL")
    
    // Assertions...
}
```

### Run Go Tests
```bash
bazel test //go_fetch:fred_test --test_output=all
```

## Python Tests (`py_api/`)

### What's Tested
- ✅ Database operations (series, observations)
- ✅ Delta calculations (MoM, YoY)
- ✅ API endpoints (health, summary, series, meta)
- ✅ Error handling (missing series, invalid data)

### Testing Strategy: Fixture Database

The `test_api.py` creates a **temporary SQLite database** with test data:

```python
@pytest.fixture
def test_db():
    """Create a test database with sample data."""
    fd, db_path = tempfile.mkstemp(suffix=".db")
    os.close(fd)
    
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    # Create schema
    cursor.execute("CREATE TABLE series (...)")
    cursor.execute("CREATE TABLE observations (...)")
    
    # Insert test data
    cursor.execute("INSERT INTO series VALUES (...)")
    cursor.execute("INSERT INTO observations VALUES (...)")
    
    conn.commit()
    conn.close()
    
    yield db_path
    
    # Cleanup
    Path(db_path).unlink(missing_ok=True)
```

### Example Test
```python
def test_api_get_summary(test_db, monkeypatch):
    """Test summary endpoint."""
    monkeypatch.setattr("py_api.main.db", EconDatabase(test_db))
    
    client = TestClient(app)
    response = client.get("/api/econ/summary")
    
    assert response.status_code == 200
    data = response.json()
    assert data["count"] == 2
    assert len(data["indicators"]) == 2
```

### Run Python Tests
```bash
bazel test //py_api:api_test --test_output=all
```

## Test Coverage

### Go Tests Coverage
- **fred.go**: FRED API client functions
  - `FetchSeries()` - Integration test
  - `fetchSeriesMetadata()` - Metadata fetching
  - `fetchObservations()` - Observations fetching
  - Error handling for bad responses

- **Not tested** (integration-level):
  - `main.go` - CLI entry point
  - `store_sqlite.go` - SQLite operations
  
  These are tested manually via `bazel run //go_fetch:refresh`

### Python Tests Coverage
- **db.py**: Database layer
  - `get_all_series()` - All series
  - `get_latest_observation()` - Latest value
  - `calculate_delta_mom()` - MoM delta
  - `calculate_delta_yoy()` - YoY delta

- **handlers.py**: API handlers
  - `get_summary()` - Summary logic
  - `get_series_data()` - Series with range

- **main.py**: FastAPI endpoints
  - `GET /api/health`
  - `GET /api/econ/summary`
  - `GET /api/econ/series`
  - `GET /api/econ/meta`

## Adding New Tests

### Adding Go Tests

1. Create test file: `go_fetch/myfeature_test.go`
2. Write test with mock HTTP:
   ```go
   func TestMyFeature(t *testing.T) {
       // Setup mock
       // Call function
       // Assert results
   }
   ```
3. Add to `BUILD.bazel`:
   ```python
   go_test(
       name = "myfeature_test",
       srcs = ["myfeature_test.go"],
       embed = [":go_fetch_lib"],
   )
   ```

### Adding Python Tests

1. Add test function to `py_api/test_api.py`:
   ```python
   def test_my_feature(test_db, monkeypatch):
       # Setup
       monkeypatch.setattr("py_api.main.db", EconDatabase(test_db))
       
       # Test
       # Assert
   ```
2. Tests are automatically discovered by pytest

## Debugging Tests

### Show Full Output
```bash
bazel test //py_api:api_test --test_output=all
```

### Run Tests Directly (Outside Bazel)
```bash
# Go tests
cd go_fetch
go test -v ./...

# Python tests
cd py_api
python -m pytest test_api.py -v
```

### Common Issues

**"No module named 'py_api'"**
- **Solution**: Run via Bazel which sets up paths correctly
- Or: `PYTHONPATH=. pytest py_api/test_api.py`

**"Database not found"**
- **Solution**: Tests create their own fixtures, check the test_db fixture

**"Import error in Go tests"**
- **Solution**: Make sure `embed = [":go_fetch_lib"]` in BUILD.bazel

## CI/CD Integration

Tests run automatically in GitHub Actions:

```yaml
- name: Run all tests
  run: bazel test //...

- name: Run Go fetch tests
  run: bazel test //go_fetch:fred_test

- name: Run Python API tests
  run: bazel test //py_api:api_test
```

See `.github/workflows/ci.yaml` for full configuration.

## Best Practices

1. ✅ **Keep tests hermetic** - No network, no external state
2. ✅ **Use fixtures/mocks** - Fast and reliable
3. ✅ **Test business logic** - Not infrastructure
4. ✅ **Name tests descriptively** - `test_api_get_summary_returns_indicators`
5. ✅ **Clean up resources** - Use `defer` (Go) or fixtures (Python)

## Further Reading

- [Bazel Testing Docs](https://bazel.build/reference/test-encyclopedia)
- [Go Testing Package](https://pkg.go.dev/testing)
- [Pytest Documentation](https://docs.pytest.org/)
- [FastAPI Testing](https://fastapi.tiangolo.com/tutorial/testing/)
