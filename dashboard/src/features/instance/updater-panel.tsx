import { RefreshCw } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';

export const UpdaterPanel = () => (
  <EmptyPanel
    icon={RefreshCw}
    title="No update check has run"
    description="Check for updates to compare the current control plane version with the configured release channel."
    actionLabel="Check updates"
  />
);
