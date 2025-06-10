<script setup lang="ts">
import ProgressTimelineItem from "@/components/ProgressTimelineItem.vue"
import TimelineCard from "@/types.d/TimelineCard"

const items = [
  {
    id: "1",
    title: "Order Shipped",
    description:
      "Pretium lectus quam id leo. Urna et pharetra aliquam vestibulum morbi blandit cursus risus.",
    date: "09/06/2023",
    state: "complete",
  },
  {
    id: "2",
    title: "In Transit",
    description:
      "Pretium lectus quam id leo. Urna et pharetra a" +
      "Pretium lectus quam id leo. Urna et pharetra a" +
      "Pretium lectus quam id leo. Urna et pharetra a" +
      "Pretium lectus quam id leo. Urna et pharetra a" +
      "Pretium lectus quam id leo. Urna et pharetra a" +
      "Pretium lectus quam id leo. Urna et pharetra a" +
      "Pretium lectus quam id leo. Urna et pharetra a" +
      "Pretium lectus quam id leo. Urna et pharetra a" +
      "Pretium lectus quam id leo. Urna et pharetra a" +
      "Pretium lectus quam id leo. Urna et pharetra a" +
      "liquam vestibulum morbi blandit cursus risus.",
    date: "10/06/2023",
    state: "complete",
  },
  {
    id: "3",
    title: "Out of Delivery",
    description:
      "Pretium lectus quam id leo. " +
      "내용은 천차만별 lectus quam id leo. " +
      "Urna et pharetra aliquam vestibulum morbi blandit cursus risus.",
    date: "12/06/2023",
    state: "incomplete",
  },
  {
    id: "4",
    title: "Delivered",
    description:
      "Pretium lectus quam id leo. Urna et pharetra aliquam vestibulum morbi blandit cursus risus.",
    date: "12/08/2023",
    state: "incomplete",
  },
] as TimelineCard[]
</script>

<template>
  <div class="flex flex-col">
    <ProgressTimelineItem
      v-for="item in items"
      :key="item.id"
      class="pb-2 md:pb-4 last-of-type:pb-0"
      :item="item"
    />
  </div>
</template>
