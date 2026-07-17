import { createFileRoute } from '@tanstack/react-router';
import { OAuthProvidersList } from '#/features/instance/oauth-providers-list';

export const Route = createFileRoute('/_dashboard/settings/oauth')({
  component: OAuthProvidersList,
});
