export default function WorkExperience() {
  return (
	<section className="border-t border-foreground">
		<div className="container mx-auto px-6 py-20 md:px-8 md:py-24">
			<div className="flex flex-col gap-5 border-b border-foreground pb-7 md:flex-row md:items-end md:justify-between">
				<div>
					<p className="font-mono text-xs uppercase tracking-[0.18em]">Employment archive</p>
					<h2 className="cobalt-section-title mt-4">Experience history.</h2>
				</div>
				<p className="text-xs font-semibold uppercase tracking-[0.18em]">Six records / 2016-2022</p>
			</div>
			<div>
				<article className="grid gap-6 border-b border-border py-12 md:grid-cols-[4rem_minmax(11rem,0.7fr)_minmax(0,1.3fr)] md:px-3">
					<p className="font-mono text-sm tracking-[0.08em]">01.</p>
					<aside className="text-xs uppercase leading-5 tracking-[0.12em]">
						<span className="inline-flex h-9 w-9 items-center justify-center border border-foreground font-mono">P</span>
						<p className="mt-5 font-semibold">Parkeringskompagniet</p>
						<p className="mt-1 font-mono opacity-65">Nov 2020 - Mar 2022</p>
						<p className="mt-1 opacity-65">1 yr 5 mos</p>
					</aside>
					<div className="space-y-10">
						<div>
							<div className="flex flex-wrap items-baseline justify-between gap-3 border-b border-border pb-4">
								<h3 className="cobalt-row-title">Technical Lead</h3>
								<p className="font-mono text-xs uppercase tracking-[0.12em]">Full-time</p>
							</div>
							<p className="mt-4 font-mono text-xs uppercase tracking-[0.1em] opacity-65">Jul 2021 - Mar 2022 - 9 mos</p>
							<p className="mt-6 max-w-3xl text-sm leading-6">My work focused on software architecture, team building, and code quality. I worked primarily with Go while restructuring the backend from a legacy Java codebase, and React with TypeScript for the customer-facing dashboard. I also architected microservices for handling camera events so parking-lot customers could get reliable real-time information.</p>
							<p className="mt-4 max-w-3xl text-sm leading-6">I hired and built a remote team of junior engineers, with a particular focus on teaching domain-driven design, Go, React, and TypeScript.</p>
							<p className="mt-6 font-mono text-xs uppercase leading-5 tracking-[0.1em]"><span className="font-medium">Focus /</span> software architecture, Go, React, TypeScript, DevOps, team building</p>
						</div>
						<div className="border-t border-border pt-10">
							<div className="flex flex-wrap items-baseline justify-between gap-3 border-b border-border pb-4">
								<h3 className="cobalt-row-title">External Consultant</h3>
								<p className="font-mono text-xs uppercase tracking-[0.12em]">Freelance</p>
							</div>
							<p className="mt-4 font-mono text-xs uppercase tracking-[0.1em] opacity-65">Nov 2020 - Aug 2021 - 10 mos</p>
							<p className="mt-6 max-w-3xl text-sm leading-6">I consulted on software architecture and design, implemented CI/CD flows, and introduced React, TypeScript, and Go into the company's codebase and engineering practices.</p>
							<p className="mt-6 font-mono text-xs uppercase leading-5 tracking-[0.1em]"><span className="font-medium">Focus /</span> architecture, implementation, CI/CD, Go, React, TypeScript</p>
						</div>
					</div>
				</article>

				<article className="grid gap-6 border-b border-border py-12 md:grid-cols-[4rem_minmax(11rem,0.7fr)_minmax(0,1.3fr)] md:px-3">
					<p className="font-mono text-sm tracking-[0.08em]">02.</p>
					<aside className="text-xs uppercase leading-5 tracking-[0.12em]">
						<span className="inline-flex h-9 w-9 items-center justify-center border border-foreground font-mono">E</span>
						<p className="mt-5 font-semibold">everphone</p>
						<p className="mt-1 opacity-65">Berlin, Germany</p>
						<p className="mt-1 font-mono opacity-65">Feb 2021 - Aug 2021</p>
					</aside>
					<div>
						<div className="flex flex-wrap items-baseline justify-between gap-3 border-b border-border pb-4">
							<h3 className="cobalt-row-title">Full Stack Engineer</h3>
							<p className="font-mono text-xs uppercase tracking-[0.12em]">Full-time</p>
						</div>
						<p className="mt-4 font-mono text-xs uppercase tracking-[0.1em] opacity-65">7 mos</p>
						<p className="mt-6 max-w-3xl text-sm leading-6">I worked on improving code quality for the company's React dashboard, particularly by using TypeScript more deliberately. I also worked on moving the API from PHP to Go and from REST to GraphQL, while improving local DevOps through containerization.</p>
						<p className="mt-6 font-mono text-xs uppercase leading-5 tracking-[0.1em]"><span className="font-medium">Focus /</span> React, TypeScript, Go, GraphQL, DevOps, test-driven development</p>
					</div>
				</article>

				<article className="grid gap-6 border-b border-border py-12 md:grid-cols-[4rem_minmax(11rem,0.7fr)_minmax(0,1.3fr)] md:px-3">
					<p className="font-mono text-sm tracking-[0.08em]">03.</p>
					<aside className="text-xs uppercase leading-5 tracking-[0.12em]">
						<span className="inline-flex h-9 w-9 items-center justify-center border border-foreground font-mono">L</span>
						<p className="mt-5 font-semibold">Legal Monster</p>
						<p className="mt-1 font-mono opacity-65">Sep 2019 - Dec 2020</p>
						<p className="mt-1 opacity-65">1 yr 4 mos</p>
					</aside>
					<div>
						<div className="flex flex-wrap items-baseline justify-between gap-3 border-b border-border pb-4">
							<h3 className="cobalt-row-title">Full Stack Engineer</h3>
							<p className="font-mono text-xs uppercase tracking-[0.12em]">Full-time</p>
						</div>
						<p className="mt-6 max-w-3xl text-sm leading-6">I worked primarily with React, TypeScript, and Ruby on Rails. I was involved in the initial stages of product discovery and product building while the company was working toward product-market fit.</p>
						<p className="mt-6 font-mono text-xs uppercase leading-5 tracking-[0.1em]"><span className="font-medium">Focus /</span> React, TypeScript, Ruby on Rails, product discovery, DevOps, test-driven development</p>
					</div>
				</article>

				<article className="grid gap-6 border-b border-border py-12 md:grid-cols-[4rem_minmax(11rem,0.7fr)_minmax(0,1.3fr)] md:px-3">
					<p className="font-mono text-sm tracking-[0.08em]">04.</p>
					<aside className="text-xs uppercase leading-5 tracking-[0.12em]">
						<span className="inline-flex h-9 w-9 items-center justify-center border border-foreground font-mono">H</span>
						<p className="mt-5 font-semibold">HS MØLLER APS</p>
						<p className="mt-1 font-mono opacity-65">Jun 2018 - Aug 2019</p>
						<p className="mt-1 opacity-65">1 yr 3 mos</p>
					</aside>
					<div>
						<div className="flex flex-wrap items-baseline justify-between gap-3 border-b border-border pb-4">
							<h3 className="cobalt-row-title">Full Stack Engineer</h3>
							<p className="font-mono text-xs uppercase tracking-[0.12em]">Full-time</p>
						</div>
						<p className="mt-6 max-w-3xl text-sm leading-6">I worked on smaller React projects, primarily with a designer to improve clients' online presence and UX/UI. I also built internal Go tools to improve DevOps workflows.</p>
						<p className="mt-6 font-mono text-xs uppercase leading-5 tracking-[0.1em]"><span className="font-medium">Focus /</span> React, TypeScript, Go, UX/UI, DevOps</p>
					</div>
				</article>

				<article className="grid gap-6 border-b border-border py-12 md:grid-cols-[4rem_minmax(11rem,0.7fr)_minmax(0,1.3fr)] md:px-3">
					<p className="font-mono text-sm tracking-[0.08em]">05.</p>
					<aside className="text-xs uppercase leading-5 tracking-[0.12em]">
						<span className="inline-flex h-9 w-9 items-center justify-center border border-foreground font-mono">E</span>
						<p className="mt-5 font-semibold">Easy-skat ApS</p>
						<p className="mt-1 opacity-65">Copenhagen, Denmark</p>
						<p className="mt-1 font-mono opacity-65">Apr 2018 - Jun 2018</p>
					</aside>
					<div>
						<div className="flex flex-wrap items-baseline justify-between gap-3 border-b border-border pb-4">
							<h3 className="cobalt-row-title">Softwareudvikler</h3>
							<p className="font-mono text-xs uppercase tracking-[0.12em]">Internship</p>
						</div>
						<p className="mt-4 font-mono text-xs uppercase tracking-[0.1em] opacity-65">3 mos</p>
						<p className="mt-6 max-w-3xl text-sm leading-6">I worked on smaller React projects, primarily with a designer to improve clients' online presence and UX/UI.</p>
						<p className="mt-6 font-mono text-xs uppercase leading-5 tracking-[0.1em]"><span className="font-medium">Focus /</span> React, TypeScript, problem solving, UX/UI</p>
					</div>
				</article>

				<article className="grid gap-6 border-b border-border py-12 md:grid-cols-[4rem_minmax(11rem,0.7fr)_minmax(0,1.3fr)] md:px-3">
					<p className="font-mono text-sm tracking-[0.08em]">06.</p>
					<aside className="text-xs uppercase leading-5 tracking-[0.12em]">
						<span className="inline-flex h-9 w-9 items-center justify-center border border-foreground font-mono">E</span>
						<p className="mt-5 font-semibold">Ennova</p>
						<p className="mt-1 opacity-65">Copenhagen Metropolitan Area</p>
						<p className="mt-1 font-mono opacity-65">Aug 2016 - Dec 2016</p>
					</aside>
					<div>
						<div className="flex flex-wrap items-baseline justify-between gap-3 border-b border-border pb-4">
							<h3 className="cobalt-row-title">Praktikant</h3>
							<p className="font-mono text-xs uppercase tracking-[0.12em]">Internship</p>
						</div>
						<p className="mt-4 font-mono text-xs uppercase tracking-[0.1em] opacity-65">5 mos</p>
					</div>
				</article>
			</div>
			<div className="flex items-center justify-between gap-5 border-b border-foreground py-8 font-mono text-xs uppercase tracking-[0.1em]">
				<span>End of work report</span>
				<span>MBV / Archive</span>
			</div>
		</div>
	</section>

  )
}
