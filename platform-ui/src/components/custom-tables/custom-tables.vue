<template>
    <div>
      <Card>
        <tables ref="tables" pageable editable searchable search-place="top" v-model="tableData"
        :addbtnable="canAdd" @on-delete="onDelete"
        :delbtnable="canDelete" @on-add="onAdd"
        :columns="insideColumns"
        @on-selection-change="onSelectionChange"
        :total="total" :pageSize="pageParams.pageSize" :pageIndex="pageParams.pageIndex" @on-pageChange="onPageChange"/>
      </Card>
      <PermissionForm v-model="modalAccess" :resourceId="curResourceId"
      :resourceType="curResourceType" :resourceName="curResourceName"
      :resourcePermissions="resourcePermissions"></PermissionForm>
    </div>
</template>
<script>
import Tables from '_c/tables'
import PermissionForm from '_c/permission-form'
import handleBtns from './handle-btns'
export default {
  name: 'CustomTables',
  components: {
    Tables,
    PermissionForm
  },
  props: {
    value: {
      type: Array,
      default () {
        return []
      }
    },
    columns: {
      type: Array,
      default () {
        return []
      }
    },
    size: String,
    width: {
      type: [Number, String]
    },
    height: {
      type: [Number, String]
    },
    stripe: {
      type: Boolean,
      default: false
    },
    border: {
      type: Boolean,
      default: false
    },
    showHeader: {
      type: Boolean,
      default: true
    },
    highlightRow: {
      type: Boolean,
      default: false
    },
    rowClassName: {
      type: Function,
      default () {
        return ''
      }
    },
    context: {
      type: Object
    },
    noDataText: {
      type: String
    },
    noFilteredDataText: {
      type: String
    },
    disabledHover: {
      type: Boolean
    },
    loading: {
      type: Boolean,
      default: false
    },
    /**
     * @description 全局设置是否可编辑
     */
    editable: {
      type: Boolean,
      default: false
    },
    /**
     * @description 是否可搜索
     */
    searchable: {
      type: Boolean,
      default: false
    },
    /**
     * @description 是否有添加按钮
     */
    addbtnable: {
      type: Boolean,
      default: false
    },
    /**
     * @description 是否有删除按钮
     */
    delbtnable: {
      type: Boolean,
      default: false
    },
    /**
     * @description 搜索控件所在位置，'top' / 'bottom'
     */
    searchPlace: {
      type: String,
      default: 'top'
    },
    canAdd: {
      type: Boolean,
      default: true
    },
    canEdit: {
      type: Boolean,
      default: true
    },
    canDelete: {
      type: Boolean,
      default: true
    },
    curResourceType: {
      type: String
    }
  },
  computed: {
  },
  watch: {
    columns (columns) {
      this.handleColumns(columns)
    }
  },
  data () {
    return {
      tableData: [],
      curResourceId: 0,
      curResourceName: '',
      resourcePermissions: [],
      modalAccess: false,
      insideColumns: [],
      selectids: [],
      pageParams: {
        pageSize: 10, // 每页显示的条数
        pageIndex: 1 // 当前页码
      },
      total: 0
    }
  },
  methods: {
    surportHandle (item) {
      let options = item.options || []
      let insideBtns = []
      options.forEach(item => {
        if (handleBtns[item]) insideBtns.push(handleBtns[item])
      })
      let btns = item.button ? [].concat(insideBtns, item.button) : insideBtns
      item.render = (h, params) => {
        params.tableData = this.value
        return h('div', btns.map(item => item(h, params, this)))
      }
      return item
    },
    handleColumns (columns) {
      this.insideColumns = columns.map((item, index) => {
        let res = item
        if (res.key === 'customhandle') {
          res = this.surportHandle(res)
        }
        return res
      })
    },
    accessPermissionCancel () {
      this.modalAccess = false
    },
    accessPermissionClick (params) {
      this.curResourceId = params.row.id
      this.curResourceName = params.row.device_name
      this.resourcePermissions = params.row.permissions
      this.modalAccess = true
    },
    onSelectionChange (selection) {
      this.selectids = []
      selection.forEach(element => {
        this.selectids.push(element.id)
      })
    },
    onRowEdit (row) {
      this.$emit('on-row-edit', row)
    },
    onAdd () {
      this.$emit('on-add')
    },
    onDelete () {
      this.$emit('on-delete', this.selectids)
    },
    onRowDelete (row) {
      this.$emit('on-delete', [row.id])
    },
    onRowPowerOn (row) {
      this.$emit('on-poweron', [row.id])
    },
    onRowPowerOff (row) {
      this.$emit('on-poweroff', [row.id])
    },
    onPageChange (searchForm) {
      this.pageParams = searchForm
      this.getTableData()
    },
    getTableData () {
      this.$emit('on-tabledata', this.pageParams, (tableData, total) => {
        this.tableData = tableData
        this.total = total
        if (this.tableData.length === 0 && this.total > 0 && this.pageParams.pageIndex > 1) {
          this.pageParams.pageIndex = this.pageParams.pageIndex - 1
          this.getTableData()
        }
      })
    }
  },
  mounted () {
    this.getTableData()
    this.handleColumns(this.columns)
  }
}
</script>
