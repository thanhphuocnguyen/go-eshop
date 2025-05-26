import { z } from 'zod';
import { BaseOptionSchema } from './common';
import { GeneralCategoryModel } from './category';

export type ProductListModel = {
  id: string;
  name: string;
  description: string;
  variantCount: number;
  slug: string;
  imageUrl?: string;
  minPrice: number;
  maxPrice: number;
  avgRating: number;
  reviewCount: number;
  sku: string;
  createdAt: Date;
  updatedAt: Date;
};

export type CategoryProductModel = {
  id: number;
  name: string;
  description: string;
  variantCount: number;
  imageUrl: string;
  priceFrom: number;
  priceTo: number;
  discountTo: number;
  createdAt: string;
};

export interface ProductCreateBody {
  name: string;
  description: string;
  price: number;
  discount?: number | null;
  stock: number;
  sku: string;
  slug: string;
  categoryId?: number;
  collectionId?: number;
  brandId?: number;
  attributes: {
    attributeId: number;
    valueIds: number[];
  }[];
}

export const BaseAttributeValueFormSchema = z.object({
  name: z.string().optional(),
  code: z.string().optional(),
  displayOrder: z.number().optional(),
  isActive: z.boolean().optional(),
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
    valueObject: BaseAttributeValueFormSchema.extend({
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
  stockQty: z.coerce.number().gte(0),
  sku: z.string().optional().readonly(),
  weight: z.coerce
    .number()
    .transform((v) => {
      if (!v) return null;
      return v;
    })
    .nullish(),
  isActive: z.boolean(),
  attributes: ProductVariantAttributeFormSchema.extend({
    id: z.string().uuid().optional(),
  }).array(),
});

export const ProductFormSchema = z.object({
  productInfo: z.object({
    name: z.string().min(3).max(100),
    description: z.string().min(10).max(5000),
    shortDescription: z.string().optional(),
    attributes: z.string().uuid().array(),
    price: z.coerce.number().gt(0),
    sku: z.string().nonempty(),
    slug: z.string().nonempty(),
    isActive: z.boolean(),
    category: BaseOptionSchema,
    brand: BaseOptionSchema,
    collection: BaseOptionSchema.nullish(),
    images: z
      .object({
        id: z.string().uuid().optional(),
        url: z.string(),
        role: z.string().nullish(),
        assignments: z.string().array(),
        isRemoved: z.boolean().optional(),
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
  displayOrder: number;
  createdAt: Date;
  isActive: boolean;
  outOfStock?: boolean;
};

export type AttributeDetailModel = {
  id: string;
  name: string;
  values: AttributeValueDetailModel[];
  createdAt: Date;
};

export type ProductVariantAttributeModel = Omit<
  AttributeDetailModel,
  'values'
> & {
  valueObject: AttributeValueDetailModel;
};

export type VariantDetailModel = {
  attributes: ProductVariantAttributeModel[];
  id: string;
  price: number;
  stockQty: number;
  sku: string;
  weight: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
};

export type AssignmentImageModel = {
  id: string;
  entityId: string;
  entityType: string;
  displayOrder: number;
  role: string;
};

export type ProductImageModel = {
  id: string;
  externalId: string;
  url: string;
  role: string;
  assignments: AssignmentImageModel[];
};

export type ProductDetailModel = {
  id: string;
  name: string;
  description: string;
  shortDescription: string;
  attributes: string[];
  slug: string;
  isActive: boolean;
  price: number;
  sku: string;
  category: GeneralCategoryModel;
  collection: GeneralCategoryModel;
  brand: GeneralCategoryModel;
  published: boolean;
  variants: VariantDetailModel[];
  maxDiscountValue?: number; // max discount value
  discountType?: string; // 'percentage' | 'fixed'
  productImages: ProductImageModel[];
  createdAt: string; // date
  updatedAt: string; // date
  ratingCount: number;
  oneStarCount: number;
  twoStarCount: number;
  threeStarCount: number;
  fourStarCount: number;
  fiveStarCount: number;
};
