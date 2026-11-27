import { routes } from '@/routes'
import type { ArticleItem } from '@/types/payloads'

export function tagsOf(item: ArticleItem): string[] {
  return Array.isArray(item.Tags) && item.Tags.every(tag => typeof tag === 'string')
    ? item.Tags.map(tag => tag.trim()).filter(Boolean)
    : []
}

export function dateLabel(value: string): string {
  // Go formats the date in the timestamp's own offset, not the browser's timezone.
  const [year, month, day] = value.slice(0, 10).split('-')
  const months = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December']
  return `${months[Number(month) - 1]} ${Number(day)}, ${year}`
}

export function publishedAt(item: ArticleItem): string | null {
  return item.PublishedAt.Valid && !item.PublishedAt.Time.startsWith('0001-01-01T00:00:00')
    ? item.PublishedAt.Time : null
}

export function Breadcrumbs({ article = false }: { article?: boolean }) {
  return <nav aria-label="Breadcrumb">
    <ol className="flex flex-wrap items-center gap-x-1.5 gap-y-1 font-mono text-xs uppercase tracking-[0.18em] text-muted-foreground">
      <li className="inline-flex items-center gap-1.5"><a href={routes.homePage.path()} className="transition-colors hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Home</a></li>
      <li className="inline-flex items-center gap-1.5"><span aria-hidden="true">/</span>{article ? <a href={routes.blogPostIndex.path()} className="transition-colors hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Blog</a> : <span className="font-medium text-foreground" aria-current="page">Blog</span>}</li>
      {article && <li className="inline-flex items-center gap-1.5"><span aria-hidden="true">/</span><span className="font-medium text-foreground" aria-current="page">Article</span></li>}
    </ol>
  </nav>
}
