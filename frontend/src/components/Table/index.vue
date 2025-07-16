<template>
  <o-search :loading="loading" @search="submitSearch" />

  <el-card v-loading="loading" element-loading-text="加载中...">
    <div class="flex justify-between mb-4">
      <slot name="table-top" />
    </div>
    <el-table :data="tableData" border v-bind="$attrs">
      <!-- 表格复选框 -->
      <el-table-column v-if="tableConfig.selection" type="selection" width="55" reserve-selection />

      <!-- 表格索引 -->
      <el-table-column v-if="tableConfig.index" type="index" width="55" label="序号" />

      <el-table-column v-for="item in tableConfig.headers" :key="item.prop" :prop="item.prop" :label="item.label"
        show-overflow-tooltip :sortable="item.sortable" :min-width="item.width ? item.width : ''">
        <template #default="scope">
          <!-- 追加插槽渲染 -->
          <slot v-if="item.slot" :row="scope.row" :name="item.slot" />

          <!-- 常规渲染 -->
          <template v-else>
            <!-- 空数据 -->
            <span v-if="scope.row[item.prop] === ''"> - </span>
            <!-- 普通数据 -->
            <span v-else> {{ scope.row[item.prop] }} </span>
          </template>
        </template>
      </el-table-column>

      <!-- 操作按钮 -->
      <el-table-column v-if="tableConfig.operations?.buttons" :fixed="tableConfig.operations.fixed || 'right'"
        :width="tableConfig.operations.width">
        <template #header>
          <span>操作</span>
        </template>
        <template #default="scope">
          <template v-for="btn in tableConfig.operations.buttons"
            :key="typeof btn === 'function' ? btn(scope).text : btn.text">
            <template v-if="isFunction(btn)">
              <div v-if="btn(scope).show" class="px-2 inline-block">
                <el-link :underline="false" :type="btn(scope).type ? btn(scope).type : 'primary'"
                  @click.stop="btn(scope).click!(scope)">
                  {{ btn(scope).text }}
                </el-link>
              </div>
            </template>
            <template v-else>
              <el-link v-if="btn.show" class="px-2 inline-block" :underline="false"
                :type="btn.type ? btn.type : 'primary'" @click.stop="btn.click!(scope)">
                {{ btn.text }}
              </el-link>
            </template>
          </template>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" class="mt-5"
      :page-sizes="pagination.sizes" :layout="pagination.layout" :total="pagination.total"
      @size-change="handleSizeChange" @current-change="handleCurrentChange" />
  </el-card>
</template>
<script lang="ts" setup>
import { type PropType, computed, onMounted, watch } from 'vue'
import { isFunction } from 'native-lodash'
import { useRoute } from 'vue-router'
import { makeRequiredProp } from '@/utils'
import { usePagination, useTable } from '@/hooks'
import { getTableData } from '@/api/common'
import { useSearchStore } from '@/stores'
import type { ITableConfig, TObject } from '@/types'

defineOptions({
  name: 'OTable',
  inheritAttrs: true
})

// Props 接收
const props = defineProps({
  tableConfig: makeRequiredProp<PropType<ITableConfig>>(Object)
})

// 路由 & Store
const route = useRoute()
const searchStore = useSearchStore()
const { pagination: pager } = usePagination()

// 初始化搜索参数
let searchData: TObject = { ...props.tableConfig.params }

// 数据与状态
const {
  loading,
  tableData,
  pagination,
  loadData,
  handleSizeChange,
  handleCurrentChange
} = useTable(getTableData, {
  api: props.tableConfig.api,
  params: searchData,
  immediate: false
})

// 计算属性：搜索状态
const searchReady = computed(() => searchStore.searchReady)

// 数据加载方法（分页 & 搜索统一入口）
const loadTableData = (search: TObject = {}) => {
  loadData({
    ...searchData,
    page: pagination.page,
    page_size: pagination.page_size,
    ...search
  })
}

// 提交搜索
const submitSearch = (search: TObject) => {
  pagination.page = 1 // 搜索时重置为第一页
  loadTableData(search)
}
// 页面加载标志
let isInitialLoad = false

// onMounted 初始化
onMounted(async () => {
  isInitialLoad = true // 标记初始加载
  const query = route.query

  // 设置分页参数
  pagination.page = Number(query.page) || pager.page
  pagination.page_size = Number(query.page_size) || pager.page_size

  Object.assign(searchData, {
    page: pagination.page,
    page_size: pagination.page_size
  })

  // 加载数据
  const ready = await searchStore.getSearchList()
  if (ready) {
    loadTableData()
  }
  isInitialLoad = false // 解除初始加载标志
})

// 监听分页参数变化
watch([() => pagination.page, () => pagination.page_size], () => {
  if (!isInitialLoad) {
    loadTableData() // 仅在非初始加载时触发
  }
})

// 公开方法
defineExpose({
  loadTableData
})
</script>
