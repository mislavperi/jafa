/**
 * E2E tests for navigation and sidebar:
 * - Sidebar links navigate to correct pages
 * - Page titles are correct
 * - Sidebar collapse/expand works
 * - Dark mode toggle is present
 *
 * These tests require authentication. They perform login before the suite runs.
 */

const TEST_USER = {
  username: process.env.E2E_USERNAME ?? 'testuser',
  password: process.env.E2E_PASSWORD ?? 'testpassword',
}

describe('Navigation', function () {
  before(function (browser) {
    browser
      .init()
      .navigateTo('/login')
      .waitForElementVisible('input[autocomplete="username"]', 5000)
      .setValue('input[autocomplete="username"]', TEST_USER.username)
      .setValue('input[autocomplete="current-password"]', TEST_USER.password)
      .click('button[type="submit"]')
      .waitForElementVisible('[data-tour="nav"]', 8000)
  })

  after(function (browser) {
    browser.end()
  })

  it('shows sidebar navigation after login', function (browser) {
    browser
      .assert.elementPresent('[data-tour="nav"]')
      .assert.textContains('[data-tour="nav"]', 'Dashboard')
      .assert.textContains('[data-tour="nav"]', 'Expenses')
      .assert.textContains('[data-tour="nav"]', 'Reports')
      .assert.textContains('[data-tour="nav"]', 'Categories')
      .assert.textContains('[data-tour="nav"]', 'Settings')
  })

  it('navigates to dashboard and shows correct page title', function (browser) {
    browser
      .click('[data-tour="nav"] a[href="/"]')
      .waitForElementVisible('[data-tour="stat-cards"]', 5000)
      .assert.urlEquals(browser.launchUrl + '/')
      .assert.textContains('body', 'Dashboard')
  })

  it('navigates to expenses page', function (browser) {
    browser
      .click('[data-tour="nav"] a[href="/expenses"]')
      .waitForElementVisible('body', 3000)
      .assert.urlContains('/expenses')
      .assert.textContains('body', 'Expenses')
  })

  it('navigates to reports page', function (browser) {
    browser
      .click('[data-tour="nav"] a[href="/reports"]')
      .waitForElementVisible('body', 3000)
      .assert.urlContains('/reports')
      .assert.textContains('body', 'Reports')
  })

  it('navigates to categories page', function (browser) {
    browser
      .click('[data-tour="nav"] a[href="/categories"]')
      .waitForElementVisible('body', 3000)
      .assert.urlContains('/categories')
      .assert.textContains('body', 'Categories')
  })

  it('navigates to settings page', function (browser) {
    browser
      .click('[data-tour="nav"] a[href="/settings"]')
      .waitForElementVisible('body', 3000)
      .assert.urlContains('/settings')
      .assert.textContains('body', 'Settings')
  })

  it('collapses sidebar and hides nav labels', function (browser) {
    browser
      .navigateTo('/')
      .waitForElementVisible('[data-tour="nav"]', 5000)
      .assert.textContains('[data-tour="nav"]', 'Dashboard')
      .click('[data-testid="sidebar-toggle"]')
      .pause(400)
      .assert.not.textContains('[data-tour="nav"]', 'Dashboard')
  })

  it('expands sidebar and shows nav labels', function (browser) {
    browser
      .click('[data-testid="sidebar-toggle"]')
      .pause(400)
      .assert.textContains('[data-tour="nav"]', 'Dashboard')
  })

  it('shows dark mode toggle button', function (browser) {
    browser
      .navigateTo('/')
      .waitForElementVisible('[data-tour="nav"]', 5000)
      .assert.elementPresent('nav button .pi-moon, nav button .pi-sun')
  })
})

// Mark this file as a module so its top-level declarations are scoped
// to the file (Nightwatch loads it for its describe() side-effects).
export {}
