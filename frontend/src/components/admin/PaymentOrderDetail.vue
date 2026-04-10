<template>
  <BaseDialog :show="true" @close="emit('close')" :title="t('payment.admin.orderDetail')" width="wide">
    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" role="status" aria-label="Loading" />
    </div>

    <div v-else-if="error" class="py-4 text-center text-sm text-red-600 dark:text-red-400">{{ error }}</div>

    <template v-else-if="detail">
      <!-- Order Fields -->
      <div class="grid grid-cols-2 gap-4 text-sm">
        <div>
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.admin.orderId') }}</span>
          <p class="font-medium text-gray-900 dark:text-slate-100">#{{ detail.order.id }}</p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.admin.user') }}</span>
          <p class="font-medium text-gray-900 dark:text-slate-100">
            {{ detail.order.user_email || detail.order.user_name || `#${detail.order.user_id}` }}
          </p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.orders.amount') }}</span>
          <p class="font-medium text-gray-900 dark:text-slate-100">¥{{ detail.order.amount }}</p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.admin.payAmount') }}</span>
          <p class="font-medium text-gray-900 dark:text-slate-100">¥{{ detail.order.pay_amount || detail.order.amount }}</p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.orders.status') }}</span>
          <p>
            <span class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium" :class="getPaymentStatusBadgeClass(detail.order.status)">
              {{ t(`payment.orderStatus.${detail.order.status}`) }}
            </span>
          </p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.orders.method') }}</span>
          <p class="font-medium text-gray-900 dark:text-slate-100">{{ detail.order.payment_type }}</p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.orders.type') }}</span>
          <p class="font-medium text-gray-900 dark:text-slate-100">{{ detail.order.order_type }}</p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.admin.rechargeCode') }}</span>
          <p class="font-mono text-sm text-gray-900 dark:text-slate-100">{{ detail.order.recharge_code }}</p>
        </div>
        <div v-if="detail.order.payment_trade_no">
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.admin.tradeNo') }}</span>
          <p class="font-mono text-sm text-gray-900 dark:text-slate-100">{{ detail.order.payment_trade_no }}</p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.orders.time') }}</span>
          <p class="text-gray-900 dark:text-slate-100">{{ formatPaymentDate(detail.order.created_at) }}</p>
        </div>
        <div v-if="detail.order.paid_at">
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.admin.paidAt') }}</span>
          <p class="text-gray-900 dark:text-slate-100">{{ formatPaymentDate(detail.order.paid_at) }}</p>
        </div>
        <div v-if="detail.order.refund_amount">
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.admin.refundAmount') }}</span>
          <p class="font-medium text-red-600 dark:text-red-400">¥{{ detail.order.refund_amount }}</p>
        </div>
        <div v-if="detail.order.failed_reason" class="col-span-2">
          <span class="text-gray-500 dark:text-slate-400">{{ t('payment.admin.failedReason') }}</span>
          <p class="text-red-600 dark:text-red-400">{{ detail.order.failed_reason }}</p>
        </div>
      </div>

      <!-- Audit Logs -->
      <div v-if="detail.audit_logs.length > 0" class="mt-6">
        <h4 class="mb-3 text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.auditLog') }}</h4>
        <div class="space-y-3">
          <div
            v-for="log in detail.audit_logs"
            :key="log.id"
            class="flex items-start gap-3 rounded-lg border px-3 py-2 text-sm border-gray-100 dark:border-slate-800"
          >
            <span class="shrink-0 text-gray-400 dark:text-slate-500">{{ formatPaymentDate(log.created_at) }}</span>
            <span class="font-medium text-gray-700 dark:text-slate-300">{{ log.action }}</span>
            <span v-if="log.detail" class="text-gray-500 dark:text-slate-400">{{ log.detail }}</span>
            <span v-if="log.operator" class="ml-auto text-xs text-gray-400 dark:text-slate-500">{{ log.operator }}</span>
          </div>
        </div>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { adminPayAPI } from '@/api/admin/pay'
import { formatPaymentDate, getPaymentStatusBadgeClass } from '@/utils/payment'
import type { AdminOrderDetail } from '@/types'

const props = defineProps<{
  orderId: number
}>()

const emit = defineEmits<{
  close: []
}>()

const { t } = useI18n()

const loading = ref(true)
const detail = ref<AdminOrderDetail | null>(null)
const error = ref('')

onMounted(async () => {
  try {
    detail.value = await adminPayAPI.getOrderDetail(props.orderId)
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : t('common.error')
  } finally {
    loading.value = false
  }
})


</script>
