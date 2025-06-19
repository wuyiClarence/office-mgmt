<template>
  <div>
    <CustomTables ref="tables" pageable editable searchable search-place="top"
      :columns="columns"
      :canDelete="canDelete" @on-delete="onDelete"
      :canAdd="canAdd" @on-add="onAdd"
      :canEdit="canEdit" @on-row-edit="OnRowEdit"
      @on-poweron="OnRowPowerOn"
      @on-poweroff="OnRowPowerOff"
      curResourceType = 'device_group'
      @on-tabledata="onTableData"/>
    <Modal  width="1500"
      :mask-closable = "false"
      v-model="modalAdd"
      :title="isAdd ? '添加设备组' : '编辑设备组'">
      <Form :model="formItem" :label-width="100"  ref="formValidate" :rules="ruleValidate">
        <Form-item label="设备组名" prop="device_group_name">
            <Input v-model="formItem.device_group_name" :disabled=!isAdd placeholder="请输入设备组名" style="width: 400px"></Input>
        </Form-item>
        <Form-item label="组内设备" prop="device_group_name">
          <CustomTransfer :columns="deviceColumns" :titles="['已选', '待选']" :operations="['新增', '移除']"
          v-model="modalTransfer"
          @on-selectTableData="onSelectTableData" @on-transfer="onTransferData" @on-tabledata="onTransferTableData"></CustomTransfer>
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
      <div v-for="device in detailDevices" :key="device.id">
        {{ device.device_name }} （{{ device.mac }}）
      </div>
    </Modal>
  </div>
</template>

<script>
import CustomTables from '_c/custom-tables'
import CustomTransfer from '_c/custom-transfer'
import { getDeviceGroup, createDeviceGroup, updateDeviceGroup, deleteDeviceGroup, powerOnDeviceGroup, powerOffDeviceGroup } from '@/api/group'
import { hasOneOf } from '@/libs/tools'
import { getDeviceList } from '@/api/device'
export default {
  name: 'device_group_page',
  components: {
    CustomTransfer,
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
        { title: '组名称', key: 'device_group_name', sortable: true },
        { title: '设备数',
          key: 'device_number',
          render: (h, params) => {
            const devices = params.row.devices
            return h('div', [
              devices.length === 0 && h('span', { style: { fontWeight: 'bold' } }, '无设备'),
              devices.length > 0 && h('span', { style: { fontWeight: 'bold' } }, devices.length + '台设备')
            ])
          }
        },
        { title: '设备',
          key: 'devicelist',
          render: (h, params) => {
            const devices = params.row.devices
            return h('div', [
              h('div', devices.slice(0, 3).map(device =>
                h('div', { style: { color: '#666' } },
                  device.mac.length > 0 ? `${device.device_name} (${device.mac})` : `${device.device_name}`
                )
              )),
              devices.length > 3 && h('a', { on: { click: () => this.showDetail(devices) } }, '查看全部...')
            ])
          }
        },
        {
          title: '操作',
          key: 'customhandle',
          width: 500,
          options: ['edit', 'permission', 'delete', 'poweron', 'poweroff']
        }
      ],
      ruleValidate: {
        device_group_name: [
          { required: true, message: '设备组名不能为空', trigger: 'blur' }
        ]
      },
      formItem: {
        id: 0,
        device_group_name: '',
        device_ids: []
      },
      isAdd: true,
      modalAdd: false,
      modalTransfer: '',
      deviceColumns: [
        {
          type: 'selection',
          width: 50,
          align: 'center'
        },
        { title: 'ID', width: 80, key: 'id', sortable: true },
        { title: '设备名称',
          key: 'device_name',
          render: (h, { row }) => {
            return h('span', {}, row.alias_name.length > 0 ? row.alias_name : row.device_name)
          }
        },
        { title: 'Mac地址', width: 150, key: 'mac' },
        { title: 'IP地址', key: 'ip' },
        { title: '类型',
          key: 'device_type',
          render: (h, { row }) => {
            return h('span', {
            }, row.device_type === 1 ? '物理机' : 'kvm虚拟机')
          }
        }
      ],
      selectTableData: [],
      transferIds: [],
      detailVisible: false,
      detailDevices: []
    }
  },
  computed: {
    access () {
      return this.$store.state.user.access
    },
    canAdd () {
      return hasOneOf(['DEVICE_GROUP_ADD'], this.access)
    },
    canDelete () {
      return hasOneOf(['DEVICE_GROUP_DEL'], this.access)
    },
    canEdit () {
      return hasOneOf(['DEVICE_GROUP_EDIT'], this.access)
    }
  },
  methods: {
    onTableData (pageParams, callback) {
      getDeviceGroup(pageParams).then(res => {
        callback(res.list, res.total)
      })
    },
    onDelete (ids) {
      if (ids.length === 0) {
        this.$Message.error('请选择要删除的设备组!')
      } else {
        deleteDeviceGroup({ ids: ids }).then(res => {
          this.$Message.success('删除成功!')
          this.$refs.tables.getTableData()
        })
      }
    },
    onAdd () {
      this.isAdd = true
      this.$refs.formValidate.resetFields()
      this.selectTableData = []
      this.modalAdd = true
      this.modalTransfer = 'show'
    },
    OnSave () {
      this.$refs.formValidate.validate((valid) => {
        if (valid) {
          if (this.isAdd) {
            let reqparams = {
              'device_group_name': this.formItem.device_group_name,
              'device_ids': this.transferIds
            }
            createDeviceGroup(reqparams).then(res => {
              this.$Message.success('提交成功!')
              this.$refs.formValidate.resetFields()
              this.modalAdd = false
              this.$refs.tables.getTableData()
            }).catch(() => {
            })
          } else {
            let reqparams = {
              'id': this.formItem.id,
              'device_ids': this.transferIds
            }
            updateDeviceGroup(reqparams).then(res => {
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
      this.modalTransfer = ''
    },
    OnRowEdit (row) {
      this.isAdd = false
      this.$refs.formValidate.resetFields()
      this.modalAdd = true
      this.modalTransfer = 'edit' + row.id
      this.formItem.id = row.id
      this.formItem.device_group_name = row.device_group_name
      this.selectTableData = JSON.parse(JSON.stringify(row.devices))
    },
    onSelectTableData (callback) {
      callback(this.selectTableData)
    },
    onTransferData (transferIds) {
      this.transferIds = transferIds
    },
    onTransferTableData (pageParams, callback) {
      getDeviceList(pageParams).then(res => {
        callback(res.list, res.total)
      })
    },
    showDetail (devices) {
      this.detailDevices = devices
      this.detailVisible = true
    },
    OnRowPowerOn (ids) {
      if (ids.length === 0) {
        this.$Message.error('请选择要开机的设备组!')
      } else {
        powerOnDeviceGroup({ ids: ids }).then(res => {
          this.$Message.success('开机操作成功!')
        })
      }
    },
    OnRowPowerOff (ids) {
      if (ids.length === 0) {
        this.$Message.error('请选择要关机的设备组!')
      } else {
        powerOffDeviceGroup({ ids: ids }).then(res => {
          this.$Message.success('关机操作成功!')
        })
      }
    }
  },
  mounted () {
  }
}
</script>
