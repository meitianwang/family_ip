/**
 * Admin Payment API endpoints
 * Handles order management, refunds, configuration, provider instances,
 * channels, subscription plans, and dashboard statistics.
 */

import { apiClient } from '../client'
import type {
  AdminPaymentOrder,
  AdminOrderDetail,
  PaymentChannel,
  PaymentSubscriptionPlan,
  PaymentProviderInstance,
  PaymentDashboardData,
  BasePaginationResponse
} from '@/types'

// ==================== Orders ====================

/**
 * List all payment orders with filtering and pagination (admin)
 * @param page - Page number (1-based)
 * @param pageSize - Items per page
 * @param filters - Optional filters by user, status, type, method, date range
 * @param options - Request options (abort signal)
 */
export async function listOrders(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    user_id?: number
    status?: string
    order_type?: string
    payment_type?: string
    date_from?: string
    date_to?: string
  },
  options?: { signal?: AbortSignal }
): Promise<BasePaginationResponse<AdminPaymentOrder>> {
  const { data } = await apiClient.get<BasePaginationResponse<AdminPaymentOrder>>(
    '/admin/pay/orders',
    { params: { page, page_size: pageSize, ...filters }, signal: options?.signal }
  )
  return data
}

/** Get order detail including audit logs */
export async function getOrderDetail(id: number): Promise<AdminOrderDetail> {
  const { data } = await apiClient.get<AdminOrderDetail>(`/admin/pay/orders/${id}`)
  return data
}

/** Admin cancel a pending order */
export async function cancelOrder(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.post<{ message: string }>(`/admin/pay/orders/${id}/cancel`)
  return data
}

/** Retry a failed recharge/subscription order */
export async function retryOrder(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.post<{ message: string }>(`/admin/pay/orders/${id}/retry`)
  return data
}

// ==================== Refunds ====================

/**
 * Process a refund for an order
 * @param req.order_id - Target order ID
 * @param req.amount - Refund amount as decimal string (must not exceed order amount)
 * @param req.reason - Optional reason
 * @param req.force - Force refund (bypass provider, deduct balance only)
 * @param req.deduct_balance - Whether to deduct from user's balance
 */
export async function processRefund(req: {
  order_id: number
  amount: string
  reason?: string
  force?: boolean
  deduct_balance?: boolean
}): Promise<{ message: string }> {
  const numAmount = parseFloat(req.amount)
  if (isNaN(numAmount) || numAmount <= 0 || !isFinite(numAmount)) {
    throw new Error('Invalid refund amount')
  }
  const { data } = await apiClient.post<{ message: string }>('/admin/pay/refund', req)
  return data
}

// ==================== Config ====================

/** Get all payment configuration key-value pairs (sensitive values masked) */
export async function getConfig(): Promise<{ configs: Record<string, string> }> {
  const { data } = await apiClient.get<{ configs: Record<string, string> }>('/admin/pay/config')
  return data
}

/**
 * Update payment configuration
 * @param configs - Key-value pairs to update (keys must have pay_ prefix)
 */
export async function updateConfig(
  configs: Record<string, string>
): Promise<{ message: string; updated: number }> {
  for (const key of Object.keys(configs)) {
    if (!key.startsWith('pay_') || !/^pay_[a-z][a-z0-9_]*$/.test(key)) {
      throw new Error(`Invalid config key: ${key}. Must start with 'pay_' and contain only lowercase letters, digits, and underscores`)
    }
    if (key.length > 64) {
      throw new Error(`Config key too long: ${key}. Maximum 64 characters`)
    }
  }
  const { data } = await apiClient.put<{ message: string; updated: number }>(
    '/admin/pay/config',
    { configs }
  )
  return data
}

// ==================== Provider Instances ====================

/** List all payment provider instances */
export async function listProviderInstances(): Promise<PaymentProviderInstance[]> {
  const { data } = await apiClient.get<PaymentProviderInstance[]>('/admin/pay/provider-instances')
  return data
}

/** Create a new payment provider instance */
export async function createProviderInstance(
  req: Omit<PaymentProviderInstance, 'id' | 'created_at' | 'updated_at'>
): Promise<PaymentProviderInstance> {
  const { data } = await apiClient.post<PaymentProviderInstance>(
    '/admin/pay/provider-instances',
    req
  )
  return data
}

/** Update an existing payment provider instance */
export async function updateProviderInstance(
  id: number,
  req: Partial<Omit<PaymentProviderInstance, 'id' | 'created_at' | 'updated_at'>>
): Promise<PaymentProviderInstance> {
  const { data } = await apiClient.put<PaymentProviderInstance>(
    `/admin/pay/provider-instances/${id}`,
    req
  )
  return data
}

/** Delete a payment provider instance */
export async function deleteProviderInstance(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(
    `/admin/pay/provider-instances/${id}`
  )
  return data
}

// ==================== Channels ====================

/** List all payment channels (admin view, includes disabled) */
export async function listChannels(): Promise<PaymentChannel[]> {
  const { data } = await apiClient.get<PaymentChannel[]>('/admin/pay/channels')
  return data
}

/** Create a new payment channel */
export async function createChannel(
  req: Omit<PaymentChannel, 'id' | 'created_at' | 'updated_at'>
): Promise<PaymentChannel> {
  const { data } = await apiClient.post<PaymentChannel>('/admin/pay/channels', req)
  return data
}

/** Update an existing payment channel */
export async function updateChannel(
  id: number,
  req: Partial<Omit<PaymentChannel, 'id' | 'created_at' | 'updated_at'>>
): Promise<PaymentChannel> {
  const { data } = await apiClient.put<PaymentChannel>(`/admin/pay/channels/${id}`, req)
  return data
}

/** Delete a payment channel */
export async function deleteChannel(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/pay/channels/${id}`)
  return data
}

// ==================== Subscription Plans ====================

/** List all subscription plans (admin view, includes not-for-sale) */
export async function listPlans(): Promise<PaymentSubscriptionPlan[]> {
  const { data } = await apiClient.get<PaymentSubscriptionPlan[]>(
    '/admin/pay/subscription-plans'
  )
  return data
}

/** Create a new subscription plan */
export async function createPlan(
  req: Omit<PaymentSubscriptionPlan, 'id' | 'created_at' | 'updated_at'>
): Promise<PaymentSubscriptionPlan> {
  const { data } = await apiClient.post<PaymentSubscriptionPlan>(
    '/admin/pay/subscription-plans',
    req
  )
  return data
}

/** Update an existing subscription plan */
export async function updatePlan(
  id: number,
  req: Partial<Omit<PaymentSubscriptionPlan, 'id' | 'created_at' | 'updated_at'>>
): Promise<PaymentSubscriptionPlan> {
  const { data } = await apiClient.put<PaymentSubscriptionPlan>(
    `/admin/pay/subscription-plans/${id}`,
    req
  )
  return data
}

/** Delete a subscription plan */
export async function deletePlan(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(
    `/admin/pay/subscription-plans/${id}`
  )
  return data
}

// ==================== Dashboard ====================

/**
 * Get payment dashboard statistics
 * @param days - Number of days for daily series (1-365, default 30)
 */
export async function getDashboard(days: number = 30): Promise<PaymentDashboardData> {
  const safeDays = Math.max(1, Math.min(365, Math.floor(days)))
  const { data } = await apiClient.get<PaymentDashboardData>('/admin/pay/dashboard', {
    params: { days: safeDays }
  })
  return data
}

export const adminPayAPI = {
  listOrders,
  getOrderDetail,
  cancelOrder,
  retryOrder,
  processRefund,
  getConfig,
  updateConfig,
  listProviderInstances,
  createProviderInstance,
  updateProviderInstance,
  deleteProviderInstance,
  listChannels,
  createChannel,
  updateChannel,
  deleteChannel,
  listPlans,
  createPlan,
  updatePlan,
  deletePlan,
  getDashboard
}

export default adminPayAPI
