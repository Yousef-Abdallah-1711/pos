'use client'

import { useEffect, useState } from 'react'
import { hasAuthenticatedSession } from '@/lib/authSession'
import { useAuthStore } from '@/stores/authStore'
import { useHasPermission } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { Loader2, Lock, Search, Plus, Edit2, Trash2, X, Eye, EyeOff } from 'lucide-react'
import { debugError } from '@/lib/debug'
import { getConfig } from '@/lib/config'
import { authFetch } from '@/lib/authFetch'

interface Role {
  id: number
  name: string
  role_name?: string
  description?: string
}

interface User {
  id: number
  name: string
  email: string
  phone?: string
  role?: string
  status?: string
  created_at?: string
  roles?: Role[]
}

type ModalType = 'create' | 'edit' | 'delete' | 'roles' | null

export default function UsersPage() {
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const hasSession = hasAuthenticatedSession(user, token)

  const { hasPermission: isAdmin, loading: permissionLoading } = useHasPermission(
    MODULE_IDS.SETTINGS,
    PERMISSION_LEVELS.WRITE
  )

  // User Management State
  const [users, setUsers] = useState<User[]>([])
  const [roles, setRoles] = useState<Role[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [success, setSuccess] = useState('')

  // Modal State
  const [activeModal, setActiveModal] = useState<ModalType>(null)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [showPassword, setShowPassword] = useState(false)

  // Form State
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    password: '',
  })

  const [selectedRoles, setSelectedRoles] = useState<number[]>([])

  // Load data
  useEffect(() => {
    if (permissionLoading) return
    if (!isAdmin) {
      setIsLoading(false)
      return
    }
    fetchData()
  }, [isAdmin, permissionLoading, token])

  const fetchData = async () => {
    if (!hasSession) return
    setIsLoading(true)
    setError('')
    const { apiUrl } = getConfig()
    try {
      // Fetch users
      const usersRes = await authFetch(`${apiUrl}/users`)
      if (usersRes.ok) {
        const data = await usersRes.json()
        const usersList = Array.isArray(data) ? data : data?.data || data?.users || []
        setUsers(usersList)
      } else {
        throw new Error('Failed to fetch users')
      }

      // Fetch roles
      const rolesRes = await authFetch(`${apiUrl}/rbac/roles`)
      if (rolesRes.ok) {
        const data = await rolesRes.json()
        const rolesList = Array.isArray(data) ? data : data?.data || data?.roles || []
        setRoles(rolesList)
      }
    } catch (err) {
      setError('Failed to load users: ' + (err instanceof Error ? err.message : 'Unknown error'))
      debugError(err)
    } finally {
      setIsLoading(false)
    }
  }

  const resetForm = () => {
    setFormData({ name: '', email: '', phone: '', password: '' })
    setSelectedRoles([])
    setSelectedUser(null)
    setShowPassword(false)
  }

  const openCreateModal = () => {
    resetForm()
    setActiveModal('create')
  }

  const openEditModal = (u: User) => {
    setSelectedUser(u)
    setFormData({
      name: u.name,
      email: u.email,
      phone: u.phone || '',
      password: '',
    })
    setActiveModal('edit')
  }

  const openRolesModal = async (u: User) => {
    setSelectedUser(u)
    setSelectedRoles([])
    try {
      const { apiUrl } = getConfig()
      const res = await authFetch(`${apiUrl}/rbac/users/${u.id}/roles`)
      if (res.ok) {
        const data = await res.json()
        setSelectedRoles((data?.data || []).map((r: Role) => r.id))
      }
    } catch (err) {
      debugError('Failed to fetch user roles:', err)
    }
    setActiveModal('roles')
  }

  const openDeleteModal = (u: User) => {
    setSelectedUser(u)
    setActiveModal('delete')
  }

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!hasSession) return

    if (!formData.email || !formData.password || !formData.name) {
      setError('Email, password, and name are required')
      return
    }

    setIsSubmitting(true)
    setError('')
    try {
      const { apiUrl } = getConfig()
      const res = await authFetch(`${apiUrl}/users`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: formData.name,
          email: formData.email,
          phone: formData.phone,
          password: formData.password,
        }),
      })

      if (!res.ok) {
        const err = await res.text()
        throw new Error(err || 'Failed to create user')
      }

      setSuccess('User created successfully!')
      setActiveModal(null)
      resetForm()
      await fetchData()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError('Error creating user: ' + (err instanceof Error ? err.message : 'Unknown error'))
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleUpdateUser = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!hasSession || !selectedUser) return

    if (!formData.name || !formData.email) {
      setError('Name and email are required')
      return
    }

    setIsSubmitting(true)
    setError('')
    try {
      const { apiUrl } = getConfig()
      const res = await authFetch(`${apiUrl}/users/${selectedUser.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: formData.name,
          phone: formData.phone,
        }),
      })

      if (!res.ok) {
        const err = await res.text()
        throw new Error(err || 'Failed to update user')
      }

      setSuccess('User updated successfully!')
      setActiveModal(null)
      resetForm()
      await fetchData()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError('Error updating user: ' + (err instanceof Error ? err.message : 'Unknown error'))
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleDeleteUser = async () => {
    if (!hasSession || !selectedUser) return

    setIsSubmitting(true)
    setError('')
    try {
      const { apiUrl } = getConfig()
      const res = await authFetch(`${apiUrl}/users/${selectedUser.id}`, {
        method: 'DELETE',
      })

      if (!res.ok) {
        const err = await res.text()
        throw new Error(err || 'Failed to delete user')
      }

      setSuccess('User deleted successfully!')
      setActiveModal(null)
      resetForm()
      await fetchData()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError('Error deleting user: ' + (err instanceof Error ? err.message : 'Unknown error'))
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleRoleToggle = (roleId: number) => {
    setSelectedRoles((prev) =>
      prev.includes(roleId) ? prev.filter((id) => id !== roleId) : [...prev, roleId]
    )
  }

  const handleSaveRoles = async () => {
    if (!hasSession || !selectedUser) return

    setIsSubmitting(true)
    setError('')
    try {
      // Fetch current roles
      const { apiUrl } = getConfig()
      const currentRes = await authFetch(`${apiUrl}/rbac/users/${selectedUser.id}/roles`)
      const currentRoles = currentRes.ok ? (await currentRes.json())?.data || [] : []
      const currentRoleIds = currentRoles.map((r: Role) => r.id)

      // Remove roles that are no longer selected
      for (const roleId of currentRoleIds) {
        if (!selectedRoles.includes(roleId)) {
          await authFetch(
            `${apiUrl}/rbac/users/${selectedUser.id}/roles/${roleId}`,
            { method: 'DELETE' }
          )
        }
      }

      // Add new roles
      for (const roleId of selectedRoles) {
        if (!currentRoleIds.includes(roleId)) {
          await authFetch(`${apiUrl}/rbac/users/${selectedUser.id}/roles`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ roleId }),
          })
        }
      }

      setSuccess('Roles updated successfully!')
      setActiveModal(null)
      resetForm()
      await fetchData()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError('Error updating roles: ' + (err instanceof Error ? err.message : 'Unknown error'))
    } finally {
      setIsSubmitting(false)
    }
  }

  // Filters
  const filteredUsers = users.filter(
    (u) =>
      u.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      u.email.toLowerCase().includes(searchQuery.toLowerCase())
  )

  // Permission checks
  if (permissionLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!isAdmin) {
    return (
      <div className="flex flex-col items-center justify-center h-96">
        <Lock className="w-16 h-16 text-gray-400 mb-4" />
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">Access Denied</h2>
        <p className="text-gray-600 dark:text-gray-400 text-center">
          You don't have permission to manage users. Please contact your administrator.
        </p>
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">User Management</h2>
          <p className="text-gray-600 dark:text-gray-400 text-sm mt-1">
            Create and manage users in your organization
          </p>
        </div>
        <button
          onClick={openCreateModal}
          className="inline-flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
        >
          <Plus className="w-4 h-4" />
          Add New User
        </button>
      </div>

      {/* Messages */}
      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-lg flex items-start justify-between">
          <span>{error}</span>
          <button onClick={() => setError('')} className="text-red-600 hover:text-red-700">
            <X className="w-4 h-4" />
          </button>
        </div>
      )}

      {success && (
        <div className="p-4 bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400 rounded-lg flex items-start justify-between">
          <span>{success}</span>
          <button onClick={() => setSuccess('')} className="text-green-600 hover:text-green-700">
            <X className="w-4 h-4" />
          </button>
        </div>
      )}

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-3 w-4 h-4 text-gray-400" />
        <input
          type="text"
          placeholder="Search users by name or email..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
        />
      </div>

      {/* Users Table */}
      {filteredUsers.length === 0 ? (
        <div className="flex flex-col items-center justify-center h-64 border-2 border-dashed border-gray-300 dark:border-gray-700 rounded-lg">
          <p className="text-gray-500 dark:text-gray-400 mb-4">No users found</p>
          {searchQuery && (
            <button
              onClick={() => setSearchQuery('')}
              className="text-primary-600 hover:text-primary-700 dark:text-primary-400 text-sm"
            >
              Clear search
            </button>
          )}
        </div>
      ) : (
        <div className="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
          <table className="w-full">
            <thead className="bg-gray-50 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white">
                  Name
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white">
                  Email
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white">
                  Phone
                </th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900 dark:text-white">
                  Status
                </th>
                <th className="px-6 py-3 text-right text-sm font-semibold text-gray-900 dark:text-white">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
              {filteredUsers.map((u) => (
                <tr key={u.id} className="hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                  <td className="px-6 py-4 text-sm text-gray-900 dark:text-white font-medium">{u.name}</td>
                  <td className="px-6 py-4 text-sm text-gray-600 dark:text-gray-400">{u.email}</td>
                  <td className="px-6 py-4 text-sm text-gray-600 dark:text-gray-400">{u.phone || '-'}</td>
                  <td className="px-6 py-4 text-sm">
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300">
                      {u.status || 'active'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex items-center justify-end gap-2">
                      <button
                        onClick={() => openRolesModal(u)}
                        className="inline-flex items-center gap-1 px-3 py-1 text-sm text-blue-600 dark:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded transition-colors"
                        title="Manage roles"
                      >
                        Roles
                      </button>
                      <button
                        onClick={() => openEditModal(u)}
                        className="inline-flex items-center gap-1 px-3 py-1 text-sm text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded transition-colors"
                        title="Edit user"
                      >
                        <Edit2 className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => openDeleteModal(u)}
                        className="inline-flex items-center gap-1 px-3 py-1 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors"
                        title="Delete user"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Create User Modal */}
      {activeModal === 'create' && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-900 rounded-lg max-w-md w-full p-6 space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Create New User</h3>
              <button onClick={() => setActiveModal(null)} className="text-gray-400 hover:text-gray-600">
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleCreateUser} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Name *
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  placeholder="John Doe"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Email *
                </label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  placeholder="john@example.com"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Phone
                </label>
                <input
                  type="tel"
                  value={formData.phone}
                  onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  placeholder="+1 (555) 123-4567"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Password *
                </label>
                <div className="relative">
                  <input
                    type={showPassword ? 'text' : 'password'}
                    value={formData.password}
                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    placeholder="••••••••"
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-2.5 text-gray-400 hover:text-gray-600"
                  >
                    {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">Minimum 6 characters</p>
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setActiveModal(null)}
                  className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
                >
                  {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                  Create User
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Edit User Modal */}
      {activeModal === 'edit' && selectedUser && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-900 rounded-lg max-w-md w-full p-6 space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Edit User</h3>
              <button onClick={() => setActiveModal(null)} className="text-gray-400 hover:text-gray-600">
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleUpdateUser} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Name *
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Email (cannot be changed)
                </label>
                <input
                  type="email"
                  value={formData.email}
                  disabled
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 cursor-not-allowed"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Phone
                </label>
                <input
                  type="tel"
                  value={formData.phone}
                  onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                />
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setActiveModal(null)}
                  className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
                >
                  {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                  Update User
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Manage Roles Modal */}
      {activeModal === 'roles' && selectedUser && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-900 rounded-lg max-w-md w-full p-6 space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                Manage Roles for {selectedUser.name}
              </h3>
              <button onClick={() => setActiveModal(null)} className="text-gray-400 hover:text-gray-600">
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-3 max-h-64 overflow-y-auto">
              {roles.length === 0 ? (
                <p className="text-sm text-gray-500 dark:text-gray-400">No roles available</p>
              ) : (
                roles.map((role) => (
                  <label
                    key={role.id}
                    className="flex items-start gap-3 p-3 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors"
                  >
                    <input
                      type="checkbox"
                      checked={selectedRoles.includes(role.id)}
                      onChange={() => handleRoleToggle(role.id)}
                      className="mt-1 w-4 h-4 rounded border-gray-300 dark:border-gray-600 text-primary-600 focus:ring-2 focus:ring-primary-500"
                    />
                    <div className="flex-1">
                      <p className="text-sm font-medium text-gray-900 dark:text-white">
                        {role.name || role.role_name}
                      </p>
                      {role.description && (
                        <p className="text-xs text-gray-600 dark:text-gray-400 mt-0.5">{role.description}</p>
                      )}
                    </div>
                  </label>
                ))
              )}
            </div>

            <div className="flex gap-3 pt-4">
              <button
                onClick={() => setActiveModal(null)}
                className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleSaveRoles}
                disabled={isSubmitting}
                className="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
              >
                {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                Save Roles
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Delete User Modal */}
      {activeModal === 'delete' && selectedUser && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-900 rounded-lg max-w-md w-full p-6 space-y-4">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Delete User?</h3>

            <p className="text-gray-600 dark:text-gray-400">
              Are you sure you want to delete <strong>{selectedUser.name}</strong> ({selectedUser.email})? This
              action cannot be undone.
            </p>

            <div className="flex gap-3 pt-4">
              <button
                onClick={() => setActiveModal(null)}
                className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleDeleteUser}
                disabled={isSubmitting}
                className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
              >
                {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
