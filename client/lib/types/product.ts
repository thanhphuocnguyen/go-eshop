export type ProductModel = {
  product_id: number;
  name: string;
  description: string;
  price: number;
  category_id: number;
  published: boolean;
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
