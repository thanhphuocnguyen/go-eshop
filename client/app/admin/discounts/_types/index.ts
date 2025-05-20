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
  description: z.string().optional(),
  discountType: z.enum(['percentage', 'fixed_amount']),
  discountValue: z
    .number()
    .min(0, 'Value must be positive')
    .refine((val) => val > 0, 'Value must be greater than 0'),
  minPurchaseAmount: z.number().nullable().optional(),
  maxDiscountAmount: z.number().nullable().optional(),
  usageLimit: z.number().nullable().optional(),
  isActive: z.boolean(),
  startsAt: z.string().min(1, 'Start date is required'),
  expiresAt: z.date(),
  products: z.array(
    z.object({
      id: z.string(),
      name: z.string(),
      price: z.number(),
    })
  ),
  categories: z.array(
    z.object({
      id: z.string(),
      name: z.string(),
    })
  ),
});

// TypeScript type derived from the schema
export type DiscountFormData = z.infer<typeof discountSchema>;
