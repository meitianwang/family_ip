<template>
  <!-- Cancel Blocked State -->
  <div v-if="cancelBlocked" class="flex flex-col items-center space-y-4 py-8">
    <div class="text-6xl text-green-600 dark:text-green-400">✓</div>
    <h2 class="text-xl font-bold text-green-600 dark:text-green-400">{{ t('payment.qr.paid') }}</h2>
    <p class="text-center text-sm text-gray-500 dark:text-slate-400">{{ t('payment.qr.paidCancelBlocked') }}</p>
    <button
      @click="emit('back')"
      class="mt-4 w-full rounded-lg py-3 font-medium text-white bg-blue-600 hover:bg-blue-700 dark:bg-blue-600/90 dark:hover:bg-blue-600"
    >
      {{ t('payment.qr.backToRecharge') }}
    </button>
  </div>

  <!-- Main Payment View -->
  <div v-else class="flex flex-col items-center space-y-4">
    <!-- Amount Display -->
    <div class="text-center">
      <div class="text-4xl font-bold text-blue-600 dark:text-blue-400">¥{{ displayAmount.toFixed(2) }}</div>
      <div v-if="hasFeeDiff" class="mt-1 text-sm text-gray-500 dark:text-slate-400">
        {{ t('payment.qr.credited') }}¥{{ amount.toFixed(2) }}
      </div>
      <div
        class="mt-1 text-sm"
        :class="expired ? 'text-red-500' : timeLeftSeconds <= 60 ? 'text-red-500 animate-pulse' : 'text-gray-500 dark:text-slate-400'"
      >
        {{ expired ? t('payment.qr.expired') : `${t('payment.qr.remaining')}: ${timeLeft}` }}
      </div>
    </div>

    <!-- Payment Content (when not expired) -->
    <template v-if="!expired">
      <!-- Auto Redirect (mobile/H5) -->
      <template v-if="shouldAutoRedirect">
        <div class="flex items-center justify-center py-6">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-blue-500 border-t-transparent" />
          <span class="ml-3 text-sm text-gray-500 dark:text-slate-400">
            {{ t('payment.qr.redirecting', { channel: channelLabel }) }}
          </span>
        </div>
        <a
          :href="payUrl!"
          target="_self"
          rel="noopener noreferrer"
          class="flex w-full items-center justify-center gap-2 rounded-lg py-3 font-medium text-white shadow-md"
          :class="getPaymentButtonClass(paymentType)"
        >
          {{ redirected ? t('payment.qr.notRedirected', { channel: channelLabel }) : t('payment.qr.goto', { channel: channelLabel }) }}
        </a>
        <p class="text-center text-sm text-gray-500 dark:text-slate-400">{{ t('payment.qr.h5Hint') }}</p>
      </template>

      <!-- QR Code Display -->
      <template v-else-if="!isStripe">
        <div
          v-if="qrDataUrl"
          class="relative rounded-lg border p-4 border-gray-200 bg-white dark:border-slate-700 dark:bg-slate-900"
        >
          <div v-if="imageLoading" class="absolute inset-0 z-10 flex items-center justify-center rounded-lg bg-black/10">
            <div class="h-8 w-8 animate-spin rounded-full border-2 border-blue-500 border-t-transparent" />
          </div>
          <img :src="qrDataUrl" alt="payment qrcode" class="h-56 w-56 rounded" />
        </div>
        <div v-else class="text-center">
          <div class="rounded-lg border-2 border-dashed p-8 border-gray-300 dark:border-slate-700">
            <p v-if="qrError" class="text-sm text-red-500 dark:text-red-400">{{ t('payment.qr.qrFailed') }}</p>
            <p v-else class="text-sm text-gray-500 dark:text-slate-400">{{ t('payment.qr.scanPay') }}</p>
          </div>
        </div>
        <p class="text-center text-sm text-gray-500 dark:text-slate-400">
          {{ t('payment.qr.openScan', { channel: channelLabel }) }}
        </p>
      </template>

      <!-- Stripe Payment (placeholder - requires Stripe JS SDK integration) -->
      <template v-else>
        <div class="w-full max-w-md space-y-4">
          <div v-if="!clientSecret" class="rounded-lg border-2 border-dashed p-8 text-center border-gray-300 dark:border-slate-700">
            <p class="text-sm text-gray-500 dark:text-slate-400">{{ t('payment.qr.initFailed') }}</p>
          </div>
          <div v-else class="rounded-lg border p-4 border-gray-200 bg-white dark:border-slate-700 dark:bg-slate-900">
            <p class="text-sm text-gray-500 dark:text-slate-400">{{ t('payment.qr.stripeRedirect') }}</p>
          </div>
        </div>
      </template>
    </template>

    <!-- Action Buttons -->
    <div class="flex w-full gap-3">
      <button
        @click="emit('back')"
        class="flex-1 rounded-lg border py-2 text-sm border-gray-300 text-gray-600 hover:bg-gray-50 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
      >
        {{ t('payment.qr.back') }}
      </button>
      <button
        v-if="!expired"
        @click="handleCancel"
        class="flex-1 rounded-lg border py-2 text-sm border-red-300 text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/30"
      >
        {{ t('payment.qr.cancelOrder') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import QRCode from 'qrcode'
import { payAPI } from '@/api/pay'
import { getPaymentButtonClass, getPaymentMethodLabel, isSafePaymentUrl } from '@/utils/payment'

const TERMINAL_STATUSES = new Set(['completed', 'expired', 'cancelled', 'failed', 'refunded', 'refund_failed'])

const props = withDefaults(
  defineProps<{
    orderId: number
    payUrl?: string
    qrCode?: string
    clientSecret?: string
    paymentType: string
    amount: number
    payAmount?: number
    expiresAt: string
    isMobile?: boolean
  }>(),
  {
    isMobile: false
  }
)

const emit = defineEmits<{
  statusChange: [status: string]
  back: []
  pollStopped: []
}>()

const { t } = useI18n()

// Timer constants
const TIMER_INTERVAL_MS = 1000
const POLL_INTERVAL_MS = 2000

const displayAmount = computed(() => props.payAmount ?? props.amount)
const hasFeeDiff = computed(() => props.payAmount !== undefined && props.payAmount !== props.amount)

const timeLeft = ref('')
const timeLeftSeconds = ref(Infinity)
const expired = ref(false)
const qrDataUrl = ref('')
const qrError = ref(false)
const imageLoading = ref(false)
const cancelBlocked = ref(false)
const redirected = ref(false)

const isStripe = computed(() => props.paymentType?.includes('stripe'))
const shouldAutoRedirect = computed(
  () => !expired.value && !isStripe.value && !!props.payUrl && (props.isMobile || !props.qrCode)
)

const channelLabel = computed(() => getPaymentMethodLabel(props.paymentType || 'alipay', t))

// Auto redirect for mobile/H5
watch(shouldAutoRedirect, (val) => {
  if (val && !redirected.value && props.payUrl && isSafePaymentUrl(props.payUrl)) {
    redirected.value = true
    window.location.replace(props.payUrl)
  }
}, { immediate: true })

// Generate QR code with cancellation support
let qrGeneration = 0
watch(
  () => props.qrCode,
  async (qrPayload) => {
    const gen = ++qrGeneration
    if (!qrPayload?.trim()) {
      qrDataUrl.value = ''
      return
    }
    imageLoading.value = true
    try {
      const url = await QRCode.toDataURL(qrPayload.trim(), {
        width: 224,
        margin: 1,
        errorCorrectionLevel: 'M'
      })
      if (gen === qrGeneration && !isUnmounted) {
        qrDataUrl.value = url
      }
    } catch (err: unknown) {
      console.warn('[PaymentQRCode] QR code generation failed:', err)
      if (gen === qrGeneration && !isUnmounted) {
        qrDataUrl.value = ''
        qrError.value = true
      }
    } finally {
      if (gen === qrGeneration && !isUnmounted) {
        imageLoading.value = false
      }
    }
  },
  { immediate: true }
)

// Timer countdown
let timerInterval: ReturnType<typeof setInterval> | null = null

function updateTimer() {
  const diff = new Date(props.expiresAt).getTime() - Date.now()
  if (diff <= 0) {
    timeLeft.value = t('payment.qr.expired')
    timeLeftSeconds.value = 0
    expired.value = true
    return
  }
  const minutes = Math.floor(diff / 60000)
  const seconds = Math.floor((diff % 60000) / 1000)
  timeLeft.value = `${minutes}:${seconds.toString().padStart(2, '0')}`
  timeLeftSeconds.value = Math.floor(diff / 1000)
}

// Status polling with abort support
let pollInterval: ReturnType<typeof setInterval> | null = null
let pollAbort: AbortController | null = null
let isUnmounted = false
let consecutivePollFailures = 0
const MAX_POLL_FAILURES = 5

async function pollStatus() {
  if (isUnmounted) return
  pollAbort?.abort()
  pollAbort = new AbortController()
  try {
    const order = await payAPI.getOrder(props.orderId, { signal: pollAbort.signal })
    consecutivePollFailures = 0
    if (!isUnmounted && (order.paid_at || TERMINAL_STATUSES.has(order.status))) {
      emit('statusChange', order.status)
    }
  } catch (err: unknown) {
    if (err instanceof DOMException && err.name === 'AbortError') return
    consecutivePollFailures++
    if (consecutivePollFailures >= MAX_POLL_FAILURES && pollInterval) {
      console.warn(`[PaymentQRCode] Polling stopped after ${MAX_POLL_FAILURES} consecutive failures`)
      clearInterval(pollInterval)
      pollInterval = null
      emit('pollStopped')
    }
  }
}

onMounted(() => {
  updateTimer()
  timerInterval = setInterval(updateTimer, TIMER_INTERVAL_MS)
  pollStatus()
  pollInterval = setInterval(pollStatus, POLL_INTERVAL_MS)
})

onUnmounted(() => {
  isUnmounted = true
  if (timerInterval) clearInterval(timerInterval)
  if (pollInterval) clearInterval(pollInterval)
  pollAbort?.abort()
  timerInterval = null
  pollInterval = null
  pollAbort = null
})

// Stop polling when expired
watch(expired, (val) => {
  if (val && pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
})

async function handleCancel() {
  try {
    // Check current status first
    const order = await payAPI.getOrder(props.orderId)
    if (order.paid_at) {
      cancelBlocked.value = true
      return
    }
    if (TERMINAL_STATUSES.has(order.status)) {
      emit('statusChange', order.status)
      return
    }

    await payAPI.cancelOrder(props.orderId)
    emit('statusChange', 'cancelled')
  } catch (err: unknown) {
    // Cancel may fail if order was paid between check and cancel — refresh status
    console.warn('[PaymentQRCode] Cancel order failed, refreshing status:', err)
    await pollStatus()
  }
}
</script>
