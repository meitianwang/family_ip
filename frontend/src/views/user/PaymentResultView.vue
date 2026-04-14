<template>
  <AppLayout>
    <div class="mx-auto max-w-md px-4 py-12">
      <!-- Loading -->
      <div v-if="loading" class="flex flex-col items-center space-y-4 py-8">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-blue-500 border-t-transparent" />
        <p class="text-sm text-gray-500 dark:text-slate-400">{{ t('payment.status.checking') }}</p>
      </div>

      <!-- Result -->
      <div v-else class="flex flex-col items-center space-y-4 py-8">
        <div class="text-6xl" :class="statusConfig.color">{{ statusConfig.icon }}</div>
        <h2 class="text-xl font-bold" :class="statusConfig.color">{{ statusConfig.label }}</h2>
        <p class="text-center text-gray-500 dark:text-slate-400">{{ statusConfig.message }}</p>

        <p v-if="pollTimedOut" class="text-center text-sm text-amber-600 dark:text-amber-400">
          {{ t('payment.status.pollTimeout') }}
        </p>

        <!-- Navigation -->
        <div class="mt-6 flex w-full flex-col gap-3">
          <!-- Proxy rental success: show "View My Proxy" button -->
          <button
            v-if="isSuccess && pendingRentalId"
            class="w-full rounded-lg py-3 text-center font-medium text-white bg-blue-600 hover:bg-blue-700"
            @click="() => { sessionStorage.removeItem('pending_rental_id'); router.push(`/proxy/rentals/${pendingRentalId}`) }"
          >
            查看我的代理
          </button>
          <router-link
            v-else
            to="/purchase"
            class="w-full rounded-lg py-3 text-center font-medium text-white bg-blue-600 hover:bg-blue-700"
          >
            {{ isSuccess ? t('payment.status.done') : t('payment.status.backToRecharge') }}
          </router-link>
          <router-link
            to="/purchase/orders"
            class="w-full rounded-lg border py-3 text-center text-sm border-gray-300 text-gray-600 hover:bg-gray-50 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
          >
            {{ t('payment.orders.title') }}
          </router-link>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import { payAPI } from '@/api/pay'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const pendingRentalId = ref<string | null>(null)

const loading = ref(true)
const currentStatus = ref('pending')
const pollTimedOut = ref(false)

const isSuccess = computed(() => currentStatus.value === 'completed')

const TERMINAL_STATUSES = new Set([
  'completed', 'expired', 'cancelled', 'failed',
  'refunded', 'refund_failed', 'partially_refunded'
])

const statusConfig = computed(() => {
  const s = currentStatus.value

  if (s === 'completed') {
    return {
      label: t('payment.status.rechargeSuccess'),
      color: 'text-green-600 dark:text-green-400',
      icon: '✓',
      message: t('payment.status.rechargeSuccessMsg')
    }
  }
  if (s === 'paid' || s === 'recharging') {
    return {
      label: t('payment.status.recharging'),
      color: 'text-blue-600 dark:text-blue-400',
      icon: '⟳',
      message: t('payment.status.rechargingMsg')
    }
  }
  if (s === 'failed') {
    return {
      label: t('payment.status.paymentFailed'),
      color: 'text-red-600 dark:text-red-400',
      icon: '✗',
      message: t('payment.status.paymentFailedMsg')
    }
  }
  if (s === 'pending') {
    return {
      label: t('payment.status.awaitingPayment'),
      color: 'text-yellow-600 dark:text-yellow-400',
      icon: '⏳',
      message: t('payment.status.awaitingPaymentMsg')
    }
  }
  if (s === 'expired') {
    return {
      label: t('payment.status.orderExpired'),
      color: 'text-gray-500 dark:text-slate-400',
      icon: '⏰',
      message: t('payment.status.orderExpiredMsg')
    }
  }
  if (s === 'cancelled') {
    return {
      label: t('payment.status.cancelled'),
      color: 'text-gray-500 dark:text-slate-400',
      icon: '✗',
      message: t('payment.status.cancelledMsg')
    }
  }
  return {
    label: t('payment.status.error'),
    color: 'text-red-600 dark:text-red-400',
    icon: '✗',
    message: t('payment.status.errorMsg')
  }
})

// Polling
const POLL_INTERVAL_MS = 3000
const POLL_TIMEOUT_MS = 30000
let pollInterval: ReturnType<typeof setInterval> | null = null
let pollTimeout: ReturnType<typeof setTimeout> | null = null
let isUnmounted = false

async function fetchOrderStatus() {
  const orderId = route.query.out_trade_no as string
  if (!orderId) return
  try {
    const order = await payAPI.getOrder(Number(orderId))
    if (isUnmounted) return
    currentStatus.value = order.status
    if (TERMINAL_STATUSES.has(order.status)) {
      stopPolling()
    }
  } catch {
    // Silently fail on poll errors
  }
}

function stopPolling() {
  if (pollInterval) { clearInterval(pollInterval); pollInterval = null }
  if (pollTimeout) { clearTimeout(pollTimeout); pollTimeout = null }
}

function startPolling() {
  fetchOrderStatus()
  pollInterval = setInterval(fetchOrderStatus, POLL_INTERVAL_MS)
  pollTimeout = setTimeout(() => {
    stopPolling()
    if (!TERMINAL_STATUSES.has(currentStatus.value)) {
      pollTimedOut.value = true
    }
  }, POLL_TIMEOUT_MS)
}

onMounted(async () => {
  pendingRentalId.value = sessionStorage.getItem('pending_rental_id')

  const orderId = route.query.out_trade_no as string
  if (!orderId) {
    loading.value = false
    return
  }
  try {
    const order = await payAPI.getOrder(Number(orderId))
    currentStatus.value = order.status
    if (!TERMINAL_STATUSES.has(order.status)) {
      startPolling()
    }
  } catch {
    // If initial fetch fails, show pending state
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  isUnmounted = true
  stopPolling()
})
</script>
