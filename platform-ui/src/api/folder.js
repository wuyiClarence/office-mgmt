import axios from '@/libs/api.request'

export const createFolder = (data) => {
  return axios.request({
    url: 'folder/create',
    method: 'post',
    data
  })
}

export const updateFolder = (data) => {
  return axios.request({
    url: 'folder/update',
    method: 'put',
    data
  })
}

export const deleteFolder = (data) => {
  return axios.request({
    url: 'folder/delete',
    method: 'delete',
    data
  })
}
