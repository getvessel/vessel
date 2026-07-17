import { createFileRoute } from '@tanstack/react-router';
import { MigrationSettings } from '#/features/instance/migration-settings';

export const Route = createFileRoute('/_dashboard/settings/migration')({
  component: () => <MigrationSettings />,
});
