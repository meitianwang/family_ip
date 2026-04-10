<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div class="flex flex-wrap items-center gap-3">
            <!-- Status Filter -->
            <div class="w-full sm:w-40">
              <select
                v-model="statusFilter"
                @change="reload()"
                class="input w-full"
              >
                <option value="">{{ t('payment.orders.allStatus') }}</option>
                <option value="pending">{{ t('payment.orderStatus.pending') }}</option>
                <option value="completed">{{ t('payment.orderStatus.completed') }}</option>
                <option value="failed">{{ t('payment.orderStatus.failed') }}</option>
                <option value="expired">{{ t('payment.orderStatus.expired') }}</option>
                <option value="cancelled">{{ t('payment.orderStatus.cancelled') }}</option>
              </select>
            </div>
          </div>
        </div>
      </template>

      <template #default>
        <!-- Loading -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" role="status" aria-label="Loading" />
        </div>

        <!-- Empty State -->
        <div v-else-if="orders.length === 0" class="py-12 text-center text-gray-500 dark:text-slate-400">
          {{ t('payment.orders.empty') }}
        </div>

        <!-- Orders Table -->
        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200 dark:border-slate-700">
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">ID</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.orders.amount') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.orders.type') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.orders.method') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.orders.status') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.orders.time') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="order in orders"
                :key="order.id"
                class="border-b border-gray-100 dark:border-slate-800"
              >
                <td class="px-4 py-3 text-gray-900 dark:text-slate-100">#{{ order.id }}</td>
                <td class="px-4 py-3 font-medium text-gray-900 dark:text-slate-100">¥{{ order.amount }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-slate-400">
                  {{ order.order_type === 'subscription' ? t('payment.orders.subscription') : t('payment.orders.balance') }}
                </td>
                <td class="px-4 py-3 text-gray-600 dark:text-slate-400">{{ getPaymentLabel(order.payment_type) }}</td>
                <td class="px-4 py-3">
                  <span
                    class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium"
                    :class="getPaymentStatusBadgeClass(order.status)"
                  >
                    {{ t(`payment.orderStatus.${order.status}`) }}
                  </span>
                </td>
                <td class="px-4 py-3 text-gray-500 dark:text-slate-400">{{ formatPaymentDate(order.created_at) }}</td>
                <td class="px-4 py-3">
                  <button
                    v-if="canRequestRefund(order)"
                    @click="openRefundDialog(order)"
                    :aria-label="t('payment.orders.requestRefund') + ' #' + order.id"
                    class="text-sm text-red-600 hover:text-red-700 dark:text-red-400"
                  >
                    {{ t('payment.orders.requestRefund') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="pagination" class="mt-4 flex items-center justify-between">
          <span class="text-sm text-gray-500 dark:text-slate-400">
            {{ t('common.total') }}: {{ pagination.total }}
          </span>
          <div class="flex gap-2">
            <button
              :disabled="pagination.page <= 1"
              @click="goToPage(pagination.page - 1)"
              class="btn btn-secondary btn-sm"
            >
              {{ t('common.prev') }}
            </button>
            <button
              :disabled="pagination.page >= pagination.pages"
              @click="goToPage(pagination.page + 1)"
              class="btn btn-secondary btn-sm"
            >
              {{ t('common.next') }}
            </button>
          </div>
        </div>
      </template>
    </TablePageLayout>

    <!-- Refund Dialog -->
    <BaseDialog :show="refundDialogOpen" @close="refundDialogOpen = false" :title="t('payment.orders.requestRefund')">
      <div class="space-y-4">
        <div>
          <label for="refund-amount" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">
            {{ t('payment.orders.refundAmount') }}
          </label>
          <input
            id="refund-amount"
            v-model="refundAmount"
            type="text"
            inputmode="decimal"
            pattern="^\d*(\.\d{0,2})?$"
            class="input w-full"
            :placeholder="t('payment.orders.refundAmountHint')"
          />
        </div>
        <div>
          <label for="refund-reason" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">
            {{ t('payment.orders.refundReason') }}
          </label>
          <textarea
            id="refund-reason"
            v-model="refundReason"
            rows="3"
            class="input w-full"
            maxlength="500"
            :placeholder="t('payment.orders.refundReasonHint')"
          />
        </div>
        <div class="flex justify-end gap-3">
          <button @click="refundDialogOpen = false" class="btn btn-secondary">{{ t('common.cancel') }}</button>
          <button @click="submitRefund" class="btn btn-primary" :disabled="refundLoading">
            {{ refundLoading ? t('payment.processing') : t('common.confirm') }}
          </button>
        </div>
      </div>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { payAPI } from '@/api/pay'
import { useAppStore } from '@/stores'
import { getPaymentMethodLabel, getPaymentStatusBadgeClass, formatPaymentDate } from '@/utils/payment'
import type { UserPaymentOrder } from '@/types'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const orders = ref<UserPaymentOrder[]>([])
const pagination = ref<{ page: number; pages: number; total: number } | null>(null)
const statusFilter = ref('')
const currentPage = ref(1)

// Refund dialog
const refundDialogOpen = ref(false)
const refundOrderId = ref(0)
const refundOrderAmount = ref('')
const refundAmount = ref('')
const refundReason = ref('')
const refundLoading = ref(false)

async function loadOrders() {
  loading.value = true
  try {
    const filters = statusFilter.value ? { status: statusFilter.value } : undefined
    const data = await payAPI.listOrders(currentPage.value, 20, filters)
    orders.value = data.items || []
    pagination.value = { page: data.page, pages: data.pages, total: data.total }
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    loading.value = false
  }
}

function reload() {
  currentPage.value = 1
  loadOrders()
}

function goToPage(page: number) {
  currentPage.value = page
  loadOrders()
}

onMounted(loadOrders)

const getPaymentLabel = (type: string) => getPaymentMethodLabel(type, t)

function canRequestRefund(order: UserPaymentOrder): boolean {
  return order.status === 'completed' && !order.refund_requested_at
}

function openRefundDialog(order: UserPaymentOrder) {
  refundOrderId.value = order.id
  refundOrderAmount.value = order.amount
  refundAmount.value = order.amount
  refundReason.value = ''
  refundDialogOpen.value = true
}

async function submitRefund() {
  if (!refundAmount.value || refundLoading.value) return
  const amount = parseFloat(refundAmount.value)
  if (isNaN(amount) || amount <= 0 || !Number.isFinite(amount)) {
    appStore.showError(t('payment.invalidAmount'))
    return
  }
  const maxAmount = parseFloat(refundOrderAmount.value)
  if (!isNaN(maxAmount) && amount > maxAmount) {
    appStore.showError(t('payment.maxAmount', { amount: refundOrderAmount.value }))
    return
  }
  // Capture values before async to avoid race conditions
  const orderId = refundOrderId.value
  const amountStr = refundAmount.value
  const reason = refundReason.value || undefined
  refundLoading.value = true
  try {
    await payAPI.requestRefund(orderId, amountStr, reason)
    appStore.showSuccess(t('payment.orders.refundRequested'))
    refundDialogOpen.value = false
    loadOrders()
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    refundLoading.value = false
  }
}
</script>
