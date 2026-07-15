<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { ref, watch } from 'vue'
import { ArrowLeft } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'

import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import RichEditor from '@/Components/RichEditor.vue'

defineOptions({ layout: AdminLayout })

interface Work {
  ID: number
  Title: string
  Slug: string
  ClientName: string | null
  ClientIndustry: string | null
  ClientUrl: string | null
  ClientLogoUrl: string | null
  Summary: string
  Content: string
  CoverImageUrl: string | null
  Specialisms: string[] | null
  Platforms: string[] | null
  Technologies: string[] | null
  Challenge: string
  Approach: string
  Deliverables: string
  Outcome: string
  StartedAt: string | null
  CompletedAt: string | null
  Status: string
  IsFeatured: boolean
  CreatedAt: string
  UpdatedAt: string
}

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Published', value: 'published' },
]

const props = defineProps<{
  item: Work
}>()

const form = useForm({
  title: props.item.Title,
  slug: props.item.Slug,
  clientName: props.item.ClientName ?? '',
  clientIndustry: props.item.ClientIndustry ?? '',
  clientUrl: props.item.ClientUrl ?? '',
  clientLogoUrl: props.item.ClientLogoUrl ?? '',
  summary: props.item.Summary,
  content: props.item.Content,
  coverImageUrl: props.item.CoverImageUrl ?? '',
  specialisms: (props.item.Specialisms ?? []).join(', '),
  platforms: (props.item.Platforms ?? []).join(', '),
  technologies: (props.item.Technologies ?? []).join(', '),
  challenge: props.item.Challenge,
  approach: props.item.Approach,
  deliverables: props.item.Deliverables,
  outcome: props.item.Outcome,
  startedAt: props.item.StartedAt ?? '',
  completedAt: props.item.CompletedAt ?? '',
  status: props.item.Status || 'draft',
  isFeatured: props.item.IsFeatured,
})

const slugManuallyEdited = ref(false)

function slugify(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}

watch(() => form.title, (val) => {
  if (!slugManuallyEdited.value) {
    form.slug = slugify(val)
  }
})

function onSlugInput() {
  slugManuallyEdited.value = true
}

function setFeatured(checked: boolean | 'indeterminate') {
  form.isFeatured = checked === true
}

function submit() {
  form
    .transform((data) => ({
      ...data,
      specialisms: typeof data.specialisms === 'string'
        ? data.specialisms.split(',').map((t) => t.trim()).filter(Boolean)
        : data.specialisms,
      platforms: typeof data.platforms === 'string'
        ? data.platforms.split(',').map((t) => t.trim()).filter(Boolean)
        : data.platforms,
      technologies: typeof data.technologies === 'string'
        ? data.technologies.split(',').map((t) => t.trim()).filter(Boolean)
        : data.technologies,
    }))
    .put(`/admin/works/${props.item.ID}`)
}
</script>

<template>
  <Head :title="`Edit ${item.Title}`" />

  <div class="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Edit Work</h1>
        <p class="text-muted-foreground">{{ item.Title }}</p>
      </div>
      <Link :href="`/admin/works/${item.ID}`">
        <Button variant="outline">
          <ArrowLeft class="mr-2 size-4" />
          Cancel
        </Button>
      </Link>
    </div>

    <form @submit.prevent="submit" class="border bg-card">
      <div class="space-y-8 p-6">
        <div
          v-if="Object.keys(form.errors).length > 0"
          class="border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive"
        >
          <p class="font-medium">Please fix the highlighted fields.</p>
        </div>

        <Accordion type="multiple" :default-value="['primary', 'content']" class="w-full">
          <AccordionItem value="primary">
            <AccordionTrigger>Primary details</AccordionTrigger>
            <AccordionContent>
              <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
                <div class="space-y-1">
                  <Label for="title">Title</Label>
                  <Input id="title" v-model="form.title" />
                  <p v-if="form.errors.title" class="text-sm text-destructive">{{ form.errors.title }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="slug">Slug</Label>
                  <Input id="slug" v-model="form.slug" @input="onSlugInput" />
                  <p v-if="form.errors.slug" class="text-sm text-destructive">{{ form.errors.slug }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="clientName">Client Name</Label>
                  <Input id="clientName" v-model="form.clientName" />
                  <p v-if="form.errors.clientName" class="text-sm text-destructive">{{ form.errors.clientName }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="clientIndustry">Client Industry</Label>
                  <Input id="clientIndustry" v-model="form.clientIndustry" />
                  <p v-if="form.errors.clientIndustry" class="text-sm text-destructive">{{ form.errors.clientIndustry }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="clientUrl">Client URL</Label>
                  <Input id="clientUrl" v-model="form.clientUrl" />
                  <p v-if="form.errors.clientUrl" class="text-sm text-destructive">{{ form.errors.clientUrl }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="clientLogoUrl">Client Logo URL</Label>
                  <Input id="clientLogoUrl" v-model="form.clientLogoUrl" />
                  <p v-if="form.errors.clientLogoUrl" class="text-sm text-destructive">{{ form.errors.clientLogoUrl }}</p>
                </div>
              </div>

              <div class="mt-6 space-y-1">
                <Label for="summary">Summary</Label>
                <Textarea id="summary" v-model="form.summary" />
                <p v-if="form.errors.summary" class="text-sm text-destructive">{{ form.errors.summary }}</p>
              </div>
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="publishing">
            <AccordionTrigger>Publishing</AccordionTrigger>
            <AccordionContent>
              <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
                <div class="space-y-1">
                  <Label for="status">Status</Label>
                  <Select id="status" v-model="form.status" :options="statusOptions" />
                  <p v-if="form.errors.status" class="text-sm text-destructive">{{ form.errors.status }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="startedAt">Started Date</Label>
                  <Input id="startedAt" type="date" v-model="form.startedAt" />
                  <p v-if="form.errors.startedAt" class="text-sm text-destructive">{{ form.errors.startedAt }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="completedAt">Completed Date</Label>
                  <Input id="completedAt" type="date" v-model="form.completedAt" />
                  <p v-if="form.errors.completedAt" class="text-sm text-destructive">{{ form.errors.completedAt }}</p>
                </div>

                <div class="space-y-2 self-end pb-2">
                  <div class="flex items-center gap-2">
                    <Checkbox
                      id="isFeatured"
                      :model-value="form.isFeatured"
                      @update:model-value="setFeatured"
                    />
                    <Label for="isFeatured">Featured work</Label>
                  </div>
                  <p v-if="form.errors.isFeatured" class="text-sm text-destructive">{{ form.errors.isFeatured }}</p>
                </div>

                <div class="space-y-1 md:col-span-2">
                  <Label for="coverImageUrl">Cover Image URL</Label>
                  <Input id="coverImageUrl" v-model="form.coverImageUrl" />
                  <p v-if="form.errors.coverImageUrl" class="text-sm text-destructive">{{ form.errors.coverImageUrl }}</p>
                </div>
              </div>
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="categorisation">
            <AccordionTrigger>Categorisation</AccordionTrigger>
            <AccordionContent>
              <div class="space-y-4">
                <div class="space-y-1">
                  <Label for="specialisms">Specialisms</Label>
                  <Input id="specialisms" v-model="form.specialisms" placeholder="Strategy, UX, Engineering" />
                  <p v-if="form.errors.specialisms" class="text-sm text-destructive">{{ form.errors.specialisms }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="platforms">Platforms</Label>
                  <Input id="platforms" v-model="form.platforms" placeholder="Web, iOS, Android" />
                  <p v-if="form.errors.platforms" class="text-sm text-destructive">{{ form.errors.platforms }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="technologies">Technologies</Label>
                  <Input id="technologies" v-model="form.technologies" placeholder="Go, Vue, PostgreSQL" />
                  <p v-if="form.errors.technologies" class="text-sm text-destructive">{{ form.errors.technologies }}</p>
                </div>
              </div>
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="case-study">
            <AccordionTrigger>Case study</AccordionTrigger>
            <AccordionContent>
              <div class="space-y-4">
                <div class="space-y-1">
                  <Label for="challenge">Challenge</Label>
                  <Textarea id="challenge" v-model="form.challenge" />
                  <p v-if="form.errors.challenge" class="text-sm text-destructive">{{ form.errors.challenge }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="approach">Approach</Label>
                  <Textarea id="approach" v-model="form.approach" />
                  <p v-if="form.errors.approach" class="text-sm text-destructive">{{ form.errors.approach }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="deliverables">Deliverables</Label>
                  <Textarea id="deliverables" v-model="form.deliverables" />
                  <p v-if="form.errors.deliverables" class="text-sm text-destructive">{{ form.errors.deliverables }}</p>
                </div>

                <div class="space-y-1">
                  <Label for="outcome">Outcome</Label>
                  <Textarea id="outcome" v-model="form.outcome" />
                  <p v-if="form.errors.outcome" class="text-sm text-destructive">{{ form.errors.outcome }}</p>
                </div>
              </div>
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="content">
            <AccordionTrigger>Content</AccordionTrigger>
            <AccordionContent>
              <div class="space-y-2">
                <RichEditor id="content" v-model="form.content" placeholder="Write the work body..." class="min-h-[32rem]" />
                <p v-if="form.errors.content" class="text-sm text-destructive">{{ form.errors.content }}</p>
              </div>
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </div>

      <div class="flex justify-end gap-2 border-t px-6 py-4">
        <Link :href="`/admin/works/${item.ID}`">
          <Button type="button" variant="outline">Cancel</Button>
        </Link>
        <Button type="submit" :disabled="form.processing">
          {{ form.processing ? 'Updating...' : 'Update Work' }}
        </Button>
      </div>
    </form>
  </div>
</template>
