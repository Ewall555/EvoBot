import http from '@/http'
import { Api } from '../common/enum'
import type { IUser } from './modal/userModel'

export function createUser(params: IUser) {
  return http.post(`${Api.user}`, params)
}

export function getUserInfo(id: number) {
  return http.get(`${Api.user}/${id}`)
}

export function updateUser(params: IUser) {
  return http.patch(`${Api.user}/${params.id}`, params)
}

export function deleteUser(id: number) {
  return http.delete(`${Api.user}/${id}`)
}
