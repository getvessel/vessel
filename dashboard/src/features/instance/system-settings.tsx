import { ServerCog } from 'lucide-react';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const SystemSettings = () => (
  <SettingsPanel
    title="System settings"
    icon={ServerCog}
    fields={[
      { label: 'Instance name', value: 'Vessl' },
      { label: 'Base URL', value: 'http://localhost:3000' },
    ]}
    toggles={[
      {
        label: 'Require secure cookies',
        description: 'Only send session cookies over HTTPS in production.',
        checked: true,
      },
    ]}
  />
);
