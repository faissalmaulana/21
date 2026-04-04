import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyTitle,
} from "@/components/ui/empty"

interface ErrorPageProps {
  status: number
  message: string
  description?: string
}

export function ErrorPage(props: ErrorPageProps) {
  return (
    <Empty className="w-full">
      <EmptyHeader>
        <EmptyTitle>
          {props.status} - {props.message}
        </EmptyTitle>
        {props.description && (
          <EmptyDescription>{props.description}</EmptyDescription>
        )}
      </EmptyHeader>
    </Empty>
  )
}
