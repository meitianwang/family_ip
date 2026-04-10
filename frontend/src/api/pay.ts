/**
 * User Payment API endpoints
 * Handles order creation, queries, cancellation, refunds, and configuration.
 */

import { apiClient } from './client'
import type {
  UserPaymentOrder,
  CreateOrderRequest,
  CreateOrderResponse,
  PaymentConfig,
  PaymentChannel,
  PaymentSubscriptionPlan,
  BasePaginationResponse
} from '@/types'

/**
 * Create a new payment order (balance recharge or subscription purchase)
 * @param req - Order creation request with amount (decimal string), payment type, and order type
 * @returns Order response with pay URL / QR code for payment
 */
export async function createOrder(req: CreateOrderRequest): Promise<CreateOrderResponse> {
  const { data } = await apiClient.post<CreateOrderResponse>('/pay/orders', req)
  return data
}

/**
 * List current user's payment orders with pagination
 * @param page - Page number (1-based)
 * @param pageSize - Items per page
 * @param filters - Optional status filter
 * @param options - Request options (abort signal)
 */
export async function listOrders(
  page: number = 1,
  pageSize: number = 20,
  filters?: { status?: string },
  options?: { signal?: AbortSignal }
): Promise<BasePaginationResponse<UserPaymentOrder>> {
  const { data } = await apiClient.get<BasePaginationResponse<UserPaymentOrder>>('/pay/orders', {
    params: { page, page_size: pageSize, ...filters },
    signal: options?.signal
  })
  return data
}

/**
 * Get a single order by ID (must belong to current user)
 * @param id - Order ID
 * @param options - Request options (abort signal)
 */
export async function getOrder(
  id: number,
  options?: { signal?: AbortSignal }
): Promise<UserPaymentOrder> {
  const { data } = await apiClient.get<UserPaymentOrder>(`/pay/orders/${id}`, {
    signal: options?.signal
  })
  return data
}

/**
 * Cancel a pending order
 * @param id - Order ID (must be in pending status)
 */
export async function cancelOrder(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.post<{ message: string }>(`/pay/orders/${id}/cancel`)
  return data
}

/**
 * Request a refund for a completed order
 * @param id - Order ID
 * @param amount - Refund amount as decimal string (e.g. "50.00"), must not exceed order amount
 * @param reason - Optional reason for the refund
 * @throws {Error} If amount is invalid or exceeds order amount
 */
export async function requestRefund(
  id: number,
  amount: string,
  reason?: string
): Promise<{ message: string }> {
  const numAmount = parseFloat(amount)
  if (isNaN(numAmount) || numAmount <= 0 || !isFinite(numAmount)) {
    throw new Error('Invalid refund amount')
  }
  if (reason !== undefined && reason.length > 500) {
    throw new Error('Reason too long (max 500 characters)')
  }
  const { data } = await apiClient.post<{ message: string }>(`/pay/orders/${id}/refund-request`, {
    amount,
    reason
  })
  return data
}

/**
 * Get payment configuration (enabled methods, limits, fee rates)
 */
export async function getConfig(): Promise<PaymentConfig> {
  const { data } = await apiClient.get<PaymentConfig>('/pay/config')
  return data
}

/**
 * List enabled payment channels (for user-facing channel selection)
 */
export async function listChannels(): Promise<PaymentChannel[]> {
  const { data } = await apiClient.get<PaymentChannel[]>('/pay/channels')
  return data
}

/**
 * List subscription plans available for purchase
 */
export async function listPlans(): Promise<PaymentSubscriptionPlan[]> {
  const { data } = await apiClient.get<PaymentSubscriptionPlan[]>('/pay/subscription-plans')
  return data
}

export const payAPI = {
  createOrder,
  listOrders,
  getOrder,
  cancelOrder,
  requestRefund,
  getConfig,
  listChannels,
  listPlans
}

export default payAPI
