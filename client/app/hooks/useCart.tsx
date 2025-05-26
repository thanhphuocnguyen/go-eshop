import useSWR from 'swr';
import { toast } from 'react-toastify';
import { clientSideFetch } from '../lib/api/apiClient';
import { CartModel } from '../lib/definitions';
import { PUBLIC_API_PATHS } from '../lib/constants/api';

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
      dedupingInterval: 60000, // 1 minute deduplication
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
