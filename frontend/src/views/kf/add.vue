<template>
    <o-form-wrap title="客服管理" @confirm="onConfirm">
        <el-form ref="kfForm" class="w-100" :model="formData" :rules="formRules" label-position="top">
            <el-form-item label="客服名称：" prop="kfname">
                <el-input v-model="formData.kfname" placeholder="请输入客服名称" />
            </el-form-item>

            <el-form-item label="客服 ID：" prop="kfid">
                <el-input v-model="formData.kfid" placeholder="请输入客服 ID" />
            </el-form-item>

            <el-form-item label="客服平台：" prop="kfplatform">
                <el-input v-model="formData.kfplatform" placeholder="请输入客服平台" />
            </el-form-item>

            <el-form-item label="机器人 ID：" prop="botid">
                <el-input v-model="formData.botid" placeholder="请输入机器人 ID" />
            </el-form-item>

            <el-form-item label="接待方式：" prop="status">
                <el-radio-group v-model="formData.status">
                    <el-radio :label="1">机器人可转人工</el-radio>
                    <el-radio :label="2">仅机器人</el-radio>
                    <el-radio :label="3">仅人工</el-radio>
                </el-radio-group>
            </el-form-item>


            <el-form-item label="是否优先上一位接待人员：" prop="receive_priority">
                <el-radio-group v-model="formData.receive_priority">
                    <el-radio :label="0">不优先</el-radio>
                    <el-radio :label="1">优先</el-radio>
                </el-radio-group>
            </el-form-item>

            <el-form-item label="接待规则：" prop="receive_rule">
                <el-radio-group v-model="formData.receive_rule">
                    <el-radio :label="1">轮流接待</el-radio>
                    <el-radio :label="2">空闲接待</el-radio>
                </el-radio-group>
            </el-form-item>


            <el-form-item label="会话超时（秒）：" prop="chat_timeout">
                <el-input-number v-model="formData.chat_timeout" placeholder="请输入会话超时时间" />
            </el-form-item>

            <el-form-item label="机器人超时（秒）：" prop="bot_timeout">
                <el-input-number v-model="formData.bot_timeout" placeholder="请输入机器人超时时间" />
            </el-form-item>

            <el-form-item label="机器人超时消息：" prop="bot_timeout_msg">
                <el-input v-model="formData.bot_timeout_msg" placeholder="请输入机器人超时消息" />
            </el-form-item>

            <el-form-item label="机器人欢迎消息（菜单）：" prop="bot_welcome_msg">
                <el-input type="textarea" :rows="9" v-model="formData.bot_welcome_msg" placeholder="请输入机器人欢迎消息" />
            </el-form-item>

            <el-form-item label="接待人员欢迎消息：" prop="staff_welcome_msg">
                <el-input v-model="formData.staff_welcome_msg" placeholder="请输入接待人员欢迎消息" />
            </el-form-item>

            <el-form-item label="无人值守消息：" prop="unmanned_msg">
                <el-input v-model="formData.unmanned_msg" placeholder="请输入无人值守消息" />
            </el-form-item>

            <el-form-item label="会话结束消息（菜单）：" prop="chatend_msg">
                <el-input type="textarea" :rows="3" v-model="formData.chatend_msg" placeholder="请输入会话结束消息" />
            </el-form-item>

            <el-form-item label="转人工关键字：" prop="transfer_keywords">
                <el-input v-model="formData.transfer_keywords" placeholder="请输入转接关键字，用逗号分隔" />
            </el-form-item>

            <el-form-item label="接待人员列表：" prop="staff_list">
                <el-select v-model="formData.staff_list" multiple placeholder="请选择接待人员" :options="staffOptions">
                    <el-option v-for="staff in staffOptions" :key="staff.uuid" :label="staff.staffname"
                        :value="staff.uuid" />
                </el-select>
            </el-form-item>
        </el-form>
    </o-form-wrap>
</template>

<script lang="ts" setup>
import { onBeforeMount, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import { createKFApi, getKFApi, updateKFApi } from '@/api/kf';
import { getStaffListApi } from '@/api/staff';
import type { IStaff } from '@/api/staff/model';
import type { IKF } from '@/api/kf/model';

import type { TLoading } from '@/types';
import type { FormInstance, FormRules } from 'element-plus';

// 路由实例
const route = useRoute();
const router = useRouter();

// 是否为编辑模式
const isEditing = ref(false);
const uuid = route.params.uuid as string;
isEditing.value = !!uuid;

// 表单引用与数据
const kfForm = ref<FormInstance>();
const formData = reactive<IKF>({
    kfname: '',
    kfid: '',
    kfplatform: '',
    botid: '',
    status: 1,
    receive_priority: 1,
    receive_rule: 1,
    chat_timeout: 30,
    bot_timeout: 30,
    bot_timeout_msg: '',
    bot_welcome_msg: '',
    staff_welcome_msg: '',
    unmanned_msg: '',
    chatend_msg: '',
    transfer_keywords: '',
    staff_list: [], // 接待人员 UUID 列表
});

// 表单校验规则
const formRules = reactive<FormRules>({
    kfname: [{ required: true, message: '请输入客服名称', trigger: 'blur' }],
    kfid: [{ required: true, message: '请输入客服 ID', trigger: 'blur' }],
    kfplatform: [{ required: true, message: '请输入客服平台', trigger: 'blur' }],
    status: [{ required: true, message: '请输入状态', trigger: 'blur' }],
    receive_priority: [{ required: true, message: '请输入接待优先级', trigger: 'blur' }],
    receive_rule: [{ required: true, message: '请输入接待规则', trigger: 'blur' }],
    staff_list: [
        {
            required: true,
            type: 'array',
            min: 1,
            message: '请至少选择一个接待人员',
            trigger: 'change',
        },
    ],
});

// 接待人员选项
const staffOptions = ref<IStaff[]>([]);

// 获取接待人员列表
const fetchStaffList = async () => {
    try {
        const { data } = await getStaffListApi();
        staffOptions.value = data.map((staff) => ({
            uuid: staff.uuid,
            staffname: staff.staffname,
        }));
    } catch (error) {
        ElMessage.error('获取接待人员列表失败');
    }
};

// 获取客服信息
const fetchKFInfo = async () => {
    try {
        const { data } = await getKFApi(uuid);
        Object.assign(formData, data);
        formData.staff_list = data.staffs?.map(staff => staff.uuid) || [];
    } catch (error) {
        ElMessage.error('获取客服信息失败');
    }
};

const onConfirm = (loading: TLoading) => {
    kfForm.value?.validate(async (valid) => {
        if (valid) {
            loading(true);
            try {
                const requestData = {
                    ...formData,
                };
                delete requestData.staffs;
                if (isEditing.value) {
                    await updateKFApi(requestData); // 编辑时发送请求
                    ElMessage.success('编辑成功');
                } else {
                    await createKFApi(requestData); // 添加时发送请求
                    ElMessage.success('添加成功');
                }
                router.back();
            } catch (error) {
                ElMessage.error('操作失败');
            } finally {
                loading(false);
            }
        }
    });
};

// 生命周期钩子
onBeforeMount(() => {
    fetchStaffList();
    if (isEditing.value) fetchKFInfo();
});
</script>
