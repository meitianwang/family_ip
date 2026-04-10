<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <select v-model="filters.status" @change="reload" class="input w-full sm:w-36">
            <option value="">{{ t('payment.orders.allStatus') }}</option>
            <option v-for="s in statusOptions" :key="s" :value="s">{{ t(`payment.orderStatus.${s}`) }}</option>
          </select>
          <select v-model="filters.order_type" @change="reload" class="input w-full sm:w-36">
            <option value="">{{ t('payment.admin.allTypes') }}</option>
            <option value="balance">{{ t('payment.orders.balance') }}</option>
            <option value="subscription">{{ t('payment.orders.subscription') }}</option>
          </select>
          <select v-model="filters.payment_type" @change="reload" class="input w-full sm:w-36">
            <option value="">{{ t('payment.admin.allMethods') }}</option>
            <option value="alipay">{{ t('payment.alipay') }}</option>
            <option value="wxpay">{{ t('payment.wechatPay') }}</option>
            <option value="stripe">Stripe</option>
          </select>
        </div>
      </template>

      <template #default>
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" role="status" aria-label="Loading" />
        </div>

        <div v-else-if="orders.length === 0" class="py-12 text-center text-gray-500 dark:text-slate-400">
          {{ t('payment.orders.empty') }}
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200 dark:border-slate-700">
                <th scope="col" class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">ID</th>
                <th scope="col" class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.user') }}</th>
                <th scope="col" class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.orders.amount') }}</th>
                <th scope="col" class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.orders.type') }}</th>
                <th scope="col" class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.orders.method') }}</th>
                <th scope="col" class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.orders.status') }}</th>
                <th scope="col" class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.orders.time') }}</th>
                <th scope="col" class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="order in orders" :key="order.id" class="border-b border-gray-100 dark:border-slate-800">
                <td class="px-3 py-3 text-gray-900 dark:text-slate-100">#{{ order.id }}</td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">{{ order.user_email || `#${order.user_id}` }}</td>
                <td class="px-3 py-3 font-medium text-gray-900 dark:text-slate-100">¥{{ order.amount }}</td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">{{ order.order_type }}</td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">{{ order.payment_type }}</td>
                <td class="px-3 py-3">
                  <span class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium" :class="getPaymentStatusBadgeClass(order.status)">
                    {{ t(`payment.orderStatus.${order.status}`) }}
                  </span>
                </td>
                <td class="px-3 py-3 text-gray-500 dark:text-slate-400">{{ formatPaymentDate(order.created_at) }}</td>
                <td class="px-3 py-3">
                  <div class="flex gap-2">
                    <button @click="viewDetail(order.id)" :aria-label="t('common.view') + ' #' + order.id" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">
                      {{ t('common.view') }}
                    </button>
                    <button
                      v-if="order.status === 'pending'"
                      @click="handleCancel(order.id)"
                      :aria-label="t('common.cancel') + ' #' + order.id"
                      class="text-sm text-amber-600 hover:text-amber-700 dark:text-amber-400"
                    >
                      {{ t('common.cancel') }}
                    </button>
                    <button
                      v-if="order.status === 'failed'"
                      @click="handleRetry(order.id)"
                      :aria-label="t('payment.admin.retry') + ' #' + order.id"
                      class="text-sm text-green-600 hover:text-green-700 dark:text-green-400"
                    >
                      {{ t('payment.admin.retry') }}
                    </button>
                    <button
                      v-if="canRefund(order.status)"
                      @click="openRefund(order)"
                      :aria-label="t('payment.admin.refund') + ' #' + order.id"
                      class="text-sm text-red-600 hover:text-red-700 dark:text-red-400"
                    >
                      {{ t('payment.admin.refund') }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="pagination" class="mt-4 flex items-center justify-between">
          <span class="text-sm text-gray-500 dark:text-slate-400">{{ t('common.total') }}: {{ pagination.total }}</span>
          <div class="flex gap-2">
            <button :disabled="pagination.page <= 1" @click="goToPage(pagination.page - 1)" class="btn btn-secondary btn-sm">
              {{ t('common.prev') }}
            </button>
            <button :disabled="pagination.page >= pagination.pages" @click="goToPage(pagination.page + 1)" class="btn btn-secondary btn-sm">
              {{ t('common.next') }}
            </button>
          </div>
        </div>
      </template>
    </TablePageLayout>

    <!-- Order Detail Modal -->
    <PaymentOrderDetail v-if="detailOrderId" :order-id="detailOrderId" @close="detailOrderId = null" />

    <!-- Refund Dialog -->
    <PaymentRefundDialog
      v-if="refundOrder"
      :order-id="refundOrder.id"
      :order-amount="refundOrder.amount"
      @close="refundOrder = null"
      @refunded="loadOrders"
    />

    <!-- Confirm Dialog for cancel/retry -->
    <ConfirmDialog
      :show="confirmAction.show"
      :title="confirmAction.title"
      :message="confirmAction.message"
      :danger="confirmAction.danger"
      @confirm="executeConfirmedAction"
      @cancel="confirmAction.show = false"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import PaymentOrderDetail from '@/components/admin/PaymentOrderDetail.vue'
import PaymentRefundDialog from '@/components/admin/PaymentRefundDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { getPaymentStatusBadgeClass, formatPaymentDate } from '@/utils/payment'
import { adminPayAPI } from '@/api/admin/pay'
import { useAppStore } from '@/stores'
import type { AdminPaymentOrder, BasePaginationResponse } from '@/types'

const { t } = useI18n()
const appStore = useAppStore()

const statusOptions = ['pending', 'paid', 'recharging', 'completed', 'expired', 'cancelled', 'failed', 'refund_requested', 'refunding', 'refunded', 'partially_refunded', 'refund_failed']

const loading = ref(true)
const orders = ref<AdminPaymentOrder[]>([])
const pagination = ref<{ page: number; pages: number; total: number } | null>(null)
const currentPage = ref(1)
const filters = reactive({ status: '', order_type: '', payment_type: '' })

const detailOrderId = ref<number | null>(null)
const refundOrder = ref<AdminPaymentOrder | null>(null)

async function loadOrders() {
  loading.value = true
  try {
    const f: { status?: string; order_type?: string; payment_type?: string } = {}
    if (filters.status) f.status = filters.status
    if (filters.order_type) f.order_type = filters.order_type
    if (filters.payment_type) f.payment_type = filters.payment_type

    const data: BasePaginationResponse<AdminPaymentOrder> = await adminPayAPI.listOrders(currentPage.value, 20, f)
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

function viewDetail(id: number) {
  detailOrderId.value = id
}

const confirmAction = reactive({
  show: false,
  title: '',
  message: '',
  danger: false,
  action: null as (() => Promise<void>) | null
})

function handleCancel(id: number) {
  confirmAction.title = t('common.cancel')
  confirmAction.message = t('payment.admin.confirmCancel')
  confirmAction.danger = true
  confirmAction.action = async () => {
    try {
      await adminPayAPI.cancelOrder(id)
      appStore.showSuccess(t('payment.admin.orderCancelled'))
      loadOrders()
    } catch (err: unknown) {
      appStore.showError(err instanceof Error ? err.message : t('common.error'))
    }
  }
  confirmAction.show = true
}

function handleRetry(id: number) {
  confirmAction.title = t('payment.admin.retry')
  confirmAction.message = t('payment.admin.confirmRetry')
  confirmAction.danger = false
  confirmAction.action = async () => {
    try {
      await adminPayAPI.retryOrder(id)
      appStore.showSuccess(t('payment.admin.retrySuccess'))
      loadOrders()
    } catch (err: unknown) {
      appStore.showError(err instanceof Error ? err.message : t('common.error'))
    }
  }
  confirmAction.show = true
}

async function executeConfirmedAction() {
  confirmAction.show = false
  if (confirmAction.action) {
    await confirmAction.action()
    confirmAction.action = null
  }
}

function canRefund(status: string): boolean {
  return ['completed', 'paid', 'refund_requested'].includes(status)
}

function openRefund(order: AdminPaymentOrder) {
  refundOrder.value = order
}


</script>
