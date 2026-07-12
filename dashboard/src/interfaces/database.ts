export interface Database {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  engine: string;
  version: string;
  port: number;
  username: string;
  password: string;
  databaseName: string;
  volumePath: string;
  containerId: string;
  status: string;
  internalDns: string;
  externalDns: string;
  customArgs: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateDatabaseRequest {
  projectId: string;
  environmentId: string;
  name: string;
  engine: string;
  version: string;
  port: number;
  username: string;
  password: string;
  databaseName: string;
  volumePath: string;
  customArgs: string;
}

export interface Storage {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  type: string;
  apiPort: number;
  consolePort: number;
  accessKey: string;
  secretKey?: string;
  bucketName: string;
  volumePath: string;
  containerId: string;
  status: string;
  internalDns: string;
  externalDns: string;
  createdAt: string;
  updatedAt: string;
}

export interface DatabaseQueryRequest {
  query: string;
}

export interface DatabaseQueryResponse {
  columns?: string[];
  rows?: Record<string, unknown>[];
  result?: unknown;
}
