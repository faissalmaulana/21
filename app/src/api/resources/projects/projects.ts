import type { Pagination, RawResponse } from "@/api/response";
import type { Project, GetProjectsParams } from "@/api/resources";
import { AppError } from "@/lib/app-error";

const PROJECTS_KEY = "projects"

/**
 * @todo
 * - use opt parameters for filtering
 */
// eslint-disable-next-line @typescript-eslint/no-unused-vars
async function getProjects(opt: GetProjectsParams): Promise<{ projects: Project[]; pagination: Pagination | null; }> {
  const response = await fetch('/dummies/projects.json');
  if (!response.ok) {
    throw new AppError({
      status: 500,
      message: "Internal Server Error"
    });
  }

  const result: RawResponse<{ projects: Project[]; pagination: Pagination | null; }> = await response.json();

  switch (result.status) {
    case 200:
      return result.data;
    default:
      throw new AppError({
        status: 500,
        message: "Internal Server Error"
      });
  }
}

export {
  PROJECTS_KEY,
  getProjects
}
