import { z } from 'zod';

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
