<template>
  <o-form-wrap title="配置管理" @confirm="onConfirm">
    <el-form
      ref="ruleForm"
      class="w-100"
      :model="formData"
      :rules="formRules"
      label-position="top"
    >
      <el-form-item label="类型：" prop="type">
        <el-input v-model="formData.type" placeholder="请输入配置类型" />
      </el-form-item>

      <el-form-item label="Agent ID：" prop="agent_id">
        <el-input v-model="formData.agent_id" placeholder="请输入Agent ID" />
      </el-form-item>

      <el-form-item label="Corp ID：" prop="corp_id">
        <el-input v-model="formData.corp_id" placeholder="请输入Corp ID" />
      </el-form-item>

      <el-form-item label="Secret：" prop="secret">
        <el-input v-model="formData.secret" placeholder="请输入Secret" />
      </el-form-item>

      <el-form-item label="Encoding AES Key：" prop="encoding_aes_key">
        <el-input v-model="formData.encoding_aes_key" placeholder="请输入Encoding AES Key" />
      </el-form-item>

      <el-form-item label="Token：" prop="token">
        <el-input v-model="formData.token" placeholder="请输入Token" />
      </el-form-item>
    </el-form>
  </o-form-wrap>
</template>

<script lang="ts" setup>
import { onBeforeMount, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'

import { getConfigApi, updateConfigApi } from '@/api/wecom/config'
import type { IConfig } from '@/api/wecom/model/configModel'

import type { TLoading } from '@/types'
import type { FormInstance, FormRules } from 'element-plus'

const route = useRoute()
const router = useRouter()

const uuid = route.params.uuid as string

const ruleForm = ref<FormInstance>()
const formData = reactive<IConfig>({
  uuid: '',
  type: '',
  agent_id: '',
  corp_id: '',
  secret: '',
  encoding_aes_key: '',
  token: '',
})

const formRules = reactive<FormRules>({
  type: [{ required: true, message: '请输入配置类型', trigger: 'blur' }],
  agent_id: [{ required: true, message: '请输入Agent ID', trigger: 'blur' }],
  corp_id: [{ required: true, message: '请输入Corp ID', trigger: 'blur' }],
  secret: [{ required: true, message: '请输入Secret', trigger: 'blur' }],
  encoding_aes_key: [{ required: true, message: '请输入Encoding AES Key', trigger: 'blur' }],
  token: [{ required: true, message: '请输入Token', trigger: 'blur' }],
})

const configInfo = async () => {
  try {
    const { data } = await getConfigApi(uuid)
    Object.assign(formData, data)
  } catch (error) {
    ElMessage.error('加载配置数据失败')
  }
}

const onConfirm = (loading: TLoading) => {
  ruleForm.value?.validate(async (valid) => {
    if (valid) {
      loading(true)
      // 更新配置
      await updateConfigApi(formData)
      ElMessage.success('编辑成功')

      loading(false)
      router.back()
    } else {
      return
    }
  })
}

onBeforeMount(() => {
  configInfo() // 只在编辑时加载数据
})
</script>
