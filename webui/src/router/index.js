import Vue from 'vue'
import Router from 'vue-router'
import i18n from '@/lang'

Vue.use(Router)

const _import = file => () => import('../view/' + file + '.vue')

const routes = [
    {
        path: '/',
        component: _import('layer'),
        name: 'dashboard',
        title: i18n.t('dashboard'),
        redirect: { name: 'dashboard' },
        children: [
            {
                path: 'dashboard',
                name: 'dashboard',
                meta: {
                    title: i18n.t('dashboard'),
                    icon: 'el-icon-monitor',
                    single: true,
                },
                component: _import('dashboard'),
            }
        ],
    },
    {
        path: '/bulk_upload',
        name: 'bulk_upload',
        title: i18n.t('bulk_upload'),
        component: _import('layer'),
        meta: {
            title: i18n.t('bulk_upload'),
            icon: 'el-icon-upload',
        },
        children: [
            {
                path: 'bulk_upload_pod',
                name: 'bulk_upload_pod',
                meta: {
                    title: i18n.t('bulk_upload_pod'),
                    icon: 'el-icon-upload2',
                },
                component: _import('bulk_upload_pod'),
            },
            {
                path: 'bulk_upload_container',
                name: 'bulk_upload_container',
                meta: {
                    title: i18n.t('bulk_upload_container'),
                    icon: 'el-icon-upload2',
                },
                component: _import('bulk_upload_container'),
            },
            {
                path: 'bulk_upload_pvc',
                name: 'bulk_upload_pvc',
                meta: {
                    title: i18n.t('bulk_upload_pvc'),
                    icon: 'el-icon-upload2',
                },
                component: _import('bulk_upload_pvc'),
            },
        ],
    }
]

const router = new Router({
    routes: routes,
    base: __dirname,
    scrollBehavior: () => ({ y: 0 }),
    mode: 'history',
})
export { routes }
export default router