import { ref } from 'vue'
import { defineStore } from 'pinia'
import searchConfig from '@/assets/enum/search'
import type { ISearchItem } from '@/types'
import { useRoute } from 'vue-router'

export const useSearchStore = defineStore('search', () => {
  const searchReady = ref(false) // 列表依赖搜索项，等搜索项获取完成后再获取列表
  const searchList = ref<ISearchItem[]>([])

  const getSearchList = async (): Promise<boolean> => {
    searchReady.value = false
    const route = useRoute()
    const searchId = searchConfig[route.path]
  
    if (!searchId) {
      console.warn('当前路径无搜索配置')
      searchReady.value = true
      return true
    }
  
    searchList.value = searchId
    searchReady.value = true
    return true
  }

  return {
    searchList,
    searchReady,
    getSearchList
  }
})
