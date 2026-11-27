import { Head, Link, useForm } from '@inertiajs/react'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import AuthLayout from '@/Layouts/AuthLayout'
import { routes } from '@/routes'
import type { PasswordResetRequestPayload } from '@/types/payloads'

export default function ResetPasswordRequest() {
  const form = useForm<PasswordResetRequestPayload>({ email: '' })
  return <AuthLayout eyebrow="Account recovery / Request" title="Reset Password." description="Enter your email address and we'll send you a link to reset your password." metaTitle="Recovery step 01" metaDescription="Email / Reset link" formLabel="Recovery email / 01">
    <Head title="Reset Password" />
    <form onSubmit={event => { event.preventDefault(); form.post(routes.passwordCreate.path()) }}>
      <fieldset className="border-0 p-0" disabled={form.processing}>
        <div className="border-b border-border">
          <Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="email">01 / Email</Label>
          <Input id="email" type="email" autoComplete="email" value={form.data.email} onChange={e => form.setData('email', e.target.value)} className="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60" required />
          {form.errors.email && <p className="px-5 pb-3 text-sm">{form.errors.email}</p>}
        </div>
        <Button type="submit" className="inline-flex min-h-14 w-full items-center justify-center border-0 border-t border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60">{form.processing ? 'Loading' : 'Send Reset Link'}</Button>
      </fieldset>
    </form>
    <p className="border-t border-border px-5 py-5 text-center text-sm">Remember your password? <Link className="font-semibold underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" href={routes.sessionNew.path()}>Login</Link></p>
  </AuthLayout>
}
