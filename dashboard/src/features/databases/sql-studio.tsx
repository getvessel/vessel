import { Play } from 'lucide-react';
import { Textarea } from '#/components/ui/textarea';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const SqlStudio = () => (
  <OperationalPage
    title="SQL studio"
    description="Write, inspect, and execute SQL with schema context."
    scope="Database"
    statusLabel="Ready"
    statusTone="healthy"
    primaryAction={{ label: 'Run query', icon: Play }}
  >
    <Textarea defaultValue="select now();" className="min-h-64 resize-none font-mono text-sm" />
  </OperationalPage>
);
