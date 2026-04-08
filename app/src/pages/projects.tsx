import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group"
import {
  Item,
  ItemActions,
  ItemContent,
  ItemGroup,
  ItemMedia,
  ItemTitle,
} from "@/components/ui/item"
import { Separator } from "@/components/ui/separator"
import { Link, useSearchParams } from "react-router"
import { useEffect, useState } from "react"
import { Switch } from "@/components/ui/switch"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Archive, Ellipsis, Hash, Plus, Search, Trash } from "lucide-react"
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
} from "@/components/ui/pagination"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Field, FieldGroup } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useQuery } from "@tanstack/react-query"
import {
  getProjects,
  PROJECTS_KEY,
  type GetProjectsParams,
} from "@/api/resources"
import { LoadingPage } from "@/components/loading-page"
import { ErrorPage } from "@/components/error-page"
import { AppError } from "@/lib/app-error"
import useDebounce from "@/hooks/use-debounce"

export default function Projects() {
  const [alertArchiveOpen, setAlertArchiveOpen] = useState(false)
  const [alertDeleteOpen, setAlertDeleteOpen] = useState(false)
  const [uRLSearchParams, setURLSearchParams] = useSearchParams()
  const [searchInput, setSearchInput] = useState<string>(uRLSearchParams.get("search") ?? "")
  const [withArchive, setWithArchive] = useState<boolean>(
    () => uRLSearchParams.get("archive") === "true"
  )

  const debouncedSearch = useDebounce<string>(searchInput, 500)

  const queryOpt: GetProjectsParams = {
    search: debouncedSearch,
    withArchive: withArchive
  }
  const { data, isPending, isError, error } = useQuery({
    queryKey: [PROJECTS_KEY, queryOpt],
    queryFn: () => getProjects(queryOpt),
  })

  const handleSearchChangeInput = (val: string) => {
    setSearchInput(val)
  }

  const handleWithArchive = (checked: boolean) => {
    setWithArchive(checked)
  }

  useEffect(() => {
    setURLSearchParams((prev) => {
      const params = new URLSearchParams(prev)

      if (!debouncedSearch) {
        params.delete("search")
      } else {
        params.set("search", debouncedSearch)
      }

      if (withArchive) {
        params.set("archive", "true")
      } else {
        params.delete("archive")
      }

      return params
    })
  }, [debouncedSearch, setURLSearchParams, withArchive])

  if (isPending) {
    return (
      <div className="flex h-screen flex-col justify-center">
        <LoadingPage title="Loading projects..." />
      </div>
    )
  }

  if (isError) {
    const appError =
      error instanceof AppError
        ? error
        : new AppError({ status: 500, message: "Something when wrong" })

    return (
      <div className="flex h-screen flex-col justify-center">
        <ErrorPage
          status={appError.status}
          message={appError.message}
          description={appError?.description}
        />
      </div>
    )
  }

  return (
    <>
      <div>
        <div className="mx-24 my-4">
          <div className="flex flex-col justify-between space-y-4">
            <h2 className="text-2xl font-semibold">My Projects</h2>
            <div className="space-y-4">
              <div>
                <InputGroup>
                  <InputGroupInput autoComplete="false" placeholder="Search..." value={searchInput} onChange={(e) => handleSearchChangeInput(e.target.value)} />
                  <InputGroupAddon>
                    <Search />
                  </InputGroupAddon>
                </InputGroup>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-x-4">
                  <span className="text-muted-foreground">Archive only</span>
                  <Switch className={"cursor-pointer"} checked={withArchive} onCheckedChange={handleWithArchive} />
                </div>
                <Dialog>
                  <form>
                    <DialogTrigger>
                      <div className="flex cursor-pointer items-center gap-x-2 rounded-md bg-primary p-1 text-primary-foreground">
                        <Plus />
                      </div>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-sm">
                      <DialogHeader>
                        <DialogTitle>Add project</DialogTitle>
                      </DialogHeader>
                      <FieldGroup>
                        <Field>
                          <Label htmlFor="name-1">Name</Label>
                          <Input
                            id="name-1"
                            name="name"
                            defaultValue="Learn Go"
                          />
                        </Field>
                      </FieldGroup>
                      <DialogFooter>
                        <DialogClose
                          render={<Button variant="outline">Cancel</Button>}
                        />
                        <Button type="submit">Save changes</Button>
                      </DialogFooter>
                    </DialogContent>
                  </form>
                </Dialog>
              </div>
            </div>
            <div>
              <h4>
                {data.pagination
                  ? `${data.pagination?.total_items_in_page} projects`
                  : "0 project"}
              </h4>
              <Separator />
              {data.projects.length !== 0 && data.pagination !== null && (
                <>
                  <div className="my-3">
                    <ItemGroup>
                      {data.projects.map((project) => (
                        <Item key={project.id}>
                          <ItemMedia variant="icon">
                            <Hash />
                          </ItemMedia>
                          <ItemContent>
                            <Link to={project.id}>
                              <ItemTitle>{project.name}</ItemTitle>
                            </Link>
                          </ItemContent>
                          <ItemActions>
                            <DropdownMenu>
                              <DropdownMenuTrigger
                                render={
                                  <Button
                                    variant={"ghost"}
                                    className={"cursor-pointer"}
                                  >
                                    <Ellipsis />
                                  </Button>
                                }
                              />
                              <DropdownMenuContent>
                                <DropdownMenuGroup>
                                  <DropdownMenuItem
                                    className={"cursor-pointer"}
                                    onClick={() => setAlertArchiveOpen(true)}
                                  >
                                    <Archive />
                                    Archive
                                  </DropdownMenuItem>
                                  <DropdownMenuItem
                                    className="cursor-pointer text-red-500"
                                    onClick={() => setAlertDeleteOpen(true)}
                                  >
                                    <Trash />
                                    Delete
                                  </DropdownMenuItem>
                                </DropdownMenuGroup>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </ItemActions>
                        </Item>
                      ))}
                    </ItemGroup>
                  </div>
                  <div className="self-end">
                    <Pagination>
                      <PaginationContent>
                        {Array.from(
                          { length: data.pagination.total_pages },
                          (_, i) => (
                            <PaginationItem key={i}>
                              <PaginationLink href={`?page=${i + 1}`}>
                                {i + 1}
                              </PaginationLink>
                            </PaginationItem>
                          )
                        )}
                      </PaginationContent>
                    </Pagination>
                  </div>
                  {/*ARCHIVE ALERT*/}
                  <AlertDialog
                    open={alertArchiveOpen}
                    onOpenChange={setAlertArchiveOpen}
                  >
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>
                          Are you absolutely sure?
                        </AlertDialogTitle>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction>Continue</AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>

                  {/*DELETE ALERT*/}
                  <AlertDialog
                    open={alertDeleteOpen}
                    onOpenChange={setAlertDeleteOpen}
                  >
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>
                          Are you absolutely sure?
                        </AlertDialogTitle>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction>Continue</AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
