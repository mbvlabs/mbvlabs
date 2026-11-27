import { Button } from "@/Components/ui/button"
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/Components/ui/table"
import { Link, router } from '@inertiajs/react'
import { useEffect } from 'react'
import { button } from '../ProjectInquiry/shared'

export type Job = { ID: number; State: string; Kind: string; Queue: string; Attempt: number; MaxAttempts: number; CreatedAt: string; ScheduledAt: string; AttemptedAt: string | null; FinalizedAt: string | null; Priority: number; Args: string; Metadata: string; AttemptedBy: string[] | null; Errors: string[] | null; Tags: string[] | null }
export type Queue = { Name: string; IsPaused: boolean; AvailableCount: number; PendingCount: number; ScheduledCount: number; RunningCount: number; RetryableCount: number; TotalCount: number; ActiveClients: number; MaxWorkers: number; UpdatedAt: string | null }
export const date = (value: string | null) => value ? new Date(value).toLocaleString() : '-'
export function usePolling(enabled: boolean, only: string[]) {
  const keys = only.join(',')
  useEffect(() => {
    if (!enabled) return
    const timer = window.setInterval(() => router.reload({ only: keys.split(',') }), 5000)
    return () => window.clearInterval(timer)
  }, [enabled, keys])
}
export function PollToggle({ enabled, setEnabled, subject }: { enabled: boolean; setEnabled: (value: boolean) => void; subject: string }) {
  return <div className="flex items-center gap-2"><span className="text-xs font-medium uppercase text-muted-foreground">{enabled ? '◌ Polling' : '◌ Paused'}</span><Button type="button" variant="outline" aria-label={`${enabled ? 'Pause' : 'Resume'} ${subject} polling`} aria-pressed={enabled} onClick={() => setEnabled(!enabled)}>{enabled ? 'Pause polling' : 'Resume polling'}</Button></div>
}
export function JobTable({ items, created = false }: { items: Job[]; created?: boolean }) {
  const headers = ['ID','Kind','Queue','State','Attempts', ...(created ? ['Created'] : []), 'Scheduled']
  return <div className="overflow-x-auto"><Table className="w-full text-left text-sm"><TableHeader><TableRow className="border-b">{headers.map(h => <TableHead key={h} className="p-4 font-medium text-muted-foreground">{h}</TableHead>)}</TableRow></TableHeader><TableBody>{items.map(j => <TableRow key={j.ID} className="cursor-pointer border-b hover:bg-muted/40" onClick={e => { if (!(e.target as Element).closest('a')) router.visit(`/admin/queue/jobs/${j.ID}`) }}><TableCell className="p-4 font-medium"><Link aria-label={`View queue job ${j.ID}`} href={`/admin/queue/jobs/${j.ID}`}>#{j.ID}</Link></TableCell><TableCell className="p-4">{j.Kind}</TableCell><TableCell className="p-4">{j.Queue}</TableCell><TableCell className="p-4">{j.State}</TableCell><TableCell className="p-4">{j.Attempt} / {j.MaxAttempts}</TableCell>{created && <TableCell className="p-4">{date(j.CreatedAt)}</TableCell>}<TableCell className="p-4">{date(j.ScheduledAt)}</TableCell></TableRow>)}{!items.length && <TableRow><TableCell colSpan={headers.length} className="h-24 text-center text-muted-foreground">No jobs found.</TableCell></TableRow>}</TableBody></Table></div>
}
