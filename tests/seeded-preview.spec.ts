import { expect, test } from '@playwright/test'

const preview = process.env.MBVLABS_PREVIEW_URL

test.skip(!preview, 'Set MBVLABS_PREVIEW_URL to a local preview seeded with demo accounts')

test('server-rendered public contact hydrates and opens the mobile menu', async ({ page, request }) => {
  const errors: string[] = []
  page.on('pageerror', error => errors.push(error.message))
  const response = await request.get(`${preview}/contact`)
  expect(response.ok()).toBe(true)
  const html = await response.text()
  expect(html).toContain('data-server-rendered="true"')
  expect(html).toContain('data-slot="input"')

  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto(`${preview}/contact`)
  const menuButton = page.getByRole('button', { name: 'Open menu' })
  await menuButton.click()
  await expect(page.getByRole('navigation', { name: 'Mobile navigation' })).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(menuButton).toBeFocused()
  await expect(menuButton).toHaveAttribute('aria-expanded', 'false')
  await expect(page.locator('form select').first()).toBeVisible()
  expect(errors).toEqual([])
})

test('sign-in and the shadcn admin sidebar work in Chromium', async ({ page }) => {
  const errors: string[] = []
  page.on('pageerror', error => errors.push(error.message))
  await page.goto(`${preview}/admin`)
  await page.getByLabel(/email/i).fill('admin@example.com')
  await page.getByLabel(/password/i).fill('password123')
  await page.getByRole('button', { name: 'Login' }).click()
  await expect(page.getByRole('heading', { name: 'Latest Inquiries' })).toBeVisible()
  const sidebar = page.locator('[data-slot="sidebar"][data-state]')
  const sidebarTrigger = page.locator('[data-slot="sidebar-trigger"]')
  await expect(sidebar).toHaveAttribute('data-state', 'expanded')
  await expect(sidebar.locator('[data-slot="sidebar-menu-button"][data-active]')).toHaveCount(1)
  await sidebarTrigger.click()
  await expect(sidebar).toHaveAttribute('data-state', 'collapsed')
  await page.getByRole('navigation', { name: 'Platform' }).getByRole('link', { name: 'Blog Posts' }).click()
  await expect(page).toHaveURL(/\/admin\/blog-posts$/)
  const post = page.getByRole('link', { name: /^View blog post / }).first()
  if (await post.count()) {
    await post.click()
    const trigger = page.getByRole('button', { name: 'Delete', exact: true })
    await trigger.click()
    const dialog = page.getByRole('alertdialog')
    const confirm = dialog.getByRole('button', { name: 'Delete blog post' })
    await expect(dialog).toBeVisible()
    await expect(confirm).toBeDisabled()
    await dialog.getByLabel('Type DELETE to confirm').fill('DELETE')
    await expect(confirm).toBeEnabled()
    await page.keyboard.press('Escape')
    await expect(dialog).toBeHidden()
    await expect(trigger).toBeFocused()
  }
  await page.setViewportSize({ width: 390, height: 844 })
  await sidebarTrigger.click()
  const sheet = page.getByRole('dialog', { name: 'Sidebar' })
  await expect(sheet).toBeVisible()
  await sheet.getByRole('link', { name: 'Diary' }).click()
  await expect(page).toHaveURL(/\/admin\/diary$/)
  await expect(sheet).toBeHidden()
  expect(errors).toEqual([])
})

test('admin can issue and revoke a content API token', async ({ page, request }) => {
  await page.goto(`${preview}/admin`)
  await page.getByLabel(/email/i).fill('admin@example.com')
  await page.getByLabel(/password/i).fill('password123')
  await page.getByRole('button', { name: 'Login' }).click()
  await page.getByRole('link', { name: 'Manage token', exact: true }).click()
  await expect(page.getByRole('heading', { name: 'Content API token', exact: true })).toBeVisible()

  const articles = `${preview}/api/content/articles`
  expect((await request.get(articles)).status()).toBe(401)
  await page.getByRole('button', { name: 'Create token', exact: true }).click()
  const revealed = page.getByLabel('New token (copy now)', { exact: true })
  await expect(revealed).toBeVisible()
  const token = await revealed.inputValue()
  expect(token.length).toBeGreaterThan(0)
  const headers = { Authorization: `Bearer ${token}` }
  expect((await request.get(articles, { headers })).status()).toBe(200)
  expect((await request.patch(`${articles}/missing-preview-article`, { headers, data: { title: 'Missing preview article' } })).status()).toBe(404)

  await page.reload()
  await expect(revealed).toHaveCount(0)
  await page.getByRole('button', { name: 'Revoke token', exact: true }).click()
  await page.getByRole('alertdialog').getByRole('button', { name: /revoke/i }).click()
  await expect(page.getByRole('status')).toHaveText('Token revoked.')
  expect((await request.get(articles, { headers })).status()).toBe(401)
})
