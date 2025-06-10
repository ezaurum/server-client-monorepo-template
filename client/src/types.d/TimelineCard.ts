import { HasID } from "@/types.d/HasID"

export default interface TimelineCard extends HasID {
  title: string
  description: string
  date: string
  state: "incomplete" | "complete"
}
