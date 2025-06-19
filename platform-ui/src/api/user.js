import axios from '@/libs/api.request'
import refreshTokenAxios from '@/libs/api.refreshtoken.request'

export const login = ({ userName, password }) => {
  const data = {
    userName,
    password
  }
  return axios.request({
    url: 'user/login',
    data,
    method: 'post'
  })
}

export const callRefreshToken = ({ refresh_token }) => {
  const data = {
    refresh_token
  }
  return refreshTokenAxios.request({
    url: 'user/refresh_token',
    data,
    method: 'post',
    timeout: 10000
  })
}

export const getUserInfo = () => {
  return axios.request({
    url: 'user/info',
    method: 'get'
  })
}

export const logout = (token) => {
  return axios.request({
    url: 'user/logout',
    method: 'post'
  })
}

export const getUnreadCount = () => {
  return axios.request({
    url: 'message/count',
    method: 'get'
  })
}

export const getMessage = () => {
  return axios.request({
    url: 'message/init',
    method: 'get'
  })
}

export const getContentByMsgId = msg_id => {
  return axios.request({
    url: 'message/content',
    method: 'get',
    params: {
      msg_id
    }
  })
}

export const hasRead = msg_id => {
  return axios.request({
    url: 'message/has_read',
    method: 'post',
    data: {
      msg_id
    }
  })
}

export const removeReaded = msg_id => {
  return axios.request({
    url: 'message/remove_readed',
    method: 'post',
    data: {
      msg_id
    }
  })
}

export const restoreTrash = msg_id => {
  return axios.request({
    url: 'message/restore',
    method: 'post',
    data: {
      msg_id
    }
  })
}

export const getUserList = () => {
  return axios.request({
    url: 'user/list',
    method: 'get'
  })
}

export const createUser = (data) => {
  return axios.request({
    url: 'user/create',
    method: 'post',
    data
  })
}

export const deleteUser = (data) => {
  return axios.request({
    url: 'user/delete',
    method: 'delete',
    data
  })
}

export const updateUser = (data) => {
  return axios.request({
    url: 'user/update_info',
    method: 'put',
    data
  })
}
