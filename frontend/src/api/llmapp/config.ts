import http from '@/http'
import { Api } from '../common/enum'
import type { IConfig } from './model/configModel'


export const createConfigApi = (params: IConfig) => {
    return http.post(`${Api.llmappconfig}`, params);
};

export const getConfigListApi = () => {
    return http.get<IConfig[]>(`${Api.llmappconfig}`)
}

export const getConfigApi = (uuid: string) => {
    return http.get(`${Api.llmappconfig}/${uuid}`)
};

export const updateConfigApi = (params: IConfig) => {
    return http.patch(`${Api.llmappconfig}/${params.uuid}`, params)
}

export const deleteConfigApi = (uuid: string) => {
    return http.delete(`${Api.llmappconfig}/${uuid}`)
}