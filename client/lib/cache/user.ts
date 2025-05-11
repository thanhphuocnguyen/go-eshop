'use server';

import { cache } from 'react';
import { apiFetch } from '../apis/api';
import { GenericResponse } from '../definitions/apiResponse';
import { UserModel } from '../definitions';
import { PUBLIC_API_PATHS } from '../constants/api';
import { cookies } from 'next/headers';

export const getUserCache = cache(async () => {
  const cookieStorage = await cookies();
  if (!cookieStorage.get('access_token')?.value) {
    return null;
  }
  const { data, error } = await apiFetch<GenericResponse<UserModel>>(
    PUBLIC_API_PATHS.USER,
    {
      method: 'GET',
      authToken: cookieStorage.get('access_token')?.value,
      nextOptions: {
        next: {
          tags: ['user'],
          revalidate: 0,
        },
      },
    }
  );
  if (error) {
    console.error(error);
    return null;
  }
  return data;
});
