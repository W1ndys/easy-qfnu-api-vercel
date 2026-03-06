import api from './request'

export function login(data) {
  return api.post('/v1/zhjw/login', {
    username: data.username,
    password: data.password
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
