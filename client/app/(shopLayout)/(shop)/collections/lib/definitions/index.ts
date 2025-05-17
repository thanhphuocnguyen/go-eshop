import { CollectionDetailModel } from '@/app/lib/definitions';
export type FilterOption = {
  id: string;
  name: string;
};

export type GetCollectionBySlugModel = {
  collection: CollectionDetailModel;
  categories: FilterOption[];
  brands: FilterOption[];
  attributes: Record<string, FilterOption[]>;
};
