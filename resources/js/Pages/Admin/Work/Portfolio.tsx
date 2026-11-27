import { routes } from '@/routes'
import { Head, Link, router, useForm } from '@inertiajs/react'
import { useState, type FormEvent, type ReactNode } from 'react'
import AdminLayout from '@/Layouts/AdminLayout'
import RichEditor from '@/Components/RichEditor'
import MarkdownPreview from '@/Components/MarkdownPreview'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Textarea } from '@/Components/ui/textarea'
import { Label } from '@/Components/ui/label'
import { NativeSelect } from '@/Components/ui/native-select'
import { Checkbox } from '@/Components/ui/checkbox'
import { Badge } from '@/Components/ui/badge'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/Components/ui/table'
import { AlertDialog, AlertDialogTrigger, AlertDialogContent, AlertDialogTitle, AlertDialogDescription, AlertDialogCancel, AlertDialogFooter } from '@/Components/ui/alert-dialog'
import type { AdminWorkData as Work, AdminProjectData as Project } from '@/types/payloads'

type Kind = 'Work' | 'Project'
const portfolioRoutes = {
  Work: { index: routes.adminWorkIndex, new: routes.adminWorkNew, show: routes.adminWorkShow, create: routes.adminWorkCreate, edit: routes.adminWorkEdit, update: routes.adminWorkUpdate, destroy: routes.adminWorkDestroy },
  Project: { index: routes.adminProjectIndex, new: routes.adminProjectNew, show: routes.adminProjectShow, create: routes.adminProjectCreate, edit: routes.adminProjectEdit, update: routes.adminProjectUpdate, destroy: routes.adminProjectDestroy },
}
const plural = (kind: Kind) => kind === 'Work' ? 'Work' : 'Projects'
const title = (item: Work | Project) => 'Title' in item ? item.Title : item.Name
const date = (value: string | null, full = false) => value ? (full ? new Date(value).toLocaleString() : new Date(value).toLocaleDateString()) : '—'
const button = 'normal-case tracking-normal inline-flex items-center justify-center border px-4 py-2 text-sm font-medium hover:bg-muted'
const input = 'flex w-full border bg-background px-3 py-2 text-sm'

function Frame({ children }: { children: ReactNode }) {
  return <AdminLayout><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6">{children}</div></AdminLayout>
}
function Header({ heading, description, children }: { heading: string; description?: string | null; children: ReactNode }) {
  return <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"><div><h1 className="text-2xl font-semibold tracking-tight">{heading}</h1>{description && <p className="text-muted-foreground">{description}</p>}</div><div className="flex gap-2">{children}</div></div>
}
export function PortfolioIndex({ kind, items }: { kind: Kind; items: (Work | Project)[] | null }) {
  const route = portfolioRoutes[kind]
  return <Frame><Head title={plural(kind)} /><Header heading={plural(kind)} description={`Manage your ${kind === 'Work' ? 'work' : 'projects'}.`}><Button asChild variant="outline" className={button}><Link href={route.new.path()}>+ New {kind}</Link></Button></Header>
    <div className="overflow-x-auto border bg-card"><Table className="w-full text-left text-sm"><TableHeader><TableRow className="border-b">{[kind === 'Work' ? 'Title' : 'Name', 'Slug', 'Status', 'Created', 'Updated', 'Published'].map(x => <TableHead className="h-10 px-4 text-sm font-medium normal-case tracking-normal text-muted-foreground" key={x}>{x}</TableHead>)}</TableRow></TableHeader><TableBody>
      {(items ?? []).map(item => <TableRow key={item.ID} className="cursor-pointer border-b hover:bg-muted/40" tabIndex={0} aria-label={`View ${kind.toLowerCase()} ${title(item)}`} onClick={() => router.visit(route.show.path(item.ID))} onKeyDown={e => { if (e.key === 'Enter') router.visit(route.show.path(item.ID)) }}><TableCell className="max-w-80 truncate p-4 font-medium"><Link href={route.show.path(item.ID)} onClick={e => e.stopPropagation()}>{title(item)}</Link></TableCell><TableCell className="max-w-64 truncate p-4 text-muted-foreground">{item.Slug}</TableCell><TableCell className="p-4"><Badge className="border px-2 py-0.5 text-xs font-semibold normal-case tracking-normal">{('Status' in item ? item.Status === 'published' : !!item.PublishedAt) ? 'Published' : 'Draft'}</Badge></TableCell><TableCell className="p-4 text-muted-foreground">{date(item.CreatedAt)}</TableCell><TableCell className="p-4 text-muted-foreground">{date(item.UpdatedAt)}</TableCell><TableCell className="p-4 text-muted-foreground">{date(item.PublishedAt)}</TableCell></TableRow>)}
      {!items?.length && <TableRow><TableCell colSpan={6} className="h-24 text-center text-muted-foreground">No {kind === 'Work' ? 'work' : 'projects'} found.</TableCell></TableRow>}
    </TableBody></Table></div></Frame>
}

const workDefaults = { title: '', slug: '', clientName: '', clientIndustry: '', clientUrl: '', clientLogoUrl: '', summary: '', content: '', coverImageUrl: '', specialisms: '', platforms: '', technologies: '', challenge: '', approach: '', deliverables: '', outcome: '', startedAt: '', completedAt: '', status: 'draft', isFeatured: false }
const projectDefaults = { name: '', slug: '', tagline: '', description: '', projectType: 'web-app', repositoryUrl: '', liveUrl: '', imageUrl: '', technologies: '', startedAt: '', launchedAt: '', publishedAt: '', isFeatured: false }
type Values = Record<string, string | boolean>
function initial(kind: Kind, item?: Work | Project): Values {
  const defaults: Values = kind === 'Work' ? workDefaults : projectDefaults
  if (!item) return { ...defaults }
  const source = item as unknown as Record<string, unknown>
  return Object.fromEntries(Object.entries(defaults).map(([key, fallback]) => {
    const value = source[key[0].toUpperCase() + key.slice(1)]
    return [key, Array.isArray(value) ? value.join(', ') : typeof value === 'string' || typeof value === 'boolean' ? value : fallback]
  }))
}
const slugify = (value: string) => value.toLowerCase().trim().replace(/[^\w\s-]/g, '').replace(/[\s_]+/g, '-').replace(/-+/g, '-').replace(/^-|-$/g, '')
type Field = { key: string; label: string; type?: 'textarea' | 'date' | 'select' | 'checkbox' | 'editor'; options?: string[]; placeholder?: string; wide?: boolean }
const f = (key: string, label: string, type?: Field['type'], extra?: Partial<Field>): Field => ({ key, label, type, ...extra })
const workSections: { name: string; fields: Field[] }[] = [
  { name: 'Primary details', fields: [f('title', 'Title'), f('slug', 'Slug'), f('clientName', 'Client Name'), f('clientIndustry', 'Client Industry'), f('clientUrl', 'Client URL'), f('clientLogoUrl', 'Client Logo URL'), f('summary', 'Summary', 'textarea', { wide: true })] },
  { name: 'Publishing', fields: [f('status', 'Status', 'select', { options: ['draft', 'published'] }), f('startedAt', 'Started Date', 'date'), f('completedAt', 'Completed Date', 'date'), f('isFeatured', 'Featured work', 'checkbox'), f('coverImageUrl', 'Cover Image URL', undefined, { wide: true })] },
  { name: 'Categorisation', fields: [f('specialisms', 'Specialisms', undefined, { placeholder: 'Strategy, UX, Engineering', wide: true }), f('platforms', 'Platforms', undefined, { placeholder: 'Web, iOS, Android', wide: true }), f('technologies', 'Technologies', undefined, { placeholder: 'Go, Vue, PostgreSQL', wide: true })] },
  { name: 'Case study', fields: ['challenge', 'approach', 'deliverables', 'outcome'].map(key => f(key, key[0].toUpperCase() + key.slice(1), 'textarea', { wide: true })) },
  { name: 'Content', fields: [f('content', 'Content', 'editor', { wide: true, placeholder: 'Write the work body...' })] },
]
const projectSections: { name: string; fields: Field[] }[] = [
  { name: 'Primary details', fields: [f('name', 'Name'), f('slug', 'Slug'), f('tagline', 'Tagline'), f('projectType', 'Project Type', 'select', { options: ['web-app', 'website', 'automation', 'integration', 'consulting', 'framework'] }), f('description', 'Description', 'editor', { wide: true, placeholder: 'Write the project description...' })] },
  { name: 'Links and media', fields: [f('repositoryUrl', 'Repository URL'), f('liveUrl', 'Live URL'), f('imageUrl', 'Image URL', undefined, { wide: true }), f('technologies', 'Technologies', undefined, { wide: true, placeholder: 'Go, Vue, PostgreSQL' })] },
  { name: 'Publishing', fields: [f('startedAt', 'Started Date', 'date'), f('launchedAt', 'Launched Date', 'date'), f('publishedAt', 'Published Date', 'date'), f('isFeatured', 'Featured project', 'checkbox')] },
]
export function PortfolioForm({ kind, item }: { kind: Kind; item?: Work | Project }) {
  const edit = !!item
  const route = portfolioRoutes[kind]
  const form = useForm<Values>(initial(kind, item))
  const [slugEdited, setSlugEdited] = useState(false)
  const cancel = item ? route.show.path(item.ID) : route.index.path()
  const sections = kind === 'Work' ? workSections : projectSections
  function change(key: string, value: string | boolean) {
    if (key === 'slug' && edit) setSlugEdited(true)
    if (key === (kind === 'Work' ? 'title' : 'name') && (!edit || !slugEdited)) form.setData(data => ({ ...data, [key]: value, slug: slugify(String(value)) }))
    else form.setData(key, value)
  }
  function submit(e: FormEvent) {
    e.preventDefault()
    if (kind === 'Work') form.transform(data => ({ ...data, ...Object.fromEntries(['specialisms', 'platforms', 'technologies'].map(key => [key, String(data[key]).split(',').map(s => s.trim()).filter(Boolean)])) }))
    if (item) form.put(route.update.path(item.ID))
    else form.post(route.create.path())
  }
  return <Frame><Head title={item ? `Edit ${title(item)}` : `New ${kind}`} /><Header heading={`${edit ? 'Edit' : 'New'} ${kind}`} description={item ? title(item) : `Add a new ${kind.toLowerCase()} to your portfolio.`}><Button asChild variant="outline" className={button}><Link href={cancel}>← {edit ? 'Cancel' : 'Back'}</Link></Button></Header>
    <form onSubmit={submit} className="border bg-card"><div className="space-y-8 p-6">{Object.keys(form.errors).length > 0 && <p className="border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">Please fix the highlighted fields.</p>}
      {sections.map((section, index) => {
        const fields = <div className="grid grid-cols-1 gap-6 md:grid-cols-2">{section.fields.map(field => <div key={field.key} className={`space-y-1 ${field.wide ? 'md:col-span-2' : ''}`}>
          {field.type === 'checkbox' ? <Label htmlFor={field.key} className="flex items-center gap-2 text-sm font-normal normal-case tracking-normal"><Checkbox id={field.key} checked={!!form.data[field.key]} onCheckedChange={checked => change(field.key, checked === true)} />{field.label}</Label> : <><Label htmlFor={field.key} className="text-sm font-medium normal-case tracking-normal">{field.label}</Label>
          {field.type === 'select' ? <NativeSelect id={field.key} className="w-full [&_select]:border-input [&_select]:bg-background [&_select]:px-3" value={String(form.data[field.key])} onChange={e => change(field.key, e.target.value)}>{field.options?.map(option => <option key={option} value={option}>{option === 'web-app' ? 'Web app' : option[0].toUpperCase() + option.slice(1)}</option>)}</NativeSelect> : field.type === 'editor' ? <RichEditor id={field.key} value={String(form.data[field.key])} placeholder={field.placeholder} className={kind === 'Work' ? 'min-h-[32rem]' : 'min-h-[24rem]'} onChange={value => change(field.key, value)} /> : field.type === 'textarea' ? <Textarea id={field.key} className={`${input} min-h-24`} value={String(form.data[field.key])} placeholder={field.placeholder} onChange={e => change(field.key, e.target.value)} /> : <Input id={field.key} type={field.type || 'text'} className={input} value={String(form.data[field.key])} placeholder={field.placeholder} onChange={e => change(field.key, e.target.value)} />}</>}
          {form.errors[field.key] && <p className="text-sm text-destructive">{form.errors[field.key]}</p>}
        </div>)}</div>
        return kind === 'Work' ? <details className="border-b py-4" key={section.name} open={index === 0 || index === 4}><summary className="cursor-pointer py-3 font-medium">{section.name}</summary>{fields}</details> : <section className="space-y-4" key={section.name}><h2 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">{section.name}</h2>{fields}</section>
      })}
    </div><div className="flex justify-end gap-2 border-t px-6 py-4"><Button asChild variant="outline" className={button}><Link href={cancel}>Cancel</Link></Button><Button variant="outline" className={button} type="submit" disabled={form.processing}>{form.processing ? (edit ? 'Updating...' : 'Creating...') : `${edit ? 'Update' : 'Create'} ${kind}`}</Button></div></form></Frame>
}

export function PortfolioShow({ kind, item }: { kind: Kind; item: Work | Project }) {
  const route = portfolioRoutes[kind]
  const [confirming, setConfirming] = useState(false)
  const [confirmation, setConfirmation] = useState('')
  const form = useForm({})
  const work = kind === 'Work' ? item as Work : null
  const project = kind === 'Project' ? item as Project : null
  const link = (url: string | null) => url ? <a className="text-primary underline" href={url} target="_blank" rel="noreferrer">{url}</a> : '—'
  const rows: [string, ReactNode][] = work ? [
    ['Client', work.ClientName || '—'], ['Industry', work.ClientIndustry || '—'],
    ['Status', <Badge className="border px-2 py-0.5 text-xs font-semibold normal-case tracking-normal">{work.Status === 'published' ? 'Published' : 'Draft'}</Badge>],
    ['Featured', <Badge className="border px-2 py-0.5 text-xs font-semibold normal-case tracking-normal">{work.IsFeatured ? 'Featured' : 'Not featured'}</Badge>],
    ['Specialisms', work.Specialisms?.join(', ') || '—'], ['Platforms', work.Platforms?.join(', ') || '—'], ['Technologies', work.Technologies?.join(', ') || '—'],
    ['Started', date(work.StartedAt)], ['Completed', date(work.CompletedAt)],
  ] : project ? [
    ['ID', project.ID], ['Name', project.Name], ['Slug', project.Slug], ['Tagline', project.Tagline || '—'],
    ['Description', project.Description ? <MarkdownPreview source={project.Description} /> : '—'],
    ['Project Type', project.ProjectType || '—'], ['Repository URL', link(project.RepositoryUrl)], ['Live URL', link(project.LiveUrl)], ['Image URL', link(project.ImageUrl)],
    ['Technologies', project.Technologies || '—'], ['Started', date(project.StartedAt)], ['Launched', date(project.LaunchedAt)],
    ['Featured', <Badge className="border px-2 py-0.5 text-xs font-semibold normal-case tracking-normal">{project.IsFeatured ? 'Featured' : 'Not featured'}</Badge>],
  ] : []
  rows.push(['Published', item.PublishedAt ? date(item.PublishedAt) : 'Not published'], ['Created', date(item.CreatedAt, true)], ['Updated', date(item.UpdatedAt, true)])
  const image = work?.CoverImageUrl || project?.ImageUrl
  return <Frame><AlertDialog open={confirming} onOpenChange={open => { setConfirming(open); if (!open) setConfirmation('') }}><Head title={title(item)} /><Header heading={title(item)} description={work?.Summary || project?.Tagline}><Button asChild variant="outline" className={button}><Link href={route.index.path()}>← Back</Link></Button><Button asChild variant="outline" className={button}><Link href={route.edit.path(item.ID)}>Edit</Link></Button><AlertDialogTrigger asChild><Button variant="outline" type="button" className={`${button} text-destructive`} disabled={form.processing}>Delete</Button></AlertDialogTrigger></Header>
    <AlertDialogContent className="gap-4 border">
      <AlertDialogTitle className="text-base font-semibold normal-case tracking-normal">Are you sure?</AlertDialogTitle>
      <AlertDialogDescription>This permanently deletes {title(item)}. This action cannot be undone.</AlertDialogDescription>
      <Label htmlFor="confirmation" className="text-sm normal-case tracking-normal">Type DELETE to confirm.</Label>
      <Input id="confirmation" autoComplete="off" className={input} value={confirmation} onChange={e => setConfirmation(e.target.value)} />
      <AlertDialogFooter><AlertDialogCancel className="text-sm font-medium normal-case tracking-normal">Cancel</AlertDialogCancel><Button variant="outline" className={`${button} text-destructive normal-case tracking-normal`} type="button" disabled={confirmation !== 'DELETE' || form.processing} onClick={() => form.delete(route.destroy.path(item.ID))}>Delete {kind.toLowerCase()}</Button></AlertDialogFooter>
    </AlertDialogContent>
    {image && <div className="overflow-hidden border bg-card"><img src={image} alt={title(item)} className="w-full object-cover" /></div>}
    <dl className="divide-y border bg-card">{rows.map(([label, value]) => <div key={label} className="grid grid-cols-3 gap-4 px-6 py-4"><dt className="text-sm font-medium text-muted-foreground">{label}</dt><dd className="col-span-2 text-sm">{value}</dd></div>)}</dl>
    {work?.Content && <div className="border bg-card"><MarkdownPreview source={work.Content} className="p-6" /></div>}
  </AlertDialog></Frame>
}
