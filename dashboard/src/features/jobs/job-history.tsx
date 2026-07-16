import { History } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';

export const JobHistory = () => (
  <EmptyPanel
    icon={History}
    title="No job runs yet"
    description="Run history will show timestamps, duration, output, status, and retry context after a job executes."
  />
);
