/**
 * Admin Proxy IP API endpoints
 * Node/product/rental management for administrators.
 */

import { apiClient } from '../client'

export interface AdminProxyNode {
  id: number
  ip_address: string
  country: string
  country_code: string
  city: string
  isp: string
  http_port: number
  vless_port: number
  vless_network: string
  vless_tls: boolean
  vless_sni: string
  vless_ws_path: string
  tags: string[]
  status: string
  description: string
}

export interface AdminProxyProduct {
  id: number
  name: string
  description: string
  duration_days: number
  traffic_limit_gb: number
  price: string
  sort_order: number
  is_active: boolean
}

export interface AdminProxyRental {
  id: number
  user_id: number
  node_id: number
  product_id: number
  status: string
  started_at?: string
  expires_at?: string
  traffic_used_bytes: number
  traffic_limit_bytes: number
  node?: AdminProxyNode
  product?: AdminProxyProduct
  credential?: {
    http_username: string
    http_password: string
    vless_uuid: string
    vless_link: string
  }
}

export interface AdminPaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  pages: number
}

// --- Nodes ---

export async function listNodes(params?: { country_code?: string; status?: string; page?: number; page_size?: number }): Promise<AdminPaginatedResponse<AdminProxyNode>> {
  const { data } = await apiClient.get<AdminPaginatedResponse<AdminProxyNode>>('/admin/proxy/nodes', { params })
  return data
}

export async function createNode(node: Omit<AdminProxyNode, 'id'>): Promise<AdminProxyNode> {
  const { data } = await apiClient.post<AdminProxyNode>('/admin/proxy/nodes', node)
  return data
}

export async function updateNode(id: number, node: Partial<AdminProxyNode>): Promise<AdminProxyNode> {
  const { data } = await apiClient.put<AdminProxyNode>(`/admin/proxy/nodes/${id}`, node)
  return data
}

export async function deleteNode(id: number): Promise<void> {
  await apiClient.delete(`/admin/proxy/nodes/${id}`)
}

// --- Products ---

export async function listProducts(): Promise<AdminProxyProduct[]> {
  const { data } = await apiClient.get<AdminProxyProduct[]>('/admin/proxy/products')
  return data
}

export async function createProduct(product: Omit<AdminProxyProduct, 'id'>): Promise<AdminProxyProduct> {
  const { data } = await apiClient.post<AdminProxyProduct>('/admin/proxy/products', product)
  return data
}

export async function updateProduct(id: number, product: Partial<AdminProxyProduct>): Promise<AdminProxyProduct> {
  const { data } = await apiClient.put<AdminProxyProduct>(`/admin/proxy/products/${id}`, product)
  return data
}

export async function deleteProduct(id: number): Promise<void> {
  await apiClient.delete(`/admin/proxy/products/${id}`)
}

// --- Rentals ---

export async function listRentals(params?: { status?: string; user_id?: number; page?: number; page_size?: number }): Promise<AdminPaginatedResponse<AdminProxyRental>> {
  const { data } = await apiClient.get<AdminPaginatedResponse<AdminProxyRental>>('/admin/proxy/rentals', { params })
  return data
}

export async function getRental(id: number): Promise<AdminProxyRental> {
  const { data } = await apiClient.get<AdminProxyRental>(`/admin/proxy/rentals/${id}`)
  return data
}

export async function updateTraffic(id: number, deltaGb: number, note?: string): Promise<void> {
  await apiClient.post(`/admin/proxy/rentals/${id}/traffic`, { delta_gb: deltaGb, note })
}

export async function forceExpire(id: number): Promise<void> {
  await apiClient.post(`/admin/proxy/rentals/${id}/expire`)
}

export async function getTrafficLogs(id: number): Promise<unknown[]> {
  const { data } = await apiClient.get<unknown[]>(`/admin/proxy/rentals/${id}/traffic`)
  return data
}

export const proxyAdminAPI = {
  listNodes,
  createNode,
  updateNode,
  deleteNode,
  listProducts,
  createProduct,
  updateProduct,
  deleteProduct,
  listRentals,
  getRental,
  updateTraffic,
  forceExpire,
  getTrafficLogs,
}

export default proxyAdminAPI
