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

export type RefreshTokenResponse = {
  access_token: string;
  access_token_expires_at: string;
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

export const SignupFormSchema = z
  .object({
    email: z.string().email(),
    fullname: z.string().min(3).max(255),
    phone: z.string().min(10).max(10),
    name: z.string().min(3).max(255),
    password: z.string().min(6).max(255),
    confirmPassword: z.string().min(6).max(255),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
  });

export type SignupFormData = z.infer<typeof SignupFormSchema>;

export type SignupFormState =
  | {
      errors?: {
        email?: string[];
        fullname?: string[];
        phone?: string[];
        name?: string[];
        password?: string[];
        confirmPassword?: string[];
      };
      message?: string;
    }
  | undefined;
