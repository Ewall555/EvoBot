import { ElMessage } from 'element-plus'
import { confirmBox, setUrlParams } from '@/utils'
import { usePagination } from '../usePagination'
import type { Ref } from 'vue'
import type { TableInstance } from '@/types'
import type { IReceptionist } from '@/api/wecom/model/receptionistModel'


// 用于列表删除，解决删除最后一页数据最后一条数据问题
export function useDelete() {
  const onDelete = (
    deleteFun: (id: number) => Promise<any>,
    id: number,
    tableRef: Ref<TableInstance | undefined>
  ) => {
    confirmBox(`请确认是否要删除该条数据？`).then(async () => {
      await deleteFun(id)
      ElMessage.success('删除成功')

      const { pagination } = usePagination()

      const pageCount = Math.ceil((pagination.total - 1) / pagination.page_size)
      const currentPage = Math.min(pageCount, pagination.page) || 1
      if (currentPage !== pagination.page) {
        setUrlParams({ page: currentPage })
        pagination.page = currentPage
      }

      tableRef.value?.loadTableData()
    })
  }

  return {
    onDelete
  }
}

export function useuuidDelete() {
  const onDelete = (
    deleteFun: (uuid: string) => Promise<any>,
    uuid: string,
    tableRef: Ref<TableInstance | undefined>
  ) => {
    confirmBox(`请确认是否要删除该条数据？`).then(async () => {
      await deleteFun(uuid)
      ElMessage.success('删除成功')

      const { pagination } = usePagination()

      const pageCount = Math.ceil((pagination.total - 1) / pagination.page_size)
      const currentPage = Math.min(pageCount, pagination.page) || 1
      if (currentPage !== pagination.page) {
        setUrlParams({ page: currentPage })
        pagination.page = currentPage
      }

      tableRef.value?.loadTableData()
    })
  }

  return {
    onDelete
  }
}

export function usekfidDelete() {
  const onDelete = (
    deleteFun: (kfid: string,params: IReceptionist) => Promise<any>,
    kfid: string,
    params: IReceptionist,
    tableRef: Ref<TableInstance | undefined>
  ) => {
    confirmBox(`请确认是否要删除该条数据？`).then(async () => {
      await deleteFun(kfid,params)
      ElMessage.success('删除成功')

      const { pagination } = usePagination()

      const pageCount = Math.ceil((pagination.total - 1) / pagination.page_size)
      const currentPage = Math.min(pageCount, pagination.page) || 1
      if (currentPage !== pagination.page) {
        setUrlParams({ page: currentPage })
        pagination.page = currentPage
      }

      tableRef.value?.loadTableData()
    })
  }

  return {
    onDelete
  }
}