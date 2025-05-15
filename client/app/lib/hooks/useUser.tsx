import useSWR from 'swr';
import { PUBLIC_API_PATHS } from '../constants/api';
import { apiFetchClientSide } from '../apis/apiClient';
import { UserModel } from '../definitions';

export const useUser = () => {
  const { data, isLoading, mutate } = useSWR(
    PUBLIC_API_PATHS.GET_ME,
    (url) =>
      apiFetchClientSide<UserModel>(url, {
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
    }
  );
  return {
    user: data,
    isLoading,
    mutateUser: mutate,
  };
};
