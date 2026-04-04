interface Project {
  id: string
  name: string
}

interface GetProjectsParams {
  search?: string
}

export type {
  Project,
  GetProjectsParams
}
