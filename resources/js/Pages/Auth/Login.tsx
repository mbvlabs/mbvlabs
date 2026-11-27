import { Head, Link, useForm } from '@inertiajs/react'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import AuthLayout from '@/Layouts/AuthLayout'
import { routes } from '@/routes'
import type { LoginPayload } from '@/types/payloads'

export default function Login() {
  const form = useForm<LoginPayload>({ email: '', password: '' })
  return <AuthLayout eyebrow="Account access / Private" title="Login to your account." description="Enter your details below to login to your account." metaTitle="Authorized access only" metaDescription="MBV Labs / Account gateway" formLabel="Credentials / 01-02">
    <Head title="Login" />
    <form onSubmit={event => { event.preventDefault(); form.post(routes.sessionCreate.path()) }}>
      <fieldset className="border-0 p-0" disabled={form.processing}>
        <div className="border-b border-border">
          <Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="email">01 / Email</Label>
          <Input id="email" type="email" autoComplete="email" value={form.data.email} onChange={e => form.setData('email', e.target.value)} className="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60" required />
          {form.errors.email && <p className="px-5 pb-3 text-sm">{form.errors.email}</p>}
        </div>
        <div className="border-b border-border">
          <Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="password">02 / Password</Label>
          <Input id="password" type="password" autoComplete="current-password" value={form.data.password} onChange={e => form.setData('password', e.target.value)} className="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60" required />
          {form.errors.password && <p className="px-5 pb-3 text-sm">{form.errors.password}</p>}
        </div>
        <div className="flex justify-end px-5 py-4"><Link className="text-sm font-semibold underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" href={routes.passwordNew.path()}>Forgot your password?</Link></div>
        <Button type="submit" className="inline-flex min-h-14 w-full items-center justify-center border-0 border-t border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60">{form.processing ? 'Loading' : 'Login'}</Button>
      </fieldset>
    </form>
  </AuthLayout>
}
