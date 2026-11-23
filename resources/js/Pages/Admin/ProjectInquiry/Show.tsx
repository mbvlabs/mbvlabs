import { Head, Link, useForm } from '@inertiajs/react'
import AdminLayout from '@/Layouts/AdminLayout'
import MarkdownPreview from '@/Components/MarkdownPreview'
import { button, Confirm, type Inquiry } from './shared'
export default function Show({ item }: { item: Inquiry }) {
  const form = useForm({})
  const domain = item.Email.slice(item.Email.lastIndexOf('@') + 1).toLowerCase()
  const path = `/admin/project-inquiries/${item.ID}`
  const rows = [['Name',item.Name],['Email',item.Email],['Company',item.Company],['Role',item.Role],['Project Type',item.ProjectType],['Timeline',item.Timeline],['Message',item.Message],['Source',item.Source],['Status',item.Status],['Metadata',item.Metadata],['Created',new Date(item.CreatedAt).toLocaleString()],['Updated',new Date(item.UpdatedAt).toLocaleString()]]
  return <AdminLayout><Head title={item.Name} /><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6"><div className="flex flex-wrap items-center justify-between gap-4"><div><h1 className="text-2xl font-semibold">{item.Name}</h1><p className="text-muted-foreground">{item.Email}</p></div><div className="flex flex-wrap gap-2"><Link href="/admin/project-inquiries" className={button}>Back</Link><Link href={`${path}/edit`} className={button}>Edit</Link>
    <Confirm label={item.DomainBanned ? 'Unban domain' : 'Ban domain'} title={`${item.DomainBanned ? 'Unban' : 'Ban'} ${domain}?`} description={`Future project inquiries from ${domain} will be ${item.DomainBanned ? 'accepted again' : 'ignored'}.`} actionLabel={item.DomainBanned ? 'Unban domain' : 'Ban domain'} confirmation={domain} disabled={form.processing} onConfirm={() => item.DomainBanned ? form.delete(`${path}/domain-ban`) : form.post(`${path}/domain-ban`)} />
    <Confirm label="Delete" title="Are you sure?" description={`This permanently deletes the inquiry from ${item.Name}. This action cannot be undone.`} actionLabel="Delete inquiry" confirmation="DELETE" disabled={form.processing} onConfirm={() => form.delete(path)} /></div></div>
    <dl className="divide-y border bg-card">{rows.map(([label,value]) => <div key={label} className="grid grid-cols-3 gap-4 px-6 py-4"><dt className="text-sm font-medium text-muted-foreground">{label}</dt><dd className="col-span-2 text-sm">{label === 'Message' ? <MarkdownPreview source={value} /> : value || '—'}</dd></div>)}</dl></div></AdminLayout>
}
