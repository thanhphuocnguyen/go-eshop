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
  discountType: z
    .object({
      id: z.enum(['percentage', 'fixed_amount']),
      name: z.enum(['Percentage', 'Fixed Amount']),
    })
    .transform((v) => v.id),
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
    .string().or(z.number())
    .transform((v) => (v === '' ? undefined : Number(v)))
    .nullish(),
  minPurchaseAmount: z
    .string().or(z.number())
    .transform((v) => (v === '' ? undefined : Number(v)))
    .nullish(),
  maxDiscountAmount: z
    .string()
    .or(z.number())
    .transform((v) => (v === '' ? undefined : Number(v)))
    .nullish(),
});

export const editDiscountSchema = createDiscountSchema.extend({
  products: z
    .array(
      z.object({
        id: z.string(),
        name: z.string(),
        price: z.number(),
      })
    )
    .optional()
    .transform((v) => (v ? v.map((e) => e.id) : [])),
  categories: z
    .array(
      z.object({
        id: z.string(),
        name: z.string(),
      })
    )
    .optional()
    .transform((v) => (v ? v.map((e) => e.id) : [])),
  users: z
    .array(
      z.object({
        id: z.string(),
        name: z.string(),
      })
    )
    .optional()
    .transform((v) => (v ? v.map((e) => e.id) : [])),
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
export type CreateDiscountFormData = z.input<typeof createDiscountSchema>;
export type CreateDiscountOutputData = z.output<typeof createDiscountSchema>;
export type EditDiscountFormData = z.input<typeof editDiscountSchema>;
export type EditDiscountOutputData = z.output<typeof editDiscountSchema>;
