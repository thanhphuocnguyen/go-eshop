import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { ProductDetailModel } from '@/app/lib/definitions';
import { toast } from 'react-toastify';
import useSWR from 'swr';

export function useProductDetail(slug: string) {
  const {
    data: productDetail,
    isLoading,
    mutate,
  } = useSWR(
    slug ? PUBLIC_API_PATHS.PRODUCT_DETAIL.replace(':id', slug) : null,
    (url) =>
      clientSideFetch<ProductDetailModel>(url).then((data) => data.data),
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
