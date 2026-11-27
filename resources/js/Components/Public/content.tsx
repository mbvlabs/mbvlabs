import type { ReactNode } from 'react'

export type NullString = { String: string; Valid: boolean }
export type NullTime = { Time: string; Valid: boolean }
export type Work = {
  ID: number; Title: string; Slug: string; Summary: string
  ClientName: NullString; ClientIndustry: NullString; ClientURL: NullString; ClientLogoURL: NullString; CoverImageURL: NullString
  Specialisms: string[] | null; Platforms: string[] | null; Technologies: string[] | null
  Challenge: string; Approach: string; Deliverables: string; Outcome: string; Content: string
  StartedAt: NullTime; CompletedAt: NullTime
}
export type Project = {
  ID: number; Name: string; Slug: string; Tagline: string; ProjectType: string
  Description: NullString; RepositoryURL: NullString; LiveURL: NullString; ImageURL: NullString
  Technologies: string[] | null; StartedAt: NullTime; LaunchedAt: NullTime
}
export const value = (v: NullString) => v?.Valid ? v.String : ''
export const initial = (s: string) => Array.from(s.trim())[0]?.toUpperCase() || 'M'
// Go formats time in its original offset, not in the browser's local timezone.
export function dateLabel(first: NullTime, second: NullTime, yearOnly = false) {
  const date = first?.Valid && !first.Time.startsWith('0001-01-01') ? first.Time : second?.Valid && !second.Time.startsWith('0001-01-01') ? second.Time : ''
  if (!date) return ''
  const match = /^(\d{4})-(\d{2})/.exec(date)
  if (!match) return ''
  return yearOnly ? match[1] : `${['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'][Number(match[2]) - 1]} ${match[1]}`
}
export const projectType = (type: string) => ({ 'web-app': 'Web app', website: 'Website', automation: 'Automation', integration: 'Integration', consulting: 'Consulting', framework: 'Framework' }[type] || type.trim())
export const technologies = (items: string[] | null) => Array.isArray(items) && items.every(item => typeof item === 'string') ? items.map(item => item.trim()).filter(Boolean) : []
export const number = (index: number) => `${String(index + 1).padStart(2, '0')}.`
export const prose = 'prose prose-sm mt-8 max-w-3xl text-foreground prose-headings:font-heading prose-headings:font-normal prose-headings:text-foreground prose-p:text-foreground prose-a:text-foreground prose-a:underline prose-strong:text-foreground prose-code:text-foreground prose-li:text-foreground prose-hr:border-foreground prose-img:rounded-none prose-pre:rounded-none prose-pre:border prose-pre:border-foreground prose-pre:bg-primary prose-pre:text-primary-foreground prose-blockquote:border-foreground prose-blockquote:text-foreground'
export const button = 'inline-flex min-h-11 items-center border border-foreground px-4 py-2 font-semibold transition-colors hover:bg-primary hover:text-primary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
export function Breadcrumbs({ section, detail }: { section: 'Work' | 'Projects'; detail?: string }) {
  const items = [{ name: 'Home', path: '/' }, { name: section, path: section === 'Work' ? '/work' : '/projects' }, ...(detail ? [{ name: detail, path: '' }] : [])]
  return <nav aria-label="Breadcrumb"><ol className="flex flex-wrap items-center gap-x-1.5 gap-y-1 font-mono text-xs uppercase tracking-[0.18em] text-muted-foreground">{items.map((item, i) => <li key={item.name} className="inline-flex items-center gap-1.5">{i > 0 && <span aria-hidden="true">/</span>}{i === items.length - 1 ? <span className="font-medium text-foreground" aria-current="page">{item.name}</span> : <a href={item.path} className="transition-colors hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">{item.name}</a>}</li>)}</ol></nav>
}
export function Tags({ items }: { items: string[] }) {
  return <div className="mt-7 flex max-w-3xl flex-wrap gap-2">{items.map((item, i) => <span key={i} className="border border-border px-3 py-2 font-mono text-xs uppercase leading-5 tracking-[0.08em]">{item}</span>)}</div>
}
export function DetailSection({ label, title, children }: { label: string; title: string; children: ReactNode }) {
  return <section className="border-t border-foreground py-10"><p className="font-mono text-xs uppercase tracking-[0.14em]">{label}</p><h2 className="cobalt-row-title mt-5">{title}</h2>{children}</section>
}
export function Continue({ project }: { project?: boolean }) {
  return <section className="border-y border-foreground"><div className="container mx-auto px-6 py-20 md:px-8 md:py-24"><p className="font-mono text-xs uppercase tracking-[0.18em]">Archive / Continue reading</p><h2 className="cobalt-section-title mt-5">More {project ? 'projects' : 'work'}.</h2><a href={project ? '/projects' : '/work'} className="group mt-10 grid gap-6 border-y border-border py-10 transition-colors hover:bg-secondary/45 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring md:grid-cols-[minmax(0,1fr)_8rem] md:items-end md:px-3"><div><p className="max-w-3xl text-sm leading-6">{project ? 'Explore open source work, commercial products, automations, internal tools, and software experiments built in and through MBV Labs.' : 'Explore selected client projects, technical leadership work, and full-time engineering experience across product delivery, architecture, backend systems, frontend applications, and operational tooling.'}</p><h3 className="cobalt-row-title mt-6">{project ? 'See all projects' : 'See all my work'}</h3></div><span className="text-xs font-semibold uppercase tracking-[0.12em] transition-transform group-hover:translate-x-1 md:text-right">Open index →</span></a></div></section>
}
