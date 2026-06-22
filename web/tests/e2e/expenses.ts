/**
 * E2E tests for the Expenses list page:
 * - Page renders correctly with title and subtitle
 * - Add Expense button opens the modal
 * - Scan Receipt button is present
 * - Recurring option in the modal works correctly
 */

const TEST_USER = {
  username: process.env.E2E_USERNAME ?? 'testuser',
  password: process.env.E2E_PASSWORD ?? 'testpassword',
}

describe('Expenses List Page', function () {
  before(function (browser) {
    browser
      .init()
      .navigateTo('/login')
      .waitForElementVisible('input[autocomplete="username"]', 5000)
      .setValue('input[autocomplete="username"]', TEST_USER.username)
      .setValue('input[autocomplete="current-password"]', TEST_USER.password)
      .click('button[type="submit"]')
      .waitForElementVisible('[data-tour="nav"]', 8000)
      .navigateTo('/expenses')
      .waitForElementVisible('body', 3000)
  })

  after(function (browser) {
    browser.end()
  })

  it('shows the expenses page heading', function (browser) {
    browser
      .assert.textContains('body', 'Expenses')
      .assert.textContains('body', 'All your transactions')
  })

  it('has an Add Expense button', function (browser) {
    browser.assert.elementPresent('[data-testid="add-expense-btn"]')
  })

  it('has a Scan Receipt button', function (browser) {
    browser.assert.textContains('body', 'Scan Receipt')
  })

  describe('Add Expense modal from Expenses page', function () {
    it('opens Add Expense modal', function (browser) {
      browser
        .click('[data-testid="add-expense-btn"]')
        .waitForElementVisible('[aria-label="Close"]', 5000)
        .assert.textContains('body', 'New Transaction')
    })

    it('shows recurring charge option', function (browser) {
      browser.assert.textContains('body', 'Recurring Charge')
    })

    it('shows frequency options after enabling recurring', function (browser) {
      browser
        .click('input[type="checkbox"]')
        .pause(300)
        .assert.textContains('body', 'monthly')
        .assert.textContains('body', 'yearly')
        .assert.textContains('body', 'Day of Month')
    })

    it('hides frequency options after disabling recurring', function (browser) {
      browser
        .click('input[type="checkbox"]')
        .pause(300)
        .assert.not.textContains('body', 'Day of Month')
    })

    it('validates required fields before submission', function (browser) {
      browser
        .click('button[type="submit"]')
        .waitForElementVisible('.text-red-500', 3000)
        .assert.textContains('.text-red-500', 'required')
    })

    it('shows split payment option', function (browser) {
      browser.assert.textContains('body', 'Split Payment')
    })

    it('shows number-of-payments field after enabling split', function (browser) {
      browser
        .click('[data-testid="split-toggle"]')
        .pause(300)
        .assert.textContains('body', 'Number of Payments')
        .assert.elementPresent('[data-testid="installment-count"]')
        .assert.elementPresent('[data-testid="per-payment"]')
    })

    it('computes the per-payment amount from cost and payment count', function (browser) {
      browser
        .clearValue('input[type="text"]')
        .setValue('input[type="text"]', 'Phone')
        .setValue('input[type="number"]', '200')
        .clearValue('[data-testid="installment-count"]')
        .setValue('[data-testid="installment-count"]', '4')
        .pause(300)
        // $200 split into 4 → $50.00 × 4
        .assert.textContains('[data-testid="per-payment"]', '50.00')
        .assert.textContains('[data-testid="per-payment"]', '× 4')
    })

    it('hides number-of-payments field after disabling split', function (browser) {
      browser
        .click('[data-testid="split-toggle"]')
        .pause(300)
        .assert.not.textContains('body', 'Number of Payments')
    })

    it('closes the modal', function (browser) {
      browser
        .click('[aria-label="Close"]')
        .waitForElementNotPresent('[aria-label="Close"]', 3000)
    })
  })
})
