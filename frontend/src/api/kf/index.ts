import { Api } from '../common/enum'
import http from '@/http'

import type { IKF } from './model'

export const createKFApi = (params: IKF) => {
    return http.post(`${Api.kf}`, params);
};

export const getKFListApi = () => {
    return http.get<IKF[]>(`${Api.kf}`);
};

export const getKFApi = (uuid: string) => {
    return http.get<IKF>(`${Api.kf}/${uuid}`)
}

export const updateKFApi = (params: IKF) => {
    return http.patch(`${Api.kf}/${params.uuid}`, params)
}

export const deleteKFApi = (uuid: string) => {
    return http.delete(`${Api.kf}/${uuid}`)
}