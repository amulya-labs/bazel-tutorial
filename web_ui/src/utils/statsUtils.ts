/**
 * Utility functions for calculating statistics on time series data
 */

import { Observation, TimeRange, TimeRangeStats } from '../types'

/**
 * Filters observations based on the selected time range
 */
export function filterObservationsByTimeRange(
  observations: Observation[],
  timeRange: TimeRange
): Observation[] {
  if (timeRange === 'All' || observations.length === 0) {
    return observations
  }

  const now = new Date()
  const cutoffDate = new Date()

  switch (timeRange) {
    case '1M':
      cutoffDate.setMonth(now.getMonth() - 1)
      break
    case '3M':
      cutoffDate.setMonth(now.getMonth() - 3)
      break
    case '6M':
      cutoffDate.setMonth(now.getMonth() - 6)
      break
    case '1Y':
      cutoffDate.setFullYear(now.getFullYear() - 1)
      break
    case '5Y':
      cutoffDate.setFullYear(now.getFullYear() - 5)
      break
  }

  return observations.filter(obs => new Date(obs.date) >= cutoffDate)
}

/**
 * Calculates statistics for a given set of observations
 */
export function calculateStats(observations: Observation[]): TimeRangeStats | null {
  if (observations.length === 0) {
    return null
  }

  const values = observations.map(obs => obs.value)
  const min = Math.min(...values)
  const max = Math.max(...values)
  const avg = values.reduce((sum, val) => sum + val, 0) / values.length

  const startValue = observations[0].value
  const endValue = observations[observations.length - 1].value
  const changeAbsolute = endValue - startValue
  const change = startValue !== 0 ? (changeAbsolute / startValue) * 100 : 0

  return {
    min,
    max,
    avg,
    change,
    changeAbsolute,
    startValue,
    endValue,
    dataPoints: observations.length,
  }
}

/**
 * Formats a number for display with appropriate precision
 */
export function formatStatValue(value: number, decimals: number = 2): string {
  if (Math.abs(value) >= 1000000) {
    return (value / 1000000).toFixed(1) + 'M'
  } else if (Math.abs(value) >= 1000) {
    return (value / 1000).toFixed(1) + 'K'
  }
  return value.toFixed(decimals)
}

/**
 * Formats percentage change with sign
 */
export function formatPercentChange(value: number): string {
  const sign = value > 0 ? '+' : ''
  return `${sign}${value.toFixed(2)}%`
}
