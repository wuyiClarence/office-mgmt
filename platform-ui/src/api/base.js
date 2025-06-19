import axios from '@/libs/api.request'
/*
 * 统一的post请求封装
 */
export function post (url, params) {
  let request = new Promise((resolve, reject) => {
    // 设置请求头的Content-Type
    axios.defaults.headers.post['Content-Type'] = 'application/json;charset=utf-8'
    axios.post(url, JSON.stringify(params)).then(response => {
      // response只返回data中的数据
      // 在axios from '../plugins/request'文件下86行，return dataAxios.data;
      // 单单返回data值
      if (response.code === 0) { // code是返回码，0代表成功
        resolve(response)
      } else {
        this.$Message.error(response.msg) // msg是返回信息，提示框弹出
      }
    }).catch(err => {
      reject(err.data)
    })
  })
  return request
}

/*
 * 统一封装get请求
 */
export function get (url, params) {
  return new Promise((resolve, reject) => {
    axios.request({
      url: url,
      params,
      method: 'get'
    }).then(res => {
      if (res.data.success === true) {
        resolve(res.data)
      } else {
        this.$Message.error(res.data.errmessage)
      }
    }).catch(err => {
      reject(err)
    })
  })
}
