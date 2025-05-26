import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { GeneralCategoryModel } from '@/app/lib/definitions';
import { toast } from 'react-toastify';
import useSWR from 'swr';

export function useBrands() {
  const { data, error, mutate } = useSWR(
    ADMIN_API_PATHS.BRANDS,
    (url) =>
      clientSideFetch<GeneralCategoryModel[]>(url, {
        method: 'GET',
      }),
    {
      refreshInterval: 0,
      revalidateOnFocus: false,
      onError: (error) => {
        toast.error(
          <div>
            Failed to fetch collections:
            <div>{JSON.stringify(error)}</div>
          </div>
        );
      },
    }
  );
  return {
    brands: data?.data,
    isLoading: !error && !data,
    isError: error,
    mutate,
  };
}
