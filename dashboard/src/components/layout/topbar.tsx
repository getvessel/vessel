import { BellIcon, Menu, PlusIcon, SearchIcon, TrainFront, Triangle, X } from 'lucide-react';
import { toast } from 'sonner';

interface TopbarProps {
  onOpenCommand: () => void;
  onMenuToggle: () => void;
  mobileMenuOpen: boolean;
}

export function Topbar({ onOpenCommand, onMenuToggle, mobileMenuOpen }: TopbarProps) {
  return (
    <header className="flex h-14 shrink-0 items-center justify-between gap-2 bg-transparent px-4 md:px-8">
      <div className="flex items-center gap-2">
        <button
          type="button"
          onClick={onMenuToggle}
          className="flex h-9 w-9 items-center justify-center rounded-xl border border-border/60 text-muted-foreground transition-colors hover:bg-muted md:hidden"
        >
          {mobileMenuOpen ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
        </button>
      </div>

      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={onOpenCommand}
          className="flex h-9 items-center gap-2 rounded-xl border border-border/60 bg-muted/40 px-3 text-muted-foreground text-sm transition-all hover:border-border hover:bg-muted hover:text-foreground active:scale-[0.97]"
        >
          <SearchIcon className="h-4 w-4 shrink-0" />
          <span className="hidden sm:inline">Search...</span>
          <kbd className="hidden rounded-md border bg-background/60 px-1.5 py-0.5 font-mono text-[11px] leading-none sm:inline">
            ⌘K
          </kbd>
        </button>

        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => toast.info('Railway Import', { description: 'Coming soon' })}
            className="flex h-9 items-center gap-2 rounded-xl border border-indigo-500/30 bg-indigo-500/10 px-3 font-semibold text-indigo-400 text-xs tracking-wider transition-all hover:bg-indigo-500/20 active:scale-[0.97] md:px-4"
          >
            <TrainFront className="h-4 w-4 shrink-0 md:h-3.5 md:w-3.5" />
            <span className="hidden md:inline">IMPORT RAILWAY</span>
          </button>

          <button
            type="button"
            onClick={() => toast.info('Vercel Import', { description: 'Coming soon' })}
            className="flex h-9 items-center gap-2 rounded-xl border border-border/60 bg-zinc-950 px-3 font-semibold text-xs text-zinc-100 tracking-wider transition-all hover:bg-zinc-900 active:scale-[0.97] md:px-4 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
          >
            <Triangle className="h-4 w-4 shrink-0 fill-current md:h-3.5 md:w-3.5" />
            <span className="hidden md:inline">IMPORT VERCEL</span>
          </button>
        </div>

        <button
          type="button"
          onClick={() => toast.info('New resource', { description: 'Creation menu coming soon' })}
          className="flex h-9 items-center gap-1.5 rounded-xl bg-primary px-4 font-semibold text-primary-foreground text-xs tracking-wider shadow-lg shadow-primary/25 transition-all hover:brightness-110 active:scale-[0.97]"
        >
          <PlusIcon className="h-4 w-4" />
          <span>NEW</span>
        </button>

        <button
          type="button"
          onClick={() => {}}
          className="relative flex h-9 w-9 items-center justify-center rounded-xl border border-border/60 transition-colors hover:bg-muted"
        >
          <BellIcon className="h-4 w-4" />
          <div className="absolute top-2 right-2 h-1.5 w-1.5 rounded-full bg-primary ring-2 ring-transparent" />
        </button>
      </div>
    </header>
  );
}
