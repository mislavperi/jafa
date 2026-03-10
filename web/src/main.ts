import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { queryClient } from './core/query'
import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(PrimeVue, {
  theme: {
    preset: Aura,
    components: {
      button: {

      }
    },
    options: {
      darkModeSelector: ".p-dark"
    }
  },
})
app.use(createPinia())
app.use(VueQueryPlugin, { queryClient })
app.use(router)

app.mount('#app')
