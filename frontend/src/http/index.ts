import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios'
import { TOOKENkEY, clearLocal, getToken } from '@/utils'
import { showMessage } from './status'

export interface ResultData<T> {
  code: number;
  message: string;
  data: T;
}
export const baseURL = import.meta.env.VITE_API_URL as string
const config = {
  baseURL,
  timeout: 50000,
  withCredentials: true,
};

class RequestHttp {
  service: AxiosInstance;
  public constructor(config: AxiosRequestConfig) {
    this.service = axios.create(config);
    this.service.interceptors.request.use(
      (config) => {
        const token = getToken()
        if (token && config && config.headers) {
          config.headers[TOOKENkEY] = `Bearer ${token}`;
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    );

    this.service.interceptors.response.use(
      (res) => {
        const { data } = res
        switch (data.code) {
          case 200:
            return data
          default:
            showMessage(data.message)
            return Promise.reject(data)
        }
      },
      (error) => {
        const { response } = error
        if (response) {
          if (response.status === 401) {
            clearLocal()
            window.location.href = '/login'
          } else if (response.message) {
            showMessage(response.message || response.status)
          }
        }
        return Promise.reject(error)
      }
    )
  }
  get<T>(url: string, params?: object, _object = {}): Promise<ResultData<T>> {
    return this.service.get(url, { params, ..._object });
  }
  post<T>(url: string, params?: object, timeout?: number): Promise<ResultData<T>> {
    return this.service.post(url, params, {
      baseURL: import.meta.env.VITE_API_URL as string,
      timeout: timeout ? timeout : 50000,
      withCredentials: true,
    });
  }
  put<T>(url: string, params?: object, _object = {}): Promise<ResultData<T>> {
    return this.service.put(url, params, _object);
  }
  patch<T>(url: string, params?: object, _object = {}): Promise<ResultData<T>> {
    return this.service.patch(url, params, _object);
  }
  delete<T>(url: string, data?: object , _object = {}): Promise<ResultData<T>> {
    return this.service.delete(url, { data, ..._object } );
  }
  download<BlobPart>(url: string, params?: object, _object = {}): Promise<BlobPart> {
    return this.service.post(url, params, _object);
  }
  upload<T>(url: string, params: object = {}, config?: AxiosRequestConfig): Promise<T> {
    return this.service.post(url, params, config);
  }
}

export default new RequestHttp(config);
