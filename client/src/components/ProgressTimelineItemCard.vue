<template>
  <div class="bg-white p-4 rounded border border-slate-200 shadow">
    <div class="flex items-center justify-between space-x-2 mb-1">
      <div class="font-bold text-slate-900">
        <slot name="title" />
      </div>
      <time class="font-caveat font-medium text-indigo-500">
        <slot name="date" />
      </time>
    </div>
    <div class="text-slate-500">
      <slot />
    </div>
  </div>
</template>
