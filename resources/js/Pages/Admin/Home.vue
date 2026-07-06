<script setup lang="ts">
import { Head } from '@inertiajs/vue3'
import {
  ArrowUpRight,
  ArrowDownRight,
  Users,
  DollarSign,
  Activity,
  CreditCard,
} from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'

defineOptions({ layout: AdminLayout })

interface Stat {
  title: string
  value: string
  change: string
  trend: 'up' | 'down'
  icon: typeof Users
  description: string
}

const stats: Stat[] = [
  {
    title: 'Total Revenue',
    value: '$45,231.89',
    change: '+20.1%',
    trend: 'up',
    icon: DollarSign,
    description: 'from last month',
  },
  {
    title: 'Subscriptions',
    value: '+2,350',
    change: '+180.1%',
    trend: 'up',
    icon: Users,
    description: 'from last month',
  },
  {
    title: 'Sales',
    value: '+12,234',
    change: '+19%',
    trend: 'up',
    icon: CreditCard,
    description: 'from last month',
  },
  {
    title: 'Active Now',
    value: '+573',
    change: '+201',
    trend: 'up',
    icon: Activity,
    description: 'since last hour',
  },
]

const recentActivity = [
  { user: 'Olivia Martin', email: 'olivia.martin@email.com', amount: '+$1,999.00' },
  { user: 'Jackson Lee', email: 'jackson.lee@email.com', amount: '+$39.00' },
  { user: 'Isabella Nguyen', email: 'isabella.nguyen@email.com', amount: '+$299.00' },
  { user: 'William Kim', email: 'will@email.com', amount: '+$99.00' },
  { user: 'Sofia Davis', email: 'sofia.davis@email.com', amount: '+$39.00' },
]
</script>

<template>
  <Head title="Admin" />

    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Admin</h1>
        <p class="text-muted-foreground">Welcome back — here is what is happening today.</p>
      </div>
    </div>

    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <div
        v-for="stat in stats"
        :key="stat.title"
        class="rounded-none border bg-card p-6 text-card-foreground shadow-sm"
      >
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-muted-foreground">{{ stat.title }}</span>
          <component :is="stat.icon" class="size-4 text-muted-foreground" />
        </div>
        <div class="mt-3 text-2xl font-bold tracking-tight">{{ stat.value }}</div>
        <div class="mt-1 flex items-center text-xs">
          <span
            :class="
              stat.trend === 'up'
                ? 'text-emerald-600 dark:text-emerald-400'
                : 'text-red-600 dark:text-red-400'
            "
            class="flex items-center font-medium"
          >
            <ArrowUpRight v-if="stat.trend === 'up'" class="mr-1 size-3" />
            <ArrowDownRight v-else class="mr-1 size-3" />
            {{ stat.change }}
          </span>
          <span class="ml-2 text-muted-foreground">{{ stat.description }}</span>
        </div>
      </div>
    </div>

    <div class="grid gap-4 lg:grid-cols-7">
      <div class="rounded-none border bg-card text-card-foreground shadow-sm lg:col-span-4">
        <div class="border-b p-6">
          <h3 class="font-semibold leading-none tracking-tight">Overview</h3>
          <p class="mt-1 text-sm text-muted-foreground">Monthly revenue for the current year.</p>
        </div>
        <div class="p-6">
          <div class="flex h-[240px] items-end gap-2">
            <div
              v-for="(height, i) in [40, 65, 45, 80, 55, 90, 70, 85, 60, 75, 50, 95]"
              :key="i"
              class="flex-1 rounded-none bg-primary/80 transition-all hover:bg-primary"
              :style="{ height: `${height}%` }"
            />
          </div>
          <div class="mt-4 flex justify-between text-xs text-muted-foreground">
            <span>Jan</span>
            <span>Dec</span>
          </div>
        </div>
      </div>

      <div class="rounded-none border bg-card text-card-foreground shadow-sm lg:col-span-3">
        <div class="border-b p-6">
          <h3 class="font-semibold leading-none tracking-tight">Recent Sales</h3>
          <p class="mt-1 text-sm text-muted-foreground">You made 265 sales this month.</p>
        </div>
        <div class="p-6">
          <div class="space-y-6">
            <div v-for="activity in recentActivity" :key="activity.email" class="flex items-center">
              <div
                class="flex size-9 items-center justify-center rounded-none bg-muted text-xs font-semibold uppercase"
              >
                {{ activity.user.charAt(0) }}{{ activity.user.split(' ')[1]?.charAt(0) }}
              </div>
              <div class="ml-4 flex-1">
                <p class="text-sm font-medium">{{ activity.user }}</p>
                <p class="text-xs text-muted-foreground">{{ activity.email }}</p>
              </div>
              <div class="text-sm font-medium">{{ activity.amount }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
</template>
