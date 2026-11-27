import { routes } from '@/routes'
import PublicLayout from '@/Layouts/PublicLayout'
import { Breadcrumbs, dateLabel, publishedAt, tagsOf } from './blogShared'
import type { PublicArticleIndexProps } from '@/types/payloads'

export default function BlogIndex({ items }: PublicArticleIndexProps) {
  const posts = items ?? []
  return <PublicLayout>
    <main id="blogpost-index-container" className="overflow-hidden bg-background font-sans text-foreground">
      <section className="border-y border-foreground">
        <div className="container mx-auto px-6 py-20 md:px-8 md:py-28">
          <Breadcrumbs />
          <h1 className="cobalt-display mt-8 max-w-6xl">Notes on building software that keeps moving.</h1>
          <p className="mt-9 max-w-2xl text-base leading-7 md:text-lg">Practical essays on product engineering, technical leadership, and the tradeoffs behind maintainable systems for focused teams.</p>
          <div className="mt-16 grid gap-5 border-t border-foreground pt-5 font-mono text-[11px] uppercase tracking-[0.12em] sm:grid-cols-3">
            <div><p className="font-medium">Format</p><p className="mt-1 opacity-65">Engineering dispatches</p></div>
            <div><p className="font-medium">Topics</p><p className="mt-1 opacity-65">Product / Systems / Teams</p></div>
            <div><p className="font-medium">Archive</p><p className="mt-1 opacity-65">Latest first</p></div>
          </div>
        </div>
      </section>
      <section className="container mx-auto px-6 py-20 md:px-8 md:py-24">
        <div className="flex flex-col gap-5 border-b border-foreground pb-7 md:flex-row md:items-end md:justify-between">
          <div><p className="font-mono text-xs uppercase tracking-[0.18em]">Dispatch archive</p><h2 className="cobalt-section-title mt-4">Latest blog posts.</h2></div>
          <p className="text-xs font-semibold uppercase tracking-[0.18em]">Practical notes / Current index</p>
        </div>
        {posts.length === 0 ? <div className="border-b border-border py-12"><p className="font-mono text-xs uppercase tracking-[0.14em]">No blog posts are published yet.</p></div> :
          <div>{posts.map((post, index) => {
            const tags = tagsOf(post)
            return <article key={post.ID}>
              <a href={routes.blogPostShow.path(encodeURIComponent(post.Slug))} className="group grid gap-6 border-b border-border py-12 transition-colors hover:bg-secondary/45 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring md:grid-cols-[4rem_minmax(10rem,0.65fr)_minmax(0,1.35fr)_7rem] md:items-start md:px-3">
                <p className="font-mono text-sm tracking-[0.08em]">{String(index + 1).padStart(2, '0')}.</p>
                <div className="font-mono text-xs uppercase leading-5 tracking-[0.1em]"><p className="font-medium">{dateLabel(publishedAt(post) ?? post.CreatedAt)}</p>
                  {tags.length > 0 && <ul className="mt-4 space-y-1 opacity-65" aria-label="Topics">{tags.map((tag, i) => <li key={i}>{tag}</li>)}</ul>}
                </div>
                <div><h3 className="cobalt-row-title">{post.Title}</h3>{post.Excerpt.Valid && post.Excerpt.String !== '' && <p className="mt-4 max-w-2xl text-sm leading-6">{post.Excerpt.String}</p>}</div>
                <span className="text-xs font-semibold uppercase tracking-[0.12em] transition-transform group-hover:translate-x-1 md:text-right">Read note →</span>
              </a>
            </article>
          })}</div>}
      </section>
      {posts.length > 0 && <>
        <section className="flex min-h-[60svh] items-center border-y border-foreground">
          <div className="container mx-auto px-6 py-24 md:px-8 md:py-28">
            <p className="font-mono text-xs uppercase tracking-[0.18em]">Working principle / 01</p>
            <blockquote className="cobalt-quote mt-10 max-w-6xl">Engineering decisions are product decisions when the team is small.</blockquote>
            <div className="mt-14 flex items-center justify-between gap-5 border-t border-foreground pt-5 font-mono text-xs uppercase tracking-[0.12em]"><p className="font-medium">MBV Labs</p><p className="opacity-65">Decide / Build / Learn</p></div>
          </div>
        </section>
        <section className="container mx-auto px-6 py-20 md:px-8 md:py-24">
          <div className="flex flex-col gap-5 border-b border-foreground pb-7 md:flex-row md:items-end md:justify-between"><h2 className="cobalt-section-title">What I write about.</h2><p className="text-xs font-semibold uppercase tracking-[0.18em]">Editorial scope / Three fields</p></div>
          <div className="grid md:grid-cols-3">
            <div className="border-b border-border py-10 md:border-r md:px-8 md:first:pl-0"><p className="font-mono text-xs tracking-[0.08em]">01.</p><h3 className="cobalt-row-title mt-6">Product engineering</h3><p className="mt-4 text-sm leading-6">Turning ambiguous goals into shipped systems.</p></div>
            <div className="border-b border-border py-10 md:border-r md:px-8"><p className="font-mono text-xs tracking-[0.08em]">02.</p><h3 className="cobalt-row-title mt-6">Technical leadership</h3><p className="mt-4 text-sm leading-6">Clear choices for founders and lean teams.</p></div>
            <div className="border-b border-border py-10 md:px-8 md:last:pr-0"><p className="font-mono text-xs tracking-[0.08em]">03.</p><h3 className="cobalt-row-title mt-6">Maintainable systems</h3><p className="mt-4 text-sm leading-6">Architecture, operations, and code that can change.</p></div>
          </div>
        </section>
      </>}
    </main>
  </PublicLayout>
}
