import { Sparkles } from 'lucide-react';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const AiSettings = () => (
  <SettingsPanel
    title="AI settings"
    icon={Sparkles}
    fields={[{ label: 'Diagnostics model', value: 'local' }]}
    toggles={[
      {
        label: 'Enable deployment diagnostics',
        description: 'Analyze failed builds and surface actionable summaries.',
        checked: false,
      },
    ]}
  />
);
