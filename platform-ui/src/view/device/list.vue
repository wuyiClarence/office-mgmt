<template>
  <div>
    <CustomTables ref="tables" pageable editable searchable search-place="top"
      :columns="columns"
      :canDelete="canDelete" @on-delete="onDelete"
      :canAdd="canAdd" @on-add="onAdd"
      :canEdit="canEdit" @on-row-edit="OnRowEdit"
      @on-poweron="OnRowPowerOn"
      @on-poweroff="OnRowPowerOff"
      curResourceType = 'device'
      @on-tabledata="onTableData"/>
      <Modal
      :mask-closable = "false"
      v-model="modalAdd"
      :title="isAdd ? '添加设备' : '编辑设备'">
        <Form :model="formItem" :label-width="100" ref="formValidate">
          <Form-item label="设备名" prop="device_name" style="width: 400px">
              <Input v-model="formItem.device_name" :disabled=!isAdd placeholder="请输入设备名" ></Input>
          </Form-item>
          <Form-item label="别名" prop="alias_name" style="width: 400px">
            <Input v-model="formItem.alias_name" placeholder="请输入设备别名" ></Input>
          </Form-item>
          <Form-item label="类型"  style="width: 400px">
            <span v-if="formItem.device_type==1">物理机</span>
            <span v-if="formItem.device_type==2">kvm虚拟机</span>
          </Form-item>
          <Form-item label="Mac地址" prop="mac" style="width: 400px">
              <Input v-model="formItem.mac" :disabled="formItem.device_type==1" ></Input>
          </Form-item>
          <Form-item label="IP地址" prop="ip" style="width: 400px">
              <Input v-model="formItem.ip" :disabled="formItem.device_type==1" ></Input>
          </Form-item>
        </Form>
        <div slot="footer">
            <Button @click.native="OnCancel">取消</Button>
            <Button type="primary" @click.native="OnSave">保存</Button>
        </div>
      </Modal>
      <Modal
      v-model="detailVisible"
      title="设备详情"
      :footer-hide="true">
      <div v-for="deviceGroup in detailDeviceGroups" :key="deviceGroup.id">
        {{ deviceGroup.device_group_name }}
      </div>
    </Modal>
  </div>
</template>

<script>
import CustomTables from '_c/custom-tables'
import { getDeviceList, deleteDevice, powerOnDevice, powerOffDevice, updateDevice } from '@/api/device'
import { hasOneOf } from '@/libs/tools'
export default {
  name: 'server_page',
  components: {
    CustomTables
  },
  data () {
    return {
      columns: [
        {
          type: 'selection',
          width: 50,
          align: 'center'
        },
        { title: 'ID', width: 80, key: 'id', sortable: true },
        { title: '名称',
          width: 200,
          key: 'device_name',
          render: (h, { row }) => {
            return h('span', {}, row.alias_name.length > 0 ? row.alias_name : row.device_name)
          }
        },
        { title: 'Mac地址', width: 150, key: 'mac' },
        { title: 'IP地址', key: 'ip' },
        { title: '系统', key: 'os_type' },
        { title: '类型',
          key: 'device_type',
          render: (h, { row }) => {
            return h('span', {
            }, row.device_type === 1 ? '物理机' : 'kvm虚拟机')
          }
        },
        {
          title: '状态',
          key: 'status',
          render: (h, { row }) => {
            return h('span', {
              style: {
                color: row.status === 1 ? '#19be6b' : '#909399',
                borderRadius: '3px'
              }
            }, row.status === 1 ? '在线' : '离线')
          }
        },
        { title: '所属设备组',
          key: 'group',
          render: (h, params) => {
            const device_group_info = params.row.device_group_info
            return h('div', [
              h('div', device_group_info.slice(0, 3).map(device =>
                h('div', { style: { color: '#666' } },
                  `${device.device_group_name}`
                )
              )),
              device_group_info.length > 3 && h('a', { on: { click: () => this.showDetail(device_group_info) } }, '查看全部...')
            ])
          }
        },
        {
          title: '操作',
          width: 500,
          align: 'center',
          key: 'customhandle',
          options: ['edit', 'permission', 'delete', 'poweron', 'poweroff']
        }
      ],
      detailVisible: false,
      detailDeviceGroups: [],
      formItem: {
        id: 0
      },
      isAdd: true,
      modalAdd: false
    }
  },
  computed: {
    access () {
      return this.$store.state.user.access
    },
    canDelete () {
      return hasOneOf(['DEVICE_DELETE'], this.access)
    },
    canEdit () {
      return hasOneOf(['DEVICE_EDIT'], this.access)
    },
    canAdd () {
      return false
    }
  },
  methods: {
    onTableData (pageParams, callback) {
      getDeviceList(pageParams).then(res => {
        callback(res.list, res.total)
      })
    },
    onDelete (ids) {
      if (ids.length === 0) {
        this.$Message.error('请选择要删除的设备组!')
      } else {
        deleteDevice({ ids: ids }).then(res => {
          this.$Message.success('删除成功!')
          this.$refs.tables.getTableData()
        })
      }
    },
    onAdd () {
      // this.isAdd = true
      // this.$refs.formValidate.resetFields()
      // this.modalAdd = true
    },
    OnSave () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          if (!this.isAdd) {
            updateDevice(this.formItem).then(res => {
              this.$Message.success('提交成功!')
              this.$refs.formValidate.resetFields()
              this.modalAdd = false
              this.$refs.tables.getTableData()
            }).catch(() => {
            })
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
    OnRowEdit (row) {
      this.isAdd = false
      this.$refs.formValidate.resetFields()
      this.formItem.id = row.id
      this.formItem.device_name = row.device_name
      this.formItem.ip = row.ip
      this.formItem.mac = row.mac
      this.formItem.alias_name = row.alias_name
      this.formItem.device_type = row.device_type

      this.modalAdd = true
    },
    OnRowPowerOn (ids) {
      if (ids.length === 0) {
        this.$Message.error('请选择要开机的设备!')
      } else {
        powerOnDevice({ ids: ids }).then(res => {
          this.$Message.success('开机操作成功!')
        })
      }
    },
    OnRowPowerOff (ids) {
      if (ids.length === 0) {
        this.$Message.error('请选择要关机的设备!')
      } else {
        powerOffDevice({ ids: ids }).then(res => {
          this.$Message.success('关机操作成功!')
        })
      }
    },
    showDetail (deviceGroups) {
      this.detailDeviceGroups = deviceGroups
      this.detailVisible = true
    }
  },
  mounted () {
  }
}
</script>

<style>

</style>
