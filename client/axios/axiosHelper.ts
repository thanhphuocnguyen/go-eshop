import { AxiosRequestConfig } from 'axios';
import axiosInstance from './instance';

export const axiosGet = async (url: string, config?: AxiosRequestConfig) => {
  try {
    const response = await axiosInstance.get(url, config);
    return response;
  } catch (error) {
    return error;
  }
};

export const axiosPost = async (
  url: string,
  data?: any,
  config?: AxiosRequestConfig
) => {
  try {
    const response = await axiosInstance.post(url, data, config);
    return response;
  } catch (error) {
    return error;
  }
};

export const axiosPut = async (
  url: string,
  data?: any,
  config?: AxiosRequestConfig
) => {
  try {
    const response = await axiosInstance.put(url, data, config);
    return response;
  } catch (error) {
    return error;
  }
};

export const axiosDelete = async (url: string, config?: AxiosRequestConfig) => {
  try {
    const response = await axiosInstance.delete(url, config);
    return response;
  } catch (error) {
    return error;
  }
};

export const axiosPatch = async (
  url: string,
  data?: any,
  config?: AxiosRequestConfig
) => {
  try {
    const response = await axiosInstance.patch(url, data, config);
    return response;
  } catch (error) {
    return error;
  }
};
