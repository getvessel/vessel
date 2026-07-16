import { Copy, Download, Search } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';

interface LogPanelProps {
  title: string;
  lines: string[];
}

export function LogPanel({ title, lines }: LogPanelProps) {
  return (
    <div className="overflow-hidden rounded-lg border bg-zinc-950 text-zinc-100">
      <div className="flex flex-wrap items-center justify-between gap-3 border-zinc-800 border-b bg-zinc-900 px-4 py-3">
        <div>
          <h2 className="font-medium text-sm">{title}</h2>
          <p className="text-xs text-zinc-400">Live stream with timestamps and filters</p>
        </div>
        <div className="flex items-center gap-2">
          <div className="relative">
            <Search className="absolute top-1/2 left-2 size-3.5 -translate-y-1/2 text-zinc-500" />
            <Input
              placeholder="Search logs"
              className="h-8 w-44 border-zinc-800 bg-zinc-950 pl-8 text-zinc-100 placeholder:text-zinc-500"
            />
          </div>
          <Button size="icon" variant="ghost" className="text-zinc-300 hover:bg-zinc-800">
            <Copy className="size-4" />
          </Button>
          <Button size="icon" variant="ghost" className="text-zinc-300 hover:bg-zinc-800">
            <Download className="size-4" />
          </Button>
        </div>
      </div>
      <pre className="max-h-[520px] overflow-auto p-4 font-mono text-xs leading-6">
        {lines.map((line) => `${line}\n`)}
      </pre>
    </div>
  );
}
