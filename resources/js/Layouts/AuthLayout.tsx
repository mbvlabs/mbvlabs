import type { ReactNode } from 'react'
import FlashMessages from '@/Components/FlashMessages'
import { routes } from '@/routes'

type Props = {
  eyebrow: string
  title: string
  description: string
  metaTitle: string
  metaDescription: string
  formLabel: string
  children: ReactNode
}

export default function AuthLayout({ eyebrow, title, description, metaTitle, metaDescription, formLabel, children }: Props) {
  return <div className="auth-cobalt flex min-h-screen flex-col bg-background font-sans text-foreground">
    <header className="border-b border-foreground">
      <div className="container mx-auto flex items-center justify-between px-6 py-5 md:px-8">
        <a href={routes.homePage.path()} className="font-mono text-xs font-medium uppercase tracking-[0.18em]">MBV Labs</a>
        <p className="font-mono text-[11px] uppercase tracking-[0.12em] opacity-65">Authorized access only</p>
      </div>
    </header>
    <main className="flex flex-1">
      <section className="flex w-full items-center border-b border-foreground">
        <div className="container mx-auto grid gap-14 px-6 py-20 md:px-8 md:py-28 lg:grid-cols-[1.1fr_0.9fr] lg:items-center lg:gap-24">
          <div>
            <p className="font-mono text-xs uppercase tracking-[0.18em]">{eyebrow}</p>
            <h1 className="cobalt-display mt-8 max-w-4xl">{title}</h1>
            <p className="mt-8 max-w-xl text-base leading-7">{description}</p>
            <div className="mt-12 max-w-xl border-t border-foreground pt-5 font-mono text-[11px] uppercase tracking-[0.12em]">
              <p className="font-medium">{metaTitle}</p><p className="mt-1 opacity-65">{metaDescription}</p>
            </div>
          </div>
          <div className="w-full border border-foreground">
            <div className="border-b border-foreground px-6 py-5"><p className="font-mono text-xs uppercase tracking-[0.14em]">{formLabel}</p></div>
            {children}
          </div>
        </div>
      </section>
    </main>
    <FlashMessages />
  </div>
}
