import useSWR from 'swr';
import { API_PATHS } from '../constants/api';
import { apiFetch } from '../apis/api';
import { GenericResponse, UserModel } from '../definitions';

export const useUser = (isUser?: boolean) => {
  const { data, isLoading, mutate } = useSWR(
    isUser ? API_PATHS.USER : null,
    (url) =>
      apiFetch<GenericResponse<UserModel>>(url, {
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
    setUser: mutate,
  };
};
