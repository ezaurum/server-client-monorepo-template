import type { Meta, StoryObj } from "@storybook/vue3"
import ProgressTimeline from "@/components/ProgressTimeline.vue"

const meta = {
  title: "Pages/ProgressTimeline",
  component: ProgressTimeline,
  render: () => ({
    components: { ProgressTimeline },
    template: "<ProgressTimeline />",
  }),
  parameters: {
    // More on how to position stories at: https://storybook.js.org/docs/vue/configure/story-layout
    layout: "fullscreen",
  },
  // This component will have an automatically generated docsPage entry: https://storybook.js.org/docs/vue/writing-docs/autodocs
  tags: ["autodocs"],
} satisfies Meta<typeof ProgressTimeline>

export default meta
type Story = StoryObj<typeof meta>

export const Progressing: Story = {}
