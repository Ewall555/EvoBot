import http from '@/http'

export interface ILogin {
  username: string
  password: string
}

export const loginApi = (params: ILogin) => {
  return http.post<any>(`/auth/login`, params);
};

export const logOutApi = () => {
  return http.post<any>(`/auth/logout`);
};
