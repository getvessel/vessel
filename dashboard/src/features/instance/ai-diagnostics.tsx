import { Sparkles } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';

export const AiDiagnostics = () => (
  <EmptyPanel
    icon={Sparkles}
    title="No diagnostics yet"
    description="Deployment analysis will appear here after a failed build or runtime health check."
  />
);
