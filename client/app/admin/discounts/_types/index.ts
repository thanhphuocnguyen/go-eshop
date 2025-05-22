import dayjs from 'dayjs';
import { z } from 'zod';

// Define types for products and categories
export interface ProductType {
  id: string;
  name: string;
  price: number;
}

export interface CategoryType {
  id: string;
  name: string;
}

// Define Zod schema for form validation
export const discountSchema = z.object({
  code: z.string().min(1, 'Discount code is required'),
  discountType: z.object({
    id: z.enum(['percentage', 'fixed_amount']),
    name: z.enum(['Percentage', 'Fixed Amount']),
  }),
  discountValue: z
    .number()
    .min(1, 'Value must be positive')
    .refine((val) => val > 0, 'Value must be greater than 0'),
  isActive: z.boolean(),
  startsAt: z.string().refine((date) => {
    return dayjs(date).isValid();
  }),
  expiresAt: z.string().refine((date) => {
    return dayjs(date).isValid();
  }),
  description: z.string().nullish(),
  usageLimit: z
    .string()
    .transform((v) => (v === '' ? undefined : Number(v)))
    .nullish(),
  minPurchaseAmount: z
    .string()
    .transform((v) => (v === '' ? undefined : Number(v)))
    .nullish(),
  maxDiscountAmount: z
    .string()
    .transform((v) => (v === '' ? undefined : Number(v)))
    .nullish(),
});

// TypeScript type derived from the schema
export type DiscountFormData = z.infer<typeof discountSchema>;
