const routes = [
    {
        path: '/',
        component: () => import('@/layouts/index-public.vue'),
        children: [
            { path: '', name: 'Tab1', component: () => import('@/views/Tab1') },
            { path: 'Tab2', name: 'Tab2', component: () => import('@/views/Tab1') }
        ]
    }
]

if (process.env.MODE !== 'ssr') {
    routes.push({
        path: '*',
        component: () => import('@/views/errors/404.vue')
    })
}

export default routes
