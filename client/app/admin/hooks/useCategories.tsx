import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { GeneralCategoryModel, GenericResponse } from '@/app/lib/definitions';
import { toast } from 'react-toastify';
import useSWR from 'swr';

export const useCategories = () => {
  const { data, error } = useSWR<GenericResponse<GeneralCategoryModel[]>>(
    ADMIN_API_PATHS.CATEGORIES,
    (url) => clientSideFetch<GeneralCategoryModel[]>(url),
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
    categories: data?.data,
    isLoading: !error && !data,
    isError: error,
  };
};
