import http from '@/http'
import type { IUser } from '../system/modal/userModel'
import type { IOptionProp } from '@/types'


// 获取当前用户信息
export const getUserInfo = () => {
  return http.get<IUser>('/auth/profile')
}

export const getTableData = (api: string, params: any) => {
  return http.get<any>(api, params)
}

// 获取指定类型 key 的下拉选项
export const getOptions = (key: string) => {
  return http.get<IOptionProp[]>(`/search/option/${key}`)
}

// 获取数据 autocomplete 下拉选项
export const getAutoComplete = (key: string, keyword: string) => {
  return http.get<any>(`search/autocomplete/${key}`, { keyword })
}
