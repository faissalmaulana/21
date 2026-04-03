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


export default function Projects() {
  const [alertArchiveOpen, setAlertArchiveOpen] = useState(false);
  const [alertDeleteOpen, setAlertDeleteOpen] = useState(false);

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
              <Dialog>
                <form>
                  <DialogTrigger>
                    <div className="flex items-center gap-x-2 bg-primary text-primary-foreground p-1 rounded-md cursor-pointer">
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
            </div>
          </div>

          <div>
            <h4>N projects</h4>
            <Separator />
            <div className="my-3">
              <ItemGroup>
                <Item>
                  <ItemMedia variant="icon">
                    <Hash />
                  </ItemMedia>
                  <ItemContent>
                    <Link to={"/"}>
                      <ItemTitle>Default Variant</ItemTitle>
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
              </ItemGroup>
            </div>

            <div className="self-end">
              <Pagination>
                <PaginationContent>
                  <PaginationItem>
                    <PaginationLink href="#">1</PaginationLink>
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationLink href="#" isActive>
                      2
                    </PaginationLink>
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationLink href="#">3</PaginationLink>
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationLink href="#">4</PaginationLink>
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationLink href="#">5</PaginationLink>
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
