import { usePage } from '@inertiajs/react'
import { useState } from 'react'
import { Button } from '@/Components/ui/button'

type FlashMessage = { id: string; type: 'success' | 'error' | 'warning' | 'info'; message: string }

export default function FlashMessages() {
  const { flash } = usePage()
  const [dismissed, setDismissed] = useState<string[]>([])
  const messages = Array.isArray(flash) ? (flash as FlashMessage[]) : []
  return <div className="pointer-events-none fixed bottom-4 right-4 z-50 grid max-w-[calc(100%-2rem)] gap-2" aria-live="polite">
    {messages.filter(({ id }) => !dismissed.includes(id)).map(({ id, type, message }) =>
      <div key={id} role={type === 'error' ? 'alert' : 'status'} className="pointer-events-auto flex items-start gap-4 border bg-background p-4 text-sm text-foreground shadow-md">
        <span>{message}</span>
        <Button type="button" variant="ghost" size="icon-xs" className="-my-1" aria-label="Dismiss notification" onClick={() => setDismissed(ids => [...ids, id])}>×</Button>
      </div>
    )}
  </div>
}
