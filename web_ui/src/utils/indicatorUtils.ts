/**
 * Utility functions for indicator classification and filtering
 */

import { Indicator } from '../types'

export type RegionFilter = 'all' | 'us' | 'world'

export interface IndicatorCategory {
  name: string
  indicators: string[]
  order: number
}

/**
 * Global commodity indicators that don't belong to a specific country
 */
const GLOBAL_COMMODITY_INDICATORS = ['POILBREUSDM'] as const

/**
 * Determines if an indicator is US-specific based on its code and name
 */
export function isUSIndicator(indicator: Indicator): boolean {
  const code = indicator.code.toUpperCase()
  const name = indicator.name.toLowerCase()

  // World Bank indicators with USA country code
  if (code.includes('WB:') && code.includes(':USA')) {
    return true
  }

  // FRED indicators are typically US-specific unless they're global commodities
  if (!code.includes('WB:') && !code.includes('OECD:')) {
    // Check if this is a global commodity indicator
    if (GLOBAL_COMMODITY_INDICATORS.includes(code as any)) {
      return false
    }
    return true
  }

  // Check name for US-specific terms
  if (name.includes('united states') || name.includes('u.s.') || name.includes('federal')) {
    return true
  }

  return false
}

/**
 * Determines if an indicator is World/Global indicator
 */
export function isWorldIndicator(indicator: Indicator): boolean {
  const code = indicator.code.toUpperCase()
  const name = indicator.name.toLowerCase()

  // World Bank indicators with WLD (World) country code
  if (code.includes('WB:') && code.includes(':WLD')) {
    return true
  }

  // World Bank indicators for other countries/regions
  if (code.includes('WB:') && (code.includes(':CHN') || code.includes(':EMU'))) {
    return true
  }

  // Check name for global terms
  if (name.includes('global') || name.includes('world') || name.includes('china') || name.includes('euro area')) {
    return true
  }

  // Check if this is a global commodity indicator
  if (GLOBAL_COMMODITY_INDICATORS.includes(code as any)) {
    return true
  }

  return false
}

/**
 * Filters indicators based on region
 */
export function filterIndicatorsByRegion(indicators: Indicator[], filter: RegionFilter): Indicator[] {
  if (filter === 'all') {
    return indicators
  }

  if (filter === 'us') {
    return indicators.filter(isUSIndicator)
  }

  if (filter === 'world') {
    return indicators.filter(isWorldIndicator)
  }

  return indicators
}

/**
 * Categories and their order for display
 */
export const INDICATOR_CATEGORIES: Record<string, IndicatorCategory> = {
  growth: {
    name: 'Economic Growth',
    indicators: ['GDPC1', 'WB:NY.GDP.MKTP.KD:WLD', 'WB:NY.GDP.MKTP.KD:USA', 'WB:NY.GDP.MKTP.KD:CHN', 'WB:NY.GDP.MKTP.KD:EMU'],
    order: 1,
  },
  inflation: {
    name: 'Inflation',
    indicators: ['CPIAUCSL', 'CPILFESL', 'WB:FP.CPI.TOTL.ZG:WLD', 'WB:FP.CPI.TOTL.ZG:CHN', 'WB:FP.CPI.TOTL.ZG:EMU'],
    order: 2,
  },
  labor: {
    name: 'Labor Market',
    indicators: ['UNRATE'],
    order: 3,
  },
  monetary: {
    name: 'Monetary Policy',
    indicators: ['FEDFUNDS'],
    order: 4,
  },
  yieldCurve: {
    name: 'Yield Curve',
    indicators: ['T10Y2Y', 'DGS10', 'DGS2'],
    order: 5,
  },
  liquidity: {
    name: 'Liquidity & Credit',
    indicators: ['M2SL', 'WB:FM.LBL.BMNY.CN:CHN', 'WB:FS.AST.PRVT.GD.ZS:WLD'],
    order: 6,
  },
  commodities: {
    name: 'Commodities & Energy',
    indicators: ['POILBREUSDM', 'WB:EG.USE.PCAP.KG.OE:WLD', 'WB:EN.ATM.CO2E.KT:WLD'],
    order: 7,
  },
  manufacturing: {
    name: 'Manufacturing',
    indicators: ['MANEMP'],
    order: 8,
  },
  sentiment: {
    name: 'Sentiment',
    indicators: ['UMCSENT'],
    order: 9,
  },
  trade: {
    name: 'Trade',
    indicators: ['IMPCH', 'EXPCH', 'WB:NE.EXP.GNFS.ZS:WLD', 'WB:NE.IMP.GNFS.ZS:WLD'],
    order: 10,
  },
}

/**
 * Gets the category for a given indicator code
 */
export function getIndicatorCategory(code: string): string {
  for (const [categoryKey, category] of Object.entries(INDICATOR_CATEGORIES)) {
    if (category.indicators.includes(code)) {
      return categoryKey
    }
  }
  return 'other'
}

/**
 * Sorts indicators by category and importance
 */
export function sortIndicatorsByCategory(indicators: Indicator[]): Indicator[] {
  return [...indicators].sort((a, b) => {
    const categoryA = getIndicatorCategory(a.code)
    const categoryB = getIndicatorCategory(b.code)

    const orderA = INDICATOR_CATEGORIES[categoryA]?.order ?? 999
    const orderB = INDICATOR_CATEGORIES[categoryB]?.order ?? 999

    if (orderA !== orderB) {
      return orderA - orderB
    }

    // Within same category, sort by the order they appear in the category
    const indexA = INDICATOR_CATEGORIES[categoryA]?.indicators.indexOf(a.code) ?? 999
    const indexB = INDICATOR_CATEGORIES[categoryB]?.indicators.indexOf(b.code) ?? 999

    return indexA - indexB
  })
}
