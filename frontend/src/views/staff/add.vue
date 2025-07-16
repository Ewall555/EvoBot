<template>
    <o-form-wrap title="人员管理" @confirm="onConfirm">
        <el-form ref="staffForm" class="w-100" :model="formData" :rules="formRules" label-position="top">
            <!-- 人员名称 -->
            <el-form-item label="人员名称：" prop="staffname">
                <el-input v-model="formData.staffname" placeholder="请输入人员名称" />
            </el-form-item>

            <!-- 接待人员 ID -->
            <el-form-item label="接待人员 ID：" prop="staffid">
                <el-input v-model="formData.staffid" placeholder="请输入接待人员 ID" />
            </el-form-item>

            <!-- 电话号码 -->
            <el-form-item label="电话号码：" prop="number">
                <el-input v-model="formData.number" placeholder="请输入电话号码" />
            </el-form-item>

            <!-- 邮箱 -->
            <el-form-item label="邮箱：" prop="email">
                <el-input v-model="formData.email" placeholder="请输入邮箱" />
            </el-form-item>

            <!-- 角色 -->
            <el-form-item label="角色：" prop="role">
                <el-input v-model="formData.role" placeholder="请输入角色" />
            </el-form-item>

            <!-- 策略列表 -->
            <el-form-item label="策略列表：" prop="policy_list">
                <el-select v-model="formData.policy_list" multiple placeholder="请选择策略" :options="policyOptions">
                    <el-option v-for="policy in policyOptions" :key="policy.uuid" :label="policy.policyname"
                        :value="policy.uuid" />
                </el-select>
            </el-form-item>
        </el-form>
    </o-form-wrap>
</template>

<script lang="ts" setup>
import { onBeforeMount, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import { createStaffApi, getStaffApi, updateStaffApi } from '@/api/staff';
import { getPolicyListApi } from '@/api/policy';
import type { IPolicy } from '@/api/policy/model';
import type { IStaff } from '@/api/staff/model';

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
const staffForm = ref<FormInstance>();
const formData = reactive<IStaff>({
    uuid:'',
    staffname: '',
    staffid: '',
    number: '',
    email: '',
    role:'',
    policy_list: [], // 存储策略 UUID 列表
});

// 表单校验规则
const formRules = reactive<FormRules>({
    staffname: [{ required: true, message: '请输入人员名称', trigger: 'blur' }],
    staffid: [{ required: true, message: '请输入接待人员 ID', trigger: 'blur' }],
    role: [{ required: true, message: '请输入角色', trigger: 'blur' }],
    policy_list: [
        {
            required: true,
            type: 'array',
            min: 1,
            message: '请至少选择一个策略',
            trigger: 'change',
        },
    ],
});

// 策略选项
const policyOptions = ref<IPolicy[]>([]);

// 获取策略列表
const fetchPolicyList = async () => {
    try {
        const { data } = await getPolicyListApi();
        policyOptions.value = data.map((policy) => ({
            uuid: policy.uuid,
            policyname: policy.policyname,
        }));
    } catch (error) {
        ElMessage.error('获取策略列表失败');
    }
};

// 获取人员信息
const fetchStaffInfo = async () => {
    try {
        const { data } = await getStaffApi(uuid);

        // 将返回的完整数据填充到表单
        Object.assign(formData, data);

        // 如果 `policies` 字段存在且是数组，直接提取 `UUID`
        formData.policy_list = data.policies?.map(policy => policy.uuid) || [];
    } catch (error) {
        ElMessage.error('获取人员信息失败');
    }
};


const onConfirm = (loading: TLoading) => {
    staffForm.value?.validate(async (valid) => {
        if (valid) {
            loading(true);
            try {
                const requestData = {
                    ...formData,
                };
                delete requestData.policies;
                if (isEditing.value) {
                    await updateStaffApi(requestData);  // 编辑时发送请求
                    ElMessage.success('编辑成功');
                } else {
                    await createStaffApi(requestData);  // 添加时发送请求
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
    fetchPolicyList();
    if (isEditing.value) fetchStaffInfo();
});
</script>
