import type { Pagination, RawResponse } from "@/api/response"
import type { Project, GetProjectsParams } from "@/api/resources"
import { AppError } from "@/lib/app-error"

const PROJECTS_KEY = "projects"

async function getProjects(
  opt: GetProjectsParams
): Promise<{ projects: Project[]; pagination: Pagination | null }> {
  const url = new URL('/dummies/projects.json', window.location.origin)

  // build query params
  if (opt.search) {
    url.searchParams.set("search", opt.search)
  }

  const response = await fetch(url.toString())

  const result: RawResponse<{
    projects: Project[]
    pagination: Pagination | null
  }> = await response.json()

  // API-level error handling
  switch (result.status) {
    case 200:
      return result.data

    case 400:
      throw new AppError({
        status: 400,
        message: result?.error?.message ?? "Bad Request",
      })

    default:
      throw new AppError({
        status: result.status ?? 500,
        message: result?.error?.message ?? "Internal Server Error",
      })
  }
}

export { PROJECTS_KEY, getProjects }
