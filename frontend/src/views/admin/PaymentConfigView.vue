<template>
  <AppLayout>
    <div class="mx-auto max-w-4xl px-4 py-6">
      <h1 class="mb-6 text-xl font-bold text-gray-900 dark:text-white">{{ t('payment.admin.paymentConfig') }}</h1>

      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" role="status" aria-label="Loading" />
      </div>

      <template v-else>
        <div class="space-y-6">
          <!-- Basic Settings -->
          <div class="rounded-xl border p-6 border-gray-200 dark:border-slate-700">
            <h3 class="mb-4 text-sm font-semibold text-gray-700 dark:text-slate-300">{{ t('payment.admin.basicSettings') }}</h3>
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label for="cfg-enabled-types" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.admin.enabledTypes') }}</label>
                <input id="cfg-enabled-types" v-model="configs['pay_enabled_payment_types']" class="input w-full" placeholder="alipay,wxpay,stripe" maxlength="200" />
              </div>
              <div>
                <label for="cfg-min-amount" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.admin.minAmount') }}</label>
                <input id="cfg-min-amount" v-model="configs['pay_min_recharge_amount']" type="number" step="0.01" min="0" class="input w-full" placeholder="1" />
              </div>
              <div>
                <label for="cfg-max-amount" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.admin.maxAmount') }}</label>
                <input id="cfg-max-amount" v-model="configs['pay_max_recharge_amount']" type="number" step="0.01" min="0" class="input w-full" placeholder="5000" />
              </div>
              <div>
                <label for="cfg-max-daily" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.admin.maxDailyAmount') }}</label>
                <input id="cfg-max-daily" v-model="configs['pay_max_daily_recharge_amount']" type="number" step="0.01" min="0" class="input w-full" placeholder="10000" />
              </div>
              <div>
                <label for="cfg-timeout" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.admin.orderTimeout') }}</label>
                <input id="cfg-timeout" v-model="configs['pay_order_timeout_minutes']" type="number" step="1" min="1" class="input w-full" placeholder="30" />
              </div>
              <div>
                <label for="cfg-max-pending" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.admin.maxPendingOrders') }}</label>
                <input id="cfg-max-pending" v-model="configs['pay_max_pending_orders']" type="number" step="1" min="0" class="input w-full" placeholder="5" />
              </div>
              <div>
                <label for="cfg-product-name" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.admin.productName') }}</label>
                <input id="cfg-product-name" v-model="configs['pay_product_name']" class="input w-full" maxlength="255" />
              </div>
              <div>
                <label for="cfg-lb-strategy" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.admin.loadBalanceStrategy') }}</label>
                <select id="cfg-lb-strategy" v-model="configs['pay_load_balance_strategy']" class="input w-full">
                  <option value="round_robin">{{ t('payment.admin.strategyRoundRobin') }}</option>
                  <option value="min_amount">{{ t('payment.admin.strategyMinAmount') }}</option>
                </select>
              </div>
            </div>
            <div class="mt-4 flex items-center gap-4">
              <label for="cfg-balance-disabled" class="flex items-center gap-2 text-sm">
                <input id="cfg-balance-disabled" type="checkbox" v-model="balanceDisabled" class="rounded" />
                {{ t('payment.admin.balanceDisabled') }}
              </label>
              <label for="cfg-auto-refund" class="flex items-center gap-2 text-sm">
                <input id="cfg-auto-refund" type="checkbox" v-model="autoRefund" class="rounded" />
                {{ t('payment.admin.autoRefund') }}
              </label>
            </div>
          </div>

          <!-- Fee Rates -->
          <div class="rounded-xl border p-6 border-gray-200 dark:border-slate-700">
            <h3 class="mb-4 text-sm font-semibold text-gray-700 dark:text-slate-300">{{ t('payment.admin.feeRates') }}</h3>
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
              <div>
                <label for="cfg-fee-alipay" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.alipay') }} (%)</label>
                <input id="cfg-fee-alipay" v-model="configs['pay_fee_rate_alipay']" type="number" step="0.01" min="0" max="100" class="input w-full" placeholder="0" />
              </div>
              <div>
                <label for="cfg-fee-wxpay" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.wechatPay') }} (%)</label>
                <input id="cfg-fee-wxpay" v-model="configs['pay_fee_rate_wxpay']" type="number" step="0.01" min="0" max="100" class="input w-full" placeholder="0" />
              </div>
              <div>
                <label for="cfg-fee-stripe" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">Stripe (%)</label>
                <input id="cfg-fee-stripe" v-model="configs['pay_fee_rate_stripe']" type="number" step="0.01" min="0" max="100" class="input w-full" placeholder="0" />
              </div>
            </div>
          </div>

          <!-- Daily Limits per Method -->
          <div class="rounded-xl border p-6 border-gray-200 dark:border-slate-700">
            <h3 class="mb-4 text-sm font-semibold text-gray-700 dark:text-slate-300">{{ t('payment.admin.dailyLimits') }}</h3>
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
              <div>
                <label for="cfg-daily-alipay" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.alipay') }}</label>
                <input id="cfg-daily-alipay" v-model="configs['pay_max_daily_amount_alipay']" type="number" step="0.01" min="0" class="input w-full" placeholder="0 (unlimited)" />
              </div>
              <div>
                <label for="cfg-daily-wxpay" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">{{ t('payment.wechatPay') }}</label>
                <input id="cfg-daily-wxpay" v-model="configs['pay_max_daily_amount_wxpay']" type="number" step="0.01" min="0" class="input w-full" placeholder="0 (unlimited)" />
              </div>
              <div>
                <label for="cfg-daily-stripe" class="mb-1 block text-sm text-gray-600 dark:text-slate-400">Stripe</label>
                <input id="cfg-daily-stripe" v-model="configs['pay_max_daily_amount_stripe']" type="number" step="0.01" min="0" class="input w-full" placeholder="0 (unlimited)" />
              </div>
            </div>
          </div>

          <!-- Save Button -->
          <div class="flex justify-end">
            <button @click="saveConfig" class="btn btn-primary" :disabled="saving">
              {{ saving ? t('payment.processing') : t('common.save') }}
            </button>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import { adminPayAPI } from '@/api/admin/pay'
import { useAppStore } from '@/stores'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const saving = ref(false)
const configs = ref<Record<string, string>>({})

const balanceDisabled = computed({
  get: () => configs.value['pay_balance_payment_disabled'] === 'true',
  set: (v) => { configs.value['pay_balance_payment_disabled'] = v ? 'true' : 'false' }
})

const autoRefund = computed({
  get: () => configs.value['pay_auto_refund_enabled'] === 'true',
  set: (v) => { configs.value['pay_auto_refund_enabled'] = v ? 'true' : 'false' }
})

onMounted(async () => {
  try {
    const data = await adminPayAPI.getConfig()
    configs.value = { ...data.configs }
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    loading.value = false
  }
})

const VALID_PAYMENT_TYPES = new Set(['alipay', 'alipay_direct', 'wxpay', 'wxpay_direct', 'stripe'])

function validateAmounts(): string | null {
  const minAmount = parseFloat(configs.value['pay_min_recharge_amount'] || '0')
  const maxAmount = parseFloat(configs.value['pay_max_recharge_amount'] || '0')
  const maxDaily = parseFloat(configs.value['pay_max_daily_recharge_amount'] || '0')

  if (minAmount < 0 || maxAmount < 0 || maxDaily < 0) return t('payment.admin.amountsMustBePositive')
  if (minAmount > 0 && maxAmount > 0 && minAmount > maxAmount) return `${t('payment.admin.minAmount')} ≤ ${t('payment.admin.maxAmount')}`
  if (maxAmount > 0 && maxDaily > 0 && maxAmount > maxDaily) return `${t('payment.admin.maxAmount')} ≤ ${t('payment.admin.maxDailyAmount')}`
  return null
}

function validatePaymentTypes(): string | null {
  const typesStr = configs.value['pay_enabled_payment_types'] || ''
  if (!typesStr) return null
  const types = typesStr.split(',').map((s) => s.trim()).filter(Boolean)
  const invalid = types.filter((s) => !VALID_PAYMENT_TYPES.has(s))
  if (invalid.length > 0) return `Unknown payment types: ${invalid.join(', ')}`
  return null
}

function validateFeeRates(): string | null {
  for (const key of Object.keys(configs.value)) {
    if (key.startsWith('pay_fee_rate_')) {
      const rate = parseFloat(configs.value[key] || '0')
      if (rate < 0 || rate > 100) return 'Fee rate must be between 0 and 100%'
    }
  }
  return null
}

async function saveConfig() {
  const amountErr = validateAmounts()
  if (amountErr) { appStore.showError(amountErr); return }
  const typeErr = validatePaymentTypes()
  if (typeErr) { appStore.showError(typeErr); return }
  const feeErr = validateFeeRates()
  if (feeErr) { appStore.showError(feeErr); return }

  saving.value = true
  try {
    const toSave: Record<string, string> = {}
    for (const [k, v] of Object.entries(configs.value)) {
      if (v !== undefined && v !== null) {
        toSave[k] = typeof v === 'string' ? v.trim() : String(v)
      }
    }
    await adminPayAPI.updateConfig(toSave)
    appStore.showSuccess(t('payment.admin.configSaved'))
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    saving.value = false
  }
}
</script>
