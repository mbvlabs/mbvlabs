import { routes } from '@/routes'
import { Input } from "@/Components/ui/input"
import { Textarea } from "@/Components/ui/textarea"
import { NativeSelect } from "@/Components/ui/native-select"
import { Button } from "@/Components/ui/button"
import { Label } from "@/Components/ui/label"
import { Head, Link, useForm } from '@inertiajs/react'
import AdminLayout from '@/Layouts/AdminLayout'
import { button } from './shared'
import type { AdminInquiryData, AdminInquiryPayload } from '@/types/payloads'

export default function Form({ item }: { item?: AdminInquiryData }) {
  const form = useForm<AdminInquiryPayload>({ name: item?.Name ?? '', email: item?.Email ?? '', company: item?.Company ?? '', role: item?.Role ?? '', projectType: item?.ProjectType ?? 'web-app', timeline: item?.Timeline ?? '', message: item?.Message ?? '', source: item?.Source ?? '', status: item?.Status ?? 'new', metadata: item?.Metadata ?? '{}' })
  const back = item ? routes.adminProjectInquiryShow.path(item.ID) : routes.adminProjectInquiryIndex.path()
  const field = (key: 'name' | 'email' | 'company' | 'role' | 'timeline' | 'source', label: string) => <div className="space-y-1"><Label htmlFor={key}>{label}</Label><Input id={key} type={key === 'email' ? 'email' : 'text'} value={form.data[key]} onChange={e => form.setData(key, e.target.value)} />{form.errors[key] && <span className="text-destructive">{form.errors[key]}</span>}</div>
  return <AdminLayout><Head title={item ? `Edit ${item.Name}` : 'New Project Inquiry'} /><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div className="flex justify-between"><div><h1 className="text-2xl font-semibold">{item ? 'Edit Project Inquiry' : 'New Project Inquiry'}</h1><p className="text-muted-foreground">{item ? item.Name : 'Add a new project inquiry.'}</p></div><Link href={back} className={button}>{item ? 'Cancel' : 'Back'}</Link></div>
    <form className="border bg-card" onSubmit={e => { e.preventDefault(); item ? form.put(routes.adminProjectInquiryUpdate.path(item.ID)) : form.post(routes.adminProjectInquiryCreate.path()) }}><div className="space-y-8 p-6">
      <section className="space-y-4"><h2 className="text-sm font-semibold uppercase text-muted-foreground">Contact</h2><div className="grid gap-6 md:grid-cols-2">{field('name', 'Name')}{field('email', 'Email')}{field('company', 'Company')}{field('role', 'Role')}</div></section>
      <section className="space-y-4"><h2 className="text-sm font-semibold uppercase text-muted-foreground">Project</h2><div className="grid gap-6 md:grid-cols-2">
        <div><Label htmlFor="projectType">Project Type</Label><NativeSelect id="projectType" className="w-full" value={form.data.projectType} onChange={e => form.setData('projectType', e.target.value)}>{[['web-app','Web app'],['website','Website'],['automation','Automation'],['integration','Integration'],['consulting','Consulting'],['open-source','Open source'],['modernization','Modernization']].map(([value, label]) => <option key={value} value={value}>{label}</option>)}</NativeSelect></div>
        <div><Label htmlFor="status">Status</Label><NativeSelect id="status" className="w-full" value={form.data.status} onChange={e => form.setData('status', e.target.value)}>{['new','contacted','qualified','proposal','won','lost'].map(s => <option key={s} value={s}>{s[0].toUpperCase() + s.slice(1)}</option>)}</NativeSelect></div>
        {field('timeline', 'Timeline')}<div className="md:col-span-2">{field('source', 'Source')}</div></div>
        <div><Label htmlFor="message">Message</Label><Textarea id="message" className="min-h-32" value={form.data.message} onChange={e => form.setData('message', e.target.value)} />{form.errors.message && <span className="text-destructive">{form.errors.message}</span>}</div>
      </section><div><Label htmlFor="metadata">Metadata</Label><Textarea id="metadata" className="min-h-40 font-mono" value={form.data.metadata} onChange={e => form.setData('metadata', e.target.value)} />{form.errors.metadata && <span className="text-destructive">{form.errors.metadata}</span>}</div>
    </div><div className="flex justify-end gap-2 border-t px-6 py-4"><Link href={back} className={button}>Cancel</Link><Button variant="outline" disabled={form.processing}>{form.processing ? item ? 'Updating...' : 'Creating...' : item ? 'Update Project Inquiry' : 'Create Project Inquiry'}</Button></div></form>
  </div></AdminLayout>
}
