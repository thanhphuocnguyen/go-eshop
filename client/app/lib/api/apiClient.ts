import { badRequestHandler } from '@/app/utils';
import { GenericResponse } from '../definitions/index';
import { refreshTokenAction } from '@/app/actions/auth';

type RequestOptions = {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  body?: Record<string, any> | FormData;
  headers?: Record<string, string>;
  authToken?: string;
  nextOptions?: RequestInit;
  retryOnUnauthorized?: boolean;
  req?: any; // For SSR cookie access
  res?: any;
  queryParams?: Record<string, any>; // Added queryParams option
};

// Helper function to serialize query parameters
function serializeQueryParams(params: Record<string, any>): string {
  if (!params || Object.keys(params).length === 0) return '';

  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value === null || value === undefined) return;

    if (Array.isArray(value)) {
      value.forEach((item) => {
        if (item !== null && item !== undefined) {
          searchParams.append(`${key}[]`, String(item));
        }
      });
    } else if (typeof value === 'object') {
      searchParams.append(key, JSON.stringify(value));
    } else {
      searchParams.append(key, String(value));
    }
  });

  return searchParams.toString();
}

export async function clientSideFetch<T = any>(
  endpoint: string,
  {
    method = 'GET',
    body,
    headers = {},
    nextOptions = {},
    retryOnUnauthorized = true,
    req,
    res,
    queryParams,
  }: RequestOptions = {}
): Promise<GenericResponse<T>> {
  // Build the URL with query parameters if provided
  let fullUrl = endpoint.startsWith('http')
    ? endpoint
    : `${process.env.NEXT_PUBLIC_API_BASE_URL || ''}${endpoint}`;

  // Append query parameters if they exist
  if (queryParams && Object.keys(queryParams).length > 0) {
    const queryString = serializeQueryParams(queryParams);
    fullUrl += fullUrl.includes('?') ? `&${queryString}` : `?${queryString}`;
  }

  const isFormData = body instanceof FormData;

  const token = localStorage.getItem('accessToken');
  const finalHeaders: Record<string, string> = {
    ...(isFormData ? {} : { 'Content-Type': 'application/json' }),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...headers,
  };

  const response = await fetch(fullUrl, {
    method,
    headers: finalHeaders,
    body: body && !isFormData ? JSON.stringify(body) : (body as BodyInit),
    ...nextOptions,
  });

  if (response.status === 401 && retryOnUnauthorized) {
    const newToken = await refreshTokenAction();
    if (newToken) {
      localStorage.setItem('accessToken', newToken);
      return clientSideFetch<T>(endpoint, {
        method,
        body,
        headers,
        nextOptions,
        retryOnUnauthorized: false,
        req,
        res,
      });
    } else {
      throw new Error('Unauthorized, redirecting to login');
    }
  }

  badRequestHandler(response);

  return response.json();
}
