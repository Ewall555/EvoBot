<template>
    <div class="account-list">
        <o-table ref="tableRef" :table-config="tableConfig">
            <template v-for="(item, index) in slotColumns" :key="index" #[item.slot!]="{ row }">
                <template v-if="item.prop === 'avatar'">
                    <img :src="row[item.prop]" alt="Avatar" class="avatar" />
                </template>
                <template v-else>
                    {{ row[item.prop] || '-' }}
                </template>
            </template>
        </o-table>
    </div>
</template>
<script lang="ts" setup>
import { Api } from '@/api/common/enum'
import { ref } from 'vue'
import router from '@/router'
import type { ITableConfig, TableInstance } from '@/types'

const tableRef = ref<TableInstance>()

// 表格配置
const tableConfig: ITableConfig = {
    api: `${Api.wecomaccount}`,
    headers: [
        { prop: 'avatar', label: '头像', slot: 'avatar' },
        { prop: 'name', label: '名称' },
        { prop: 'open_kfid', label: 'Open KFID' }
    ],
    operations: {
        width: 120,
        buttons: [
            {
                text: '查看详情',
                show: true,
                click: ({ row }) => handleViewDetails(row.open_kfid)
            }
        ]
    }
}

const slotColumns = tableConfig.headers.filter((header) => header.slot)

// 加载账号列表
const accountList = ref<any[]>([])

// 跳转详情页面
const handleViewDetails = (kfid: string) => {
    router.push(`/wecom/account/details/${kfid}`)
}


</script>
<style scoped>
.avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
}
</style>
