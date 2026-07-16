import { createFileRoute } from '@tanstack/react-router';
import { Box, Database, Globe, Plus, Server, Star } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { StatusBadge } from '#/features/dashboard/status-badge';

export const Route = createFileRoute('/_dashboard/templates')({
  component: TemplatesPage,
});

const templates = [
  {
    id: 'node-express',
    name: 'Node.js Express',
    description: 'Minimal HTTP service with sensible build and start commands.',
    category: 'Framework',
    icon: Server,
  },
  {
    id: 'go-fiber',
    name: 'Go Fiber',
    description: 'Small Go API template with fast container startup.',
    category: 'Framework',
    icon: Box,
  },
  {
    id: 'python-fastapi',
    name: 'Python FastAPI',
    description: 'Typed Python API starter with health checks and docs.',
    category: 'Framework',
    icon: Globe,
  },
  {
    id: 'postgres',
    name: 'Managed PostgreSQL',
    description: 'Provision a Postgres database with connection metadata.',
    category: 'Database',
    icon: Database,
  },
];

function TemplatesPage() {
  return (
    <OperationalPage
      title="Templates"
      description="Start from vetted service and data templates, then tune the generated deployment settings before launch."
      scope="Create"
      statusLabel="Catalog ready"
      statusTone="healthy"
      primaryAction={{ label: 'New template', icon: Plus }}
      metrics={[
        {
          label: 'Templates',
          value: templates.length.toString(),
          detail: 'Deployable starters',
          icon: Star,
        },
        {
          label: 'Managed resources',
          value: '1',
          detail: 'Database templates',
          icon: Database,
        },
      ]}
    >
      <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
        {templates.map((template) => (
          <Card key={template.id} className="shadow-none">
            <CardHeader className="space-y-3">
              <div className="flex items-center justify-between">
                <div className="flex size-9 items-center justify-center rounded-md border bg-muted/40">
                  <template.icon className="size-4 text-muted-foreground" />
                </div>
                <StatusBadge
                  label={template.category}
                  tone={template.category === 'Database' ? 'info' : 'neutral'}
                />
              </div>
              <CardTitle className="text-base">{template.name}</CardTitle>
            </CardHeader>
            <CardContent className="grid gap-4">
              <p className="min-h-12 text-muted-foreground text-sm leading-6">
                {template.description}
              </p>
              <Button variant="outline">Use template</Button>
            </CardContent>
          </Card>
        ))}
      </div>
    </OperationalPage>
  );
}
