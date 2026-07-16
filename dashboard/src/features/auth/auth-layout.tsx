import { ShieldCheck } from 'lucide-react';
import type { ReactNode } from 'react';
import { AuthBrandPanel } from './auth-brand-panel';

export function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="grid min-h-screen bg-zinc-950 text-white lg:grid-cols-[minmax(0,1fr)_minmax(420px,0.86fr)]">
      <AuthBrandPanel />
      <main className="flex min-h-screen items-center justify-center border-white/10 border-l bg-zinc-950 px-5 py-8 sm:px-8">
        <section className="w-full max-w-[440px] rounded-lg border border-white/10 bg-white/[0.035] p-6 shadow-2xl shadow-black/30 backdrop-blur-sm sm:p-7">
          <div className="mb-6 flex items-center justify-between border-white/10 border-b pb-4 lg:hidden">
            <div className="flex items-center gap-3">
              <div className="flex size-9 items-center justify-center rounded-md bg-white text-zinc-950">
                <span className="font-semibold text-sm">V</span>
              </div>
              <div>
                <p className="font-semibold leading-none">Vessl</p>
                <p className="mt-1 text-white/45 text-xs">Control plane</p>
              </div>
            </div>
            <ShieldCheck className="size-4 text-emerald-300" />
          </div>
          {children}
          <div className="mt-6 flex items-center justify-between border-white/10 border-t pt-4 text-white/45 text-xs">
            <span>Session policy active</span>
            <span>Local auth</span>
          </div>
        </section>
      </main>
    </div>
  );
}
