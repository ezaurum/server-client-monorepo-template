<script setup lang="ts">
import ProgressTimelineItemCard from "@/components/ProgressTimelineItemCard.vue"
import { computed, ref } from "vue"
import TimelineCard from "@/types.d/TimelineCard"

const props = withDefaults(
  defineProps<{
    item: TimelineCard
  }>(),
  {
    item: () => ({
      id: "",
      title: "Order Shipped",
      description:
        "Pretium lectus quam id leo. Urna et pharetra aliquam vestibulum morbi blandit cursus risus.",
      date: "09/06/2023",
      state: "incomplete",
    }),
  }
)

const title = ref(props.item.title)
const description = ref(props.item.description)
const date = ref(props.item.date)
const state = ref(props.item.state)
const isComplete = computed(() => state.value === "complete")
</script>

<template>
  <div
    class="relative flex items-center gap-1 md:gap-4 before:content-[attr(data-state)] before:block before:w-2 before:h-full before:bg-amber-300 before:absolute before:-z-10 before:left-4 before:top-1/2 first-of-type:before:top-1/2 before:last-of-type:before:h-1/2 last-of-type:before:bottom-1/2 max-w-full md:max-w-2xl"
    :class="isComplete ? 'is-active' : null"
  >
    <!-- Icon -->
    <div
      class="flex items-center justify-center w-10 h-10 rounded-full border border-white group-[.is-active]:bg-emerald-500 group-[.is-active]:text-emerald-50 text-slate-50 bg-slate-300 shadow shrink-0"
    >
      <svg
        class="fill-current"
        xmlns="http://www.w3.org/2000/svg"
        width="12"
        height="10"
      >
        <path
          fill-rule="nonzero"
          d="M10.422 1.257 4.655 7.025 2.553 4.923A.916.916 0 0 0 1.257 6.22l2.75 2.75a.916.916 0 0 0 1.296 0l6.415-6.416a.916.916 0 0 0-1.296-1.296Z"
        />
      </svg>
    </div>

    <!-- Card -->
    <ProgressTimelineItemCard>
      <template #title> {{ title }}</template>
      <template #default>
        {{ description }}
      </template>
      <template #date> {{ date }}</template>
    </ProgressTimelineItemCard>
  </div>
</template>
