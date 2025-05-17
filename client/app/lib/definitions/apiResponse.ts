// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type GenericResponse<T = any> = {
  success: boolean;
  message: string;
  data: T;
  pagination?: Pagination;
  meta?: Meta;
  error?: ErrorResponse | null;
};

export type ErrorResponse = {
  code: string;
  details: string;
  stack: string;
};

export type Meta = {
  timestamp: string;
  requestId: string;
  path: string;
  method: string;
};

export type Pagination = {
  page: number;
  pageSize: number;
  totalPages: number;
  total: number;
  hasNextPage: boolean;
  hasPreviousPage: boolean;
};
