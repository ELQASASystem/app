import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'Index',
        component: () => import('../views/Index/Index')
    },
    {
        path: '/answer/tea/list',
        name: 'List',
        component: () => import('../views/Answer/Teacher/List')
    },
    {
        path: '/answer/tea/a/:id',
        name: 'Detail',
        component: () => import('../views/Answer/Teacher/Detail.vue')
    },
    {
        path: '/answer/tea/new',
        name: 'New',
        component: () => import('../views/Answer/Teacher/New')
    },
]

const router = new VueRouter({
    routes
})

export default router
