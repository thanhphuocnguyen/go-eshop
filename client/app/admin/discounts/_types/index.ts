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
export const createDiscountSchema = z.object({
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

export const editDiscountSchema = createDiscountSchema.extend({
  products: z.array(z.string()).optional(),
  categories: z.array(z.string()).optional(),
  users: z.array(z.string()).optional(),
});

export const discountTypes = z.enum(['percentage', 'fixed_amount']);
export const discountTypeNames = z.enum(['Percentage', 'Fixed Amount']);

export const discountTypeOptions = [
  {
    id: discountTypes.Enum.percentage,
    name: discountTypeNames.Enum.Percentage,
  },
  {
    id: discountTypes.Enum.fixed_amount,
    name: discountTypeNames.Enum['Fixed Amount'],
  },
];

// TypeScript type derived from the schema
export type CreateDiscountFormData = z.infer<typeof createDiscountSchema>;
export type EditDiscountFormData = z.infer<typeof editDiscountSchema>;
