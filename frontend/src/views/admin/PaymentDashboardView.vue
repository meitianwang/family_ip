<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl px-4 py-6">
      <h1 class="mb-6 text-xl font-bold text-gray-900 dark:text-white">{{ t('payment.admin.dashboard') }}</h1>

      <!-- Time Range Selector -->
      <div class="mb-6 flex gap-2">
        <button
          v-for="d in [7, 30, 90]"
          :key="d"
          @click="selectDays(d)"
          :aria-label="t('payment.admin.showLastDays', { days: d })"
          :disabled="loading"
          class="rounded-lg px-4 py-2 text-sm font-medium transition-colors"
          :class="days === d
            ? 'bg-blue-600 text-white'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700'"
        >
          {{ d }}{{ t('payment.admin.days') }}
        </button>
      </div>

      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" role="status" aria-label="Loading" />
      </div>

      <template v-else-if="data">
        <!-- Stats Cards -->
        <div class="mb-8 grid grid-cols-2 gap-4 lg:grid-cols-4">
          <div class="rounded-xl border p-4 border-gray-200 dark:border-slate-700">
            <p class="text-sm text-gray-500 dark:text-slate-400">{{ t('payment.admin.todayRevenue') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">¥{{ data.today_amount }}</p>
          </div>
          <div class="rounded-xl border p-4 border-gray-200 dark:border-slate-700">
            <p class="text-sm text-gray-500 dark:text-slate-400">{{ t('payment.admin.todayOrders') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ data.today_order_count }}</p>
          </div>
          <div class="rounded-xl border p-4 border-gray-200 dark:border-slate-700">
            <p class="text-sm text-gray-500 dark:text-slate-400">{{ t('payment.admin.totalRevenue') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">¥{{ data.total_amount }}</p>
          </div>
          <div class="rounded-xl border p-4 border-gray-200 dark:border-slate-700">
            <p class="text-sm text-gray-500 dark:text-slate-400">{{ t('payment.admin.totalOrders') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ data.total_order_count }}</p>
          </div>
        </div>

        <!-- Daily Revenue Chart -->
        <div v-if="data.daily_series.length > 0" class="mb-8 rounded-xl border p-6 border-gray-200 dark:border-slate-700">
          <h3 class="mb-4 text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.dailyRevenue') }}</h3>
          <div class="h-64">
            <canvas ref="chartCanvas" />
          </div>
        </div>

        <!-- Payment Method Distribution -->
        <div v-if="data.payment_methods.length > 0" class="rounded-xl border p-6 border-gray-200 dark:border-slate-700">
          <h3 class="mb-4 text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.paymentMethods') }}</h3>
          <div class="space-y-3">
            <div
              v-for="pm in data.payment_methods"
              :key="pm.payment_type"
              class="flex items-center justify-between"
            >
              <span class="text-sm text-gray-700 dark:text-slate-300">{{ pm.payment_type }}</span>
              <div class="flex items-center gap-4">
                <span class="text-sm font-medium text-gray-900 dark:text-white">¥{{ pm.amount }}</span>
                <span class="text-sm text-gray-500 dark:text-slate-400">{{ pm.count }} {{ t('payment.admin.ordersUnit') }}</span>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { Chart, registerables } from 'chart.js'
import AppLayout from '@/components/layout/AppLayout.vue'
import { adminPayAPI } from '@/api/admin/pay'
import { useAppStore } from '@/stores'
import type { PaymentDashboardData } from '@/types'

Chart.register(...registerables)

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const days = ref(30)
const data = ref<PaymentDashboardData | null>(null)
const chartCanvas = ref<HTMLCanvasElement | null>(null)
let chartInstance: Chart | null = null

function selectDays(d: number) {
  if (loading.value || days.value === d) return
  days.value = d
  loadDashboard()
}

async function loadDashboard() {
  loading.value = true
  try {
    data.value = await adminPayAPI.getDashboard(days.value)
    await nextTick()
    renderChart()
  } catch (err: unknown) {
    data.value = null
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    loading.value = false
  }
}

onBeforeUnmount(() => {
  if (chartInstance) {
    chartInstance.destroy()
    chartInstance = null
  }
})

function renderChart() {
  if (!chartCanvas.value || !data.value?.daily_series.length) return
  if (chartInstance) {
    try {
      chartInstance.destroy()
    } catch (err: unknown) {
      // Chart.js may throw if canvas was detached during a re-render cycle.
      // Safe to ignore — a new chart will be created below.
      console.warn('[PaymentDashboard] Chart cleanup warning:', err)
    }
    chartInstance = null
  }

  const isDark = document.documentElement.classList.contains('dark')
  // Tailwind color tokens: blue-500, slate-400/500, slate-700/200
  const lineColor = isDark ? '#60a5fa' : '#3b82f6' // blue-400 / blue-500
  const tickColor = isDark ? '#94a3b8' : '#64748b' // slate-400 / slate-500
  const gridColor = isDark ? '#334155' : '#e2e8f0' // slate-700 / slate-200

  chartInstance = new Chart(chartCanvas.value, {
    type: 'line',
    data: {
      labels: data.value.daily_series.map((d) => d.date),
      datasets: [
        {
          label: t('payment.admin.revenue'),
          data: data.value.daily_series.map((d) => parseFloat(d.amount)),
          borderColor: lineColor,
          backgroundColor: lineColor + '1a', // 10% opacity
          fill: true,
          tension: 0.3
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        x: { ticks: { color: tickColor }, grid: { color: gridColor } },
        y: { ticks: { color: tickColor }, grid: { color: gridColor } }
      },
      plugins: {
        legend: { display: false }
      }
    }
  })
}

onMounted(loadDashboard)
</script>
