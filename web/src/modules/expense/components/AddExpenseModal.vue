<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Message from 'primevue/message'
import Popover from 'primevue/popover'
import { Form, FormField } from '@primevue/forms'
import { useCreateExpense, useUpdateExpense } from '../composables/useExpenses'
import { useAllTags, useCreateTag, useAddTagToExpense, useRemoveTagFromExpense } from '../composables/useTags'
import { getTagsForExpense } from '../api/tag'
import type { Tag, RecurrenceInterval, RecurringSchedule, Expense, ExpenseKind } from '../models/expense'
import { tagColor } from '../constants'
import {
  dividerClass,
  inputClass,
  stampBase,
  stampSel,
  freqBase,
  freqSel,
  recurringOnClass,
  tenderClass,
} from '../constants/receiptClasses'
import { useReceiptStamp } from '../composables/useReceiptStamp'
import { useThemeStore } from '@/stores/theme'
import { currencySymbol } from '@/core/currency'

const theme = useThemeStore()
const symbol = computed(() => currencySymbol(theme.currency))

const props = defineProps<{ visible: boolean; expense?: Expense | null }>()
const emit = defineEmits<{ 'update:visible': [value: boolean] }>()

const editing = computed(() => !!props.expense)

const kind = ref<ExpenseKind>('expense')
const isIncome = computed(() => kind.value === 'income')
const name = ref('')
const cost = ref<string>('')
const date = ref<string>(new Date().toISOString().slice(0, 10))
const selectedTagIds = ref<Set<number>>(new Set())
const originalTagIds = ref<Set<number>>(new Set())
const newTagName = ref('')
const recurring = ref(false)
const frequency = ref<RecurrenceInterval>('monthly')
const recurrenceDay = ref<number>(new Date().getDate())
const split = ref(false)
const installmentCount = ref<number>(2)
const submitError = ref<string | null>(null)
const tagsOpen = ref(false)
const tagPopover = ref<InstanceType<typeof Popover> | null>(null)

function toggleTags(event: Event) {
  tagPopover.value?.toggle(event)
}

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

// Per-payment amount when the cost is split into installments (e.g. a $200
// phone split into 4 → $50). Clamped so we never divide by < 1.
const paymentAmount = computed(() => {
  const n = Math.max(1, Math.floor(installmentCount.value || 0))
  return costNum.value / n
})

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
      kind.value = props.expense.kind ?? 'expense'
      name.value = props.expense.name
      cost.value = String(props.expense.cost ?? '')
      if (props.expense.recurringSchedule) {
        recurring.value = true
        frequency.value = props.expense.recurringSchedule.interval
        recurrenceDay.value = props.expense.recurringSchedule.dayOfMonth
        date.value = props.expense.recurringSchedule.startDate
      }
      if (props.expense.installmentPlan) {
        split.value = true
        installmentCount.value = props.expense.installmentPlan.count
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
  kind.value = 'expense'
  name.value = ''
  cost.value = ''
  date.value = new Date().toISOString().slice(0, 10)
  selectedTagIds.value = new Set()
  originalTagIds.value = new Set()
  newTagName.value = ''
  recurring.value = false
  frequency.value = 'monthly'
  recurrenceDay.value = new Date().getDate()
  split.value = false
  installmentCount.value = 2
  submitError.value = null
  tagsOpen.value = false
}

function close() {
  reset()
  emit('update:visible', false)
}

// PrimeVue Forms resolver. It validates the receipt fields (name, cost) and
// maps messages back onto the matching FormField via its error key.
function resolver() {
  const errors: Record<string, { message: string }[]> = {}
  if (!name.value.trim()) errors.name = [{ message: 'Name required' }]
  if (isNaN(costNum.value) || costNum.value <= 0) errors.cost = [{ message: 'Cost required' }]
  if (split.value && (!Number.isFinite(installmentCount.value) || installmentCount.value < 2)) {
    errors.installmentCount = [{ message: 'At least 2 payments required' }]
  }
  return { values: { name: name.value, cost: cost.value }, errors }
}

async function handleSubmit({ valid }: { valid: boolean }) {
  submitError.value = null
  if (!valid) return
  try {
    let recurringSchedule: RecurringSchedule | undefined
    if (recurring.value) {
      recurringSchedule = {
        interval: frequency.value,
        dayOfMonth: recurrenceDay.value,
        startDate: date.value,
      }
    }
    let splitCount: number | undefined
    if (split.value && !isIncome.value) {
      splitCount = Math.floor(installmentCount.value)
    }
    const payload = {
      name: name.value.trim(),
      kind: kind.value,
      amount: 1,
      cost: costNum.value,
      recurringSchedule,
      installmentCount: splitCount,
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
    <Form
      :resolver="resolver"
      class="exp-receipt relative w-full font-mono px-7 pt-9 pb-6"
      @submit="handleSubmit"
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

      <!-- Entry type: expense vs income -->
      <div class="grid grid-cols-2 gap-1.5 mb-3">
        <button
          type="button"
          data-testid="kind-expense"
          :class="[freqBase, kind === 'expense' ? freqSel : '']"
          @click="kind = 'expense'"
        >
          Expense
        </button>
        <button
          type="button"
          data-testid="kind-income"
          :class="[freqBase, kind === 'income' ? freqSel : '']"
          @click="kind = 'income'"
        >
          Income
        </button>
      </div>

      <!-- Item Description -->
      <FormField v-slot="$field" name="name" class="mb-3 block">
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
        <div v-if="$field.invalid" class="mt-1 text-[11px] text-red-500 font-medium">
          {{ $field.error?.message }}
        </div>
      </FormField>

      <!-- Unit Price + Date -->
      <div class="grid grid-cols-2 gap-4 mb-3">
        <FormField v-slot="$field" name="cost" class="block">
          <label class="flex justify-between items-baseline text-[10px] tracking-[0.16em] uppercase opacity-70 mb-1">
            <span>{{ isIncome ? 'Amount' : 'Unit Price' }}</span>
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
          <div v-if="$field.invalid" class="mt-1 text-[11px] text-red-500 font-medium">
            {{ $field.error?.message }}
          </div>
        </FormField>
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
          @click="toggleTags"
        >
          <span
            v-if="selectedTags.length"
            class="w-2.5 h-2.5 rounded-[2px] shrink-0"
            :style="{ background: tagColor(selectedTags[0]!.id) }"
          />
          <span class="flex-1 text-left truncate">{{ departmentLabel }}</span>
          <span class="text-[11px] opacity-70">{{ tagsOpen ? '▴' : '▾' }}</span>
        </button>

        <Popover
          ref="tagPopover"
          @show="tagsOpen = true"
          @hide="tagsOpen = false"
          :pt="{ content: { class: '!p-0' } }"
        >
          <div
            class="w-[18rem] max-w-[90vw] p-2.5 border-[1.5px] border-dashed border-[var(--exp-receipt-border)] bg-[var(--exp-receipt-bg)] text-[var(--exp-receipt-text)] font-mono"
          >
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
        </Popover>
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

      <hr v-if="!isIncome" :class="dividerClass" />

      <!-- Split into payments (expenses only) -->
      <div
        v-if="!isIncome"
        class="p-3.5 border-[1.5px] transition"
        :class="split ? recurringOnClass : 'border-dashed border-[var(--exp-receipt-border)] bg-[var(--exp-receipt-text)]/[0.025]'"
      >
        <label class="flex items-center justify-between gap-3 cursor-pointer select-none">
          <div class="flex flex-col gap-0.5">
            <div class="text-[12.5px] font-extrabold tracking-[0.1em] uppercase">Split Payment</div>
            <div class="text-[10px] tracking-[0.1em] uppercase opacity-60">Divide cost into equal payments</div>
          </div>
          <span class="relative inline-block w-9 h-5 shrink-0">
            <input
              v-model="split"
              type="checkbox"
              data-testid="split-toggle"
              class="opacity-0 w-0 h-0 peer"
            />
            <span
              class="absolute inset-0 rounded-full transition cursor-pointer"
              :class="split ? 'bg-[var(--jafa-accent)]' : 'bg-[var(--exp-receipt-border)]'"
            />
            <span
              class="absolute top-0.5 left-0.5 w-4 h-4 rounded-full bg-[var(--exp-receipt-bg)] transition"
              :class="{ 'translate-x-4': split }"
            />
          </span>
        </label>

        <div v-if="split" class="mt-2.5 pt-2.5 border-t-[1.5px] border-dashed border-[var(--exp-receipt-border)]">
          <label class="block text-[10px] tracking-[0.16em] uppercase opacity-70 mb-1">Number of Payments</label>
          <input
            v-model.number="installmentCount"
            type="number"
            min="2"
            step="1"
            data-testid="installment-count"
            :class="inputClass"
            class="w-full font-mono tabular-nums"
          />
          <div v-if="split && installmentCount < 2" class="mt-1 text-[11px] text-red-500 font-medium">
            At least 2 payments required
          </div>
          <div class="mt-2 flex justify-between items-baseline text-[12px]">
            <span class="text-[10.5px] tracking-[0.14em] uppercase opacity-60">Per Payment</span>
            <span class="font-bold tabular-nums" data-testid="per-payment">
              {{ installmentCount >= 2 ? `${symbol}${paymentAmount.toFixed(2)} × ${Math.floor(installmentCount)}` : '—' }}
            </span>
          </div>
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
    </Form>
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
