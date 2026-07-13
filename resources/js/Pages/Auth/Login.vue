<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'

import AuthLayout from '@/Layouts/AuthLayout.vue'
import { routes } from '@/routes'

const form = useForm({
  email: '',
  password: '',
})

function submit() {
  form.post(routes.sessionCreate())
}
</script>

<template>
  <Head title="Login" />

  <AuthLayout
    eyebrow="Account access / Private"
    title="Login to your account."
    description="Enter your details below to login to your account."
    meta-title="Authorized access only"
    meta-description="MBV Labs / Account gateway"
    form-label="Credentials / 01-02"
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
        <div class="border-b border-border">
          <label class="block px-5 pt-4 font-mono text-xs uppercase tracking-[0.12em]" for="password">02 / Password</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            autocomplete="current-password"
            class="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60"
            required
          >
          <p v-if="form.errors.password" class="px-5 pb-3 text-sm">{{ form.errors.password }}</p>
        </div>
        <div class="flex justify-end px-5 py-4">
          <Link class="text-sm font-semibold underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" :href="routes.passwordNew()">
            Forgot your password?
          </Link>
        </div>
        <button
          type="submit"
          class="inline-flex min-h-14 w-full items-center justify-center border-t border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="form.processing"
        >
          {{ form.processing ? 'Loading' : 'Login' }}
        </button>
      </fieldset>
    </form>
  </AuthLayout>
</template>
