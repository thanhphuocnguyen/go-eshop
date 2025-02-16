import { API_PUBLIC_PATHS } from '@/lib/constants/api';
import { RefreshTokenResponse } from '@/lib/definitions/auth';
import { NextRequest, NextResponse } from 'next/server';

export async function middleware(request: NextRequest) {
  if (request.cookies.get('token')?.value) {
    return NextResponse.next();
  }
  if (!request.cookies.get('refresh_token')?.value) {
    return NextResponse.next();
  }
  const response = NextResponse.next();
  try {
    const data: RefreshTokenResponse = await fetch(
      process.env.NEXT_API_URL + API_PUBLIC_PATHS.REFRESH_TOKEN,
      {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${
            request.cookies.get('refresh_token')?.value
          }`,
        },
      }
    ).then((res) => res.json());
    response.cookies.set('token', data.access_token);
  } catch (error) {
    console.error(error);
  }

  return response;
}
