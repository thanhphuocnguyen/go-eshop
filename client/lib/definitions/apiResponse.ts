export type GenericResponse<T> = {
  success: boolean;
  message: string;
  data: T;
  pagination?: Pagination;
  meta?: Meta;
  error?: ErrorResponse | null;
};

export type ErrorResponse = {
  code: number;
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
  totalItems: number;
};
