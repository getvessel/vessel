import { createFileRoute } from '@tanstack/react-router';
import { FileCode2, Play, UploadCloud } from 'lucide-react';
import { Textarea } from '#/components/ui/textarea';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const Route = createFileRoute('/_dashboard/projects/$projectId/compose')({
  component: ProjectComposePage,
});

function ProjectComposePage() {
  const { projectId } = Route.useParams();

  return (
    <OperationalPage
      title="Compose deploy"
      description="Paste or import a compose file, review generated resources, and deploy them into this project."
      scope={projectId}
      statusLabel="Ready to parse"
      statusTone="info"
      primaryAction={{ label: 'Deploy compose', icon: Play }}
      secondaryActions={[{ label: 'Import file', icon: UploadCloud, variant: 'outline' }]}
      metrics={[
        {
          label: 'Source',
          value: 'Compose',
          detail: 'Docker-compatible spec',
          icon: FileCode2,
        },
      ]}
    >
      <Textarea
        defaultValue={
          'services:\n  web:\n    image: ghcr.io/acme/app:latest\n    ports:\n      - "3000:3000"'
        }
        className="min-h-[420px] resize-none font-mono text-sm"
      />
    </OperationalPage>
  );
}
