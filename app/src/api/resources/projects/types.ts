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

interface UpdateProjectBodyParam {
  name?: string
  to_be_archived?: boolean
}

export type { Project, GetProjectsParams, PostProjectBodyParam, UpdateProjectBodyParam }
