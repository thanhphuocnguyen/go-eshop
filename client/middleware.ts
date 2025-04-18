import { API_PATHS } from '@/lib/constants/api';
import { RefreshTokenResponse } from '@/lib/definitions/auth';
import { NextRequest, NextResponse } from 'next/server';
import { GenericResponse } from './lib/definitions';

export async function middleware(request: NextRequest) {
  const path = request.nextUrl.pathname;

  const isProtectedRoute = path.startsWith('/admin');

  if (isProtectedRoute && request.cookies.get('user_role')?.value !== 'admin') {
    return NextResponse.redirect(new URL('/not-found', request.nextUrl));
  }

  if (request.cookies.get('token')?.value) {
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
    response.cookies.set('token', data.access_token, {
      expires: new Date(data.access_token_expires_at),
      sameSite: 'lax',
      httpOnly: true,
    });
  } catch (error) {
    console.error(error);
  }

  return response;
}
