export default [
    {
      path: '/staff',
      children: [
        {
          path: '',
          component: () => import('@/views/staff/index.vue'),
          meta: {
            title: '坐席人员管理',
            activePath: '/staff'
          }
        },
        {
          path: 'add',
          component: () => import('@/views/staff/add.vue'),
          meta: {
            title: '坐席人员管理-添加',
            activePath: '/staff'
          }
        },
        {
          path: 'edit/:uuid',
          component: () => import('@/views/staff/add.vue'),
          meta: {
            title: '坐席人员管理-编辑',
            activePath: '/staff'
          }
        }
      ]
    }
  ]
  