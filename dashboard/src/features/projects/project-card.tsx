import { FolderKanban } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';

export const ProjectCard = () => (
  <EmptyPanel
    icon={FolderKanban}
    title="Project summary"
    description="Open a project to inspect its services, jobs, canvas, variables, and settings."
  />
);
