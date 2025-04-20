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
  price: number;
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
  value: z.string().optional(),
  display_value: z.string().optional(),
  display_order: z.number().optional(),
  is_active: z.boolean().optional(),
});

export const BaseAttributeFormSchema = z.object({
  name: z.string(),
});

export const AttributeFormSchema = BaseAttributeFormSchema.extend({
  values: BaseAttributeValueFormSchema.array().min(1, {
    message: 'At least one value is required',
  }),
});

export const ProductVariantAttributeFormSchema = BaseAttributeFormSchema.extend(
  {
    value: BaseAttributeValueFormSchema.extend({
      id: z.number().min(1, {
        message: 'Value is required',
      }),
    }).nullish(),
  }
);

export const VariantFormSchema = z.object({
  id: z.string().optional(),
  price: z.coerce.number().gt(0),
  stock: z.coerce.number().gte(0),
  weight: z.preprocess((value) => {
    if (value === '') return undefined;
    return Number(value);
  }, z.number().gte(0).nullish()),
  is_active: z.boolean().default(true),
  attributes: z
    .array(
      ProductVariantAttributeFormSchema.extend({
        id: z.number().optional(),
      })
    )
    .min(1),
  sku: z.string().readonly().optional(),
});

export const ProductFormSchema = z.object({
  id: z.string().optional(),
  name: z.string().min(3).max(100),
  description: z.string().min(10).max(5000),
  price: z.coerce.number().gt(0),
  sku: z.string().nonempty(),
  slug: z.string().nonempty(),
  is_active: z.boolean().optional().default(true),
  category: BaseOptionSchema,
  collection: BaseOptionSchema.nullish(),
  images: z
    .array(
      z.object({
        is_deleted: z.boolean().optional(),
        id: z.number().optional(),
        image_url: z.string(),
      })
    )
    .nullish(),
  brand: BaseOptionSchema.nullish(),
  variants: VariantFormSchema.array(),
  removed_images: z.array(z.number()).optional(),
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
  id: number;
  value: string;
  display_value: string;
  display_order: number;
  created_at: Date;
  is_active: boolean;
};
export type AttributeDetailModel = {
  id: number;
  name: string;
  created_at: Date;
  values: AttributeValueDetailModel[];
};

export type ProductVariantAttributeModel = Omit<
  AttributeDetailModel,
  'values'
> & {
  value: AttributeValueDetailModel;
};

export type VariantDetailModel = {
  attributes: ProductVariantAttributeModel[];
  id: string;
  price: number;
  stock_qty: number;
  sku: string;
  weight: number;
  is_active: boolean;
  image_url: string;
  image_id: number;
  created_at: string;
  updated_at: string;
};

export type AssignmentImageModel = {
  id: number;
  entity_id: number;
  entity_type: string;
  display_order: number;
  role: string;
};

export type ProductImageModel = {
  id: number;
  external_id: string;
  url: string;
  mime_type?: string;
  file_size?: number;
  entity_id?: number;
  entity_type?: string;
  display_order?: number;
  role?: string;
};

export type VariantImageMode = {
  assignments: AssignmentImageModel[];
  id: number;
  external_id: string;
  url: string;
};

export type ProductDetailModel = {
  id: string;
  name: string;
  description: string;
  slug: string;
  is_active: boolean;
  price: number;
  sku: string;
  category: GeneralCategoryModel;
  collection: GeneralCategoryModel;
  brand: GeneralCategoryModel;
  published: boolean;
  variants: VariantDetailModel[];
  images: ProductImageModel[];
  variant_images: VariantImageMode[];
  created_at: string; // date
  updated_at: string; // date
};
