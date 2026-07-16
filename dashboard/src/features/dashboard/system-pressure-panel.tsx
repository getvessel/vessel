import type { LucideIcon } from 'lucide-react';
import { Cpu, HardDrive, MemoryStick, ShieldCheck } from 'lucide-react';
import { PolarAngleAxis, RadialBar, RadialBarChart } from 'recharts';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { ChartContainer } from '#/components/ui/chart';
import { Skeleton } from '#/components/ui/skeleton';
import { formatGb, formatMb } from '#/features/dashboard/dashboard-format';
import type { SystemStats } from '#/interfaces/system';

interface SystemPressurePanelProps {
  stats?: SystemStats;
  state: string;
}

interface PressureDetail {
  label: string;
  value?: number;
  detail: string;
  icon: LucideIcon;
  color: string;
  tone: string;
}

export function SystemPressurePanel({ stats, state }: SystemPressurePanelProps) {
  const details: PressureDetail[] = [
    {
      label: 'CPU',
      value: stats?.cpu.percent,
      detail: stats ? `${stats.cpu.cores} cores available` : 'Waiting for daemon metrics',
      icon: Cpu,
      color: 'var(--primary)',
      tone: 'text-primary',
    },
    {
      label: 'Memory',
      value: stats?.memory.percent,
      detail: stats
        ? `${formatMb(stats.memory.freeMb)} free of ${formatMb(stats.memory.totalMb)}`
        : 'Waiting for daemon metrics',
      icon: MemoryStick,
      color: 'var(--secondary)',
      tone: 'text-teal-300',
    },
    {
      label: 'Disk',
      value: stats?.disk.percent,
      detail: stats
        ? `${formatGb(stats.disk.freeGb)} free of ${formatGb(stats.disk.totalGb)}`
        : 'Waiting for daemon metrics',
      icon: HardDrive,
      color: 'var(--warning)',
      tone: 'text-yellow-300',
    },
  ];

  return (
    <Card className="rounded-lg shadow-none">
      <CardHeader className="border-b">
        <CardTitle className="flex items-center gap-2 text-base">
          <ShieldCheck className="size-4 text-muted-foreground" />
          System pressure
        </CardTitle>
        <CardDescription>{state}. Current utilization with warning threshold.</CardDescription>
      </CardHeader>
      <CardContent className="grid gap-4">
        {stats ? (
          <div className="grid gap-3 md:grid-cols-3">
            {details.map((item) => {
              const value = Math.round(item.value ?? 0);
              const gaugeData = [{ name: item.label, value, fill: item.color }];

              return (
                <div key={item.label} className="rounded-lg border bg-muted/10 p-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <item.icon className="size-4 text-muted-foreground" />
                      <span className="font-medium text-sm">{item.label}</span>
                    </div>
                    <span className="text-muted-foreground text-xs">85% limit</span>
                  </div>

                  <div className="relative mx-auto mt-4 h-36 max-w-48">
                    <ChartContainer
                      config={{
                        value: {
                          label: item.label,
                          color: item.color,
                        },
                      }}
                      className="h-full w-full"
                    >
                      <RadialBarChart
                        accessibilityLayer
                        data={gaugeData}
                        startAngle={180}
                        endAngle={0}
                        innerRadius="72%"
                        outerRadius="100%"
                      >
                        <PolarAngleAxis
                          type="number"
                          domain={[0, 100]}
                          angleAxisId={0}
                          tick={false}
                        />
                        <RadialBar
                          dataKey="value"
                          background={{ fill: 'var(--muted)' }}
                          cornerRadius={8}
                        />
                      </RadialBarChart>
                    </ChartContainer>
                    <div className="absolute inset-x-0 bottom-4 text-center">
                      <p className={`font-semibold text-3xl ${item.tone}`}>{value}%</p>
                      <p className="mt-1 text-muted-foreground text-xs">current</p>
                    </div>
                  </div>

                  <div className="mt-3 border-t pt-3">
                    <p className="text-muted-foreground text-xs">{item.detail}</p>
                  </div>
                </div>
              );
            })}
          </div>
        ) : (
          <div className="grid gap-3 md:grid-cols-3">
            {details.map((item) => (
              <Skeleton key={item.label} className="h-56 w-full" />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
