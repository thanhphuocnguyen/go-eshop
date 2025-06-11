import useSWR from 'swr';
import { PUBLIC_API_PATHS } from '../lib/constants/api';
import { clientSideFetch } from '../lib/api/apiClient';
import { UserModel } from '../lib/definitions';

export const useUser = () => {
  const { data, isLoading, mutate } = useSWR(
    PUBLIC_API_PATHS.GET_ME,
    (url) =>
      clientSideFetch<UserModel>(url, {
        method: 'GET',
      }).then((res) => {
        if (res.error) {
          throw res.error;
        }
        return res.data;
      }),
    {
      refreshInterval: 0,
      dedupingInterval: 60000, // 1 minute deduplication
      revalidateOnFocus: false,
      shouldRetryOnError: false,
    }
  );
  return {
    user: data,
    isLoading,
    mutateUser: mutate,
  };
};
