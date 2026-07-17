import { createFileRoute } from '@tanstack/react-router';
import { GeneralSettings } from '#/features/instance/general-settings';

export const Route = createFileRoute('/_dashboard/settings/general')({
  component: GeneralSettings,
});
