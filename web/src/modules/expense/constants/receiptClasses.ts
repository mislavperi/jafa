// Receipt-styled control classes for AddExpenseModal, expressed as Tailwind
// utilities. Kept out of the component so the markup stays readable and the
// scoped <style> can stay minimal.

export const dividerClass =
  'border-0 border-t-[1.5px] border-dashed border-[var(--exp-receipt-border)] my-2.5'

export const inputClass =
  'bg-transparent border-0 border-b-[1.5px] border-dashed border-[var(--exp-receipt-border)] pt-1 px-0.5 pb-1.5 text-[13px] text-[var(--exp-receipt-text)] outline-none focus:border-b-[var(--exp-receipt-text)] placeholder:text-[color-mix(in_srgb,var(--exp-receipt-text)_40%,transparent)]'

export const stampBase =
  'flex items-center gap-[7px] px-[9px] py-[7px] bg-[var(--exp-receipt-bg)] border-[1.5px] border-[var(--exp-receipt-border)] text-[var(--exp-receipt-text)] text-[10.5px] tracking-[0.08em] cursor-pointer text-left uppercase font-semibold transition-all duration-[120ms] hover:border-[var(--exp-receipt-text)]'

export const stampSel =
  'bg-[var(--exp-receipt-text)] border-[var(--exp-receipt-text)] text-[var(--exp-receipt-bg)]'

export const freqBase =
  'py-[7px] px-1 bg-transparent border-[1.5px] border-[var(--exp-receipt-text)] text-[var(--exp-receipt-text)] text-[9.5px] tracking-[0.14em] cursor-pointer uppercase font-bold transition-all hover:bg-[color-mix(in_srgb,var(--exp-receipt-text)_5%,transparent)]'

export const freqSel = 'bg-[var(--exp-receipt-text)] text-[var(--exp-receipt-bg)]'

export const recurringOnClass =
  'bg-[color-mix(in_srgb,var(--jafa-accent)_8%,transparent)] border-[var(--jafa-accent)] border-solid'

export const tenderClass =
  'bg-[var(--exp-receipt-text)] text-[var(--exp-receipt-bg)] border-0 cursor-pointer enabled:hover:opacity-85 disabled:bg-[var(--exp-receipt-border)]'
