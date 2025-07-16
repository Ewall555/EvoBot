<template>
    <o-form-wrap title="llmapp" @confirm="onConfirm">
        <el-form ref="staffForm" class="w-100" :model="formData" :rules="formRules" label-position="top">
            <!-- 人员名称 -->
            <el-form-item label="名称：" prop="config_name">
                <el-input v-model="formData.config_name" placeholder="请输入名称" />
            </el-form-item>

            <!-- 类型 -->
            <el-form-item label="类型：" prop="llmapp_type">
                <el-input v-model="formData.llmapp_type" placeholder="请输入类型" />
            </el-form-item>

            <!-- base_url -->
            <el-form-item label="base_url：" prop="base_url">
                <el-input v-model="formData.base_url" placeholder="请输入base_url" />
            </el-form-item>

            <!-- api_key -->
            <el-form-item label="api_key：" prop="api_key">
                <el-input v-model="formData.api_key" placeholder="请输入api_key" />
            </el-form-item>
        </el-form>
    </o-form-wrap>
</template>

<script lang="ts" setup>
import { onBeforeMount, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import { createConfigApi, getConfigApi, updateConfigApi } from '@/api/llmapp/config';
import { getConfigListApi } from '@/api/llmapp/config';
import type { IPolicy } from '@/api/policy/model';
import type { IConfig } from '@/api/llmapp/model/configModel';

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
const formData = reactive<IConfig>({
    uuid:'',
    llmapp_type: '',
    config_name: '',
    api_key: '',
    base_url: '',
});

// 表单校验规则
const formRules = reactive<FormRules>({
    config_name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
    llmapp_type: [{ required: true, message: '请输入类型', trigger: 'blur' }],
    api_key: [{ required: true, message: '请输入api_key', trigger: 'blur' }],
    base_url: [{ required: true, message: '请输入base_url', trigger: 'blur' }],
});

// 选项
const llmappOptions = ref<IConfig[]>([]);

// 获取列表
const fetchLLMAppList = async () => {
    try {
        const { data } = await getConfigListApi();
        llmappOptions.value = data.map((llmapp) => ({
            uuid: llmapp.uuid,
            config_name: llmapp.config_name,
            llmapp_type: llmapp.llmapp_type,
            api_key: llmapp.api_key,
            base_url: llmapp.base_url,
        }));
    } catch (error) {
        ElMessage.error('获取列表失败');
    }
};

// 获取信息
const fetchLLMAppInfo = async () => {
    try {
        const { data } = await getConfigApi(uuid);

        // 将返回的完整数据填充到表单
        Object.assign(formData, data);

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
                if (isEditing.value) {
                    await updateConfigApi(requestData);  // 编辑时发送请求
                    ElMessage.success('编辑成功');
                } else {
                    await createConfigApi(requestData);  // 添加时发送请求
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
    if (isEditing.value) fetchLLMAppInfo();
});
</script>
