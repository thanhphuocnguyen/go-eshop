import { z } from 'zod';

export type BrandRequest = {
  name: string;
  slug: string;
  image_url: string | null;
  description: string;
};

export const BrandFormSchema = z.object({
  brand_id: z.string().optional(),
  name: z.string().nonempty(),
  description: z.string().nonempty(),
  image_url: z.string().optional(),
});
