import { CategoryProductModel } from './product';

export type Category = {
  category_id: number;
  name: string;
  description: string;
  sort_order: number;
  published: boolean;
  created_at: Date;
  updated_at: Date;
  products: CategoryProductModel[];
};
