import { routes } from '@/routes'
import { Input } from "@/Components/ui/input"
import { NativeSelect } from "@/Components/ui/native-select"
import { Button } from "@/Components/ui/button"
import { Head, Link, router } from '@inertiajs/react'
import { useState } from 'react'
import AdminLayout from '@/Layouts/AdminLayout'
import { button } from '../ProjectInquiry/shared'
import { JobTable, usePolling } from './shared'
import type { AdminQueueJobsProps } from '@/types/payloads'
export default function Jobs({ items, queues: queueRows, filters, pagination }: AdminQueueJobsProps) {
  const queues = queueRows ?? []
  const [values, setValues] = useState({ state: filters.State, queue: filters.Queue, kind: filters.Kind })
  usePolling(true, ['items','queues','pagination'])
  const apply = (v: typeof values) => router.get(routes.adminQueueJobs.path(), v, { preserveScroll: true, preserveState: true })
  const pageHref = (page: number) => { const p = new URLSearchParams(); Object.entries(values).forEach(([k,v]) => { if (v) p.set(k,v) }); p.set('page',String(page)); p.set('per_page',String(pagination.PageSize)); return `${routes.adminQueueJobs.path()}?${p}` }
  return <AdminLayout><Head title="Queue Jobs" /><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6"><div><h1 className="text-2xl font-semibold">Queue Jobs</h1><p className="text-muted-foreground">Inspect and manage River jobs.</p></div>
    <form className="grid gap-3 border bg-card p-4 md:grid-cols-[1fr_1fr_2fr_auto_auto]" onSubmit={e => { e.preventDefault(); apply(values) }}><NativeSelect aria-label="State" className="w-full" value={values.state} onChange={e => setValues({ ...values, state: e.target.value })}><option value="">All states</option>{['available','cancelled','completed','discarded','pending','retryable','running','scheduled'].map(s => <option key={s}>{s}</option>)}</NativeSelect><NativeSelect aria-label="Queue" className="w-full" value={values.queue} onChange={e => setValues({ ...values, queue: e.target.value })}><option value="">All queues</option>{queues.map(q => <option key={q.Name}>{q.Name}</option>)}</NativeSelect><Input className="w-full" placeholder="Filter by kind" value={values.kind} onChange={e => setValues({ ...values, kind: e.target.value })} /><Button variant="outline">Filter</Button><Button type="button" variant="outline" onClick={() => { const empty = { state:'', queue:'', kind:'' }; setValues(empty); apply(empty) }}>Clear</Button></form>
    <div className="border bg-card"><JobTable items={items} created /></div><div className="flex flex-wrap items-center justify-between gap-3"><p className="text-sm text-muted-foreground">Showing page {pagination.Page} of {Math.max(pagination.TotalPages,1)} for {pagination.TotalCount} jobs.</p><div className="flex gap-2">{pagination.Page > 1 ? <Link className={button} href={pageHref(pagination.Page-1)}>Previous</Link> : <Button disabled variant="outline">Previous</Button>}{pagination.Page < pagination.TotalPages ? <Link className={button} href={pageHref(pagination.Page+1)}>Next</Link> : <Button disabled variant="outline">Next</Button>}</div></div></div></AdminLayout>
}
