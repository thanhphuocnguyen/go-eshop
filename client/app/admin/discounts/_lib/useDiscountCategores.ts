import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';
import useSWR from 'swr';

export default function useDiscountCategories(id: string) {
  const { data, isLoading } = useSWR(
    ADMIN_API_PATHS.DISCOUNT_CATEGORIES.replaceAll(':id', id),
    (url) =>
      clientSideFetch<{ id: string; name: string }[]>(url, {}).then((res) => {
        if (res.error) {
          throw new Error(res.error.details);
        }
        return res.data;
      }),
    {
      revalidateOnFocus: false,
      onError: (error) => {
        toast.error(`Failed to fetch discount categories: ${error.message}`, {
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
    categoryDiscounts: data,
    isLoading,
  };
}
