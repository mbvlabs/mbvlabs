import { Head, useForm } from '@inertiajs/react'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import AuthLayout from '@/Layouts/AuthLayout'
import { routes } from '@/routes'
import type { PasswordResetPageProps, PasswordResetPayload } from '@/types/payloads'

export default function ResetPassword({ token }: PasswordResetPageProps) {
  const form = useForm<PasswordResetPayload>({ resetPasswordToken: token, password: '', confirmPassword: '' })
  return <AuthLayout eyebrow="Account recovery / New password" title="Reset Your Password." description="Enter your new password below." metaTitle="Recovery step 02" metaDescription="New password / Confirmation" formLabel="New credentials / 01-02">
    <Head title="Reset Password" />
    <form onSubmit={event => { event.preventDefault(); form.put(routes.passwordUpdate.path()) }}>
      <fieldset className="border-0 p-0" disabled={form.processing}>
        {form.errors.resetPasswordToken && <p className="px-5 pb-3 text-sm">{form.errors.resetPasswordToken}</p>}
        <div className="border-b border-border">
          <Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="password">01 / New Password</Label>
          <Input id="password" type="password" autoComplete="new-password" minLength={8} value={form.data.password} onChange={e => form.setData('password', e.target.value)} className="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60" required />
          {form.errors.password && <p className="px-5 pb-3 text-sm">{form.errors.password}</p>}
        </div>
        <div className="border-b border-border">
          <Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="confirmPassword">02 / Confirm New Password</Label>
          <Input id="confirmPassword" type="password" autoComplete="new-password" minLength={8} value={form.data.confirmPassword} onChange={e => form.setData('confirmPassword', e.target.value)} className="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60" required />
          {form.errors.confirmPassword && <p className="px-5 pb-3 text-sm">{form.errors.confirmPassword}</p>}
        </div>
        <Button type="submit" className="inline-flex min-h-14 w-full items-center justify-center border-0 border-t border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60">{form.processing ? 'Loading' : 'Reset Password'}</Button>
      </fieldset>
    </form>
  </AuthLayout>
}
