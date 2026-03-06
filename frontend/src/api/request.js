import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true
})

function clearSessionAndRedirect() {
  localStorage.removeItem('zhjw_cookie')
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

api.interceptors.response.use(
  (response) => {
    if (response.data?.code === 401) {
      clearSessionAndRedirect()
      return Promise.reject(new Error('登录状态已过期，请重新登录'))
    }

    if (response.config.responseType === 'blob') {
      return response.data
    }

    return response.data
  },
  (error) => {
    const status = error.response?.status
    const code = error.response?.data?.code
    if (status === 401 || code === 401) {
      clearSessionAndRedirect()
    }

    return Promise.reject(error)
  }
)

export default api
