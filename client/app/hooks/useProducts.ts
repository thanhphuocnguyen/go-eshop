import useSWR from 'swr';
import { PUBLIC_API_PATHS } from '../lib/constants/api';
import { clientSideFetch } from '../lib/api/apiClient';
import { ProductListModel } from '../lib/definitions';

export const useProducts = ({
  page,
  limit,
  debouncedSearch,
}: {
  page: number;
  limit: number;
  debouncedSearch: string;
}) => {
  const { data, isLoading, mutate } = useSWR(
    [PUBLIC_API_PATHS.PRODUCTS, page, limit, debouncedSearch],
    ([url, page, limit, search]) =>
      clientSideFetch<ProductListModel[]>(
        `${url}?page=${page}&pageSize=${limit}&search=${search}`,
        {}
      ).then((res) => {
        if (res.error) {
          throw new Error(res.error.details);
        }
        const { data, pagination } = res;
        return {
          data,
          pagination,
        };
      }),
    {
      revalidateOnFocus: false,
      onError: (err) => {
        throw err;
      },
    }
  );
  return {
    isLoading,
    mutate,
    products: data?.data,
    pagination: data?.pagination,
  };
};
