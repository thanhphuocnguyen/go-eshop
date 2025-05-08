// Main definitions index file - exports all application types

export * from './category';
export * from './apiResponse';
export * from './common';
export * from './product';
export * from './collection';
export * from './user';
export * from './auth';
export * from './brand';
export * from './form.type';
export * from './user';
export * from './cart';
export * from './checkout';
export * from './order'; // Only export from order.ts
export * from './image';

/**
 * Generic response structure for API calls
 */
export interface GenericResponse<T> {
  data: T;
  error?: {
    details: string;
  };
  pagination?: {
    totalPages: number;
    currentPage: number;
    totalItems: number;
  };
}
