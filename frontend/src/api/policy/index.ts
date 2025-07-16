import { Api } from '../common/enum'
import http from '@/http'

import type { IPolicy } from './model'

export const createPolicyApi = (params: IPolicy) => {
    return http.post(`${Api.policy}`, params);
};

export const getPolicyListApi = () => {
    return http.get<IPolicy[]>(`${Api.policy}`);
};

export const getPolicyApi = (uuid: string) => {
    return http.get<IPolicy>(`${Api.policy}/${uuid}`)
}

export const updatePolicyApi = (params: IPolicy) => {
    return http.patch(`${Api.policy}/${params.uuid}`, params)
}

export const deletePolicyApi = (uuid: string) => {
    return http.delete(`${Api.policy}/${uuid}`)
}