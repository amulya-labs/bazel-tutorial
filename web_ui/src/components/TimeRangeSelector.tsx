import { TimeRange } from '../types'
import './TimeRangeSelector.css'

interface TimeRangeSelectorProps {
  selectedRange: TimeRange
  onRangeChange: (range: TimeRange) => void
}

const TIME_RANGES: TimeRange[] = ['1M', '3M', '6M', '1Y', '5Y', 'All']

function TimeRangeSelector({ selectedRange, onRangeChange }: TimeRangeSelectorProps) {
  return (
    <div className="time-range-selector">
      {TIME_RANGES.map((range) => (
        <button
          key={range}
          className={`time-range-button ${selectedRange === range ? 'active' : ''}`}
          onClick={() => onRangeChange(range)}
        >
          {range}
        </button>
      ))}
    </div>
  )
}

export default TimeRangeSelector
