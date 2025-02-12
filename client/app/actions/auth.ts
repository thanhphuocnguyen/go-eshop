import apiClient from '@/axios/axios';
import {
  LoginFormData,
  LoginFormSchema,
  LoginFormState,
  LoginResponse,
} from '@/lib/definitions/auth';
import { AxiosError } from 'axios';
import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';

export async function login(state: LoginFormState, formData: FormData) {
  const validatedFields = LoginFormSchema.safeParse({
    name: formData.get('username') as string,
    password: formData.get('password') as string,
  });
  const formState: LoginFormState = {};
  if (!validatedFields.success) {
    formState.errors = validatedFields.error.flatten().fieldErrors;
    return formState;
  }
  const { username, password } = validatedFields.data;

  try {
    const data = await apiClient.post<LoginFormData, LoginResponse>(
      '/auth/login',
      {
        password: password,
        username: username,
      }
    );
    (await cookies()).set('refresh_token', data.refresh_token, {
      expires: data.refresh_token_expire_at,
      secure: true,
      value: data.refresh_token,
    });
    (await cookies()).set('session_id', data.session_id, {
      expires: data.token_expire_at,
      secure: true,
      value: data.session_id,
    });
    (await cookies()).set('token', data.token, {
      expires: data.token_expire_at,
      secure: true,
      value: data.token,
    });
    redirect('/dashboard');
  } catch (error) {
    const errorResponse = error as AxiosError<any>;
    formState.message = errorResponse.response?.data.message;
    return formState;
  }
}
