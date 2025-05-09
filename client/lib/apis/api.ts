/* eslint-disable @typescript-eslint/no-explicit-any */
// lib/api.ts
import { getCookie, setCookie, deleteCookie } from 'cookies-next';

import { PUBLIC_API_PATHS } from '../constants/api';

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
      value.forEach(item => {
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

let isRefreshing = false;
let refreshPromise: Promise<string | null> | null = null;

async function getAccessToken(
  req?: any,
  res?: any
): Promise<string | undefined> {
  const token = await getCookie('access_token', { req, res })?.toString();
  return token;
}

function storeAccessToken(token: string, res?: any) {
  setCookie('access_token', token, {
    maxAge: 60 * 60, // 1 hour
    path: '/',
    ...(res ? { res } : {}),
  });
}

function clearAccessToken(res?: any) {
  deleteCookie('access_token', { path: '/', ...(res ? { res } : {}) });
}

async function refreshAccessToken(
  refreshToken: string,
  req?: any,
  res?: any
): Promise<string | null> {
  if (isRefreshing && refreshPromise) return refreshPromise;

  isRefreshing = true;
  refreshPromise = fetch(PUBLIC_API_PATHS.REFRESH_TOKEN, {
    method: 'POST',
    // credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${refreshToken}`,
    },
  })
    .then(async (res) => {
      if (!res.ok) throw new Error('Refresh failed');
      const data = await res.json();
      const newToken = data.accessToken;
      if (!newToken) throw new Error('No token returned');

      storeAccessToken(newToken, res);
      return newToken;
    })
    .catch(() => {
      clearAccessToken(res);
      return null;
    })
    .finally(() => {
      isRefreshing = false;
      refreshPromise = null;
    });

  return refreshPromise;
}

export async function apiFetch<T = any>(
  endpoint: string,
  {
    method = 'GET',
    body,
    headers = {},
    authToken,
    nextOptions = {},
    retryOnUnauthorized = true,
    req,
    res,
    queryParams,
  }: RequestOptions = {}
): Promise<T> {
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

  const token = authToken || (await getAccessToken(req, res));
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

  const refreshToken = await getCookie('refresh_token');

  if (
    response.status === 401 &&
    retryOnUnauthorized &&
    !token &&
    refreshToken
  ) {
    const newToken = await refreshAccessToken(req, res);
    if (newToken) {
      return apiFetch<T>(endpoint, {
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

  const contentType = response.headers.get('content-type');
  if (contentType?.includes('application/json')) {
    return response.json() as Promise<T>;
  }

  return response.text() as unknown as Promise<T>;
}
