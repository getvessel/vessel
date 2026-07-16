import { createFileRoute, Link } from '@tanstack/react-router';
import { ArrowLeft, SearchX } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { EmptyPanel } from '#/features/dashboard/empty-panel';

export const Route = createFileRoute('/$')({
  component: NotFoundPage,
});

function NotFoundPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-6">
      <div className="w-full max-w-xl">
        <EmptyPanel
          icon={SearchX}
          title="Page not found"
          description="The route you opened does not exist in this dashboard or has moved to another workspace section."
        />
        <div className="mt-4 flex justify-center">
          <Button asChild>
            <Link to="/">
              <ArrowLeft className="size-4" />
              Return to dashboard
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
}
