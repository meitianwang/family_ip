<template>
  <AppLayout>
    <div class="px-4 pt-4">
      <PaymentAdminNav />
    </div>
    <TablePageLayout>
      <template #filters>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('payment.admin.channels') }}</h2>
          <div class="flex gap-2">
            <button @click="openSyncDialog" class="btn btn-secondary btn-sm">{{ t('payment.admin.syncFromGroups') }}</button>
            <button @click="openCreateDialog" class="btn btn-primary btn-sm">{{ t('common.create') }}</button>
          </div>
        </div>
      </template>

      <template #default>
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" role="status" aria-label="Loading" />
        </div>

        <div v-else-if="channels.length === 0" class="py-12 text-center text-gray-500 dark:text-slate-400">
          {{ t('payment.admin.noChannels') }}
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200 dark:border-slate-700">
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">ID</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.channelName') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.platform') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.channel.rate') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.groupId') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.enabled') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="ch in channels" :key="ch.id" class="border-b border-gray-100 dark:border-slate-800">
                <td class="px-4 py-3 text-gray-900 dark:text-slate-100">{{ ch.id }}</td>
                <td class="px-4 py-3 font-medium text-gray-900 dark:text-slate-100">{{ ch.name }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-slate-400">{{ ch.platform }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-slate-400">{{ ch.rate_multiplier }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-slate-400">{{ ch.group_id || '-' }}</td>
                <td class="px-4 py-3">
                  <span :class="ch.enabled ? 'text-green-600 dark:text-green-400' : 'text-gray-400'" role="img" :aria-label="ch.enabled ? t('payment.admin.enabled') : t('payment.admin.disabled')">
                    {{ ch.enabled ? '✓' : '✗' }}
                  </span>
                </td>
                <td class="px-4 py-3">
                  <div class="flex gap-2">
                    <button @click="openEditDialog(ch)" :aria-label="t('common.edit') + ' ' + ch.name" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">
                      {{ t('common.edit') }}
                    </button>
                    <button @click="handleDelete(ch.id)" :aria-label="t('common.delete') + ' ' + ch.name" class="text-sm text-red-600 hover:text-red-700 dark:text-red-400">
                      {{ t('common.delete') }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>
    </TablePageLayout>

    <!-- Create/Edit Dialog -->
    <BaseDialog :show="dialogOpen" @close="dialogOpen = false" :title="editingChannel ? t('common.edit') : t('common.create')">
      <div class="space-y-4">
        <div>
          <label for="ch-name" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.channelName') }}</label>
          <input id="ch-name" v-model="form.name" class="input w-full" maxlength="100" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="ch-platform" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.platform') }}</label>
            <input id="ch-platform" v-model="form.platform" class="input w-full" placeholder="default" maxlength="50" />
          </div>
          <div>
            <label for="ch-rate" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.channel.rate') }}</label>
            <input id="ch-rate" v-model="form.rate_multiplier" type="number" step="0.0001" min="0.0001" class="input w-full" placeholder="1.0" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="ch-group" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.groupId') }}</label>
            <input id="ch-group" v-model.number="form.group_id" type="number" class="input w-full" />
          </div>
          <div>
            <label for="ch-sort" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.sortOrder') }}</label>
            <input id="ch-sort" v-model.number="form.sort_order" type="number" class="input w-full" />
          </div>
        </div>
        <div>
          <label for="ch-desc" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.description') }}</label>
          <textarea id="ch-desc" v-model="form.description" rows="2" class="input w-full" maxlength="500" />
        </div>
        <div>
          <label for="ch-models" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.models') }}</label>
          <input id="ch-models" v-model="form.models" class="input w-full" placeholder="model-a,model-b" maxlength="1000" />
        </div>
        <div>
          <label for="ch-features" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.channel.features') }}</label>
          <input id="ch-features" v-model="form.features" class="input w-full" maxlength="500" />
        </div>
        <label for="ch-enabled" class="flex items-center gap-2 text-sm">
          <input id="ch-enabled" type="checkbox" v-model="form.enabled" class="rounded" />
          {{ t('payment.admin.enabled') }}
        </label>
        <div class="flex justify-end gap-3">
          <button @click="dialogOpen = false" class="btn btn-secondary">{{ t('common.cancel') }}</button>
          <button @click="handleSave" class="btn btn-primary" :disabled="formLoading">
            {{ formLoading ? t('payment.processing') : t('common.save') }}
          </button>
        </div>
      </div>
    </BaseDialog>

    <ConfirmDialog
      :show="showDeleteConfirm"
      :title="t('common.delete')"
      :message="t('common.confirmDelete')"
      :danger="true"
      @confirm="confirmDelete"
      @cancel="showDeleteConfirm = false"
    />

    <!-- Sync from Groups Dialog -->
    <BaseDialog :show="syncDialogOpen" @close="syncDialogOpen = false" :title="t('payment.admin.syncFromGroups')">
      <p class="mb-4 text-sm text-gray-500 dark:text-slate-400">{{ t('payment.admin.selectGroupsToSync') }}</p>
      <div v-if="syncLoading" class="flex items-center justify-center py-8">
        <div class="h-6 w-6 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" />
      </div>
      <div v-else-if="syncGroups.length === 0" class="py-8 text-center text-gray-500 dark:text-slate-400">
        {{ t('payment.admin.noGroupsToSync') }}
      </div>
      <div v-else class="max-h-80 space-y-2 overflow-y-auto">
        <label
          v-for="g in syncGroups"
          :key="g.id"
          class="flex items-center gap-3 rounded-lg border px-3 py-2 cursor-pointer"
          :class="syncExistingIds.has(g.id)
            ? 'border-gray-200 bg-gray-50 dark:border-slate-700 dark:bg-slate-800/50 opacity-60'
            : syncSelectedIds.has(g.id)
              ? 'border-blue-400 bg-blue-50 dark:border-blue-600 dark:bg-blue-950'
              : 'border-gray-200 dark:border-slate-700 hover:border-gray-300 dark:hover:border-slate-600'"
        >
          <input
            type="checkbox"
            :checked="syncSelectedIds.has(g.id)"
            :disabled="syncExistingIds.has(g.id)"
            @change="toggleSyncGroup(g.id)"
            class="rounded"
          />
          <div class="flex-1">
            <span class="text-sm font-medium text-gray-900 dark:text-slate-100">{{ g.name }}</span>
            <span class="ml-2 text-xs text-gray-500 dark:text-slate-400">#{{ g.id }}</span>
            <span v-if="g.platform" class="ml-2 text-xs text-gray-400 dark:text-slate-500">{{ g.platform }}</span>
          </div>
          <span v-if="syncExistingIds.has(g.id)" class="text-xs text-gray-400 dark:text-slate-500">{{ t('payment.admin.alreadyImported') }}</span>
        </label>
      </div>
      <div class="mt-4 flex justify-end gap-3">
        <button @click="syncDialogOpen = false" class="btn btn-secondary">{{ t('common.cancel') }}</button>
        <button @click="handleSync" :disabled="syncSelectedIds.size === 0 || syncImporting" class="btn btn-primary">
          {{ syncImporting ? t('payment.processing') : t('payment.admin.syncFromGroups') }} ({{ syncSelectedIds.size }})
        </button>
      </div>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import PaymentAdminNav from '@/components/admin/PaymentAdminNav.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { adminPayAPI } from '@/api/admin/pay'
// Groups API removed (AI-specific); sync feature disabled
const listGroups = async () => [] as any[]
import { useAppStore } from '@/stores'
import type { PaymentChannel } from '@/types'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const channels = ref<PaymentChannel[]>([])
const dialogOpen = ref(false)
const editingChannel = ref<PaymentChannel | null>(null)
const formLoading = ref(false)

interface ChannelForm {
  name: string
  platform: string
  rate_multiplier: string
  group_id: number | undefined
  sort_order: number
  description: string
  models: string
  features: string
  enabled: boolean
}

const form = reactive<ChannelForm>({
  name: '',
  platform: 'default',
  rate_multiplier: '1.0',
  group_id: undefined,
  sort_order: 0,
  description: '',
  models: '',
  features: '',
  enabled: true
})

async function loadChannels() {
  loading.value = true
  try {
    channels.value = await adminPayAPI.listChannels()
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    loading.value = false
  }
}

onMounted(loadChannels)

function resetForm() {
  form.name = ''
  form.platform = 'claude'
  form.rate_multiplier = '1.0'
  form.group_id = undefined
  form.sort_order = 0
  form.description = ''
  form.models = ''
  form.features = ''
  form.enabled = true
}

function openCreateDialog() {
  editingChannel.value = null
  resetForm()
  dialogOpen.value = true
}

function openEditDialog(ch: PaymentChannel) {
  editingChannel.value = ch
  form.name = ch.name
  form.platform = ch.platform
  form.rate_multiplier = ch.rate_multiplier
  form.group_id = ch.group_id
  form.sort_order = ch.sort_order
  form.description = ch.description || ''
  form.models = ch.models || ''
  form.features = ch.features || ''
  form.enabled = ch.enabled
  dialogOpen.value = true
}

async function handleSave() {
  // Trim string fields
  form.name = form.name.trim()
  form.platform = form.platform.trim()
  form.description = form.description.trim()
  form.models = form.models.trim()
  form.features = form.features.trim()

  if (!form.name) {
    appStore.showError(t('common.nameRequired'))
    return
  }
  const rate = parseFloat(form.rate_multiplier)
  if (isNaN(rate) || rate <= 0) {
    appStore.showError(t('payment.channel.rate') + ' > 0')
    return
  }
  if (form.group_id !== undefined && form.group_id !== null && form.group_id < 0) {
    form.group_id = undefined
  }

  formLoading.value = true
  try {
    const payload: Omit<PaymentChannel, 'id' | 'created_at' | 'updated_at'> = {
      name: form.name,
      platform: form.platform,
      rate_multiplier: form.rate_multiplier,
      group_id: form.group_id,
      sort_order: form.sort_order,
      description: form.description || undefined,
      models: form.models || undefined,
      features: form.features || undefined,
      enabled: form.enabled
    }
    if (editingChannel.value) {
      await adminPayAPI.updateChannel(editingChannel.value.id, payload)
    } else {
      await adminPayAPI.createChannel(payload)
    }
    dialogOpen.value = false
    loadChannels()
    appStore.showSuccess(t('common.saved'))
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    formLoading.value = false
  }
}

const showDeleteConfirm = ref(false)
const deletingId = ref<number | null>(null)

function handleDelete(id: number) {
  deletingId.value = id
  showDeleteConfirm.value = true
}

async function confirmDelete() {
  showDeleteConfirm.value = false
  if (!deletingId.value) return
  try {
    await adminPayAPI.deleteChannel(deletingId.value)
    loadChannels()
    appStore.showSuccess(t('common.deleted'))
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  }
}

// Sync from Groups
const syncDialogOpen = ref(false)
const syncLoading = ref(false)
const syncImporting = ref(false)
const syncGroups = ref<Array<{ id: number; name: string; platform?: string }>>([])
const syncSelectedIds = ref<Set<number>>(new Set())
const syncExistingIds = computed(() => new Set(channels.value.map(ch => ch.group_id).filter(Boolean)))

async function openSyncDialog() {
  syncDialogOpen.value = true
  syncSelectedIds.value = new Set()
  syncLoading.value = true
  try {
    const groups = await listGroups()
    syncGroups.value = (groups || []).map((g: any) => ({ id: g.id, name: g.name, platform: g.platform }))
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    syncLoading.value = false
  }
}

function toggleSyncGroup(id: number) {
  const s = new Set(syncSelectedIds.value)
  if (s.has(id)) s.delete(id)
  else s.add(id)
  syncSelectedIds.value = s
}

async function handleSync() {
  syncImporting.value = true
  try {
    const selected = syncGroups.value.filter(g => syncSelectedIds.value.has(g.id))
    for (const g of selected) {
      await adminPayAPI.createChannel({
        name: g.name,
        platform: g.platform || 'default',
        rate_multiplier: '1.0',
        group_id: g.id,
        sort_order: 0,
        enabled: true
      })
    }
    syncDialogOpen.value = false
    loadChannels()
    appStore.showSuccess(t('payment.admin.syncSuccess'))
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    syncImporting.value = false
  }
}
</script>
