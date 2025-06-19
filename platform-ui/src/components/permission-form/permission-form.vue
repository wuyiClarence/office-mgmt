<template>
  <Modal
      :mask-closable = "false"
      v-model="localValue"
      :scrollable = "false"
      :closable="false"
      title="权限设置"
      ok-text="关闭"
      width="800"
      @on-ok="modalOk"
      @on-cancel="modalCancel">
      <div>
        <Card :bordered="false">
            <p slot="title" v-if="resourceType=='role'">角色:{{ resourceName }}</p>
            <p slot="title" v-if="resourceType=='user'">用户:{{ resourceName }}</p>
            <p>管理用户和角色可以对{{ resourceName }}进行的操作 </p>
        </Card>
    </div>
      <Tabs v-model="tabValue" @on-click="tabClick">
          <TabPane label="所有人" name="alluser" v-if="canSetAllUser">
            <CheckboxGroup v-model="userKeysValue" @on-change="checkGroupChange">
              <Checkbox v-for="item in resourcePermissions"  :key="item.key" :label="item.key">{{ item.name }}</Checkbox>
            </CheckboxGroup>
          </TabPane>
          <TabPane label="指定用户" name="user">
            <tables ref="userPermissionTables" pageable editable searchable search-place="top" v-model="userPermissionTableData" :columns="userPermissionColumns"
                :totalPage="permissionTotal" :pageSize="permissionPageParams.pageSize" :pageIndex="permissionPageParams.pageIndex" @on-pageChange="handlePermissionPageChange"/>
          </TabPane>
          <TabPane label="指定角色" name="role">
            <tables ref="rolePermissionTables" pageable editable searchable search-place="top" v-model="rolePermissionTableData" :columns="rolePermissionColumns"
                :totalPage="permissionTotal" :pageSize="permissionPageParams.pageSize" :pageIndex="permissionPageParams.pageIndex" @on-pageChange="handlePermissionPageChange"/>
          </TabPane>
      </Tabs>
  </Modal>
  </template>
<script>
import Tables from '../tables'
import { getPermissionAllUser, getPermissionUser, getPermissionRole, putPermissionAllUser, putPermissionUser, putPermissionRole } from '@/api/permission'
export default {
  name: 'PermissionForm',
  components: {
    Tables
  },
  props: {
    value: {
      type: Boolean,
      default: false
    },
    resourceId: {
      type: Number,
      default: 0
    },
    resourceType: {
      type: String
    },
    resourceName: {
      type: String
    },
    resourcePermissions: {
      type: Array
    }
  },
  data () {
    return {
      canSetAllUser: 'false',
      borderTitle: '',
      localValue: this.value,
      tabValue: 'alluser',
      permissionTotal: 0, // 总条数
      permissionPageParams: {
        pageSize: 10, // 每页显示的条数
        pageIndex: 1 // 当前页码
      },
      userKeysValue: [],
      userKeys: [],
      userPermissionTableData: [],
      userPermissionColumns: [
        { title: '账号', width: 130, align: 'center', key: 'user_name' },
        { title: '姓名', width: 130, align: 'center', key: 'user_display_name' },
        { title: '权限',
          align: 'center',
          key: 'permissions',
          render: (h, params) => {
            const self = this
            return h('CheckboxGroup', {
              props: {
                value: params.row.keys
              },
              on: {
                'on-change': (keys) => {
                  keys.forEach(el => {
                    if (params.row.keys.indexOf(el) === -1) {
                      // add 权限
                      this.userCheckBoxChange(el, true, params.row.user_id)
                    }
                  })
                  params.row.keys.forEach(el => {
                    if (keys.indexOf(el) === -1) {
                      // delete 权限
                      this.userCheckBoxChange(el, false, params.row.user_id)
                    }
                  })
                  this.$set(params.row, 'keys', keys)
                }
              }
            }, self.resourcePermissions.map(item => {
              return h('Checkbox', {
                props: {
                  label: item.key,
                  key: item.key
                }
              }, item.name)
            }))
          }
        }
      ],
      rolePermissionTableData: [],
      rolePermissionColumns: [
        { title: '角色名称', width: 130, align: 'center', key: 'role_name' },
        { title: '权限',
          align: 'center',
          key: 'permissions',
          render: (h, params) => {
            const self = this
            return h('CheckboxGroup', {
              props: {
                value: params.row.keys
              },
              on: {
                'on-change': (keys) => {
                  keys.forEach(el => {
                    if (params.row.keys.indexOf(el) === -1) {
                      // add 权限
                      this.roleCheckBoxChange(el, true, params.row.role_id)
                    }
                  })
                  params.row.keys.forEach(el => {
                    if (keys.indexOf(el) === -1) {
                      // delete 权限
                      this.roleCheckBoxChange(el, false, params.row.role_id)
                    }
                  })
                  this.$set(params.row, 'keys', keys)
                }
              }
            }, self.resourcePermissions.map(item => {
              return h('Checkbox', {
                props: {
                  label: item.key
                }
              }, item.name)
            }))
          }
        }
      ]
    }
  },
  watch: {
    resourceId (newId) {
      this.tabValue = 'alluser'
      this.getAllUserPermission()
    },
    value (newVal) {
      this.localValue = newVal
    },
    localValue (newVal) {
      this.$emit('input', newVal)
    }
  },
  computed: {
  },
  methods: {
    userCheckBoxChange (key, status, userId) {
      let reqparams = {
        'resource_id': this.resourceId,
        'resource_type': this.resourceType,
        'permission_key': key,
        'enable': status,
        'user_id': userId
      }
      putPermissionUser(reqparams).then(res => {
        this.$Message.success('提交成功!')
      })
    },
    roleCheckBoxChange (key, status, roleId) {
      let reqparams = {
        'resource_id': this.resourceId,
        'resource_type': this.resourceType,
        'permission_key': key,
        'enable': status,
        'role_id': roleId
      }
      putPermissionRole(reqparams).then(res => {
        this.$Message.success('提交成功!')
      })
    },
    checkGroupChange (data) {
      this.userKeys.forEach(oldval => {
        var indexnotfound = this.userKeysValue.indexOf(oldval)
        if (indexnotfound === -1) {
          // delete
          let reqparams = {
            'resource_id': this.resourceId,
            'resource_type': this.resourceType,
            'permission_key': oldval,
            'enable': false
          }
          putPermissionAllUser(reqparams).then(res => {
            this.$Message.success('提交成功!')
            this.getAllUserPermission()
          })
        }
      })
      this.userKeysValue.forEach(newval => {
        var indexnotfound = this.userKeys.indexOf(newval)
        if (indexnotfound === -1) {
          // add
          let reqparams = {
            'resource_id': this.resourceId,
            'resource_type': this.resourceType,
            'permission_key': newval,
            'enable': true
          }
          putPermissionAllUser(reqparams).then(res => {
            this.$Message.success('提交成功!')
            this.getAllUserPermission()
          })
        }
      })
    },
    handlePermissionPageChange (params) {
      this.permissionPageParams = params
      if (this.tabValue === 'alluser') {
        this.getAllUserPermission()
      } else if (this.tabValue === 'user') {
        this.getUserPermission()
      } else if (this.tabValue === 'role') {
        this.getRolePermission()
      }
    },
    getAllUserPermission () {
      let reqparams = {
        'resource_id': this.resourceId,
        'resource_type': this.resourceType,
        'pageIndex': this.permissionPageParams.pageIndex,
        'pageSize': this.permissionPageParams.pageSize
      }
      getPermissionAllUser(reqparams).then(res => {
        if (res.can_set_alluser === false) {
          this.tabValue = 'user'
          this.canSetAllUser = false
          this.getUserPermission()
        } else {
          this.canSetAllUser = true
        }
        this.userKeysValue = res.keys
        this.userKeys = res.keys
      })
    },
    getUserPermission () {
      let reqparams = {
        'resource_id': this.resourceId,
        'resource_type': this.resourceType,
        'pageIndex': this.permissionPageParams.pageIndex,
        'pageSize': this.permissionPageParams.pageSize
      }
      getPermissionUser(reqparams).then(res => {
        this.permissionTotal = res.total
        this.userPermissionTableData = res.user_permission
      })
    },
    getRolePermission () {
      let reqparams = {
        'resource_id': this.resourceId,
        'resource_type': this.resourceType,
        'pageIndex': this.permissionPageParams.pageIndex,
        'pageSize': this.permissionPageParams.pageSize
      }
      getPermissionRole(reqparams).then(res => {
        this.permissionTotal = res.total
        this.rolePermissionTableData = res.role_permission
      })
    },
    tabClick (name) {
      this.permissionTotal = 0
      if (name === 'alluser') {
        this.getAllUserPermission()
      } else if (name === 'user') {
        this.getUserPermission()
      } else if (name === 'role') {
        this.getRolePermission()
      }
    },
    modalOk () {
    },
    modalCancel () {
    }
  },
  created () {
  }
}
</script>
