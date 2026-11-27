import { Head, usePage } from '@inertiajs/react'
import { useEffect, useRef, useState, type ReactNode } from 'react'
import { Button } from '@/Components/ui/button'

type SEO = {
  title: string
  description: string
  canonical: string
  image: string
  imageAlt: string
  metaType: string
  extraMeta?: { name?: string; property?: string; content: string }[]
  structuredData?: Record<string, unknown>
}

type SharedProps = {
  seo: SEO
  publicSite: { production: boolean; year: number }
  auth?: { isAuthenticated: boolean }
}

type Flash = { id: string; type: string; message: string }

const navigation = [
  ['Services', '/services'],
  ['Work', '/work'],
  ['Projects', '/projects'],
  ['Blog', '/blog'],
  ['About', '/about'],
] as const

function PublicHead({ seo, production }: { seo: SEO; production: boolean }) {
  const favicon = 'https://media.mbvlabs.com/favicons'
  return (
    <Head title={seo.title}>
      <link head-key="rss" rel="alternate" type="application/rss+xml" title="MBV Labs Blog" href="/rss.xml" />
      <meta head-key="description" name="description" content={seo.description} />
      <link head-key="canonical" rel="canonical" href={seo.canonical} />
      <meta head-key="og:type" property="og:type" content={seo.metaType} />
      <meta head-key="og:title" property="og:title" content={seo.title} />
      <meta head-key="og:description" property="og:description" content={seo.description} />
      <meta head-key="og:url" property="og:url" content={seo.canonical} />
      <meta head-key="og:site_name" property="og:site_name" content="MBV Labs" />
      <meta head-key="og:locale" property="og:locale" content="en_US" />
      <meta head-key="twitter:card" name="twitter:card" content="summary_large_image" />
      <meta head-key="twitter:creator" name="twitter:creator" content="@mbvisti" />
      <meta head-key="twitter:title" name="twitter:title" content={seo.title} />
      <meta head-key="twitter:description" name="twitter:description" content={seo.description} />
      {seo.image && <meta head-key="twitter:image" name="twitter:image" content={seo.image} />}
      {seo.image && <meta head-key="og:image" property="og:image" content={seo.image} />}
      {seo.image && seo.imageAlt && <meta head-key="og:image:alt" property="og:image:alt" content={seo.imageAlt} />}
      {(seo.extraMeta ?? []).filter(meta => meta.content && (meta.name || meta.property)).map((meta, index) =>
        <meta key={index} head-key={`extra-meta-${index}`} {...(meta.name ? { name: meta.name } : {})} {...(meta.property ? { property: meta.property } : {})} content={meta.content} />
      )}
      <link head-key="favicon-ico" rel="icon" type="image/x-icon" href={`${favicon}/favicon.ico`} />
      <link head-key="favicon-16" rel="icon" type="image/png" sizes="16x16" href={`${favicon}/favicon-16x16.png`} />
      <link head-key="favicon-32" rel="icon" type="image/png" sizes="32x32" href={`${favicon}/favicon-32x32.png`} />
      <link head-key="favicon-base" rel="icon" type="image/png" href={`${favicon}/mbvlabs-favi-base.png`} />
      <link head-key="apple-icon" rel="apple-touch-icon" sizes="180x180" href={`${favicon}/apple-touch-icon.png`} />
      <link head-key="android-192" rel="icon" type="image/png" sizes="192x192" href={`${favicon}/android-chrome-192x192.png`} />
      <link head-key="android-256" rel="icon" type="image/png" sizes="256x256" href={`${favicon}/android-chrome-256x256.png`} />
      <link head-key="mask-icon" rel="mask-icon" href={`${favicon}/safari-pinned-tab.svg`} color="#18181b" />
      <link head-key="manifest" rel="manifest" href={`${favicon}/site.webmanifest`} />
      {seo.structuredData && <script head-key="schema-org" id="schema-org" type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(seo.structuredData).replace(/</g, '\\u003c') }} />}
      {production && <script head-key="analytics" defer src="https://analytics.mbvlabs.com/t/script.js" data-website-id="d6e66aea-49c8-4db8-964e-d0fbc98aed85" />}
    </Head>
  )
}

function SiteHeader({ authenticated }: { authenticated: boolean }) {
  const [open, setOpen] = useState(false)
  const opener = useRef<HTMLButtonElement>(null)
  const closer = useRef<HTMLButtonElement>(null)
  function closeMenu() {
    setOpen(false)
    opener.current?.focus()
  }
  useEffect(() => {
    if (!open) return
    closer.current?.focus()
    const onKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') closeMenu()
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [open])
  const mobileNavigation = [['Home', '/'], navigation[1], navigation[2], navigation[0], navigation[4], navigation[3], ['Contact', '/contact'], ...(authenticated ? [['Admin', '/admin']] : [])]
  return (
    <header className="sticky top-0 z-50 border-b border-border bg-background/90 backdrop-blur">
      <div className="container mx-auto flex items-center justify-between px-6 py-5 md:px-8">
        <a href="/" className="inline-flex items-center gap-2 text-sm font-bold tracking-tight text-foreground" aria-label="mbvlabs home">
          <span className="inline-flex h-4 w-5 items-center justify-center"><svg className="h-4 w-5" viewBox="0 0 24 18" fill="none" aria-hidden="true"><path d="M2 12.5L8 6.5L13 11.5L22 2.5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" /><path d="M2 17H8.5L14 11.5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" opacity=".55" /></svg></span>
          <span>mbvlabs.</span>
        </a>
        <nav className="flex items-center gap-5" aria-label="Primary navigation">
          {navigation.map(([label, href]) => <a key={href} href={href} className="hidden text-[11px] font-semibold uppercase tracking-[0.14em] text-muted-foreground transition-colors hover:text-foreground md:inline-flex">{label}</a>)}
          <a href="/contact" data-palantir-event="contact-header" className="border border-primary bg-primary px-4 py-2 text-[11px] font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary">Contact</a>
          {authenticated && <a href="/admin" className="hidden text-[11px] font-semibold uppercase tracking-[0.14em] md:inline-flex">Admin</a>}
          <Button variant="ghost" size="icon-sm" ref={opener} className="inline-flex h-9 w-9 items-center justify-center text-foreground transition-opacity hover:bg-transparent hover:text-inherit hover:opacity-60 aria-expanded:bg-transparent aria-expanded:text-inherit md:hidden" type="button" aria-label="Open menu" aria-controls="mobile-navigation" aria-expanded={open} onClick={() => setOpen(true)}><svg className="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M5 8H19M5 16H19" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg></Button>
        </nav>
      </div>
      <div id="mobile-navigation" className="fixed inset-0 z-[60] min-h-screen overflow-y-auto bg-primary px-6 py-7 text-primary-foreground md:hidden" hidden={!open}>
        <div className="flex items-center justify-between">
          <a href="/" className="inline-flex items-center gap-2 text-sm font-bold" aria-label="mbvlabs home">mbvlabs.</a>
          <Button variant="ghost" size="icon-sm" ref={closer} className="inline-flex h-9 w-9 items-center justify-center transition-opacity hover:bg-transparent hover:text-inherit hover:opacity-60 aria-expanded:bg-transparent aria-expanded:text-inherit [&_svg]:size-4" type="button" aria-label="Close menu" onClick={closeMenu}><svg className="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M6 6L18 18M18 6L6 18" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg></Button>
        </div>
        <nav className="mt-16" aria-label="Mobile navigation"><ul>
          {mobileNavigation.map(([label, href]) => <li key={href}><a href={href} data-palantir-event={label === 'Contact' ? 'contact-header' : undefined} className={`block border-t border-primary-foreground/20 py-5 font-heading text-4xl${label === 'Contact' ? ' border-b' : ''}`}>{label}</a></li>)}
        </ul></nav>
      </div>
    </header>
  )
}

function SiteFooter({ year }: { year: number }) {
  return <footer id="contact" className="border-t border-foreground bg-background text-foreground">
    <div className="container mx-auto px-6 py-16 md:px-8 md:py-24">
      <div className="grid gap-10 border-b border-foreground pb-16 lg:grid-cols-[1.4fr_0.6fr] lg:items-end">
        <div><p className="text-xs font-semibold uppercase tracking-[0.18em]">Next engagement / Open</p><h2 className="mt-7 max-w-4xl font-heading text-4xl font-normal leading-[0.98] md:text-5xl lg:text-6xl">Tell me about the software you need to ship.</h2></div>
        <div><p className="max-w-md text-sm leading-6">Senior engineering judgment, hands-on implementation, and a clear path from product ambiguity to maintainable production software.</p><a href="/contact" data-palantir-event="contact-footer" className="mt-8 inline-flex min-h-11 items-center border border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary">Start an inquiry</a></div>
      </div>
      <div className="grid gap-10 py-12 text-xs sm:grid-cols-2 lg:grid-cols-5">
        <div><p className="font-semibold uppercase tracking-[0.14em]">Work</p><ul className="mt-4 space-y-3 opacity-70"><li><a href="/work">Client work</a></li><li><a href="/projects">Projects</a></li></ul></div>
        <div><p className="font-semibold uppercase tracking-[0.14em]">Company</p><ul className="mt-4 space-y-3 opacity-70"><li><a href="/services">Services</a></li><li><a href="/about">About</a></li></ul></div>
        <div><p className="font-semibold uppercase tracking-[0.14em]">Connect</p><ul className="mt-4 space-y-3 opacity-70"><li><a href="/blog">Blog</a></li><li><a href="/rss.xml" data-palantir-event="rss-click">RSS</a></li><li><a href="mailto:hello@mbvlabs.com" data-palantir-event="contact-email">Email</a></li></ul></div>
        <div><p className="font-semibold uppercase tracking-[0.14em]">Social</p><ul className="mt-4 space-y-3 opacity-70"><li><a href="https://www.youtube.com/@mbvlabs" target="_blank" rel="noreferrer">YouTube</a></li><li><a href="https://x.com/mbvlabs" target="_blank" rel="noreferrer">X</a></li><li><a href="https://www.linkedin.com/in/mortenvistisen" target="_blank" rel="noreferrer">LinkedIn</a></li><li><a href="https://github.com/mbvlabs" target="_blank" rel="noreferrer">GitHub</a></li></ul></div>
        <div><p className="font-semibold uppercase tracking-[0.14em]">Field notes</p><p className="mt-4 max-w-sm leading-5 opacity-70">Pragmatic engineering decisions, shipped systems, and the tradeoffs behind founder-led products.</p></div>
      </div>
      <div className="flex flex-col gap-4 border-t border-border pt-6 font-mono text-[11px] tracking-[0.08em] sm:flex-row sm:items-center sm:justify-between"><a href="/" className="font-semibold" aria-label="mbvlabs home">mbvlabs.</a><p>&copy; {year} / FIELD OFFICE</p></div>
    </div>
  </footer>
}

export default function PublicLayout({ children }: { children: ReactNode }) {
  const { props, flash } = usePage<SharedProps>()
  const flashes = Array.isArray(flash) ? flash as Flash[] : []
  return <div className="cobalt-grid flex min-h-screen flex-col bg-background text-foreground antialiased">
    <PublicHead seo={props.seo} production={props.publicSite.production} />
    <SiteHeader authenticated={props.auth?.isAuthenticated ?? false} />
    <div className="flex-1">{children}</div>
    <SiteFooter year={props.publicSite.year} />
    <div id="flashContainer" className="fixed bottom-4 right-4 z-50 flex flex-col gap-2" aria-live="polite">
      {flashes.map((item, index) => <div key={item.id ?? index} className={`max-w-sm border px-4 py-3 text-sm font-semibold shadow-lg ${item.type === 'success' ? 'border-green-600 bg-green-600 text-white' : item.type === 'error' ? 'border-red-600 bg-red-600 text-white' : 'border-zinc-950 bg-zinc-950 text-white'}`}>{item.message}</div>)}
    </div>
  </div>
}
