import { useEffect, useState } from 'react'
import IndicatorCard from './components/IndicatorCard'
import LineChart from './components/LineChart'
import './App.css'

interface Indicator {
  code: string
  name: string
  unit: string
  source: string
  last_updated: string
  value: number
  delta_mom: number | null
  delta_yoy: number | null
}

interface SummaryResponse {
  indicators: Indicator[]
  count: number
  last_refresh: string | null
}

interface SeriesData {
  code: string
  name: string
  unit: string
  observations: Array<{ date: string; value: number }>
  latest_value: number
  delta_mom: number | null
  delta_yoy: number | null
}

function App() {
  const [summary, setSummary] = useState<SummaryResponse | null>(null)
  const [selectedSeries, setSelectedSeries] = useState<SeriesData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchSummary()
  }, [])

  const fetchSummary = async () => {
    try {
      setLoading(true)

      // Try static JSON first (for GitHub Pages deployment)
      // Use relative path to respect Vite's base configuration
      let response = await fetch('data/summary.json')

      // Fall back to API if static file not found
      if (!response.ok && response.status === 404) {
        response = await fetch('/api/econ/summary')
      }

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data = await response.json()
      setSummary(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch data')
    } finally {
      setLoading(false)
    }
  }

  const fetchSeriesData = async (code: string, range: string = '5y') => {
    try {
      // Try static JSON first (for GitHub Pages deployment)
      // Use relative path to respect Vite's base configuration
      let response = await fetch(`data/series-${code}.json`)

      // Fall back to API if static file not found
      if (!response.ok && response.status === 404) {
        response = await fetch(`/api/econ/series?code=${code}&range=${range}`)
      }

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data = await response.json()
      setSelectedSeries(data)
    } catch (err) {
      console.error('Failed to fetch series data:', err)
    }
  }

  const handleCardClick = (code: string) => {
    fetchSeriesData(code)
  }

  const handleBackClick = () => {
    setSelectedSeries(null)
  }

  if (loading) {
    return (
      <div className="app">
        <div className="loading">Loading economic data...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="app">
        <div className="error">
          <h2>Error</h2>
          <p>{error}</p>
          <p>Make sure to run: <code>bazel run //go_fetch:refresh</code> first to populate the database.</p>
          <p>Then start the API server: <code>bazel run //py_api:server</code></p>
        </div>
      </div>
    )
  }

  if (selectedSeries) {
    return (
      <div className="app">
        <header className="header">
          <div className="header-content">
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '1rem' }}>
              <button onClick={handleBackClick} className="back-button">
                ← Back to Dashboard
              </button>
              <a
                href="/bazel-tutorial/"
                style={{
                  padding: '0.5rem 1rem',
                  background: 'rgba(255,255,255,0.2)',
                  color: 'white',
                  textDecoration: 'none',
                  borderRadius: '4px',
                  fontSize: '0.9rem',
                  fontWeight: 500,
                  border: '1px solid rgba(255,255,255,0.3)'
                }}
              >
                ← Back to Docs
              </a>
            </div>
            <h1>{selectedSeries.name}</h1>
          </div>
        </header>
        <main className="main">
          <div className="series-details">
            <div className="series-info">
              <div className="info-item">
                <span className="label">Code:</span>
                <span className="value">{selectedSeries.code}</span>
              </div>
              <div className="info-item">
                <span className="label">Unit:</span>
                <span className="value">{selectedSeries.unit}</span>
              </div>
              <div className="info-item">
                <span className="label">Latest Value:</span>
                <span className="value">{selectedSeries.latest_value.toFixed(2)}</span>
              </div>
              <div className="info-item">
                <span className="label">MoM Change:</span>
                <span className={`value ${selectedSeries.delta_mom !== null && selectedSeries.delta_mom !== 0 ? (selectedSeries.delta_mom > 0 ? 'positive' : 'negative') : 'neutral'}`}>
                  {selectedSeries.delta_mom !== null ? `${selectedSeries.delta_mom > 0 ? '+' : ''}${selectedSeries.delta_mom.toFixed(2)}%` : 'N/A'}
                </span>
              </div>
              <div className="info-item">
                <span className="label">YoY Change:</span>
                <span className={`value ${selectedSeries.delta_yoy !== null && selectedSeries.delta_yoy !== 0 ? (selectedSeries.delta_yoy > 0 ? 'positive' : 'negative') : 'neutral'}`}>
                  {selectedSeries.delta_yoy !== null ? `${selectedSeries.delta_yoy > 0 ? '+' : ''}${selectedSeries.delta_yoy.toFixed(2)}%` : 'N/A'}
                </span>
              </div>
            </div>
            <div className="chart-container">
              <LineChart data={selectedSeries.observations} />
            </div>
          </div>
        </main>
      </div>
    )
  }

  return (
    <div className="app">
      <header className="header">
        <div className="header-content">
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '1rem' }}>
            <h1 style={{ margin: 0 }}>📊 Economic Indicators Dashboard</h1>
            <a
              href="/bazel-tutorial/dashboard/"
              style={{
                padding: '0.5rem 1rem',
                background: 'rgba(255,255,255,0.2)',
                color: 'white',
                textDecoration: 'none',
                borderRadius: '4px',
                fontSize: '0.9rem',
                fontWeight: 500,
                border: '1px solid rgba(255,255,255,0.3)'
              }}
            >
              ← Back to Docs
            </a>
          </div>
          <p className="subtitle">Real-time economic data from FRED</p>
          {summary?.last_refresh && (
            <p className="refresh-time">
              Last updated: {new Date(summary.last_refresh).toLocaleString()}
            </p>
          )}
        </div>
      </header>
      <main className="main">
        {summary && summary.indicators.length === 0 ? (
          <div className="empty-state">
            <h2>No Data Available</h2>
            <p>The database is empty. Please populate it with economic data:</p>
            <ol style={{ textAlign: 'left', margin: '1rem auto', maxWidth: '500px' }}>
              <li>Get a free FRED API key at <a href="https://fred.stlouisfed.org/docs/api/api_key.html" target="_blank" rel="noopener noreferrer">https://fred.stlouisfed.org</a></li>
              <li>Set the environment variable: <code>export FRED_API_KEY=your_key_here</code></li>
              <li>Run: <code>bazel run //go_fetch:refresh</code></li>
              <li>Refresh this page</li>
            </ol>
          </div>
        ) : (
          <div className="dashboard-grid">
            {summary?.indicators.map((indicator) => (
              <IndicatorCard
                key={indicator.code}
                indicator={indicator}
                onClick={() => handleCardClick(indicator.code)}
              />
            ))}
          </div>
        )}
      </main>
    </div>
  )
}

export default App
