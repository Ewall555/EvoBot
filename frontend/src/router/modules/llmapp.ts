export default [
    {
      path: '/llmapp',
      children: [
        {
          path: '',
          component: () => import('@/views/llmapp/index.vue'),
          meta: {
            title: 'llm 应用管理',
            activePath: '/llmapp'
          }
        },
        {
          path: 'add',
          component: () => import('@/views/llmapp/add.vue'),
          meta: {
            title: 'llm 应用管理-添加',
            activePath: '/llmapp'
          }
        },
        {
          path: 'edit/:uuid',
          component: () => import('@/views/llmapp/add.vue'),
          meta: {
            title: 'llm 应用管理-编辑',
            activePath: '/llmapp'
          }
        }
    ]
}
]