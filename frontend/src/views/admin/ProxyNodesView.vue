<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex items-center justify-between gap-3 flex-wrap">
          <div class="flex gap-2 flex-wrap">
            <select v-model="filterStatus" @change="reload" class="input w-36">
              <option value="">全部状态</option>
              <option value="available">可用</option>
              <option value="rented">已租出</option>
              <option value="offline">下线</option>
            </select>
          </div>
          <button @click="openCreate" class="btn btn-primary btn-sm">添加节点</button>
        </div>
      </template>

      <template #default>
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" />
        </div>

        <div v-else-if="nodes.length === 0" class="py-12 text-center text-gray-500 dark:text-slate-400">暂无节点</div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200 dark:border-slate-700">
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">ID</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">IP 地址</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">位置</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">ISP</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">HTTP 端口</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">VLESS 端口</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">状态</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="node in nodes" :key="node.id" class="border-b border-gray-100 dark:border-slate-800">
                <td class="px-3 py-3 text-gray-900 dark:text-slate-100">{{ node.id }}</td>
                <td class="px-3 py-3 font-mono text-gray-900 dark:text-slate-100">{{ node.ip_address }}</td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">
                  {{ countryFlag(node.country_code) }} {{ node.city || node.country || node.country_code }}
                </td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">{{ node.isp || '-' }}</td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">{{ node.http_port }}</td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">{{ node.vless_port }}</td>
                <td class="px-3 py-3">
                  <span :class="statusClass(node.status)" class="text-xs px-2 py-0.5 rounded-full">{{ statusLabel(node.status) }}</span>
                </td>
                <td class="px-3 py-3">
                  <div class="flex gap-3">
                    <button @click="openEdit(node)" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">编辑</button>
                    <button @click="handleDelete(node.id)" class="text-sm text-red-600 hover:text-red-700 dark:text-red-400">删除</button>
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

    <!-- Create/Edit Dialog -->
    <BaseDialog :show="dialogOpen" @close="closeDialog" :title="editingNode ? '编辑节点' : '添加节点'">
      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div class="col-span-2">
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">IP 地址 *</label>
            <input v-model="form.ip_address" class="input w-full" placeholder="1.2.3.4" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">国家</label>
            <input v-model="form.country" class="input w-full" placeholder="United States" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">国家代码</label>
            <input v-model="form.country_code" class="input w-full" placeholder="US" maxlength="2" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">城市</label>
            <input v-model="form.city" class="input w-full" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">ISP</label>
            <input v-model="form.isp" class="input w-full" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">HTTP 端口</label>
            <input v-model.number="form.http_port" type="number" class="input w-full" placeholder="3128" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">VLESS 端口</label>
            <input v-model.number="form.vless_port" type="number" class="input w-full" placeholder="443" />
          </div>
        </div>

        <div class="border-t border-gray-200 dark:border-slate-700 pt-4">
          <div class="text-xs font-semibold text-gray-500 dark:text-slate-400 uppercase tracking-wide mb-3">VLESS 配置</div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">传输协议</label>
              <select v-model="form.vless_network" class="input w-full">
                <option value="tcp">TCP</option>
                <option value="ws">WebSocket</option>
                <option value="grpc">gRPC</option>
              </select>
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">TLS</label>
              <select v-model="form.vless_tls" class="input w-full">
                <option :value="false">关闭</option>
                <option :value="true">开启</option>
              </select>
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">SNI</label>
              <input v-model="form.vless_sni" class="input w-full" placeholder="example.com" />
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">WS Path</label>
              <input v-model="form.vless_ws_path" class="input w-full" placeholder="/" />
            </div>
          </div>
        </div>

        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">标签（逗号分隔）</label>
          <input v-model="tagsInput" class="input w-full" placeholder="住宅, 美国, 高速" />
        </div>
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">状态</label>
          <select v-model="form.status" class="input w-full">
            <option value="available">可用</option>
            <option value="offline">下线</option>
          </select>
        </div>
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">备注</label>
          <input v-model="form.description" class="input w-full" />
        </div>
      </div>
      <template #footer>
        <button @click="closeDialog" class="btn btn-secondary">取消</button>
        <button @click="save" :disabled="saving" class="btn btn-primary">{{ saving ? '保存中...' : '保存' }}</button>
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
import type { AdminProxyNode } from '@/api/admin/proxy'

const nodes = ref<AdminProxyNode[]>([])
const loading = ref(true)
const filterStatus = ref('')
const currentPage = ref(1)
const totalPages = ref(1)

const dialogOpen = ref(false)
const editingNode = ref<AdminProxyNode | null>(null)
const saving = ref(false)
const tagsInput = ref('')

const defaultForm = (): Omit<AdminProxyNode, 'id'> => ({
  ip_address: '',
  country: '',
  country_code: '',
  city: '',
  isp: '',
  http_port: 3128,
  vless_port: 443,
  vless_network: 'tcp',
  vless_tls: false,
  vless_sni: '',
  vless_ws_path: '',
  tags: [],
  status: 'available',
  description: '',
})

const form = ref(defaultForm())

async function reload() {
  loading.value = true
  try {
    const result = await proxyAdminAPI.listNodes({
      status: filterStatus.value || undefined,
      page: currentPage.value,
      page_size: 20,
    })
    nodes.value = result.items
    totalPages.value = result.pages
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingNode.value = null
  form.value = defaultForm()
  tagsInput.value = ''
  dialogOpen.value = true
}

function openEdit(node: AdminProxyNode) {
  editingNode.value = node
  form.value = { ...node }
  tagsInput.value = node.tags.join(', ')
  dialogOpen.value = true
}

function closeDialog() {
  dialogOpen.value = false
}

async function save() {
  saving.value = true
  try {
    form.value.tags = tagsInput.value.split(',').map(t => t.trim()).filter(Boolean)
    if (editingNode.value) {
      await proxyAdminAPI.updateNode(editingNode.value.id, form.value)
    } else {
      await proxyAdminAPI.createNode(form.value)
    }
    closeDialog()
    await reload()
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm('确认删除此节点？')) return
  try {
    await proxyAdminAPI.deleteNode(id)
    await reload()
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '删除失败')
  }
}

function countryFlag(code: string): string {
  if (!code || code.length !== 2) return '🌐'
  return code.toUpperCase().replace(/./g, c => String.fromCodePoint(c.charCodeAt(0) + 0x1F1A5))
}

function statusLabel(status: string): string {
  const map: Record<string, string> = { available: '可用', rented: '已租出', offline: '下线' }
  return map[status] ?? status
}

function statusClass(status: string): string {
  const map: Record<string, string> = {
    available: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300',
    rented: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300',
    offline: 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400',
  }
  return map[status] ?? 'bg-gray-100 text-gray-500'
}

onMounted(reload)
</script>
