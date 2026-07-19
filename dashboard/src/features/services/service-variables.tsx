import { FileSearch, Loader2, Plus, Trash } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { useCreate, useDelete, useList, useEnvSuggestions } from '#/hooks/useServices';

export function ServiceVariables({
  app,
}: {
  app: any /* biome-ignore lint/suspicious/noExplicitAny: any */;
}) {
  const { data, isLoading } = useList(app.id);
  const [scanEnabled, setScanEnabled] = useState(false);
  const { data: suggestionsData, isLoading: isScanning } = useEnvSuggestions(app.id, scanEnabled);
  const createVar = useCreate();
  const deleteVar = useDelete();
  const [key, setKey] = useState('');
  const [val, setVal] = useState('');

  if (isLoading) return <Loader2 className="animate-spin" />;

  const vars = data?.data || [];
  const suggestions = suggestionsData?.data || [];
  
  // Filter out suggestions that already exist in vars
  const availableSuggestions = suggestions.filter(
    (s: any) => !vars.find((v: any) => v.key === s.key)
  );

  return (
    <div className="space-y-6">
      <div className="flex gap-2">
        <Input placeholder="Key" value={key} onChange={(e) => setKey(e.target.value)} />
        <Input placeholder="Value" value={val} onChange={(e) => setVal(e.target.value)} />
        <Button
          onClick={() => {
            if (!key) return;
            createVar.mutate({ serviceId: app.id, payload: { key, value: val, isSecret: false } });
            setKey('');
            setVal('');
          }}
        >
          Add
        </Button>
      </div>

      <div className="space-y-2">
        {vars.length === 0 ? (
          <div className="text-center p-8 border border-dashed rounded-lg bg-gray-50/50">
            <p className="text-sm text-gray-500">No environment variables configured</p>
          </div>
        ) : (
          vars.map((v: any) => (
            <div key={v.id} className="flex items-center justify-between rounded border p-3 bg-white shadow-sm">
              <div className="font-mono text-sm">
                <span className="font-semibold text-primary">{v.key}</span>
                <span className="mx-2 text-gray-400">=</span>
                <span className="text-gray-600 truncate max-w-xs">{v.isSecret ? '********' : v.value}</span>
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => deleteVar.mutate({ serviceId: app.id, id: v.id })}
                disabled={deleteVar.isPending}
              >
                <Trash className="h-4 w-4 text-red-500" />
              </Button>
            </div>
          ))
        )}
      </div>

      {app.repositoryUrl && (
        <div className="mt-8 rounded-lg border p-4 bg-slate-50/50">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h3 className="text-sm font-medium flex items-center gap-2">
                <FileSearch className="h-4 w-4 text-primary" />
                Auto-fill from repository
              </h3>
              <p className="text-xs text-gray-500 mt-1">
                Scan your repository for .env.example files to prepopulate variables
              </p>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setScanEnabled(true)}
              disabled={isScanning || availableSuggestions.length > 0}
            >
              {isScanning ? (
                <Loader2 className="h-4 w-4 animate-spin mr-2" />
              ) : null}
              {availableSuggestions.length > 0 ? 'Scanned' : 'Scan Repository'}
            </Button>
          </div>

          {availableSuggestions.length > 0 && (
            <div className="space-y-2 mt-4 pt-4 border-t border-slate-200">
              {availableSuggestions.map((s: any) => (
                <div key={s.key} className="flex items-center justify-between bg-white rounded border border-slate-200 p-2 text-sm">
                  <div className="font-mono">
                    <span className="font-semibold text-slate-700">{s.key}</span>
                    {s.value && (
                      <>
                        <span className="mx-2 text-slate-400">=</span>
                        <span className="text-slate-500">{s.value}</span>
                      </>
                    )}
                  </div>
                  <Button
                    size="sm"
                    variant="secondary"
                    className="h-7 text-xs gap-1"
                    onClick={() => {
                      createVar.mutate({ serviceId: app.id, payload: { key: s.key, value: s.value, isSecret: false } });
                    }}
                  >
                    <Plus className="h-3 w-3" /> Add
                  </Button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
