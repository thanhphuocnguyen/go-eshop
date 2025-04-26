import { API_PATHS } from '@/lib/constants/api';
import { RefreshTokenResponse } from '@/lib/definitions/auth';
import { NextRequest, NextResponse } from 'next/server';
import { GenericResponse } from './lib/definitions';

const AdminPath = '/admin';
export async function middleware(request: NextRequest) {
  const privatePaths = ['/dashboard', '/profile'];
  const path = request.nextUrl.pathname;

  if (
    path.startsWith(AdminPath) &&
    request.cookies.get('role')?.value !== 'admin'
  ) {
    return NextResponse.redirect(new URL('/not-found', request.nextUrl));
  }
  const isProtectedRoute = privatePaths.some((route) => path.startsWith(route));

  if (isProtectedRoute && !request.cookies.get('access_token')?.value) {
    return NextResponse.redirect(new URL('/login', request.nextUrl));
  }

  if (request.cookies.get('access_token')?.value) {
    return NextResponse.next();
  }
  if (!request.cookies.get('refresh_token')?.value) {
    return NextResponse.next();
  }

  const response = NextResponse.next();
  try {
    const refreshResult = await fetch(API_PATHS.REFRESH_TOKEN, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${request.cookies.get('refresh_token')?.value}`,
      },
    });
    if (!refreshResult.ok) {
      console.error(refreshResult.statusText);
      return response;
    }
    const { data }: GenericResponse<RefreshTokenResponse> =
      await refreshResult.json();
    response.cookies.set('access_token', data.access_token, {
      expires: new Date(data.access_token_expires_at),
    });
  } catch (error) {
    console.error(error);
  }

  return response;
}
