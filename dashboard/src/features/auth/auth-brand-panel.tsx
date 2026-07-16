import {
  Activity,
  Boxes,
  CheckCircle2,
  Clock3,
  Database,
  GitBranch,
  LockKeyhole,
  Server,
  ShieldCheck,
} from 'lucide-react';
import { Badge } from '#/components/ui/badge';

const capabilities = [
  { label: 'Deployments', value: 'Tracked', icon: GitBranch },
  { label: 'Databases', value: 'Managed', icon: Database },
  { label: 'Services', value: 'Live', icon: Boxes },
];

const healthChecks = [
  { label: 'Gateway', value: 'Reachable', icon: Server },
  { label: 'Secrets', value: 'Sealed', icon: LockKeyhole },
  { label: 'Jobs', value: 'Idle', icon: Clock3 },
];

export function AuthBrandPanel() {
  return (
    <section className="hidden min-h-screen bg-zinc-950 text-white lg:flex">
      <div className="flex w-full flex-col justify-between p-10">
        <div className="flex items-center gap-3">
          <div className="flex size-9 items-center justify-center rounded-md bg-white text-zinc-950">
            <Activity className="size-5" />
          </div>
          <div>
            <p className="font-semibold text-lg leading-none">Vessl</p>
            <p className="mt-1 text-white/55 text-xs">Control plane</p>
          </div>
        </div>

        <div className="grid max-w-3xl gap-8">
          <div className="max-w-xl">
            <Badge className="mb-5 border-white/10 bg-white/10 text-white hover:bg-white/10">
              Self-hosted PaaS
            </Badge>
            <h1 className="max-w-2xl font-semibold text-4xl leading-tight tracking-tight">
              Enter a calm control plane for live infrastructure.
            </h1>
            <p className="mt-5 max-w-lg text-sm text-white/60 leading-6">
              Review deploys, databases, domains, and instance health without losing operational
              context.
            </p>
          </div>

          <div className="grid max-w-2xl gap-4">
            <div className="rounded-lg border border-white/10 bg-white/[0.035]">
              <div className="flex items-center justify-between border-white/10 border-b px-4 py-3">
                <div>
                  <p className="font-medium text-sm">Instance state</p>
                  <p className="mt-1 text-white/45 text-xs">Last checked moments ago</p>
                </div>
                <Badge className="border-emerald-400/20 bg-emerald-400/10 text-emerald-200 hover:bg-emerald-400/10">
                  <CheckCircle2 className="size-3" />
                  Healthy
                </Badge>
              </div>
              <div className="grid grid-cols-3 divide-x divide-white/10">
                {healthChecks.map((item) => (
                  <div key={item.label} className="grid gap-3 p-4">
                    <item.icon className="size-4 text-white/45" />
                    <div>
                      <p className="text-white/45 text-xs">{item.label}</p>
                      <p className="mt-1 font-medium text-sm">{item.value}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="grid gap-3">
              {capabilities.map((item) => (
                <div
                  key={item.label}
                  className="flex items-center justify-between rounded-lg border border-white/10 bg-white/[0.035] px-4 py-3"
                >
                  <div className="flex items-center gap-3">
                    <item.icon className="size-4 text-white/65" />
                    <span className="font-medium text-sm">{item.label}</span>
                  </div>
                  <div className="flex items-center gap-2 text-white/55 text-xs">
                    <span>{item.value}</span>
                    <ShieldCheck className="size-4 text-emerald-300" />
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="flex items-center justify-between border-white/10 border-t pt-6 text-white/45 text-xs">
          <span>Local instance</span>
          <span>API gateway :8080</span>
        </div>
      </div>
    </section>
  );
}
