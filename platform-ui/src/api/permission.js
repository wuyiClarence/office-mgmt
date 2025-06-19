import axios from '@/libs/api.request'

export const getPermissionInfo = (params) => {
  return axios.request({
    url: 'permission/info',
    params,
    method: 'get'
  })
}

export const getPermissionAllUser = (params) => {
  return axios.request({
    url: 'permission/alluser',
    params,
    method: 'get'
  })
}

export const getPermissionUser = (params) => {
  return axios.request({
    url: 'permission/user',
    params,
    method: 'get'
  })
}

export const getPermissionRole = (params) => {
  return axios.request({
    url: 'permission/role',
    params,
    method: 'get'
  })
}

export const getPermissionUserOwn = (params) => {
  return axios.request({
    url: 'permission/user/own',
    params,
    method: 'get'
  })
}

export const getPermissionRoleOwn = (params) => {
  return axios.request({
    url: 'permission/role/own',
    params,
    method: 'get'
  })
}
export const postPermissionUserOwn = (data) => {
  return axios.request({
    url: 'permission/user/own',
    data,
    method: 'post'
  })
}
export const postPermissionRoleOwn = (data) => {
  return axios.request({
    url: 'permission/role/own',
    data,
    method: 'post'
  })
}
export const putPermissionAllUser = (data) => {
  return axios.request({
    url: 'permission/alluser',
    data,
    method: 'put'
  })
}

export const putPermissionUser = (data) => {
  return axios.request({
    url: 'permission/user',
    data,
    method: 'put'
  })
}

export const putPermissionRole = (data) => {
  return axios.request({
    url: 'permission/role',
    data,
    method: 'put'
  })
}
