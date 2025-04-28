import useSWR from 'swr';
import { API_PATHS } from '../constants/api';
import { apiFetch } from '../apis/api';
import { GenericResponse, UserModel } from '../definitions';

export const useUser = (accessToken?: string) => {
  const { data, isLoading, mutate } = useSWR(
    accessToken ? [API_PATHS.USER, accessToken] : null,
    ([url, accessToken]) =>
      apiFetch<GenericResponse<UserModel>>(url, {
        method: 'GET',
        authToken: accessToken,
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
