import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/login'
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue')
    },
    {
      path: '/home',
      name: 'home',
      component: () => import('@/views/HomeView.vue')
    },
    {
      path: '/grade',
      name: 'grade',
      component: () => import('@/views/GradeView.vue')
    },
    {
      path: '/schedule',
      name: 'schedule',
      component: () => import('@/views/ScheduleView.vue')
    },
    {
      path: '/exam',
      name: 'exam',
      component: () => import('@/views/ExamView.vue')
    },
    {
      path: '/selection',
      name: 'selection',
      component: () => import('@/views/SelectionView.vue')
    },
    {
      path: '/course-plan',
      name: 'course-plan',
      component: () => import('@/views/CoursePlanView.vue')
    }
  ]
})

router.beforeEach((to) => {
  const cookie = localStorage.getItem('zhjw_cookie')

  if (to.path !== '/login' && !cookie) {
    return '/login'
  }

  if (to.path === '/login' && cookie) {
    return '/home'
  }

  return true
})

export default router
