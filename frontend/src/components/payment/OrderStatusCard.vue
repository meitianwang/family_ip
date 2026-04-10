<template>
  <div class="flex flex-col items-center space-y-4 py-8">
    <div class="text-6xl" :class="config.color" role="img" :aria-label="config.label">{{ config.icon }}</div>
    <h2 class="text-xl font-bold" :class="config.color">{{ config.label }}</h2>
    <p class="text-center text-gray-500 dark:text-slate-400">{{ config.message }}</p>
    <button
      @click="emit('back')"
      class="mt-4 w-full rounded-lg py-3 font-medium text-white bg-blue-600 hover:bg-blue-700 dark:hover:bg-blue-500"
    >
      {{ isSuccess ? t('payment.status.done') : t('payment.status.backToRecharge') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { payAPI } from '@/api/pay'

const props = defineProps<{
  orderId: number
  status: string
  paidAt?: string
  completedAt?: string
  failedReason?: string
}>()

const emit = defineEmits<{
  back: []
  stateChange: [status: string]
  pollStopped: []
}>()

const { t } = useI18n()

const currentStatus = ref(props.status)

watch(() => props.status, (val) => {
  currentStatus.value = val
})

const isSuccess = computed(() => currentStatus.value === 'completed')

const config = computed(() => {
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

// Timer constants
const POLL_INTERVAL_MS = 3000
const POLL_TIMEOUT_MS = 30000

const TERMINAL_STATUSES = new Set(['completed', 'expired', 'cancelled', 'failed', 'refunded', 'refund_failed'])

// Poll for status updates when in non-terminal state
let pollInterval: ReturnType<typeof setInterval> | null = null
let pollTimeout: ReturnType<typeof setTimeout> | null = null
let pollAbort: AbortController | null = null
let isUnmounted = false
let consecutiveFailures = 0
const MAX_FAILURES = 5

async function refreshOrder() {
  if (isUnmounted) return
  pollAbort?.abort()
  pollAbort = new AbortController()
  try {
    const order = await payAPI.getOrder(props.orderId, { signal: pollAbort.signal })
    if (isUnmounted) return
    consecutiveFailures = 0
    currentStatus.value = order.status
    emit('stateChange', order.status)
  } catch (err: unknown) {
    if (err instanceof DOMException && err.name === 'AbortError') return
    consecutiveFailures++
    if (consecutiveFailures >= MAX_FAILURES && pollInterval) {
      console.warn(`[OrderStatusCard] Polling stopped after ${MAX_FAILURES} consecutive failures`)
      clearInterval(pollInterval)
      pollInterval = null
      emit('pollStopped')
    }
  }
}

onMounted(() => {
  const s = currentStatus.value
  // Only poll for intermediate states
  if (s === 'paid' || s === 'recharging') {
    refreshOrder()
    pollInterval = setInterval(refreshOrder, POLL_INTERVAL_MS)
    pollTimeout = setTimeout(() => {
      if (pollInterval) clearInterval(pollInterval)
    }, POLL_TIMEOUT_MS)
  }
})

onUnmounted(() => {
  isUnmounted = true
  if (pollInterval) clearInterval(pollInterval)
  if (pollTimeout) clearTimeout(pollTimeout)
  pollAbort?.abort()
  pollInterval = null
  pollTimeout = null
  pollAbort = null
})

// Stop polling when reaching terminal state
watch(currentStatus, (val) => {
  if (TERMINAL_STATUSES.has(val) && pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
})
</script>
