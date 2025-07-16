import { Api } from '../common/enum'
import http from '@/http'

import type { IConfig } from './model/configModel'

// export const createConfigApi = (params: IConfig) => {
//     return http.post(`${Api.wecom}${Api.config}`, params);
// };

// export const getConfigListApi = () => {
//     return http.get<IConfig[]>(`${Api.wecom}${Api.config}`);
// };

export const getConfigApi = (uuid: string) => {
    return http.get<IConfig>(`${Api.wecomconfig}/${uuid}`)
}

export const updateConfigApi = (params: IConfig) => {
    return http.patch(`${Api.wecomconfig}/${params.uuid}`, params)
}

// export const deleteConfigApi = (uuid: string) => {
//     return http.delete(`${Api.wecom}${Api.config}/${uuid}`)
// }