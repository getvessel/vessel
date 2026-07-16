import { Mail } from 'lucide-react';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const EmailSettings = () => (
  <SettingsPanel
    title="Email delivery"
    icon={Mail}
    fields={[
      { label: 'SMTP host', value: '' },
      { label: 'From address', value: 'noreply@localhost' },
    ]}
    actionLabel="Save email"
  />
);
