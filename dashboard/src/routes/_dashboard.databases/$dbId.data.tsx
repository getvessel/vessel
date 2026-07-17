import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/databases/$dbId/data')({
  component: () => null,
});
