import { useState } from 'react'
import { Head, Link } from '@inertiajs/react'
import { Button } from '@/Components/ui/button'
import { Card } from '@/Components/ui/card'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/Components/ui/alert-dialog'
import AdminLayout from '@/Layouts/AdminLayout'

export default function ContentToken({ expiresAt: initialExpiresAt, csrfToken }: { expiresAt: string | null; csrfToken: string }) {
  const [expiresAt, setExpiresAt] = useState(initialExpiresAt)
  const [token, setToken] = useState<string | null>(null)
  const [pending, setPending] = useState<'create' | 'revoke' | null>(null)
  const [confirming, setConfirming] = useState<'create' | 'revoke' | null>(null)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  async function changeToken(action: 'create' | 'revoke') {
    setPending(action)
    setError('')
    setMessage('')
    setToken(null)
    try {
      const response = await fetch('/admin/content-token', {
        method: action === 'create' ? 'POST' : 'DELETE',
        credentials: 'same-origin',
        redirect: 'error',
        headers: { 'X-CSRF-Token': csrfToken, Accept: 'application/json' },
      })
      if (action === 'create') {
        if (response.status !== 201) throw new Error()
        const result: unknown = await response.json()
        if (!result || typeof result !== 'object' || !('token' in result) || typeof result.token !== 'string' || !result.token || !('expiresAt' in result) || typeof result.expiresAt !== 'string' || !result.expiresAt) throw new Error()
        setToken(result.token)
        setExpiresAt(result.expiresAt)
        setMessage('Token created. Copy it now; it cannot be viewed later.')
      } else {
        if (response.status !== 204) throw new Error()
        setExpiresAt(null)
        setMessage('Token revoked.')
      }
    } catch {
      setError(`Could not confirm whether the token was ${action === 'create' ? 'created' : 'revoked'}. Refresh this page to check its status before trying again.`)
    } finally {
      setPending(null)
    }
  }

  async function copyToken() {
    if (!token) return
    try {
      await navigator.clipboard.writeText(token)
      setError('')
      setMessage('Token copied. Store it securely; it cannot be viewed later.')
    } catch {
      setError('Could not copy the token. Select and copy it manually.')
    }
  }

  return <AdminLayout><Head title="Content API token" /><div className="mx-auto w-full max-w-3xl space-y-6">
    <div><Link href="/admin" className="text-sm text-muted-foreground underline">Back to admin</Link><h1 className="mt-3 text-2xl font-semibold tracking-tight">Content API token</h1><p className="text-muted-foreground">Manage access to the content API.</p></div>
    <Card className="gap-4 border p-5">
      <div><h2 className="font-semibold">Token status</h2><p className="mt-1 text-sm text-muted-foreground">{expiresAt ? `Token expires at ${new Date(expiresAt).toLocaleString()}.` : 'No active token.'}</p></div>
      <p className="text-sm text-muted-foreground">Tokens expire after one hour. The plaintext token is shown only when created and cannot be viewed later. Creating another token immediately revokes the previous one; revoking also immediately stops clients using it.</p>
      <div className="flex flex-wrap gap-2"><Button type="button" disabled={pending !== null} onClick={() => expiresAt ? setConfirming('create') : void changeToken('create')}>{pending === 'create' ? 'Creating…' : expiresAt ? 'Create new token' : 'Create token'}</Button><Button type="button" variant="destructive" disabled={pending !== null || !expiresAt} onClick={() => setConfirming('revoke')}>{pending === 'revoke' ? 'Revoking…' : 'Revoke token'}</Button></div>
      <AlertDialog open={confirming !== null} onOpenChange={open => { if (!open) setConfirming(null) }}>
        <AlertDialogContent>
          <AlertDialogHeader><AlertDialogTitle>{confirming === 'create' ? 'Create a new token?' : 'Revoke token?'}</AlertDialogTitle><AlertDialogDescription>{confirming === 'create' ? 'Creating a new token immediately revokes the previous token and stops clients using it.' : 'Revoking this token immediately stops clients using it.'}</AlertDialogDescription></AlertDialogHeader>
          <AlertDialogFooter><AlertDialogCancel>Cancel</AlertDialogCancel><AlertDialogAction variant="destructive" onClick={() => { if (confirming) void changeToken(confirming) }}>{confirming === 'create' ? 'Create new token' : 'Revoke token'}</AlertDialogAction></AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
      {token && <div className="space-y-2"><label htmlFor="content-token" className="text-sm font-medium">New token (copy now)</label><textarea id="content-token" readOnly value={token} rows={3} className="w-full border bg-background p-3 font-mono text-sm" /><Button type="button" variant="outline" disabled={pending !== null} onClick={() => void copyToken()}>Copy token</Button></div>}
      {message && <p role="status" className="text-sm">{message}</p>}
      {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
    </Card>
  </div></AdminLayout>
}
