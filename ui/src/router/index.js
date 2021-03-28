import { createRouter, createWebHistory } from '@ionic/vue-router'
// import routes from './routes'
const routes = [
    {
        path: '/',
        component: () => import('@/layouts/index-public.vue'),
        xcomponent: () => import('@/views/Home'),
        children: [
            { path: '', name: 'Home', component: () => import('@/views/Home') },
            { path: 'wind', name: '3dWind', component: () => import('@/views/3d_wind_home') }
        ]
    }
]

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes
})

export default router
