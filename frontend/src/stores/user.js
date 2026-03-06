import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const cookie = ref(localStorage.getItem('zhjw_cookie') || '')
  const isLoggedIn = ref(!!cookie.value)

  function setCookie(newCookie) {
    cookie.value = newCookie
    isLoggedIn.value = true
    localStorage.setItem('zhjw_cookie', newCookie)
  }

  function clearCookie() {
    cookie.value = ''
    isLoggedIn.value = false
    localStorage.removeItem('zhjw_cookie')
  }

  return {
    cookie,
    isLoggedIn,
    setCookie,
    clearCookie
  }
})
