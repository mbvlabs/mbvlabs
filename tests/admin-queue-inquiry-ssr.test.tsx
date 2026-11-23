import { strict as assert } from 'node:assert'
import { test } from 'node:test'
import React from 'react'
import { createInertiaApp } from '@inertiajs/react'
import { renderToString } from 'react-dom/server'
import InquiryForm from '../resources/js/Pages/Admin/ProjectInquiry/Form'
import InquiryShow from '../resources/js/Pages/Admin/ProjectInquiry/Show'
import QueueJobs from '../resources/js/Pages/Admin/Queue/Jobs'
import QueueShow from '../resources/js/Pages/Admin/Queue/Show'

const inquiry = { ID: 1, Name: 'Ada', Email: 'ada@example.com', Company: null, Role: null, ProjectType: 'website', Timeline: null, Message: 'Hello', Source: null, Status: 'new', Metadata: '{}', CreatedAt: '2026-01-01', UpdatedAt: '2026-01-01', DomainBanned: false }
const job = { ID: 2, State: 'retryable', Kind: 'mail', Queue: 'default', Attempt: 1, MaxAttempts: 3, CreatedAt: '2026-01-01', ScheduledAt: '2026-01-01', AttemptedAt: null, FinalizedAt: null, Priority: 1, Args: '{}', Metadata: '{}', AttemptedBy: null, Errors: null, Tags: null }

async function render(component: typeof InquiryForm | typeof InquiryShow | typeof QueueJobs | typeof QueueShow, props: object) {
  const result = await createInertiaApp({
    page: { component: 'Test', props: { ...props, errors: {} }, url: '/admin', version: null, rescuedProps: [], flash: {}, rememberedState: {} },
    resolve: () => component,
    render: renderToString,
    setup: ({ App, props: appProps }) => <App {...appProps} />,
  })
  return result.body
}

test('admin inquiries and queue render shadcn controls with existing actions and filters', async () => {
  const form = await render(InquiryForm, { item: inquiry })
  assert.match(form, /data-slot="input"[^>]*id="name"/)
  assert.match(form, /data-slot="native-select"[^>]*id="projectType"/)
  assert.match(form, /Update Project Inquiry/)
  const show = await render(InquiryShow, { item: inquiry })
  assert.match(show, /data-slot="alert-dialog-trigger"/)
  assert.match(show, /Ban domain/)
  assert.match(show, /Delete/)
  const jobs = await render(QueueJobs, { items: [job], queues: [], filters: { State: 'retryable', Queue: '', Kind: '' }, pagination: { Page: 1, PageSize: 20, TotalCount: 1, TotalPages: 1 } })
  assert.match(jobs, /data-slot="table"/)
  assert.match(jobs, /aria-label="State"/)
  assert.match(jobs, /\/admin\/queue\/jobs\/2/)
  const jobShow = await render(QueueShow, { item: job })
  assert.match(jobShow, /Pause polling/)
  assert.match(jobShow, /data-slot="alert-dialog-trigger"/)
  assert.match(jobShow, /Retry/)
})
