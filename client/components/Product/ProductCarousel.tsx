import apiClient from '@/axios/axios';
import { GenericResponse } from '@/lib/types';
import { ProductModel } from '@/lib/types/product';

export default async function ProductCarousel({
  categoryID,
}: {
  categoryID: number;
}) {
  await apiClient.get<GenericResponse<ProductModel>>(
    `/products?categoryID=${categoryID}`,
    {}
  );
}
