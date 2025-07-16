import { Api } from '../common/enum'
import http from '@/http'

import type { IStaff } from './model'

export const createStaffApi = (params: IStaff) => {
    return http.post(`${Api.staff}`, params);
};

export const getStaffListApi = () => {
    return http.get<IStaff[]>(`${Api.staff}`)
}

export const getStaffApi = (uuid: string) => {
    return http.get<IStaff>(`${Api.staff}/${uuid}`)
}

export const updateStaffApi = (params: IStaff) => {
    return http.patch(`${Api.staff}/${params.uuid}`, params)
}

export const deleteStaffApi = (uuid: string) => {
    return http.delete(`${Api.staff}/${uuid}`)
}