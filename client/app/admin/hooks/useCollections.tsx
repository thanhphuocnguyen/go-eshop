import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { GeneralCategoryModel } from '@/app/lib/definitions';
import { toast } from 'react-toastify';
import useSWR from 'swr';

export function useCollections() {
  const { data, error } = useSWR(
    PUBLIC_API_PATHS.COLLECTIONS,
    (url) =>
      clientSideFetch<GeneralCategoryModel[]>(url, {}).then((res) => {
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
