import {
  Empty,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/components/ui/empty"
import { Spinner } from "@/components/ui/spinner"

export function LoadingPage(props: { title: string }) {
  return (
    <Empty className="w-full">
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <Spinner className="size-6" />
        </EmptyMedia>
        <EmptyTitle>{props.title}</EmptyTitle>
      </EmptyHeader>
    </Empty>
  )
}
