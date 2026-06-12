/**
 * E2E tests for the Dashboard page:
 * - Stat cards are displayed
 * - "Add Expense" button opens the modal
 * - Expense modal has correct fields and validates required inputs
 * - Modal can be closed without submitting
 * - Recent expenses and breakdown sections are present
 */

const TEST_USER = {
  username: process.env.E2E_USERNAME ?? 'testuser',
  password: process.env.E2E_PASSWORD ?? 'testpassword',
}

describe('Dashboard', function () {
  before(function (browser) {
    browser
      .init()
      .navigateTo('/login')
      .waitForElementVisible('input[autocomplete="username"]', 5000)
      .setValue('input[autocomplete="username"]', TEST_USER.username)
      .setValue('input[autocomplete="current-password"]', TEST_USER.password)
      .click('button[type="submit"]')
      .waitForElementVisible('[data-tour="stat-cards"]', 8000)
  })

  after(function (browser) {
    browser.end()
  })

  it('shows the dashboard heading', function (browser) {
    browser
      .assert.textContains('body', 'Dashboard')
      .assert.textContains('body', 'Your spending at a glance')
  })

  it('renders stat cards', function (browser) {
    browser
      .assert.elementPresent('[data-tour="stat-cards"]')
      .assert.textContains('[data-tour="stat-cards"]', 'Budget')
      .assert.textContains('[data-tour="stat-cards"]', 'Left to spend')
  })

  it('shows the recent expenses section', function (browser) {
    browser
      .assert.elementPresent('[data-tour="recent-expenses"]')
      .assert.textContains('[data-tour="recent-expenses"]', 'Recent Expenses')
  })

  it('shows the expense breakdown section', function (browser) {
    browser
      .assert.elementPresent('[data-tour="breakdown"]')
      .assert.textContains('[data-tour="breakdown"]', 'Expense Breakdown')
  })

  it('shows the upcoming bills section', function (browser) {
    browser
      .assert.elementPresent('[data-tour="upcoming-bills"]')
      .assert.textContains('[data-tour="upcoming-bills"]', 'Upcoming Bills')
  })

  it('has Add Expense and Scan Receipt action buttons', function (browser) {
    browser
      .assert.elementPresent('[data-tour="add-expense"]')
      .assert.elementPresent('[data-tour="scan-receipt"]')
      .assert.textContains('[data-tour="add-expense"]', 'Add Expense')
      .assert.textContains('[data-tour="scan-receipt"]', 'Scan Receipt')
  })

  describe('Add Expense modal', function () {
    it('opens the add expense modal on button click', function (browser) {
      browser
        .click('[data-tour="add-expense"]')
        .waitForElementVisible('[aria-label="Close"]', 5000)
        .assert.textContains('body', 'New Transaction')
        .assert.textContains('body', 'Item Description')
        .assert.textContains('body', 'Unit Price')
    })

    it('shows validation errors when submitting empty modal', function (browser) {
      browser
        .click('button[type="submit"]')
        .waitForElementVisible('.text-red-500', 3000)
        .assert.textContains('.text-red-500', 'required')
    })

    it('fills in expense name and cost', function (browser) {
      browser
        .clearValue('input[placeholder="e.g. Weekly groceries"]')
        .setValue('input[placeholder="e.g. Weekly groceries"]', 'Test Expense')
        .clearValue('input[type="number"]')
        .setValue('input[type="number"]', '25.50')
        .assert.value('input[placeholder="e.g. Weekly groceries"]', 'Test Expense')
    })

    it('shows correct total after entering cost', function (browser) {
      browser.assert.textContains('body', '25.50')
    })

    it('closes the modal via close button without saving', function (browser) {
      browser
        .click('[aria-label="Close"]')
        .waitForElementNotPresent('[aria-label="Close"]', 3000)
        .assert.not.textContains('body', 'New Transaction')
    })

    it('closes the modal via void transaction button', function (browser) {
      browser
        .click('[data-tour="add-expense"]')
        .waitForElementVisible('[aria-label="Close"]', 5000)
        .assert.textContains('body', 'New Transaction')
        .useXpath()
        .click('//button[contains(.,"Void transaction")]')
        .useCss()
        .waitForElementNotPresent('[aria-label="Close"]', 3000)
    })
  })

  describe('Recurring expense option', function () {
    before(function (browser) {
      browser
        .click('[data-tour="add-expense"]')
        .waitForElementVisible('[aria-label="Close"]', 5000)
    })

    after(function (browser) {
      browser
        .click('[aria-label="Close"]')
        .waitForElementNotPresent('[aria-label="Close"]', 3000)
    })

    it('shows recurring charge toggle', function (browser) {
      browser.assert.textContains('body', 'Recurring Charge')
    })

    it('shows frequency options after enabling recurring toggle', function (browser) {
      browser
        .click('input[type="checkbox"]')
        .pause(300)
        .assert.textContains('body', 'monthly')
        .assert.textContains('body', 'yearly')
        .assert.textContains('body', 'Day of Month')
    })

    it('hides frequency options after disabling recurring toggle', function (browser) {
      browser
        .click('input[type="checkbox"]')
        .pause(300)
        .assert.not.textContains('body', 'Day of Month')
    })
  })
})
