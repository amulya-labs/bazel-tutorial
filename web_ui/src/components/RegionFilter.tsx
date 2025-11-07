import { RegionFilter as RegionFilterType } from '../utils/indicatorUtils'
import './RegionFilter.css'

interface RegionFilterProps {
  selectedRegion: RegionFilterType
  onRegionChange: (region: RegionFilterType) => void
  counts: {
    all: number
    us: number
    world: number
  }
}

function RegionFilter({ selectedRegion, onRegionChange, counts }: RegionFilterProps) {
  return (
    <div className="region-filter" role="group" aria-label="Region filter">
      <button
        className={`region-button ${selectedRegion === 'all' ? 'active' : ''}`}
        onClick={() => onRegionChange('all')}
        aria-pressed={selectedRegion === 'all'}
        aria-label={`Show all indicators (${counts.all})`}
      >
        <span className="region-label">All Indicators</span>
        <span className="region-count">{counts.all}</span>
      </button>
      <button
        className={`region-button ${selectedRegion === 'us' ? 'active' : ''}`}
        onClick={() => onRegionChange('us')}
        aria-pressed={selectedRegion === 'us'}
        aria-label={`Show US indicators (${counts.us})`}
      >
        <span className="region-label">US</span>
        <span className="region-count">{counts.us}</span>
      </button>
      <button
        className={`region-button ${selectedRegion === 'world' ? 'active' : ''}`}
        onClick={() => onRegionChange('world')}
        aria-pressed={selectedRegion === 'world'}
        aria-label={`Show world indicators (${counts.world})`}
      >
        <span className="region-label">World</span>
        <span className="region-count">{counts.world}</span>
      </button>
    </div>
  )
}

export default RegionFilter
