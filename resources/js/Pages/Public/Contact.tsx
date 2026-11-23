import PublicLayout from '@/Layouts/PublicLayout'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import { NativeSelect } from '@/Components/ui/native-select'

// The public stylesheet already supplies the legacy select chevron.
const selectClass = 'w-full [&_select]:h-12 [&_select]:border-0 [&_select]:px-5 [&_select]:text-foreground [&_select:focus-visible]:ring-2 [&_select:focus-visible]:ring-inset [&_select:focus-visible]:ring-ring [&_svg]:hidden'
import { Textarea } from '@/Components/ui/textarea'

export default function Contact({ csrfToken }: { csrfToken?: string }) {
  return (
    <PublicLayout>
			<main id="projectinquiry-index-container" className="overflow-hidden bg-background font-sans text-foreground">
				<section className="border-y border-foreground">
					<div className="container mx-auto px-6 py-20 md:px-8 md:py-28">
						<p className="font-mono text-xs uppercase tracking-[0.18em]">Intake sheet / Project inquiry</p>
						<h1 className="cobalt-display mt-8 max-w-6xl">Let's work together.</h1>
						<p className="mt-9 max-w-2xl text-base leading-7 md:text-lg">Send a few details about what you are building, where the technical pressure is, and what kind of help would move the project forward.</p>
						<div className="mt-16 grid gap-5 border-t border-foreground pt-5 font-mono text-[11px] uppercase tracking-[0.12em] sm:grid-cols-3">
							<div><p className="font-medium">Channel</p><p className="mt-1 opacity-65">Project inquiry</p></div>
							<div><p className="font-medium">Location</p><p className="mt-1 opacity-65">Remote</p></div>
							<div><p className="font-medium">Inbox</p><p className="mt-1 opacity-65">hello@mbvlabs.com</p></div>
						</div>
					</div>
				</section>

				<section className="container mx-auto px-6 py-20 md:px-8 md:py-24">
					<div className="grid gap-16 lg:grid-cols-[0.65fr_1.35fr] lg:gap-24">
						<aside>
							<p className="font-mono text-xs uppercase tracking-[0.18em]">Direct contact / 01</p>
							<h2 className="cobalt-section-title mt-5">Email me.</h2>
							<p className="mt-6 max-w-md text-sm leading-6">Prefer to skip the form? Send the context directly and I will reply from the same inbox.</p>
							<a href="mailto:hello@mbvlabs.com" data-palantir-event="contact-email" className="mt-6 inline-flex min-h-11 items-center border border-foreground px-4 py-2 text-sm font-semibold transition-colors hover:bg-primary hover:text-primary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">hello@mbvlabs.com</a>
							<div className="mt-12 border-t border-foreground pt-7">
								<p className="font-mono text-xs uppercase tracking-[0.14em]">No public office</p>
								<p className="mt-5 max-w-md text-sm leading-6">MBV Labs works remotely with founders and lean teams. There is no public office address, careers page, or press desk to route through.</p>
							</div>
						</aside>

						<div>
							<div className="flex items-end justify-between gap-5 border-b border-foreground pb-6">
								<h2 className="cobalt-section-title">Project inquiry.</h2>
								<p className="hidden font-mono text-xs uppercase tracking-[0.12em] sm:block">Fields / 01-07</p>
							</div>
							<form method="post" action="/contact" className="mt-8 border border-foreground">
								{csrfToken && <input type="hidden" name="_csrf" value={csrfToken} />}
								<div className="divide-y divide-border">
									<div>
										<Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="name">01 / Name</Label>
										<Input id="name" name="name" type="text" autoComplete="name" required minLength={2} maxLength={255} placeholder="Name" className="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none placeholder:text-muted-foreground focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring"/>
									</div>
									<div>
										<Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="email">02 / Email</Label>
										<Input id="email" name="email" type="email" autoComplete="email" required maxLength={255} placeholder="Email" className="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none placeholder:text-muted-foreground focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring"/>
									</div>
									<div>
										<Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="company">03 / Company</Label>
										<Input id="company" name="company" type="text" autoComplete="organization" maxLength={255} placeholder="Company" className="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none placeholder:text-muted-foreground focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring"/>
									</div>
									<div>
										<Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="role">04 / Role</Label>
										<Input id="role" name="role" type="text" autoComplete="organization-title" maxLength={255} placeholder="Role" className="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none placeholder:text-muted-foreground focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring"/>
									</div>
									<div className="grid md:grid-cols-2 md:divide-x md:divide-border">
										<div>
											<Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="projectType">05 / Project type</Label>
											<NativeSelect id="projectType" name="projectType" className={selectClass}>
												<option value="">Project type</option>
												<option value="web-app">Web app</option>
												<option value="website">Website</option>
												<option value="automation">Automation</option>
												<option value="integration">Integration</option>
												<option value="consulting">Consulting</option>
												<option value="modernization">Modernization</option>
											</NativeSelect>
										</div>
										<div className="border-t border-border md:border-t-0">
											<Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="timeline">06 / Timeline</Label>
											<NativeSelect id="timeline" name="timeline" className={selectClass}>
												<option value="">Timeline</option>
												<option value="asap">ASAP</option>
												<option value="1-3-months">1-3 months</option>
												<option value="3-6-months">3-6 months</option>
												<option value="flexible">Flexible</option>
											</NativeSelect>
										</div>
									</div>
									<div>
										<Label className="block px-5 pt-4 font-mono text-xs font-normal uppercase tracking-[0.12em]" htmlFor="message">07 / Message</Label>
										<Textarea id="message" name="message" required minLength={10} rows={5} placeholder="Message" className="min-h-36 w-full field-sizing-fixed resize-y border-0 bg-transparent px-5 py-4 text-sm leading-6 text-foreground outline-none placeholder:text-muted-foreground focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring"></Textarea>
									</div>
								</div>
								<Button type="submit" data-palantir-event="inquiry-submit-click" className="inline-flex min-h-14 w-full items-center justify-center border-0 border-t border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring">
									Send inquiry
								</Button>
							</form>
						</div>
					</div>
				</section>

				<section className="flex min-h-[60svh] items-center border-y border-foreground">
					<div className="container mx-auto px-6 py-24 md:px-8 md:py-28">
						<p className="font-mono text-xs uppercase tracking-[0.18em]">Intake principle / 01</p>
						<blockquote className="cobalt-quote mt-10 max-w-6xl">Good project notes make the first technical conversation much sharper.</blockquote>
						<div className="mt-14 flex items-center justify-between gap-5 border-t border-foreground pt-5 font-mono text-xs uppercase tracking-[0.12em]">
							<p className="font-medium">MBV Labs</p>
							<p className="opacity-65">Context / Pressure / Outcome</p>
						</div>
					</div>
				</section>
			</main>
    </PublicLayout>
  )
}
