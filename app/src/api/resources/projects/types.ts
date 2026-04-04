interface Project {
  id: string
  name: string
}

interface GetProjectsParams {
  search?: string
  withArchive?: boolean
}

export type { Project, GetProjectsParams }
