import { z } from 'zod';
import { BaseOptionSchema } from './common';
import { GeneralCategoryModel } from './category';

export type ProductListModel = {
  id: number;
  name: string;
  description: string;
  variant_count: number;
  slug: string;
  image_url: string;
  min_price: number;
  max_price: number;
  sku: string;
  created_at: Date;
  updated_at: Date;
};

export type CategoryProductModel = {
  id: number;
  name: string;
  description: string;
  variant_count: number;
  image_url: string;
  price_from: number;
  price_to: number;
  discount_to: number;
  created_at: string;
};

export interface ProductCreateBody {
  name: string;
  description: string;
  price: number;
  discount?: number | null;
  stock: number;
  sku: string;
  slug: string;
  category_id?: number;
  collection_id?: number;
  brand_id?: number;
  attributes: {
    attribute_id: number;
    value_ids: number[];
  }[];
}

export const BaseAttributeValueFormSchema = z.object({
  name: z.string().optional(),
  code: z.string().optional(),
  display_order: z.number().optional(),
  is_active: z.boolean().optional(),
});

export const BaseAttributeFormSchema = z.object({
  name: z.string().optional(),
});

export const AttributeFormSchema = BaseAttributeFormSchema.extend({
  values: BaseAttributeValueFormSchema.extend({
    id: z.string().optional(),
  })
    .array()
    .min(1, {
      message: 'At least one value is required',
    }),
});

export const ProductVariantAttributeFormSchema = BaseAttributeFormSchema.extend(
  {
    value_object: BaseAttributeValueFormSchema.extend({
      id: z
        .string()
        .uuid({
          message: 'Value is required',
        })
        .optional(),
    }),
  }
);

export const VariantFormSchema = z.object({
  id: z.string().uuid().optional(),
  price: z.coerce.number().gt(0),
  stock_qty: z.coerce.number().gte(0),
  sku: z.string().optional().readonly(),
  weight: z.coerce
    .number()
    .transform((v) => {
      if (!v) return null;
      return v;
    })
    .nullish(),
  is_active: z.boolean(),
  attributes: ProductVariantAttributeFormSchema.extend({
    id: z.string().uuid().optional(),
  }).array(),
});

export const ProductFormSchema = z.object({
  product_info: z.object({
    name: z.string().min(3).max(100),
    description: z.string().min(10).max(5000),
    short_description: z.string().optional(),
    attributes: z.string().uuid().array(),
    price: z.coerce.number().gt(0),
    sku: z.string().nonempty(),
    slug: z.string().nonempty(),
    is_active: z.boolean(),
    category: BaseOptionSchema,
    brand: BaseOptionSchema,
    collection: BaseOptionSchema.nullish(),
    images: z
      .object({
        id: z.string().uuid().optional(),
        url: z.string(),
        role: z.string().nullish(),
        assignments: z.string().array(),
        is_removed: z.boolean().optional(),
      })
      .array(),
  }),
  variants: z.array(VariantFormSchema),
});

export type ProductModelForm = z.infer<typeof ProductFormSchema>;
export type VariantModelForm = z.infer<typeof VariantFormSchema>;
export type AttributeFormModel = z.infer<typeof AttributeFormSchema>;
export type ProductVariantAttributeFormModel = z.infer<
  typeof ProductVariantAttributeFormSchema
>;

export type AttributeValueFormModel = z.infer<
  typeof BaseAttributeValueFormSchema
>;

export type AttributeValueDetailModel = {
  id: string;
  name: string;
  code: string;
  display_order: number;
  created_at: Date;
  is_active: boolean;
  out_of_stock?: boolean;
};

export type AttributeDetailModel = {
  id: string;
  name: string;
  values: AttributeValueDetailModel[];
  created_at: Date;
};

export type ProductVariantAttributeModel = Omit<
  AttributeDetailModel,
  'values'
> & {
  value_object: AttributeValueDetailModel;
};

export type VariantDetailModel = {
  attributes: ProductVariantAttributeModel[];
  id: string;
  price: number;
  stock_qty: number;
  sku: string;
  weight: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type AssignmentImageModel = {
  id: number;
  entity_id: string;
  entity_type: string;
  display_order: number;
  role: string;
};

export type ProductImageModel = {
  id: string;
  external_id: string;
  url: string;
  role: string;
  assignments: AssignmentImageModel[];
};

export type ProductDetailModel = {
  id: string;
  name: string;
  description: string;
  short_description: string;
  attributes: string[];
  slug: string;
  is_active: boolean;
  price: number;
  sku: string;
  category: GeneralCategoryModel;
  collection: GeneralCategoryModel;
  brand: GeneralCategoryModel;
  published: boolean;
  variants: VariantDetailModel[];
  product_images: ProductImageModel[];
  created_at: string; // date
  updated_at: string; // date
};
