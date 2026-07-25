import {defineConfig} from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  timeout: 30_000,
  use: {baseURL: 'http://127.0.0.1:7331', channel: 'chrome', trace: 'retain-on-failure'},
  webServer: {
    command: 'go run ../../../../cmd/forge workbench serve',
    cwd: '.',
    url: 'http://127.0.0.1:7331/api/v1/overview',
    reuseExistingServer: true,
    timeout: 120_000
  }
})
