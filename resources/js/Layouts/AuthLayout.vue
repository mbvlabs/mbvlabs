<script setup lang="ts">
import { Link } from '@inertiajs/vue3'

import { routes } from '@/routes'

defineProps<{
  eyebrow: string
  title: string
  description: string
  metaTitle: string
  metaDescription: string
  formLabel: string
}>()
</script>

<template>
  <div class="cobalt-grid flex min-h-screen flex-col bg-background font-sans text-foreground">
    <header class="border-b border-foreground">
      <div class="container mx-auto flex items-center justify-between px-6 py-5 md:px-8">
        <Link :href="routes.homePage()" class="font-mono text-xs font-medium uppercase tracking-[0.18em]">
          MBV Labs
        </Link>
        <p class="font-mono text-[11px] uppercase tracking-[0.12em] opacity-65">Authorized access only</p>
      </div>
    </header>

    <main class="flex flex-1">
      <section class="flex w-full items-center border-b border-foreground">
        <div class="container mx-auto grid gap-14 px-6 py-20 md:px-8 md:py-28 lg:grid-cols-[1.1fr_0.9fr] lg:items-center lg:gap-24">
          <div>
            <p class="font-mono text-xs uppercase tracking-[0.18em]">{{ eyebrow }}</p>
            <h1 class="cobalt-display mt-8 max-w-4xl">{{ title }}</h1>
            <p class="mt-8 max-w-xl text-base leading-7">{{ description }}</p>
            <div class="mt-12 max-w-xl border-t border-foreground pt-5 font-mono text-[11px] uppercase tracking-[0.12em]">
              <p class="font-medium">{{ metaTitle }}</p>
              <p class="mt-1 opacity-65">{{ metaDescription }}</p>
            </div>
          </div>

          <div class="w-full border border-foreground">
            <div class="border-b border-foreground px-6 py-5">
              <p class="font-mono text-xs uppercase tracking-[0.14em]">{{ formLabel }}</p>
            </div>
            <slot />
          </div>
        </div>
      </section>
    </main>
  </div>
</template>
