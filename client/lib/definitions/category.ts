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
  products?: ProductCategory[];
};

export type ProductCategory = {
  id: string;
  name: string;
  image_url?: string;
  description: string;
  slug: string;
};
