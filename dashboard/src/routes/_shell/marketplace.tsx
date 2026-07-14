import { createFileRoute } from '@tanstack/react-router';
import { Box, Database, Globe, Server, Star } from 'lucide-react';

export const Route = createFileRoute('/_shell/marketplace')({
  component: MarketplacePage,
});

const TEMPLATES = [
  {
    id: 'node-express',
    name: 'Node.js Express',
    description: 'A minimal Node.js Express server template.',
    category: 'Framework',
    icon: <Server className="w-8 h-8 text-green-500" />,
    repo: 'https://github.com/vesslhq/vessl-examples.git',
    branch: 'main',
    rootDirectory: 'node-express',
  },
  {
    id: 'go-fiber',
    name: 'Go Fiber',
    description: 'A lightning fast Go web framework template.',
    category: 'Framework',
    icon: <Box className="w-8 h-8 text-blue-400" />,
    repo: 'https://github.com/vesslhq/vessl-examples.git',
    branch: 'main',
    rootDirectory: 'go-fiber',
  },
  {
    id: 'python-fastapi',
    name: 'Python FastAPI',
    description: 'A high performance Python API framework template.',
    category: 'Framework',
    icon: <Globe className="w-8 h-8 text-yellow-500" />,
    repo: 'https://github.com/vesslhq/vessl-examples.git',
    branch: 'main',
    rootDirectory: 'python-fastapi',
  },
  {
    id: 'ruby-sinatra',
    name: 'Ruby Sinatra',
    description: 'A lightweight Ruby web application template.',
    category: 'Framework',
    icon: <Box className="w-8 h-8 text-red-400" />,
    repo: 'https://github.com/vesslhq/vessl-examples.git',
    branch: 'main',
    rootDirectory: 'ruby-sinatra',
  },
  {
    id: 'php-basic',
    name: 'PHP Basic',
    description: 'A classic PHP web application template.',
    category: 'Framework',
    icon: <Globe className="w-8 h-8 text-indigo-400" />,
    repo: 'https://github.com/vesslhq/vessl-examples.git',
    branch: 'main',
    rootDirectory: 'php-basic',
  },
];

const MANAGED_RESOURCES = [
  {
    id: 'redis',
    name: 'Managed Redis',
    description: 'High-performance managed Redis instance (via Upstash/Vessl Cloud).',
    category: 'Database',
    icon: <Database className="w-8 h-8 text-red-500" />,
    isManaged: true,
  },
  {
    id: 'postgres',
    name: 'Managed PostgreSQL',
    description: 'Scalable serverless PostgreSQL database (via Neon/Vessl Cloud).',
    category: 'Database',
    icon: <Database className="w-8 h-8 text-blue-500" />,
    isManaged: true,
  },
];

function MarketplacePage() {
  return (
    <div className="flex-1 space-y-4 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <div>
          <h2 className="text-3xl font-bold tracking-tight text-purple-50">
            Marketplace & Templates
          </h2>
          <p className="text-muted-foreground mt-1">
            Deploy fullstack starter templates or provision managed resources with one click.
          </p>
        </div>
      </div>

      <div className="mt-8 space-y-8">
        {/* Templates Section */}
        <section>
          <div className="flex items-center gap-2 mb-4">
            <Star className="w-5 h-5 text-purple-400" />
            <h3 className="text-xl font-semibold text-purple-50">Starter Templates</h3>
          </div>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {TEMPLATES.map((tpl) => (
              <div
                key={tpl.id}
                className="group relative flex flex-col justify-between rounded-xl border border-white/10 bg-black/40 p-6 backdrop-blur-xl transition-all hover:border-purple-500/50 hover:bg-purple-500/5"
              >
                <div className="flex items-start gap-4">
                  <div className="rounded-lg bg-white/5 p-3">{tpl.icon}</div>
                  <div>
                    <h4 className="font-semibold text-purple-50">{tpl.name}</h4>
                    <span className="text-xs font-medium text-purple-300/80 uppercase tracking-wider">
                      {tpl.category}
                    </span>
                  </div>
                </div>
                <p className="mt-4 text-sm text-gray-400 line-clamp-2">{tpl.description}</p>
                <div className="mt-6">
                  <button
                    type="button"
                    className="w-full rounded-md bg-purple-600/20 px-4 py-2 text-sm font-medium text-purple-300 transition-colors hover:bg-purple-600/30 hover:text-purple-100"
                  >
                    Deploy Template
                  </button>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Managed Resources Section */}
        <section>
          <div className="flex items-center gap-2 mb-4">
            <Database className="w-5 h-5 text-purple-400" />
            <h3 className="text-xl font-semibold text-purple-50">Managed Resources</h3>
          </div>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {MANAGED_RESOURCES.map((res) => (
              <div
                key={res.id}
                className="group relative flex flex-col justify-between rounded-xl border border-white/10 bg-black/40 p-6 backdrop-blur-xl transition-all hover:border-blue-500/50 hover:bg-blue-500/5"
              >
                <div className="absolute -top-3 -right-3 rounded-full bg-blue-600 px-3 py-1 text-xs font-bold text-white shadow-lg">
                  Managed
                </div>
                <div className="flex items-start gap-4">
                  <div className="rounded-lg bg-white/5 p-3">{res.icon}</div>
                  <div>
                    <h4 className="font-semibold text-purple-50">{res.name}</h4>
                    <span className="text-xs font-medium text-blue-300/80 uppercase tracking-wider">
                      {res.category}
                    </span>
                  </div>
                </div>
                <p className="mt-4 text-sm text-gray-400 line-clamp-2">{res.description}</p>
                <div className="mt-6">
                  <button
                    type="button"
                    className="w-full rounded-md bg-blue-600/20 px-4 py-2 text-sm font-medium text-blue-300 transition-colors hover:bg-blue-600/30 hover:text-blue-100"
                  >
                    Provision Resource
                  </button>
                </div>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}
