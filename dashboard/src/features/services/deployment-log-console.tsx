import { Copy, Download, Search } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';

const logLines = [
  { time: '10:42:01', level: 'info', text: 'Starting deployment worker' },
  { time: '10:42:04', level: 'info', text: 'Cloning repository from main' },
  {
    time: '10:42:16',
    level: 'info',
    text: 'Installing dependencies with detected package manager',
  },
  { time: '10:43:08', level: 'info', text: 'Build completed successfully' },
  { time: '10:43:21', level: 'info', text: 'Container health check passed' },
  { time: '10:43:28', level: 'success', text: 'Release promoted to production' },
];

export function DeploymentLogConsole() {
  return (
    <section className="overflow-hidden rounded-lg border bg-card">
      <div className="flex flex-col gap-3 border-b p-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h2 className="font-semibold text-base">Live logs</h2>
          <p className="text-muted-foreground text-sm">Streaming build and release output.</p>
        </div>
        <div className="flex gap-2">
          <div className="relative">
            <Search className="absolute top-1/2 left-2.5 size-3.5 -translate-y-1/2 text-muted-foreground" />
            <Input className="h-8 w-48 pl-8" placeholder="Search logs" />
          </div>
          <Button type="button" size="sm" variant="outline">
            <Copy className="size-4" />
          </Button>
          <Button type="button" size="sm" variant="outline">
            <Download className="size-4" />
          </Button>
        </div>
      </div>
      <div className="max-h-[320px] overflow-auto bg-[#07070b] p-4 font-mono text-xs">
        {logLines.map((line) => (
          <div
            key={`${line.time}-${line.text}`}
            className="grid grid-cols-[72px_64px_1fr] gap-3 py-1"
          >
            <span className="text-zinc-500">{line.time}</span>
            <span className={line.level === 'success' ? 'text-emerald-400' : 'text-primary'}>
              {line.level}
            </span>
            <span className="text-zinc-200">{line.text}</span>
          </div>
        ))}
      </div>
    </section>
  );
}
