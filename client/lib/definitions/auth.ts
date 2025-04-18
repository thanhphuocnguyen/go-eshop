import { z } from 'zod';
import { UserModel } from '../definitions';

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
    username: z.string().min(3).max(255),
    email: z.string().email(),
    fullname: z.string().min(3).max(255),
    phone: z.string().min(10).max(10),
    password: z.string().min(6).max(255),
    confirmPassword: z.string().min(6).max(255),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
  });

export type SignupFormData = z.infer<typeof SignupFormSchema>;

export type SignupFormState =
  | {
      data?: SignupFormData;
      errors?: {
        email?: string[];
        fullname?: string[];
        phone?: string[];
        username?: string[];
        password?: string[];
        confirmPassword?: string[];
      };
      message?: string;
    }
  | undefined;

export const registerSchema = z
  .object({
    username: z.string().min(3).max(20),
    password: z.string().min(8).max(100),
    confirmPassword: z.string().min(8).max(100),
    email: z.string().email(),
    phone: z.string().min(10).max(15),
    fullname: z.string().min(3).max(100),
  })
  .superRefine((data) => {
    if (data.password !== data.confirmPassword) {
      return { confirmPassword: 'Passwords do not match' };
    }
    return {};
  });

export type RegisterForm = z.infer<typeof registerSchema>;
