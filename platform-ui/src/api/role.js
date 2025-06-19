import axios from '@/libs/api.request'

export const getRoleMenuPermission = (params) => {
  return axios.request({
    url: 'role/menu_permission',
    params,
    method: 'get'
  })
}

export const putRoleMenuPermission = (data) => {
  return axios.request({
    url: 'role/menu_permission',
    method: 'put',
    data
  })
}

export const getRoleList = (params) => {
  return axios.request({
    url: 'role/list',
    params,
    method: 'get'
  })
}

export const createRole = (data) => {
  return axios.request({
    url: 'role/create',
    method: 'post',
    data
  })
}

export const deleteRole = (data) => {
  return axios.request({
    url: 'role/delete',
    method: 'delete',
    data
  })
}

export const updateRole = (data) => {
  return axios.request({
    url: 'role/update',
    method: 'put',
    data
  })
}
