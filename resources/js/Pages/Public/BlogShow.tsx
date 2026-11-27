import { routes } from '@/routes'
import PublicLayout from '@/Layouts/PublicLayout'
import { Breadcrumbs, dateLabel, publishedAt, tagsOf } from './blogShared'
import type { PublicArticleShowProps } from '@/types/payloads'

export default function BlogShow({ item, markdownHTML, relatedArticles }: PublicArticleShowProps) {
  const tags = tagsOf(item)
  const published = publishedAt(item)
  const publishedLabel = dateLabel(published ?? item.CreatedAt)
  const showUpdated = published !== null && new Date(item.UpdatedAt) > new Date(published)
  const excerpt = item.Excerpt.Valid ? item.Excerpt.String : ''
  const cover = item.CoverImageUrl.Valid ? item.CoverImageUrl.String : ''
  return <PublicLayout>
    <main id="blogpost-show-container" className="overflow-hidden bg-background font-sans text-foreground">
      <section className="border-y border-foreground">
        <div className="container mx-auto px-6 py-20 md:px-8 md:py-28">
          <Breadcrumbs article />
          <h1 className="cobalt-display mt-8 max-w-6xl">{item.Title}</h1>
          {excerpt !== '' && <p className="mt-9 max-w-2xl text-base leading-7 md:text-lg">{excerpt}</p>}
        </div>
      </section>
      <section className="border-b border-foreground">
        <div className="container mx-auto grid px-6 md:grid-cols-3 md:px-8">
          <div className="border-b border-border py-6 md:border-b-0 md:border-r md:pr-8"><p className="font-mono text-xs uppercase tracking-[0.14em]">01 / Published</p><p className="mt-3 text-sm leading-6">{publishedLabel}</p></div>
          <div className="border-b border-border py-6 md:border-b-0 md:border-r md:px-8">{showUpdated ? <><p className="font-mono text-xs uppercase tracking-[0.14em]">02 / Updated</p><p className="mt-3 text-sm leading-6">{dateLabel(item.UpdatedAt)}</p></> : <><p className="font-mono text-xs uppercase tracking-[0.14em]">02 / Category</p><p className="mt-3 text-sm leading-6">Blog</p></>}</div>
          <div className="py-6 md:pl-8"><p className="font-mono text-xs uppercase tracking-[0.14em]">03 / Topics</p><p className="mt-3 text-sm leading-6">{tags.length ? tags.join(', ') : 'Product engineering'}</p></div>
        </div>
      </section>
      {cover !== '' && <section className="border-b border-foreground bg-secondary/45"><div className="container mx-auto px-6 py-10 md:px-8 md:py-14"><img src={cover} alt={item.Title} className="w-full border border-foreground bg-muted object-cover" /></div></section>}
      <section className="container mx-auto grid gap-16 px-6 py-20 md:px-8 md:py-24 lg:grid-cols-[minmax(12rem,0.45fr)_minmax(0,1.55fr)] lg:gap-24">
        <aside className="text-xs uppercase leading-5 tracking-[0.12em] lg:sticky lg:top-24 lg:self-start">
          <p className="font-mono opacity-65">Article record</p>
          <div className="mt-5 border-t border-foreground pt-5">
            <p className="font-semibold">{publishedLabel}</p>
            {showUpdated && <p className="mt-2 font-mono opacity-65">Updated {dateLabel(item.UpdatedAt)}</p>}
            {tags.length > 0 && <div className="mt-6"><p className="font-mono opacity-65">Topics</p><ul className="mt-3 space-y-1">{tags.map((tag, i) => <li key={i}>{tag}</li>)}</ul></div>}
          </div>
          <a href={routes.blogPostIndex.path()} className="mt-8 inline-flex min-h-11 items-center border border-foreground px-4 py-2 font-semibold transition-colors hover:bg-primary hover:text-primary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">← All blog posts</a>
        </aside>
        <div className="min-w-0">
          {item.Body !== '' ? <article className="prose max-w-3xl border-t border-foreground pt-10 text-foreground prose-headings:font-heading prose-headings:font-bold prose-headings:text-foreground prose-p:text-foreground prose-a:text-foreground prose-a:underline prose-a:decoration-1 prose-a:underline-offset-4 prose-strong:text-foreground prose-code:text-foreground prose-code:before:content-none prose-code:after:content-none prose-li:text-foreground prose-ul:marker:text-foreground prose-ol:marker:text-foreground prose-hr:border-foreground prose-img:rounded-none prose-img:border prose-img:border-foreground prose-pre:rounded-none prose-pre:border prose-pre:border-foreground prose-pre:bg-primary prose-pre:text-primary-foreground prose-blockquote:border-foreground prose-blockquote:font-heading prose-blockquote:font-normal prose-blockquote:text-foreground" dangerouslySetInnerHTML={{ __html: markdownHTML }} /> :
            excerpt !== '' && <article className="max-w-3xl border-t border-foreground pt-10 text-base leading-7"><p>{excerpt}</p></article>}
        </div>
      </section>
      <section className="border-y border-foreground">
        <div className="container mx-auto px-6 py-20 md:px-8 md:py-24">
          <p className="font-mono text-xs uppercase tracking-[0.18em]">Archive / Continue reading</p>
          <h2 className="cobalt-section-title mt-5">More from the blog.</h2>
          {!!relatedArticles?.length && <ul className="mt-10 border-t border-border">
            {relatedArticles.map(article => <li key={article.slug} className="border-b border-border">
              <a href={routes.blogPostShow.path(encodeURIComponent(article.slug))} className="group flex min-h-11 items-center justify-between gap-6 py-6 text-base font-semibold underline underline-offset-4 transition-colors hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring md:px-3">
                {article.title}<span aria-hidden="true" className="transition-transform group-hover:translate-x-1">→</span>
              </a>
            </li>)}
          </ul>}
          <a href={routes.blogPostIndex.path()} className="group mt-10 grid gap-6 border-y border-border py-10 transition-colors hover:bg-secondary/45 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring md:grid-cols-[minmax(0,1fr)_8rem] md:items-end md:px-3">
            <div><p className="max-w-3xl text-sm leading-6">Read more practical notes on product engineering, technical leadership, software architecture, and the delivery tradeoffs behind maintainable systems for focused teams.</p><h3 className="cobalt-row-title mt-6">Read all blog posts</h3></div>
            <span className="text-xs font-semibold uppercase tracking-[0.12em] transition-transform group-hover:translate-x-1 md:text-right">Open index →</span>
          </a>
        </div>
      </section>
    </main>
  </PublicLayout>
}
