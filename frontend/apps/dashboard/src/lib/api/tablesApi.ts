import { api } from '@/lib/api'

export type TableStatus = 'available' | 'occupied' | 'reserved' | 'maintenance'
export type TableLocation = 'indoor' | 'outdoor' | 'vip' | 'rooftop' | 'private'
export type TableShape = 'square' | 'circle' | 'rectangle'

export interface RestaurantTable {
  id: number
  tenant_id: number
  restaurant_id: number
  table_number: number
  name?: string
  capacity: number
  location: TableLocation
  status: TableStatus
  description?: string
  qr_code_config?: unknown
  position_x: number
  position_y: number
  shape: TableShape
  layout_id?: number | null
  created_at: string
  updated_at: string
}

export interface TableStats {
  total_tables: number
  available: number
  occupied: number
  reserved: number
  maintenance: number
  total_capacity: number
}

export interface CreateTablePayload {
  table_number: number
  name?: string
  capacity: number
  location: TableLocation
  shape: TableShape
  description?: string
}

export type UpdateTablePayload = Partial<CreateTablePayload> & {
  status?: TableStatus
  position_x?: number
  position_y?: number
  layout_id?: number
}

export interface TablePositionPayload {
  table_id: number
  position_x: number
  position_y: number
}

export interface TableLayout {
  id: number
  tenant_id: number
  restaurant_id: number
  name: string
  floor_number: number
  canvas_width: number
  canvas_height: number
  background_image?: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface CreateLayoutPayload {
  name: string
  floor_number: number
  canvas_width: number
  canvas_height: number
  background_image?: string
  is_default?: boolean
}

const defaultStats: TableStats = {
  total_tables: 0,
  available: 0,
  occupied: 0,
  reserved: 0,
  maintenance: 0,
  total_capacity: 0,
}

export const tablesApi = {
  async list(status?: TableStatus): Promise<RestaurantTable[]> {
    const response = await api.get('/tables', {
      params: status ? { status } : undefined,
    })
    const payload = response.data?.data
    return Array.isArray(payload) ? (payload as RestaurantTable[]) : []
  },

  async get(id: number): Promise<RestaurantTable> {
    const response = await api.get(`/tables/${id}`)
    return response.data?.data as RestaurantTable
  },

  async create(data: CreateTablePayload): Promise<RestaurantTable> {
    const response = await api.post('/tables', data)
    return response.data?.data as RestaurantTable
  },

  async update(id: number, data: UpdateTablePayload): Promise<void> {
    await api.put(`/tables/${id}`, data)
  },

  async remove(id: number): Promise<void> {
    await api.delete(`/tables/${id}`)
  },

  async updateStatus(id: number, status: TableStatus): Promise<void> {
    await api.put(`/tables/${id}/status`, { status })
  },

  async updatePositions(positions: TablePositionPayload[]): Promise<void> {
    await api.put('/tables/positions', { positions })
  },

  async getStats(): Promise<TableStats> {
    const response = await api.get('/tables/stats')
    const payload = response.data?.data
    if (!payload || typeof payload !== 'object') {
      return defaultStats
    }

    return {
      total_tables: Number(payload.total_tables) || 0,
      available: Number(payload.available) || 0,
      occupied: Number(payload.occupied) || 0,
      reserved: Number(payload.reserved) || 0,
      maintenance: Number(payload.maintenance) || 0,
      total_capacity: Number(payload.total_capacity) || 0,
    }
  },

  async getLayouts(): Promise<TableLayout[]> {
    const response = await api.get('/tables/layouts')
    const payload = response.data?.data
    return Array.isArray(payload) ? (payload as TableLayout[]) : []
  },

  async createLayout(data: CreateLayoutPayload): Promise<TableLayout> {
    const response = await api.post('/tables/layouts', data)
    return response.data?.data as TableLayout
  },
}

