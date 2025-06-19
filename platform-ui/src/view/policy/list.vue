<template>
  <div>
    <CustomTables ref="tables" pageable editable searchable search-place="top"
      :columns="columns"
      :canDelete="canDelete" @on-delete="onDelete"
      :canAdd="canAdd" @on-add="onAdd"
      :canEdit="canEdit" @on-row-edit="OnRowEdit"
      curResourceType = 'policy'
      @on-tabledata="onTableData"/>
      <Modal width="1500"
      :mask-closable = "false"
      v-model="modalAdd"
      :title="isAdd ? '添加策略' : '编辑策略'">
        <Form :model="formItem" :label-width="100" ref="formValidate" :rules="ruleValidate">
          <Form-item label="策略名" prop="policy_name" style="width: 400px">
              <Input v-model="formItem.policy_name" :disabled=!isAdd placeholder="请输入策略名" ></Input>
          </Form-item>
          <Form-item label="状态" prop="status">
            <i-switch v-model="formItem.status">
              <span slot="open">开启</span>
              <span slot="close">关闭</span>
            </i-switch>
          </Form-item>
          <Form-item label="执行动作" prop="action_type" style="width: 400px">
            <Select v-model="formItem.action_type">
              <Option v-for="item in actionTypes" :value="item.value" :key="item.value">{{ item.label }}</Option>
            </Select>
          </Form-item>
          <Form-item label="执行时间" prop="execute_time">
            <TimePicker v-model="formItem.execute_time" @on-change="handleExecuteTimeChange" format="HH:mm:ss" placeholder="选择时间" style="width: 168px"></TimePicker>
          </Form-item>
          <Form-item label="重复" prop="execute_type">
            <RadioGroup v-model.number="formItem.execute_type" type="button">
              <Radio :label="1"><span>仅一次</span></Radio>
              <Radio :label="2"><span>每天</span></Radio>
              <Radio :label="3"><span>星期几</span></Radio>
            </RadioGroup>
          </Form-item>
          <Form-item v-if="formItem.execute_type == 3" label="指定星期" prop="week">
            <CheckboxGroup v-model="formItem.week">
              <Checkbox  label="Monday">
                <span>星期一</span>
              </Checkbox>
              <Checkbox label="Tuesday">
                <span>星期二</span>
              </Checkbox>
              <Checkbox label="Wednesday">
                <span>星期三</span>
              </Checkbox>
              <Checkbox label="Thursday">
                <span>星期四</span>
              </Checkbox>
              <Checkbox label="Friday">
                <span>星期五</span>
              </Checkbox>
              <Checkbox label="Saturday">
                <span>星期六</span>
              </Checkbox>
              <Checkbox label="Sunday">
                <span>星期日</span>
              </Checkbox>
            </CheckboxGroup>
          </Form-item>
          <Form-item label="限定日期" prop="datetime">
            <Col span="12">
              <DatePicker v-model="formItem.datetime" value-format="yyyy-MM-dd HH:mm:ss"  type="daterange" split-panels show-week-numbers confirm placement="bottom-end" placeholder="指定日期" style="width: 200px"></DatePicker>
            </Col>
          </Form-item>
          <Form-item label="关联" prop="associate_type">
            <Col span="12">
              <RadioGroup v-model="formItem.associate_type"  @on-change="onRadioChange">
                <Radio :label="1">
                    <span>设备</span>
                </Radio>
                <Radio :label="2">
                    <span>设备组</span>
                </Radio>
            </RadioGroup>
            </Col>
          </Form-item>
          <Form-item v-show="formItem.associate_type == 1" label="设备" prop="">
            <CustomTransfer :columns="deviceColumns" :titles="['已选', '待选']" :operations="['新增', '移除']"
            v-model="modalTransfer"
            @on-selectTableData="onSelectDeviceTableData" @on-transfer="onTransferDeviceData" @on-tabledata="onTransferDeviceTableData"></CustomTransfer>
          </Form-item>
          <Form-item v-show="formItem.associate_type == 2" label="设备组" prop="">
            <CustomTransfer :columns="deviceGroupColumns" :titles="['已选', '待选']" :operations="['新增', '移除']"
            v-model="modalTransfer"
            @on-selectTableData="onSelectDeviceGroupTableData" @on-transfer="onTransferDeviceGroupData" @on-tabledata="onTransferDeviceGroupTableData"></CustomTransfer>
          </Form-item>
        </Form>
        <div slot="footer">
            <Button @click.native="OnCancel">取消</Button>
            <Button type="primary" @click.native="OnSave">保存</Button>
        </div>
      </Modal>
    <Modal
      v-model="detailDeviceVisible"
      title="设备详情"
      :footer-hide="true">
      <div v-for="device in detailDevices" :key="device.id">
        {{ device.device_name }} （{{ device.mac }}）
      </div>
    </Modal>
    <Modal
      v-model="detailDeviceGroupVisible"
      title="设备组详情"
      :footer-hide="true">
      <div v-for="deviceGroup in detailDeviceGroups" :key="deviceGroup.id">
      {{ deviceGroup.device_group_name }}
      </div>
    </Modal>
  </div>
</template>

<script>
import dayjs from 'dayjs'
import dayjsPluginUTC from 'dayjs-plugin-utc'
import CustomTables from '_c/custom-tables'
import CustomTransfer from '_c/custom-transfer'
import { getPolicyList, createPolicy, updatePolicy, deletePolicy } from '@/api/policy'
import { hasOneOf } from '@/libs/tools'
import { getDeviceList } from '@/api/device'
import { getDeviceGroup } from '@/api/group'
// import { getDeviceGroup } from '@/api/group'
export default {
  name: 'server_page',
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
        { title: '策略名称', width: 100, key: 'policy_name', sortable: true, editable: true },
        { title: '执行时间',
          width: 100,
          key: 'execute_time'
        },
        { title: '执行动作',
          key: 'action_type',
          width: 100,
          render: (h, { row }) => {
            const action = this.actionTypes.find(item => item.value === row.action_type)
            return h('span', {}, action ? action.label : '未知操作')
          }
        },
        {
          title: '状态',
          key: 'status',
          width: 80,
          render: (h, { row }) => {
            return h('span', {
              style: {
                color: row.status === 1 ? '#19be6b' : '#909399',
                borderRadius: '3px'
              }
            }, row.status === 1 ? '开启' : '关闭')
          }
        },
        {
          title: '重复',
          key: 'execute_type',
          width: 200,
          render: (h, { row }) => {
            if (row.execute_type === 1) {
              return h('span', {}, '仅一次')
            } else if (row.execute_type === 2) {
              return h('span', {}, '每天')
            } else if (row.execute_type === 3) {
              let dislable = ''
              let dot = 0
              if ((row.day_of_week & 0x02) === 0x02) {
                dislable += '星期一'
                dot = 1
              }
              if ((row.day_of_week & 0x04) === 0x04) {
                if (dot === 1) {
                  dislable += ','
                }
                dislable += '星期二'
              }
              if ((row.day_of_week & 0x08) === 0x08) {
                if (dot === 1) {
                  dislable += ','
                }
                dislable += '星期三'
              }
              if ((row.day_of_week & 0x10) === 0x10) {
                if (dot === 1) {
                  dislable += ','
                }
                dislable += '星期四'
              }
              if ((row.day_of_week & 0x20) === 0x20) {
                if (dot === 1) {
                  dislable += ','
                }
                dislable += '星期五'
              }
              if ((row.day_of_week & 0x40) === 0x40) {
                if (dot === 1) {
                  dislable += ','
                }
                dislable += '星期六'
              }
              if ((row.day_of_week & 0x80) === 0x80) {
                if (dot === 1) {
                  dislable += ','
                }
                dislable += '星期日'
              }
              return h('span', {}, dislable)
            }
          }
        },
        {
          title: '限定日期',
          key: 'daterange',
          width: 200,
          render: (h, { row }) => {
            let dislable = ''
            if (row.start_date == null && row.end_date == null) {
              dislable = '无限制'
            } else {
              dislable = dayjs(row.start_date).local().format('YYYY/MM/DD') + ' - ' + dayjs(row.end_date).local().format('YYYY/MM/DD')
            }
            return h('span', {}, dislable)
          }
        },
        {
          title: '关联',
          key: 'associate_type',
          render: (h, { row }) => {
            let dislable = ''
            if (row.associate_type === 1) {
              dislable = '设备'
            } else {
              dislable = '设备组'
            }
            return h('span', {}, dislable)
          }
        },
        { title: '设备/设备组',
          key: 'devices',
          width: 200,
          render: (h, params) => {
            if (params.row.associate_type === 1) {
              const devices = params.row.devices
              return h('div', [
                h('div', devices.slice(0, 3).map(device =>
                  h('div', { style: { color: '#666' } },
                    device.mac.length > 0 ? `${device.device_name} (${device.mac})` : `${device.device_name}`
                  )
                )),
                devices.length > 3 && h('a', { on: { click: () => this.showDeviceDetail(devices) } }, '查看全部...')
              ])
            }
            if (params.row.associate_type === 2) {
              const device_groups = params.row.device_groups
              return h('div', [
                h('div', device_groups.slice(0, 3).map(device_group =>
                  h('div', { style: { color: '#666' } },
                    `${device_group.device_group_name}`
                  )
                )),
                device_groups.length > 3 && h('a', { on: { click: () => this.showDeviceGroupDetail(device_groups) } }, '查看全部...')
              ])
            }
          }
        },
        { title: '策略最后执行时间',
          key: 'group'
        },
        {
          title: '操作',
          width: 300,
          align: 'center',
          key: 'customhandle',
          options: ['edit', 'permission', 'delete']
        }
      ],
      ruleValidate: {
        policy_name: [
          { required: true, message: '策略名称不能为空', trigger: 'blur' }
        ],
        action_type: [
          { required: true, message: '执行动作不能为空', trigger: 'blur' }
        ],
        execute_time: [
          { required: true, message: '执行时间不能为空', trigger: 'blur' }
        ]
      },
      formItem: {
        id: 0,
        policy_name: '',
        action_type: 'power_on',
        execute_time: '',
        status: false,
        execute_type: 1,
        associate_type: 1,
        day_of_week: 1,
        week: ['Monday'],
        start_date: null,
        end_date: null,
        datetime: []
      },
      actionTypes: [
        {
          value: 'power_on',
          label: '开机'
        },
        {
          value: 'power_off',
          label: '关机'
        }
      ],
      isAdd: true,
      modalAdd: false,

      modalTransfer: '',
      detailDeviceVisible: false,
      detailDevices: [],
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
        { title: 'Mac地址', key: 'mac' },
        { title: 'IP地址', key: 'ip' }
      ],
      selectDeviceTableData: [],
      transferDeviceIds: [],
      detailDeviceGroupVisible: false,
      detailDeviceGroups: [],
      deviceGroupColumns: [
        {
          type: 'selection',
          width: 50,
          align: 'center'
        },
        { title: 'ID', width: 80, key: 'id', sortable: true },
        { title: '设备组名称', key: 'device_group_name', sortable: true }
      ],
      selectDeviceGroupTableData: [],
      transferDeviceGroupIds: []
    }
  },
  computed: {
    access () {
      return this.$store.state.user.access
    },
    canDelete () {
      return hasOneOf(['POLICY_DEL'], this.access)
    },
    canEdit () {
      return hasOneOf(['POLICY_EDIT'], this.access)
    },
    canAdd () {
      return hasOneOf(['POLICY_ADD'], this.access)
    }
  },
  methods: {
    onTableData (pageParams, callback) {
      getPolicyList(pageParams).then(res => {
        const today = dayjs().format('YYYY-MM-DD')
        res.list.forEach(element => {
          const fullUTCStr = `${today}T${element.execute_time}Z`
          element.execute_time = dayjs(fullUTCStr).local().format('HH:mm:ss')
          // 转换为本地时间显示
        })
        callback(res.list, res.total)
      })
    },
    onDelete (ids) {
      if (ids.length === 0) {
        this.$Message.error('请选择要删除的策略!')
      } else {
        deletePolicy({ ids: ids }).then(res => {
          this.$Message.success('删除成功!')
          this.$refs.tables.getTableData()
        })
      }
    },
    onAdd () {
      this.isAdd = true
      this.modalTransfer = 'adddevice'
      this.selectDeviceTableData = []
      this.selectDeviceGroupTableData = []
      this.$refs.formValidate.resetFields()
      this.modalAdd = true
    },
    OnSave () {
      this.$refs.formValidate.validate((valid) => {
        this.formItem.start_date = null
        this.formItem.end_date = null
        if (this.formItem.datetime.length === 2) {
          if (this.formItem.datetime[0].length !== 0) {
            this.formItem.start_date = this.formItem.datetime[0]
          }
          if (this.formItem.datetime[1].length !== 0) {
            this.formItem.end_date = this.formItem.datetime[1]
          }
        }

        this.formItem.day_of_week = 0
        if (this.formItem.week.indexOf('Monday') !== -1) {
          this.formItem.day_of_week = this.formItem.day_of_week | 0x02
        }
        if (this.formItem.week.indexOf('Tuesday') !== -1) {
          this.formItem.day_of_week = this.formItem.day_of_week | 0x04
        }
        if (this.formItem.week.indexOf('Wednesday') !== -1) {
          this.formItem.day_of_week = this.formItem.day_of_week | 0x08
        }
        if (this.formItem.week.indexOf('Thursday') !== -1) {
          this.formItem.day_of_week = this.formItem.day_of_week | 0x10
        }
        if (this.formItem.week.indexOf('Friday') !== -1) {
          this.formItem.day_of_week = this.formItem.day_of_week | 0x20
        }
        if (this.formItem.week.indexOf('Saturday') !== -1) {
          this.formItem.day_of_week = this.formItem.day_of_week | 0x40
        }
        if (this.formItem.week.indexOf('Sunday') !== -1) {
          this.formItem.day_of_week = this.formItem.day_of_week | 0x80
        }
        // 保存为UTC时间
        const today = dayjs().format('YYYY-MM-DD')
        const fullUTCStr = `${today} ${this.formItem.execute_time}`
        let execute_time = dayjs(fullUTCStr).utc().format('HH:mm:ss')
        let reqparams = {
          'policy_name': this.formItem.policy_name,
          'status': this.formItem.status ? 1 : 0,
          'action_type': this.formItem.action_type,
          'associate_type': this.formItem.associate_type,
          'execute_time': execute_time,
          'execute_type': this.formItem.execute_type,
          'day_of_week': this.formItem.day_of_week,
          'start_date': this.formItem.start_date,
          'end_date': this.formItem.end_date,
          'device_ids': this.transferDeviceIds,
          'device_group_ids': this.transferDeviceGroupIds
        }
        if (valid) {
          if (this.isAdd) {
            createPolicy(reqparams).then(res => {
              this.$Message.success('提交成功!')
              this.$refs.formValidate.resetFields()
              this.modalAdd = false
              this.$refs.tables.getTableData()
            }).catch(() => {
            })
          } else {
            reqparams.id = this.formItem.id
            updatePolicy(reqparams).then(res => {
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
      this.deviceTransfer = false
      this.deviceGroupTransfer = false
    },
    OnRowEdit (row) {
      this.$refs.formValidate.resetFields()
      if (row.associate_type === 1) {
        this.modalTransfer = 'editdevice' + row.id
        this.selectDeviceGroupTableData = []
        this.selectDeviceTableData = JSON.parse(JSON.stringify(row.devices))
      } else {
        this.modalTransfer = 'editdeviceGroup' + row.id
        this.selectDeviceGroupTableData = JSON.parse(JSON.stringify(row.device_groups))
        this.selectDeviceTableData = []
      }
      this.formItem.id = row.id
      this.formItem.policy_name = row.policy_name
      if (row.status === 1) {
        this.formItem.status = true
      } else {
        this.formItem.status = false
      }
      this.formItem.day_of_week = row.day_of_week
      this.formItem.week = []
      if ((this.formItem.day_of_week & 0x02) === 0x02) {
        this.formItem.week.push('Monday')
      }
      if ((this.formItem.day_of_week & 0x04) === 0x04) {
        this.formItem.week.push('Tuesday')
      }
      if ((this.formItem.day_of_week & 0x08) === 0x08) {
        this.formItem.week.push('Wednesday')
      }
      if ((this.formItem.day_of_week & 0x10) === 0x10) {
        this.formItem.week.push('Thursday')
      }
      if ((this.formItem.day_of_week & 0x20) === 0x20) {
        this.formItem.week.push('Friday')
      }
      if ((this.formItem.day_of_week & 0x40) === 0x40) {
        this.formItem.week.push('Saturday')
      }
      if ((this.formItem.day_of_week & 0x80) === 0x80) {
        this.formItem.week.push('Sunday')
      }
      this.formItem.action_type = row.action_type
      this.formItem.associate_type = row.associate_type
      this.formItem.execute_time = row.execute_time
      this.formItem.execute_type = row.execute_type
      this.formItem.start_date = row.start_date
      this.formItem.end_date = row.end_date
      this.formItem.datetime = []
      if (row.start_date !== null) {
        this.formItem.datetime.push(dayjs(row.start_date).local().format())
      }
      if (row.end_date !== null) {
        this.formItem.datetime.push(dayjs(row.end_date).local().format())
      }
      // this.transferDeviceIds = row.associate_type
      // this.transferDeviceGroupIds = row.associate_type
      this.isAdd = false
      this.modalAdd = true
    },
    showDeviceDetail (devices) {
      this.detailDevices = devices
      this.detailDeviceVisible = true
    },
    showDeviceGroupDetail (deviceGroups) {
      this.detailDeviceGroups = deviceGroups
      this.detailDeviceGroupVisible = true
    },
    handleExecuteTimeChange (time) {
      this.formItem.execute_time = time
    },
    onSelectDeviceTableData (callback) {
      callback(this.selectDeviceTableData)
    },
    onTransferDeviceData (transferIds) {
      this.transferDeviceIds = transferIds
    },
    onTransferDeviceTableData (pageParams, callback) {
      getDeviceList(pageParams).then(res => {
        callback(res.list, res.total)
      })
    },
    onSelectDeviceGroupTableData (callback) {
      callback(this.selectDeviceGroupTableData)
    },
    onTransferDeviceGroupData (transferIds) {
      this.transferDeviceGroupIds = transferIds
    },
    onTransferDeviceGroupTableData (pageParams, callback) {
      getDeviceGroup(pageParams).then(res => {
        callback(res.list, res.total)
      })
    },
    onRadioChange (value) {
      if (value === 1) {
        this.modalTransfer = 'adddevice'
      } else {
        this.modalTransfer = 'adddevicegroup'
      }
    }
  },
  mounted () {
    dayjs.extend(dayjsPluginUTC)
  }
}
</script>

<style>

</style>
