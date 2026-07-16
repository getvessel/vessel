import type { LucideIcon } from 'lucide-react';
import type { FormEvent } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import { Textarea } from '#/components/ui/textarea';

interface SettingsField {
  label: string;
  value: string;
  type?: 'input' | 'textarea';
}

interface SettingsToggle {
  label: string;
  description: string;
  checked?: boolean;
}

interface SettingsPanelProps {
  title: string;
  icon: LucideIcon;
  fields?: SettingsField[];
  toggles?: SettingsToggle[];
  actionLabel?: string;
}

export function SettingsPanel({
  title,
  icon: Icon,
  fields = [],
  toggles = [],
  actionLabel = 'Save changes',
}: SettingsPanelProps) {
  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    toast.info('Changes staged', {
      description: 'Saving will be enabled when this settings endpoint is connected.',
    });
  };

  return (
    <Card className="rounded-lg shadow-none">
      <CardHeader className="border-b">
        <CardTitle className="flex items-center gap-2 text-base">
          <Icon className="size-4 text-muted-foreground" />
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="grid gap-5">
          {fields.length > 0 ? (
            <div className="grid gap-4 md:grid-cols-2">
              {fields.map((field) => (
                <div
                  key={field.label}
                  className={field.type === 'textarea' ? 'grid gap-2 md:col-span-2' : 'grid gap-2'}
                >
                  <Label>{field.label}</Label>
                  {field.type === 'textarea' ? (
                    <Textarea defaultValue={field.value} className="min-h-24 resize-none" />
                  ) : (
                    <Input defaultValue={field.value} />
                  )}
                </div>
              ))}
            </div>
          ) : null}
          {toggles.map((toggle) => (
            <div
              key={toggle.label}
              className="flex items-start justify-between gap-6 rounded-lg border p-3"
            >
              <div className="space-y-1">
                <p className="font-medium text-sm">{toggle.label}</p>
                <p className="text-muted-foreground text-sm leading-5">{toggle.description}</p>
              </div>
              <Switch defaultChecked={toggle.checked} />
            </div>
          ))}
          <div>
            <Button type="submit">{actionLabel}</Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
