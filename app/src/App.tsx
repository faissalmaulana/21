import { Project, Projects, Root, Task, Tasks } from "@/pages"
import { BrowserRouter, Routes, Route } from "react-router"

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route>
          <Route index element={<Root />} />
        </Route>

        <Route path="tasks">
          <Route index element={<Tasks />} />
          <Route path=":id" element={<Task />} />
        </Route>

        <Route path="projects">
          <Route index element={<Projects />} />
          <Route path=":id" element={<Project />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
export default App
