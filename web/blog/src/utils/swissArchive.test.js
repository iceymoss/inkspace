import { describe, expect, it } from 'vitest'
import { formatSwissCode, getSwissPlateSpan } from './swissArchive'

describe('Swiss archive helpers', () => {
  it('formats real ids without truncating them', () => {
    expect(formatSwissCode('W–', 42)).toBe('W–042')
    expect(formatSwissCode('K–', 1234)).toBe('K–1234')
  })

  it('uses an explicit placeholder for invalid ids', () => {
    expect(formatSwissCode('PL.', 0)).toBe('PL.—')
    expect(formatSwissCode('P–', undefined)).toBe('P–—')
  })

  it('keeps the twelve-column plate rhythm stable', () => {
    expect(Array.from({ length: 6 }, (_, index) => getSwissPlateSpan(index))).toEqual([6, 3, 3, 3, 3, 6])
  })
})
