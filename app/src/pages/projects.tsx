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
import { useEffect, useRef, useState, type SubmitEvent } from "react"
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
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import {
  deleteProject,
  getProjects,
  postProject,
  PROJECTS_KEY,
  updateProject,
  type GetProjectsParams,
} from "@/api/resources"
import { LoadingPage } from "@/components/loading-page"
import { ErrorPage } from "@/components/error-page"
import { AppError } from "@/lib/app-error"
import useDebounce from "@/hooks/use-debounce"
import { toast } from "sonner"

export default function Projects() {
  const [alertArchiveOpen, setAlertArchiveOpen] = useState(false)
  const [alertDeleteOpen, setAlertDeleteOpen] = useState(false)
  const [uRLSearchParams, setURLSearchParams] = useSearchParams()
  const [searchInput, setSearchInput] = useState<string>(
    uRLSearchParams.get("search") ?? ""
  )
  const [withArchive, setWithArchive] = useState<boolean>(
    () => uRLSearchParams.get("archive") === "true"
  )
  const [page, setPage] = useState<string>(
    uRLSearchParams.get("page") ?? "1"
  )
  const headingRef = useRef<HTMLHeadingElement | null>(null)
  const prevSearchRef = useRef<string>("")
  const prevArchiveRef = useRef<boolean>(withArchive)

  const queryClient = useQueryClient()
  const [projectName, setProjectName] = useState("")
  const [openAddDialog, setOpenAddDialog] = useState(false)
  const [deleteProjectId, setDeleteProjectId] = useState("")
  const [updateProjectId, setUpdateProjectId] = useState("")

  const createProjectMutation = useMutation({
    mutationFn: postProject,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PROJECTS_KEY] })
      setProjectName("")
      setOpenAddDialog(false)
    },
  })

  const deleteProjectMutation = useMutation({
    mutationFn: deleteProject,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PROJECTS_KEY] })
      setAlertDeleteOpen(false)
      toast.success("Delete project succesfully")
    },
  })

  const updateProjectMutation = useMutation({
    mutationFn: updateProject,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [PROJECTS_KEY] })
      setAlertArchiveOpen(false)
      toast.success("Update project succesfully")
    },
    onError: () => {
      setAlertArchiveOpen(false)
      toast.error("Failed update project")
    }
  })

  const debouncedSearch = useDebounce<string>(searchInput, 500)

  const queryOpt: GetProjectsParams = {
    search: debouncedSearch,
    withArchive: withArchive,
    page,
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

  const handlePostProjectSubmit = (e: SubmitEvent<HTMLFormElement>) => {
    e.preventDefault()

    createProjectMutation.mutate({
      name: projectName,
    })
  }

  const handleDeleteProject = (id: string) => {
    deleteProjectMutation.mutate(id)
  }

  const handleUpdateProjectToBeArchived = (id: string) => {
    updateProjectMutation.mutate({ id, body: { to_be_archived: true } })
  }

  const handlePage = (val: number) => {
    setPage(String(val))

    headingRef.current?.scrollIntoView({
      behavior: "smooth",
      block: "start",
    })
  }

  useEffect(() => {
    setURLSearchParams((prev) => {
      const params = new URLSearchParams(prev)

      const prevSearch = prevSearchRef.current
      const prevArchive = prevArchiveRef.current

      const isNewSearch = debouncedSearch !== prevSearch
      const isArchiveChanged = withArchive !== prevArchive

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

      if ((isNewSearch && debouncedSearch) || isArchiveChanged) {
        params.delete("page")
        handlePage(1)
      } else {
        if (!page || page === "1") {
          params.delete("page")
        } else {
          params.set("page", page)
        }
      }

      prevSearchRef.current = debouncedSearch
      prevArchiveRef.current = withArchive

      return params
    })
  }, [debouncedSearch, setURLSearchParams, withArchive, page])

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
            <h2 ref={headingRef} className="text-2xl font-semibold">
              My Projects
            </h2>
            <div className="space-y-4">
              <div>
                <InputGroup>
                  <InputGroupInput
                    autoComplete="false"
                    placeholder="Search..."
                    value={searchInput}
                    onChange={(e) =>
                      handleSearchChangeInput(e.target.value)
                    }
                  />
                  <InputGroupAddon>
                    <Search />
                  </InputGroupAddon>
                </InputGroup>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-x-4">
                  <span className="text-muted-foreground">
                    Archive only
                  </span>
                  <Switch
                    className={"cursor-pointer"}
                    checked={withArchive}
                    onCheckedChange={handleWithArchive}
                  />
                </div>
                <Dialog
                  open={openAddDialog}
                  onOpenChange={(next) => {
                    if (createProjectMutation.isPending) return

                    if (!next) {
                      createProjectMutation.reset()
                      setProjectName("")
                    }

                    setOpenAddDialog(next)
                  }}
                >
                  <DialogTrigger>
                    <div className="flex cursor-pointer items-center gap-x-2 rounded-md bg-primary p-1 text-primary-foreground">
                      <Plus />
                    </div>
                  </DialogTrigger>
                  <DialogContent className="sm:max-w-sm">
                    <form onSubmit={handlePostProjectSubmit}>
                      <DialogHeader>
                        <DialogTitle>Add project</DialogTitle>
                      </DialogHeader>
                      <FieldGroup
                        className={
                          createProjectMutation.isError
                            ? "gap-2 mt-2"
                            : "gap-5"
                        }
                      >
                        {createProjectMutation.isError && (
                          <p className="text-sm text-red-500">
                            {
                              (createProjectMutation.error as AppError)
                                .message
                            }
                          </p>
                        )}
                        <Field>
                          <Input
                            className={
                              createProjectMutation.isError
                                ? "mt-0"
                                : "mt-4"
                            }
                            name="name"
                            autoFocus
                            value={projectName}
                            onChange={(e) => {
                              setProjectName(e.target.value)
                              if (createProjectMutation.isError) {
                                createProjectMutation.reset()
                              }
                            }}
                          />
                        </Field>
                      </FieldGroup>
                      <DialogFooter>
                        <DialogClose
                          render={
                            <Button variant="outline">
                              Cancel
                            </Button>
                          }
                        />
                        <Button
                          type="submit"
                          disabled={createProjectMutation.isPending}
                        >
                          {createProjectMutation.isPending
                            ? "Saving..."
                            : "Save changes"}
                        </Button>
                      </DialogFooter>
                    </form>
                  </DialogContent>
                </Dialog>
              </div>
            </div>
            <div>
              <h4>
                {`${data.pagination?.total_items_in_page ?? 0} projects`}
              </h4>
              <Separator />
              {data.projects.length !== 0 &&
                data.pagination !== null && (
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
                                <ItemTitle>
                                  {project.name}
                                </ItemTitle>
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
                                      onClick={() => {
                                        setAlertArchiveOpen(true)
                                        setUpdateProjectId(project.id)
                                      }}
                                    >
                                      <Archive />
                                      Archive
                                    </DropdownMenuItem>
                                    <DropdownMenuItem
                                      className="cursor-pointer text-red-500"
                                      onClick={() => {
                                        setAlertDeleteOpen(true)
                                        setDeleteProjectId(project.id)
                                      }}
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
                            {
                              length:
                                data.pagination?.total_pages ?? 1,
                            },
                            (_, i) => (
                              <PaginationItem key={i}>
                                <Button
                                  onClick={() =>
                                    handlePage(i + 1)
                                  }
                                  variant={"ghost"}
                                  className="cursor-pointer inline-flex items-center justify-center rounded-md text-sm font-medium size-9 hover:bg-muted hover:text-foreground"
                                >
                                  {i + 1}
                                </Button>
                              </PaginationItem>
                            )
                          )}
                        </PaginationContent>
                      </Pagination>
                    </div>

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
                          <AlertDialogCancel>
                            Cancel
                          </AlertDialogCancel>
                          <AlertDialogAction onClick={() => handleUpdateProjectToBeArchived(updateProjectId)}>
                            Continue
                          </AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>

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
                          <AlertDialogCancel>
                            Cancel
                          </AlertDialogCancel>
                          <AlertDialogAction onClick={() => handleDeleteProject(deleteProjectId)}>
                            Continue
                          </AlertDialogAction>
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
