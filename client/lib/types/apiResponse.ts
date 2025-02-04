export type GenericResponse<T> = {
  data: T;
  status: number;
  message: string;
};

export type ErrorResponse = {
  status: number;
  message: string;
};

export type GenericListResponse<T> = {
  data: T[];
  total: number;
  status: number;
  message: string;
};
