<template>
  <o-form-wrap title="工作策略管理" @confirm="onConfirm">
    <el-form ref="ruleForm" class="w-100" :model="formData" :rules="formRules" label-position="top">
      <el-form-item label="策略名称：" prop="policyname">
        <el-input v-model="formData.policyname" placeholder="请输入策略名称" />
      </el-form-item>

      <el-form-item label="最大接待数量：" prop="max_count">
        <el-input-number v-model="formData.max_count" placeholder="请输入最大接待数量" />
      </el-form-item>

      <el-form-item label="重复策略：" prop="repeat">
        <el-select v-model="formData.repeat" placeholder="请选择重复策略">
          <el-option v-for="option in repeatOptions" :key="option.value" :value="option.value" :label="option.label" />
        </el-select>
      </el-form-item>

      <!-- 只有当 repeat === 1 时，显示 week 字段 -->
      <el-form-item v-if="formData.repeat === 1" label="工作日：" prop="week">
        <!-- 使用CheckboxGroup显示工作日选择框 -->
        <el-checkbox-group v-model="selectedDays" @change="updateWeek">
          <el-checkbox label="0" name="day">周日</el-checkbox>
          <el-checkbox label="1" name="day">周一</el-checkbox>
          <el-checkbox label="2" name="day">周二</el-checkbox>
          <el-checkbox label="3" name="day">周三</el-checkbox>
          <el-checkbox label="4" name="day">周四</el-checkbox>
          <el-checkbox label="5" name="day">周五</el-checkbox>
          <el-checkbox label="6" name="day">周六</el-checkbox>
        </el-checkbox-group>
      </el-form-item>

      <el-form-item label="工作时间：" prop="work_times">
        <el-button type="primary" @click="addWorkTime">添加工作时间</el-button>
          <div v-for="(time, index) in formData.work_times" :key="index" class="work-time">
            <el-time-picker
              v-model="tempTimeRanges[index]"
              is-range
              range-separator="至"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
              format="HH:mm:ss"
              value-format="HH:mm:ss"
              @change="(val) => handleTimeChange(index, val)"
            ></el-time-picker>
            <el-button @click="removeWorkTime(index)" type="danger">删除</el-button>
          </div>
      </el-form-item>
    </el-form>
  </o-form-wrap>
</template>

<script lang="ts" setup>
import { onBeforeMount, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'

import { createPolicyApi, getPolicyApi, updatePolicyApi } from '@/api/policy'
import type { IPolicy } from '@/api/policy/model'

import type { TLoading } from '@/types'
import type { FormInstance, FormRules } from 'element-plus'

const route = useRoute()
const router = useRouter()

const isEditing = ref(false)
const uuid = route.params.uuid as string

isEditing.value = !!uuid

const ruleForm = ref<FormInstance>()
const formData = reactive<IPolicy>({
  uuid: '',
  policyname: '',
  max_count: 100,
  repeat: 1,
  week: '',
  work_times: [{ start_time: '08:00:00', end_time: '18:00:00' }]
})

const tempTimeRanges = ref<[string, string][]>([
  ['08:00:00', '18:00:00'] 
]);

const handleTimeChange = (index: number, val: [string, string] | null) => {
  if (Array.isArray(formData.work_times)) {
    if (val) {
      formData.work_times[index].start_time = val[0];
      formData.work_times[index].end_time = val[1];
      tempTimeRanges.value[index] = val;
    }
  }
};

const removeWorkTime = (index: number) => {
  if (Array.isArray(formData.work_times)) {
    formData.work_times.splice(index, 1);
    tempTimeRanges.value.splice(index, 1);
  }
};

interface RepeatOption {
  value: number
  label: string
}

const repeatOptions: RepeatOption[] = [
  { value: 1, label: '自定义工作日' },
  { value: 2, label: '每天' },
  { value: 3, label: '周一至周五' },
  { value: 4, label: '法定工作日（跳过法定节假日）' },
  { value: 5, label: '法定节假日（跳过法定工作日）' }
]

const formRules = reactive<FormRules>({
  policyname: [{ required: true, message: '请输入策略名称', trigger: 'blur' }],
  max_count: [{ required: true, message: '请输入最大接待数量', trigger: 'blur' }],
  repeat: [{ required: true, message: '请选择重复策略', trigger: 'change' }],
  week: [{ required: true, message: '当重复策略为自定义时，工作日必填', trigger: 'blur' }],
  work_times: [{ required: true, type: 'array', min: 1, message: '请添加至少一个工作时间', trigger: 'blur' }]
})

const policyInfo = async () => {
  const { data } = await getPolicyApi(uuid)
  Object.assign(formData, data)

  // 根据 week 标记更新 selectedDays
  if (formData.week) {
    selectedDays.value = formData.week
      .split('')
      .map((day, index) => (day === '1' ? String(index) : null))
      .filter(Boolean) as string[]  // 将勾选的工作日的索引添加到 selectedDays 中
  }
  if (Array.isArray(formData.work_times)) {
    tempTimeRanges.value = formData.work_times.map(t => [t.start_time, t.end_time] as [string, string]);
  }
}

const onConfirm = (loading: TLoading) => {
  ruleForm.value?.validate(async (valid) => {
    if (valid) {
      loading(true)
      if (isEditing.value) {
        await updatePolicyApi(formData)
        ElMessage.success('编辑成功')
      } else {
        await createPolicyApi(formData)
        ElMessage.success('添加成功')
      }

      loading(false)
      router.back()
    } else {
      return
    }
  })
}

const addWorkTime = () => {
  if (Array.isArray(formData.work_times)) {
    formData.work_times.push({ start_time: '', end_time: '' })
  }
}

// const removeWorkTime = async (index: number) => {
//   try {
//     await ElMessageBox.confirm('确定要删除该工作时间吗？', '删除确认', {
//       confirmButtonText: '确定',
//       cancelButtonText: '取消',
//       type: 'warning'
//     })
//     formData.work_times?.splice(index, 1)
//     ElMessage.success('工作时间已删除')
//   } catch { }
// }

// 当用户勾选工作日时，更新week字段
const selectedDays = ref<string[]>([])  // 用于保存勾选的工作日

const updateWeek = () => {
  // 将勾选的工作日按顺序拼接成一个字符串
  formData.week = Array(7).fill('0').map((_, index) =>
    selectedDays.value.includes(String(index)) ? '1' : '0'
  ).join('')
}

onBeforeMount(() => {
  if (isEditing.value) policyInfo()
})
</script>
