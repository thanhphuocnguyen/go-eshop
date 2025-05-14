import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { RefreshTokenResponse } from '@/app/lib/definitions/auth';
import { NextRequest, NextResponse } from 'next/server';
import { GenericResponse } from './app/lib/definitions';
import { jwtDecode, JwtPayload } from 'jwt-decode';
import { revalidateTag } from 'next/cache';

export type JwtModel = {
  role: string;
  username: string;
  user_id: string;
  id: string;
} & JwtPayload;

const AdminPath = '/admin';
export async function middleware(request: NextRequest) {
  const privatePaths = ['/profile', '/checkout', '/cart', '/orders'];
  const path = request.nextUrl.pathname;
  const accessToken = request.cookies.get('access_token')?.value;
  const refreshToken = request.cookies.get('refresh_token')?.value;
  const response = NextResponse.next();
  if (!accessToken && refreshToken) {
    const refreshResult = await fetch(PUBLIC_API_PATHS.REFRESH_TOKEN, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${refreshToken}`,
      },
    });
    if (!refreshResult.ok) {
      const response = NextResponse.redirect(
        new URL('/login', request.nextUrl)
      );
      response.cookies.delete('access_token');
      response.cookies.delete('refresh_token');
      response.cookies.delete('session_id');
      return response;
    }

    const { data, error }: GenericResponse<RefreshTokenResponse> =
      await refreshResult.json();
    if (error) {
      request.cookies.delete('access_token');
      request.cookies.delete('refresh_token');
      request.cookies.delete('session_id');
      return NextResponse.redirect(new URL('/login', request.nextUrl));
    }

    response.cookies.set('access_token', data.access_token, {
      expires: new Date(data.access_token_expires_at),
    });
    revalidateTag('user');
  }
  if (path.startsWith(AdminPath) && accessToken) {
    const decode = jwtDecode<JwtModel>(accessToken || '');
    if (decode['role'] !== 'admin') {
      return NextResponse.redirect(new URL('/not-found', request.nextUrl));
    }
  }
  const isProtectedRoute = privatePaths.some((route) => path.startsWith(route));

  if (isProtectedRoute && !request.cookies.get('access_token')?.value) {
    return NextResponse.redirect(new URL('/login', request.nextUrl));
  }

  return response;
}
