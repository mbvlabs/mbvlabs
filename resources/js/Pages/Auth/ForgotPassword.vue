<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'

import AuthLayout from '@/Layouts/AuthLayout.vue'
import { routes } from '@/routes'

const form = useForm({ email: '' })

function submit() {
  form.post(routes.passwordCreate())
}
</script>

<template>
  <Head title="Reset Password" />

  <AuthLayout
    eyebrow="Account recovery / Request"
    title="Reset Password."
    description="Enter your email address and we'll send you a link to reset your password."
    meta-title="Recovery step 01"
    meta-description="Email / Reset link"
    form-label="Recovery email / 01"
  >
    <form @submit.prevent="submit">
      <fieldset class="border-0 p-0" :disabled="form.processing">
        <div class="border-b border-border">
          <label class="block px-5 pt-4 font-mono text-xs uppercase tracking-[0.12em]" for="email">01 / Email</label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            autocomplete="email"
            class="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60"
            required
          >
          <p v-if="form.errors.email" class="px-5 pb-3 text-sm">{{ form.errors.email }}</p>
        </div>
        <button
          type="submit"
          class="inline-flex min-h-14 w-full items-center justify-center border-t border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="form.processing"
        >
          {{ form.processing ? 'Loading' : 'Send Reset Link' }}
        </button>
      </fieldset>
    </form>
    <p class="border-t border-border px-5 py-5 text-center text-sm">
      Remember your password?
      <Link class="font-semibold underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" :href="routes.sessionNew()">Login</Link>
    </p>
  </AuthLayout>
</template>
