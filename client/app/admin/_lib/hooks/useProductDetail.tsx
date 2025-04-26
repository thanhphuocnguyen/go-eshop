import { apiFetch } from '@/lib/apis/api';
import { API_PATHS } from '@/lib/constants/api';
import { GenericResponse, ProductDetailModel } from '@/lib/definitions';
import { toast } from 'react-toastify';
import useSWR from 'swr';

export function useProductDetail(slug: string | number) {
  const {
    data: productDetail,
    isLoading,
    mutate,
  } = useSWR(
    API_PATHS.PRODUCT_DETAIL.replace(':id', slug),
    (url) =>
      apiFetch<GenericResponse<ProductDetailModel>>(url).then(
        (data) => data.data
      ),
    {
      refreshInterval: 0,
      revalidateOnFocus: false,
      onError: (err) => {
        toast.error(
          <div>
            Failed to fetch product detail:
            <div>{JSON.stringify(err)}</div>
          </div>
        );
      },
    }
  );
  return {
    productDetail,
    isLoading,
    mutate,
  };
}
