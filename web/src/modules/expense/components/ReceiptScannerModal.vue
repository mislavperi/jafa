<script setup lang="ts">
import { ref, computed } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import ProgressSpinner from 'primevue/progressspinner'
import Select from 'primevue/select'
import Checkbox from 'primevue/checkbox'
import Message from 'primevue/message'
import { useCreateExpense } from '../composables/useExpenses'
import { useAllTags, useCreateTag, useAddTagToExpense } from '../composables/useTags'
import type { Tag } from '../models/expense'
import type { ScanStep, Sample, ReviewItem } from '../models/receiptScan'
import { SAMPLES, CATEGORY_HINTS, guessCategory } from '../constants/receiptSamples'
import { useThemeStore } from '@/stores/theme'
import { formatCurrency } from '@/core/currency'

const theme = useThemeStore()

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{ 'update:visible': [value: boolean] }>()

const step = ref<ScanStep>('upload')
const imageData = ref<string | null>(null)
const parsed = ref<Sample | null>(null)
const items = ref<ReviewItem[]>([])
const dragging = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const submitting = ref(false)
const submitError = ref<string | null>(null)
// Live OCR is not implemented yet. A real upload is previewed but NOT parsed —
// only the explicit "try a sample" receipts run the (demo) scan flow, so we
// never present fabricated line items as if they were read from the user's image.
const liveScanUnavailable = ref(false)

const { mutateAsync: createExpense } = useCreateExpense()
const { data: allTags } = useAllTags()
const { mutateAsync: createTag } = useCreateTag()
const { mutateAsync: addTag } = useAddTagToExpense()

function reset() {
  step.value = 'upload'
  imageData.value = null
  parsed.value = null
  items.value = []
  dragging.value = false
  submitting.value = false
  submitError.value = null
  liveScanUnavailable.value = false
}

function close() {
  reset()
  emit('update:visible', false)
}

function handleFile(file: File | null | undefined) {
  if (!file || !file.type.startsWith('image/')) return
  const reader = new FileReader()
  reader.onload = (e) => {
    imageData.value = (e.target?.result as string) ?? null
    // No real OCR yet — preview the image and tell the user to use a sample
    // instead of fabricating items from a random canned receipt.
    liveScanUnavailable.value = true
  }
  reader.readAsDataURL(file)
}

function useSample(s: Sample) {
  imageData.value = null
  liveScanUnavailable.value = false
  startScan(s)
}

function startScan(sample: Sample) {
  step.value = 'scanning'
  setTimeout(() => {
    parsed.value = sample
    items.value = sample.items.map((it, i) => {
      const tag = it.category ?? guessCategory(it.name)
      return {
        ...it,
        id: i,
        included: true,
        suggestedTag: tag,
        finalTag: tag,
        needsReview: it.confidence < 0.6,
      }
    })
    step.value = 'review'
  }, 1500)
}

function onDrop(e: DragEvent) {
  e.preventDefault()
  dragging.value = false
  handleFile(e.dataTransfer?.files[0])
}

function onPick(e: Event) {
  const t = e.target as HTMLInputElement
  handleFile(t.files?.[0])
}

function toggleItem(id: number) {
  const it = items.value.find((i) => i.id === id)
  if (it) it.included = !it.included
}

function setTag(id: number, tag: string) {
  const it = items.value.find((i) => i.id === id)
  if (it) {
    it.finalTag = tag
    it.needsReview = false
  }
}

const needsReview = computed(() => items.value.filter((i) => i.needsReview))
const autoItems = computed(() => items.value.filter((i) => !i.needsReview))
const includedItems = computed(() => items.value.filter((i) => i.included && i.amount > 0))
const includedTotal = computed(() => includedItems.value.reduce((s, i) => s + i.amount, 0))

const tagOptions = computed(() => {
  const known = new Set<string>([...Object.keys(CATEGORY_HINTS), 'other'])
  for (const t of allTags.value ?? []) known.add(t.name.toLowerCase())
  for (const i of items.value) known.add(i.finalTag)
  return [...known].map((name) => ({ label: name, value: name }))
})

async function ensureTag(name: string): Promise<Tag> {
  const existing = (allTags.value ?? []).find((t) => t.name.toLowerCase() === name.toLowerCase())
  if (existing) return existing
  return await createTag(name)
}

async function commit() {
  if (!parsed.value || includedItems.value.length === 0) return
  submitting.value = true
  submitError.value = null
  const failed: string[] = []
  for (const it of includedItems.value) {
    try {
      const expense = await createExpense({ name: it.name, amount: 1, cost: it.amount })
      try {
        const tag = await ensureTag(it.finalTag)
        await addTag({ expenseId: expense.id, tagId: tag.id })
      } catch {
        // tag failure shouldn't block expense creation
      }
    } catch {
      failed.push(it.name)
    }
  }
  submitting.value = false
  if (failed.length === includedItems.value.length) {
    submitError.value = 'Failed to import items'
    return
  }
  if (failed.length) {
    submitError.value = `Imported ${includedItems.value.length - failed.length}/${includedItems.value.length}. Failed: ${failed.join(', ')}`
    return
  }
  close()
}

function fmt(n: number) {
  const sign = n < 0 ? '-' : ''
  return `${sign}${formatCurrency(Math.abs(n), theme.currency)}`
}

function confidenceTone(c: number) {
  if (c > 0.85) return 'bg-emerald-500'
  if (c > 0.6) return 'bg-amber-500'
  return 'bg-red-500'
}
</script>

<template>
  <Dialog
    :visible="props.visible"
    @update:visible="emit('update:visible', $event)"
    modal
    :closable="false"
    :show-header="false"
    :style="{ width: '60rem', maxWidth: '95vw' }"
    content-class="p-0"
  >
    <div class="bg-[var(--jafa-surface)] text-[var(--jafa-text)]">
      <!-- Header -->
      <div class="flex items-start justify-between px-6 py-4 border-b border-[var(--jafa-border)]">
        <div>
          <h3 class="text-[calc(15px*var(--jafa-text-scale,1))] font-semibold tracking-[0.04em]">SCAN RECEIPT</h3>
          <p class="text-[calc(11px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] uppercase tracking-[0.1em] mt-0.5">
            <span v-if="step === 'upload'">Step 1 / 3 · Upload</span>
            <span v-else-if="step === 'scanning'">Step 2 / 3 · Scanning</span>
            <span v-else>Step 3 / 3 · Review items</span>
          </p>
        </div>
        <Button icon="pi pi-times" text severity="secondary" rounded @click="close" />
      </div>

      <!-- Progress strip -->
      <div class="flex border-b border-[var(--jafa-border)] text-[calc(11px*var(--jafa-text-scale,1))] uppercase tracking-[0.12em]">
        <div
          v-for="(label, i) in ['Upload', 'Scan', 'Review']"
          :key="label"
          class="flex-1 px-4 py-2 text-center border-r border-[var(--jafa-border)] last:border-r-0"
          :class="{
            'text-[var(--jafa-accent)]': (i === 0 && step === 'upload') || (i === 1 && step === 'scanning') || (i === 2 && step === 'review'),
            'text-emerald-400': (i === 0 && step !== 'upload') || (i === 1 && step === 'review'),
            'text-[var(--jafa-text-dim)]': (i === 1 && step === 'upload') || (i === 2 && step !== 'review'),
          }"
        >
          <span v-if="(i === 0 && step !== 'upload') || (i === 1 && step === 'review')" class="mr-1">✓</span>{{ label }}
        </div>
      </div>

      <!-- Step 1: Upload -->
      <div v-if="step === 'upload'" class="p-6">
        <div
          class="border-2 border-dashed rounded-[14px] p-10 flex flex-col items-center justify-center cursor-pointer transition"
          :class="dragging ? 'border-[var(--jafa-accent)] bg-[var(--jafa-accent)]/5' : 'border-[var(--jafa-border-strong)] hover:border-[#52525a]'"
          @dragover.prevent="dragging = true"
          @dragleave="dragging = false"
          @drop="onDrop"
          @click="fileInput?.click()"
        >
          <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="onPick" />
          <i class="pi pi-receipt text-[calc(40px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] mb-3" />
          <div class="text-[calc(13px*var(--jafa-text-scale,1))] font-semibold tracking-[0.12em] uppercase">Drop receipt here</div>
          <div class="text-[calc(11px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] mt-1 tracking-wider uppercase">PNG · JPG · HEIC · max 10MB</div>
          <Button label="Choose file" icon="pi pi-plus" severity="secondary" size="small" class="mt-4" @click.stop="fileInput?.click()" />
        </div>

        <Message v-if="liveScanUnavailable" severity="info" :closable="false" class="mt-4">
          Automatic receipt reading isn't available yet — your uploaded image won't be parsed.
          Pick one of the sample receipts below to preview the import flow.
        </Message>

        <div class="flex items-center gap-3 my-6 text-[calc(11px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] uppercase tracking-[0.14em]">
          <div class="h-px flex-1 bg-[var(--jafa-border)]" />
          <span>Or try a sample</span>
          <div class="h-px flex-1 bg-[var(--jafa-border)]" />
        </div>

        <div class="grid grid-cols-3 gap-3">
          <button
            v-for="s in SAMPLES"
            :key="s.id"
            type="button"
            class="text-left border border-[var(--jafa-border)] rounded-[12px] p-3 bg-[#0e0e10] hover:border-[var(--jafa-accent)]/60 hover:bg-[#16161a] transition"
            @click="useSample(s)"
          >
            <div class="bg-[#f5f1e6] text-[#1a1a1a] rounded-md p-3 mb-2 font-mono text-[calc(10px*var(--jafa-text-scale,1))] leading-tight">
              <div class="text-center font-bold tracking-widest">{{ s.merchant }}</div>
              <div class="text-center text-[#7a7864] text-[calc(9px*var(--jafa-text-scale,1))] mt-0.5">{{ s.date }}</div>
              <hr class="my-1.5 border-0 border-t border-dashed border-[#c7c5b8]" />
              <div v-for="(it, idx) in s.items.slice(0, 3)" :key="idx" class="flex justify-between">
                <span class="truncate pr-1">{{ it.name }}</span>
                <span class="font-semibold">{{ fmt(it.amount) }}</span>
              </div>
              <div v-if="s.items.length > 3" class="text-[#8a8878] italic">… +{{ s.items.length - 3 }} more</div>
              <hr class="my-1.5 border-0 border-t border-dashed border-[#c7c5b8]" />
              <div class="flex justify-between font-bold"><span>TOTAL</span><span>{{ fmt(s.total) }}</span></div>
            </div>
            <div class="text-[calc(12px*var(--jafa-text-scale,1))] font-semibold text-[var(--jafa-text)] truncate">{{ s.merchant }}</div>
            <div class="text-[calc(11px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)]">{{ s.items.length }} items · {{ fmt(s.total) }}</div>
          </button>
        </div>
      </div>

      <!-- Step 2: Scanning -->
      <div v-else-if="step === 'scanning'" class="flex flex-col items-center justify-center py-20">
        <ProgressSpinner style="width:48px;height:48px" stroke-width="3" />
        <div class="text-[calc(13px*var(--jafa-text-scale,1))] font-semibold tracking-[0.18em] uppercase mt-5 text-[var(--jafa-accent)]">Scanning…</div>
        <div class="text-[calc(12px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] mt-1">Reading receipt and extracting items</div>
      </div>

      <!-- Step 3: Review -->
      <div v-else-if="step === 'review' && parsed" class="grid grid-cols-[1fr_1.6fr] max-h-[70vh]">
        <!-- Left: receipt preview + meta -->
        <div class="p-5 border-r border-[var(--jafa-border)] overflow-auto">
          <div class="bg-[#f5f1e6] text-[#1a1a1a] rounded-md p-4 font-mono text-[calc(11px*var(--jafa-text-scale,1))]">
            <img v-if="imageData" :src="imageData" alt="receipt" class="w-full rounded" />
            <template v-else>
              <div class="text-center font-bold tracking-widest text-[calc(12px*var(--jafa-text-scale,1))]">{{ parsed.merchant }}</div>
              <div class="text-center text-[#7a7864] text-[calc(10px*var(--jafa-text-scale,1))] mt-0.5">{{ parsed.address }}</div>
              <div class="text-center text-[#7a7864] text-[calc(10px*var(--jafa-text-scale,1))]">{{ parsed.date }}</div>
              <hr class="my-2 border-0 border-t border-dashed border-[#c7c5b8]" />
              <div v-for="(it, i) in parsed.items" :key="i" class="flex justify-between gap-2">
                <span class="truncate">{{ it.name }}</span>
                <span class="font-semibold shrink-0">{{ fmt(it.amount) }}</span>
              </div>
              <hr class="my-2 border-0 border-t border-dashed border-[#c7c5b8]" />
              <div class="flex justify-between font-bold"><span>TOTAL</span><span>{{ fmt(parsed.total) }}</span></div>
              <div class="text-center mt-3 text-[calc(10px*var(--jafa-text-scale,1))] tracking-widest">★ THANK YOU ★</div>
            </template>
          </div>
          <div class="mt-4 flex flex-col gap-1.5 text-[calc(12px*var(--jafa-text-scale,1))]">
            <div class="flex justify-between"><span class="text-[var(--jafa-text-muted)] uppercase tracking-wider text-[calc(10px*var(--jafa-text-scale,1))]">Merchant</span><span class="text-[var(--jafa-text)] font-medium">{{ parsed.merchant }}</span></div>
            <div class="flex justify-between"><span class="text-[var(--jafa-text-muted)] uppercase tracking-wider text-[calc(10px*var(--jafa-text-scale,1))]">Date</span><span class="text-[var(--jafa-text)] tabular-nums">{{ parsed.date }}</span></div>
            <div class="flex justify-between"><span class="text-[var(--jafa-text-muted)] uppercase tracking-wider text-[calc(10px*var(--jafa-text-scale,1))]">Receipt Total</span><span class="text-[var(--jafa-text)] tabular-nums font-medium">{{ fmt(parsed.total) }}</span></div>
          </div>
        </div>

        <!-- Right: items -->
        <div class="flex flex-col overflow-hidden">
          <div class="flex-1 overflow-auto p-5 flex flex-col gap-5">
            <!-- Needs review -->
            <section v-if="needsReview.length" class="border border-dashed border-amber-500/50 rounded-[12px] p-4 bg-amber-500/[0.04]">
              <header class="mb-3">
                <h4 class="text-[calc(11px*var(--jafa-text-scale,1))] font-bold tracking-[0.1em] uppercase text-[var(--jafa-accent)] flex items-center gap-1.5">
                  <i class="pi pi-exclamation-triangle text-[calc(11px*var(--jafa-text-scale,1))]" />
                  Needs review · {{ needsReview.length }}
                </h4>
                <p class="text-[calc(11px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] mt-0.5">Couldn't categorize confidently. Pick a tag or skip.</p>
              </header>
              <div class="flex flex-col gap-2">
                <div v-for="it in needsReview" :key="it.id" class="flex items-center gap-3 py-2 border-b border-[var(--jafa-border)] last:border-0" :class="{ 'opacity-50': !it.included }">
                  <Checkbox :model-value="it.included" :binary="true" @update:model-value="toggleItem(it.id)" />
                  <div class="flex-1 min-w-0">
                    <div class="text-[calc(13px*var(--jafa-text-scale,1))] font-medium text-[var(--jafa-text)] truncate">{{ it.name }}</div>
                    <div class="flex items-center gap-1.5 mt-0.5">
                      <span class="w-1.5 h-1.5 rounded-full" :class="confidenceTone(it.confidence)" />
                      <span class="text-[calc(10px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] uppercase tracking-wider">{{ Math.round(it.confidence * 100) }}% confidence</span>
                    </div>
                  </div>
                  <Select
                    :model-value="it.finalTag"
                    :options="tagOptions"
                    option-label="label"
                    option-value="value"
                    placeholder="Tag"
                    size="small"
                    class="w-36"
                    @update:model-value="(v) => setTag(it.id, v)"
                  />
                  <div class="w-20 text-right tabular-nums text-[calc(13px*var(--jafa-text-scale,1))] font-semibold">{{ fmt(it.amount) }}</div>
                </div>
              </div>
            </section>

            <!-- Auto-categorized -->
            <section class="border border-[var(--jafa-border)] rounded-[12px] p-4">
              <header class="mb-3">
                <h4 class="text-[calc(11px*var(--jafa-text-scale,1))] font-bold tracking-[0.1em] uppercase text-emerald-400 flex items-center gap-1.5">
                  <i class="pi pi-check text-[calc(11px*var(--jafa-text-scale,1))]" />
                  Auto-categorized · {{ autoItems.length }}
                </h4>
                <p class="text-[calc(11px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] mt-0.5">Click a tag to change.</p>
              </header>
              <div class="flex flex-col gap-2">
                <div v-for="it in autoItems" :key="it.id" class="flex items-center gap-3 py-2 border-b border-[var(--jafa-border)] last:border-0" :class="{ 'opacity-50': !it.included }">
                  <Checkbox :model-value="it.included" :binary="true" @update:model-value="toggleItem(it.id)" />
                  <div class="flex-1 min-w-0">
                    <div class="text-[calc(13px*var(--jafa-text-scale,1))] font-medium text-[var(--jafa-text)] truncate">{{ it.name }}</div>
                    <div class="flex items-center gap-1.5 mt-0.5">
                      <span class="w-1.5 h-1.5 rounded-full" :class="confidenceTone(it.confidence)" />
                      <span class="text-[calc(10px*var(--jafa-text-scale,1))] text-[var(--jafa-text-muted)] uppercase tracking-wider">{{ Math.round(it.confidence * 100) }}% confidence</span>
                    </div>
                  </div>
                  <Select
                    :model-value="it.finalTag"
                    :options="tagOptions"
                    option-label="label"
                    option-value="value"
                    size="small"
                    class="w-36"
                    @update:model-value="(v) => setTag(it.id, v)"
                  />
                  <div class="w-20 text-right tabular-nums text-[calc(13px*var(--jafa-text-scale,1))] font-semibold">{{ fmt(it.amount) }}</div>
                </div>
              </div>
            </section>
          </div>

          <!-- Footer summary + actions -->
          <div class="border-t border-[var(--jafa-border)] p-4 flex flex-col gap-3 bg-[#0f0f12]">
            <Message v-if="submitError" severity="error" :closable="false">{{ submitError }}</Message>
            <div class="flex items-center justify-between">
              <div class="flex flex-col">
                <span class="text-[calc(10px*var(--jafa-text-scale,1))] uppercase tracking-[0.12em] text-[var(--jafa-text-muted)]">Selected</span>
                <span class="text-[calc(12px*var(--jafa-text-scale,1))] tabular-nums">{{ includedItems.length }} items</span>
              </div>
              <div class="flex flex-col items-end">
                <span class="text-[calc(10px*var(--jafa-text-scale,1))] uppercase tracking-[0.12em] text-[var(--jafa-text-muted)]">Total to add</span>
                <span class="text-[calc(16px*var(--jafa-text-scale,1))] font-semibold tabular-nums">{{ fmt(includedTotal) }}</span>
              </div>
            </div>
            <div class="flex gap-2 justify-end">
              <Button label="Scan another" icon="pi pi-chevron-left" severity="secondary" size="small" @click="reset" />
              <Button
                :label="`Add ${includedItems.length} expense${includedItems.length !== 1 ? 's' : ''}`"
                icon="pi pi-chevron-right"
                icon-pos="right"
                size="small"
                :disabled="includedItems.length === 0"
                :loading="submitting"
                @click="commit"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </Dialog>
</template>
