import { Users } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';

export const UserManagement = () => (
  <EmptyPanel
    icon={Users}
    title="User management"
    description="Invite users, assign roles, and review pending access requests from the instance settings."
    actionLabel="Invite user"
  />
);
