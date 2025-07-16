export default [
    {
      path: '/kf',
      children: [
        {
          path: '',
          component: () => import('@/views/kf/index.vue'),
          meta: {
            title: '客服应用管理',
            activePath: '/kf'
          }
        },
        {
          path: 'add',
          component: () => import('@/views/kf/add.vue'),
          meta: {
            title: '客服应用管理-添加',
            activePath: '/kf'
          }
        },
        {
          path: 'edit/:uuid',
          component: () => import('@/views/kf/add.vue'),
          meta: {
            title: '客服应用管理-编辑',
            activePath: '/kf'
          }
        }
      ]
    }
  ]
  