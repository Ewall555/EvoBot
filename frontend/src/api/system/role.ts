import http from '@/http'
import { Api } from '../common/enum'
import type { IRole } from './modal/roleModel'

export function createRole(params: IRole) {
  return http.post(`${Api.role}`, params)
}

export function getRoleInfo(id: number) {
  return http.get(`${Api.role}/${id}`)
}

export function updateRole(params: IRole) {
  return http.patch(`${Api.role}/${params.id}`, params)
}

export function deleteRole(id: number) {
  return http.delete(`${Api.role}/${id}`)
}
