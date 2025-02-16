'use server';

import { API_PUBLIC_PATHS } from '@/lib/constants/api';
import {
  LoginFormSchema,
  LoginFormState,
  LoginResponse,
  RefreshTokenResponse,
  SignupFormSchema,
  SignupFormState,
} from '@/lib/definitions/auth';
import { GenericResponse } from '@/lib/types';
import { cookies } from 'next/headers';
import { redirect, RedirectType } from 'next/navigation';

export async function login(state: LoginFormState, formData: FormData) {
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

  const data: GenericResponse<LoginResponse> = await fetch(
    process.env.NEXT_API_URL + API_PUBLIC_PATHS.LOGIN,
    {
      method: 'POST',
      body: JSON.stringify({
        password: password,
        username: username,
      }),
    }
  ).then((res) => res.json());

  if (!data.data) {
    formState.message = 'Invalid username or password';
    return formState;
  }
  const ck = await cookies();
  ck.set('refresh_token', data.data.refresh_token, {
    expires: new Date(data.data.refresh_token_expire_at),
    secure: true,
    httpOnly: true,
    value: data.data.refresh_token,
  });
  ck.set('session_id', data.data.session_id, {
    expires: new Date(data.data.token_expire_at),
    secure: true,
    httpOnly: true,
    value: data.data.session_id,
  });
  ck.set('token', data.data.token, {
    expires: new Date(data.data.token_expire_at),
    secure: true,
    httpOnly: true,
    sameSite: 'lax',
    value: data.data.token,
  });

  redirect('/', RedirectType.replace);
}

export async function logout() {
  (await cookies()).delete('refresh_token');
  (await cookies()).delete('session_id');
  (await cookies()).delete('token');
  redirect('/login');
}

export async function register(state: SignupFormState, formData: FormData) {
  const validatedFields = SignupFormSchema.safeParse({
    email: formData.get('email') as string,
    name: formData.get('name') as string,
    fullname: formData.get('fullname') as string,
    phone: formData.get('phone') as string,
    password: formData.get('password') as string,
  });
  const formState: SignupFormState = {};
  if (!validatedFields.success) {
    formState.errors = validatedFields.error.flatten().fieldErrors;
    return formState;
  }
  const { email, name, password } = validatedFields.data;

  await fetch('/auth/register', {
    method: 'POST',
    body: JSON.stringify({
      email: email,
      name: name,
      password: password,
    }),
  });

  redirect('/login');
}
