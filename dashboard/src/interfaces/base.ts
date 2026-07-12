export interface BaseResponse<T = unknown> {
  status: "success" | "error" | "warning";
  message: string;
  data: T;
  path?: string;
  executionTime?: number;
}

export interface PaginatedData<T> {
  records: T[];
  total: number;
  page: number;
  totalPages: number;
}

export interface BaseEntity {
  id: string;
  createdAt: string;
  updatedAt: string;
}
