'use server';

import { LoginPayload } from './types';

export async function login(currentState: boolean, formData: LoginPayload) {
  return true;
}
