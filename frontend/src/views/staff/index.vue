<template>
  <div class="staff-page">
    <!-- 表格组件 -->
    <o-table ref="tableRef" :table-config="tableConfig">
      <!-- 渲染 slotColumns -->
      <template v-for="(item, index) in slotColumns" :key="index" #[item.slot!]="{ row }">
        <!-- 渲染 policy_list -->
        <span v-if="item.prop === 'policy_list'">
          <!-- 显式声明 policy 为 IPolicy 类型 -->
          {{ (row.policies as IPolicy[])?.map((policy: IPolicy) => policy.policyname).join('; ') || '-' }}
        </span>
        <!-- 渲染其他字段 -->
        <span v-else>
          {{ row[item.prop] || '-' }}
        </span>
      </template>

      <!-- 表格顶部按钮 -->
      <template #table-top>
        <el-button type="primary" @click="handleAdd">添加</el-button>
      </template>
    </o-table>
  </div>
</template>


<script lang="ts" setup>
import { ref } from 'vue'
import router from '@/router'
import { useuuidDelete } from '@/hooks'
import { Api } from '@/api/common/enum'
import { deleteStaffApi } from '@/api/staff'
import type { ITableConfig, TableInstance } from '@/types'
import type { IPolicy } from '@/api/policy/model'

// 配置表格信息
const tableConfig: ITableConfig = {
  api: `${Api.staff}`,
  headers: [
    { prop: 'staffname', label: '人员名称' },
    { prop: 'staffid', label: '接待人员 ID' },
    { prop: 'number', label: '电话号码' },
    { prop: 'email', label: '邮箱' },
    { prop: 'role', label: '角色' },
    { prop: 'policy_list', label: '策略列表', slot: 'policy_list' },
  ],
  operations: {
    width: 115,
    buttons: [
      {
        text: '编辑',
        show: true,
        click: ({ row }) => {
          router.push(`/staff/edit/${row.uuid}`)
        },
      },
      {
        text: '删除',
        type: 'danger',
        show: true,
        click: ({ row }) => {
          handleDelete(row.uuid)
        },
      },
    ],
  },
}

// 筛选带有 slot 的字段
const slotColumns = tableConfig.headers.filter((header) => header.slot)

// 表格实例
const tableRef = ref<TableInstance>()

// 添加人员
const handleAdd = () => {
  router.push('/staff/add')
}

// 删除人员
const handleDelete = (uuid: string) => {
  const { onDelete } = useuuidDelete()
  onDelete(deleteStaffApi, uuid, tableRef)
}
</script>