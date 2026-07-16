import type { ReactNode } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';

export interface ResourceColumn<T> {
  key: string;
  label: string;
  render: (row: T) => ReactNode;
}

interface ResourceTableProps<T> {
  columns: ResourceColumn<T>[];
  rows: T[];
  getRowKey: (row: T) => string;
  caption?: string;
}

export function ResourceTable<T>({ columns, rows, getRowKey, caption }: ResourceTableProps<T>) {
  return (
    <div className="overflow-hidden rounded-lg border bg-card">
      <div className="flex items-center justify-between border-b bg-muted/20 px-3 py-2">
        <p className="text-muted-foreground text-xs">
          {rows.length.toLocaleString()} {rows.length === 1 ? 'record' : 'records'}
        </p>
        <p className="text-muted-foreground text-xs">Keyboard focus enabled</p>
      </div>
      <Table>
        {caption ? <caption className="sr-only">{caption}</caption> : null}
        <TableHeader>
          <TableRow className="bg-muted/40 hover:bg-muted/40">
            {columns.map((column) => (
              <TableHead key={column.key} className="border-r last:border-r-0">
                {column.label}
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((row) => (
            <TableRow key={getRowKey(row)} tabIndex={0} className="focus-visible:bg-muted/50">
              {columns.map((column) => (
                <TableCell key={column.key} className="border-r last:border-r-0">
                  {column.render(row)}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
