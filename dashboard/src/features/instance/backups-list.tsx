import { format } from 'date-fns';
import { Check, Download, Trash2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Checkbox } from '#/components/ui/checkbox';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useCreate, useDelete, useList, useListRecords, useTrigger } from '#/hooks/useBackups';

export function BackupsList() {
  const { data: configsData, isLoading } = useList('global');
  const configs = configsData?.data || [];
  const config = configs[0];

  const createBackup = useCreate();
  const triggerBackup = useTrigger();
  const deleteBackup = useDelete();

  const { data: recordsData, isLoading: isLoadingRecords } = useListRecords(config?.id || '');
  const records = recordsData?.data || [];

  const [name, setName] = useState('vessl-db');
  const [description, setDescription] = useState('Vessl database');
  const [schedule, setSchedule] = useState('0 0 * * *');
  const [retentionDays, setRetentionDays] = useState('7');

  useEffect(() => {
    if (config) {
      setName(config.name);
      setSchedule(config.schedule);
      setRetentionDays(config.retentionDays.toString());
    }
  }, [config]);

  const handleSave = async (e?: React.FormEvent) => {
    e?.preventDefault();
    try {
      if (config) {
        await deleteBackup.mutateAsync({ id: config.id });
      }
      await createBackup.mutateAsync({
        payload: {
          projectId: 'global',
          name,
          schedule,
          retentionDays: parseInt(retentionDays, 10),
        },
      });
      toast.success('Backup configuration saved');
    } catch {
      toast.error('Failed to save backup configuration');
    }
  };

  const handleTrigger = async () => {
    if (!config) {
      toast.error('Please save the configuration first');
      return;
    }
    try {
      await triggerBackup.mutateAsync({ id: config.id });
      toast.success('Backup triggered successfully');
    } catch {
      toast.error('Failed to trigger backup');
    }
  };

  if (isLoading) {
    return <div className="p-6 text-muted-foreground">Loading backups...</div>;
  }

  return (
    <div className="flex flex-col gap-8 pb-12">
      <div className="flex flex-col gap-4">
        <div className="flex items-center gap-4">
          <h2 className="font-bold text-2xl tracking-tight">Backup</h2>
          <Button
            size="sm"
            onClick={handleSave}
            disabled={createBackup.isPending || deleteBackup.isPending}
          >
            <Check className="mr-2 h-4 w-4" />
            {createBackup.isPending || deleteBackup.isPending ? 'Saving...' : 'Save'}
          </Button>
        </div>
        <p className="text-muted-foreground text-sm">Backup configuration for Vessl instance.</p>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
        <div className="space-y-2">
          <Label>UUID</Label>
          <Input disabled value={config?.id || 'Pending...'} />
        </div>
        <div className="space-y-2">
          <Label>Name</Label>
          <Input value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className="space-y-2">
          <Label>Description</Label>
          <Input value={description} onChange={(e) => setDescription(e.target.value)} />
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
        <div className="space-y-2">
          <Label>User</Label>
          <Input value="vessl" disabled />
        </div>
        <div className="space-y-2">
          <Label>Password</Label>
          <Input type="password" value="********" disabled />
        </div>
      </div>

      <div className="mt-4 flex flex-col gap-4">
        <div className="flex items-center gap-4">
          <h3 className="font-bold text-xl">Scheduled Backup</h3>
          <Button
            variant="outline"
            size="sm"
            onClick={handleSave}
            disabled={createBackup.isPending || deleteBackup.isPending}
          >
            Save
          </Button>
          <Button
            variant="secondary"
            size="sm"
            onClick={handleTrigger}
            disabled={triggerBackup.isPending || !config}
          >
            Backup Now
          </Button>
        </div>

        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Checkbox id="backup-enabled" defaultChecked />
            <Label htmlFor="backup-enabled" className="cursor-pointer">
              Backup Enabled
            </Label>
          </div>
          <div className="flex items-center gap-2">
            <Checkbox id="s3-enabled" />
            <Label htmlFor="s3-enabled" className="cursor-pointer">
              S3 Enabled
            </Label>
          </div>
          <div className="flex items-center gap-2">
            <Checkbox id="disable-local" />
            <Label htmlFor="disable-local" className="cursor-pointer">
              Disable Local Backup
            </Label>
          </div>
        </div>
      </div>

      <div className="mt-4 flex flex-col gap-4">
        <h3 className="font-bold text-xl">Settings</h3>
        <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
          <div className="space-y-2">
            <Label>Frequency *</Label>
            <Input value={schedule} onChange={(e) => setSchedule(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Timezone *</Label>
            <Input value="UTC" disabled />
          </div>
          <div className="space-y-2">
            <Label>Timeout *</Label>
            <Input value="3600" disabled />
          </div>
        </div>
      </div>

      <div className="mt-4 flex flex-col gap-4">
        <h3 className="font-bold text-xl">Backup Retention Settings</h3>
        <ul className="list-disc space-y-1 pl-5 text-muted-foreground text-sm">
          <li>Setting a value to 0 means unlimited retention.</li>
          <li>
            The retention rules work independently - whichever limit is reached first will trigger
            cleanup.
          </li>
        </ul>

        <h4 className="mt-2 font-bold">Local Backup Retention</h4>
        <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
          <div className="space-y-2">
            <Label>Number of backups to keep *</Label>
            <Input value="0" disabled />
          </div>
          <div className="space-y-2">
            <Label>Days to keep backups *</Label>
            <Input value={retentionDays} onChange={(e) => setRetentionDays(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Maximum storage (GB) *</Label>
            <Input value="0" disabled />
          </div>
        </div>
      </div>

      <div className="mt-8 flex flex-col gap-4">
        <div className="flex items-center gap-4">
          <h3 className="font-bold text-xl">Executions ({records.length})</h3>
          <Button variant="outline" size="sm">
            Cleanup Failed Backups
          </Button>
          <Button variant="destructive" size="sm">
            Cleanup Deleted
          </Button>
        </div>

        <div className="flex flex-col gap-4">
          {isLoadingRecords ? (
            <div className="py-8 text-muted-foreground">Loading executions...</div>
          ) : records.length === 0 ? (
            <div className="py-8 text-muted-foreground">No executions yet.</div>
          ) : (
            records.map((record) => (
              <div
                key={record.id}
                className="flex flex-col gap-3 rounded-lg border border-border bg-card/50 p-4"
              >
                <Badge
                  variant="outline"
                  className={
                    record.status === 'completed'
                      ? 'w-fit border-green-500/20 bg-green-500/10 text-green-500'
                      : record.status === 'failed'
                        ? 'w-fit border-red-500/20 bg-red-500/10 text-red-500'
                        : 'w-fit border-yellow-500/20 bg-yellow-500/10 text-yellow-500'
                  }
                >
                  {record.status === 'completed'
                    ? 'Success'
                    : record.status === 'failed'
                      ? 'Failed'
                      : 'Running'}
                </Badge>

                <div className="text-muted-foreground text-sm">
                  {record.startedAt
                    ? format(new Date(record.startedAt), 'MMM d, HH:mm')
                    : 'Unknown time'}{' '}
                  • Database: vessl • Size: {(record.fileSizeBytes / 1024 / 1024).toFixed(2)} MB
                  <br />
                  Location: {record.filePath}
                </div>

                <div className="flex items-center gap-2 text-sm">
                  <span className="text-muted-foreground">Backup Availability:</span>
                  <Badge
                    variant="outline"
                    className="gap-1 border-green-500/20 bg-green-500/10 text-green-500"
                  >
                    <Check className="h-3 w-3" /> Local Storage
                  </Badge>
                </div>

                <div className="mt-1 flex items-center gap-2">
                  <Button variant="outline" size="sm" asChild disabled={!record.s3Url}>
                    {record.s3Url ? (
                      <a href={record.s3Url} target="_blank" rel="noreferrer">
                        <Download className="mr-2 h-4 w-4" />
                        Download
                      </a>
                    ) : (
                      <span>
                        <Download className="mr-2 h-4 w-4" />
                        Download
                      </span>
                    )}
                  </Button>
                  <Button variant="destructive" size="sm">
                    <Trash2 className="mr-2 h-4 w-4" />
                    Delete
                  </Button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
