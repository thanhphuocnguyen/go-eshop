import axios from 'axios';
import { cookies } from 'next/headers';

const axiosInstance = axios.create({
  baseURL: process.env.NEXT_API_URL!, // Replace with your API base URL
  timeout: 1000, // Set a timeout if needed
  headers: { 'Content-Type': 'application/json' },
});

// You can also add interceptors if needed
axiosInstance.interceptors.request.use(
  async (config) => {
    // Do something before request is sent
    const pasetoToken = (await cookies()).get('token');
    if (pasetoToken) {
      config.headers['Authorization'] = `Bearer ${pasetoToken}`;
    }
    return config;
  },
  (error) => {
    // Do something with request error
    return Promise.reject(error);
  }
);

axiosInstance.interceptors.response.use(
  (response) => {
    // Do something with response data
    return response;
  },
  (error) => {
    // Do something with response error
    return Promise.reject(error);
  }
);

export default axiosInstance;
