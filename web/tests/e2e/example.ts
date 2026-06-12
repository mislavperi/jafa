/**
 * Smoke test: verifies the app loads and redirects unauthenticated users to /login.
 */
describe('Smoke Test', function () {
  before(function (browser) {
    browser.init()
  })

  it('loads and shows the login page for unauthenticated users', function (browser) {
    browser
      .waitForElementVisible('input[autocomplete="username"]', 5000)
      .assert.urlContains('/login')
      .assert.textContains('body', 'Customer Sign-In')
  })

  after(function (browser) {
    browser.end()
  })
})
