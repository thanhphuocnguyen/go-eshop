export type GeneralCategoryModel = {
  id: string;
  name: string;
  description: string;
  slug: string;
  image_url: string;
  published: boolean;
  remarkable: boolean;
  created_at: Date;
  updated_at: Date;
};

export type ProductCategory = {
  id: string;
  name: string;
  image_url?: string;
  min_price: number;
  max_price: number;
  slug: string;
  variant_count: number;
};
