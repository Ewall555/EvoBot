import { Api } from '../common/enum'
import http from '@/http'

import type { IReceptionist } from './model/receptionistModel'

export const createReceptionistApi = (kfid: string, params: IReceptionist) => {
    return http.post(`${Api.wecomreceptionist}/${kfid}`, params);
}

export const getReceptionistListApi = (kfid: string) => {
    return http.get<IReceptionist[]>(`${Api.wecomreceptionist}/${kfid}`);
}

export const deleteReceptionistApi = (kfid: string, params: IReceptionist) => {
    return http.delete(`${Api.wecomreceptionist}/${kfid}`, params)
}