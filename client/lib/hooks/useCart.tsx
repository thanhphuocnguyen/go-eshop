import useSWR from 'swr';
import { API_PATHS } from '../constants/api';
import { apiFetch } from '../apis/api';
import { GenericResponse } from '../definitions';
import { CartModel } from '../definitions/cart';
import { toast } from 'react-toastify';

export const useCart = (userId?: string) => {
  const { data, isLoading, error, mutate } = useSWR(
    userId ? [API_PATHS.CART, userId] : null,
    ([url]) =>
      apiFetch<GenericResponse<CartModel>>(url, {
        method: 'GET',
      }).then((res) => {
        if (res.error) {
          throw res.error;
        }
        return res.data;
      }),
    {
      refreshInterval: 0,
      revalidateOnFocus: false,
      onError: (error) => {
        toast.error(
          <div>
            Failed to fetch cart:
            <div>{JSON.stringify(error)}</div>
          </div>
        );
      },
    }
  );
  return {
    cart: data,
    cartLoading: isLoading,
    isError: error,
    mutateCart: mutate,
  };
};
