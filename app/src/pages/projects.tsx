import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuGroup, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { InputGroup, InputGroupAddon, InputGroupInput } from "@/components/ui/input-group";
import { Item, ItemActions, ItemContent, ItemGroup, ItemMedia, ItemTitle } from "@/components/ui/item";
import { Separator } from "@/components/ui/separator";
import { Link } from "react-router";
import { useState } from "react";
import { Switch } from "@/components/ui/switch";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import {
  Archive,
  Ellipsis,
  FolderDot,
  Hash,
  Plus,
  Search,
  Trash
} from 'lucide-react';
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
import { Field, FieldGroup } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useQuery } from "@tanstack/react-query";
import { getProjects, PROJECTS_KEY } from "@/api/resources";

import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/components/ui/empty"


export default function Projects() {
  const { data, isPending, isError, error } = useQuery({
    queryKey: [PROJECTS_KEY],
    queryFn: getProjects
  })


  const [alertArchiveOpen, setAlertArchiveOpen] = useState(false);
  const [alertDeleteOpen, setAlertDeleteOpen] = useState(false);

  if (isPending) {
    return <span>Loading...</span>
  }

  if (isError) {
    return <span>Error: {error.message}</span>
  }

  if (data.projects.length === 0 || data.pagination === null) {
    return (
      <div>
        <EmptyProjects />
      </div>
    )
  }

  return (
    <div>
      <div className="my-4 mx-24">
        <div className="flex flex-col space-y-4 justify-between">
          <h2 className="text-2xl font-semibold">My Projects</h2>
          <div className="space-y-4">
            <div>
              <InputGroup>
                <InputGroupInput placeholder="Search..." />
                <InputGroupAddon>
                  <Search />
                </InputGroupAddon>
              </InputGroup>
            </div>
            <div className="flex justify-between items-center">
              <div className="flex items-center gap-x-4">
                <span className="text-muted-foreground">Archive only</span>
                <Switch className={"cursor-pointer"} />
              </div>
              <CreateNewProjectForm>
                <div className="flex items-center gap-x-2 bg-primary text-primary-foreground p-1 rounded-md cursor-pointer">
                  <Plus />
                </div>
              </CreateNewProjectForm>
            </div>
          </div>

          <div>
            <h4>{data.pagination.total_items_in_page} projects</h4>
            <Separator />
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
                      {/*ARCHIVE ALERT*/}
                      <AlertDialog open={alertArchiveOpen} onOpenChange={setAlertArchiveOpen}>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction>Continue</AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>

                      {/*DELETE ALERT*/}
                      <AlertDialog open={alertDeleteOpen} onOpenChange={setAlertDeleteOpen}>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction>Continue</AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>

                      <DropdownMenu>
                        <DropdownMenuTrigger render={<Button variant={"ghost"} className={"cursor-pointer"}><Ellipsis /></Button>} />
                        <DropdownMenuContent>
                          <DropdownMenuGroup>
                            <DropdownMenuItem className={"cursor-pointer"} onClick={() => setAlertArchiveOpen(true)}>
                              <Archive />
                              Archive
                            </DropdownMenuItem>

                            <DropdownMenuItem className="text-red-500 cursor-pointer" onClick={() => setAlertDeleteOpen(true)}>
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
                  {Array.from({ length: data.pagination.total_pages }, (_, i) => (
                    <PaginationItem key={i}>
                      <PaginationLink href={`?page=${i + 1}`}>
                        {i + 1}
                      </PaginationLink>
                    </PaginationItem>
                  ))}
                </PaginationContent>
              </Pagination>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}



export function EmptyProjects() {
  return (
    <Empty>
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <FolderDot />
        </EmptyMedia>
        <EmptyTitle>No Projects Yet</EmptyTitle>
        <EmptyDescription>
          You haven&apos;t created any projects yet. Get started by creating
          your first project.
        </EmptyDescription>
      </EmptyHeader>
      <EmptyContent className="flex-row justify-center gap-2">
        <CreateNewProjectForm>
          <div className="flex items-center gap-x-2 bg-primary text-primary-foreground p-2 rounded-md cursor-pointer">Create Project</div>
        </CreateNewProjectForm>
      </EmptyContent>
    </Empty>
  )
}

interface CreateNewProjectFormProps {
  // this is for trigger dialog,
  // it should be not button
  children: React.ReactElement
}

export function CreateNewProjectForm(prop: CreateNewProjectFormProps) {
  return (
    <Dialog>
      <form>
        <DialogTrigger>
          {prop.children}
        </DialogTrigger>
        <DialogContent className="sm:max-w-sm">
          <DialogHeader>
            <DialogTitle>Add project</DialogTitle>
          </DialogHeader>
          <FieldGroup>
            <Field>
              <Label htmlFor="name-1">Name</Label>
              <Input id="name-1" name="name" defaultValue="Learn Go" />
            </Field>
          </FieldGroup>
          <DialogFooter>
            <DialogClose render={<Button variant="outline">Cancel</Button>} />
            <Button type="submit">Save changes</Button>
          </DialogFooter>
        </DialogContent>
      </form>
    </Dialog>
  )
}
