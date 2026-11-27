import { routes } from '@/routes'
import { Head, Link, useForm } from '@inertiajs/react'
import { useState } from 'react'
import AdminLayout from '@/Layouts/AdminLayout'
import MarkdownPreview from '@/Components/MarkdownPreview'
import type { AdminBlogItemProps } from '@/types/payloads'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import { AlertDialog, AlertDialogTrigger, AlertDialogContent, AlertDialogTitle, AlertDialogDescription, AlertDialogCancel, AlertDialogFooter } from '@/Components/ui/alert-dialog'

export default function Show({ item }: AdminBlogItemProps) {
  const form = useForm({})
  const [confirming, setConfirming] = useState(false)
  const [confirmation, setConfirmation] = useState('')
  const close = () => { setConfirming(false); setConfirmation('') }
  const details = [['Slug', item.Slug], ['Status', item.Status === 'draft' && item.PublicationSchedule ? 'Scheduled' : item.Status], ['Tags', item.Tags || '—'], ...(item.PublicationSchedule ? [['Scheduled for', new Date(item.PublicationSchedule.scheduled_at).toLocaleString()]] : []), ['Published', item.PublishedAt ? new Date(item.PublishedAt).toLocaleDateString() : 'Not published'], ['Created', new Date(item.CreatedAt).toLocaleString()], ['Updated', new Date(item.UpdatedAt).toLocaleString()]]
  return <AdminLayout><AlertDialog open={confirming} onOpenChange={open => { if (open) setConfirming(true); else close() }}><Head title={item.Title} /><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div className="flex flex-wrap items-center justify-between gap-4"><div><h1 className="text-2xl font-semibold tracking-tight">{item.Title}</h1><p className="text-muted-foreground">{item.Excerpt}</p></div><div className="flex gap-2"><Link className="border px-4 py-2 text-sm" href={routes.adminBlogPostIndex.path()}>← Back</Link><Link className="border px-4 py-2 text-sm" href={routes.adminBlogPostEdit.path(item.ID)}>Edit</Link><AlertDialogTrigger asChild><Button variant="destructive" className="bg-destructive px-4 py-2 text-sm font-normal normal-case tracking-normal text-destructive-foreground" disabled={form.processing}>Delete</Button></AlertDialogTrigger></div></div>
    {item.CoverImageUrl && <div className="overflow-hidden border bg-card"><img src={item.CoverImageUrl} alt={item.Title} className="w-full object-cover" /></div>}
    <dl className="divide-y border bg-card">{details.map(([label, value]) => <div key={label} className="grid grid-cols-3 gap-4 px-6 py-4"><dt className="text-sm font-medium text-muted-foreground">{label}</dt><dd className="col-span-2 text-sm">{value}</dd></div>)}</dl>
    {item.Body && <div className="border bg-card"><MarkdownPreview source={item.Body} className="p-6" /></div>}
  </div><AlertDialogContent className="gap-4 border">
    <AlertDialogTitle className="text-base font-semibold normal-case tracking-normal">Are you sure?</AlertDialogTitle>
    <AlertDialogDescription>This permanently deletes {item.Title}. This action cannot be undone.</AlertDialogDescription>
    <Label htmlFor="confirmation" className="text-sm font-medium normal-case tracking-normal">Type DELETE to confirm</Label>
    <Input id="confirmation" autoFocus autoComplete="off" className="border bg-background px-3 py-2" value={confirmation} onChange={e => setConfirmation(e.target.value)} />
    <AlertDialogFooter><AlertDialogCancel className="text-sm font-normal normal-case tracking-normal">Cancel</AlertDialogCancel><Button variant="destructive" className="bg-destructive px-4 py-2 text-sm font-normal normal-case tracking-normal text-destructive-foreground" disabled={confirmation !== 'DELETE' || form.processing} onClick={() => form.delete(routes.adminBlogPostDestroy.path(item.ID))}>Delete blog post</Button></AlertDialogFooter>
  </AlertDialogContent></AlertDialog></AdminLayout>
}
