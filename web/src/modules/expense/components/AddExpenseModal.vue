<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Message from 'primevue/message'
import { useCreateExpense, useUpdateExpense } from '../composables/useExpenses'
import { useAllTags, useCreateTag, useAddTagToExpense, useRemoveTagFromExpense } from '../composables/useTags'
import { getTagsForExpense } from '../api/tag'
import type { Tag, RecurrenceInterval, RecurringSchedule, Expense } from '../models/expense'

const props = defineProps<{ visible: boolean; expense?: Expense | null }>()
const emit = defineEmits<{ 'update:visible': [value: boolean] }>()

const editing = computed(() => !!props.expense)

const TAG_COLORS = ['#f5c518', '#f97316', '#22c55e', '#3b82f6', '#a855f7', '#ec4899', '#14b8a6', '#ef4444']

const name = ref('')
const cost = ref<string>('')
const date = ref<string>(new Date().toISOString().slice(0, 10))
const selectedTagIds = ref<Set<number>>(new Set())
const originalTagIds = ref<Set<number>>(new Set())
const newTagName = ref('')
const recurring = ref(false)
const frequency = ref<RecurrenceInterval>('monthly')
const recurrenceDay = ref<number>(new Date().getDate())
const submitError = ref<string | null>(null)
const fieldError = ref<string | null>(null)
const tagsOpen = ref(false)

const { mutateAsync: createExpense, isPending: creating } = useCreateExpense()
const { mutateAsync: updateExpense, isPending: updating } = useUpdateExpense()
const isPending = computed(() => creating.value || updating.value)
const { data: allTags } = useAllTags()
const { mutateAsync: createTag } = useCreateTag()
const { mutateAsync: addTag } = useAddTagToExpense()
const { mutateAsync: removeTag } = useRemoveTagFromExpense()

const tagList = computed<Tag[]>(() => allTags.value ?? [])

function colorFor(id: number) {
  return TAG_COLORS[id % TAG_COLORS.length]
}

const selectedTags = computed<Tag[]>(() =>
  tagList.value.filter((t) => selectedTagIds.value.has(t.id)),
)

const departmentLabel = computed(() => {
  if (selectedTags.value.length === 0) return 'Select…'
  if (selectedTags.value.length === 1) return selectedTags.value[0]!.name
  return `${selectedTags.value.length} selected`
})

const departmentCode = computed(() => {
  if (!selectedTags.value.length) return 'NONE'
  return selectedTags.value[0]!.name.slice(0, 4).toUpperCase()
})

const FREQ_OPTIONS: { value: RecurrenceInterval; label: string }[] = [
  { value: 'monthly', label: 'monthly' },
  { value: 'yearly', label: 'yearly' },
]

const costNum = computed(() => parseFloat(cost.value) || 0)

const displayDate = computed(() => {
  if (!date.value) return '—'
  const [y, m, d] = date.value.split('-')
  return `${m}/${d}/${y?.slice(2) ?? ''}`
})

// Stable txn id per modal session
const txnId = ref('TXN-' + Math.random().toString(36).slice(2, 8).toUpperCase())

const stampDate = ref('')
const stampTime = ref('')
function refreshStamps() {
  const n = new Date()
  stampDate.value = n.toLocaleDateString('en-US', { month: '2-digit', day: '2-digit', year: '2-digit' })
  stampTime.value = n.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
}
refreshStamps()

const barcodeBars = computed(() => {
  const arr: number[] = []
  let seed = 23
  for (let i = 0; i < 48; i++) {
    seed = (seed * 9301 + 49297) % 233280
    const r = seed / 233280
    arr.push(r < 0.5 ? 1 : r < 0.85 ? 2 : 3)
  }
  return arr
})

watch(
  () => props.visible,
  async (v) => {
    if (!v) return
    refreshStamps()
    txnId.value = 'TXN-' + Math.random().toString(36).slice(2, 8).toUpperCase()
    if (props.expense) {
      name.value = props.expense.name
      cost.value = String(props.expense.cost ?? '')
      if (props.expense.recurringSchedule) {
        recurring.value = true
        frequency.value = props.expense.recurringSchedule.interval
        recurrenceDay.value = props.expense.recurringSchedule.dayOfMonth
        date.value = props.expense.recurringSchedule.startDate
      }
      try {
        const tags = await getTagsForExpense(props.expense.id)
        const ids = new Set(tags.map((t) => t.id))
        originalTagIds.value = ids
        selectedTagIds.value = new Set(ids)
      } catch {
        originalTagIds.value = new Set()
      }
    }
  },
)

function toggleTag(id: number) {
  const next = new Set(selectedTagIds.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  selectedTagIds.value = next
}

async function handleCreateTag() {
  const n = newTagName.value.trim()
  if (!n) return
  const tag = await createTag(n)
  const next = new Set(selectedTagIds.value)
  next.add(tag.id)
  selectedTagIds.value = next
  newTagName.value = ''
}

function reset() {
  name.value = ''
  cost.value = ''
  date.value = new Date().toISOString().slice(0, 10)
  selectedTagIds.value = new Set()
  originalTagIds.value = new Set()
  newTagName.value = ''
  recurring.value = false
  frequency.value = 'monthly'
  recurrenceDay.value = new Date().getDate()
  submitError.value = null
  fieldError.value = null
  tagsOpen.value = false
}

function close() {
  reset()
  emit('update:visible', false)
}

async function submit() {
  submitError.value = null
  fieldError.value = null
  if (!name.value.trim()) {
    fieldError.value = 'Name required'
    return
  }
  if (isNaN(costNum.value) || costNum.value <= 0) {
    fieldError.value = 'Cost required'
    return
  }
  try {
    let recurringSchedule: RecurringSchedule | undefined
    if (recurring.value) {
      recurringSchedule = {
        interval: frequency.value,
        dayOfMonth: recurrenceDay.value,
        startDate: date.value,
      }
    }
    const payload = {
      name: name.value.trim(),
      amount: 1,
      cost: costNum.value,
      recurringSchedule,
    }
    const expense = props.expense
      ? await updateExpense({ id: props.expense.id, ...payload })
      : await createExpense(payload)
    const toAdd = [...selectedTagIds.value].filter((id) => !originalTagIds.value.has(id))
    const toRemove = [...originalTagIds.value].filter((id) => !selectedTagIds.value.has(id))
    await Promise.all([
      ...toAdd.map((tagId) => addTag({ expenseId: expense.id, tagId })),
      ...toRemove.map((tagId) => removeTag({ expenseId: expense.id, tagId })),
    ])
    close()
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Failed to save expense'
  }
}

const dayOptions = Array.from({ length: 28 }, (_, i) => i + 1)
</script>

<template>
  <Dialog
    :visible="props.visible"
    @update:visible="emit('update:visible', $event)"
    modal
    :closable="false"
    :show-header="false"
    :style="{ width: '26rem', maxWidth: '95vw' }"
    content-class="!p-0 !bg-transparent"
    :pt="{ root: { class: '!shadow-none !bg-transparent !border-0' }, mask: { class: '!bg-black/70 backdrop-blur-[2px]' } }"
  >
    <form
      class="exp-receipt relative w-full font-mono px-7 pt-9 pb-6"
      @submit.prevent="submit"
      @click.stop
    >
      <button
        type="button"
        aria-label="Close"
        class="absolute top-3.5 right-3.5 w-7 h-7 rounded-full border-[1.5px] border-current flex items-center justify-center text-[18px] leading-none cursor-pointer hover:bg-[var(--exp-receipt-text)] hover:text-[var(--exp-receipt-bg)] transition"
        @click="close"
      >
        ×
      </button>

      <!-- Header -->
      <div class="text-center">
        <div class="inline-flex items-center gap-2 font-bold text-[16px] tracking-[0.12em]">
          <img src="/icon.png" class="w-6 h-6" alt="" />
          <span>JAFA</span>
        </div>
        <div class="mt-1.5 text-[10px] tracking-[0.18em] uppercase opacity-70">
          ★ EXPENSE TRACKING CO. ★
        </div>
        <div class="mt-2.5 flex justify-between text-[10px] tracking-wide opacity-70">
          <span>{{ stampDate }} {{ stampTime }}</span>
          <span>#{{ txnId }}</span>
        </div>
      </div>

      <hr class="exp-divider" />

      <div class="text-center my-1">
        <div class="text-[17px] font-extrabold tracking-[0.08em] uppercase">
          {{ editing ? 'Edit Transaction' : 'New Transaction' }}
        </div>
        <div class="mt-1 text-[10.5px] tracking-[0.16em] uppercase opacity-60">
          {{ editing ? 'Update the entry below' : 'Please fill in the details' }}
        </div>
      </div>

      <hr class="exp-divider" />

      <!-- Item Description -->
      <div class="mb-3">
        <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase opacity-70 mb-1">
          <span>1× Item Description</span>
          <span class="font-bold">REQUIRED</span>
        </label>
        <input
          v-model="name"
          type="text"
          placeholder="e.g. Weekly groceries"
          class="exp-input w-full"
          autofocus
        />
      </div>

      <!-- Unit Price + Date -->
      <div class="grid grid-cols-2 gap-4 mb-3">
        <div>
          <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase opacity-70 mb-1">
            <span>Unit Price</span>
            <span class="font-bold">EUR</span>
          </label>
          <div class="relative">
            <span class="absolute left-0 bottom-1.5 font-bold text-[14px]">€</span>
            <input
              v-model="cost"
              type="number"
              step="0.01"
              min="0"
              placeholder="0.00"
              class="exp-input w-full !pl-4 tabular-nums"
            />
          </div>
        </div>
        <div>
          <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase opacity-70 mb-1">
            <span>Date Stamped</span>
            <span class="font-bold">{{ displayDate }}</span>
          </label>
          <input v-model="date" type="date" class="exp-input w-full" />
        </div>
      </div>

      <!-- Department / Tags -->
      <div class="mb-3 relative">
        <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase opacity-70 mb-1">
          <span>Department</span>
          <span class="font-bold">{{ departmentCode }}</span>
        </label>
        <button
          type="button"
          class="exp-stamp-trigger w-full flex items-center gap-2.5 px-3 py-2 border-[1.5px] border-current bg-transparent font-mono text-[12px] tracking-[0.06em] uppercase font-bold cursor-pointer hover:bg-[var(--exp-receipt-text)]/[0.04]"
          @click="tagsOpen = !tagsOpen"
        >
          <span
            v-if="selectedTags.length"
            class="w-2.5 h-2.5 rounded-[2px] shrink-0"
            :style="{ background: colorFor(selectedTags[0]!.id) }"
          />
          <span class="flex-1 text-left truncate">{{ departmentLabel }}</span>
          <span class="text-[11px] opacity-70">{{ tagsOpen ? '▴' : '▾' }}</span>
        </button>

        <div v-if="tagsOpen" class="mt-2 p-2.5 border-[1.5px] border-dashed border-[var(--exp-receipt-border)] bg-[var(--exp-receipt-text)]/[0.03]">
          <div v-if="tagList.length" class="grid grid-cols-2 gap-1.5">
            <button
              v-for="t in tagList"
              :key="t.id"
              type="button"
              class="exp-stamp"
              :class="{ sel: selectedTagIds.has(t.id) }"
              @click="toggleTag(t.id)"
            >
              <span class="w-2 h-2 rounded-[2px] shrink-0" :style="{ background: colorFor(t.id) }" />
              <span class="truncate">{{ t.name }}</span>
            </button>
          </div>
          <div v-else class="text-[10.5px] tracking-[0.1em] uppercase opacity-60 py-1">No tags yet</div>

          <div class="flex gap-1.5 mt-2 pt-2 border-t-[1.5px] border-dashed border-[var(--exp-receipt-border)]">
            <InputText
              v-model="newTagName"
              placeholder="New tag…"
              size="small"
              class="flex-1 !text-[12px]"
              @keyup.enter.prevent="handleCreateTag"
            />
            <button
              type="button"
              class="exp-add-tag px-3 text-[10px] tracking-[0.14em] uppercase font-bold cursor-pointer border-[1.5px] border-current bg-transparent hover:bg-[var(--exp-receipt-text)] hover:text-[var(--exp-receipt-bg)] transition disabled:opacity-40 disabled:cursor-not-allowed"
              :disabled="!newTagName.trim()"
              @click="handleCreateTag"
            >
              Add
            </button>
          </div>
        </div>
      </div>

      <hr class="exp-divider" />

      <!-- Recurring -->
      <div
        class="exp-recurring p-3.5 border-[1.5px] transition"
        :class="recurring ? 'exp-recurring-on border-solid' : 'border-dashed border-[var(--exp-receipt-border)] bg-[var(--exp-receipt-text)]/[0.025]'"
      >
        <label class="flex items-center justify-between gap-3 cursor-pointer select-none">
          <div class="flex flex-col gap-0.5">
            <div class="text-[12.5px] font-extrabold tracking-[0.1em] uppercase">Recurring Charge</div>
            <div class="text-[10px] tracking-[0.1em] uppercase opacity-60">Auto-repeats on a schedule</div>
          </div>
          <span class="relative inline-block w-9 h-5 shrink-0">
            <input v-model="recurring" type="checkbox" class="opacity-0 w-0 h-0 peer" />
            <span
              class="absolute inset-0 rounded-full transition cursor-pointer"
              :class="recurring ? 'bg-[var(--jafa-accent)]' : 'bg-[var(--exp-receipt-border)]'"
            />
            <span
              class="absolute top-0.5 left-0.5 w-4 h-4 rounded-full bg-[var(--exp-receipt-bg)] transition"
              :class="{ 'translate-x-4': recurring }"
            />
          </span>
        </label>

        <div v-if="recurring" class="mt-2.5 pt-2.5 border-t-[1.5px] border-dashed border-[var(--exp-receipt-border)]">
          <div class="grid grid-cols-2 gap-1.5 mb-2.5">
            <button
              v-for="f in FREQ_OPTIONS"
              :key="f.value"
              type="button"
              class="exp-freq-stamp"
              :class="{ sel: frequency === f.value }"
              @click="frequency = f.value"
            >
              {{ f.label }}
            </button>
          </div>
          <label class="block text-[10px] tracking-[0.16em] uppercase opacity-70 mb-1">Day of Month</label>
          <select v-model.number="recurrenceDay" class="exp-input w-full font-mono">
            <option v-for="d in dayOptions" :key="d" :value="d">{{ d }}</option>
          </select>
        </div>
      </div>

      <hr class="exp-divider" />

      <!-- Subtotal lines -->
      <div class="flex flex-col">
        <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
          <span class="text-[10.5px] tracking-[0.14em] uppercase opacity-60">Subtotal</span>
          <span class="font-semibold tabular-nums">€{{ costNum.toFixed(2) }}</span>
        </div>
        <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
          <span class="text-[10.5px] tracking-[0.14em] uppercase opacity-60">Category Tax</span>
          <span class="font-semibold tabular-nums">€0.00</span>
        </div>
        <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
          <span class="text-[10.5px] tracking-[0.14em] uppercase opacity-60">Schedule</span>
          <span class="font-semibold uppercase">{{ recurring ? frequency : 'one-time' }}</span>
        </div>
      </div>

      <hr class="exp-divider" />

      <!-- TOTAL -->
      <div class="flex justify-between items-center gap-3 py-1 text-[16px] font-extrabold tracking-[0.04em]">
        <span>TOTAL</span>
        <span class="tabular-nums">€{{ costNum.toFixed(2) }}</span>
      </div>

      <Message v-if="fieldError" severity="error" :closable="false" class="!mt-2">{{ fieldError }}</Message>
      <Message v-if="submitError" severity="error" :closable="false" class="!mt-2">{{ submitError }}</Message>

      <button
        type="submit"
        :disabled="isPending"
        class="exp-tender w-full mt-3 flex items-center justify-center gap-1.5 py-3 px-4 font-bold text-[13px] tracking-[0.18em] uppercase transition disabled:cursor-not-allowed"
      >
        <span v-if="isPending">PROCESSING…</span>
        <template v-else>
          <span>{{ editing ? 'Update Entry' : 'Ring it up' }}</span>
          <span>→</span>
        </template>
      </button>

      <button
        type="button"
        class="exp-void w-full mt-1.5 py-2 text-[10.5px] tracking-[0.18em] uppercase font-semibold opacity-70 hover:opacity-100 hover:text-[var(--jafa-accent)] transition cursor-pointer"
        @click="close"
      >
        ✕ Void transaction
      </button>

      <hr class="exp-divider mt-4" />

      <div class="text-center text-[11px] tracking-[0.16em] uppercase font-bold">
        ★ Thank you for tracking ★
      </div>
      <div class="text-center text-[9.5px] tracking-[0.1em] opacity-60 mt-1">
        Keep this receipt for your records
      </div>

      <div class="flex justify-center gap-[1.5px] my-3 h-[34px] items-stretch">
        <span
          v-for="(w, i) in barcodeBars"
          :key="i"
          class="block bg-[var(--exp-receipt-text)]"
          :style="{ width: `${w}px` }"
        />
      </div>
      <div class="text-center text-[10.5px] tracking-[0.3em] opacity-70">
        JAFA · {{ departmentCode }} · {{ txnId.split('-')[1] }}
      </div>
    </form>
  </Dialog>
</template>

<style scoped>
.exp-receipt {
  background: var(--exp-receipt-bg);
  color: var(--exp-receipt-text);
  box-shadow:
    0 30px 80px -20px rgba(0, 0, 0, 0.7),
    0 50px 120px -40px rgba(245, 197, 24, 0.15);
  background-image:
    radial-gradient(color-mix(in srgb, var(--exp-receipt-text) 4%, transparent) 1px, transparent 1px),
    radial-gradient(color-mix(in srgb, var(--exp-receipt-text) 4%, transparent) 1px, transparent 1px);
  background-size: 3px 3px, 7px 7px;
  background-position: 0 0, 1.5px 1.5px;
}

.exp-divider {
  border: none;
  border-top: 1.5px dashed var(--exp-receipt-border);
  margin: 10px 0;
}

.exp-input {
  background: transparent;
  border: none;
  border-bottom: 1.5px dashed var(--exp-receipt-border);
  padding: 4px 2px 6px;
  font-family: inherit;
  font-size: 13px;
  color: var(--exp-receipt-text);
  outline: none;
}
.exp-input:focus {
  border-bottom-color: var(--exp-receipt-text);
}
.exp-input::placeholder {
  color: color-mix(in srgb, var(--exp-receipt-text) 40%, transparent);
}

.exp-stamp {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 7px 9px;
  background: var(--exp-receipt-bg);
  border: 1.5px solid var(--exp-receipt-border);
  color: var(--exp-receipt-text);
  font-family: inherit;
  font-size: 10.5px;
  letter-spacing: 0.08em;
  cursor: pointer;
  text-align: left;
  text-transform: uppercase;
  font-weight: 600;
  transition: all 0.12s;
}
.exp-stamp:hover {
  border-color: var(--exp-receipt-text);
}
.exp-stamp.sel {
  background: var(--exp-receipt-text);
  border-color: var(--exp-receipt-text);
  color: var(--exp-receipt-bg);
}

.exp-freq-stamp {
  padding: 7px 4px;
  background: transparent;
  border: 1.5px solid var(--exp-receipt-text);
  color: var(--exp-receipt-text);
  font-family: inherit;
  font-size: 9.5px;
  letter-spacing: 0.14em;
  cursor: pointer;
  text-transform: uppercase;
  font-weight: 700;
  transition: all 0.1s;
}
.exp-freq-stamp:hover {
  background: color-mix(in srgb, var(--exp-receipt-text) 5%, transparent);
}
.exp-freq-stamp.sel {
  background: var(--exp-receipt-text);
  color: var(--exp-receipt-bg);
}

.exp-recurring-on {
  background: color-mix(in srgb, var(--jafa-accent) 8%, transparent) !important;
  border-color: var(--jafa-accent) !important;
}

.exp-tender {
  background: var(--exp-receipt-text);
  color: var(--exp-receipt-bg);
  border: none;
  font-family: inherit;
  cursor: pointer;
}
.exp-tender:hover:not(:disabled) {
  opacity: 0.85;
}
.exp-tender:disabled {
  background: var(--exp-receipt-border);
}
</style>
