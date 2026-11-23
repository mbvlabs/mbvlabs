import { marked } from 'marked'
import sanitizeHTML from 'sanitize-html'

export default function MarkdownPreview({ source, className = '' }: { source?: string | null; className?: string }) {
  const html = sanitizeHTML(marked.parse(source ?? '', { async: false }), {
    allowedTags: [...sanitizeHTML.defaults.allowedTags, 'img', 'del', 'details', 'summary'],
    allowedAttributes: { ...sanitizeHTML.defaults.allowedAttributes, code: ['class'] },
  })
  return <article className={`markdown-preview ${className}`} dangerouslySetInnerHTML={{ __html: html }} />
}
