<template>
  <AppLayout>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">我的代理</h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">管理你已购买的家庭 IP 代理</p>
      </div>
      <RouterLink
        to="/proxy/marketplace"
        class="px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium transition-colors"
      >
        购买代理
      </RouterLink>
    </div>

    <!-- Status filter -->
    <div class="flex gap-2 mb-6 flex-wrap">
      <button
        v-for="s in statusTabs"
        :key="s.value"
        class="px-3 py-1.5 rounded-full text-sm transition-colors"
        :class="filterStatus === s.value
          ? 'bg-blue-600 text-white'
          : 'bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700'"
        @click="filterStatus = s.value; reload()"
      >
        {{ s.label }}
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-12 text-gray-400">加载中...</div>

    <!-- Empty -->
    <div v-else-if="rentals.length === 0" class="text-center py-16 text-gray-400">
      <div class="text-4xl mb-3">📭</div>
      <p>暂无代理记录</p>
    </div>

    <!-- List -->
    <div v-else class="space-y-3">
      <div
        v-for="rental in rentals"
        :key="rental.id"
        class="border border-gray-200 dark:border-gray-700 rounded-xl p-4 bg-white dark:bg-gray-800 hover:border-blue-300 transition-colors cursor-pointer"
        @click="router.push(`/proxy/rentals/${rental.id}`)"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span class="text-base font-semibold text-gray-900 dark:text-white">
                {{ rental.node ? `${countryFlag(rental.node.country_code)} ${rental.node.city || rental.node.country}` : `节点 #${rental.node_id}` }}
              </span>
              <span :class="statusBadgeClass(rental.status)" class="text-xs px-2 py-0.5 rounded-full">
                {{ statusLabel(rental.status) }}
              </span>
            </div>
            <div v-if="rental.node" class="text-xs text-gray-400 font-mono mb-2">
              {{ rental.node.ip_address }}
            </div>
            <div class="flex flex-wrap gap-x-4 gap-y-1 text-xs text-gray-500 dark:text-gray-400">
              <span v-if="rental.product">{{ rental.product.name }}</span>
              <span v-if="rental.expires_at">到期：{{ formatDate(rental.expires_at) }}</span>
            </div>
          </div>
          <div class="text-right shrink-0">
            <div class="text-sm text-gray-500 dark:text-gray-400 mb-1">流量</div>
            <div class="text-sm font-medium text-gray-900 dark:text-white">
              {{ formatBytes(rental.traffic_used_bytes) }}
              <span v-if="rental.traffic_limit_bytes > 0">
                / {{ formatBytes(rental.traffic_limit_bytes) }}
              </span>
              <span v-else>/ 不限</span>
            </div>
            <div v-if="rental.traffic_limit_bytes > 0" class="mt-1 h-1.5 w-24 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
              <div
                class="h-full rounded-full"
                :class="trafficPercent(rental) > 90 ? 'bg-red-500' : 'bg-blue-500'"
                :style="`width: ${Math.min(trafficPercent(rental), 100)}%`"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="mt-6 flex justify-center gap-2">
      <button
        v-for="p in totalPages"
        :key="p"
        class="px-3 py-1.5 rounded-lg text-sm transition-colors"
        :class="p === currentPage
          ? 'bg-blue-600 text-white'
          : 'bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300 hover:bg-gray-200'"
        @click="currentPage = p; reload()"
      >{{ p }}</button>
    </div>
  </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { listRentals } from '@/api/proxy'
import type { ProxyRental } from '@/api/proxy'

const router = useRouter()

const rentals = ref<ProxyRental[]>([])
const loading = ref(true)
const filterStatus = ref('')
const currentPage = ref(1)
const totalPages = ref(1)

const statusTabs = [
  { value: '', label: '全部' },
  { value: 'active', label: '使用中' },
  { value: 'pending_payment', label: '待支付' },
  { value: 'expired', label: '已到期' },
  { value: 'cancelled', label: '已取消' },
]

async function reload() {
  loading.value = true
  try {
    const result = await listRentals(currentPage.value, 20, filterStatus.value)
    rentals.value = result.items
    totalPages.value = result.pages
  } finally {
    loading.value = false
  }
}

function countryFlag(code: string): string {
  if (!code || code.length !== 2) return '🌐'
  return code.toUpperCase().replace(/./g, c => String.fromCodePoint(c.charCodeAt(0) + 0x1F1A5))
}

function statusLabel(status: string): string {
  const map: Record<string, string> = {
    active: '使用中',
    pending_payment: '待支付',
    expired: '已到期',
    cancelled: '已取消',
  }
  return map[status] ?? status
}

function statusBadgeClass(status: string): string {
  const map: Record<string, string> = {
    active: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300',
    pending_payment: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-300',
    expired: 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400',
    cancelled: 'bg-red-100 text-red-600 dark:bg-red-900/40 dark:text-red-300',
  }
  return map[status] ?? 'bg-gray-100 text-gray-500'
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('zh-CN')
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${units[i]}`
}

function trafficPercent(rental: ProxyRental): number {
  if (!rental.traffic_limit_bytes) return 0
  return (rental.traffic_used_bytes / rental.traffic_limit_bytes) * 100
}

onMounted(reload)
</script>
