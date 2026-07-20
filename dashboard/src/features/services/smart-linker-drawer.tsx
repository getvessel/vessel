import { useState } from 'react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '#/components/ui/sheet';

interface SmartLinkerDrawerProps {
  projectId: string;
  onLinkVariable: (key: string, value: string) => void;
  exampleEnvValues?: Record<string, string>;
}

export function SmartLinkerDrawer({ onLinkVariable, exampleEnvValues }: SmartLinkerDrawerProps) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Sheet open={isOpen} onOpenChange={setIsOpen}>
      <SheetTrigger asChild>
        <Button variant="outline" size="sm">
          Smart Linker
        </Button>
      </SheetTrigger>
      <SheetContent className="w-100 overflow-y-auto sm:w-135">
        <SheetHeader>
          <SheetTitle>Smart Variable Linker</SheetTitle>
          <SheetDescription>
            Auto-link environment variables from databases or other services in this project.
          </SheetDescription>
        </SheetHeader>

        <div className="mt-6 space-y-6">
          {exampleEnvValues && Object.keys(exampleEnvValues).length > 0 && (
            <div className="space-y-3">
              <h3 className="font-medium text-sm">.env.example Suggestions</h3>
              <div className="flex flex-wrap gap-2">
                {Object.entries(exampleEnvValues).map(([key, value]) => (
                  <Badge
                    key={key}
                    variant="secondary"
                    className="cursor-pointer hover:bg-zinc-200 dark:hover:bg-zinc-800"
                    onClick={() => {
                      onLinkVariable(key, value || '');
                      setIsOpen(false);
                    }}
                  >
                    {key}
                  </Badge>
                ))}
              </div>
            </div>
          )}

          <div className="space-y-4">
            <h3 className="font-medium text-sm">Available Resources</h3>
            <p className="text-sm text-zinc-500">
              Feature disabled: No suitable endpoint exists to list project resources.
            </p>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
