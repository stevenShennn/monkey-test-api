import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Execute from '../views/Execute.vue'
import SiteDetail from '../views/SiteDetail.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/execute',
    name: 'Execute',
    component: Execute
  },
  {
    path: '/site/:domain',
    name: 'SiteDetail',
    component: SiteDetail,
    props: true
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 