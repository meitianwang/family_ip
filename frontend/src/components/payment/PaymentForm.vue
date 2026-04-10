<template>
  <form @submit.prevent="handleSubmit" class="space-y-6" :aria-label="t('payment.rechargeAccount')">
    <!-- Account Info -->
    <div class="rounded-xl border p-4 border-slate-200 bg-slate-50 dark:border-slate-700 dark:bg-slate-800/80">
      <div class="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
        {{ t('payment.rechargeAccount') }}
      </div>
      <div class="mt-1 text-base font-medium text-slate-900 dark:text-slate-100">
        {{ userName || `${t('payment.user')} #${userId}` }}
      </div>
      <div v-if="userBalance !== undefined" class="mt-1 text-sm text-slate-500 dark:text-slate-400">
        {{ t('payment.currentBalance') }}
        <span class="font-medium text-green-600">{{ (typeof userBalance === 'number' ? userBalance : 0).toFixed(2) }}</span>
      </div>
    </div>

    <!-- Fixed Amount Display -->
    <div
      v-if="fixedAmount"
      class="rounded-xl border p-4 text-center border-slate-200 bg-slate-50 dark:border-slate-700 dark:bg-slate-800/60"
    >
      <div class="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
        {{ t('payment.rechargeAmount') }}
      </div>
      <div class="mt-1 text-3xl font-bold text-emerald-600 dark:text-emerald-400">
        ¥{{ fixedAmount.toFixed(2) }}
      </div>
    </div>

    <!-- Amount Selection -->
    <template v-else>
      <div>
        <label class="mb-2 block text-sm font-medium text-slate-700 dark:text-slate-200">
          {{ t('payment.rechargeAmount') }}
        </label>
        <div class="grid grid-cols-3 gap-2">
          <button
            v-for="val in filteredQuickAmounts"
            :key="val"
            type="button"
            @click="handleQuickAmount(val)"
            :aria-label="`¥${val}`"
            :aria-pressed="amount === val"
            class="rounded-lg border-2 px-4 py-3 text-center font-medium transition-colors"
            :class="
              amount === val
                ? 'border-blue-500 bg-blue-50 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300'
                : 'border-gray-200 bg-white text-gray-700 hover:border-gray-300 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-200 dark:hover:border-slate-500'
            "
          >
            ¥{{ val }}
          </button>
        </div>
      </div>

      <div>
        <label class="mb-2 block text-sm font-medium text-slate-700 dark:text-slate-200">
          {{ t('payment.customAmount') }}
        </label>
        <div class="relative">
          <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-slate-500">¥</span>
          <input
            type="text"
            inputmode="decimal"
            pattern="^\d*(\.\d{0,2})?$"
            name="payment-amount"
            autocomplete="off"
            :value="customAmount"
            :disabled="loading"
            @input="handleCustomAmountChange(($event.target as HTMLInputElement).value)"
            :placeholder="`${effectiveMin} - ${effectiveMax}`"
            class="w-full rounded-lg border py-3 pl-8 pr-4 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 border-gray-300 bg-white text-gray-900 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100 disabled:opacity-50"
          />
        </div>
      </div>
    </template>

    <!-- Amount Validation Error -->
    <div
      v-if="!fixedAmount && customAmount !== '' && !isValid"
      class="text-xs text-amber-700 dark:text-amber-300"
    >
      {{ amountErrorMsg }}
    </div>

    <!-- Payment Method Selection -->
    <div v-if="enabledPaymentTypes.length > 1">
      <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-200">
        {{ t('payment.paymentMethod') }}
      </label>
      <div class="grid grid-cols-2 gap-3 sm:flex" role="radiogroup" :aria-label="t('payment.paymentMethod')">
        <button
          v-for="type in enabledPaymentTypes"
          :key="type"
          type="button"
          :disabled="isMethodUnavailable(type)"
          @click="!isMethodUnavailable(type) && (paymentType = type)"
          class="relative flex h-[58px] flex-col items-center justify-center rounded-lg border px-3 transition-all sm:flex-1"
          :class="getMethodButtonClass(type)"
          :title="isMethodUnavailable(type) ? t('payment.dailyLimitReached') : undefined"
          :aria-label="t('payment.selectMethod', { method: getPaymentMethodLabel(type, t) })"
          :aria-pressed="effectivePaymentType === type"
          role="radio"
        >
          <span class="flex items-center gap-2">
            <!-- Payment Icon -->
            <span
              v-if="getPaymentIconType(type) === 'alipay'"
              class="flex h-8 w-8 items-center justify-center rounded-md bg-[#00AEEF] text-xl font-bold leading-none text-white"
            >
              {{ t('payment.alipayIconText') }}
            </span>
            <span
              v-else-if="getPaymentIconType(type) === 'wxpay'"
              class="flex h-8 w-8 items-center justify-center rounded-full bg-[#07C160] text-white"
            >
              <svg viewBox="0 0 24 24" class="h-5 w-5" fill="currentColor">
                <path d="M10 3C6.13 3 3 5.58 3 8.75c0 1.7.84 3.23 2.17 4.29l-.5 2.21 2.4-1.32c.61.17 1.25.27 1.93.27.22 0 .43-.01.64-.03C9.41 13.72 9 12.88 9 12c0-3.31 3.13-6 7-6 .26 0 .51.01.76.03C15.96 3.98 13.19 3 10 3z" />
                <path d="M16 8c-3.31 0-6 2.24-6 5s2.69 5 6 5c.67 0 1.31-.1 1.9-.28l2.1 1.15-.55-2.44C20.77 15.52 22 13.86 22 12c0-2.21-2.69-4-6-4z" />
              </svg>
            </span>
            <span
              v-else-if="getPaymentIconType(type) === 'stripe'"
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#635bff] text-white"
            >
              <svg viewBox="0 0 24 24" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <rect x="2" y="5" width="20" height="14" rx="2" />
                <path d="M2 10h20" />
              </svg>
            </span>

            <span class="flex flex-col items-start leading-none">
              <span class="text-xl font-semibold tracking-tight">{{ getPaymentMethodLabel(type, t) }}</span>
              <span
                v-if="isMethodUnavailable(type)"
                class="text-[10px] tracking-wide text-red-400"
              >
                {{ t('payment.dailyLimitReachedShort') }}
              </span>
            </span>
          </span>
        </button>
      </div>

      <p
        v-if="methodLimits && methodLimits[effectivePaymentType] && !methodLimits[effectivePaymentType].available"
        class="mt-2 text-xs text-amber-600 dark:text-amber-300"
      >
        {{ t('payment.methodLimitReached') }}
      </p>
    </div>

    <!-- Fee Breakdown -->
    <div
      v-if="feeRate > 0 && selectedAmount > 0"
      class="rounded-xl border px-4 py-3 text-sm border-slate-200 bg-slate-50 text-slate-600 dark:border-slate-700 dark:bg-slate-800/60 dark:text-slate-300"
    >
      <div class="flex items-center justify-between">
        <span>{{ t('payment.rechargeAmount') }}</span>
        <span>¥{{ selectedAmount.toFixed(2) }}</span>
      </div>
      <div class="mt-1 flex items-center justify-between">
        <span>{{ t('payment.fee', { rate: feeRate }) }}</span>
        <span>¥{{ feeAmount.toFixed(2) }}</span>
      </div>
      <div class="mt-1.5 flex items-center justify-between border-t pt-1.5 font-medium border-slate-200 text-slate-900 dark:border-slate-700 dark:text-slate-100">
        <span>{{ t('payment.amountToPay') }}</span>
        <span>¥{{ payAmount.toFixed(2) }}</span>
      </div>
    </div>

    <!-- Pending Orders Warning -->
    <div
      v-if="pendingBlocked"
      class="rounded-lg border p-3 text-sm border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-700 dark:bg-amber-900/30 dark:text-amber-300"
    >
      {{ t('payment.pendingBlocked', { count: pendingCount }) }}
    </div>

    <!-- Submit Button -->
    <button
      type="submit"
      :disabled="!isValid || loading || pendingBlocked"
      class="w-full rounded-lg py-3 text-center font-medium transition-colors"
      :class="
        isValid && !loading && !pendingBlocked
          ? 'text-white ' + getPaymentButtonClass(effectivePaymentType)
          : 'cursor-not-allowed bg-gray-300 text-gray-500 dark:bg-slate-700 dark:text-slate-400'
      "
    >
      <template v-if="loading">{{ t('payment.processing') }}</template>
      <template v-else-if="pendingBlocked">{{ t('payment.tooManyPending') }}</template>
      <template v-else>{{ t('payment.rechargeNow') }} ¥{{ displayPayAmount.toFixed(2) }}</template>
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { MethodLimit } from '@/types'
import { getPaymentIconType, getPaymentMethodLabel, getPaymentButtonClass, getPaymentSelectedClass, AMOUNT_TEXT_PATTERN } from '@/utils/payment'

const props = withDefaults(
  defineProps<{
    userId: number
    userName?: string
    userBalance?: number
    enabledPaymentTypes: string[]
    methodLimits?: Record<string, MethodLimit>
    minAmount: number
    maxAmount: number
    loading?: boolean
    pendingBlocked?: boolean
    pendingCount?: number
    fixedAmount?: number
  }>(),
  {
    loading: false,
    pendingBlocked: false,
    pendingCount: 0
  }
)

const emit = defineEmits<{
  submit: [amount: number, paymentType: string]
}>()

const { t } = useI18n()

const QUICK_AMOUNTS = [10, 20, 50, 100, 200, 500, 1000, 2000]

const amount = ref<number | ''>(props.fixedAmount ?? '')
const paymentType = ref(props.enabledPaymentTypes[0] || 'alipay')
const customAmount = ref(props.fixedAmount ? String(props.fixedAmount) : '')

// Sync refs when fixedAmount prop changes
watch(() => props.fixedAmount, (val) => {
  if (val !== undefined) {
    amount.value = val
    customAmount.value = String(val)
  }
})

const effectivePaymentType = computed(() =>
  props.enabledPaymentTypes.includes(paymentType.value)
    ? paymentType.value
    : props.enabledPaymentTypes[0] || ''
)

const selectedAmount = computed(() => (typeof amount.value === 'number' ? amount.value : 0))

const methodSingleMax = computed(() => {
  const v = props.methodLimits?.[effectivePaymentType.value]?.single_max
  return v ? parseFloat(v) : 0
})
const methodSingleMin = computed(() => {
  const v = props.methodLimits?.[effectivePaymentType.value]?.single_min
  return v ? parseFloat(v) : 0
})
const effectiveMax = computed(() =>
  methodSingleMax.value > 0 ? methodSingleMax.value : props.maxAmount
)
const effectiveMin = computed(() =>
  methodSingleMin.value > 0 ? Math.max(methodSingleMin.value, props.minAmount) : props.minAmount
)

const feeRate = computed(() => {
  const v = props.methodLimits?.[effectivePaymentType.value]?.fee_rate
  return v ? parseFloat(v) : 0
})
// Calculate fee in cents to avoid floating-point precision issues
const feeAmount = computed(() => {
  if (feeRate.value <= 0 || selectedAmount.value <= 0) return 0
  const amountCents = Math.round(selectedAmount.value * 100)
  const feeCents = Math.ceil((amountCents * feeRate.value) / 100)
  return feeCents / 100
})
const payAmount = computed(() => {
  if (feeRate.value <= 0 || selectedAmount.value <= 0) return selectedAmount.value
  const amountCents = Math.round(selectedAmount.value * 100)
  const feeCents = Math.round(feeAmount.value * 100)
  return (amountCents + feeCents) / 100
})

const displayPayAmount = computed(() =>
  feeRate.value > 0 && selectedAmount.value > 0 ? payAmount.value : selectedAmount.value || 0
)

const isMethodAvailable = computed(
  () => !props.methodLimits || props.methodLimits[effectivePaymentType.value]?.available !== false
)

const isValid = computed(() => {
  const sa = selectedAmount.value
  return (
    sa >= effectiveMin.value &&
    sa <= effectiveMax.value &&
    Math.abs(Math.round(sa * 100) - sa * 100) < 1e-8 &&
    isMethodAvailable.value
  )
})

const filteredQuickAmounts = computed(() =>
  QUICK_AMOUNTS.filter((val) => val >= effectiveMin.value && val <= effectiveMax.value)
)

const amountErrorMsg = computed(() => {
  const num = parseFloat(customAmount.value)
  if (!isNaN(num)) {
    if (num < props.minAmount) return t('payment.minAmount', { amount: props.minAmount })
    if (num > effectiveMax.value) return t('payment.maxAmount', { amount: effectiveMax.value })
  }
  return t('payment.invalidAmount')
})

function handleQuickAmount(val: number) {
  amount.value = val
  customAmount.value = String(val)
}

function handleCustomAmountChange(val: string) {
  if (!AMOUNT_TEXT_PATTERN.test(val)) return
  customAmount.value = val
  if (val === '') {
    amount.value = ''
    return
  }
  const num = parseFloat(val)
  if (!isNaN(num) && num > 0 && Math.abs(Math.round(num * 100) - num * 100) < 1e-8) {
    amount.value = num
  } else {
    amount.value = ''
  }
}

function handleSubmit() {
  if (!isValid.value || props.loading) return
  emit('submit', selectedAmount.value, effectivePaymentType.value)
}

function isMethodUnavailable(type: string): boolean {
  return props.methodLimits?.[type]?.available === false
}

function getMethodButtonClass(type: string): string {
  if (isMethodUnavailable(type)) return 'cursor-not-allowed border-gray-200 bg-gray-50 opacity-50 dark:border-slate-700 dark:bg-slate-800/50'
  if (effectivePaymentType.value === type) return getPaymentSelectedClass(type)
  return 'border-gray-300 bg-white text-slate-700 hover:border-gray-400 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-200 dark:hover:border-slate-500'
}
</script>
