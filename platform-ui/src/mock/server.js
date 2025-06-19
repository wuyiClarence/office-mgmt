import Mock from 'mockjs'
import { getParams, doCustomTimes } from '@/libs/util'

function randomIp () {
  return `${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}`
}

function randomMac () {
  const hexDigits = '0123456789ABCDEF'
  let mac = ''
  for (let i = 0; i < 6; i++) {
    if (i !== 0) {
      mac += ':' // 使用冒号分隔
    }
    mac += hexDigits.charAt(Math.floor(Math.random() * 16)) + hexDigits.charAt(Math.floor(Math.random() * 16))
  }
  return mac
}

function randomServerType () {
  const types = ['实体设备', '虚拟设备']
  return types[Math.floor(Math.random() * types.length)]
}

function randomStatus () {
  const types = ['运行中', '关机']
  return types[Math.floor(Math.random() * types.length)]
}

function randomGroup () {
  const types = ['办公室', '机房']
  return types[Math.floor(Math.random() * types.length)]
}

export const server_list = req => {
  let tableData = []
  const params = getParams(req.url)
  doCustomTimes(params.pageSize, () => {
    tableData.push(Mock.mock({
      name: '@name',
      mac: randomMac(),
      ipaddress: randomIp(),
      type: randomServerType(),
      status: randomStatus(),
      group: randomGroup()
    }))
  })
  let responsedata = {
    success: true,
    data: {
      total: 200,
      list: tableData
    },
    errcode: 0,
    errmessage: null
  }
  return responsedata
}

export const server_group_list = req => {
  let tableData = []
  const params = getParams(req.url)
  doCustomTimes(params.pageSize, () => {
    tableData.push(Mock.mock({
      groupname: '@name',
      serverlist: randomMac()
    }))
  })
  let responsedata = {
    success: true,
    data: {
      total: 20,
      list: tableData
    },
    errcode: 0,
    errmessage: null
  }
  return responsedata
}
