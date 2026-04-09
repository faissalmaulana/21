import type { Pagination, RawResponse } from "@/api/response"
import type { Project, GetProjectsParams, PostProjectBodyParam } from "@/api/resources"
import { AppError } from "@/lib/app-error"

const PROJECTS_KEY = "projects"

async function getProjects(
  opt: GetProjectsParams
): Promise<{ projects: Project[]; pagination: Pagination | null }> {
  // in production (window.location.origin) should be changed to host of the server
  const url = new URL('http://localhost:8080/api/projects', window.location.origin)

  // build query params
  if (opt.search) {
    url.searchParams.set("search", opt.search)
  }

  if (opt.withArchive) {
    url.searchParams.set("archive", "true")
  }

  if (opt.page !== "1" && opt.page !== "") {
    url.searchParams.set("page", opt.page)
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

async function postProject(body: PostProjectBodyParam): Promise<string> {
  let response: Response

  try {
    response = await fetch("http://localhost:8080/api/projects", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    })
  } catch {
    throw new AppError({
      status: 0,
      message: "Network error",
    })
  }

  let result: RawResponse<string>

  try {
    result = await response.json()
  } catch {
    throw new AppError({
      status: response.status,
      message: "Invalid JSON response",
    })
  }

  switch (result.status) {
    case 201:
      return result.data

    case 400:
      throw new AppError({
        status: 400,
        message: result?.error?.message ?? "Bad Request",
      })

    default:
      throw new AppError({
        status: result.status ?? response.status ?? 500,
        message: result?.error?.message ?? "Internal Server Error",
      })
  }
}

export { PROJECTS_KEY, getProjects, postProject }
