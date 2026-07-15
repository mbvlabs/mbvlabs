<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { TableKit } from '@tiptap/extension-table'
import Placeholder from '@tiptap/extension-placeholder'
import { Markdown, type MarkdownStorage } from 'tiptap-markdown'
import { onBeforeUnmount, watch } from 'vue'

import { Button } from '@/components/ui/button'
import { Toggle } from '@/components/ui/toggle'

const props = defineProps<{
  modelValue: string
  placeholder?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit.configure({
      heading: { levels: [2, 3, 4] },
      link: { openOnClick: false },
    }),
    TableKit,
    Placeholder.configure({
      placeholder: props.placeholder || 'Start writing...',
    }),
    Markdown,
  ],
  onUpdate: () => {
    emit('update:modelValue', getMarkdown())
  },
})

watch(() => props.modelValue, (val) => {
  if (editor.value && val !== getMarkdown()) {
    editor.value.commands.setContent(val, { emitUpdate: false })
  }
})

onBeforeUnmount(() => {
  editor.value?.destroy()
})

function toggleBold() {
  editor.value?.chain().focus().toggleBold().run()
}

function toggleItalic() {
  editor.value?.chain().focus().toggleItalic().run()
}

function toggleHeading(level: 2 | 3 | 4) {
  editor.value?.chain().focus().toggleHeading({ level }).run()
}

function toggleBulletList() {
  editor.value?.chain().focus().toggleBulletList().run()
}

function toggleOrderedList() {
  editor.value?.chain().focus().toggleOrderedList().run()
}

function toggleBlockquote() {
  editor.value?.chain().focus().toggleBlockquote().run()
}

function toggleHorizontalRule() {
  editor.value?.chain().focus().setHorizontalRule().run()
}

function setLink() {
  const url = window.prompt('Link URL:', editor.value?.getAttributes('link').href || '')
  if (url === null) return
  if (url === '') {
    editor.value?.chain().focus().extendMarkRange('link').unsetLink().run()
    return
  }
  editor.value?.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
}

function insertTable() {
  editor.value?.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()
}

function getMarkdown() {
  const storage = editor.value?.storage as { markdown?: MarkdownStorage } | undefined
  return storage?.markdown?.getMarkdown() ?? ''
}
</script>

<template>
  <div class="flex h-160! flex-col overflow-hidden rounded-none border bg-card">
    <div class="flex flex-wrap gap-1 border-b px-3 py-2">
      <Toggle size="sm" :pressed="editor?.isActive('bold')" @click="toggleBold">
        <span class="font-bold">B</span>
      </Toggle>
      <Toggle size="sm" :pressed="editor?.isActive('italic')" @click="toggleItalic">
        <span class="italic">I</span>
      </Toggle>
      <div class="mx-1 w-px bg-border" />
      <Toggle size="sm" :pressed="editor?.isActive('heading', { level: 2 })" @click="toggleHeading(2)">
        H2
      </Toggle>
      <Toggle size="sm" :pressed="editor?.isActive('heading', { level: 3 })" @click="toggleHeading(3)">
        H3
      </Toggle>
      <Toggle size="sm" :pressed="editor?.isActive('heading', { level: 4 })" @click="toggleHeading(4)">
        H4
      </Toggle>
      <div class="mx-1 w-px bg-border" />
      <Toggle size="sm" :pressed="editor?.isActive('bulletList')" @click="toggleBulletList">
        <span class="text-xs">•≡</span>
      </Toggle>
      <Toggle size="sm" :pressed="editor?.isActive('orderedList')" @click="toggleOrderedList">
        <span class="text-xs">1.</span>
      </Toggle>
      <Toggle size="sm" :pressed="editor?.isActive('blockquote')" @click="toggleBlockquote">
        <span class="text-xs">„"</span>
      </Toggle>
      <Toggle size="sm" @click="toggleHorizontalRule">
        <span class="text-xs">—</span>
      </Toggle>
      <div class="mx-1 w-px bg-border" />
      <Toggle
        size="sm"
        :pressed="editor?.isActive('link')"
        aria-label="Add or edit link"
        title="Add or edit link"
        @click="setLink"
      >
        Link
      </Toggle>
      <Toggle
        size="sm"
        :disabled="!editor?.isActive('link')"
        aria-label="Remove link"
        title="Remove link"
        @click="editor?.chain().focus().extendMarkRange('link').unsetLink().run()"
      >
        Unlink
      </Toggle>
      <div class="mx-1 w-px bg-border" />
      <Toggle
        size="sm"
        :disabled="editor?.isActive('table')"
        aria-label="Insert table"
        title="Insert 3 by 3 table"
        @click="insertTable"
      >
        Table
      </Toggle>
      <Toggle
        size="sm"
        :disabled="!editor?.isActive('table')"
        aria-label="Add table row"
        title="Add row after"
        @click="editor?.chain().focus().addRowAfter().run()"
      >
        + Row
      </Toggle>
      <Toggle
        size="sm"
        :disabled="!editor?.isActive('table')"
        aria-label="Add table column"
        title="Add column after"
        @click="editor?.chain().focus().addColumnAfter().run()"
      >
        + Column
      </Toggle>
      <Toggle
        size="sm"
        :disabled="!editor?.isActive('table')"
        aria-label="Delete table row"
        title="Delete row"
        @click="editor?.chain().focus().deleteRow().run()"
      >
        − Row
      </Toggle>
      <Toggle
        size="sm"
        :disabled="!editor?.isActive('table')"
        aria-label="Delete table column"
        title="Delete column"
        @click="editor?.chain().focus().deleteColumn().run()"
      >
        − Column
      </Toggle>
      <Toggle
        size="sm"
        :disabled="!editor?.isActive('table')"
        aria-label="Delete table"
        title="Delete table"
        @click="editor?.chain().focus().deleteTable().run()"
      >
        Delete table
      </Toggle>
    </div>
    <div class="flex min-h-0 flex-1 overflow-y-auto p-4">
      <EditorContent :editor="editor" class="prose prose-sm flex-1 max-w-none dark:prose-invert focus:outline-none" />
    </div>
  </div>
</template>

<style>
.ProseMirror {
  outline: none;
  box-sizing: border-box;
  min-height: 100%;
}
.ProseMirror > :last-child { margin-bottom: 0; }
.ProseMirror .ProseMirror-gapcursor { display: none; }
.ProseMirror p.is-editor-empty:first-child::before {
  color: hsl(var(--muted-foreground));
  content: attr(data-placeholder);
  float: left;
  height: 0;
  pointer-events: none;
}
.ProseMirror h2 { font-size: 1.5rem; font-weight: 600; margin-top: 1.5rem; margin-bottom: 0.5rem; }
.ProseMirror h3 { font-size: 1.25rem; font-weight: 600; margin-top: 1.25rem; margin-bottom: 0.5rem; }
.ProseMirror h4 { font-size: 1.1rem; font-weight: 600; margin-top: 1rem; margin-bottom: 0.5rem; }
.ProseMirror p { margin-bottom: 0.75rem; line-height: 1.7; }
.ProseMirror ul, .ProseMirror ol { padding-left: 1.5rem; margin-bottom: 0.75rem; }
.ProseMirror li { margin-bottom: 0.25rem; }
.ProseMirror blockquote {
  border-left: 3px solid hsl(var(--border));
  padding-left: 1rem;
  margin-left: 0;
  margin-right: 0;
  font-style: italic;
  color: hsl(var(--muted-foreground));
}
.ProseMirror hr { border: none; border-top: 1px solid hsl(var(--border)); margin: 1.5rem 0; }
.ProseMirror a { color: hsl(var(--primary)); text-decoration: underline; cursor: pointer; }
.ProseMirror .tableWrapper { margin: 1rem 0; overflow-x: auto; }
.ProseMirror table { border-collapse: collapse; table-layout: fixed; width: 100%; }
.ProseMirror th, .ProseMirror td { border: 1px solid hsl(var(--border)); min-width: 6rem; padding: 0.5rem; vertical-align: top; }
.ProseMirror th { background: hsl(var(--muted)); font-weight: 600; text-align: left; }
.ProseMirror .selectedCell { background: hsl(var(--muted) / 0.5); }
.ProseMirror img { max-width: 100%; height: auto; border-radius: 0; margin: 1rem 0; }
</style>
