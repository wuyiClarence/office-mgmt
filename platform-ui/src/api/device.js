import axios from '@/libs/api.request'

export const getDeviceList = params => {
  return axios.request({
    url: 'device/list',
    method: 'get',
    params: params
  })
}

export const deleteDevice = (data) => {
  return axios.request({
    url: 'device/delete',
    method: 'delete',
    data
  })
}

export const powerOnDevice = (data) => {
  return axios.request({
    url: 'device/poweron',
    method: 'post',
    data
  })
}

export const powerOffDevice = (data) => {
  return axios.request({
    url: 'device/poweroff',
    method: 'post',
    data
  })
}

export const updateDevice = (data) => {
  return axios.request({
    url: 'device/update',
    method: 'put',
    data
  })
}
