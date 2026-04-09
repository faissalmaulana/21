interface Project {
  id: string
  name: string
}

interface GetProjectsParams {
  search?: string
  withArchive?: boolean
  page: string
}

interface PostProjectBodyParam {
  name: string
}

export type { Project, GetProjectsParams, PostProjectBodyParam }
