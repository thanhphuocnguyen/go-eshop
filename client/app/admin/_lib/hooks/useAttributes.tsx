import useSWR from 'swr';
import { toast } from 'react-toastify';
import { apiFetch } from '@/lib/api/api';
import { API_PATHS } from '@/lib/constants/api';
import { AttributeDetailModel, GenericResponse } from '@/lib/definitions';

export function useAttributes(ids: number[] = [], productId?: number) {
  const { data, isLoading, error } = useSWR(
    [API_PATHS.ATTRIBUTES, ids, productId],
    ([url]) => {
      return apiFetch<GenericResponse<AttributeDetailModel[]>>(
        `${url}?ids[]=${ids.join(',')}`,
        {}
      ).then((res) => {
        if (res.error) {
          throw new Error(res.error.stack, {
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
