import { createFileRoute, Link } from '@tanstack/react-router';
import { Bell, DatabaseBackup, GitBranch, Globe2, Shield, Users, Wrench } from 'lucide-react';
import { Card, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const Route = createFileRoute('/_dashboard/settings/')({
  component: SettingsIndexPage,
});

const settingsLinks = [
  {
    title: 'Users',
    description: 'Invite operators, adjust roles, and review access.',
    to: '/settings/users',
    icon: Users,
  },
  {
    title: 'DNS',
    description: 'Configure wildcard domains and routing defaults.',
    to: '/settings/dns',
    icon: Globe2,
  },
  {
    title: 'Git apps',
    description: 'Connect providers for repository imports and deploys.',
    to: '/settings/git-apps',
    icon: GitBranch,
  },
  {
    title: 'Backups',
    description: 'Define snapshot retention and destinations.',
    to: '/settings/backups',
    icon: DatabaseBackup,
  },
  {
    title: 'Updates',
    description: 'Check versions and deploy control plane upgrades.',
    to: '/settings/updates',
    icon: Wrench,
  },
  {
    title: 'OAuth',
    description: 'Manage identity providers and login policy.',
    to: '/settings/oauth',
    icon: Shield,
  },
];

function SettingsIndexPage() {
  return (
    <OperationalPage
      title="Instance settings"
      description="Control identity, DNS, Git integration, backups, maintenance, updates, and migration settings for this Vessl instance."
      scope="Admin"
      statusLabel="Configurable"
      statusTone="healthy"
      metrics={[
        {
          label: 'Areas',
          value: settingsLinks.length.toString(),
          detail: 'Admin sections',
          icon: Bell,
        },
      ]}
    >
      <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
        {settingsLinks.map((item) => (
          <Link key={item.to} to={item.to as never} className="group">
            <Card className="h-full shadow-none transition-colors group-hover:border-primary/50">
              <CardHeader>
                <item.icon className="size-4 text-muted-foreground group-hover:text-primary" />
                <CardTitle className="text-base">{item.title}</CardTitle>
                <CardDescription>{item.description}</CardDescription>
              </CardHeader>
            </Card>
          </Link>
        ))}
      </div>
    </OperationalPage>
  );
}
