import { routes } from '@/routes'
import PublicLayout from '@/Layouts/PublicLayout'

export default function About() {
  return <PublicLayout>
    <main id="about-container" className="overflow-hidden bg-background font-sans text-foreground">
      <section className="border-y border-foreground">
        <div className="container mx-auto px-6 py-20 md:px-8 md:py-28">
          <p className="font-mono text-xs uppercase tracking-[0.18em]">Field report / About MBV Labs</p>
          <h1 className="cobalt-display mt-8 max-w-6xl">A fractional tech lead who still ships.</h1>
          <p className="mt-9 max-w-2xl text-base leading-7 md:text-lg">MBV Labs helps founders and lean teams turn ambiguous software goals into production web apps, internal tools, automation, and developer infrastructure without hiring a full-time technical leader too early.</p>
          <div className="mt-16 grid gap-5 border-t border-foreground pt-5 font-mono text-[11px] uppercase tracking-[0.12em] sm:grid-cols-3">
            <div><p className="font-medium">Role</p><p className="mt-1 opacity-65">Fractional tech lead</p></div>
            <div><p className="font-medium">Practice</p><p className="mt-1 opacity-65">Plan / Build / Improve</p></div>
            <div><p className="font-medium">Model</p><p className="mt-1 opacity-65">Senior / Hands-on</p></div>
          </div>
        </div>
      </section>
      <section className="container mx-auto px-6 py-20 md:px-8 md:py-24">
        <div className="flex flex-col gap-5 border-b border-foreground pb-7 md:flex-row md:items-end md:justify-between">
          <h2 className="cobalt-section-title">What I do.</h2>
          <p className="text-xs font-semibold uppercase tracking-[0.18em]">Delivery / Three fields</p>
        </div>
        <div className="grid md:grid-cols-3">
          <div className="border-b border-border py-10 md:border-r md:pr-8">
            <p className="font-mono text-xs tracking-[0.08em]">01.</p><h3 className="cobalt-row-title mt-6">Shape the work</h3>
            <p className="mt-4 text-sm leading-6">Turn fuzzy product or technical goals into scope, architecture, delivery tradeoffs, and a practical build plan.</p>
          </div>
          <div className="border-b border-border py-10 md:border-r md:px-8">
            <p className="font-mono text-xs tracking-[0.08em]">02.</p><h3 className="cobalt-row-title mt-6">Build the system</h3>
            <p className="mt-4 text-sm leading-6">Ship full-stack features, backend systems, internal tools, integrations, automations, and developer infrastructure.</p>
          </div>
          <div className="border-b border-border py-10 md:pl-8">
            <p className="font-mono text-xs tracking-[0.08em]">03.</p><h3 className="cobalt-row-title mt-6">Leave it maintainable</h3>
            <p className="mt-4 text-sm leading-6">Make choices that fit the team's future operating reality, not just the fastest path to a demo.</p>
          </div>
        </div>
      </section>
      <section className="flex min-h-[60svh] items-center border-y border-foreground">
        <div className="container mx-auto px-6 py-24 md:px-8 md:py-28">
          <p className="font-mono text-xs uppercase tracking-[0.18em]">Founder brief / 01</p>
          <blockquote className="cobalt-quote mt-10 max-w-6xl">&ldquo;I need technical help, but I do not want to hire full-time yet.&rdquo;</blockquote>
          <div className="mt-14 flex items-center justify-between gap-5 border-t border-foreground pt-5 font-mono text-xs uppercase tracking-[0.12em]"><p className="font-medium">MBV Labs</p><p className="opacity-65">Senior help / Flexible model</p></div>
        </div>
      </section>
      <section className="container mx-auto px-6 py-20 md:px-8 md:py-24">
        <div className="grid gap-12 lg:grid-cols-[0.65fr_1.35fr] lg:gap-24">
          <div><p className="font-mono text-xs uppercase tracking-[0.18em]">Audience / 01</p><h2 className="cobalt-section-title mt-5">Who it is for.</h2></div>
          <div className="border-t border-foreground pt-7">
            <h3 className="cobalt-row-title max-w-3xl">Founders and lean teams with important software to ship.</h3>
            <p className="mt-6 max-w-3xl text-base leading-7">The fit is focused, high-leverage engineering work where senior judgment matters more than large headcount: MVPs, product buildouts, internal operations tools, integrations, codebase modernization, and technical experiments.</p>
            <div className="mt-10 grid border-y border-border md:grid-cols-2">
              <div className="border-b border-border py-7 md:border-b-0 md:border-r md:pr-8"><p className="font-mono text-xs uppercase tracking-[0.12em]">When teams call</p><p className="mt-4 text-sm leading-6">Progress is too slow, the existing team lacks bandwidth, a system is getting harder to change, or a product bottleneck needs software to remove it.</p></div>
              <div className="py-7 md:pl-8"><p className="font-mono text-xs uppercase tracking-[0.12em]">What they get</p><p className="mt-4 text-sm leading-6">Direct access to senior engineering judgment, hands-on implementation, clear communication, and software that can keep moving after launch.</p></div>
            </div>
          </div>
        </div>
      </section>
      <section className="border-y border-foreground">
        <div className="container mx-auto px-6 py-20 md:px-8 md:py-24">
          <p className="font-mono text-xs uppercase tracking-[0.18em]">Working style / Direct ownership</p>
          <h2 className="cobalt-section-title mt-5 max-w-5xl">Pragmatic software, clear tradeoffs, direct senior ownership.</h2>
          <p className="mt-7 max-w-3xl text-base leading-7">MBV Labs is not positioned as cheap implementation labor or a bloated agency process. It is a focused technical partner for teams that need someone who can understand the business context, make practical technical decisions, communicate clearly, and personally ship durable software.</p>
          <a href={routes.projectInquiryIndex.path()} data-palantir-event="contact-about" className="mt-9 inline-flex min-h-11 items-center border border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Start an inquiry</a>
        </div>
      </section>
    </main>
  </PublicLayout>
}
