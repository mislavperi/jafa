import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { queryClient } from './core/query'
import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura'
import { definePreset } from '@primeuix/themes'

import App from './App.vue'
import router from './router'

const MONO_STACK = `'Geist Mono', 'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, monospace`

const JafaPreset = definePreset(Aura, {
  primitive: {
    borderRadius: {
      none: '0',
      xs: '3px',
      sm: '5px',
      md: '8px',
      lg: '12px',
      xl: '16px',
    },
  },

  extend: {
    jafa: {
      brand:        '#f5c518',
      brandHover:   '#ffd233',
      brandDim:     '#3d2f0a',
      brandInk:     '#0a0a0b',
      bg:           '#0a0a0b',
      surface:      '#131316',
      surface2:     '#18181c',
      surface3:     '#1f1f24',
      border:       '#26262c',
      borderStrong: '#32323a',
      fontMono:     MONO_STACK,
    },
  },

  semantic: {
    primary: {
      50: '{amber.50}',
      100: '{amber.100}',
      200: '{amber.200}',
      300: '{amber.300}',
      400: '{amber.400}',
      500: '{amber.500}',
      600: '{amber.600}',
      700: '{amber.700}',
      800: '{amber.800}',
      900: '{amber.900}',
      950: '{amber.950}',
    },

    colorScheme: {
      light: {
        surface: {
          0:   '#ffffff',
          50:  '{zinc.50}',
          100: '{zinc.100}',
          200: '{zinc.200}',
          300: '{zinc.300}',
          400: '{zinc.400}',
          500: '{zinc.500}',
          600: '{zinc.600}',
          700: '{zinc.700}',
          800: '{zinc.800}',
          900: '{zinc.900}',
          950: '{zinc.950}',
        },
        content: {
          background: '{surface.50}',
          hoverBackground: '{surface.100}',
          borderColor: '{surface.200}',
          color: '{text.color}',
          hoverColor: '{text.hover.color}',
        },
        text: {
          color: '{surface.700}',
          hoverColor: '{surface.900}',
          mutedColor: '{surface.400}',
          hoverMutedColor: '{surface.500}',
        },
      },
      dark: {
        surface: {
          0:   '#ffffff',
          50:  '#f5f5f7',
          100: '#e4e4e7',
          200: '#a1a1aa',
          300: '#71717a',
          400: '#52525b',
          500: '#3f3f46',
          600: '#32323a',
          700: '#26262c',
          800: '#131316',
          900: '#0a0a0b',
          950: '#000000',
        },
        content: {
          background: '#131316',
          hoverBackground: '#1f1f24',
          borderColor: '#26262c',
          color: '{text.color}',
          hoverColor: '{text.hover.color}',
        },
        text: {
          color: '#f5f5f7',
          hoverColor: '#ffffff',
          mutedColor: '#a1a1aa',
          hoverMutedColor: '#d4d4d8',
        },
      },
    },

    formField: {
      fontFamily: MONO_STACK,
    },
  },

  css: ({ dt }: { dt: (token: string) => string }) => `
    html, body, #app { height: 100%; }

    * { font-family: ${dt('jafa.fontMono')}; }

    body {
      background: ${dt('jafa.bg')};
      color: #f5f5f7;
      font-size: 13.5px;
      line-height: 1.5;
      -webkit-font-smoothing: antialiased;
      -moz-osx-font-smoothing: grayscale;
    }

    .p-datatable td,
    .tabular { font-variant-numeric: tabular-nums; }

    .p-datatable .p-datatable-thead > tr > th {
      text-transform: uppercase;
      letter-spacing: 0.06em;
      font-size: 0.7rem;
      font-weight: 600;
    }

    .p-panel .p-panel-header .p-panel-title {
      text-transform: uppercase;
      letter-spacing: 0.06em;
      font-size: 0.7rem;
    }

    *:focus-visible {
      outline: 2px solid ${dt('jafa.brand')};
      outline-offset: 2px;
    }

    ::-webkit-scrollbar { width: 10px; height: 10px; }
    ::-webkit-scrollbar-track { background: transparent; }
    ::-webkit-scrollbar-thumb {
      background: ${dt('jafa.border')};
      border-radius: 8px;
      border: 2px solid ${dt('jafa.bg')};
    }
    ::-webkit-scrollbar-thumb:hover { background: ${dt('jafa.borderStrong')}; }

    .receipt-shell {
      position: relative;
      width: 100%;
      max-width: 380px;
      background: #fafaf6;
      color: #1a1a1a;
      font-family: ${dt('jafa.fontMono')};
      padding: 40px 32px 28px;
      box-shadow:
        0 1px 0 rgba(255,255,255,0.04),
        0 30px 80px -20px rgba(0,0,0,0.7),
        0 50px 120px -40px rgba(245, 197, 24, 0.15);
      background-image:
        radial-gradient(rgba(0,0,0,0.015) 1px, transparent 1px),
        radial-gradient(rgba(0,0,0,0.015) 1px, transparent 1px);
      background-size: 3px 3px, 7px 7px;
      background-position: 0 0, 1.5px 1.5px;
    }
    .receipt-shell::before, .receipt-shell::after {
      content: '';
      position: absolute;
      left: 0;
      right: 0;
      height: 14px;
      background: radial-gradient(circle at 7px 0, transparent 6px, #fafaf6 6.5px) 0 0 / 14px 14px repeat-x;
    }
    .receipt-shell::before { top: -1px; transform: rotate(180deg); }
    .receipt-shell::after { bottom: -13px; }

    .receipt-divider {
      border: none;
      border-top: 1.5px dashed #c7c5b8;
      margin: 14px 0;
    }

    .receipt-input {
      width: 100%;
      background: transparent;
      border: none;
      border-bottom: 1.5px solid #1a1a1a;
      padding: 4px 0 6px;
      font-family: ${dt('jafa.fontMono')};
      font-size: 14px;
      color: #1a1a1a;
      outline: none;
      letter-spacing: 0.02em;
    }
    .receipt-input::placeholder { color: #b8b6a8; }
    .receipt-input:focus { border-bottom-color: #d97706; }

    .receipt-btn {
      width: 100%;
      padding: 14px 16px;
      background: #1a1a1a;
      color: #fafaf6;
      border: none;
      font-family: ${dt('jafa.fontMono')};
      font-size: 13px;
      font-weight: 700;
      letter-spacing: 0.24em;
      cursor: pointer;
      transition: background 0.15s, transform 0.05s;
      text-transform: uppercase;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 10px;
      white-space: nowrap;
    }
    .receipt-btn:hover:not(:disabled) { background: #d97706; }
    .receipt-btn:active { transform: translateY(1px); }
    .receipt-btn:disabled { opacity: 0.5; cursor: not-allowed; }

    .receipt-link {
      color: #d97706;
      text-decoration: none;
      font-weight: 700;
      border-bottom: 1.5px solid #d97706;
      padding-bottom: 1px;
    }
    .receipt-link:hover { background: #d97706; color: #fafaf6; }

    @keyframes drift {
      0%   { transform: translateY(0);     opacity: 0; }
      10%  { opacity: 0.25; }
      90%  { opacity: 0.25; }
      100% { transform: translateY(-60vh); opacity: 0; }
    }
    .drift { animation: drift linear infinite; }
  `,

  components: {
    panel: {
      header: { padding: '0.55rem 0.875rem' },
      toggleableHeader: { padding: '0.3rem 0.875rem' },
      content: { padding: '0.875rem' },
      footer: { padding: '0 0.875rem 0.875rem 0.875rem' },
    },

    datatable: {
      headerCell: { padding: '0.5rem 0.75rem' },
      bodyCell: { padding: '0.55rem 0.75rem' },
      footerCell: { padding: '0.45rem 0.75rem' },
    },

    card: {
      root: { shadow: 'none', borderRadius: '{border.radius.lg}' },
      body: { padding: '1rem 1.125rem', gap: '0.25rem' },
      title: { fontSize: '0.7rem', fontWeight: '600' },
      subtitle: { color: '{text.muted.color}' },
    },

    chip: {
      root: {
        borderRadius: '999px',
        paddingX: '0.55rem',
        paddingY: '0.15rem',
        gap: '0.35rem',
      },
    },

    button: {
      root: {
        sm: { paddingX: '0.625rem', paddingY: '0.35rem' },
        borderRadius: '{border.radius.md}',
        letterSpacing: '0.02em',
      },
      label: { fontWeight: '600' },
    },

    dialog: { root: { borderRadius: '{border.radius.lg}' } },
    select: { root: { borderRadius: '{border.radius.md}' } },
    inputtext: { root: { borderRadius: '{border.radius.md}' } },
    inputnumber: { root: { borderRadius: '{border.radius.md}' } },
    password: { root: { borderRadius: '{border.radius.md}' } },
    menu: {
      root: { borderRadius: '{border.radius.md}' },
      item: { padding: '0.5rem 0.75rem' },
    },
  },
})

const app = createApp(App)

app.use(PrimeVue, {
  theme: {
    preset: JafaPreset,
    options: {
      darkModeSelector: '.p-dark',
    },
  },
})
app.use(createPinia())
app.use(VueQueryPlugin, { queryClient })
app.use(router)

// Initialize theme store so CSS vars apply before mount
import { useThemeStore } from '@/stores/theme'
useThemeStore()

app.mount('#app')
