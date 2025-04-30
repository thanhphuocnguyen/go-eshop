import { apiFetch } from '@/lib/apis/api';
import { ADMIN_API_PATHS } from '@/lib/constants/api';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import { toast } from 'react-toastify';
import useSWR from 'swr';

export function useCollections() {
  const { data, error } = useSWR(
    ADMIN_API_PATHS.COLLECTIONS,
    (url) =>
      apiFetch<GenericResponse<GeneralCategoryModel[]>>(url, {}).then((res) => {
        if (res.error) {
          throw new Error(res.error.stack);
        }
        return res.data;
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
    collections: data,
    isLoading: !error && !data,
    isError: error,
  };
}
