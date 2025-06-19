import Mock from 'mockjs'
import { doCustomTimes } from '@/libs/util'
import orgData from './data/org-data'
import { treeData } from './data/tree-select'
const Random = Mock.Random

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

export const getTableData = req => {
  let tableData = []
  doCustomTimes(10, () => {
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

export const getDragList = req => {
  let dragList = []
  doCustomTimes(5, () => {
    dragList.push(Mock.mock({
      name: Random.csentence(10, 13),
      id: Random.increment(10)
    }))
  })
  return dragList
}

export const uploadImage = req => {
  return Promise.resolve()
}

export const getOrgData = req => {
  return orgData
}

export const getTreeSelectData = req => {
  return treeData
}
