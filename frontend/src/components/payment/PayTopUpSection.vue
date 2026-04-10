<template>
  <!-- Channel Grid -->
  <div v-if="channels.length > 0 && !showForm" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
    <ChannelCard
      v-for="channel in channels"
      :key="channel.id"
      :channel="channel"
      @top-up="showForm = true"
    />
  </div>

  <!-- Payment Form -->
  <div v-if="showForm || channels.length === 0" class="mx-auto max-w-lg">
    <button
      v-if="channels.length > 0"
      @click="showForm = false"
      class="mb-4 flex items-center gap-1 text-sm text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200"
    >
      <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
      </svg>
      {{ t('payment.backToChannels') }}
    </button>

    <PaymentForm
      :user-id="userId"
      :user-name="userName"
      :user-balance="userBalance"
      :enabled-payment-types="enabledPaymentTypes"
      :method-limits="methodLimits"
      :min-amount="minAmount"
      :max-amount="maxAmount"
      :loading="loading"
      :pending-blocked="pendingBlocked"
      :pending-count="pendingCount"
      @submit="(amount, type) => emit('submit', amount, type)"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ChannelCard from '@/components/payment/ChannelCard.vue'
import PaymentForm from '@/components/payment/PaymentForm.vue'
import type { PaymentChannel, MethodLimit } from '@/types'

defineProps<{
  channels: PaymentChannel[]
  userId: number
  userName?: string
  userBalance?: number
  enabledPaymentTypes: string[]
  methodLimits?: Record<string, MethodLimit>
  minAmount: number
  maxAmount: number
  loading: boolean
  pendingBlocked: boolean
  pendingCount: number
}>()

const emit = defineEmits<{
  submit: [amount: number, paymentType: string]
}>()

const { t } = useI18n()
const showForm = ref(false)

defineExpose({ resetForm: () => { showForm.value = false } })
</script>
