import { strict as assert } from 'node:assert'
import { test } from 'node:test'
import React from 'react'
import { createInertiaApp } from '@inertiajs/react'
import { renderToString } from 'react-dom/server'
import WorkIndex from '../resources/js/Pages/Public/WorkIndex'
import WorkShow from '../resources/js/Pages/Public/WorkShow'
import ProjectIndex from '../resources/js/Pages/Public/ProjectIndex'
import ProjectShow from '../resources/js/Pages/Public/ProjectShow'

const nil = { String: '', Valid: false }
const noDate = { Time: '0001-01-01T00:00:00Z', Valid: false }
const date = { Time: '2023-02-01T00:00:00+02:00', Valid: true }
const work = { ID: 1, Title: 'Client system', Slug: 'client-system', Summary: 'Summary', ClientName: nil, ClientIndustry: nil, ClientURL: nil, ClientLogoURL: nil, CoverImageURL: nil, Specialisms: ['Engineering'], Platforms: null, Technologies: ['Go'], Challenge: 'Challenge', Approach: '', Deliverables: '', Outcome: '', Content: 'Notes', StartedAt: noDate, CompletedAt: date }
const project = { ID: 2, Name: 'A project', Slug: 'a-project', Tagline: 'Tagline', ProjectType: 'web-app', Description: { String: 'Description', Valid: true }, RepositoryURL: { String: 'https://example.com/repo', Valid: true }, LiveURL: nil, ImageURL: nil, Technologies: [' Go ', ''], StartedAt: noDate, LaunchedAt: date }
const shared = { seo: { title: 'Page', description: '', canonical: 'https://example.com', image: '', imageAlt: '', metaType: 'website' }, publicSite: { production: false, year: 2026 } }

async function render(component: typeof WorkIndex | typeof WorkShow | typeof ProjectIndex | typeof ProjectShow, props: object) {
  const result = await createInertiaApp({
    page: { component: 'Test', props: { ...shared, ...props }, url: '/', version: null },
    resolve: () => component,
    render: renderToString,
    setup: ({ App, props }) => <App {...props}/>,
  })
  return result.body
}

test('public work and project pages render raw entities and sanitized markdown on the server', async () => {
  const workIndex = await render(WorkIndex, { items: [work] })
  assert.match(workIndex, /href="\/work\/client-system"/)
  assert.match(workIndex, /February 2023/)
  assert.match(workIndex, /Parkeringskompagniet/)
  assert.match(await render(WorkIndex, { items: [] }), /No work is published yet/)
  const workShow = await render(WorkShow, { item: work, markdownHTML: '<p>Safe <strong>notes</strong></p>' })
  assert.match(workShow, /Engineering, Go/)
  assert.match(workShow, /<strong>notes<\/strong>/)
  const projectIndex = await render(ProjectIndex, { items: [project] })
  assert.match(projectIndex, /href="\/projects\/a-project"/)
  assert.match(projectIndex, /Open source/)
  assert.match(await render(ProjectIndex, { items: null }), /No projects are listed yet/)
  const projectShow = await render(ProjectShow, { item: project, markdownHTML: '<p>Safe description</p>' })
  assert.match(projectShow, /Safe description/)
  assert.match(projectShow, /Private or not publicly hosted/)
  assert.match(projectShow, /https:\/\/example.com\/repo/)
})
