<template>
  <AppLayout>
    <div class="mx-auto max-w-4xl px-4 py-6">
      <!-- Loading -->
      <div v-if="pageLoading" class="flex items-center justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" role="status" aria-label="Loading" />
      </div>

      <!-- Balance Disabled & No Plans -->
      <div v-else-if="config?.balance_payment_disabled && !hasPlans" class="flex items-center justify-center py-12">
        <div class="max-w-md text-center">
          <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 dark:bg-dark-700">
            <Icon name="creditCard" size="lg" class="text-gray-400" />
          </div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('payment.notEnabled') }}</h3>
        </div>
      </div>

      <!-- Main Content -->
      <template v-else>
        <!-- Step: Form -->
        <div v-if="step === 'form'">
          <!-- Tabs -->
          <div v-if="showTabs" class="mb-6 flex items-center justify-center gap-2">
            <button
              v-if="showTopUpTab"
              @click="mainTab = 'topup'"
              class="rounded-full px-6 py-2 text-sm font-medium transition-colors"
              :class="mainTab === 'topup'
                ? 'bg-blue-600 text-white shadow-sm'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700'"
            >
              {{ t('payment.tabs.topUp') }}
            </button>
            <button
              v-if="showSubscribeTab"
              @click="mainTab = 'subscribe'"
              class="rounded-full px-6 py-2 text-sm font-medium transition-colors"
              :class="mainTab === 'subscribe'
                ? 'bg-blue-600 text-white shadow-sm'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700'"
            >
              {{ t('payment.tabs.subscribe') }}
            </button>
          </div>

          <PayTopUpSection
            v-if="mainTab === 'topup'"
            ref="topUpRef"
            :channels="channels"
            :user-id="authStore.user!.id"
            :user-name="authStore.user?.username"
            :user-balance="authStore.user?.balance"
            :enabled-payment-types="config?.enabled_payment_types || []"
            :method-limits="methodLimitsMap"
            :min-amount="parseFloat(config?.min_recharge_amount || '1')"
            :max-amount="parseFloat(config?.max_recharge_amount || '1000')"
            :loading="orderLoading"
            :pending-blocked="pendingBlocked"
            :pending-count="pendingCount"
            @submit="handleTopUpSubmit"
          />

          <PaySubscribeSection
            v-if="mainTab === 'subscribe'"
            ref="subscribeRef"
            :plans="plans"
            :enabled-payment-types="config?.enabled_payment_types || []"
            :loading="orderLoading"
            @submit="handleSubscribeSubmit"
          />

          <!-- Recent Orders -->
          <div v-if="recentOrders.length > 0" class="mt-8">
            <div class="mb-3 flex items-center justify-between">
              <h3 class="text-sm font-medium text-slate-700 dark:text-slate-300">{{ t('payment.recentOrders') }}</h3>
              <router-link to="/purchase/orders" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">
                {{ t('payment.viewAll') }}
              </router-link>
            </div>
            <div class="space-y-2">
              <div
                v-for="order in recentOrders"
                :key="order.id"
                class="flex items-center justify-between rounded-lg border px-4 py-3 border-slate-200 dark:border-slate-700"
              >
                <div>
                  <span class="text-sm font-medium text-slate-900 dark:text-slate-100">¥{{ order.amount }}</span>
                  <span class="ml-2 text-xs text-slate-500 dark:text-slate-400">{{ formatPaymentDate(order.created_at) }}</span>
                </div>
                <span
                  class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium"
                  :class="getPaymentStatusBadgeClass(order.status)"
                >
                  <span class="h-1.5 w-1.5 rounded-full" :class="getStatusDotClass(order.status)" />
                  {{ t(`payment.orderStatus.${order.status}`) }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Step: Paying -->
        <div v-else-if="step === 'paying' && orderResult" class="mx-auto max-w-md">
          <PaymentQRCode
            :order-id="orderResult.order_id"
            :pay-url="orderResult.pay_url"
            :qr-code="orderResult.qr_code"
            :client-secret="orderResult.client_secret"
            :payment-type="orderResult.payment_type"
            :amount="parseFloat(orderResult.amount)"
            :pay-amount="orderResult.pay_amount ? parseFloat(orderResult.pay_amount) : undefined"
            :expires-at="orderResult.expires_at"
            :is-mobile="isMobile"
            @status-change="handleStatusChange"
            @back="resetToForm"
          />
        </div>

        <!-- Step: Result -->
        <div v-else-if="step === 'result'" class="mx-auto max-w-md">
          <OrderStatusCard
            :order-id="orderResult?.order_id || 0"
            :status="finalStatus"
            @back="resetToForm"
            @state-change="handleFinalStateChange"
          />
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import PayTopUpSection from '@/components/payment/PayTopUpSection.vue'
import PaySubscribeSection from '@/components/payment/PaySubscribeSection.vue'
import PaymentQRCode from '@/components/payment/PaymentQRCode.vue'
import OrderStatusCard from '@/components/payment/OrderStatusCard.vue'
import { payAPI } from '@/api/pay'
import { formatPaymentDate, getPaymentStatusBadgeClass } from '@/utils/payment'
import type {
  PaymentConfig,
  PaymentChannel,
  PaymentSubscriptionPlan,
  CreateOrderResponse,
  UserPaymentOrder,
  MethodLimit
} from '@/types'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const pageLoading = ref(true)
const orderLoading = ref(false)
const isMobile = ref(false)

const step = ref<'form' | 'paying' | 'result'>('form')
const mainTab = ref<'topup' | 'subscribe'>('topup')

const topUpRef = ref<InstanceType<typeof PayTopUpSection> | null>(null)
const subscribeRef = ref<InstanceType<typeof PaySubscribeSection> | null>(null)

const config = ref<PaymentConfig | null>(null)
const channels = ref<PaymentChannel[]>([])
const plans = ref<PaymentSubscriptionPlan[]>([])
const orderResult = ref<CreateOrderResponse | null>(null)
const finalStatus = ref('pending')
const recentOrders = ref<UserPaymentOrder[]>([])
const pendingCount = ref(0)

const hasPlans = computed(() => plans.value.length > 0)
const showTopUpTab = computed(() => !config.value?.balance_payment_disabled)
const showSubscribeTab = computed(() => hasPlans.value)
const showTabs = computed(() => showTopUpTab.value && showSubscribeTab.value)

const pendingBlocked = computed(() => {
  const max = config.value?.max_pending_orders || 0
  return max > 0 && pendingCount.value >= max
})

const methodLimitsMap = computed(() => {
  if (!config.value?.method_limits) return undefined
  const map: Record<string, MethodLimit> = {}
  for (const ml of config.value.method_limits) {
    map[ml.payment_type] = ml
  }
  return map
})

onMounted(async () => {
  isMobile.value = /Android|iPhone|iPad|iPod/i.test(navigator.userAgent)

  try {
    const [configData, channelsData, plansData, ordersData] = await Promise.all([
      payAPI.getConfig(),
      payAPI.listChannels(),
      payAPI.listPlans(),
      payAPI.listOrders(1, 5)
    ])

    config.value = configData
    channels.value = channelsData
    plans.value = plansData
    recentOrders.value = ordersData.items || []
    pendingCount.value = (ordersData.items || []).filter((o) => o.status === 'pending').length

    if (configData.balance_payment_disabled) {
      mainTab.value = 'subscribe'
    }
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    pageLoading.value = false
  }
})

async function handleTopUpSubmit(amount: number, paymentType: string) {
  if (!config.value?.enabled_payment_types?.includes(paymentType)) {
    appStore.showError(t('payment.methodLimitReached'))
    return
  }
  orderLoading.value = true
  try {
    const result = await payAPI.createOrder({
      amount: amount.toFixed(2),
      payment_type: paymentType,
      order_type: 'balance',
      is_mobile: isMobile.value
    })
    orderResult.value = result
    step.value = 'paying'
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('payment.createOrderFailed'))
  } finally {
    orderLoading.value = false
  }
}

async function handleSubscribeSubmit(paymentType: string) {
  const plan = subscribeRef.value?.selectedPlan
  if (!plan) return
  if (!config.value?.enabled_payment_types?.includes(paymentType)) {
    appStore.showError(t('payment.methodLimitReached'))
    return
  }
  orderLoading.value = true
  try {
    const result = await payAPI.createOrder({
      amount: plan.price,
      payment_type: paymentType,
      order_type: 'subscription',
      plan_id: plan.id,
      is_mobile: isMobile.value
    })
    orderResult.value = result
    step.value = 'paying'
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('payment.createOrderFailed'))
  } finally {
    orderLoading.value = false
  }
}

function handleStatusChange(status: string) {
  finalStatus.value = status
  step.value = 'result'
}

function handleFinalStateChange(status: string) {
  finalStatus.value = status
}

function resetToForm() {
  step.value = 'form'
  orderResult.value = null
  topUpRef.value?.resetForm()
  subscribeRef.value?.resetSelection()
  authStore.refreshUser()
  payAPI.listOrders(1, 5).then((data) => {
    recentOrders.value = data.items || []
    pendingCount.value = (data.items || []).filter((o) => o.status === 'pending').length
  }).catch(() => {
    // Non-critical: order list reload failure doesn't block user flow
  })
}

function getStatusDotClass(status: string): string {
  if (status === 'completed') return 'bg-green-500'
  if (status === 'pending') return 'bg-yellow-500'
  if (status === 'paid' || status === 'recharging') return 'bg-blue-500'
  if (status === 'failed') return 'bg-red-500'
  return 'bg-gray-400'
}
</script>
