import axios from '@/libs/api.request'

export const getNameSpace = params => {
  return axios.request({
    url: 'namespace/list',
    method: 'get',
    params: params
  })
}

export const createNameSpace = (data) => {
  return axios.request({
    url: 'namespace/create',
    method: 'post',
    data
  })
}

export const deleteNameSpace = (data) => {
  return axios.request({
    url: 'namespace/delete',
    method: 'delete',
    data
  })
}

export const updateNameSpace = (data) => {
  return axios.request({
    url: 'namespace/update',
    method: 'put',
    data
  })
}
