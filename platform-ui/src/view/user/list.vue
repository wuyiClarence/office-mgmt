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
      :title="isAdd ? '添加用户' : '编辑用户'">
      <Form :model="formItem" :label-width="80"  ref="formValidate" :rules="ruleValidate">
        <Form-item label="用户名" prop="user_name">
            <Input v-model="formItem.user_name" :disabled=!isAdd placeholder="请输入用户名"></Input>
        </Form-item>
        <Form-item type="password" label="密码" prop="password">
            <Input type="password"  v-model="formItem.password" placeholder="请输入密码"></Input>
        </Form-item>
        <Form-item type="password" label="确认密码" prop="passwordCheck">
            <Input type="password" v-model="formItem.passwordCheck" placeholder="请输入确认密码"></Input>
        </Form-item>
        <Form-item label="角色" prop="role_ids">
            <Select v-model="formItem.role_ids" multiple  placeholder="请选择">
              <Option v-for="item in roleList" :value="item.id" :key="item.id">{{ item.role_name }}</Option>
            </Select>
        </Form-item>
        <Form-item label="姓名" prop="user_display_name">
            <Input v-model="formItem.user_display_name" placeholder="请输入姓名"></Input>
        </Form-item>
        <Form-item label="手机号" prop="phone_number">
            <Input v-model="formItem.phone_number" placeholder="请输入电话"></Input>
        </Form-item>
        <Form-item label="邮箱" prop="email">
            <Input v-model="formItem.email" placeholder="请输入邮箱"></Input>
        </Form-item>
      </Form>
      <div slot="footer">
          <Button @click.native="OnCancel">取消</Button>
          <Button type="primary" @click.native="OnOk">保存</Button>
        </div>
    </Modal>
    <PermissionUserRoleForm v-model="modalAccess" :id="curResourceId" :type="'user'"  :disname="curResourceName" ></PermissionUserRoleForm>
  </div>
</template>

<script>
import PermissionUserRoleForm from '_c/permission-userrole-form'
import TipModal from '_c/tipmodal'
import Tables from '_c/tables'
import { getUserList, createUser, deleteUser, updateUser } from '@/api/user'
import { getRoleList } from '@/api/role'
import { hasOneOf } from '@/libs/tools'
export default {
  name: 'user_page',
  components: {
    Tables,
    TipModal,
    PermissionUserRoleForm
  },
  data () {
    return {
      curResourceId: 0,
      curResourceName: '',
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
        { title: 'ID', width: 80, key: 'user_id', sortable: true },
        { title: '账号', width: 130, align: 'center', key: 'user_name' },
        { title: '姓名', width: 130, align: 'center', key: 'user_display_name' },
        { title: '角色',
          minWidth: 150,
          align: 'center',
          key: 'role_infos',
          render: (h, params) => {
            let disrole = ''
            let i = 0
            params.row.role_infos.forEach(element => {
              if (i > 0) {
                disrole = disrole + ','
              }
              disrole = disrole + element.role_name
              i++
            })
            return h('div', disrole)
          }
        },
        { title: '邮箱', minWidth: 180, align: 'center', key: 'email' },
        { title: '手机号', minWidth: 110, align: 'center', key: 'phone_number' },
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
                // 删除按钮，带确认提示
                this.renderEditButton(h, params, vm),
                this.renderPermissionButton(h, params, vm),
                this.renderDeleteButton(h, params, vm)
                // 操作按钮
              ])
            }
          ]
        }
      ],
      modalAdd: false,
      ruleValidate: {
        user_name: [
          { required: true, message: '用户名不能为空', trigger: 'blur' }
        ],
        password: [
          {
            validator: (rule, value, callback) => {
              if (!value) {
                if (this.isAdd) {
                  callback(new Error('用户名不能为空'))
                }
                return callback()
              }

              if (value.length < 6) {
                callback(new Error('密码长度不能少于6位'))
              } else {
                callback()
              }
            },
            trigger: 'blur'
          }
        ],
        passwordCheck: [
          {
            validator: (rule, value, callback) => {
              if (!value) {
                if (this.isAdd) {
                  callback(new Error('用户名不能为空'))
                }
                if (!this.formItem.password) {
                  return callback()
                }
              }

              if (value.length < 6) {
                callback(new Error('密码长度不能少于6位'))
              } else if (value !== this.formItem.password) {
                callback(new Error('两次密码输入不一致'))
              } else {
                callback() // 验证通过
              }
            },
            trigger: 'blur'
          }
        ],
        email: [
          // { required: true, message: '邮箱不能为空', trigger: 'blur' },
          { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
        ],
        phone_number: [
          {
            validator: (rule, value, callback) => {
              let reg = /^1[3-9]\d{9}$/
              if (!value) {
                return callback()
              }
              if (!reg.test(value)) {
                return callback(new Error('手机号格式不正确'))
              } else {
                callback()
              }
            },
            trigger: 'blur'
          }
        ],
        role_ids: [
          { required: true, type: 'array', min: 1, message: '至少选择一个角色', trigger: 'change' }
        ],
        user_display_name: [
          { required: true, message: '姓名不能为空', trigger: 'blur' }
        ]
      },
      roleList: [
      ],
      formItem: {
        user_name: '',
        password: '',
        passwordCheck: '',
        role_ids: [],
        user_display_name: '',
        phone_number: '',
        email: ''
      },
      tableData: [],
      modalAccess: false,
      selectids: [],
      isAdd: true
    }
  },
  computed: {
    access () {
      return this.$store.state.user.access
    },
    canAdd () {
      return hasOneOf(['USER_ADD'], this.access)
    },
    canDelete () {
      return hasOneOf(['USER_DEL'], this.access)
    },
    canEdit () {
      return hasOneOf(['USER_EDIT'], this.access)
    }
  },
  watch: {
  // 监听 password 的变化
    'formItem.password' (newPassword) {
      // 只有在密码更改时才手动触发验证
      if (this.formItem.passwordCheck !== undefined && this.formItem.passwordCheck.length > 0) {
        this.$nextTick(() => {
          this.$refs['formValidate'].validateField('passwordCheck') // 手动验证 passwordCheck 字段
        })
      }
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
        this.selectids.push(element.user_id)
      })
    },
    handlePageChange (params) {
      this.pageParams = params
      this.getUserList() // 切换页码时重新获取数据
    },
    handleDelete () {
      if (this.selectids.length === 0) {
        this.$Message.error('请选择要删除的用户!')
      } else {
        this.deleteUser(this.selectids)
      }
    },
    handleRowEdit (row) {
      this.isAdd = false
      this.$refs.formValidate.resetFields()
      this.formItem.user_id = row.user_id
      this.formItem.user_name = row.user_name
      this.formItem.password = ''
      this.formItem.passwordCheck = ''
      row.role_infos.forEach(el => {
        this.formItem.role_ids.push(el.role_id)
      })
      this.formItem.user_display_name = row.user_display_name
      this.formItem.phone_number = row.phone_number
      this.formItem.email = row.email
      this.modalAdd = true
    },
    handleAdd () {
      this.isAdd = true
      this.$refs.formValidate.resetFields()
      this.modalAdd = true
    },
    handleRowDelete (row) {
      this.deleteUser([row.user_id])
    },
    deleteUser (userIds) {
      deleteUser({ user_ids: userIds }).then(res => {
        this.$Message.success('删除成功!')
        this.getUserList()
      })
    },
    OnOk () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          if (this.isAdd) {
            createUser(this.formItem).then(res => {
              this.$Message.success('提交成功!')
              this.$refs.formValidate.resetFields()
              this.modalAdd = false
              this.getUserList()
            }).catch(() => {
            })
          } else {
            updateUser(this.formItem).then(res => {
              this.$Message.success('提交成功!')
              this.$refs.formValidate.resetFields()
              this.modalAdd = false
              this.getUserList()
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
      return h('Button', {
        props: {
          type: 'primary' // 设置按钮样式
        },
        on: {
          click: () => {
            this.accessPermissionClick(params)
            // vm.$emit('on-action', params) // 操作按钮的点击逻辑
          }
        },
        style: {
          width: '80px',
          marginLeft: '8px' // 设置按钮的间距
        }
      }, '管理权限')
    },
    getUserList () {
      getUserList(this.pageParams).then(res => {
        this.tableData = res.list
        this.total = res.total
      })
    },
    getRoleList () {
      getRoleList().then(res => {
        this.roleList = []
        this.roleList = res.list
      })
    },
    accessPermissionCancel () {
      this.modalAccess = false
    },
    accessPermissionClick (params) {
      this.curResourceId = params.row.user_id
      this.curResourceType = 'role'
      this.curResourceName = params.row.user_name
      this.modalAccess = true
    }
  },
  mounted () {
    this.getUserList()
    this.getRoleList()
  }
}
</script>

<style>

</style>
