<template>
  <div class="kf-page">
    <!-- 表格组件 -->
    <o-table ref="tableRef" :table-config="tableConfig">
      <!-- 渲染 slotColumns -->
      <template v-for="(item, index) in slotColumns" :key="index" #[item.slot!]="{ row }">
        <!-- 自定义渲染逻辑 -->
        <span v-if="item.prop === 'staffs'">
          {{ (row.staffs as Array<IStaff>)?.map((staff) => staff.staffname).join('; ') || '-' }}
        </span>
        <span v-else-if="item.prop === 'transfer_keywords'">
          {{ row[item.prop]?.split(',').join('; ') || '-' }}
        </span>
        <!-- 格式化接待方式 -->
        <span v-else-if="item.prop === 'status'">
          {{ formatStatus(row.status) || '-' }}
        </span>
        <!-- 格式化是否优先上一位接待人员 -->
        <span v-else-if="item.prop === 'receive_priority'">
          {{ formatReceivePriority(row.receive_priority) || '-' }}
        </span>
        <!-- 格式化接待规则 -->
        <span v-else-if="item.prop === 'receive_rule'">
          {{ formatReceiveRule(row.receive_rule) || '-' }}
        </span>
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
import { ref } from 'vue';
import router from '@/router';
import { useuuidDelete } from '@/hooks';
import { Api } from '@/api/common/enum';
import { deleteKFApi } from '@/api/kf';
import type { ITableConfig, TableInstance } from '@/types';
import type { IStaff } from '@/api/staff/model';

// 配置表格信息
const tableConfig: ITableConfig = {
  api: `${Api.kf}`, // 使用 kf 的接口路径
  headers: [
    { prop: 'kfname', label: '客服名称' },
    // { prop: 'kfid', label: '客服 ID' },
    { prop: 'kfplatform', label: '客服平台' },
    // { prop: 'botid', label: '机器人 ID' },
    { prop: 'status', label: '接待方式', slot: 'status' },
    { prop: 'receive_priority', label: '优先上一位接待人员', slot: 'receive_priority' },
    { prop: 'receive_rule', label: '接待规则', slot: 'receive_rule' },
    { prop: 'chat_timeout', label: '会话超时（秒）' },
    { prop: 'bot_timeout', label: '机器人超时（秒）' },
    { prop: 'transfer_keywords', label: '转接关键词', slot: 'transfer_keywords' },
    { prop: 'staffs', label: '接待人员列表', slot: 'staffs' },
  ],
  operations: {
    width: 150, // 操作列宽度
    buttons: [
      {
        text: '编辑',
        show: true,
        click: ({ row }) => {
          router.push(`/kf/edit/${row.uuid}`);
        },
      },
      {
        text: '删除',
        type: 'danger',
        show: true,
        click: ({ row }) => {
          handleDelete(row.uuid);
        },
      },
    ],
  },
};

// 筛选带有 slot 的字段
const slotColumns = tableConfig.headers.filter((header) => header.slot);

// 表格实例
const tableRef = ref<TableInstance>();

// 格式化接待方式
const formatStatus = (status: number) => {
  switch (status) {
    case 1:
      return '机器人可转人工';
    case 2:
      return '仅机器人';
    case 3:
      return '仅人工';
    default:
      return '-';
  }
};

// 格式化是否优先上一位接待人员
const formatReceivePriority = (receivePriority: number) => {
  return receivePriority === 1 ? '优先' : '不优先';
};

// 格式化接待规则
const formatReceiveRule = (receiveRule: number) => {
  switch (receiveRule) {
    case 1:
      return '轮流接待';
    case 2:
      return '空闲接待';
    default:
      return '-';
  }
};

// 添加客服
const handleAdd = () => {
  router.push('/kf/add');
};

// 删除客服
const handleDelete = (uuid: string) => {
  const { onDelete } = useuuidDelete();
  onDelete(deleteKFApi, uuid, tableRef); // 使用 kf 的删除接口
};
</script>