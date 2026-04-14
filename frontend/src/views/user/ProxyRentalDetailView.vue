<template>
  <AppLayout>
  <div class="max-w-2xl mx-auto px-4 py-8">
    <!-- Back -->
    <button
      @click="router.back()"
      class="mb-6 flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
    >
      ← 返回
    </button>

    <div v-if="loading" class="text-center py-12 text-gray-400">加载中...</div>

    <div v-else-if="!rental" class="text-center py-12 text-gray-400">租约不存在</div>

    <template v-else>
      <!-- Header -->
      <div class="border border-gray-200 dark:border-gray-700 rounded-2xl p-5 bg-white dark:bg-gray-800 mb-4">
        <div class="flex items-start justify-between mb-3">
          <div>
            <h1 class="text-xl font-bold text-gray-900 dark:text-white">
              {{ rental.node ? `${countryFlag(rental.node.country_code)} ${rental.node.city || rental.node.country}` : `节点 #${rental.node_id}` }}
            </h1>
            <div v-if="rental.node" class="text-sm text-gray-400 font-mono mt-0.5">{{ rental.node.ip_address }}</div>
          </div>
          <span :class="statusBadgeClass(rental.status)" class="text-xs px-3 py-1 rounded-full font-medium">
            {{ statusLabel(rental.status) }}
          </span>
        </div>

        <div class="grid grid-cols-2 gap-3 text-sm">
          <div>
            <div class="text-gray-500 dark:text-gray-400">套餐</div>
            <div class="font-medium text-gray-900 dark:text-white mt-0.5">{{ rental.product?.name ?? '-' }}</div>
          </div>
          <div>
            <div class="text-gray-500 dark:text-gray-400">到期时间</div>
            <div class="font-medium text-gray-900 dark:text-white mt-0.5">{{ rental.expires_at ? formatDate(rental.expires_at) : '-' }}</div>
          </div>
        </div>

        <!-- Traffic -->
        <div class="mt-4">
          <div class="flex justify-between text-sm mb-1">
            <span class="text-gray-500 dark:text-gray-400">流量使用</span>
            <span class="text-gray-900 dark:text-white">
              {{ formatBytes(rental.traffic_used_bytes) }}
              {{ rental.traffic_limit_bytes > 0 ? `/ ${formatBytes(rental.traffic_limit_bytes)}` : '/ 不限' }}
            </span>
          </div>
          <div v-if="rental.traffic_limit_bytes > 0" class="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all"
              :class="trafficPercent > 90 ? 'bg-red-500' : 'bg-blue-500'"
              :style="`width: ${Math.min(trafficPercent, 100)}%`"
            />
          </div>
        </div>

        <!-- Cancel button -->
        <button
          v-if="rental.status === 'pending_payment'"
          class="mt-4 w-full py-2 rounded-lg border border-red-300 text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/20 text-sm transition-colors disabled:opacity-60"
          :disabled="cancelling"
          @click="cancel"
        >
          {{ cancelling ? '取消中...' : '取消订单' }}
        </button>
      </div>

      <!-- Credentials -->
      <div v-if="rental.credential" class="border border-gray-200 dark:border-gray-700 rounded-2xl p-5 bg-white dark:bg-gray-800 mb-4">
        <h2 class="font-semibold text-gray-900 dark:text-white mb-4">连接凭证</h2>

        <!-- HTTP Proxy -->
        <div class="mb-5">
          <div class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">HTTP 代理</div>
          <div class="space-y-2">
            <CredentialRow label="主机" :value="`${rental.node?.ip_address ?? ''}:${rental.node?.http_port ?? ''}`" />
            <CredentialRow label="用户名" :value="rental.credential.http_username" />
            <CredentialRow label="密码" :value="rental.credential.http_password" secret />
          </div>
        </div>

        <!-- VLESS -->
        <div>
          <div class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">VLESS（Shadowrocket）</div>
          <div class="flex items-start gap-2">
            <div class="flex-1 min-w-0 bg-gray-50 dark:bg-gray-900 rounded-lg px-3 py-2 font-mono text-xs text-gray-700 dark:text-gray-300 break-all">
              {{ rental.credential.vless_link }}
            </div>
            <button
              class="shrink-0 px-3 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 text-white text-xs transition-colors"
              @click="copyText(rental.credential!.vless_link)"
            >
              {{ copied === 'vless' ? '已复制' : '复制' }}
            </button>
          </div>
        </div>
      </div>

      <!-- No credential yet -->
      <div v-else-if="rental.status === 'pending_payment'" class="border border-yellow-200 dark:border-yellow-800 rounded-2xl p-5 bg-yellow-50 dark:bg-yellow-900/20 text-sm text-yellow-700 dark:text-yellow-300">
        订单待支付，支付完成后凭证将自动生成。
      </div>
    </template>
  </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, defineComponent, h } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useRouter, useRoute } from 'vue-router'
import { getRental, cancelRental } from '@/api/proxy'
import type { ProxyRental } from '@/api/proxy'

const router = useRouter()
const route = useRoute()

const rental = ref<ProxyRental | null>(null)
const loading = ref(true)
const cancelling = ref(false)
const copied = ref('')

// Inline credential row component
const CredentialRow = defineComponent({
  props: {
    label: String,
    value: String,
    secret: Boolean,
  },
  setup(props) {
    const shown = ref(false)
    const copiedLocal = ref(false)

    function copy() {
      if (!props.value) return
      navigator.clipboard.writeText(props.value).then(() => {
        copiedLocal.value = true
        setTimeout(() => { copiedLocal.value = false }, 1500)
      })
    }

    return () => h('div', { class: 'flex items-center gap-2' }, [
      h('span', { class: 'text-xs text-gray-500 dark:text-gray-400 w-12 shrink-0' }, props.label),
      h('span', {
        class: 'flex-1 font-mono text-xs bg-gray-50 dark:bg-gray-900 rounded px-2 py-1 text-gray-800 dark:text-gray-200'
      }, props.secret && !shown.value ? '••••••••' : props.value),
      props.secret
        ? h('button', {
            class: 'text-xs text-gray-400 hover:text-gray-600 dark:hover:text-gray-200',
            onClick: () => { shown.value = !shown.value }
          }, shown.value ? '隐藏' : '显示')
        : null,
      h('button', {
        class: 'text-xs px-2 py-0.5 rounded bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600',
        onClick: copy
      }, copiedLocal.value ? '✓' : '复制'),
    ])
  }
})

const trafficPercent = computed(() => {
  if (!rental.value || !rental.value.traffic_limit_bytes) return 0
  return (rental.value.traffic_used_bytes / rental.value.traffic_limit_bytes) * 100
})

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
  return new Date(iso).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${units[i]}`
}

function copyText(text: string) {
  navigator.clipboard.writeText(text).then(() => {
    copied.value = 'vless'
    setTimeout(() => { copied.value = '' }, 1500)
  })
}

async function cancel() {
  if (!rental.value) return
  if (!confirm('确认取消此订单？')) return
  cancelling.value = true
  try {
    await cancelRental(rental.value.id)
    rental.value = await getRental(rental.value.id)
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '取消失败')
  } finally {
    cancelling.value = false
  }
}

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    rental.value = await getRental(id)
  } finally {
    loading.value = false
  }
})
</script>
