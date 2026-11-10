/**
 * Custom React Hooks for RBAC
 */

import { useEffect, useState, useCallback } from 'react'
import {
  rbacApi,
  permissionCache,
  hasPermissionLevel,
  PERMISSION_LEVELS,
  MODULE_IDS,
  UserPermissions,
  Role,
  RolePermission,
  PermissionCheckResponse,
} from '@/lib/rbac'
import { useAuthStore } from '@/stores/authStore'

const CURRENT_USER_PERMISSION_CACHE_KEY = 'current_user'
const PRIVILEGED_KDS_ROLES = new Set(['owner', 'admin', 'manager', 'super_admin'])

const normalizePermissionLevelString = (value: unknown): string | null => {
  if (typeof value !== 'string') {
    return null
  }

  const normalized = value.trim().toUpperCase()
  const validLevels = new Set([
    PERMISSION_LEVELS.NONE,
    PERMISSION_LEVELS.READ,
    PERMISSION_LEVELS.EXPORT,
    PERMISSION_LEVELS.WRITE,
    PERMISSION_LEVELS.DELETE,
    PERMISSION_LEVELS.ADMIN,
  ])

  return (validLevels as ReadonlySet<string>).has(normalized) ? normalized : null
}

const extractPermissionLevel = (value: unknown): string | null => {
  const direct = normalizePermissionLevelString(value)
  if (direct) {
    return direct
  }

  if (value && typeof value === 'object') {
    const record = value as Record<string, unknown>
    return (
      extractPermissionLevel(record.permission_level) ||
      extractPermissionLevel(record.permissionLevel) ||
      extractPermissionLevel(record.level)
    )
  }

  return null
}

const normalizePermissionsPayload = (payload: unknown): UserPermissions => {
  const normalized: UserPermissions = {}
  const assign = (moduleIdValue: unknown, levelValue: unknown) => {
    const moduleId = Number(moduleIdValue)
    const level = extractPermissionLevel(levelValue)
    if (!Number.isFinite(moduleId) || moduleId <= 0 || !level) {
      return
    }
    normalized[moduleId] = level
  }

  const processArray = (items: unknown[]) => {
    items.forEach((item) => {
      if (!item || typeof item !== 'object') {
        return
      }
      const record = item as Record<string, unknown>
      assign(
        record.module_id ?? record.moduleId ?? record.id,
        record.permission_level ?? record.permissionLevel ?? record.level ?? record.permission
      )
    })
  }

  const resolved = (payload as { data?: unknown })?.data ?? payload
  if (Array.isArray(resolved)) {
    processArray(resolved)
    return normalized
  }

  if (resolved && typeof resolved === 'object') {
    const record = resolved as Record<string, unknown>

    if (Array.isArray(record.permissions)) {
      processArray(record.permissions)
    }
    if (Array.isArray(record.data)) {
      processArray(record.data)
    }

    Object.entries(record).forEach(([key, value]) => {
      assign(key, value)
    })
  }

  return normalized
}

const readCurrentUserRole = (): string | null => {
  const role = useAuthStore.getState().user?.role
  return role ? String(role).toLowerCase() : null
}

const canUseKDSRoleFallback = () => {
  const role = readCurrentUserRole()
  return role ? PRIVILEGED_KDS_ROLES.has(role) : false
}

const applyKDSRoleFallback = (permissions: UserPermissions): UserPermissions => {
  if (!canUseKDSRoleFallback()) {
    return permissions
  }

  const currentLevel = extractPermissionLevel(permissions[MODULE_IDS.KDS])
  if (currentLevel) {
    return permissions
  }

  return {
    ...permissions,
    [MODULE_IDS.KDS]: PERMISSION_LEVELS.ADMIN,
  }
}

/**
 * Hook to check if user has a specific permission
 * Returns true if user has the required permission level for the module
 */
export const useHasPermission = (
  moduleId: number,
  requiredLevel: string = PERMISSION_LEVELS.READ
) => {
  const initialCached = applyKDSRoleFallback(
    permissionCache.get(CURRENT_USER_PERMISSION_CACHE_KEY) || {}
  )
  const [hasPermission, setHasPermission] = useState<boolean>(
    hasPermissionLevel(initialCached?.[moduleId], requiredLevel)
  )
  const [loading, setLoading] = useState(initialCached && Object.keys(initialCached).length === 0)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    let active = true

    const checkPermission = async () => {
      const cached = permissionCache.get(CURRENT_USER_PERMISSION_CACHE_KEY)
      if (cached && active) {
        setHasPermission(hasPermissionLevel(cached[moduleId], requiredLevel))
        setLoading(false)
      }

      try {
        if (!cached && active) {
          setLoading(true)
        }
        const response = await rbacApi.checkPermission(moduleId, requiredLevel)
        if (!active) {
          return
        }
        if (moduleId === MODULE_IDS.KDS && canUseKDSRoleFallback()) {
          setHasPermission(true)
        } else {
          setHasPermission(response.has_permission || false)
        }
        setError(null)
      } catch (err) {
        if (!active) {
          return
        }
        if (!cached) {
          if (moduleId === MODULE_IDS.KDS && canUseKDSRoleFallback()) {
            setHasPermission(true)
            setError(null)
            return
          }
          setError(err instanceof Error ? err : new Error(String(err)))
          setHasPermission(false)
        }
      } finally {
        if (active) {
          setLoading(false)
        }
      }
    }

    void checkPermission()
    return () => {
      active = false
    }
  }, [moduleId, requiredLevel])

  return { hasPermission, loading, error }
}

/**
 * Hook to get all user permissions
 */
export const useUserPermissions = () => {
  const cachedPermissions = applyKDSRoleFallback(
    permissionCache.get(CURRENT_USER_PERMISSION_CACHE_KEY) || {}
  )
  const [permissions, setPermissions] = useState<UserPermissions>(cachedPermissions)
  const [loading, setLoading] = useState(Object.keys(cachedPermissions).length === 0)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    let active = true

    const fetchPermissions = async () => {
      const cached = permissionCache.get(CURRENT_USER_PERMISSION_CACHE_KEY)
      if (cached && active) {
        // Serve cached permissions immediately, then revalidate in background.
        setPermissions(cached)
        setLoading(false)
      }

      try {
        if (!cached && active) {
          setLoading(true)
        }
        const response = await rbacApi.getUserPermissions()
        const resolved = applyKDSRoleFallback(normalizePermissionsPayload(response))
        if (!active) {
          return
        }
        setPermissions(resolved)
        // Cache the permissions
        permissionCache.set(CURRENT_USER_PERMISSION_CACHE_KEY, resolved)
        setError(null)
      } catch (err) {
        if (!active) {
          return
        }

        // Fail soft: keep cached permissions if available.
        const fallback = permissionCache.get(CURRENT_USER_PERMISSION_CACHE_KEY)
        if (fallback) {
          setPermissions(applyKDSRoleFallback(fallback))
        } else {
          if (canUseKDSRoleFallback()) {
            const roleFallbackPermissions: UserPermissions = {
              [MODULE_IDS.KDS]: PERMISSION_LEVELS.ADMIN,
            }
            setPermissions(roleFallbackPermissions)
            permissionCache.set(CURRENT_USER_PERMISSION_CACHE_KEY, roleFallbackPermissions)
            setError(null)
          } else {
            setPermissions({})
            setError(err instanceof Error ? err : new Error(String(err)))
          }
        }
      } finally {
        if (active) {
          setLoading(false)
        }
      }
    }

    void fetchPermissions()

    return () => {
      active = false
    }
  }, [])

  const checkPermission = useCallback(
    (moduleId: number, requiredLevel: string = PERMISSION_LEVELS.READ): boolean => {
      const userLevel = permissions[moduleId]
      return hasPermissionLevel(userLevel, requiredLevel)
    },
    [permissions]
  )

  return { permissions, checkPermission, loading, error }
}

export const warmupCurrentUserPermissions = async () => {
  try {
    const response = await rbacApi.getUserPermissions()
    const resolved = applyKDSRoleFallback(normalizePermissionsPayload(response))
    permissionCache.set(CURRENT_USER_PERMISSION_CACHE_KEY, resolved)
    return resolved
  } catch {
    return null
  }
}

/**
 * Hook to manage roles with CRUD operations
 */
export const useRoles = () => {
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)

  const fetchRoles = useCallback(async (tenantId: number) => {
    try {
      setLoading(true)
      setError(null)
      const response = await rbacApi.listRoles(tenantId)
      setRoles(Array.isArray(response) ? response : response.data || [])
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
    } finally {
      setLoading(false)
    }
  }, [])

  const getRole = useCallback(async (roleId: number) => {
    try {
      setLoading(true)
      const response = await rbacApi.getRole(roleId)
      return response.data || response
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
      throw error
    } finally {
      setLoading(false)
    }
  }, [])

  const createRole = useCallback(
    async (data: { role_name: string; description?: string }) => {
      try {
        setLoading(true)
        const response = await rbacApi.createRole(data)
        const newRole = response.data || response
        setRoles([...roles, newRole])
        return newRole
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [roles]
  )

  const updateRole = useCallback(
    async (roleId: number, data: { role_name?: string; description?: string }) => {
      try {
        setLoading(true)
        const response = await rbacApi.updateRole(roleId, data)
        const updatedRole = response.data || response
        setRoles(roles.map((r) => (r.id === roleId ? updatedRole : r)))
        return updatedRole
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [roles]
  )

  const deleteRole = useCallback(
    async (roleId: number) => {
      try {
        setLoading(true)
        await rbacApi.deleteRole(roleId)
        setRoles(roles.filter((r) => r.id !== roleId))
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [roles]
  )

  return { roles, loading, error, fetchRoles, getRole, createRole, updateRole, deleteRole }
}

/**
 * Hook to manage role permissions
 */
export const useRolePermissions = (roleId: number) => {
  const [permissions, setPermissions] = useState<RolePermission[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    const fetchPermissions = async () => {
      try {
        setLoading(true)
        setError(null)
        const response = await rbacApi.getRolePermissions(roleId)
        setPermissions(Array.isArray(response) ? response : response.data || [])
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
      } finally {
        setLoading(false)
      }
    }

    if (roleId) {
      fetchPermissions()
    }
  }, [roleId])

  const assignPermission = useCallback(
    async (moduleId: number, permissionLevel: string) => {
      try {
        setLoading(true)
        await rbacApi.assignPermission({
          role_id: roleId,
          module_id: moduleId,
          permission_level: permissionLevel,
        })
        // Refetch permissions after assignment
        const response = await rbacApi.getRolePermissions(roleId)
        setPermissions(Array.isArray(response) ? response : response.data || [])
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [roleId]
  )

  const revokePermission = useCallback(
    async (moduleId: number) => {
      try {
        setLoading(true)
        await rbacApi.revokePermission(roleId, moduleId)
        setPermissions(permissions.filter((p) => p.module_id !== moduleId))
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [roleId, permissions]
  )

  return { permissions, loading, error, assignPermission, revokePermission }
}

/**
 * Hook to manage user roles
 */
export const useUserRoles = (userId: number) => {
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    const fetchRoles = async () => {
      try {
        setLoading(true)
        setError(null)
        const response = await rbacApi.getUserRoles(userId)
        setRoles(Array.isArray(response) ? response : response.data || [])
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
      } finally {
        setLoading(false)
      }
    }

    if (userId) {
      fetchRoles()
    }
  }, [userId])

  const assignRole = useCallback(
    async (roleId: number) => {
      try {
        setLoading(true)
        await rbacApi.assignRoleToUser(userId, roleId)
        // Refetch roles after assignment
        const response = await rbacApi.getUserRoles(userId)
        setRoles(Array.isArray(response) ? response : response.data || [])
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [userId]
  )

  const removeRole = useCallback(
    async (roleId: number) => {
      try {
        setLoading(true)
        await rbacApi.removeRoleFromUser(userId, roleId)
        setRoles(roles.filter((r) => r.id !== roleId))
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [userId, roles]
  )

  return { roles, loading, error, assignRole, removeRole }
}

/**
 * Hook to manage user roles by email
 */
export const useUserRolesByEmail = (email: string) => {
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    const fetchRoles = async () => {
      try {
        setLoading(true)
        setError(null)
        const response = await rbacApi.getUserRolesByEmail(email)
        setRoles(Array.isArray(response) ? response : response.data || [])
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
      } finally {
        setLoading(false)
      }
    }

    if (email) {
      fetchRoles()
    }
  }, [email])

  const assignRole = useCallback(
    async (roleId: number) => {
      try {
        setLoading(true)
        await rbacApi.assignRoleToUserByEmail(email, roleId)
        // Refetch roles after assignment
        const response = await rbacApi.getUserRolesByEmail(email)
        setRoles(Array.isArray(response) ? response : response.data || [])
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [email]
  )

  const removeRole = useCallback(
    async (roleId: number) => {
      try {
        setLoading(true)
        await rbacApi.removeRoleFromUserByEmail(email, roleId)
        setRoles(roles.filter((r) => r.id !== roleId))
      } catch (err) {
        const error = err instanceof Error ? err : new Error(String(err))
        setError(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [email, roles]
  )

  return { roles, loading, error, assignRole, removeRole }
}
