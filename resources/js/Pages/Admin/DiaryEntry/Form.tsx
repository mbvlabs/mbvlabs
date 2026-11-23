import { Head, Link, useForm } from '@inertiajs/react'
import { useEffect, useState } from 'react'
import AdminLayout from '@/Layouts/AdminLayout'
import RichEditor from '@/Components/RichEditor'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'

export type DiaryEntry = { ID: number; EntryDate: string; MorningThoughts: string; EveningThoughts: string; CreatedAt: string; UpdatedAt: string }
type Props = { item?: DiaryEntry; today?: string; focus?: 'morning' | 'evening' | ''; create?: boolean }
export default function Form({ item, today, focus, create = false }: Props) {
  const editing = !!item && !create
  const [period, setPeriod] = useState<'morning' | 'evening'>(focus === 'evening' ? 'evening' : 'morning')
  useEffect(() => { if (focus === 'morning' || focus === 'evening') setPeriod(focus) }, [focus])
  const form = useForm({ entryDate: item?.EntryDate ?? today ?? '', morningThoughts: item?.MorningThoughts ?? '', eveningThoughts: item?.EveningThoughts ?? '' })
  const field = period === 'morning' ? 'morningThoughts' : 'eveningThoughts'
  const cancel = editing ? `/admin/diary/${item!.ID}` : '/admin/diary'
  const outline = 'inline-flex items-center border bg-background px-4 py-2 text-sm font-medium hover:bg-muted'
  return <AdminLayout><Head title={editing ? `Edit Diary Entry ${item!.EntryDate}` : 'New Diary Entry'} /><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"><div><h1 className="text-2xl font-semibold tracking-tight">{editing ? 'Edit Diary Entry' : 'New Diary Entry'}</h1><p className="text-muted-foreground">{editing ? item!.EntryDate : 'Capture morning and evening notes for one day.'}</p></div><Button asChild variant="outline" className={`${outline} normal-case tracking-normal`}><Link href={cancel}>← {editing ? 'Cancel' : 'Back'}</Link></Button></div>
    <form className="border bg-card" onSubmit={e => { e.preventDefault(); if (editing) form.put(`/admin/diary/${item!.ID}`); else form.post('/admin/diary') }}><div className="space-y-8 p-6">
      {Object.keys(form.errors).length > 0 && <div className="border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">Please fix the highlighted fields.</div>}
      <section className="space-y-4"><div className="max-w-sm space-y-1"><Label htmlFor="entryDate" className="text-sm font-medium normal-case tracking-normal">Entry Date</Label><Input id="entryDate" type="date" className="flex h-9 w-full border border-input bg-background px-3 py-1 text-sm" value={form.data.entryDate} onChange={e => form.setData('entryDate', e.target.value)} />{form.errors.entryDate && <p className="text-sm text-destructive">{form.errors.entryDate}</p>}</div></section>
      <section className="space-y-3"><div className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"><Label htmlFor={field} className="text-sm font-medium normal-case tracking-normal">{period === 'morning' ? 'Morning Thoughts' : 'Evening Thoughts'}</Label><div className="inline-flex w-fit border bg-background" role="group" aria-label="Thought period">{(['morning', 'evening'] as const).map(p => <Button variant="ghost" key={p} type="button" aria-label={`Show ${p} thoughts editor`} aria-pressed={period === p} onClick={() => setPeriod(p)} className={`min-w-28 border px-3 py-1 text-sm font-normal normal-case tracking-normal ${period === p ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:bg-muted'}`}>{p === 'morning' ? 'Morning' : 'Evening'}</Button>)}</div></div>
        <RichEditor key={period} id={field} value={form.data[field]} onChange={value => form.setData(field, value)} placeholder={`Write the ${period} thoughts...`} className="h-[30rem]" />
        {form.errors[field] && <p className="text-sm text-destructive">{form.errors[field]}</p>}
      </section>
    </div><div className="flex justify-end gap-2 border-t px-6 py-4"><Button asChild variant="outline" className={`${outline} normal-case tracking-normal`}><Link href={cancel}>Cancel</Link></Button><Button type="submit" disabled={form.processing} className="bg-primary px-4 py-2 text-sm font-medium normal-case tracking-normal text-primary-foreground disabled:opacity-50">{form.processing ? editing ? 'Saving...' : 'Creating...' : editing ? 'Save Diary Entry' : 'Create Diary Entry'}</Button></div></form>
  </div></AdminLayout>
}
