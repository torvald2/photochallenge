import Vue from 'vue'
import VueRouter from 'vue-router'
import HomeView from '../views/ConnectWallet.vue'
import CreateView from '../views/CreateView.vue'
import CompleteView from '../views/CompleteView.vue'
Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
    props: true

  },
  {
    path: '/create',
    name: 'create',
    component: CreateView,
    props: true

  },
  {
    path: '/complete',
    name: 'complete',
    component: CompleteView,
    props: true
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
