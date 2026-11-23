import { Head, Link, useForm } from '@inertiajs/react'
import { useState } from 'react'
import AdminLayout from '@/Layouts/AdminLayout'
import { button, Confirm } from '../ProjectInquiry/shared'
import { date, PollToggle, usePolling, type Job } from './shared'
function json(value: string) { try { return JSON.stringify(JSON.parse(value || '{}'), null, 2) } catch { return value } }
export default function Show({ item }: { item: Job }) {
  const [polling, setPolling] = useState(true)
  const form = useForm({})
  usePolling(polling && !form.processing, ['item'])
  const errors = (item.Errors ?? []).map((raw,index) => { let parsed: { attempt?: number; at?: string } = {}; try { parsed = JSON.parse(raw) ?? {} } catch { /* retain raw error */ } return { raw, number: index+1, ...parsed } }).sort((a,b) => a.at && b.at ? new Date(b.at).getTime() - new Date(a.at).getTime() : b.number-a.number)
  const actions = [{ label:'Retry', title:'Retry job', description:`Move job ${item.ID} back to available so River can process it again.`, action:'retry', actionLabel:'Retry' },{ label:'Cancel', title:'Cancel job', description:`Cancel job ${item.ID}. River will treat it as finalized and it will not run.`, action:'cancel', actionLabel:'Cancel job' },{ label:'Discard', title:'Discard job', description:`Discard job ${item.ID}. This finalizes the job without running it.`, action:'discard', actionLabel:'Discard' }]
  const groups = [[['State',item.State],['Queue',item.Queue],['Attempts',`${item.Attempt} / ${item.MaxAttempts}`],['Priority',item.Priority]], [['Created',date(item.CreatedAt)],['Scheduled',date(item.ScheduledAt)],['Attempted',date(item.AttemptedAt)],['Finalized',date(item.FinalizedAt)]]]
  return <AdminLayout><Head title={`Queue Job ${item.ID}`} /><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6"><div className="flex flex-wrap justify-between gap-4"><div><h1 className="text-2xl font-semibold">Job #{item.ID}</h1><p className="text-muted-foreground">{item.Kind}</p></div><PollToggle enabled={polling} setEnabled={setPolling} subject="job" /></div>
    <div className="flex flex-wrap gap-2 border-b pb-4"><Link className={button} href="/admin/queue/jobs">Back</Link>{actions.map(a => <Confirm key={a.action} {...a} disabled={form.processing} onConfirm={() => form.post(`/admin/queue/jobs/${item.ID}/${a.action}`, { preserveScroll: true })} />)}</div>
    <section className="grid border bg-card md:grid-cols-2">{groups.map((group,i) => <dl key={i} className="space-y-4 p-6">{group.map(([label,value]) => <div key={label}><dt className="text-sm font-medium text-muted-foreground">{label}</dt><dd className="mt-1 text-sm">{value}</dd></div>)}</dl>)}</section>
    {(['Args','Metadata'] as const).map(key => <section key={key} className="border bg-card p-6"><h2 className="font-semibold">{key}</h2><pre className="mt-4 max-h-96 overflow-auto whitespace-pre-wrap font-mono text-xs">{json(item[key])}</pre></section>)}
    <section className="border bg-card p-6"><h2 className="font-semibold">Errors</h2>{errors.length ? <div className="mt-4 space-y-3">{errors.map(e => <div key={e.number} className="border bg-muted/30"><div className="flex gap-2 border-b px-3 py-2 text-xs text-muted-foreground"><strong>Error {e.number}</strong>{e.at && <span>{date(e.at)}</span>}{e.attempt != null && <span>Attempt {e.attempt}</span>}</div><pre className="max-h-80 overflow-auto whitespace-pre-wrap p-3 font-mono text-xs">{json(e.raw)}</pre></div>)}</div> : <p className="mt-4 text-sm text-muted-foreground">No recorded errors.</p>}</section>
    <section className="border bg-card p-6"><h2 className="font-semibold">Worker Data</h2><dl className="mt-4 grid gap-4 md:grid-cols-2">{([['Attempted By',item.AttemptedBy],['Tags',item.Tags]] as const).map(([label,values]) => <div key={label}><dt className="text-sm font-medium text-muted-foreground">{label}</dt><dd className="mt-1 text-sm">{values?.length ? values.join(', ') : '-'}</dd></div>)}</dl></section></div></AdminLayout>
}
