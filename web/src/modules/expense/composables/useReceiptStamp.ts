import { ref, computed } from 'vue'

// Decorative receipt chrome for the expense modal: a per-session transaction
// id, the printed date/time stamp and the barcode bars. Extracted from the
// component to keep the view focused on the form itself.
export function useReceiptStamp(barcodeSeed = 23) {
  function makeTxnId() {
    return 'TXN-' + Math.random().toString(36).slice(2, 8).toUpperCase()
  }

  const txnId = ref(makeTxnId())
  const stampDate = ref('')
  const stampTime = ref('')

  function refreshStamps() {
    const n = new Date()
    stampDate.value = n.toLocaleDateString('en-US', { month: '2-digit', day: '2-digit', year: '2-digit' })
    stampTime.value = n.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
  }
  refreshStamps()

  // Re-stamp the receipt for a fresh session (new date/time + transaction id).
  function newSession() {
    refreshStamps()
    txnId.value = makeTxnId()
  }

  const barcodeBars = computed(() => {
    let seed = barcodeSeed
    return Array.from({ length: 48 }, () => {
      seed = (seed * 9301 + 49297) % 233280
      const r = seed / 233280
      return r < 0.5 ? 1 : r < 0.85 ? 2 : 3
    })
  })

  return { txnId, stampDate, stampTime, refreshStamps, newSession, barcodeBars }
}
