import { GeneralCategoryModel } from './category';
import { ProductDetailModel } from './product';

export type CollectionDetailModel = GeneralCategoryModel & {
  products: ProductDetailModel[];
};
