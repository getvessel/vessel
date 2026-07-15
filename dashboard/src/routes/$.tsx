import { createFileRoute, Link } from '@tanstack/react-router';
import { AlertTriangle, ArrowLeft } from 'lucide-react';

export const Route = createFileRoute('/$')({
  component: NotFound,
});

function NotFound() {
  return (
    <div className="flex h-full w-full flex-col items-center justify-center p-8 text-center text-zinc-100">
      <div className="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-red-500/10">
        <AlertTriangle className="h-10 w-10 text-red-500" />
      </div>
      <h1 className="mb-2 font-bold text-4xl tracking-tight">404 - Page Not Found</h1>
      <p className="mb-8 max-w-md text-zinc-400">
        The page you are looking for doesn't exist, has been moved, or you don't have access to it.
      </p>
      <Link
        to="/"
        className="inline-flex h-10 items-center justify-center rounded-md bg-zinc-50 px-4 py-2 font-medium text-sm text-zinc-900 transition-colors hover:bg-zinc-50/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-zinc-400 disabled:pointer-events-none disabled:opacity-50"
      >
        <ArrowLeft className="mr-2 h-4 w-4" />
        Back to Dashboard
      </Link>
    </div>
  );
}
