import { Api } from '../common/enum'
import http from '@/http'

import type { IAccount } from './model/accountModel'

export const getAccountListApi = () => {
    return http.get<IAccount>(`${Api.wecomaccount}`)
}

export const getAccountAddContactWayApi = (kfid: string) => {
    return http.get<string>(`${Api.wecomaccount}/${kfid}`)
}