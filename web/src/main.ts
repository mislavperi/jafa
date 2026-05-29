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
