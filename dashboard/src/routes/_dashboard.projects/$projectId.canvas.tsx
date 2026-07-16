import { createFileRoute } from '@tanstack/react-router';
import { Boxes, Database, Network } from 'lucide-react';
import { CanvasGridPanel } from '#/features/canvas/canvas-grid-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const Route = createFileRoute('/_dashboard/projects/$projectId/canvas')({
  component: ProjectCanvasPage,
});

function ProjectCanvasPage() {
  const { projectId } = Route.useParams();

  return (
    <OperationalPage
      title="Environment canvas"
      description="Visualize services, databases, storage, and routing relationships for this project environment."
      scope={projectId}
      statusLabel="Canvas ready"
      statusTone="info"
      metrics={[
        {
          label: 'Services',
          value: '0',
          detail: 'Application nodes',
          icon: Boxes,
        },
        {
          label: 'Databases',
          value: '0',
          detail: 'Data nodes',
          icon: Database,
        },
        {
          label: 'Routes',
          value: '0',
          detail: 'Network edges',
          icon: Network,
        },
      ]}
    >
      <CanvasGridPanel />
    </OperationalPage>
  );
}
