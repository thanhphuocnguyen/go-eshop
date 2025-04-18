import { API_PATHS } from '@/lib/constants/api';
import useSWR from 'swr';
import { toast } from 'react-toastify';
import { BaseOption, GenericResponse } from '@/lib/definitions';
import { apiFetch } from '@/lib/api/api';

export interface AttributeResponse extends BaseOption {
  id: number;
  name: string;
  values: {
    id: number;
    value: string;
  }[];
  // created_at: string;
}
export function useAttributes() {
  const { data, isLoading, error } = useSWR(
    API_PATHS.ATTRIBUTES,
    (url) => {
      return apiFetch<GenericResponse<AttributeResponse[]>>(url, {}).then(
        (res) => {
          if (res.error) {
            throw new Error(res.error.stack);
          }
          return res.data;
        }
      );
    },
    {
      refreshInterval: 0,
      revalidateOnFocus: false,
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
