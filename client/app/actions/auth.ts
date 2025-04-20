'use server';

import { API_PATHS } from '@/lib/constants/api';
import {
  LoginFormSchema,
  LoginFormState,
  LoginResponse,
  SignupFormSchema,
  SignupFormState,
} from '@/lib/definitions/auth';
import { GenericResponse } from '@/lib/definitions';
import { cookies } from 'next/headers';
import { redirect, RedirectType } from 'next/navigation';
import { apiFetch } from '@/lib/api/api';

export async function login(_: LoginFormState, formData: FormData) {
  const validatedFields = LoginFormSchema.safeParse({
    username: formData.get('username') as string,
    password: formData.get('password') as string,
  });
  const formState: LoginFormState = {};
  if (!validatedFields.success) {
    formState.errors = validatedFields.error.flatten().fieldErrors;
    return formState;
  }
  const { username, password } = validatedFields.data;

  const { data, error } = await apiFetch<GenericResponse<LoginResponse>>(
    API_PATHS.LOGIN,
    {
      method: 'POST',
      body: {
        password: password,
        username: username,
      },
    }
  );

  if (error) {
    formState.message = 'Invalid username or password';
    return formState;
  }

  if (!data) {
    formState.message = 'Invalid username or password';
    return formState;
  }
  const cookieStorage = await cookies();
  cookieStorage.set('token', data.token, {
    expires: new Date(data.token_expire_at),
    value: data.token,
  });

  cookieStorage.set('refresh_token', data.refresh_token, {
    expires: new Date(data.refresh_token_expire_at),
    value: data.refresh_token,
  });

  cookieStorage.set('session_id', data.session_id, {
    expires: new Date(data.refresh_token_expire_at),
    value: data.session_id,
  });
  cookieStorage.set('user_role', JSON.stringify(data.user), {
    expires: new Date(data.refresh_token_expire_at),
    value: data.user.role,
  });
  cookieStorage.set('user_id', JSON.stringify(data.user), {
    expires: new Date(data.refresh_token_expire_at),
    value: data.user.id,
  });
  cookieStorage.set('user_name', JSON.stringify(data.user), {
    expires: new Date(data.refresh_token_expire_at),
    value: data.user.fullname,
  });

  cookieStorage.set('user_email', JSON.stringify(data.user), {
    expires: new Date(data.refresh_token_expire_at),
    value: data.user.username,
  });

  redirect('/', RedirectType.replace);
}

export async function logout() {
  const cookieStorage = await cookies();
  cookieStorage.delete('token');
  cookieStorage.delete('refresh_token');
  cookieStorage.delete('session_id');
  cookieStorage.delete('user_role');
  cookieStorage.delete('user_id');
  cookieStorage.delete('user_name');
  redirect('/login');
}

export async function register(_: SignupFormState, formData: FormData) {
  const validatedFields = SignupFormSchema.safeParse({
    email: formData.get('email') as string,
    username: formData.get('username') as string,
    fullname: formData.get('fullname') as string,
    phone: formData.get('phone') as string,
    password: formData.get('password') as string,
    confirmPassword: formData.get('confirmPassword') as string,
  });

  const formState: SignupFormState = {};

  if (!validatedFields.success) {
    formState.errors = validatedFields.error.flatten().fieldErrors;
    formState.data = validatedFields.data;
    formState.message = 'Invalid form data';
    return formState;
  }

  const { email, username, password, fullname, phone } = validatedFields.data;
  console.log('Sending request');
  const resp = await apiFetch(API_PATHS.REGISTER, {
    method: 'POST',
    body: {
      email,
      username,
      password,
      fullname,
      phone,
    },
  });
  const data = await resp.json();
  if (!resp.ok) {
    console.error(data);
    formState.message = 'Failed to register';
    return formState;
  }

  redirect('/login');
}
