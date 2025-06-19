import axios from '@/libs/api.request'

export const getPolicyList = params => {
  return axios.request({
    url: 'policy/list',
    method: 'get',
    params: params
  })
}

export const updatePolicy = (data) => {
  return axios.request({
    url: 'policy/update',
    method: 'put',
    data
  })
}

export const createPolicy = (data) => {
  return axios.request({
    url: 'policy/create',
    method: 'post',
    data
  })
}

export const deletePolicy = (data) => {
  return axios.request({
    url: 'policy/delete',
    method: 'delete',
    data
  })
}
