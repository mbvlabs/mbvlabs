import { routes } from '@/routes'
import { Head, Link, useForm } from '@inertiajs/react'
import { useState } from 'react'
import AdminLayout from '@/Layouts/AdminLayout'
import MarkdownPreview from '@/Components/MarkdownPreview'
import type { AdminDiaryItemProps } from '@/types/payloads'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import { AlertDialog, AlertDialogTrigger, AlertDialogContent, AlertDialogTitle, AlertDialogDescription, AlertDialogCancel, AlertDialogFooter } from '@/Components/ui/alert-dialog'

export default function Show({ item }: AdminDiaryItemProps) {
  const form = useForm({})
  const [confirming, setConfirming] = useState(false)
  const [confirmation, setConfirmation] = useState('')
  const close = () => { setConfirming(false); setConfirmation('') }
  return <AdminLayout><AlertDialog open={confirming} onOpenChange={open => { if (open) setConfirming(true); else close() }}><Head title={`Diary Entry ${item.EntryDate}`} /><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"><div><h1 className="text-2xl font-semibold tracking-tight">{item.EntryDate}</h1><p className="text-muted-foreground">Diary entry</p></div><div className="flex gap-2"><Link href={routes.adminDiaryEntryIndex.path()} className="border px-4 py-2 text-sm">← Back</Link><Link href={routes.adminDiaryEntryEdit.path(item.ID)} className="border px-4 py-2 text-sm">Edit</Link><AlertDialogTrigger asChild><Button variant="destructive" disabled={form.processing} className="bg-destructive px-4 py-2 text-sm font-normal normal-case tracking-normal text-destructive-foreground">Delete</Button></AlertDialogTrigger></div></div>
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">{([['Morning', item.MorningThoughts], ['Evening', item.EveningThoughts]] as const).map(([label, text]) => <section key={label} className="border bg-card p-6"><h2 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">{label}</h2>{text ? <MarkdownPreview source={text} className="mt-4" /> : <p className="mt-4 text-sm text-muted-foreground">—</p>}</section>)}</div>
    <dl className="divide-y border bg-card">{([['Created', item.CreatedAt], ['Updated', item.UpdatedAt]] as const).map(([label, date]) => <div key={label} className="grid grid-cols-3 gap-4 px-6 py-4"><dt className="text-sm font-medium text-muted-foreground">{label}</dt><dd className="col-span-2 text-sm">{new Date(date).toLocaleString()}</dd></div>)}</dl>
  </div><AlertDialogContent className="gap-4 border">
    <AlertDialogTitle className="text-base font-semibold normal-case tracking-normal">Are you sure?</AlertDialogTitle>
    <AlertDialogDescription>This permanently deletes the diary entry for {item.EntryDate}. This action cannot be undone.</AlertDialogDescription>
    <Label htmlFor="confirmation" className="text-sm font-medium normal-case tracking-normal">Type DELETE to confirm</Label>
    <Input id="confirmation" autoFocus autoComplete="off" className="border bg-background px-3 py-2" value={confirmation} onChange={e => setConfirmation(e.target.value)} />
    <AlertDialogFooter><AlertDialogCancel className="text-sm font-normal normal-case tracking-normal">Cancel</AlertDialogCancel><Button variant="destructive" className="bg-destructive px-4 py-2 text-sm font-normal normal-case tracking-normal text-destructive-foreground" disabled={confirmation !== 'DELETE' || form.processing} onClick={() => form.delete(routes.adminDiaryEntryDestroy.path(item.ID))}>Delete diary entry</Button></AlertDialogFooter>
  </AlertDialogContent></AlertDialog></AdminLayout>
}
