import './IndicatorCard.css'

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

interface IndicatorCardProps {
  indicator: Indicator
  onClick: () => void
}

function IndicatorCard({ indicator, onClick }: IndicatorCardProps) {
  const formatValue = (value: number): string => {
    if (value >= 1000) {
      return value.toFixed(1)
    }
    return value.toFixed(2)
  }

  const formatDelta = (delta: number | null): string => {
    if (delta === null) return 'N/A'
    const sign = delta > 0 ? '+' : ''
    return `${sign}${delta.toFixed(2)}%`
  }

  const getDeltaClass = (delta: number | null): string => {
    if (delta === null) return 'neutral'
    return delta > 0 ? 'positive' : 'negative'
  }

  return (
    <div className="indicator-card" onClick={onClick}>
      <div className="card-header">
        <h3 className="card-title">{indicator.name}</h3>
        <span className="card-code">{indicator.code}</span>
      </div>
      
      <div className="card-value">
        <span className="value">{formatValue(indicator.value)}</span>
        <span className="unit">{indicator.unit}</span>
      </div>

      <div className="card-deltas">
        <div className="delta-item">
          <span className="delta-label">MoM</span>
          <span className={`delta-value ${getDeltaClass(indicator.delta_mom)}`}>
            {formatDelta(indicator.delta_mom)}
          </span>
        </div>
        <div className="delta-item">
          <span className="delta-label">YoY</span>
          <span className={`delta-value ${getDeltaClass(indicator.delta_yoy)}`}>
            {formatDelta(indicator.delta_yoy)}
          </span>
        </div>
      </div>

      <div className="card-footer">
        <span className="last-updated">
          Updated: {new Date(indicator.last_updated).toLocaleDateString()}
        </span>
      </div>
    </div>
  )
}

export default IndicatorCard
