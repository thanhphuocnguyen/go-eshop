import { z } from 'zod';
import { UserModel } from '../types/user';

export const LoginFormSchema = z.object({
  username: z.string().min(3).max(255),
  password: z.string().min(6).max(255),
});

export type LoginFormData = z.infer<typeof LoginFormSchema>;

export type LoginResponse = {
  session_id: string;
  token: string;
  token_expire_at: Date;
  refresh_token: string;
  refresh_token_expire_at: Date;
  user: UserModel;
};

export type LoginFormState =
  | {
      errors?: {
        username?: string[];
        password?: string[];
      };
      message?: string;
    }
  | undefined;
