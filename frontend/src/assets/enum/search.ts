/**
 * @description: 每个列表页面的搜索项 id
 * key => 对应的每个页面route的 path 属性
 * value => 对应的每个列表页面获取搜索项的 id 值，由后台提供
 */

import type { ISearchItem } from '@/types'

const searchConfig: Record<string, ISearchItem[]> = {
  '/policy': [
    // {
    //   id: 'policy_name_text',
    //   name: '策略名称',
    //   type: 'text',
    //   hint: '请输入策略名称'
    // }
  ],
  '/staff': [
    // {
    //   id: 'policy_name_text',
    //   name: '策略名称',
    //   type: 'text',
    //   hint: '请输入策略名称'
    // }
  ],
  '/kf': [
    // {
    //   id: 'policy_name_text',
    //   name: '策略名称',
    //   type: 'text',
    //   hint: '请输入策略名称'
    // }
  ],
  '/llmapp': [
    // {
    //   id: 'policy_name_text',
    //   name: '策略名称',
    //   type: 'text',
    //   hint: '请输入策略名称'
    // }
  ],
  '/wecom/config': [
    // {
    //   id: 'policy_name_text',
    //   name: '策略名称',
    //   type: 'text',
    //   hint: '请输入策略名称'
    // }
  ],
  '/wecom/account': [
    // {
    //   id: 'policy_name_text',
    //   name: '策略名称',
    //   type: 'text',
    //   hint: '请输入策略名称'
    // }
  ],
  '/user': [
    {
      id: 'user_name_text',
      name: '用户名称',
      type: 'text',
      hint: '请输入用户名称'
    },
    {
      id: 'user_mobile',
      name: '手机号',
      type: 'auto_complete',
      hint: '请输入手机号'
    }
  ],
  '/role': [
    {
      id: 'role_name_text',
      name: '角色名称',
      type: 'text',
      hint: '请输入角色名称'
    },
    {
      id: 'role_status',
      name: '角色状态',
      type: 'select',
      hint: '请选择角色状态',
      multi_select: false,
      options: [
        { key: 'enable', value: '已启用' },
        { key: 'disabled', value: '已禁用' }
      ]
    }
  ]
}

export default searchConfig

