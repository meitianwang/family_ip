/**
 * Proxy IP API endpoints
 * Handles node listing, product listing, rental creation and management.
 */

import { apiClient } from './client'

export interface ProxyNode {
  id: number
  ip_address: string
  country: string
  country_code: string
  city: string
  isp: string
  http_port: number
  vless_port: number
  tags: string[]
  status: string
  description: string
}

export interface ProxyProduct {
  id: number
  name: string
  description: string
  duration_days: number
  traffic_limit_gb: number
  price: string
  sort_order: number
}

export interface ProxyCredential {
  http_host: string
  http_port: number
  http_username: string
  http_password: string
  vless_link: string
}

export interface ProxyRental {
  id: number
  node_id: number
  product_id: number
  status: string
  started_at?: string
  expires_at?: string
  traffic_used_bytes: number
  traffic_limit_bytes: number
  created_at: string
  node?: ProxyNode
  product?: ProxyProduct
  credential?: ProxyCredential
}

export interface CreateRentalResponse {
  rental_id: number
  order_id: number
  pay_url?: string
  qr_code?: string
  amount: string
  expires_at: string
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  pages: number
}

export async function listNodes(params?: { country_code?: string; tag?: string }): Promise<ProxyNode[]> {
  const { data } = await apiClient.get<ProxyNode[]>('/proxy/nodes', { params })
  return data
}

export async function getNode(id: number): Promise<ProxyNode> {
  const { data } = await apiClient.get<ProxyNode>(`/proxy/nodes/${id}`)
  return data
}

export async function listProducts(): Promise<ProxyProduct[]> {
  const { data } = await apiClient.get<ProxyProduct[]>('/proxy/products')
  return data
}

export async function createRental(nodeId: number, productId: number, payType: string): Promise<CreateRentalResponse> {
  const { data } = await apiClient.post<CreateRentalResponse>('/proxy/rentals', {
    node_id: nodeId,
    product_id: productId,
    pay_type: payType
  })
  return data
}

export async function listRentals(page = 1, pageSize = 20): Promise<PaginatedResponse<ProxyRental>> {
  const { data } = await apiClient.get<PaginatedResponse<ProxyRental>>('/proxy/rentals', {
    params: { page, page_size: pageSize }
  })
  return data
}

export async function getRental(id: number): Promise<ProxyRental> {
  const { data } = await apiClient.get<ProxyRental>(`/proxy/rentals/${id}`)
  return data
}

export async function cancelRental(id: number): Promise<void> {
  await apiClient.post(`/proxy/rentals/${id}/cancel`)
}

export const proxyAPI = {
  listNodes,
  getNode,
  listProducts,
  createRental,
  listRentals,
  getRental,
  cancelRental
}

export default proxyAPI
