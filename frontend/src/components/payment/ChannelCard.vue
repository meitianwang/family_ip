<template>
  <div
    class="flex flex-col rounded-2xl border p-6 transition-shadow hover:shadow-lg border-slate-200 bg-white dark:border-slate-700 dark:bg-slate-800/70"
  >
    <!-- Header -->
    <div class="mb-4">
      <div class="mb-3 flex items-center gap-2">
        <span
          class="rounded-full px-2.5 py-0.5 text-xs font-semibold"
          :class="platformBadgeClass"
        >
          {{ channel.platform }}
        </span>
        <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">{{ channel.name }}</h3>
      </div>

      <!-- Rate -->
      <div class="mb-3">
        <div class="flex items-baseline gap-2">
          <span class="text-sm text-slate-500 dark:text-slate-400">{{ t('payment.channel.rate') }}</span>
          <div class="flex items-baseline">
            <span class="text-xl font-bold" :class="accentClass">1</span>
            <span class="mx-1.5 text-lg text-slate-400 dark:text-slate-500">:</span>
            <span class="text-xl font-bold" :class="accentClass">{{ channel.rate_multiplier }}</span>
          </div>
        </div>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
          {{ t('payment.channel.quotaHint', { quota: usableQuota }) }}
        </p>
      </div>

      <!-- Description -->
      <p v-if="channel.description" class="text-sm leading-relaxed text-slate-500 dark:text-slate-400">
        {{ channel.description }}
      </p>
    </div>

    <!-- Models -->
    <div v-if="models.length > 0" class="mb-4">
      <p class="mb-2 text-xs text-slate-400 dark:text-slate-500">{{ t('payment.channel.supportedModels') }}</p>
      <div class="flex flex-wrap gap-1.5">
        <span
          v-for="model in models"
          :key="model"
          class="inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1 text-xs"
          :class="modelTagClass"
        >
          <span class="h-1.5 w-1.5 rounded-full" :class="dotClass" />
          {{ model }}
        </span>
      </div>
    </div>

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

    <!-- Top-up Button -->
    <button
      type="button"
      @click="emit('topUp')"
      class="mt-2 inline-flex w-full items-center justify-center gap-2 rounded-xl py-3 text-sm font-semibold text-white transition-colors"
      :class="buttonClass"
    >
      <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
      </svg>
      {{ t('payment.channel.topUpNow') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { PaymentChannel } from '@/types'
import { usePlatformStyle } from '@/composables/usePlatformStyle'

const props = defineProps<{
  channel: PaymentChannel
}>()

const emit = defineEmits<{
  topUp: []
}>()

const { t } = useI18n()

const usableQuota = computed(() => {
  const multiplier = parseFloat(props.channel.rate_multiplier)
  if (!multiplier || multiplier <= 0 || !isFinite(multiplier)) return '0.00'
  const quota = 1 / multiplier
  return isFinite(quota) ? quota.toFixed(2) : '0.00'
})

const models = computed(() =>
  props.channel.models ? props.channel.models.split(',').map((m) => m.trim()).filter(Boolean) : []
)

const features = computed(() =>
  props.channel.features ? props.channel.features.split(',').map((f) => f.trim()).filter(Boolean) : []
)

const { badgeClass: platformBadgeClass, accentClass, modelTagClass, dotClass, buttonClass } =
  usePlatformStyle(() => props.channel.platform)
</script>
