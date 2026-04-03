import type { Pagination, RawResponse } from "@/api/response";
import type { Project } from "@/api/resources";

const PROJECTS_KEY = "projects"

async function getProjects(): Promise<{ projects: Project[]; pagination: Pagination | null; }> {
  const response = await fetch('/dummies/projects.json');
  if (!response.ok) {
    throw new Error('Internal Server Error');
  }

  const result: RawResponse<{ projects: Project[]; pagination: Pagination | null; }> = await response.json();

  switch (result.status) {
    case 200:
      return result.data;
    default:
      throw new Error('Internal Server Error');
  }
}

export {
  PROJECTS_KEY,
  getProjects
}
