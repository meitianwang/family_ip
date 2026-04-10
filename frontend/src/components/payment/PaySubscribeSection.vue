<template>
  <!-- Subscription Confirm -->
  <SubscriptionConfirm
    v-if="selectedPlan"
    :plan="selectedPlan"
    :payment-types="enabledPaymentTypes"
    :loading="loading"
    @back="selectedPlan = null"
    @submit="(type) => emit('submit', type)"
  />

  <!-- Plans Grid -->
  <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
    <SubscriptionPlanCard
      v-for="plan in plans"
      :key="plan.id"
      :plan="plan"
      @subscribe="selectPlan"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import SubscriptionConfirm from '@/components/payment/SubscriptionConfirm.vue'
import SubscriptionPlanCard from '@/components/payment/SubscriptionPlanCard.vue'
import type { PaymentSubscriptionPlan } from '@/types'

const props = defineProps<{
  plans: PaymentSubscriptionPlan[]
  enabledPaymentTypes: string[]
  loading: boolean
}>()

const emit = defineEmits<{
  submit: [paymentType: string]
}>()

const selectedPlan = ref<PaymentSubscriptionPlan | null>(null)

function selectPlan(planId: number) {
  const plan = props.plans.find((p) => p.id === planId)
  if (plan) selectedPlan.value = plan
}

defineExpose({
  selectedPlan,
  resetSelection: () => { selectedPlan.value = null }
})
</script>
