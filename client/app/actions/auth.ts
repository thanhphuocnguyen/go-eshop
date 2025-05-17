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
  return fetch(PUBLIC_API_PATHS.REFRESH_TOKEN, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${cookieStore.get('refresh_token')?.value}`,
    },
  })
    .then<GenericResponse<RefreshTokenResponse>>((response) => {
      if (!response.ok) {
        throw new Error('Failed to refresh token');
      }
      return response.json();
    })
    .then((data) => {
      if (data.error) {
        cookieStore.delete('access_token');
        cookieStore.delete('refresh_token');
        cookieStore.delete('session_id');
        redirect('/login');
      }
      if (data) {
        const { access_token, access_token_expires_at } = data.data;
        cookieStore.set('access_token', access_token, {
          expires: new Date(access_token_expires_at),
        });
        return access_token;
      }
    });
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
        access_token,
        access_token_expires_in,
        refresh_token,
        refresh_token_expires_at,
        session_id,
      } = data;
      const cookieStore = await cookies();
      cookieStore.set('access_token', access_token, {
        expires: new Date(access_token_expires_in),
      });
      cookieStore.set('refresh_token', refresh_token, {
        expires: new Date(refresh_token_expires_at),
      });
      cookieStore.set('session_id', session_id, {
        expires: new Date(refresh_token_expires_at),
      });
      return {
        access_token,
        refresh_token,
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
  cookieStore.delete('access_token');
  cookieStore.delete('refresh_token');
  cookieStore.delete('session_id');
}
