import useSWR from 'swr';
import { toast } from 'react-toastify';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { AttributeDetailModel } from '@/app/lib/definitions';

export function useAttributes(ids: number[] = [], productId?: number) {
  const { data, isLoading, error } = useSWR(
    [ADMIN_API_PATHS.ATTRIBUTES, ids, productId],
    ([url]) => {
      return clientSideFetch<AttributeDetailModel[]>(
        `${url}?ids[]=${ids.join(',')}`,
        {}
      ).then((res) => {
        if (res.error) {
          throw new Error(res.error.details, {
            cause: res.error,
          });
        }
        return res.data;
      });
    },
    {
      revalidateOnFocus: false,
      revalidateOnReconnect: false,
      onError: (error) => {
        toast.error(
          <div>
            Failed to fetch attributes:
            <div>{JSON.stringify(error)}</div>
          </div>
        );
      },
    }
  );
  return {
    attributes: data,
    attributesLoading: isLoading,
    attributesError: error,
  };
}
