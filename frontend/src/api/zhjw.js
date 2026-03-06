import api from './request'

const ZHJW_BASE = 'http://zhjw.qfnu.edu.cn'

export function getInitCookie() {
  return api.get('/v1/zhjw/init')
}

export function login(data) {
  return api.post('/v1/zhjw/login', data)
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

export function getZhjwCaptchaUrl() {
  return `${ZHJW_BASE}/jsxsd/verifycode.servlet?t=${Date.now()}`
}
