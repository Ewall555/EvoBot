export default [
    {
      path: '/policy',
      children: [
        {
          path: '',
          component: () => import('@/views/policy/index.vue'),
          meta: {
            title: '工作策略管理',
            activePath: '/policy'
          }
        },
        {
          path: 'add',
          component: () => import('@/views/policy/add.vue'),
          meta: {
            title: '工作策略管理-添加',
            activePath: '/policy'
          }
        },
        {
          path: 'edit/:uuid',
          component: () => import('@/views/policy/add.vue'),
          meta: {
            title: '工作策略管理-编辑',
            activePath: '/policy'
          }
        }
      ]
    }
  ]
  