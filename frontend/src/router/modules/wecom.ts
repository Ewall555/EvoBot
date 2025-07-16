export default [
    {
      path: '/wecom/config',
      children: [
        {
          path: '',
          component: () => import('@/views/wecom/config/index.vue'),
          meta: {
            title: '自建应用配置',
            activePath: '/wecom/config'
          }
        },
        {
          path: 'edit/:uuid',
          component: () => import('@/views/wecom/config/edit.vue'),
          meta: {
            title: '自建应用配置-编辑',
            activePath: '/wecom/config'
          }
        }
      ]
    },
    {
      path: '/wecom/account',
      children: [
        {
          path: '',
          component: () => import('@/views/wecom/account/index.vue'),
          meta: {
            title: '客服账号管理',
            activePath: '/wecom/account'
          }
        },
        {
          path: 'details/:kfid',
          component: () => import('@/views/wecom/account/details.vue'),
          meta: {
            title: '客服账号管理-接待人员',
            activePath: '/wecom/account'
          }
        },
        {
          path: 'add',
          component: () => import('@/views/wecom/account/add.vue'),
          meta: {
            title: '客服账号管理-添加',
            activePath: '/wecom/account'
          }
        },
        {
          path: 'edit/:kfid',
          component: () => import('@/views/wecom/account/add.vue'),
          meta: {
            title: '客服账号管理-编辑',
            activePath: '/wecom/account'
          }
        }
      ]
    },
  ]
  