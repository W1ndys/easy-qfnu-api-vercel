import api from './request'
import axios from 'axios'

const proxyApi = axios.create({
  baseURL: '/api/proxy',
  timeout: 30000,
  withCredentials: true,
  headers: {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
  }
})

export function getCaptchaUrl() {
  return '/api/proxy/jsxsd/verifycode.servlet'
}

export function initSession() {
  return proxyApi.get('/jsxsd/')
}

export function login(data) {
  const username = data.username
  const password = data.password
  const captcha = data.captcha
  
  const encoded = btoa(username) + '%%%' + btoa(password)
  
  const formData = new URLSearchParams()
  formData.append('encoded', encoded)
  formData.append('RANDOMCODE', captcha)
  
  return proxyApi.post('/jsxsd/xk/LoginToXkLdap', formData, {
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    }
  }).then(response => {
    const body = response.data
    if (typeof body === 'string') {
      if (body.includes('验证码错误')) {
        throw new Error('验证码错误')
      }
      if (body.includes('密码错误') || body.includes('用户名或密码错误')) {
        throw new Error('用户名或密码错误')
      }
    }
    return proxyApi.get('/jsxsd/framework/xsMain.jsp').then(verifyRes => {
      if (typeof verifyRes.data === 'string' && verifyRes.data.includes('用户登录')) {
        throw new Error('登录验证失败')
      }
      return { success: true }
    })
  })
}

export function getGrade(cookie) {
  return api.get('/v1/zhjw/grade', {
    headers: { Authorization: cookie }
  })
}

export function getCoursePlan(cookie) {
  return api.get('/v1/zhjw/course-plan', {
    headers: { Authorization: cookie }
  })
}

export function getExamSchedules(cookie) {
  return api.get('/v1/zhjw/exam', {
    headers: { Authorization: cookie }
  })
}

export function getSelectionResults(cookie) {
  return api.get('/v1/zhjw/selection', {
    headers: { Authorization: cookie }
  })
}

export function getClassSchedules(cookie) {
  return api.get('/v1/zhjw/schedule', {
    headers: { Authorization: cookie }
  })
}
