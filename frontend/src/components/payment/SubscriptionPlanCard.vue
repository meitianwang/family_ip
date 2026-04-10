<template>
  <div
    class="flex flex-col rounded-2xl border p-6 transition-shadow hover:shadow-lg border-slate-200 bg-white dark:border-slate-700 dark:bg-slate-800/70"
  >
    <!-- Header: Platform badge + Name + Period -->
    <div class="mb-4">
      <div class="mb-3 flex flex-wrap items-center gap-2">
        <span
          v-if="plan.platform"
          class="rounded-full px-2.5 py-0.5 text-xs font-semibold"
          :class="platformBadgeClass"
        >
          {{ plan.platform }}
        </span>
        <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">{{ plan.name }}</h3>
        <span class="rounded-full px-2.5 py-0.5 text-xs font-medium bg-emerald-50 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300">
          {{ periodLabel }}
        </span>
      </div>

      <!-- Price -->
      <div class="flex items-baseline gap-2">
        <span v-if="plan.original_price" class="text-sm line-through text-slate-400 dark:text-slate-500">
          ¥{{ plan.original_price }}
        </span>
        <span class="text-3xl font-bold" :class="accentClass">¥{{ plan.price }}</span>
        <span class="text-sm text-slate-500 dark:text-slate-400">{{ periodSuffix }}</span>
      </div>
    </div>

    <!-- Description -->
    <p v-if="plan.description" class="mb-4 text-sm leading-relaxed text-slate-500 dark:text-slate-400">
      {{ plan.description }}
    </p>

    <!-- Features -->
    <div v-if="features.length > 0" class="mb-5">
      <p class="mb-2 text-xs text-slate-400 dark:text-slate-500">{{ t('payment.channel.features') }}</p>
      <div class="flex flex-wrap gap-1.5">
        <span
          v-for="feature in features"
          :key="feature"
          class="rounded-md px-2 py-1 text-xs bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400"
        >
          {{ feature }}
        </span>
      </div>
    </div>

    <div class="flex-1" />

    <!-- Subscribe Button -->
    <button
      type="button"
      @click="emit('subscribe', plan.id)"
      class="mt-2 inline-flex w-full items-center justify-center gap-2 rounded-xl py-3 text-sm font-semibold text-white transition-colors"
      :class="buttonClass"
    >
      <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
      </svg>
      {{ t('payment.plan.subscribeNow') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { PaymentSubscriptionPlan } from '@/types'
import { usePlatformStyle } from '@/composables/usePlatformStyle'
import { formatPeriodLabel, formatPeriodSuffix } from '@/utils/payment'

const props = defineProps<{
  plan: PaymentSubscriptionPlan
}>()

const emit = defineEmits<{
  subscribe: [planId: number]
}>()

const { t } = useI18n()

const features = computed(() =>
  props.plan.features ? props.plan.features.split(',').map((f) => f.trim()).filter(Boolean) : []
)

const periodLabel = computed(() => formatPeriodLabel(props.plan.validity_days, props.plan.validity_unit, t))
const periodSuffix = computed(() => formatPeriodSuffix(props.plan.validity_days, props.plan.validity_unit, t))

const { badgeClass: platformBadgeClass, accentClass, buttonClass } =
  usePlatformStyle(() => props.plan.platform || '')
</script>
