import { useEffect } from 'react'
import { EditorContent, useEditor } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import { TableKit } from '@tiptap/extension-table'
import Placeholder from '@tiptap/extension-placeholder'
import { Markdown, type MarkdownStorage } from 'tiptap-markdown'
import { Button } from '@/Components/ui/button'

type Props = { id: string; value: string; onChange: (value: string) => void; placeholder?: string; className?: string }

export default function RichEditor({ id, value, onChange, placeholder, className = '' }: Props) {
  const editor = useEditor({
    immediatelyRender: false,
    content: value,
    extensions: [StarterKit.configure({ heading: { levels: [2, 3, 4] }, link: { openOnClick: false } }), TableKit, Placeholder.configure({ placeholder: placeholder || 'Start writing...' }), Markdown],
    editorProps: { attributes: { id, 'aria-label': placeholder || 'Editor' } },
    onUpdate: ({ editor }) => onChange((editor.storage as unknown as { markdown: MarkdownStorage }).markdown.getMarkdown()),
  })
  useEffect(() => {
    if (editor && value !== (editor.storage as unknown as { markdown: MarkdownStorage }).markdown.getMarkdown()) editor.commands.setContent(value, { emitUpdate: false })
  }, [editor, value])
  if (!editor) return <div className={`h-160 border bg-card ${className}`} />
  const button = (label: string, action: () => void, active = false, disabled = false, title = label) =>
    <Button variant="ghost" type="button" key={label} onClick={action} aria-label={title} title={title} aria-pressed={active} disabled={disabled} className={`h-auto min-h-8 px-2 text-sm font-normal normal-case tracking-normal hover:bg-muted disabled:opacity-40 ${active ? 'bg-muted' : ''}`}>{label}</Button>
  const link = () => {
    const url = window.prompt('Link URL:', editor.getAttributes('link').href || '')
    if (url === null) return
    if (!url) editor.chain().focus().extendMarkRange('link').unsetLink().run()
    else editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
  }
  const table = editor.isActive('table')
  return <div className={`flex h-160 flex-col overflow-hidden border bg-card ${className}`}>
    <div className="flex flex-wrap gap-1 border-b px-3 py-2" role="toolbar" aria-label="Formatting">
      {button('B', () => { editor.chain().focus().toggleBold().run() }, editor.isActive('bold'), false, 'Bold')}
      {button('I', () => { editor.chain().focus().toggleItalic().run() }, editor.isActive('italic'), false, 'Italic')}
      {([2, 3, 4] as const).map(level => button(`H${level}`, () => { editor.chain().focus().toggleHeading({ level }).run() }, editor.isActive('heading', { level })))}
      {button('•≡', () => { editor.chain().focus().toggleBulletList().run() }, editor.isActive('bulletList'), false, 'Bullet list')}
      {button('1.', () => { editor.chain().focus().toggleOrderedList().run() }, editor.isActive('orderedList'), false, 'Ordered list')}
      {button('„”', () => { editor.chain().focus().toggleBlockquote().run() }, editor.isActive('blockquote'), false, 'Blockquote')}
      {button('—', () => { editor.chain().focus().setHorizontalRule().run() }, false, false, 'Horizontal rule')}
      {button('Link', link, editor.isActive('link'), false, 'Add or edit link')}
      {button('Unlink', () => { editor.chain().focus().extendMarkRange('link').unsetLink().run() }, false, !editor.isActive('link'), 'Remove link')}
      {button('Table', () => { editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run() }, false, table, 'Insert table')}
      {button('+ Row', () => { editor.chain().focus().addRowAfter().run() }, false, !table, 'Add table row')}
      {button('+ Column', () => { editor.chain().focus().addColumnAfter().run() }, false, !table, 'Add table column')}
      {button('− Row', () => { editor.chain().focus().deleteRow().run() }, false, !table, 'Delete table row')}
      {button('− Column', () => { editor.chain().focus().deleteColumn().run() }, false, !table, 'Delete table column')}
      {button('Delete table', () => { editor.chain().focus().deleteTable().run() }, false, !table)}
    </div>
    <div className="flex min-h-0 flex-1 overflow-y-auto p-4"><EditorContent editor={editor} className="flex-1" /></div>
  </div>
}
