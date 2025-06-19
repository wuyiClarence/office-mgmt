const btns = {
  edit: (h, params, vm) => {
    if (!vm.canEdit) {
      return null
    }
    const hasPermission = params.row.permissions.find(permission => permission.key === 'edit')
    if (!hasPermission) {
      return null
    }
    return h('Button', {
      props: {
        type: 'primary' // 设置按钮样式
      },
      on: {
        click: () => {
          // 操作按钮的点击逻辑
          vm.onRowEdit(params.row)
        }
      },
      style: {
        width: '80px',
        marginLeft: '8px' // 设置两个按钮的间距
      }
    }, '编辑')
  },
  delete: (h, params, vm) => {
    if (!vm.canDelete) {
      return null
    }
    const hasPermission = params.row.permissions.find(permission => permission.key === 'delete')
    if (!hasPermission) {
      return null
    }
    return h('Poptip', {
      props: {
        confirm: true,
        transfer: true,
        title: '你确定要删除吗?' // 确认提示内容
      },
      on: {
        'on-ok': () => {
          vm.onRowDelete(params.row) // 用户点击确认后触发删除逻辑
        }
      }
    }, [
      h('Button', {
        props: {
          type: 'error', // 设置按钮样式
          icon: 'md-trash' // 按钮图标
        },
        style: {
          width: '80px',
          marginLeft: '8px' // 设置按钮的间距
        }
      }, '删除')
    ])
  },
  permission: (h, params, vm) => {
    const hasPermission = params.row.permissions.find(permission => permission.key === 'permissionmgmt')
    if (!hasPermission) {
      return null
    }
    return h('Button', {
      props: {
        type: 'primary' // 设置按钮样式
      },
      on: {
        click: () => {
          vm.accessPermissionClick(params)
          // vm.$emit('on-action', params) // 操作按钮的点击逻辑
        }
      },
      style: {
        width: '80px',
        marginLeft: '8px' // 设置按钮的间距
      }
    }, '管理权限')
  },
  poweron: (h, params, vm) => {
    if (!vm.canEdit) {
      return null
    }
    const hasPermission = params.row.permissions.find(permission => permission.key === 'poweron')
    if (!hasPermission) {
      return null
    }
    return h('Button', {
      props: {
        type: 'success' // 设置按钮样式
      },
      on: {
        click: () => {
          // 操作按钮的点击逻辑
          vm.onRowPowerOn(params.row)
        }
      },
      style: {
        width: '80px',
        marginLeft: '8px' // 设置两个按钮的间距
      }
    }, '开机')
  },
  poweroff: (h, params, vm) => {
    if (!vm.canEdit) {
      return null
    }
    const hasPermission = params.row.permissions.find(permission => permission.key === 'poweroff')
    if (!hasPermission) {
      return null
    }
    return h('Button', {
      props: {
        type: 'warning' // 设置按钮样式
      },
      on: {
        click: () => {
          // 操作按钮的点击逻辑
          vm.onRowPowerOff(params.row)
        }
      },
      style: {
        width: '80px',
        marginLeft: '8px' // 设置两个按钮的间距
      }
    }, '关机')
  }
}

export default btns
