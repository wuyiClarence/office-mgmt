// import Mock from 'mockjs'
// import {
//   getParams,
//   doCustomTimes
// } from '@/libs/util'

// function randomMac () {
//   const hexDigits = '0123456789ABCDEF'
//   let mac = ''
//   for (let i = 0; i < 6; i++) {
//     if (i !== 0) {
//       mac += ':' // 使用冒号分隔
//     }
//     mac += hexDigits.charAt(Math.floor(Math.random() * 16)) + hexDigits.charAt(Math.floor(Math.random() * 16))
//   }
//   return mac
// }

// export const policy_list = req => {
//   let tableData = []
//   const params = getParams(req.url)
//   doCustomTimes(params.pageSize, () => {
//     tableData.push(Mock.mock({
//       name: '@name',
//       serverlist: randomMac()
//     }))
//   })
//   let responsedata = {
//     success: true,
//     data: {
//       total: 20,
//       list: tableData
//     },
//     errcode: 0,
//     errmessage: null
//   }
//   return responsedata
// }
