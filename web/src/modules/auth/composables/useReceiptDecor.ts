import { computed } from 'vue'

// PrimeVue passthrough that makes an <InputText> look like a dashed-underline
// receipt field. Shared by the login and register forms.
export const RECEIPT_INPUT_PT = {
  root: {
    class:
      '!bg-transparent !border-0 !border-b !border-dashed !border-[#8a8878] !rounded-none !text-[#1a1a1a] !font-mono !px-1 !py-1.5 focus:!border-[#1a1a1a] focus:!shadow-none placeholder:!text-[#a8a692]',
  },
}

export interface ReceiptDecorOptions {
  /** Currency symbol used in the floating backdrop amounts. */
  currency?: string
  /** Seeds for the deterministic PRNG so each page looks slightly different. */
  barcodeSeed?: number
  floatSeed?: number
  tickerSeed?: number
}

// Deterministic LCG so the decoration is stable across renders.
function lcg(seed: number): () => number {
  let s = seed
  return () => {
    s = (s * 9301 + 49297) % 233280
    return s / 233280
  }
}

/**
 * Receipt metadata + decorative backdrop (barcode bars, floating amounts and
 * the ticker path) shared by the auth pages. Extracted so the views only hold
 * their form logic.
 */
export function useReceiptDecor(opts: ReceiptDecorOptions = {}) {
  const { currency = '$', barcodeSeed = 7, floatSeed = 13, tickerSeed = 99 } = opts

  const lastUpdate = new Date(__LAST_UPDATE__)
  const stampDate = lastUpdate.toLocaleDateString('en-US', { month: '2-digit', day: '2-digit', year: '2-digit' })
  const stampTime = lastUpdate.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
  const version = __APP_VERSION__

  const barcodeBars = computed(() => {
    const next = lcg(barcodeSeed)
    return Array.from({ length: 52 }, () => {
      const r = next()
      return r < 0.5 ? 1 : r < 0.85 ? 2 : 3
    })
  })

  const floatingNums = computed(() => {
    const samples = [
      `${currency}12.40`, `−${currency}87.40`, `${currency}2,148`, `${currency}59.99`,
      `+${currency}480`, `${currency}10.99`, `${currency}1,850`, `−${currency}28.50`,
      `${currency}87.40`, `${currency}45.00`, `−${currency}15.49`,
    ]
    const next = lcg(floatSeed)
    return Array.from({ length: 18 }, () => {
      const val = samples[Math.floor(next() * samples.length)]!
      return {
        val,
        left: next() * 100,
        delay: next() * -30,
        dur: 20 + next() * 25,
        size: 11 + next() * 4,
        color: val.startsWith('+')
          ? 'rgba(34, 197, 94, 0.45)'
          : val.startsWith('−')
            ? 'rgba(239, 68, 68, 0.4)'
            : 'rgba(245, 197, 24, 0.4)',
      }
    })
  })

  const tickerPath = computed(() => {
    const next = lcg(tickerSeed)
    const points: [number, number][] = []
    let y = 60
    for (let x = 0; x <= 100; x += 4) {
      y += (next() - 0.5) * 16
      y = Math.max(20, Math.min(100, y))
      points.push([x, y])
    }
    return 'M ' + points.map((p) => `${p[0]} ${p[1]}`).join(' L ')
  })

  return { stampDate, stampTime, version, barcodeBars, floatingNums, tickerPath }
}
