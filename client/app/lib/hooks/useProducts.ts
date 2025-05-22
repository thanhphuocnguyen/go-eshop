import useSWR from 'swr';
import { PUBLIC_API_PATHS } from '../constants/api';
import { clientSideFetch } from '../apis/apiClient';
import { ProductListModel } from '../definitions';

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
      ),
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
