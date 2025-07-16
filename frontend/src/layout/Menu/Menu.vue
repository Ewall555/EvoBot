<template>
  <el-scrollbar class="menu-container">
    <el-menu
      :collapse="layoutStore.isCollapse"
      :default-active="activeMenu"
      unique-opened
    >
      <o-menu-item :menu-list="menuList" />
    </el-menu>
  </el-scrollbar>
</template>

<script lang="ts" setup>
import { onBeforeMount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useLayoutStore } from '@/stores'
import OMenuItem from './MenuItem.vue'
import type { IMenuItem } from '../types'

defineOptions({
  name: 'OMenu'
})

const route = useRoute()
const layoutStore = useLayoutStore()

const activeMenu = ref('')
const menuList = ref<IMenuItem[]>([])

const defaultMenu = [
  {
    id: 'home',
    name: '首页',
    icon: 'HomeFilled',
    path: '/home',
  },
  {
    id: 'kf',
    name: '客服应用管理',
    icon: 'Management',
    path: '/kf',
  },
  {
    id: 'staff',
    name: '坐席人员管理',
    icon: 'Management',
    path: '/staff',
  },
  {
    id: 'policy',
    name: '工作策略管理',
    icon: 'Management',
    path: '/policy',
  },
  {
    id: 'llmapp',
    name: '知识库管理',
    icon: 'Management',
    path: '/llmapp',
  },
  {
    id: 'wecom',
    name: '企业微信管理',
    icon: 'Management',
    children: [
      {
        id: 'account',
        name: '客服账号管理',
        path: '/wecom/account',
      },
      {
        id: 'config',
        name: '自建应用配置',
        path: '/wecom/config',
      },
    ],
  },
  // {
  //   id: 'system',
  //   name: '系统管理',
  //   icon: 'Platform',
  //   children: [
  //     {
  //       id: 'user',
  //       name: '用户管理',
  //       path: '/system/user',
  //     },
  //     {
  //       id: 'role',
  //       name: '角色管理',
  //       path: '/system/role',
  //     },
  //   ],
  // },
];

onBeforeMount(() => {
  menuList.value = defaultMenu;
})

watch(
  () => route,
  (val) => {
    const { meta } = val
    activeMenu.value = meta.activePath as string
  },
  { immediate: true, deep: true }
)
</script>

<style lang="scss" scoped>
.menu-container {
  position: fixed;
  left: 0;
  top: 60px;
  bottom: 0;
  z-index: 2000;
  height: calc(100% - 60px);
  overflow-y: auto;
  box-sizing: border-box;
  background-color: #fff;
  box-shadow: 0 2px 10px 0 rgba(0, 0, 0, 0.1);
}

.el-menu {
  border-right: 0;

  &:not(.el-menu--collapse) {
    width: 220px;
  }
}
</style>
@/stores/modules/layout
