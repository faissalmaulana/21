interface Project {
  id: string
  name: string
}

interface GetProjectsParams {
  search?: string
  withArchive?: boolean
}

interface PostProjectBodyParam {
  name: string
}

export type { Project, GetProjectsParams, PostProjectBodyParam }
