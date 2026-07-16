import { FolderOpen } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';

export const ObjectBrowser = () => (
  <EmptyPanel
    icon={FolderOpen}
    title="No objects selected"
    description="Choose a bucket to browse folders, inspect metadata, download objects, and manage retention."
  />
);
