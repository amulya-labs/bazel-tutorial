import { TimeRangeStats } from '../types'
import { formatStatValue, formatPercentChange } from '../utils/statsUtils'
import './StatsDisplay.css'

interface StatsDisplayProps {
  stats: TimeRangeStats | null
  unit: string
}

function StatsDisplay({ stats, unit }: StatsDisplayProps) {
  if (!stats) {
    return (
      <div className="stats-display">
        <div className="stats-message">No data available for this time range</div>
      </div>
    )
  }

  return (
    <div className="stats-display">
      <div className="stat-item">
        <span className="stat-label">Min</span>
        <span className="stat-value">{formatStatValue(stats.min)} <span className="stat-unit">{unit}</span></span>
      </div>
      <div className="stat-item">
        <span className="stat-label">Max</span>
        <span className="stat-value">{formatStatValue(stats.max)} <span className="stat-unit">{unit}</span></span>
      </div>
      <div className="stat-item">
        <span className="stat-label">Average</span>
        <span className="stat-value">{formatStatValue(stats.avg)} <span className="stat-unit">{unit}</span></span>
      </div>
      <div className="stat-item">
        <span className="stat-label">Change</span>
        <span className={`stat-value ${stats.change !== null ? (stats.change > 0 ? 'positive' : stats.change < 0 ? 'negative' : 'neutral') : 'neutral'}`}>
          {stats.change !== null ? formatPercentChange(stats.change) : 'N/A'}
        </span>
      </div>
      <div className="stat-item">
        <span className="stat-label">Data Points</span>
        <span className="stat-value">{stats.dataPoints}</span>
      </div>
    </div>
  )
}

export default StatsDisplay
