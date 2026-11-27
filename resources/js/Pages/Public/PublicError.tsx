import PublicLayout from '@/Layouts/PublicLayout'

export default function PublicError({ status }: { status: number }) {
	const notFound = status === 404
	const badRequest = status === 400
	const id = badRequest ? 'bad-request-container' : notFound ? 'not-found-container' : 'internal-error-container'
	const title = badRequest ? 'The request made was invalid.' : notFound ? 'Not found.' : 'Something went wrong on our end.'
	const description = badRequest
		? 'The request could not be completed. Return home and try the action again when you are ready.'
		: notFound
			? 'The page you requested is not in this field report. Return home or continue through the public archive.'
			: 'The request could not be completed right now. Return home or contact MBV Labs if you need help continuing.'

	return (
		<PublicLayout>
			<main id={id} className="flex flex-1 bg-background font-sans text-foreground">
				<section className="flex w-full items-center border-y border-foreground">
					<div className="container mx-auto px-6 py-24 md:px-8 md:py-32">
						<p className="font-mono text-xs uppercase tracking-[0.18em]">System note / {status}</p>
						<h1 className="cobalt-display mt-8 max-w-6xl">{title}</h1>
						<p className="mt-8 max-w-xl text-base leading-7">{description}</p>
						<nav className="mt-10 flex flex-wrap gap-3" aria-label="Recovery navigation">
							<a href="/" className="inline-flex min-h-11 items-center border border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Return home</a>
							{notFound && <a href="/work" className="inline-flex min-h-11 items-center border border-foreground px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] transition-colors hover:bg-primary hover:text-primary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">View work</a>}
							<a href="/contact" className="inline-flex min-h-11 items-center border border-foreground px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] transition-colors hover:bg-primary hover:text-primary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Contact</a>
						</nav>
						<div className="mt-16 flex flex-col gap-2 border-t border-foreground pt-5 font-mono text-[11px] uppercase tracking-[0.12em] sm:flex-row sm:items-center sm:justify-between">
							<p className="font-medium">{badRequest ? 'Request / Not completed' : notFound ? 'Requested page / Unavailable' : 'Service / Temporarily unavailable'}</p>
							<p className="opacity-65">{badRequest ? 'No submitted details displayed' : notFound ? 'MBV Labs / Public system' : 'Safe recovery / Home or contact'}</p>
						</div>
					</div>
				</section>
			</main>
		</PublicLayout>
	)
}
