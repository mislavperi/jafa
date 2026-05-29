<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Message from 'primevue/message'
import { useCreateExpense, useUpdateExpense } from '../composables/useExpenses'
import { useAllTags, useCreateTag, useAddTagToExpense, useRemoveTagFromExpense } from '../composables/useTags'
import { getTagsForExpense } from '../api/tag'
import type { Tag, RecurrenceInterval, RecurringSchedule, Expense } from '../models/expense'
import { tagColor } from '../constants'
import { useReceiptStamp } from '../composables/useReceiptStamp'
import { useThemeStore } from '@/stores/theme'
import { currencySymbol } from '@/core/currency'

const theme = useThemeStore()
const symbol = computed(() => currencySymbol(theme.currency))

// Receipt-styled control classes expressed as Tailwind utilities (kept here so
// the markup stays readable and the scoped <style> can stay minimal).
const dividerClass = 'border-0 border-t-[1.5px] border-dashed border-[var(--exp-receipt-border)] my-2.5'
const inputClass =
  'bg-transparent border-0 border-b-[1.5px] border-dashed border-[var(--exp-receipt-border)] pt-1 px-0.5 pb-1.5 text-[13px] text-[var(--exp-receipt-text)] outline-none focus:border-b-[var(--exp-receipt-text)] placeholder:text-[color-mix(in_srgb,var(--exp-receipt-text)_40%,transparent)]'
const stampBase =
  'flex items-center gap-[7px] px-[9px] py-[7px] bg-[var(--exp-receipt-bg)] border-[1.5px] border-[var(--exp-receipt-border)] text-[var(--exp-receipt-text)] text-[10.5px] tracking-[0.08em] cursor-pointer text-left uppercase font-semibold transition-all duration-[120ms] hover:border-[var(--exp-receipt-text)]'
const stampSel = 'bg-[var(--exp-receipt-text)] border-[var(--exp-receipt-text)] text-[var(--exp-receipt-bg)]'
const freqBase =
  'py-[7px] px-1 bg-transparent border-[1.5px] border-[var(--exp-receipt-text)] text-[var(--exp-receipt-text)] text-[9.5px] tracking-[0.14em] cursor-pointer uppercase font-bold transition-all hover:bg-[color-mix(in_srgb,var(--exp-receipt-text)_5%,transparent)]'
const freqSel = 'bg-[var(--exp-receipt-text)] text-[var(--exp-receipt-bg)]'
const recurringOnClass = 'bg-[color-mix(in_srgb,var(--jafa-accent)_8%,transparent)] border-[var(--jafa-accent)] border-solid'
const tenderClass =
  'bg-[var(--exp-receipt-text)] text-[var(--exp-receipt-bg)] border-0 cursor-pointer enabled:hover:opacity-85 disabled:bg-[var(--exp-receipt-border)]'

const props = defineProps<{ visible: boolean; expense?: Expense | null }>()
const emit = defineEmits<{ 'update:visible': [value: boolean] }>()

const editing = computed(() => !!props.expense)

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

// Decorative receipt chrome (txn id, date/time stamp, barcode).
const { txnId, stampDate, stampTime, newSession, barcodeBars } = useReceiptStamp(23)

watch(
  () => props.visible,
  async (v) => {
    if (!v) return
    newSession()
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

      <hr :class="dividerClass" />

      <div class="text-center my-1">
        <div class="text-[17px] font-extrabold tracking-[0.08em] uppercase">
          {{ editing ? 'Edit Transaction' : 'New Transaction' }}
        </div>
        <div class="mt-1 text-[10.5px] tracking-[0.16em] uppercase opacity-60">
          {{ editing ? 'Update the entry below' : 'Please fill in the details' }}
        </div>
      </div>

      <hr :class="dividerClass" />

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
          :class="inputClass" class="w-full"
          autofocus
        />
      </div>

      <!-- Unit Price + Date -->
      <div class="grid grid-cols-2 gap-4 mb-3">
        <div>
          <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase opacity-70 mb-1">
            <span>Unit Price</span>
            <span class="font-bold">{{ theme.currency }}</span>
          </label>
          <div class="relative">
            <span class="absolute left-0 bottom-1.5 font-bold text-[14px]">{{ symbol }}</span>
            <input
              v-model="cost"
              type="number"
              step="0.01"
              min="0"
              placeholder="0.00"
              :class="inputClass" class="w-full !pl-4 tabular-nums"
            />
          </div>
        </div>
        <div>
          <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase opacity-70 mb-1">
            <span>Date Stamped</span>
            <span class="font-bold">{{ displayDate }}</span>
          </label>
          <input v-model="date" type="date" :class="inputClass" class="w-full" />
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
          class="w-full flex items-center gap-2.5 px-3 py-2 border-[1.5px] border-current bg-transparent font-mono text-[12px] tracking-[0.06em] uppercase font-bold cursor-pointer hover:bg-[var(--exp-receipt-text)]/[0.04]"
          @click="tagsOpen = !tagsOpen"
        >
          <span
            v-if="selectedTags.length"
            class="w-2.5 h-2.5 rounded-[2px] shrink-0"
            :style="{ background: tagColor(selectedTags[0]!.id) }"
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
              :class="[stampBase, selectedTagIds.has(t.id) ? stampSel : '']"
              @click="toggleTag(t.id)"
            >
              <span class="w-2 h-2 rounded-[2px] shrink-0" :style="{ background: tagColor(t.id) }" />
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
              class="px-3 text-[10px] tracking-[0.14em] uppercase font-bold cursor-pointer border-[1.5px] border-current bg-transparent hover:bg-[var(--exp-receipt-text)] hover:text-[var(--exp-receipt-bg)] transition disabled:opacity-40 disabled:cursor-not-allowed"
              :disabled="!newTagName.trim()"
              @click="handleCreateTag"
            >
              Add
            </button>
          </div>
        </div>
      </div>

      <hr :class="dividerClass" />

      <!-- Recurring -->
      <div
        class="p-3.5 border-[1.5px] transition"
        :class="recurring ? recurringOnClass : 'border-dashed border-[var(--exp-receipt-border)] bg-[var(--exp-receipt-text)]/[0.025]'"
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
              :class="[freqBase, frequency === f.value ? freqSel : '']"
              @click="frequency = f.value"
            >
              {{ f.label }}
            </button>
          </div>
          <label class="block text-[10px] tracking-[0.16em] uppercase opacity-70 mb-1">Day of Month</label>
          <select v-model.number="recurrenceDay" :class="inputClass" class="w-full font-mono">
            <option v-for="d in dayOptions" :key="d" :value="d">{{ d }}</option>
          </select>
        </div>
      </div>

      <hr :class="dividerClass" />

      <!-- Subtotal lines -->
      <div class="flex flex-col">
        <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
          <span class="text-[10.5px] tracking-[0.14em] uppercase opacity-60">Subtotal</span>
          <span class="font-semibold tabular-nums">{{ symbol }}{{ costNum.toFixed(2) }}</span>
        </div>
        <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
          <span class="text-[10.5px] tracking-[0.14em] uppercase opacity-60">Category Tax</span>
          <span class="font-semibold tabular-nums">{{ symbol }}0.00</span>
        </div>
        <div class="flex justify-between items-baseline gap-4 py-0.5 text-[12px]">
          <span class="text-[10.5px] tracking-[0.14em] uppercase opacity-60">Schedule</span>
          <span class="font-semibold uppercase">{{ recurring ? frequency : 'one-time' }}</span>
        </div>
      </div>

      <hr :class="dividerClass" />

      <!-- TOTAL -->
      <div class="flex justify-between items-center gap-3 py-1 text-[16px] font-extrabold tracking-[0.04em]">
        <span>TOTAL</span>
        <span class="tabular-nums">{{ symbol }}{{ costNum.toFixed(2) }}</span>
      </div>

      <Message v-if="fieldError" severity="error" :closable="false" class="!mt-2">{{ fieldError }}</Message>
      <Message v-if="submitError" severity="error" :closable="false" class="!mt-2">{{ submitError }}</Message>

      <button
        type="submit"
        :disabled="isPending"
        :class="tenderClass"
        class="w-full mt-3 flex items-center justify-center gap-1.5 py-3 px-4 font-bold text-[13px] tracking-[0.18em] uppercase transition disabled:cursor-not-allowed"
      >
        <span v-if="isPending">PROCESSING…</span>
        <template v-else>
          <span>{{ editing ? 'Update Entry' : 'Ring it up' }}</span>
          <span>→</span>
        </template>
      </button>

      <button
        type="button"
        class="w-full mt-1.5 py-2 text-[10.5px] tracking-[0.18em] uppercase font-semibold opacity-70 hover:opacity-100 hover:text-[var(--jafa-accent)] transition cursor-pointer"
        @click="close"
      >
        ✕ Void transaction
      </button>

      <hr :class="dividerClass" class="mt-4" />

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
/*
 * The only rule kept in CSS: the "thermal paper" treatment for the receipt
 * (dotted texture + layered shadow). Everything else is Tailwind. This is hard
 * to express cleanly as utilities because of the multi-layer background-image
 * and color-mix dotting, so it stays here intentionally.
 */
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
</style>
