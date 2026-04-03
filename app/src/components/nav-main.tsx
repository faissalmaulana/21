import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { NavLink } from "react-router"

export function NavMain({
  items,
}: {
  items: {
    title: string
    url: string
  }[]
}) {
  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2">
        <SidebarMenu>
          {items.map((item) => (
            <NavLink
              to={item.url}
              key={item.title}
              className={({ isActive }) =>
                `rounded-md px-3 py-1.5
                ${isActive
                  ? "bg-sidebar-primary text-sidebar-primary-foreground"
                  : "hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 "
                }`
              }
            >
              <SidebarMenuItem>
                <span>{item.title}</span>
              </SidebarMenuItem>
            </NavLink>
          ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}
