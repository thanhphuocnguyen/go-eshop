import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

class HttpClient {
  private instance: AxiosInstance;
  constructor(baseURL: string) {
    this.instance = axios.create({
      baseURL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Interceptors for request
    this.instance.interceptors.request.use(
      (config) => {
        // Add auth token or modify headers if needed
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Interceptors for response
    this.instance.interceptors.response.use(
      (response) => response.data,
      (error) => {
        // Handle global errors
        return Promise.reject(error);
      }
    );
  }

  get<T>(
    path: string,
    params: Record<string, string | number>,
    config: AxiosRequestConfig = {}
  ) {
    return this.instance.get<T, T>(path, { params, ...config });
  }

  post<D, T>(path: string, data: D, config: AxiosRequestConfig = {}) {
    return this.instance.post<T, T>(path, data, { ...config });
  }

  put<D, T>(path: string, data: D, config: AxiosRequestConfig = {}) {
    return this.instance.put<T, T>(path, data, { ...config });
  }

  patch<D, T>(path: string, data: D, config: AxiosRequestConfig = {}) {
    return this.instance.patch<T, T>(path, data, { ...config });
  }

  delete<T>(path: string, config: AxiosRequestConfig = {}) {
    return this.instance.delete<T, T>(path, { ...config });
  }
}

// Example usage
const apiClient = new HttpClient(process.env.NEXT_API_URL as string);

export default apiClient;
