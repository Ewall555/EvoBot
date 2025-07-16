<template>
  <o-table ref="tableRef" :table-config="tableConfig">
    <template v-for="(item, index) in slotColumns" :key="index" #[item.slot!]="{ row }">
      <!-- 处理字段展示 -->
      <span v-if="row[item.prop]">{{ row[item.prop] }}</span>
      <span v-else>-</span>
    </template>

    <template #table-top>
      <span>配置列表</span>
    </template>
  </o-table>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Api } from '@/api/common/enum'
import type { ITableConfig, TableInstance } from '@/types'

// 获取 Vue Router
const router = useRouter()

// 配置表格展示
const tableConfig: ITableConfig = {
  api: `${Api.wecomconfig}`, // 获取配置数据的接口路径
  headers: [
    {
      prop: 'type',
      label: '类型'
    },
    {
      prop: 'agent_id',
      label: 'Agent ID'
    },
    {
      prop: 'corp_id',
      label: '企业 ID'
    },
    {
      prop: 'encoding_aes_key',
      label: 'Encoding AES Key'
    },
    {
      prop: 'token',
      label: 'Token'
    }
  ],
  operations: {
    width: 80, // 操作列宽度
    buttons: [
      {
        text: '编辑',
        show: true,
        click: ({ row }) => {
          // 编辑按钮点击事件
          router.push(`/wecom/config/edit/${row.uuid}`)
        }
      }
    ]
  }
}

// 用于过滤表格中的slot列
const slotColumns = tableConfig.headers.filter((header) => header.slot)

const tableRef = ref<TableInstance>()
</script>