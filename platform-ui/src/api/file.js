import axios from '@/libs/api.request'

export const getFile = params => {
  return axios.request({
    url: 'file/list',
    method: 'get',
    params: params
  })
}

export const uploadFile = (data, config) => {
  return axios.request({
    url: 'file/upload',
    method: 'post',
    data,
    ...config
  })
}

export const deleteFile = (data) => {
  return axios.request({
    url: 'file/delete',
    method: 'delete',
    data
  })
}
