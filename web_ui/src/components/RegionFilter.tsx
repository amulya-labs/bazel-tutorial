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
    <div className="region-filter">
      <button
        className={`region-button ${selectedRegion === 'all' ? 'active' : ''}`}
        onClick={() => onRegionChange('all')}
      >
        <span className="region-label">All Indicators</span>
        <span className="region-count">{counts.all}</span>
      </button>
      <button
        className={`region-button ${selectedRegion === 'us' ? 'active' : ''}`}
        onClick={() => onRegionChange('us')}
      >
        <span className="region-label">US</span>
        <span className="region-count">{counts.us}</span>
      </button>
      <button
        className={`region-button ${selectedRegion === 'world' ? 'active' : ''}`}
        onClick={() => onRegionChange('world')}
      >
        <span className="region-label">World</span>
        <span className="region-count">{counts.world}</span>
      </button>
    </div>
  )
}

export default RegionFilter
