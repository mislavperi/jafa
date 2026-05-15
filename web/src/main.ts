import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { queryClient } from './core/query'
import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura'
import { definePreset } from '@primeuix/themes'

import App from './App.vue'
import router from './router'

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
          background: '{surface.800}',
          hoverBackground: '{surface.700}',
          borderColor: '{surface.600}',
          color: '{text.color}',
          hoverColor: '{text.hover.color}',
        },
      },
    },
  },

  components: {
    panel: {
      header: {
        padding: '0.55rem 0.875rem',
      },
      toggleableHeader: {
        padding: '0.3rem 0.875rem',
      },
      content: {
        padding: '0.875rem',
      },
      footer: {
        padding: '0 0.875rem 0.875rem 0.875rem',
      },
    },

    datatable: {
      headerCell: {
        padding: '0.45rem 0.75rem',
      },
      bodyCell: {
        padding: '0.45rem 0.75rem',
      },
      footerCell: {
        padding: '0.45rem 0.75rem',
      },
    },

    card: {
      root: {
        shadow: 'none',
        borderRadius: '{border.radius.md}',
      },
      body: {
        padding: '0.875rem',
        gap: '0.25rem',
      },
      title: {
        fontSize: '0.7rem',
        fontWeight: '700',
      },
      subtitle: {
        color: '{text.muted.color}',
      },
    },

    chip: {
      root: {
        borderRadius: '4px',
        paddingX: '0.45rem',
        paddingY: '0.15rem',
        gap: '0.25rem',
      },
    },

    button: {
      root: {
        sm: {
          paddingX: '0.5rem',
          paddingY: '0.3rem',
        },
      },
      label: {
        fontWeight: '500',
      },
    },

    dialog: {
      root: {
        borderRadius: '{border.radius.md}',
      },
    },

    select: {
      root: {
        borderRadius: '{border.radius.md}',
      },
    },

    inputtext: {
      root: {
        borderRadius: '{border.radius.md}',
      },
    },

    inputnumber: {
      root: {
        borderRadius: '{border.radius.md}',
      },
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

app.mount('#app')
