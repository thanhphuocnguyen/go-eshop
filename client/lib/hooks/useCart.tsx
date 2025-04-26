import useSWR from 'swr';
import { API_PATHS } from '../constants/api';
import { apiFetch } from '../apis/api';
import { GenericResponse } from '../definitions';
import { CartModel } from '../definitions/cart';
import { toast } from 'react-toastify';

export const useCart = () => {
  const { data, error } = useSWR(
    API_PATHS.CART,
    (url) =>
      apiFetch<GenericResponse<CartModel[]>>(url, {
        method: 'GET',
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
    cart: data?.data,
    isLoading: !error && !data,
    isError: error,
  };
};
