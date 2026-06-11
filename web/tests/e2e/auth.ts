/**
 * E2E tests for authentication flows:
 * - Login page renders correctly
 * - Validation errors on empty submit
 * - Invalid credentials show error
 * - Register page renders and validates
 * - Navigation between login and register
 * - Unauthenticated users are redirected to /login
 */

describe('Authentication', function () {
  before(function (browser) {
    browser.init()
  })

  after(function (browser) {
    browser.end()
  })

  describe('Login page', function () {
    before(function (browser) {
      browser.navigateTo('/login')
    })

    it('renders the login form', function (browser) {
      browser
        .waitForElementVisible('input[autocomplete="username"]', 5000)
        .assert.elementPresent('input[autocomplete="current-password"]')
        .assert.textContains('button[type="submit"]', 'Sign In')
    })

    it('shows receipt-themed UI elements', function (browser) {
      browser
        .assert.textContains('body', 'Customer Sign-In')
        .assert.textContains('body', 'Please enter your credentials')
    })

    it('shows validation errors when submitting empty form', function (browser) {
      browser
        .click('button[type="submit"]')
        .waitForElementVisible('.text-red-700', 3000)
        .assert.textContains('.text-red-700', 'required')
    })

    it('links to the register page', function (browser) {
      browser
        .assert.elementPresent('a[href="/register"]')
        .click('a[href="/register"]')
        .waitForElementVisible('input[autocomplete="given-name"]', 5000)
        .assert.urlContains('/register')
    })
  })

  describe('Register page', function () {
    before(function (browser) {
      browser.navigateTo('/register')
    })

    it('renders the registration form', function (browser) {
      browser
        .waitForElementVisible('input[autocomplete="username"]', 5000)
        .assert.elementPresent('input[autocomplete="new-password"]')
        .assert.elementPresent('input[autocomplete="email"]')
        .assert.elementPresent('input[autocomplete="given-name"]')
        .assert.elementPresent('input[autocomplete="family-name"]')
    })

    it('shows receipt-themed UI elements', function (browser) {
      browser
        .assert.textContains('body', 'Open an Account')
        .assert.textContains('body', 'Create Account')
    })

    it('shows validation errors when submitting without required fields', function (browser) {
      browser
        .click('button[type="submit"]')
        .waitForElementVisible('.text-red-700', 3000)
        .assert.textContains('.text-red-700', 'required')
    })

    it('shows email validation error for invalid email', function (browser) {
      browser
        .clearValue('input[autocomplete="email"]')
        .setValue('input[autocomplete="email"]', 'not-an-email')
        .setValue('input[autocomplete="username"]', 'testuser')
        .setValue('input[autocomplete="new-password"]', 'testpass')
        .click('button[type="submit"]')
        .waitForElementVisible('.text-red-700', 3000)
        .assert.textContains('.text-red-700', 'valid email')
    })

    it('links back to login page', function (browser) {
      browser
        .assert.elementPresent('a[href="/login"]')
        .click('a[href="/login"]')
        .waitForElementVisible('input[autocomplete="current-password"]', 5000)
        .assert.urlContains('/login')
    })
  })

  describe('Route guards', function () {
    it('redirects unauthenticated users to login when accessing dashboard', function (browser) {
      browser
        .navigateTo('/')
        .waitForElementVisible('input[autocomplete="username"]', 5000)
        .assert.urlContains('/login')
    })

    it('redirects unauthenticated users to login when accessing expenses', function (browser) {
      browser
        .navigateTo('/expenses')
        .waitForElementVisible('input[autocomplete="username"]', 5000)
        .assert.urlContains('/login')
    })

    it('redirects unauthenticated users to login when accessing reports', function (browser) {
      browser
        .navigateTo('/reports')
        .waitForElementVisible('input[autocomplete="username"]', 5000)
        .assert.urlContains('/login')
    })
  })
})
