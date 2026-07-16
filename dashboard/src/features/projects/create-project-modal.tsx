import { FolderPlus } from 'lucide-react';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const CreateProjectForm = () => (
  <SettingsPanel
    title="Create project"
    icon={FolderPlus}
    fields={[
      { label: 'Project name', value: '' },
      { label: 'Description', value: '', type: 'textarea' },
    ]}
    actionLabel="Create project"
  />
);
