/**
 * Utility functions for calculating statistics on time series data
 */

import { Observation, TimeRange, TimeRangeStats } from '../types'

/**
 * Filters observations based on the selected time range
 * Calculates the cutoff date from the most recent observation, not the current system date
 */
export function filterObservationsByTimeRange(
  observations: Observation[],
  timeRange: TimeRange
): Observation[] {
  if (timeRange === 'All' || observations.length === 0) {
    return observations
  }

  // Find the most recent observation date to calculate the cutoff from
  const mostRecentDate = new Date(observations[observations.length - 1].date)

  // Calculate how many months or years to subtract
  let monthsToSubtract = 0
  let yearsToSubtract = 0

  switch (timeRange) {
    case '1M':
      monthsToSubtract = 1
      break
    case '3M':
      monthsToSubtract = 3
      break
    case '6M':
      monthsToSubtract = 6
      break
    case '1Y':
      yearsToSubtract = 1
      break
    case '5Y':
      yearsToSubtract = 5
      break
  }

  // Calculate cutoff date using more robust arithmetic
  const cutoffDate = new Date(mostRecentDate)
  if (yearsToSubtract > 0) {
    cutoffDate.setFullYear(cutoffDate.getFullYear() - yearsToSubtract)
  } else if (monthsToSubtract > 0) {
    // Handle month subtraction more carefully to avoid overflow issues
    const targetMonth = cutoffDate.getMonth() - monthsToSubtract
    const targetYear = cutoffDate.getFullYear() + Math.floor(targetMonth / 12)
    const normalizedMonth = ((targetMonth % 12) + 12) % 12
    cutoffDate.setFullYear(targetYear)
    cutoffDate.setMonth(normalizedMonth)
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
  // Return null for undefined percentage change (division by zero)
  const change = startValue !== 0 ? (changeAbsolute / startValue) * 100 : null

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
