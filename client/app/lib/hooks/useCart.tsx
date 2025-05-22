import useSWR from 'swr';
import { PUBLIC_API_PATHS } from '../constants/api';
import { clientSideFetch } from '../apis/apiClient';
import { CartModel } from '../definitions/cart';
import { toast } from 'react-toastify';

export const useCart = (userId?: string) => {
  const { data, isLoading, error, mutate } = useSWR(
    userId ? [PUBLIC_API_PATHS.CART, userId] : null,
    ([url]) =>
      clientSideFetch<CartModel>(url, {
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
