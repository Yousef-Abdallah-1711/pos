/**
 * RBAC (Role-Based Access Control) Utility Functions
 * Handles permission checking, role management, and access control logic
 */

import { api } from './api'

const RBAC_PERMISSIONS_TIMEOUT_MS = 2500

// Permission level constants (must match backend)
export const PERMISSION_LEVELS = {
  NONE: 'NONE',
  READ: 'READ',
  EXPORT: 'EXPORT',
  WRITE: 'WRITE',
  DELETE: 'DELETE',
  ADMIN: 'ADMIN',
} as const

// Module IDs (must match backend migrations)
export const MODULE_IDS = {
  PRODUCTS: 1,
  CATEGORIES: 2,
  HR: 3,
  NOTIFICATIONS: 4,
  SETTINGS: 5,
  USERS: 6,
  ROLES: 7,
  REPORTS: 8,
  THEMES: 9,
  ORDERS: 10,
  LEAVES: 11,
  SALARIES: 12,
  KDS: 13,
  TABLES: 14,
} as const

// Permission level hierarchy (higher number = more permissions)
const PERMISSION_HIERARCHY: Record<string, number> = {
  [PERMISSION_LEVELS.NONE]: 0,
  [PERMISSION_LEVELS.READ]: 1,
  [PERMISSION_LEVELS.EXPORT]: 2,
  [PERMISSION_LEVELS.WRITE]: 3,
  [PERMISSION_LEVELS.DELETE]: 4,
  [PERMISSION_LEVELS.ADMIN]: 5,
}

const normalizePermissionLevelInput = (value: unknown): string | undefined => {
  if (typeof value === 'string') {
    const normalized = value.trim().toUpperCase()
    return normalized in PERMISSION_HIERARCHY ? normalized : undefined
  }

  if (value && typeof value === 'object') {
    const obj = value as Record<string, unknown>
    const nested =
      obj.permission_level ??
      obj.permissionLevel ??
      obj.level
    return normalizePermissionLevelInput(nested)
  }

  return undefined
}

/**
 * Check if permission level A is greater than or equal to permission level B
 */
export const hasPermissionLevel = (
  userLevel: unknown,
  requiredLevel: string
): boolean => {
  const normalizedUserLevel = normalizePermissionLevelInput(userLevel)
  const normalizedRequiredLevel = normalizePermissionLevelInput(requiredLevel)

  if (!normalizedUserLevel || !normalizedRequiredLevel) {
    return false
  }

  const userValue = PERMISSION_HIERARCHY[normalizedUserLevel] ?? -1
  const requiredValue = PERMISSION_HIERARCHY[normalizedRequiredLevel] ?? -1
  return userValue >= requiredValue
}

/**
 * Compare two permission levels
 */
export const comparePermissionLevels = (levelA: string, levelB: string): number => {
  const valueA = PERMISSION_HIERARCHY[levelA] ?? -1
  const valueB = PERMISSION_HIERARCHY[levelB] ?? -1
  return valueA - valueB
}

/**
 * Get next permission level in hierarchy
 */
export const getNextPermissionLevel = (current: string): string | null => {
  const levels = [
    PERMISSION_LEVELS.NONE,
    PERMISSION_LEVELS.READ,
    PERMISSION_LEVELS.EXPORT,
    PERMISSION_LEVELS.WRITE,
    PERMISSION_LEVELS.DELETE,
    PERMISSION_LEVELS.ADMIN,
  ]
  const currentIndex = levels.indexOf(current)
  return currentIndex < levels.length - 1 ? levels[currentIndex + 1] : null
}

/**
 * Get previous permission level in hierarchy
 */
export const getPreviousPermissionLevel = (current: string): string | null => {
  const levels = [
    PERMISSION_LEVELS.NONE,
    PERMISSION_LEVELS.READ,
    PERMISSION_LEVELS.EXPORT,
    PERMISSION_LEVELS.WRITE,
    PERMISSION_LEVELS.DELETE,
    PERMISSION_LEVELS.ADMIN,
  ]
  const currentIndex = levels.indexOf(current)
  return currentIndex > 0 ? levels[currentIndex - 1] : null
}

// API Types
export interface Role {
  id: number
  tenant_id: number
  role_name: string
  description?: string
  is_system: boolean
  created_at?: string
  updated_at?: string
}

export interface RolePermission {
  id: number
  tenant_id: number
  role_id: number
  module_id: number
  permission_level: string
  created_at?: string
  updated_at?: string
}

export interface UserRole {
  id: number
  tenant_id: number
  user_id: number
  role_id: number
  assigned_by?: number
  created_at?: string
  updated_at?: string
}

export interface PermissionCheckResponse {
  has_permission: boolean
  module_id: number
  permission_level: string
}

export interface UserPermissions {
  [moduleId: number]: string // permission level for each module
}

/**
 * API Client for RBAC Operations
 */
export const rbacApi = {
  // Role Management
  listRoles: async (tenantId: number) => {
    const response = await api.get(`/rbac/roles`, {
      params: { 'X-Tenant-ID': tenantId },
    })
    return response.data
  },

  getRole: async (roleId: number) => {
    const response = await api.get(`/rbac/roles/${roleId}`)
    return response.data
  },

  createRole: async (data: { role_name: string; description?: string }) => {
    const response = await api.post('/rbac/roles', data)
    return response.data
  },

  updateRole: async (
    roleId: number,
    data: { role_name?: string; description?: string }
  ) => {
    const response = await api.put(`/rbac/roles/${roleId}`, data)
    return response.data
  },

  deleteRole: async (roleId: number) => {
    const response = await api.delete(`/rbac/roles/${roleId}`)
    return response.data
  },

  // Permission Management
  assignPermission: async (data: {
    role_id: number
    module_id: number
    permission_level: string
  }) => {
    const response = await api.post('/rbac/permissions/assign', data)
    return response.data
  },

  revokePermission: async (roleId: number, moduleId: number) => {
    const response = await api.delete(
      `/rbac/permissions/${roleId}/${moduleId}`
    )
    return response.data
  },

  getRolePermissions: async (roleId: number) => {
    const response = await api.get(`/rbac/roles/${roleId}/permissions`)
    return response.data
  },

  checkPermission: async (
    moduleId: number,
    permissionLevel: string = PERMISSION_LEVELS.READ
  ) => {
    const response = await api.get('/rbac/check-permission', {
      timeout: RBAC_PERMISSIONS_TIMEOUT_MS,
      params: {
        moduleId,
        permissionLevel,
      },
    })
    return response.data as PermissionCheckResponse
  },

  getUserPermissions: async () => {
    const response = await api.get('/rbac/me/permissions', {
      timeout: RBAC_PERMISSIONS_TIMEOUT_MS,
    })
    return response.data
  },

  // User-Role Management
  assignRoleToUser: async (userId: number, roleId: number) => {
    const response = await api.post(`/rbac/users/${userId}/roles`, {
      role_id: roleId,
    })
    return response.data
  },

  removeRoleFromUser: async (userId: number, roleId: number) => {
    const response = await api.delete(`/rbac/users/${userId}/roles/${roleId}`)
    return response.data
  },

  getUserRoles: async (userId: number) => {
    const response = await api.get(`/rbac/users/${userId}/roles`)
    return response.data
  },

  getMyRoles: async () => {
    const response = await api.get('/rbac/me/roles')
    return response.data
  },

  assignRoleToUserByEmail: async (email: string, roleId: number) => {
    const response = await api.post(
      `/rbac/users-by-email/${email}/roles`,
      {
        role_id: roleId,
      }
    )
    return response.data
  },

  removeRoleFromUserByEmail: async (email: string, roleId: number) => {
    const response = await api.delete(
      `/rbac/users-by-email/${email}/roles/${roleId}`
    )
    return response.data
  },

  getUserRolesByEmail: async (email: string) => {
    const response = await api.get(`/rbac/users-by-email/${email}/roles`)
    return response.data
  },
}

/**
 * Local permission cache management
 */
class PermissionCache {
  private cache: Map<string, { permissions: UserPermissions; timestamp: number }>
  private ttl: number = 5 * 60 * 1000 // 5 minutes
  private storageKey = 'rbac-permission-cache-v1'

  constructor() {
    this.cache = new Map()
  }

  private getStorageItem(userId: string) {
    if (typeof window === 'undefined') {
      return null
    }

    try {
      const raw = window.localStorage.getItem(this.storageKey)
      if (!raw) {
        return null
      }

      const parsed = JSON.parse(raw) as Record<
        string,
        { permissions: UserPermissions; timestamp: number }
      >
      const entry = parsed[userId]
      if (!entry || typeof entry !== 'object') {
        return null
      }
      return entry
    } catch {
      return null
    }
  }

  private setStorageItem(userId: string, payload: { permissions: UserPermissions; timestamp: number }) {
    if (typeof window === 'undefined') {
      return
    }

    try {
      const raw = window.localStorage.getItem(this.storageKey)
      const parsed = raw
        ? (JSON.parse(raw) as Record<string, { permissions: UserPermissions; timestamp: number }>)
        : {}

      parsed[userId] = payload
      window.localStorage.setItem(this.storageKey, JSON.stringify(parsed))
    } catch {
      // Ignore localStorage failures.
    }
  }

  private removeStorageItem(userId?: string) {
    if (typeof window === 'undefined') {
      return
    }

    try {
      if (!userId) {
        window.localStorage.removeItem(this.storageKey)
        return
      }

      const raw = window.localStorage.getItem(this.storageKey)
      if (!raw) {
        return
      }
      const parsed = JSON.parse(raw) as Record<
        string,
        { permissions: UserPermissions; timestamp: number }
      >
      delete parsed[userId]
      window.localStorage.setItem(this.storageKey, JSON.stringify(parsed))
    } catch {
      // Ignore localStorage failures.
    }
  }

  private isExpired(timestamp: number) {
    return Date.now() - timestamp > this.ttl
  }

  set(userId: string, permissions: UserPermissions): void {
    const payload = { permissions, timestamp: Date.now() }
    this.cache.set(userId, payload)
    this.setStorageItem(userId, payload)
  }

  get(userId: string): UserPermissions | null {
    const entry = this.cache.get(userId)
    if (entry && !this.isExpired(entry.timestamp)) {
      return entry.permissions
    }
    if (entry && this.isExpired(entry.timestamp)) {
      this.cache.delete(userId)
    }

    const storedEntry = this.getStorageItem(userId)
    if (!storedEntry) {
      return null
    }
    if (this.isExpired(storedEntry.timestamp)) {
      this.removeStorageItem(userId)
      return null
    }

    this.cache.set(userId, storedEntry)
    return storedEntry.permissions
  }

  clear(userId?: string): void {
    if (userId) {
      this.cache.delete(userId)
      this.removeStorageItem(userId)
    } else {
      this.cache.clear()
      this.removeStorageItem()
    }
  }

  setTTL(ttl: number): void {
    this.ttl = ttl
  }
}

export const permissionCache = new PermissionCache()

/**
 * Get module name from ID
 */
export const getModuleName = (moduleId: number): string => {
  const names: Record<number, string> = {
    [MODULE_IDS.PRODUCTS]: 'Products',
    [MODULE_IDS.CATEGORIES]: 'Categories',
    [MODULE_IDS.HR]: 'HR',
    [MODULE_IDS.NOTIFICATIONS]: 'Notifications',
    [MODULE_IDS.SETTINGS]: 'Settings',
    [MODULE_IDS.USERS]: 'Users',
    [MODULE_IDS.ROLES]: 'Roles',
    [MODULE_IDS.REPORTS]: 'Reports',
    [MODULE_IDS.THEMES]: 'Themes',
    [MODULE_IDS.ORDERS]: 'Orders',
    [MODULE_IDS.LEAVES]: 'Leaves',
    [MODULE_IDS.SALARIES]: 'Payroll',
    [MODULE_IDS.KDS]: 'Kitchen Display',
    [MODULE_IDS.TABLES]: 'Tables',
  }
  return names[moduleId] || `Module ${moduleId}`
}

/**
 * Get permission level description
 */
export const getPermissionDescription = (level: string): string => {
  const descriptions: Record<string, string> = {
    [PERMISSION_LEVELS.NONE]: 'No access',
    [PERMISSION_LEVELS.READ]: 'View only',
    [PERMISSION_LEVELS.EXPORT]: 'View and export',
    [PERMISSION_LEVELS.WRITE]: 'View and edit',
    [PERMISSION_LEVELS.DELETE]: 'View, edit, and delete',
    [PERMISSION_LEVELS.ADMIN]: 'Full access (admin)',
  }
  return descriptions[level] || 'Unknown'
}
