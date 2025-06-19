import axios from '@/libs/api.request'

export const getDeviceGroup = params => {
  return axios.request({
    url: 'device_group/list',
    method: 'get',
    params: params
  })
}

export const createDeviceGroup = (data) => {
  return axios.request({
    url: 'device_group/create',
    method: 'post',
    data
  })
}

export const deleteDeviceGroup = (data) => {
  return axios.request({
    url: 'device_group/delete',
    method: 'delete',
    data
  })
}

export const updateDeviceGroup = (data) => {
  return axios.request({
    url: 'device_group/update',
    method: 'put',
    data
  })
}

export const powerOnDeviceGroup = (data) => {
  return axios.request({
    url: 'device_group/poweron',
    method: 'post',
    data
  })
}

export const powerOffDeviceGroup = (data) => {
  return axios.request({
    url: 'device_group/poweroff',
    method: 'post',
    data
  })
}
