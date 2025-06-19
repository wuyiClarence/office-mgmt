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
            <p slot="title" v-if="type == 'role'">给角色{{ disname }}分配可操作的权限</p>
            <p slot="title" v-if="type == 'user'">给用户{{ disname }}分配可操作的权限</p>
        </Card>
    </div>
      <Tabs v-model="tabValue" @on-click="tabClick">
          <TabPane label="设备" name="device">
            <tables ref="deviceTables" pageable editable searchable search-place="top" v-model="tableData" :columns="tableColumns"
                :totalPage="permissionTotal" :pageSize="permissionPageParams.pageSize" :pageIndex="permissionPageParams.pageIndex" @on-pageChange="handlePermissionPageChange"/>
          </TabPane>
          <TabPane label="设备组" name="device_group">
            <tables ref="deviceGroupTables" pageable editable searchable search-place="top" v-model="tableData" :columns="tableColumns"
                :totalPage="permissionTotal" :pageSize="permissionPageParams.pageSize" :pageIndex="permissionPageParams.pageIndex" @on-pageChange="handlePermissionPageChange"/>
          </TabPane>
          <TabPane label="策略" name="policy">
            <tables ref="policys" pageable editable searchable search-place="top" v-model="tableData" :columns="tableColumns"
                :totalPage="permissionTotal" :pageSize="permissionPageParams.pageSize" :pageIndex="permissionPageParams.pageIndex" @on-pageChange="handlePermissionPageChange"/>
          </TabPane>
      </Tabs>
  </Modal>
  </template>
<script>
import Tables from '../tables'
import { getPermissionRoleOwn, getPermissionUserOwn, postPermissionRoleOwn, postPermissionUserOwn } from '@/api/permission'
export default {
  name: 'PermissionUserRoleForm',
  components: {
    Tables
  },
  props: {
    value: {
      type: Boolean,
      default: false
    },
    id: {
      type: Number,
      default: 0
    },
    type: {
      type: String
    },
    disname: {
      type: String
    }
  },
  data () {
    return {
      borderTitle: '',
      localValue: this.value,
      tabValue: 'device',
      permissionTotal: 0, // 总条数
      permissionPageParams: {
        pageSize: 10, // 每页显示的条数
        pageIndex: 1 // 当前页码
      },
      userKeysValue: [],
      userKeys: [],
      tableData: [],
      tableColumns: [
        { title: 'ID', width: 130, align: 'center', key: 'resource_id' },
        { title: '名称', width: 130, align: 'center', key: 'resource_name' },
        { title: '权限',
          align: 'center',
          key: 'permissions',
          render: (h, params) => {
            return h('CheckboxGroup', {
              props: {
                value: params.row.assigned_own_keys.map(item => item.key)
              },
              on: {
                'on-change': (keys) => {
                  const assignedKeys = params.row.assigned_own_keys
                  const updatedKeys = []

                  // 检查新增项
                  keys.forEach(key => {
                    if (!assignedKeys.some(el => el.key === key)) {
                      this.checkBoxChange(key, true, params.row.resource_id)
                      updatedKeys.push({ key, name: (params.row.user_own_keys.find(item => item.key === key) || {}).name })
                    }
                  })

                  // 保留删除逻辑
                  assignedKeys.forEach(item => {
                    if (!keys.includes(item.key)) {
                      this.checkBoxChange(item.key, false, params.row.resource_id)
                    } else {
                      updatedKeys.push(item)
                    }
                  })
                  // 更新数据
                  this.$set(params.row, 'assigned_own_keys', updatedKeys)
                }
              }
            }, params.row.user_own_keys.map(item => {
              return h('Checkbox', {
                props: {
                  label: item.key,
                  key: item.key
                }
              }, item.name)
            }))
          }
        }
      ]
    }
  },
  watch: {
    id (newId) {
      this.tabValue = 'device'
      this.getPermssion()
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
    getPermssion () {
      if (this.type === 'role') {
        let reqparams = {
          'role_id': this.id,
          'resource_type': this.tabValue,
          'pageIndex': this.permissionPageParams.pageIndex,
          'pageSize': this.permissionPageParams.pageSize
        }
        getPermissionRoleOwn(reqparams).then(res => {
          this.permissionTotal = res.total
          this.tableData = res.resource_permission
        })
      }
      if (this.type === 'user') {
        let reqparams = {
          'user_id': this.id,
          'resource_type': this.tabValue,
          'pageIndex': this.permissionPageParams.pageIndex,
          'pageSize': this.permissionPageParams.pageSize
        }
        getPermissionUserOwn(reqparams).then(res => {
          this.permissionTotal = res.total
          this.tableData = res.resource_permission
        })
      }
    },
    checkBoxChange (key, status, resourceId) {
      if (this.type === 'role') {
        let reqparams = {
          'role_id': this.id,
          'resource_id': resourceId,
          'resource_type': this.tabValue,
          'assigned_key': key,
          'enable': status
        }
        postPermissionRoleOwn(reqparams).then(res => {
          this.$Message.success('提交成功!')
        })
      }
      if (this.type === 'user') {
        let reqparams = {
          'user_id': this.id,
          'resource_id': resourceId,
          'resource_type': this.tabValue,
          'assigned_key': key,
          'enable': status
        }
        postPermissionUserOwn(reqparams).then(res => {
          this.$Message.success('提交成功!')
        })
      }
    },
    handlePermissionPageChange (params) {
      this.permissionPageParams = params
      this.getPermssion()
    },
    tabClick (name) {
      this.tabValue = name
      this.getPermssion()
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
