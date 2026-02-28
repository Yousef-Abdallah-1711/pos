import { getTableStatusNotice, normalizeTableStatus } from './table-ordering'

describe('table-ordering', () => {
  describe('normalizeTableStatus', () => {
    it('normalizes valid status values', () => {
      expect(normalizeTableStatus('occupied')).toBe('occupied')
      expect(normalizeTableStatus(' RESERVED ')).toBe('reserved')
    })

    it('returns null for unknown status', () => {
      expect(normalizeTableStatus('unknown')).toBeNull()
      expect(normalizeTableStatus('')).toBeNull()
      expect(normalizeTableStatus(undefined)).toBeNull()
    })
  })

  describe('getTableStatusNotice', () => {
    it('returns warning notice for occupied table', () => {
      expect(getTableStatusNotice('occupied')).toEqual({
        status: 'occupied',
        level: 'warning',
        translationKey: 'tableOrdering.statusWarnings.occupied',
      })
    })

    it('returns critical notice for maintenance table', () => {
      expect(getTableStatusNotice('maintenance')).toEqual({
        status: 'maintenance',
        level: 'critical',
        translationKey: 'tableOrdering.statusWarnings.maintenance',
      })
    })

    it('returns null for available or invalid statuses', () => {
      expect(getTableStatusNotice('available')).toBeNull()
      expect(getTableStatusNotice('bad-status')).toBeNull()
      expect(getTableStatusNotice(undefined)).toBeNull()
    })
  })
})
