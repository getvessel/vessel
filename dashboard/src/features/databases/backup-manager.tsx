import { DatabaseBackup } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';

export const BackupManager = () => (
  <EmptyPanel
    icon={DatabaseBackup}
    title="No backups recorded"
    description="Create a snapshot or enable scheduled backups to protect this database before risky operations."
    actionLabel="Create backup"
  />
);
