import { PUBLIC_API_PATHS } from '@/lib/constants/api';
import { RefreshTokenResponse } from '@/lib/definitions/auth';
import { NextRequest, NextResponse } from 'next/server';
import { GenericResponse } from './lib/definitions';
import { jwtDecode, JwtPayload } from 'jwt-decode';

type JwtMode = {
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
  if (path.startsWith(AdminPath) && accessToken) {
    const decode = jwtDecode<JwtMode>(accessToken || '');
    console.log({ decode });
    if (decode['role'] !== 'admin') {
      return NextResponse.redirect(new URL('/not-found', request.nextUrl));
    }
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
    const refreshResult = await fetch(PUBLIC_API_PATHS.REFRESH_TOKEN, {
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
