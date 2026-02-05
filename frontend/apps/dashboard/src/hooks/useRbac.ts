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

/**
 * Hook to check if user has a specific permission
 * Returns true if user has the required permission level for the module
 */
export const useHasPermission = (
  moduleId: number,
  requiredLevel: string = PERMISSION_LEVELS.READ
) => {
  const [hasPermission, setHasPermission] = useState<boolean>(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    const checkPermission = async () => {
      try {
        setLoading(true)
        const response = await rbacApi.checkPermission(moduleId, requiredLevel)
        setHasPermission(response.has_permission || false)
      } catch (err) {
        setError(err instanceof Error ? err : new Error(String(err)))
        setHasPermission(false)
      } finally {
        setLoading(false)
      }
    }

    checkPermission()
  }, [moduleId, requiredLevel])

  return { hasPermission, loading, error }
}

/**
 * Hook to get all user permissions
 */
export const useUserPermissions = () => {
  const [permissions, setPermissions] = useState<UserPermissions>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    const fetchPermissions = async () => {
      try {
        setLoading(true)
        const response = await rbacApi.getUserPermissions()
        const resolved = (response as { data?: UserPermissions })?.data ?? response
        setPermissions(resolved || {})
        // Cache the permissions
        permissionCache.set('current_user', (resolved as UserPermissions) || {})
      } catch (err) {
        setError(err instanceof Error ? err : new Error(String(err)))
      } finally {
        setLoading(false)
      }
    }

    fetchPermissions()
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
