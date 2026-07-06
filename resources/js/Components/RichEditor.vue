<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import { Markdown } from 'tiptap-markdown'
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
    }),
    Placeholder.configure({
      placeholder: props.placeholder || 'Start writing...',
    }),
    Markdown,
  ],
  onUpdate: () => {
    emit('update:modelValue', editor.value?.storage.markdown.getMarkdown() ?? '')
  },
})

watch(() => props.modelValue, (val) => {
  if (editor.value && val !== editor.value.storage.markdown.getMarkdown()) {
    editor.value.commands.setContent(val, false)
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
  const url = window.prompt('Link URL:')
  if (url) {
    editor.value?.chain().focus().setLink({ href: url }).run()
  }
}
</script>

<template>
  <div class="flex h-full flex-col rounded-none border bg-card">
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
      <Toggle size="sm" :pressed="editor?.isActive('link')" @click="setLink">
        <span class="text-xs underline">⧉</span>
      </Toggle>
    </div>
    <div class="flex min-h-0 flex-1 p-4">
      <EditorContent :editor="editor" class="prose prose-sm flex-1 max-w-none dark:prose-invert focus:outline-none" />
    </div>
  </div>
</template>

<style>
.ProseMirror {
  outline: none;
  height: 100%;
}
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
.ProseMirror img { max-width: 100%; height: auto; border-radius: 0; margin: 1rem 0; }
</style>
