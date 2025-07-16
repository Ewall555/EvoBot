<template>
  <div class="staff-page">
    <!-- 表格组件 -->
    <o-table ref="tableRef" :table-config="tableConfig">
      <!-- 渲染 slotColumns -->
      <template v-for="(item, index) in slotColumns" :key="index" #[item.slot!]="{ row }">
        <!-- 渲染 llmapp_list -->
        <span v-if="item.prop === 'llmapp_list'">
          <!-- 显式声明 llmapp 为 IConfig 类型 -->
          {{ (row.llmappconfigs as IConfig[])?.map((llmappconfig: IConfig) => llmappconfig.config_name).join('; ') || '-' }}
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
import { deleteConfigApi } from '@/api/llmapp/config'
import type { ITableConfig, TableInstance } from '@/types'
import type { IConfig } from '@/api/llmapp/model/configModel'

// 配置表格信息
const tableConfig: ITableConfig = {
  api: `${Api.llmappconfig}`,
  headers: [
    { prop: 'config_name', label: '名称' },
    { prop: 'llmapp_type', label: '类型' },
    { prop: 'base_url', label: 'base url' },
    { prop: 'api_key', label: 'api key' },
  ],
  operations: {
    width: 115,
    buttons: [
      {
        text: '编辑',
        show: true,
        click: ({ row }) => {
          router.push(`/llmapp/edit/${row.uuid}`)
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

// 添加
const handleAdd = () => {
  router.push('/llmapp/add')
}

// 删除
const handleDelete = (uuid: string) => {
  const { onDelete } = useuuidDelete()
  onDelete(deleteConfigApi, uuid, tableRef)
}
</script>