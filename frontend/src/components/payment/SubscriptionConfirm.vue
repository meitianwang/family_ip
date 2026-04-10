<template>
  <div class="mx-auto max-w-lg space-y-6">
    <!-- Back link -->
    <button
      type="button"
      @click="emit('back')"
      class="flex items-center gap-1 text-sm transition-colors text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200"
    >
      <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
      </svg>
      {{ t('payment.confirm.backToPlans') }}
    </button>

    <!-- Title -->
    <h2 class="text-xl font-semibold text-slate-900 dark:text-slate-100">
      {{ t('payment.confirm.title') }}
    </h2>

    <!-- Plan detail card -->
    <div class="rounded-2xl border p-5 border-slate-200 bg-white dark:border-slate-700 dark:bg-slate-800/80">
      <div class="mb-3 flex flex-wrap items-center gap-2">
        <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">{{ plan.name }}</h3>
        <span class="rounded-full px-2.5 py-0.5 text-xs font-medium bg-emerald-50 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300">
          {{ periodLabel }}
        </span>
      </div>
      <div class="flex items-baseline gap-2">
        <span v-if="plan.original_price" class="text-sm line-through text-slate-400 dark:text-slate-500">
          ¥{{ plan.original_price }}
        </span>
        <span class="text-3xl font-bold text-emerald-600 dark:text-emerald-400">¥{{ plan.price }}</span>
      </div>
      <p v-if="plan.description" class="mt-3 text-sm leading-relaxed text-slate-500 dark:text-slate-400">
        {{ plan.description }}
      </p>
    </div>

    <!-- Payment method selector -->
    <div>
      <label class="mb-2 block text-sm font-medium text-slate-700 dark:text-slate-200">
        {{ t('payment.paymentMethod') }}
      </label>
      <div class="space-y-2">
        <button
          v-for="type in paymentTypes"
          :key="type"
          type="button"
          role="radio"
          :aria-checked="selectedPayment === type"
          :aria-label="getPaymentMethodLabel(type, t)"
          @click="selectedPayment = type"
          class="flex w-full items-center gap-3 rounded-xl border-2 px-4 py-3 text-left transition-all"
          :class="getMethodClass(type)"
        >
          <!-- Radio indicator -->
          <span
            class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full border-2"
            :class="selectedPayment === type ? getPaymentRadioBorderClass(type) : 'border-slate-300 dark:border-slate-600'"
          >
            <span v-if="selectedPayment === type" class="h-2.5 w-2.5 rounded-full" :style="{ backgroundColor: getPaymentBrandColor(type) }" />
          </span>
          <!-- Label -->
          <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ getPaymentMethodLabel(type, t) }}</span>
        </button>
      </div>
    </div>

    <!-- Amount due -->
    <div class="flex items-center justify-between rounded-xl border px-4 py-3 border-slate-200 bg-slate-50 dark:border-slate-700 dark:bg-slate-800/60">
      <span class="text-sm font-medium text-slate-600 dark:text-slate-300">{{ t('payment.confirm.amountDue') }}</span>
      <span class="text-xl font-bold text-emerald-500">¥{{ plan.price }}</span>
    </div>

    <!-- Submit button -->
    <button
      type="button"
      :disabled="!selectedPayment || loading"
      @click="handleSubmit"
      class="w-full rounded-xl py-3 text-sm font-bold text-white transition-colors"
      :class="
        selectedPayment && !loading
          ? 'bg-emerald-500 hover:bg-emerald-600 active:bg-emerald-700'
          : 'cursor-not-allowed bg-slate-200 text-slate-400 dark:bg-slate-700 dark:text-slate-400'
      "
    >
      {{ loading ? t('payment.processing') : t('payment.confirm.buyNow') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { PaymentSubscriptionPlan } from '@/types'
import { formatPeriodLabel, getPaymentMethodLabel, getPaymentBrandColor, getPaymentRadioBorderClass, getPaymentConfirmSelectedClass } from '@/utils/payment'

const props = defineProps<{
  plan: PaymentSubscriptionPlan
  paymentTypes: string[]
  loading: boolean
}>()

const emit = defineEmits<{
  back: []
  submit: [paymentType: string]
}>()

const { t } = useI18n()

const selectedPayment = ref(props.paymentTypes[0] || '')

const periodLabel = computed(() => formatPeriodLabel(props.plan.validity_days, props.plan.validity_unit, t))

function handleSubmit() {
  if (selectedPayment.value && !props.loading) {
    emit('submit', selectedPayment.value)
  }
}

function getMethodClass(type: string): string {
  if (selectedPayment.value === type) return getPaymentConfirmSelectedClass(type)
  return 'border-slate-200 hover:border-slate-300 bg-white dark:border-slate-700 dark:hover:border-slate-600 dark:bg-slate-800/60'
}
</script>
