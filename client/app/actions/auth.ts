'use server';

import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import {
  GenericResponse,
  LoginResponse,
  RefreshTokenResponse,
} from '@/app/lib/definitions';
import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';

export async function refreshTokenAction() {
  const cookieStore = await cookies();
  const refreshToken = cookieStore.get('refreshToken')?.value;
  if (!refreshToken) {
    redirect('/login');
  }
  const response = await fetch(PUBLIC_API_PATHS.REFRESH_TOKEN, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${refreshToken}`,
    },
  });
  if (!response.ok) {
    throw new Error('Failed to refresh token');
  }
  const data = (await response.json()) as GenericResponse<RefreshTokenResponse>;

  if (data.error) {
    cookieStore.delete('accessToken');
    cookieStore.delete('refreshToken');
    cookieStore.delete('sessionId');
    redirect('/login');
  }
  if (data) {
    const { accessToken, accessTokenExpiresAt } = data.data;
    cookieStore.set('accessToken', accessToken, {
      expires: new Date(accessTokenExpiresAt),
    });
    return accessToken;
  }
}

export async function loginAction(data: FormData) {
  const result = await fetch(PUBLIC_API_PATHS.LOGIN, {
    method: 'POST',
    body: data,
  });
  if (!result.ok) {
    return {
      error: {
        message: 'Login failed',
        status: result.status,
      },
    };
  } else {
    const { data } = (await result.json()) as GenericResponse<LoginResponse>;
    if (data) {
      const {
        accessToken,
        accessTokenExpiresIn,
        refreshToken,
        refreshTokenExpiresAt,
        sessionId,
      } = data;
      const cookieStore = await cookies();
      cookieStore.set('accessToken', accessToken, {
        expires: new Date(accessTokenExpiresIn),
      });
      cookieStore.set('refreshToken', refreshToken, {
        expires: new Date(refreshTokenExpiresAt),
      });
      cookieStore.set('sessionId', sessionId, {
        expires: new Date(refreshTokenExpiresAt),
      });
      return {
        accessToken,
        refreshToken,
      };
    } else {
      return {
        error: {
          message: 'Login failed',
          status: 401,
        },
      };
    }
  }
}

export async function logoutAction() {
  const cookieStore = await cookies();
  cookieStore.delete('accessToken');
  cookieStore.delete('refreshToken');
  cookieStore.delete('sessionId');
}
