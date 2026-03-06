export function resolveRequestError(error, emptyMessage = '暂无数据') {
  const status = error?.response?.status
  const code = error?.response?.data?.code

  if (status === 401 || code === 401) {
    return { message: '', isEmpty: false }
  }

  if (status === 404 || code === 404) {
    return { message: '', isEmpty: true, emptyMessage }
  }

  if (error?.code === 'ERR_NETWORK' || !error?.response) {
    return { message: '网络异常，请检查网络连接', isEmpty: false }
  }

  return {
    message: error?.response?.data?.msg || error?.message || '请求失败，请稍后重试',
    isEmpty: false
  }
}
