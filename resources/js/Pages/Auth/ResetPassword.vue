<script setup lang="ts">
import { Head, useForm } from '@inertiajs/vue3'

import AuthLayout from '@/Layouts/AuthLayout.vue'
import { routes } from '@/routes'

const props = defineProps<{ token: string }>()

const form = useForm({
  resetPasswordToken: props.token,
  password: '',
  confirmPassword: '',
})

function submit() {
  form.put(routes.passwordUpdate())
}
</script>

<template>
  <Head title="Reset Password" />

  <AuthLayout
    eyebrow="Account recovery / New password"
    title="Reset Your Password."
    description="Enter your new password below."
    meta-title="Recovery step 02"
    meta-description="New password / Confirmation"
    form-label="New credentials / 01-02"
  >
    <form @submit.prevent="submit">
      <fieldset class="border-0 p-0" :disabled="form.processing">
        <div class="border-b border-border">
          <label class="block px-5 pt-4 font-mono text-xs uppercase tracking-[0.12em]" for="password">01 / New Password</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            autocomplete="new-password"
            minlength="8"
            class="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60"
            required
          >
          <p v-if="form.errors.password" class="px-5 pb-3 text-sm">{{ form.errors.password }}</p>
        </div>
        <div class="border-b border-border">
          <label class="block px-5 pt-4 font-mono text-xs uppercase tracking-[0.12em]" for="confirmPassword">02 / Confirm New Password</label>
          <input
            id="confirmPassword"
            v-model="form.confirmPassword"
            type="password"
            autocomplete="new-password"
            minlength="8"
            class="h-12 w-full border-0 bg-transparent px-5 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60"
            required
          >
          <p v-if="form.errors.confirmPassword" class="px-5 pb-3 text-sm">{{ form.errors.confirmPassword }}</p>
        </div>
        <button
          type="submit"
          class="inline-flex min-h-14 w-full items-center justify-center border-t border-primary bg-primary px-5 py-3 text-xs font-semibold uppercase tracking-[0.14em] text-primary-foreground transition-colors hover:bg-transparent hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="form.processing"
        >
          {{ form.processing ? 'Loading' : 'Reset Password' }}
        </button>
      </fieldset>
    </form>
  </AuthLayout>
</template>
