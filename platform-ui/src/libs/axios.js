import axios from 'axios'
import store from '@/store'
// import { Spin } from 'iview'
import { callRefreshToken } from '../api/user'
import { setToken, setRefreshToken } from '@/libs/util'
import { Message } from 'iview'

let requests = []
let isRefreshing = false

class HttpRequest {
  constructor (baseUrl = baseURL) {
    this.baseUrl = baseUrl
    this.queue = {}
  }
  getInsideConfig () {
    const config = {
      baseURL: this.baseUrl,
      headers: {
        //
      }
    }
    return config
  }
  destroy (url) {
    delete this.queue[url]
    if (!Object.keys(this.queue).length) {
      // Spin.hide()
    }
  }
  interceptors (instance, url) {
    // 请求拦截
    instance.interceptors.request.use(config => {
      // 添加全局的loading...
      if (!Object.keys(this.queue).length) {
        // Spin.show() // 不建议开启，因为界面不友好
      }

      this.queue[url] = true
      const token = store.state.user.token
      if (token) {
        // 添加 token 到请求头
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    }, error => {
      return Promise.reject(error)
    })
    // 响应拦截
    instance.interceptors.response.use(async (response) => {
      let res = response.data
      // 如果是login接口 直接放过
      if (response.config.url.includes('user/login')) {
        return res.data
      }
      if (res.code === 401) {
        if (!isRefreshing) {
          isRefreshing = true
          try {
            const refreshToken = store.state.user.refreshToken
            if (!refreshToken || refreshToken.length === 0) {
              throw new Error('Refresh token is missing, redirecting to login')
            }

            const refreshTokenRes = await callRefreshToken({ refresh_token: refreshToken })
            const newAccessToken = refreshTokenRes.data.data.access_token
            setToken({ token: newAccessToken })
            store.state.user.token = newAccessToken
            requests.forEach((callback) => callback(newAccessToken))
            requests = []

            response.config.headers['Authorization'] = `Bearer ${newAccessToken}`
            isRefreshing = false
            return instance(response.config)
          } catch (refreshError) {
            requests = []
            isRefreshing = false
            setToken('')
            setRefreshToken('')
            window.location.href = '/login'
            return Promise.reject(refreshError)
          }
        } else {
          return new Promise(resolve => {
            // 用函数形式将 resolve 存入，等待刷新后再执行
            requests.push(token => {
              response.config.headers.Authorization = `Bearer ${token}`
              resolve(instance(response.config))
            })
          })
        }
      } else if (res.code !== 0) {
        this.destroy(url)
        Message.error(res.msg)
        return Promise.reject(res.msg)
      }
      if (res.code !== 0 && res.code !== 401) {
        const errorMsg = res.msg || '未知错误'
        const error = new Error(errorMsg) // 创建一个错误对象
        error.code = res.code // 添加错误码
        return Promise.reject(error)
      }
      return res.data
    }, error => {
      window.location.href = '/login'
      return Promise.reject(error)
    })
  }
  request (options) {
    const instance = axios.create()
    options = Object.assign(this.getInsideConfig(), options)
    this.interceptors(instance, options.url)
    return instance(options)
  }
}
export default HttpRequest
