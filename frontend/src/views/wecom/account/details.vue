<template>
    <div class="receptionist-details">
        <!-- 接待人员列表 -->
        <o-table ref="tableRef" :table-config="tableConfig">
            <template v-for="(item, index) in slotColumns" :key="index" #[item.slot!]="{ row }">
                <span v-if="item.prop === 'status'">
                    <!-- 根据 status 显示不同状态 -->
                    <span :class="{ 'status-active': row[item.prop] === 0, 'status-disabled': row[item.prop] === 1 }">
                        {{ row[item.prop] === 0 ? '接待中' : '停止接待' }}
                    </span>
                </span>
                <span v-else>
                    <!-- 处理接待人员字段 -->
                    <template v-if="item.prop === 'userid'">
                        {{ row[item.prop] }}
                    </template>
                </span>
            </template>
            <!-- 表格顶部按钮 -->
            <template #table-top>
                <el-button type="primary" @click="handleAdd">添加</el-button>
                <div v-if="addContactWayUrl" class="add-contact-way-url">
                    客服地址:<a :href="addContactWayUrl" target="_blank" class="link-button">
                        {{ addContactWayUrl }}
                    </a>
                </div>
            </template>
        </o-table>
    </div>
</template>
<script lang="ts" setup>
import { Api } from '@/api/common/enum'
import { onMounted, ref } from 'vue'
import { usekfidDelete } from '@/hooks'
import type { ITableConfig, TableInstance } from '@/types'
import { useRoute } from 'vue-router';
import { createReceptionistApi, deleteReceptionistApi } from '@/api/wecom/receptionist';
import type { IReceptionist } from '@/api/wecom/model/receptionistModel'
import { ElMessage, ElMessageBox } from 'element-plus';
import { getAccountAddContactWayApi } from '@/api/wecom/account';

const route = useRoute()
const tableRef = ref<TableInstance>()
const kfid = route.params.kfid as string
const addContactWayUrl = ref<string | null>(null)

// 表格配置
const tableConfig: ITableConfig = {
    api: `${Api.wecomreceptionist}/${kfid}`,
    headers: [
        { prop: 'userid', label: '接待人员ID' },
        { prop: 'status', label: '状态', slot: 'status' }
    ],
    operations: {
        width: 120,
        buttons: [
            {
                text: '删除',
                type: 'danger',
                show: true,
                click: ({ row }) => {
                    handleDelete(kfid, row.userid)
                },
            },
        ]
    }
}

const slotColumns = tableConfig.headers.filter((header) => header.slot)

const handleDelete = (kfid: string, userid: string) => {
    const { onDelete } = usekfidDelete()
    const params: IReceptionist = { userid_list: [userid] }
    onDelete(deleteReceptionistApi, kfid, params, tableRef)
}
const handleAdd = () => {
    ElMessageBox.prompt('请输入接待人员ID列表，多个ID用逗号分隔', '添加接待人员', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
    }).then(({ value }) => {
        if (!value) {
            ElMessage.error('请输入有效的接待人员ID列表');
            return;
        }

        const useridList = value.split(',').map(id => id.trim());
        if (useridList.length === 0 || !useridList.every(id => id)) {
            ElMessage.error('请输入有效的接待人员ID列表');
            return;
        }

        const params: IReceptionist = { userid_list: useridList };
        createReceptionistApi(kfid, params).then(() => {
            ElMessage.success('接待人员添加成功');
            // 刷新表格数据
            tableRef.value?.loadTableData();
        }).catch((error) => {
            ElMessage.error(`添加接待人员失败: ${error.message}`);
        });
    }).catch(() => {
        ElMessage.info('已取消添加');
    });
}
const fetchAddContactWayUrl = async () => {
    try {
        const response = await getAccountAddContactWayApi(kfid);
        const url = new URL(response.data);
        addContactWayUrl.value = url.origin + url.pathname;
    } catch (error) {
    }
}

// 组件挂载时调用
onMounted(() => {
    fetchAddContactWayUrl();
})
</script>
