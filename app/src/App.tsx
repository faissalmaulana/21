import { Project, Projects, Root, Task, Tasks } from "@/pages"
import { BrowserRouter, Routes, Route } from "react-router"
import { AppLayout } from "@/components/app-layout"

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppLayout />}>
          <Route index element={<Root />} />

          <Route path="tasks">
            <Route index element={<Tasks />} />
            <Route path=":id" element={<Task />} />
          </Route>

          <Route path="projects">
            <Route index element={<Projects />} />
            <Route path=":id" element={<Project />} />
          </Route>
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
export default App
