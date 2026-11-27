import PublicLayout from '@/Layouts/PublicLayout'

type PublicWork = {
  ID: number; Slug: string; Title: string; Summary: string
  ClientName: { String: string; Valid: boolean }
  CompletedAt: { Time: string; Valid: boolean }
}
type PublicProject = {
  ID: number; Slug: string; Name: string; Tagline: string; ProjectType: string
  RepositoryURL: { String: string; Valid: boolean }
}

const philosophy = [
  ['Use agents to explore better options', 'The advantage is not producing more code faster. It is exploring more approaches, testing more assumptions, and comparing tradeoffs before committing to a path.'],
  ['Stay hypothesis-driven', 'Good software still starts with a clear hypothesis. Agents accelerate the loop, but the work remains grounded in product context, technical judgment, and evidence.'],
  ['Iterate quickly without losing control', 'AI can generate options, prototypes, tests, migrations, refactors, and implementation plans. The discipline is framing the problem, inspecting the output, and turning useful experiments into maintainable software.'],
  ['Build systems teams can own', 'The goal is not a pile of AI-generated code. The goal is finished software your team can understand, operate, and extend after the engagement ends.'],
]
const services = [
  ['Shape the technical plan', 'Turn ambiguous business and product goals into a practical implementation path, with tradeoffs visible before work starts.'],
  ['Build production software', 'Ship web apps, internal tools, automations, integrations, and backend systems using maintainable full-stack engineering.'],
  ['Improve existing systems', 'Audit, modernize, refactor, or extend codebases so your team can move faster without compounding technical debt.'],
]
const faqs = [
  ['What does a fractional tech lead do?', 'I help turn an unclear product or engineering goal into a practical plan, then stay hands-on to build and ship it. You get senior technical judgment without hiring a full-time lead.'],
  ['What kinds of projects are a good fit?', 'Focused, high-leverage work where architecture and execution both matter: production web apps, internal tools, automations, integrations, developer infrastructure, MVPs, and codebase improvements.'],
  ['Do you only advise, or do you write code?', 'Both. I can shape scope and architecture, then implement the software myself or work alongside your team.'],
  ['Can you work with an existing team or codebase?', 'Yes. I can join an existing delivery effort, extend or modernize a codebase, and leave the team with software it can understand and operate.'],
  ['Can one person handle a serious project?', 'The best fit is focused work where senior judgment matters more than headcount. Larger scopes can be phased, delivered with your team, or scoped into the right delivery model.'],
  ['How does an engagement start?', 'Start with a project inquiry. We will clarify the goal, constraints, current system, and the smallest useful first phase before committing to implementation.'],
]
const projectTypes: Record<string, string> = { 'web-app': 'Web app', website: 'Website', automation: 'Automation', integration: 'Integration', consulting: 'Consulting', framework: 'Framework' }
const number = (index: number) => `${String(index + 1).padStart(2, '0')}.`
const rowClass = 'group grid gap-6 border-b border-border py-12 transition-colors hover:bg-secondary/45 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring md:grid-cols-[4rem_minmax(0,0.72fr)_minmax(0,1.35fr)_7rem] md:items-start md:px-3'

export default function Home({ works, projects }: { works: PublicWork[] | null; projects: PublicProject[] | null }) {
  return <PublicLayout>
    <main id="landing-page-fragment-container" className="overflow-hidden bg-background font-sans text-foreground">
      <section className="relative isolate min-h-[680px] border-y border-foreground md:min-h-[calc(100svh-72px)]">
        <div className="container relative mx-auto flex min-h-[680px] flex-col px-6 py-8 md:min-h-[calc(100svh-72px)] md:px-8 md:py-10">
          <div className="cobalt-scanlines pointer-events-none absolute bottom-24 right-8 top-10 z-0 hidden w-[28%] max-w-[24rem] md:block" aria-hidden="true"></div>
          <div className="h-px w-full bg-foreground"></div>
          <div className="flex flex-1 items-center py-16 md:py-20"><div className="relative z-10 w-full md:max-w-[62%] xl:max-w-[52rem]">
            <p className="text-xs font-semibold uppercase tracking-[0.18em]">Field office / Senior product engineering</p>
            <h1 className="cobalt-display mt-7">Fractional tech lead for founders who need software shipped.</h1>
            <p className="mt-8 max-w-xl text-base leading-7 md:text-lg">I help lean teams shape ambiguous product goals, make practical architecture decisions, and build production web apps, internal tools, automations, and developer infrastructure.</p>
            <div className="mt-9 flex flex-wrap gap-3">
              <a href="/contact" data-palantir-event="contact-home-hero" className="inline-flex min-h-11 items-center border border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Start an inquiry</a>
              <a href="/work" className="inline-flex min-h-11 items-center border border-foreground px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] transition-colors hover:bg-primary hover:text-primary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">View work</a>
            </div>
          </div></div>
          <div className="grid gap-5 border-t border-foreground pt-5 text-[11px] uppercase tracking-[0.14em] sm:grid-cols-3">
            <div><p className="font-semibold">Issue 01</p><p className="mt-1 opacity-70">MBV Labs / 2026</p></div>
            <div><p className="font-semibold">Focus</p><p className="mt-1 opacity-70">Plan / Build / Improve</p></div>
            <div><p className="font-semibold">Working model</p><p className="mt-1 opacity-70">Senior judgment / Hands-on delivery</p></div>
          </div>
        </div>
      </section>
      <section className="border-b border-foreground"><div className="container mx-auto px-6 py-16 md:px-8 md:py-20">
        <div className="flex items-end justify-between gap-6"><div><p className="text-xs font-semibold uppercase tracking-[0.18em]">Trusted technical work</p><h2 className="cobalt-section-title mt-4">Selected client ledger.</h2></div><p className="hidden font-mono text-xs tracking-[0.08em] md:block">CLIENTS / 001-005</p></div>
        <div className="mt-10 grid border-y border-foreground md:grid-cols-5">{['YLD', 'everphone', 'Kuatro Group', 'Chatsheet', 'Candy'].map((name, index) => <div key={name} className="flex min-h-36 flex-col justify-between border-b border-border px-5 py-6 last:border-b-0 md:border-b-0 md:border-r md:last:border-r-0"><p className="font-mono text-xs tracking-[0.08em]">{number(index)}</p><p className="font-heading text-2xl font-normal md:text-3xl">{name}</p></div>)}</div>
      </div></section>
      <section className="container mx-auto px-6 py-20 md:px-8 md:py-24">
        <div className="flex flex-col gap-5 border-b border-foreground pb-7 md:flex-row md:items-end md:justify-between"><h2 className="cobalt-section-title max-w-4xl">Selected work, in three entries.</h2><p className="shrink-0 text-xs font-semibold uppercase tracking-[0.18em]">Production systems / Current index</p></div>
        {!works?.length ? <div className="border-b border-border py-10"><p className="text-sm">No featured work is published yet.</p></div> : <div>{works.map((work, index) => <a key={work.ID} href={`/work/${work.Slug}`} className={rowClass}>
          <p className="font-mono text-sm tracking-[0.08em]">{number(index)}</p>
          <div className="text-xs uppercase leading-5 tracking-[0.12em]">{work.ClientName.Valid && work.ClientName.String && <p className="font-semibold">{work.ClientName.String}</p>}{work.CompletedAt.Valid && !work.CompletedAt.Time.startsWith('0001-01-01') && <p className="mt-1 opacity-65">{work.CompletedAt.Time.slice(0, 4)}</p>}</div>
          <div><h3 className="cobalt-row-title">{work.Title}</h3><p className="mt-4 max-w-2xl text-sm leading-6">{work.Summary}</p></div>
          <span className="text-xs font-semibold uppercase tracking-[0.12em] transition-transform group-hover:translate-x-1 md:text-right">Read case →</span>
        </a>)}</div>}
        <div className="flex justify-end pt-8"><a href="/work" className="inline-flex min-h-11 items-center border border-foreground px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] transition-colors hover:bg-primary hover:text-primary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">View the complete work index</a></div>
      </section>
      <section className="relative flex min-h-[72svh] items-center border-y border-foreground"><div className="container mx-auto px-6 py-24 md:px-8 md:py-28">
        <p className="text-xs font-semibold uppercase tracking-[0.18em]">The founder problem / 01</p><blockquote className="cobalt-quote mt-10 max-w-6xl">&ldquo;We need someone senior who can take this from unclear product goal to shipped, maintainable software.&rdquo;</blockquote>
        <div className="mt-14 flex flex-col gap-3 border-t border-foreground pt-5 text-xs uppercase tracking-[0.14em] sm:flex-row sm:items-center sm:justify-between"><p className="font-semibold">The reason MBV Labs exists</p><p className="font-mono opacity-70">AMBIGUITY → DELIVERY → OWNERSHIP</p></div>
      </div></section>
      <section className="container mx-auto px-6 py-20 md:px-8 md:py-24"><div className="grid gap-16 lg:grid-cols-[0.95fr_1.05fr] lg:gap-24">
        <div><p className="text-xs font-semibold uppercase tracking-[0.18em]">Development philosophy</p><h2 className="cobalt-section-title mt-6">Agent-first development for getting to the right answer faster.</h2><p className="mt-7 max-w-lg text-base leading-7">AI is most valuable when it helps a senior engineer explore the problem space faster: more options, faster experiments, clearer tradeoffs, and quicker convergence on software that fits the business.</p></div>
        <div className="grid sm:grid-cols-2">{philosophy.map(([title, body], index) => <div key={title} className="border-t border-border py-10 sm:px-8 sm:odd:border-r sm:odd:pl-0 sm:even:pr-0"><p className="font-mono text-sm tracking-[0.08em]">{number(index)}</p><h3 className="cobalt-row-title mt-6">{title}</h3><p className="mt-4 text-sm leading-6">{body}</p></div>)}</div>
      </div></section>
      <section className="border-y border-foreground"><div className="container mx-auto px-6 py-20 md:px-8 md:py-24">
        <div className="flex flex-col gap-5 border-b border-foreground pb-7 md:flex-row md:items-end md:justify-between"><h2 className="cobalt-section-title max-w-4xl">Technical leadership, measured by shipped outcomes.</h2><p className="shrink-0 text-xs font-semibold uppercase tracking-[0.18em]">Services / Field report</p></div>
        <div className="mt-12 grid gap-14 lg:grid-cols-[1.1fr_0.9fr] lg:items-end lg:gap-20">
          <div><div className="cobalt-pixel-chart" aria-hidden="true">{Array.from({ length: 8 }, (_, index) => <div key={index} className="cobalt-pixel-bar"></div>)}</div><div className="mt-3 flex justify-between font-mono text-[10px] tracking-[0.08em] md:text-xs"><span>UNCLEAR GOAL</span><span>SHIPPED SYSTEM</span></div></div>
          <div>{services.map(([title, body], index) => <div key={title} className="grid grid-cols-[4rem_1fr] gap-5 border-t border-border py-10 first:border-t-0 lg:first:border-t"><p className="font-mono text-sm tracking-[0.08em]">{number(index)}</p><div><h3 className="cobalt-row-title">{title}</h3><p className="mt-4 max-w-xl text-sm leading-6">{body}</p></div></div>)}
            <a href="/services" className="mt-8 inline-flex min-h-11 items-center border border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Explore services</a></div>
        </div>
      </div></section>
      <section className="container mx-auto px-6 py-20 md:px-8 md:py-24">
        <div className="flex flex-col gap-5 border-b border-foreground pb-7 md:flex-row md:items-end md:justify-between"><h2 className="cobalt-section-title">Featured project ledger.</h2><p className="shrink-0 text-xs font-semibold uppercase tracking-[0.18em]">Independent builds / Selected</p></div>
        {!projects?.length ? <div className="border-b border-border py-10"><p className="text-sm">No featured projects are published yet.</p></div> : <div>{projects.map((project, index) => <a key={project.ID} href={`/projects/${project.Slug}`} className={rowClass}>
          <p className="font-mono text-sm tracking-[0.08em]">{number(index)}</p><div><p className="text-xs font-semibold uppercase leading-5 tracking-[0.12em]">{projectTypes[project.ProjectType] ?? project.ProjectType.trim()}</p><p className="mt-1 text-xs uppercase tracking-[0.12em] opacity-65">{project.RepositoryURL.Valid && project.RepositoryURL.String ? 'Open source' : 'Closed-sourced'}</p></div>
          <div><h3 className="cobalt-row-title">{project.Name}</h3><p className="mt-4 max-w-2xl text-sm leading-6">{project.Tagline}</p></div><span className="text-xs font-semibold uppercase tracking-[0.12em] transition-transform group-hover:translate-x-1 md:text-right">View project →</span>
        </a>)}</div>}
        <div className="flex items-center justify-between gap-5 border-b border-foreground py-8 font-mono text-xs tracking-[0.08em]"><span>END OF LANDING REPORT</span><span>MBV / 2026</span></div>
      </section>
      <section className="border-t border-foreground"><div className="container mx-auto px-6 py-20 md:px-8 md:py-24">
        <div className="flex flex-col gap-5 border-b border-foreground pb-7 md:flex-row md:items-end md:justify-between"><h2 className="cobalt-section-title max-w-4xl">Frequently asked questions.</h2><p className="shrink-0 text-xs font-semibold uppercase tracking-[0.18em]">Before an engagement / Common questions</p></div>
        <div className="mt-8">{faqs.map(([question, answer]) => <details key={question} className="border-t border-border last:border-b"><summary className="cursor-pointer py-7 font-heading text-2xl font-normal leading-tight focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:text-3xl">{question}</summary><p className="max-w-2xl pb-7 pr-6 text-sm leading-6">{answer}</p></details>)}</div>
      </div></section>
    </main>
  </PublicLayout>
}
