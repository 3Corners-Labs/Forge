import AxeBuilder from '@axe-core/playwright'
import {expect, test} from '@playwright/test'

test('moves one outcome from conversation to a reviewable verified candidate', async ({page}) => {
  await page.goto('/')
  await expect(page.getByRole('heading', {name:'What should we build next?'})).toBeVisible()
  await expect(page.getByText('Intent first · no execution authority')).toBeVisible()

  await page.getByRole('button', {name:'Start discovery →'}).click()
  await expect(page.getByText('Your judgement is needed')).toBeVisible()
  await page.getByRole('button', {name:'Review & approve →'}).click()
  await expect(page.getByText('Approve bounded plan')).toBeVisible()
  await page.getByRole('button', {name:'Review & approve →'}).click()
  await expect(page.getByText('Publish reviewable change')).toBeVisible()
  await page.getByRole('button', {name:'Review & approve →'}).click()

  await expect(page.getByText('No decisions waiting')).toBeVisible()
  await expect(page.getByText('100%', {exact:true})).toBeVisible()
  await page.getByRole('tab', {name:'tests'}).click()
  await expect(page.getByText('2 passed')).toBeVisible()
  await page.getByRole('tab', {name:'evidence'}).click()
  await expect(page.getByText('Publication', {exact:true})).toBeVisible()

  const violations = await new AxeBuilder({page}).analyze()
  expect(violations.violations).toEqual([])
})
