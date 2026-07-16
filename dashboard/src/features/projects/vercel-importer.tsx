import { UploadCloud } from 'lucide-react';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const VercelImporter = () => (
  <SettingsPanel
    title="Vercel import"
    icon={UploadCloud}
    fields={[
      { label: 'Access token', value: '' },
      { label: 'Project name', value: '' },
    ]}
    actionLabel="Validate import"
  />
);
