<template>
    <div class="parent">
      <div>
        <Tag type="dot" size="large" color="primary">{{titles[0]}}</Tag>
        <tables ref="select_tables"  v-model="selectTableData"
          :width="600"
          size="small"
          :columns="columns"
          @on-selection-change="onSelectSelectionChange"
          />
          <div>总计:{{ selectTableData.length }}条</div>
      </div>
      <div class ="oprator">
        <Button type="primary" @click="addSelect"><Icon type="md-arrow-back" size="18"/>{{operations[0]}}</Button>
        <Button type="primary" @click="delSelect">{{operations[1]}}<Icon type="md-arrow-forward" size="18"/></Button>
      </div>
      <div>
        <Tag  type="dot" size="large" color="primary">{{titles[1]}}</Tag>
        <tables ref="candidate_tables" pageable search-place="top" v-model="candidateTableData"
          :width="650"
          size="small"
          :columns="columns"
          @on-selection-change="onCandidateSelectionChange"
          :total="total" :pageSize="pageParams.pageSize" :pageIndex="pageParams.pageIndex" @on-pageChange="onPageChange"/>
      </div>
    </div>
</template>
<script>
import Tables from '_c/tables'
export default {
  name: 'CustomTransfer',
  components: {
    Tables
  },
  props: {
    value: {
      type: String,
      default: ''
    },
    columns: {
      type: Array,
      default () {
        return []
      }
    },
    titles: {
      type: Array,
      default () {
        return ['已选', '待选']
      }
    },
    operations: {
      type: Array,
      default () {
        return ['新增', '移除']
      }
    }
  },
  computed: {
  },
  watch: {
    value (val) {
      this.total = 0
      this.pageParams.pageIndex = 1
      this.pageParams.pageSize = 10
      this.getTableData()
      this.getSelectTableData()
    },
    selectTableData () {
      this.updateCanditaeTableStatus()
      let selectids = []
      this.selectTableData.forEach(element => {
        selectids.push(element.id)
      })
      this.$emit('on-transfer', selectids)
    }
  },
  data () {
    return {
      total: 0,
      pageParams: {
        pageSize: 10,
        pageIndex: 1
      },
      selectTableData: [],
      selectSelection: [],
      candidateSelection: [],
      candidateTableData: []
    }
  },
  methods: {
    addSelect () {
      this.candidateSelection.forEach(element => {
        if (this.selectTableData.findIndex(item => item.id === element.id) === -1) {
          this.selectTableData.push(element)
        }
      })
    },
    delSelect () {
      let newSelectTableData = []
      this.selectTableData.forEach(element => {
        if (this.selectSelection.findIndex(item => item.id === element.id) === -1) {
          newSelectTableData.push(element)
        }
      })
      this.selectTableData = newSelectTableData
    },
    updateCanditaeTableStatus () {
      this.candidateTableData.forEach(el => {
        el._checked = false
        if (this.selectTableData.findIndex(item => item.id === el.id) === -1) {
          el._disabled = false
        } else {
          el._disabled = true
        }
      })
      this.candidateTableData = [...this.candidateTableData]
    },
    onSelectSelectionChange (selection) {
      this.selectSelection = selection
    },
    onCandidateSelectionChange (selection) {
      this.candidateSelection = selection
    },
    onPageChange (searchForm) {
      this.pageParams = searchForm
      this.getTableData()
    },
    getSelectTableData () {
      this.$emit('on-selectTableData', (tableData) => {
        this.selectTableData = tableData
      })
    },
    getTableData () {
      this.$emit('on-tabledata', this.pageParams, (tableData, total) => {
        this.candidateTableData = tableData
        this.updateCanditaeTableStatus()
        this.total = total
        if (this.candidateTableData.length === 0 && this.candidateTableData > 0 && this.pageParams.pageIndex > 1) {
          this.pageParams.pageIndex = this.pageParams.pageIndex - 1
          this.getTableData()
        }
      })
    }
  },
  mounted () {
  }
}
</script>
<style>
.parent {
  display: flex; /* 启用Flex布局 */
  justify-content: flex-start; /* 子元素从左到右排列（默认值，可省略） */
  gap: 10px; /* 可选：设置子元素之间的间距 */
}
.oprator {
  display: flex; /* 启用Flex布局 */
  flex-direction: column; /* 垂直排列按钮 */
  justify-content: center; /* 垂直居中 */
  gap: 8px; /* 按钮之间的间距（可选） */
}
</style>
