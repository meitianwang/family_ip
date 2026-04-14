<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <select v-model="filterStatus" @change="reload" class="input w-36">
            <option value="">全部状态</option>
            <option value="active">使用中</option>
            <option value="pending_payment">待支付</option>
            <option value="expired">已到期</option>
            <option value="cancelled">已取消</option>
          </select>
          <input
            v-model="filterUserId"
            @change="reload"
            type="number"
            min="1"
            placeholder="用户 ID"
            class="input w-32"
          />
        </div>
      </template>

      <template #default>
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" />
        </div>

        <div v-else-if="rentals.length === 0" class="py-12 text-center text-gray-500 dark:text-slate-400">暂无租约</div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200 dark:border-slate-700">
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">ID</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">用户</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">节点</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">套餐</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">状态</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">到期</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">流量</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="r in rentals" :key="r.id" class="border-b border-gray-100 dark:border-slate-800">
                <td class="px-3 py-3 text-gray-900 dark:text-slate-100">{{ r.id }}</td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">#{{ r.user_id }}</td>
                <td class="px-3 py-3 font-mono text-xs text-gray-600 dark:text-slate-400">
                  {{ r.node ? r.node.ip_address : `#${r.node_id}` }}
                </td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">{{ r.product?.name ?? `#${r.product_id}` }}</td>
                <td class="px-3 py-3">
                  <span :class="statusClass(r.status)" class="text-xs px-2 py-0.5 rounded-full">{{ statusLabel(r.status) }}</span>
                </td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">{{ r.expires_at ? formatDate(r.expires_at) : '-' }}</td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">
                  {{ formatBytes(r.traffic_used_bytes) }}
                  <span v-if="r.traffic_limit_bytes > 0"> / {{ formatBytes(r.traffic_limit_bytes) }}</span>
                </td>
                <td class="px-3 py-3">
                  <div class="flex gap-2 flex-wrap">
                    <button @click="openDetail(r)" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">详情</button>
                    <button
                      v-if="r.status === 'active'"
                      @click="openTrafficUpdate(r)"
                      class="text-sm text-green-600 hover:text-green-700 dark:text-green-400"
                    >流量</button>
                    <button
                      v-if="r.status === 'active'"
                      @click="handleForceExpire(r.id)"
                      class="text-sm text-red-600 hover:text-red-700 dark:text-red-400"
                    >强制到期</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="px-3 py-3 flex gap-2 justify-end">
          <button
            v-for="p in totalPages" :key="p"
            @click="currentPage = p; reload()"
            class="px-2.5 py-1 rounded text-xs"
            :class="p === currentPage ? 'bg-primary-600 text-white' : 'bg-gray-100 dark:bg-slate-700 text-gray-600 dark:text-slate-300'"
          >{{ p }}</button>
        </div>
      </template>
    </TablePageLayout>

    <!-- Detail Dialog -->
    <BaseDialog :show="detailOpen" @close="detailOpen = false" title="租约详情">
      <div v-if="selectedRental" class="space-y-4 text-sm">
        <div class="grid grid-cols-2 gap-3">
          <div><div class="text-gray-500 dark:text-slate-400">ID</div><div class="font-medium mt-0.5">{{ selectedRental.id }}</div></div>
          <div><div class="text-gray-500 dark:text-slate-400">用户</div><div class="font-medium mt-0.5">#{{ selectedRental.user_id }}</div></div>
          <div><div class="text-gray-500 dark:text-slate-400">节点 IP</div><div class="font-mono mt-0.5">{{ selectedRental.node?.ip_address ?? '-' }}</div></div>
          <div><div class="text-gray-500 dark:text-slate-400">套餐</div><div class="mt-0.5">{{ selectedRental.product?.name ?? '-' }}</div></div>
          <div><div class="text-gray-500 dark:text-slate-400">状态</div><div class="mt-0.5">{{ statusLabel(selectedRental.status) }}</div></div>
          <div><div class="text-gray-500 dark:text-slate-400">到期时间</div><div class="mt-0.5">{{ selectedRental.expires_at ? formatDate(selectedRental.expires_at) : '-' }}</div></div>
        </div>

        <div v-if="selectedRental.credential" class="border-t border-gray-200 dark:border-slate-700 pt-4">
          <div class="text-xs font-semibold text-gray-500 dark:text-slate-400 uppercase tracking-wide mb-3">连接凭证</div>
          <div class="space-y-2 font-mono text-xs">
            <div><span class="text-gray-400">HTTP 用户名：</span>{{ selectedRental.credential.http_username }}</div>
            <div><span class="text-gray-400">HTTP 密码：</span>{{ selectedRental.credential.http_password }}</div>
            <div><span class="text-gray-400">VLESS UUID：</span>{{ selectedRental.credential.vless_uuid }}</div>
            <div class="break-all"><span class="text-gray-400">VLESS Link：</span>{{ selectedRental.credential.vless_link }}</div>
          </div>
        </div>
      </div>
      <template #footer>
        <button @click="detailOpen = false" class="btn btn-secondary">关闭</button>
      </template>
    </BaseDialog>

    <!-- Traffic Update Dialog -->
    <BaseDialog :show="trafficOpen" @close="trafficOpen = false" title="更新流量">
      <div class="space-y-4">
        <p class="text-sm text-gray-500 dark:text-slate-400">为租约 #{{ selectedRental?.id }} 添加/扣减流量（GB）</p>
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">流量变化（GB，负数为扣减）</label>
          <input v-model.number="trafficDeltaGb" type="number" step="0.1" class="input w-full" />
        </div>
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">备注</label>
          <input v-model="trafficNote" class="input w-full" />
        </div>
      </div>
      <template #footer>
        <button @click="trafficOpen = false" class="btn btn-secondary">取消</button>
        <button @click="submitTraffic" :disabled="trafficSaving" class="btn btn-primary">
          {{ trafficSaving ? '提交中...' : '提交' }}
        </button>
      </template>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { proxyAdminAPI } from '@/api/admin/proxy'
import type { AdminProxyRental } from '@/api/admin/proxy'

const rentals = ref<AdminProxyRental[]>([])
const loading = ref(true)
const filterStatus = ref('')
const filterUserId = ref<number | ''>('')
const currentPage = ref(1)
const totalPages = ref(1)

const detailOpen = ref(false)
const trafficOpen = ref(false)
const selectedRental = ref<AdminProxyRental | null>(null)
const trafficDeltaGb = ref(0)
const trafficNote = ref('')
const trafficSaving = ref(false)

async function reload() {
  loading.value = true
  try {
    const result = await proxyAdminAPI.listRentals({
      status: filterStatus.value || undefined,
      user_id: filterUserId.value || undefined,
      page: currentPage.value,
      page_size: 20,
    })
    rentals.value = result.items
    totalPages.value = result.pages
  } finally {
    loading.value = false
  }
}

async function openDetail(rental: AdminProxyRental) {
  try {
    selectedRental.value = await proxyAdminAPI.getRental(rental.id)
  } catch {
    selectedRental.value = rental
  }
  detailOpen.value = true
}

function openTrafficUpdate(rental: AdminProxyRental) {
  selectedRental.value = rental
  trafficDeltaGb.value = 0
  trafficNote.value = ''
  trafficOpen.value = true
}

async function submitTraffic() {
  if (!selectedRental.value) return
  trafficSaving.value = true
  try {
    await proxyAdminAPI.updateTraffic(selectedRental.value.id, trafficDeltaGb.value, trafficNote.value)
    trafficOpen.value = false
    await reload()
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '操作失败')
  } finally {
    trafficSaving.value = false
  }
}

async function handleForceExpire(id: number) {
  if (!confirm('确认强制使该租约到期？')) return
  try {
    await proxyAdminAPI.forceExpire(id)
    await reload()
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '操作失败')
  }
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

function statusClass(status: string): string {
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

onMounted(reload)
</script>
