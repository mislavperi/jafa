import { describe, it, expect } from 'vitest'

import { mount } from '@vue/test-utils'
import App from '../App.vue'

describe('App', () => {
  it('mounts and renders the router view', () => {
    const wrapper = mount(App, {
      global: {
        stubs: {
          RouterView: { template: '<div data-testid="router-view" />' },
        },
      },
    })
    expect(wrapper.find('[data-testid="router-view"]').exists()).toBe(true)
  })
})
