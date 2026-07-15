import type { BaseResponse } from './base';

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

export interface TableSchema {
  name: string;
  columns: ColumnSchema[];
}

export interface ColumnSchema {
  name: string;
  type: string;
  isNullable: boolean;
  isPrimary: boolean;
}

export interface ImportDatabaseRequest {
  sourceUrl: string;
}

export type TableRowPayload = Record<string, unknown>;

export type GetDatabasesResponse = BaseResponse<Database[]>;
export type GetDatabaseResponse = BaseResponse<Database>;
export type CreateDatabaseResponse = BaseResponse<Database>;
export type DatabaseQueryResponseType = BaseResponse<DatabaseQueryResponse>;
export type ListTablesResponse = BaseResponse<TableSchema[]>;
export type ImportDatabaseResponse = BaseResponse<void>;
export type DeleteDatabaseResponse = BaseResponse<void>;
