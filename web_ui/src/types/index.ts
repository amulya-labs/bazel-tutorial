/**
 * Core type definitions for the Economic Dashboard
 */

export interface Indicator {
  code: string
  name: string
  unit: string
  source: string
  last_updated: string
  value: number
  delta_mom: number | null
  delta_yoy: number | null
}

export interface SummaryResponse {
  indicators: Indicator[]
  count: number
  last_refresh: string | null
}

export interface Observation {
  date: string
  value: number
}

export interface SeriesData {
  code: string
  name: string
  unit: string
  observations: Observation[]
  latest_value: number
  delta_mom: number | null
  delta_yoy: number | null
}

export interface TimeRangeStats {
  min: number
  max: number
  avg: number
  change: number // Percentage change from start to end
  changeAbsolute: number // Absolute change from start to end
  startValue: number
  endValue: number
  dataPoints: number
}

export type TimeRange = '1M' | '3M' | '6M' | '1Y' | '5Y' | 'All'
