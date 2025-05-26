import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';
import useSWR from 'swr';

export default function useDiscountProducts(id: string) {
  const { data, isLoading } = useSWR(
    ADMIN_API_PATHS.DISCOUNT_PRODUCTS.replace(':id', id),
    (url) =>
      clientSideFetch<{ id: string; name: string; price: number }[]>(
        url,
        {}
      ).then((res) => {
        if (res.error) {
          throw new Error(res.error.details);
        }
        return res.data;
      }),
    {
      revalidateOnFocus: false,
      onError: (error) => {
        toast.error(`Failed to fetch discount products: ${error.message}`, {
          position: 'top-right',
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
          progress: undefined,
        });
      },
    }
  );

  return {
    productDiscounts: data,
    isLoading,
  };
}
