<template>
  <div>
    <Card>
      <tables ref="tables" pageable editable searchable search-place="top" v-model="tableData" :columns="columns"
      :addbtnable="canAdd" @on-delete="handleDelete"
      :delbtnable="canDelete" @on-add="handleAdd"
      @on-selection-change="onSelectionChange"
      :totalPage="total" :pageSize="pageParams.pageSize" :pageIndex="pageParams.pageIndex" @on-pageChange="handlePageChange"/>
    </Card>
  <Modal
      :mask-closable = "false"
      v-model="modalAdd"
      :scrollable = "false"
      :title="isAdd ? '添加角色' : '编辑角色'">
      <Form :model="formItem" :label-width="80"  ref="formValidate" :rules="ruleValidate">
        <Form-item label="角色名称" prop="role_name">
            <Input v-model="formItem.role_name" :disabled=!isAdd placeholder="请输入角色名称"></Input>
        </Form-item>
        <Form-item label="描述" prop="description">
            <Input v-model="formItem.description" placeholder="描述"></Input>
        </Form-item>
        <tree-table expand-key="name"
        :show-header="false"
        :expand-type="false"
        :is-fold="false"
        :selectable="false"
        :columns="treecolumns"
        :border="false"
        max-height="600px"
        :data="treeData" >
          <template slot="likes" slot-scope="scope">
            <!-- <Button @click="handle(scope)">123</Button> -->
            <i-switch  v-model="scope.row.checked" @on-change="(value)=>permissionChange(value,scope)"></i-switch>
          </template>
        </tree-table>
      </Form>
      <div slot="footer">
          <Button @click.native="OnCancel">取消</Button>
          <Button type="primary" @click.native="OnOk">保存</Button>
      </div>
  </Modal>
  <PermissionUserRoleForm v-model="modalAccess" :id="curResourceId" :type="'role'" :resourceType="curResourceType" :disname="curResourceName" :resourcePermissions="resourcePermissions"></PermissionUserRoleForm>
  </div>
</template>

<script>
import PermissionUserRoleForm from '_c/permission-userrole-form'
import Tables from '_c/tables'
import { getRoleList, getRoleMenuPermission, createRole, deleteRole, updateRole } from '@/api/role'
import { hasOneOf } from '@/libs/tools'
export default {
  name: 'tables_page',
  components: {
    Tables,
    PermissionUserRoleForm
  },
  data () {
    return {
      curResourceId: 0,
      curResourceType: 'role',
      curResourceName: '',
      resourcePermissions: [],
      selectedMenu: [],
      treecolumns: [
        {
          title: 'name',
          key: 'name',
          width: '200px'
        },
        {
          title: 'likes',
          key: 'checked',
          minWidth: '200px',
          type: 'template',
          template: 'likes'
        }
      ],
      treeData: [],
      modalAccess: false,
      total: 0, // 总条数
      pageParams: {
        pageSize: 10, // 每页显示的条数
        pageIndex: 1 // 当前页码
      },
      columns: [
        {
          type: 'selection',
          width: 50,
          align: 'center'
        },
        { title: 'ID', width: 80, align: 'center', key: 'id', sortable: true },
        { title: '角色名称', width: 130, align: 'center', key: 'role_name' },
        { title: '描述', width: 130, align: 'center', key: 'description' },
        {
          title: '创建者',
          key: 'create_user_name',
          minWidth: 200,
          align: 'center',
          render (h, params) {
            return h('tag', {
              props: {
                type: 'border',
                color: 'blue'
              }
            }, params.row.create_user_name + ' ' + '(' + params.row.create_user_display_name + ')')
          }
        },
        {
          title: '创建时间',
          minWidth: 200,
          align: 'center',
          key: 'created_at',
          render: (h, params) => {
            return h('div', this.formatDate(params.row.created_at))
          }
        },
        {
          title: '操作',
          width: 300,
          align: 'center',
          key: 'handle',
          fixed: 'right',
          button: [
            (h, params, vm) => {
              return h('div', [
                this.renderEditButton(h, params, vm),
                this.renderPermissionButton(h, params, vm),
                this.renderDeleteButton(h, params, vm)
              ])
            }
          ]
        }
      ],
      modalAdd: false,
      ruleValidate: {
        role_name: [
          { required: true, message: '角色名称不能为空', trigger: 'blur' }
        ]
      },
      formItem: {
        role_name: '',
        description: '',
        menu: []
      },
      tableData: [],
      selectids: [],
      isAdd: true
    }
  },
  computed: {
    access () {
      return this.$store.state.user.access
    },
    canAdd () {
      return hasOneOf(['ROLE_ADD'], this.access)
    },
    canDelete () {
      return hasOneOf(['ROLE_DEL'], this.access)
    },
    canEdit () {
      return hasOneOf(['ROLE_EDIT'], this.access)
    }
  },
  methods: {
    formatDate (date) {
      const convdate = new Date(date)
      return new Intl.DateTimeFormat('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
        // timeZone: 'UTC' // 如果需要使用本地时间，删除此行
      }).format(convdate)
    },
    onSelectionChange (selection) {
      this.selectids = []
      selection.forEach(element => {
        this.selectids.push(element.id)
      })
    },
    handlePageChange (params) {
      this.pageParams = params
      this.getRoleList() // 切换页码时重新获取数据
    },
    handleDelete () {
      if (this.selectids.length === 0) {
        this.$Message.error('请选择要删除的角色!')
      } else {
        this.deleteRole(this.selectids)
      }
    },
    handleRowEdit (row) {
      this.isAdd = false
      this.$refs.formValidate.resetFields()
      this.formItem.role_id = row.id
      this.formItem.role_name = row.role_name
      this.formItem.description = row.description
      getRoleMenuPermission({ 'role_id': row.id }).then(res => {
        this.treeData = []
        this.treeData = res.permission_tree
        this.buildTree(this.treeData, res.menu)
        this.modalAdd = true
      })
    },
    handleAdd () {
      this.treeData = []
      this.treeData = this.$store.state.user.permission_tree
      this.buildTree(this.treeData, this.$store.state.user.access)
      this.isAdd = true
      this.$refs.formValidate.resetFields()
      this.modalAdd = true
    },
    handleRowDelete (row) {
      this.deleteRole([row.id])
    },
    deleteRole (roleIds) {
      deleteRole({ role_ids: roleIds }).then(res => {
        this.$Message.success('删除成功!')
        this.getRoleList()
      })
    },
    OnOk () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          this.selectedMenu = []
          this.formItem.menu = []
          this.getTreeCheckedItem(this.treeData)
          this.formItem.menu = this.selectedMenu
          if (this.isAdd) {
            createRole(this.formItem).then(res => {
              this.$Message.success('提交成功!')
              this.$refs.formValidate.resetFields()
              this.modalAdd = false
              this.getRoleList()
            }).catch(() => {
            })
          } else {
            updateRole(this.formItem).then(res => {
              this.$Message.success('提交成功!')
              this.$refs.formValidate.resetFields()
              this.modalAdd = false
              this.getRoleList()
            }).catch(() => {
            })
            // 编辑用户
          }
        } else {
          this.$Message.error('表单验证失败!')
        }
      })
    },
    OnCancel () {
      this.$refs.formValidate.resetFields()
      this.modalAdd = false
    },
    renderEditButton (h, params, vm) {
      if (!this.canEdit) {
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
            this.handleRowEdit(params.row)
          }
        },
        style: {
          width: '80px',
          marginLeft: '8px' // 设置两个按钮的间距
        }
      }, '编辑')
    },
    renderDeleteButton (h, params, vm) {
      if (params.row.system_create) return null
      if (!this.canDelete) {
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
            this.handleRowDelete(params.row) // 用户点击确认后触发删除逻辑
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
    renderPermissionButton (h, params, vm) {
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
            // vm.$emit('on-action', params) // 操作按钮的点击逻辑
            this.accessPermissionClick(params)
          }
        },
        style: {
          width: '80px',
          marginLeft: '8px' // 设置按钮的间距
        }
      }, '管理权限')
    },
    updateChildPermission (tree, status) {
      for (const el of tree) {
        el.checked = status
        if (el.children !== undefined) {
          this.updateChildPermission(el.children, status)
        }
      }
    },
    updateTreePermission (tree, key, status) {
      for (const el of tree) {
        if (el.key === key) {
          el.checked = status
          if (status === false) {
            if (el.children !== undefined) {
              this.updateChildPermission(el.children, status)
            }
          }
          return true
        }
        if (el.children !== undefined) {
          const updated = this.updateTreePermission(el.children, key, status)
          if (status === true && updated === true) {
            el.checked = true
          }
        }
      }
    },
    permissionChange (status, scope) {
      var data = this.treeData
      this.treeData = []
      this.$nextTick(() => {
        this.updateTreePermission(data, scope.row.key, status)
        this.treeData = data
      })
    },
    buildTree (treeTop, keys) {
      treeTop.forEach(el => {
        el.expand = true
        el.title = el.name
        if (keys.includes(el.key)) {
          el.checked = true
        } else {
          el.checked = false
        }
        if (el.children !== undefined) {
          this.buildTree(el.children, keys)
        }
      })
    },
    getTreeCheckedItem (treeTop) {
      treeTop.forEach(el => {
        if (el.checked === true) {
          this.selectedMenu.push(el.key)
          if (el.children !== undefined) {
            this.getTreeCheckedItem(el.children)
          }
        }
      })
    },
    getRoleList () {
      getRoleList(this.pageParams).then(res => {
        this.tableData = res.list
        this.total = res.total
      })
    },
    accessPermissionCancel () {
      this.modalAccess = false
    },
    accessPermissionClick (params) {
      this.curResourceId = params.row.id
      this.curResourceType = 'role'
      this.curResourceName = params.row.role_name
      this.resourcePermissions = params.row.permissions
      this.modalAccess = true
    }
  },
  mounted () {
    this.getRoleList()
  }
}
</script>

<style>

</style>
