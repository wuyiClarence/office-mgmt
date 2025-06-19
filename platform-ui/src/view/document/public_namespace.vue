<style lang="less" scoped>
  @import './public_namespace.less';
</style>

<template>
  <div class="file-explorer-page">
    <div class="flex-container">
      <!-- 左侧空间菜单 -->
      <div class="space-menu-container">
        <div class="menu-header">
          <h3 class="menu-title">
            <Icon type="md-cube" size="20" />
            空间管理
          </h3>
          <Button type="primary" size="small" @click="handleNameSpaceAdd">
            <Icon type="md-add" />
            新建空间
          </Button>
        </div>

        <div class="space-list">
          <div
            v-for="space in namespaceList"
            :key="space.id"
            class="space-item"
            :class="{ active: selectedSpace && selectedSpace.id === space.id }"
            @click="selectSpace(space)"
          >
            <div class="space-main">
              <div class="space-icon">
                <Icon type="md-cube" size="24" />
              </div>
              <div class="space-info">
                <div class="space-name">{{ space.space_name }}</div>
                <div class="space-desc">{{ space.description || '暂无描述' }}</div>
              </div>
            </div>
            <div class="space-actions">
              <Icon
                type="md-more"
                size="20"
                @click.stop="showSpaceMenu($event, space)"
                class="action-btn"
              />
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="namespaceList.length === 0" class="empty-spaces">
            <Icon type="md-cube" size="48" color="#ccc" />
            <p>暂无空间</p>
            <Button type="primary" @click="handleNameSpaceAdd">创建第一个空间</Button>
          </div>
        </div>
      </div>

      <!-- 右侧内容区域 -->
      <div class="content-container">
        <!-- 头部信息 -->
        <div class="content-header">
          <div class="current-space-info" v-if="selectedSpace">
            <div class="space-title">
              <Icon type="md-cube" size="20" />
              <span>{{ selectedSpace.space_name }}</span>
            </div>
            <div class="breadcrumb-container">
              <Breadcrumb separator=">">
                <Breadcrumb-item>{{ selectedSpace.space_name }}</Breadcrumb-item>
                <Breadcrumb-item v-for="(crumb, index) in breadcrumbList" :key="index">
                  <span>{{ crumb.name }}</span>
                </Breadcrumb-item>
              </Breadcrumb>
            </div>
          </div>

          <div v-else class="no-space-selected">
            <Icon type="md-information-circle" size="20" />
            <span>请选择一个空间来查看内容</span>
          </div>
        </div>

        <!-- 工具栏 -->
        <div class="toolbar" v-if="selectedSpace">
          <Button-group>
            <Button type="primary" @click="refreshContent">
              <Icon type="ios-refresh" />刷新
            </Button>
            <Button type="primary" @click="handleCreateFolder">
              <Icon type="ios-folder" />新建文件夹
            </Button>
            <Button type="primary" @click="$refs.fileInput.click()">
              <Icon type="ios-cloud-upload" />上传文件
            </Button>
          </Button-group>

          <div class="search-box">
            <Input
              v-model="searchKeyword"
              placeholder="搜索文件或文件夹"
              @on-enter="searchFiles"
            >
              <Icon slot="prepend" type="ios-search" />
            </Input>
          </div>
        </div>

        <!-- 内容列表 -->
        <div class="content-list" v-if="selectedSpace">
          <div v-if="loading" class="loading-container">
            <div class="overlay">
              <Spin />
            </div>
          </div>

          <div v-else-if="contentList.length === 0" class="empty-container">
            <div class="empty-icon">
              <Icon type="ios-folder-outline" size="60" />
              <p class="empty-text">此文件夹为空</p>
              <Button type="primary" @click="handleCreateFolder">
                <Icon type="ios-folder" />新建文件夹
              </Button>
            </div>
          </div>

          <div v-else class="file-grid">
            <div
              v-for="item in contentList"
              :key="item.file_id"
              class="file-item"
              @click="item.is_dir ? openFolder(item) : openFile(item)"
              @contextmenu.prevent="showContextMenu($event, item)"
            >
              <input type="checkbox" v-model="selectedItems[item.file_id]" @change="handleSelectionChange(item)" @click.stop/>
              <div class="file-icon">
                <Icon
                  :type="getIconType(item)"
                  :size="40"
                  :color="item.is_dir ? '#3688FF' : '#666'"
                />
              </div>
              <div class="file-name" :title="item.file_name">{{ item.file_name }}</div>
              <div class="file-info">
                <span>{{ formatDate(item.updated_at) }}</span>
                <span v-if="!item.is_dir">{{ formatSize(item.size) }}</span>
              </div>
              <Icon type="ios-more" @click="showItemActions($event, item)" class="action-icon" />
            </div>
          </div>
        </div>

        <!-- 未选择空间时的提示 -->
        <div v-else class="no-space-placeholder">
          <div class="placeholder-content">
            <Icon type="md-cube" size="80" color="#ccc" />
            <h3>选择一个空间开始管理文件</h3>
            <p>从左侧选择一个空间，或者创建一个新的空间来开始管理您的文件。</p>
            <Button type="primary" @click="handleNameSpaceAdd">
              <Icon type="md-add" />
              创建新空间
            </Button>
          </div>
        </div>
      </div>
    </div>

    <input
      type="file"
      ref="fileInput"
      style="display: none"
      @change="handleFileSelected"
      accept="*/*"
    />
    <Modal v-model="isUploading" title="文件上传中" :footer-hide="true" :closable="false" :mask-closable="false" class-name="upload-progress-modal">
      <div style="text-align:center;">
        <Progress :percent="uploadProgress" :stroke-width="20" />
        <p>{{ uploadProgress }}%</p>
        <Button @click="cancelUpload" type="error" style="margin-top: 15px;">取消上传</Button>
      </div>
    </Modal>
    <!-- 空间操作菜单 -->
    <Menu
      ref="spaceMenu"
      theme="dark"
      :style="{ left: spaceMenuLeft + 'px', top: spaceMenuTop + 'px' }"
      v-if="spaceMenuVisible"
      @on-select="handleSpaceMenuSelect"
    >
      <MenuItem name="edit">
        <Icon type="ios-create" />编辑空间
      </MenuItem>
      <MenuItem name="delete">
        <Icon type="ios-trash" />删除空间
      </MenuItem>
    </Menu>

    <!-- 文件右键菜单 -->
    <Menu
      ref="contextMenu"
      theme="dark"
      :active-name="activeName"
      :style="{ left: contextMenuLeft + 'px', top: contextMenuTop + 'px' }"
      v-if="contextMenuVisible"
      :key="contextMenuKey"
      @on-select="handleContextMenuSelect"
    >
      <MenuItem name="open" v-if="selectedItem.type === 'folder'">
        <Icon type="ios-folder-open" />打开
      </MenuItem>
      <MenuItem name="download" v-if="selectedItem.type === 'file'">
        <Icon type="ios-cloud-download" />下载
      </MenuItem>
      <MenuItem name="rename">
        <Icon type="ios-create" />重命名
      </MenuItem>
      <MenuItem name="delete">
        <Icon type="ios-trash" />删除
      </MenuItem>
    </Menu>

    <!-- 文件操作菜单 -->
    <Menu
      ref="itemActionsMenu"
      theme="dark"
      :style="{ left: itemActionsMenuLeft + 'px', top: itemActionsMenuTop + 'px' }"
      v-if="itemActionsMenuVisible"
      @on-select="handleItemActionsMenuSelect"
    >
      <MenuItem name="open" v-if="selectedItem.type === 'folder'">
        <Icon type="ios-folder-open" />打开
      </MenuItem>
      <MenuItem name="download" v-if="selectedItem.type === 'file'">
        <Icon type="ios-cloud-download" />下载
      </MenuItem>
      <MenuItem name="rename">
        <Icon type="ios-create" />重命名
      </MenuItem>
      <MenuItem name="delete">
        <Icon type="ios-trash" />删除
      </MenuItem>
    </Menu>

    <!-- 空间创建/编辑模态框 -->
    <Modal
      :mask-closable="false"
      v-model="modalSpace"
      :title="isAddSpace ? '创建空间' : '编辑空间'"
    >
      <Form :model="formItemSpace" :label-width="80" ref="formValidateSpace" :rules="ruleValidateSpace">
        <Form-item label="空间名" prop="space_name">
          <Input v-model="formItemSpace.space_name" placeholder="请输入空间名"></Input>
        </Form-item>
        <Form-item label="描述" prop="description">
          <Input v-model="formItemSpace.description" type="textarea" :rows="3" placeholder="请输入描述"></Input>
        </Form-item>
      </Form>
      <div slot="footer">
        <Button @click.native="OnCancelNameSpace">取消</Button>
        <Button type="primary" @click.native="OnOkNameSpace">保存</Button>
      </div>
    </Modal>

    <!-- 文件夹创建/编辑模态框 -->
    <Modal
      :mask-closable="false"
      v-model="modalFolder"
      :title="isAddFolder ? '创建文件夹' : '编辑文件夹'"
    >
      <Form :model="formItemFolder" :label-width="80" ref="formValidateFolder" :rules="ruleValidateFolder">
        <Form-item label="文件夹名" prop="folder_name">
          <Input v-model="formItemFolder.folder_name" placeholder="请输入文件夹名"></Input>
        </Form-item>
      </Form>
      <div slot="footer">
        <Button @click.native="OnCancelFolder">取消</Button>
        <Button type="primary" @click.native="OnOkFolder">保存</Button>
      </div>
    </Modal>
  </div>
</template>

<script>
import axios from 'axios'
import { getNameSpace, createNameSpace, deleteNameSpace, updateNameSpace } from '@/api/namespace'
import { getFile, deleteFile, uploadFile } from '@/api/file'
import { createFolder, updateFolder, deleteFolder } from '../../api/folder'

export default {
  name: 'FileExplorer',
  data () {
    return {
      namespaceList: [],
      selectedSpace: null,
      breadcrumbList: [],
      contentList: [],
      loading: false,
      searchKeyword: '',

      currentParentId: 0,

      // 菜单相关
      spaceMenuVisible: false,
      spaceMenuLeft: 0,
      spaceMenuTop: 0,
      contextMenuVisible: false,
      contextMenuLeft: 0,
      contextMenuTop: 0,
      itemActionsMenuVisible: false,
      itemActionsMenuLeft: 0,
      itemActionsMenuTop: 0,
      contextMenuKey: 0,
      activeName: null,
      selectedItem: {},
      selectedSpaceForMenu: {},
      currentPath: '',

      fileIcons: {
        pdf: 'ios-document',
        doc: 'ios-paper',
        docx: 'ios-paper',
        xls: 'ios-stats',
        xlsx: 'ios-stats',
        ppt: 'ios-film',
        pptx: 'ios-film',
        jpg: 'ios-image',
        jpeg: 'ios-image',
        png: 'ios-image',
        gif: 'ios-image',
        txt: 'ios-document-text',
        default: 'ios-document'
      },

      selectedItems: {},

      // 空间模态框相关
      isAddSpace: false,
      modalSpace: false,
      ruleValidateSpace: {
        space_name: [
          { required: true, message: '空间名不能为空', trigger: 'blur' }
        ]
      },
      formItemSpace: {
        space_name: '',
        description: ''
      },
      isAddFolder: false,
      modalFolder: false,
      ruleValidateFolder: {
        folder_name: [
          { required: true, message: '文件夹名不能为空', trigger: 'blur' }
        ]
      },
      formItemFolder: {
        space_name: ''
      },
      uploadProgress: 0, // 上传进度百分比
      isUploading: false, // 是否正在上传
      uploadCancelToken: null // 用于取消上传的 token
    }
  },

  created () {
    this.getNameSpaceList()
  },

  methods: {
    // 获取空间列表
    getNameSpaceList () {
      getNameSpace().then(res => {
        this.namespaceList = res.list || []
        // 如果有空间且当前没有选中的空间，默认选中第一个
        if (this.namespaceList.length > 0 && !this.selectedSpace) {
          this.selectSpace(this.namespaceList[0])
        }
      })
    },

    // 选择空间
    selectSpace (space) {
      this.selectedSpace = space
      this.currentPath = ''
      this.breadcrumbList = []

      this.currentParentId = 0
      this.loadContent()
    },

    // 显示空间菜单
    showSpaceMenu (event, space) {
      this.selectedSpaceForMenu = space
      this.spaceMenuLeft = event.clientX
      this.spaceMenuTop = event.clientY
      this.spaceMenuVisible = true
      event.stopPropagation()
    },

    // 处理空间菜单选择
    handleSpaceMenuSelect (name) {
      this.spaceMenuVisible = false
      switch (name) {
        case 'edit':
          this.editSpace(this.selectedSpaceForMenu)
          break
        case 'delete':
          this.confirmDeleteSpace(this.selectedSpaceForMenu)
          break
      }
    },

    // 编辑空间
    editSpace (space) {
      this.isAddSpace = false
      this.formItemSpace = {
        id: space.id,
        space_name: space.space_name,
        description: space.description || ''
      }
      this.modalSpace = true
    },

    // 确认删除空间
    confirmDeleteSpace (space) {
      this.$Modal.confirm({
        title: '确认删除',
        content: `确定要删除空间 "${space.space_name}" 吗？此操作不可撤销。`,
        onOk: () => {
          this.deleteNameSpace([space.id])
        }
      })
    },

    // 新建空间
    handleNameSpaceAdd () {
      this.isAddSpace = true
      this.$refs.formValidateSpace.resetFields()
      this.formItemSpace = {
        space_name: '',
        description: ''
      }
      this.modalSpace = true
    },

    // 取消空间操作
    OnCancelNameSpace () {
      this.$refs.formValidateSpace.resetFields()
      this.modalSpace = false
    },

    // 保存空间
    OnOkNameSpace () {
      this.$refs.formValidateSpace.validate((valid) => {
        if (valid) {
          if (this.isAddSpace) {
            createNameSpace(this.formItemSpace).then(res => {
              this.$Message.success('创建成功!')
              this.$refs.formValidateSpace.resetFields()
              this.modalSpace = false
              this.getNameSpaceList()
            }).catch(() => {})
          } else {
            updateNameSpace(this.formItemSpace).then(res => {
              this.$Message.success('更新成功!')
              this.$refs.formValidateSpace.resetFields()
              this.modalSpace = false
              this.getNameSpaceList()
            }).catch(() => {})
          }
        } else {
          this.$Message.error('表单验证失败!')
        }
      })
    },

    // 删除空间
    deleteNameSpace (spaceIds) {
      deleteNameSpace({ space_ids: spaceIds }).then(res => {
        this.$Message.success('删除成功!')
        // 如果删除的是当前选中的空间，清空选择
        if (this.selectedSpace && spaceIds.includes(this.selectedSpace.id)) {
          this.selectedSpace = null
          this.contentList = []
        }
        this.getNameSpaceList()
      })
    },
    // 取消文件夹操作
    OnCancelFolder () {
      this.$refs.formValidateFolder.resetFields()
      this.modalFolder = false
    },

    // 保存文件夹操作
    OnOkFolder () {
      this.$refs.formValidateFolder.validate((valid) => {
        if (valid) {
          if (this.isAddFolder) {
            let req = {
              folder_name: this.formItemFolder.folder_name,
              namespace_id: this.selectedSpace.id,
              parent_id: this.currentParentId
            }
            createFolder(req).then(res => {
              this.$Message.success('创建成功!')
              this.$refs.formValidateFolder.resetFields()
              this.modalFolder = false
              this.loadContent()
            }).catch(() => {})
          } else {
            updateFolder(this.formItemFolder).then(res => {
              this.$Message.success('更新成功!')
              this.$refs.formValidateFolder.resetFields()
              this.modalFolder = false
              this.loadContent()
            }).catch(() => {})
          }
        } else {
          this.$Message.error('表单验证失败!')
        }
      })
    },
    // 加载内容
    loadContent () {
      if (!this.selectedSpace) return

      this.loading = true
      let req = {
        namespace_id: this.selectedSpace.id,
        parent_id: this.currentParentId,
        pageSize: 10, // 每页显示的条数
        pageIndex: 1 // 当前页码
      }
      getFile(req).then(res => {
        this.contentList = res.list || []
        this.loading = false
      })
    },

    // // 生成模拟内容
    // generateMockContent (path) {
    //   const depth = path.split('/').filter(p => p).length
    //   const mockFolders = [
    //     { name: '文档', type: 'folder', id: `folder-docs-${depth}`, updatedAt: '2025-05-27' },
    //     { name: '图片', type: 'folder', id: `folder-images-${depth}`, updatedAt: '2025-05-26' },
    //     { name: '视频', type: 'folder', id: `folder-videos-${depth}`, updatedAt: '2025-05-25' },
    //     { name: '项目', type: 'folder', id: `folder-projects-${depth}`, updatedAt: '2025-05-24' }
    //   ]
    //   const mockFiles = [
    //     { name: '需求文档.pdf', type: 'file', id: `file-req-${depth}`, size: 2145893, updatedAt: '2025-05-27' },
    //     { name: '会议记录.docx', type: 'file', id: `file-meeting-${depth}`, size: 102400, updatedAt: '2025-05-26' },
    //     { name: '项目计划.xlsx', type: 'file', id: `file-plan-${depth}`, size: 567890, updatedAt: '2025-05-25' },
    //     { name: '架构设计图.png', type: 'file', id: `file-arch-${depth}`, size: 3214567, updatedAt: '2025-05-24' }
    //   ]

    //   const foldersCount = Math.floor(Math.random() * 4)
    //   const filesCount = Math.floor(Math.random() * 6)
    //   return [
    //     ...mockFolders.slice(0, foldersCount),
    //     ...mockFiles.slice(0, filesCount)
    //   ]
    // },

    // 打开文件夹
    openFolder (item) {
      const path = this.currentPath ? `${this.currentPath}/${item.file_name}` : item.file_name
      this.currentPath = path
      this.currentParentId = item.file_id
      this.loadContent()
      this.breadcrumbList.push({
        name: item.file_name,
        path: path
      })
    },

    // 打开文件
    openFile (file) {
      this.$Message.info(`打开文件: ${file.name}`)
    },

    // 获取文件图标
    getIconType (item) {
      if (item.is_dir) {
        return 'ios-folder'
      } else {
        const ext = item.file_name.split('.').pop().toLowerCase()
        return this.fileIcons[ext] || this.fileIcons.default
      }
    },

    // 刷新内容
    refreshContent () {
      this.loadContent()
    },

    // 创建文件夹
    handleCreateFolder () {
      this.isAddFolder = true
      this.$refs.formValidateFolder.resetFields()
      this.formItemFolder = {
        folder_name: ''
      }
      this.modalFolder = true
    },
    handleFileSelected (event) {
      const file = event.target.files[0]
      if (!file) return

      // 开始上传逻辑
      this.startFileUpload(file)
      event.target.value = ''
    },
    // 上传文件
    startFileUpload (file) {
      const CancelToken = axios.CancelToken
      const source = CancelToken.source()

      this.isUploading = true
      this.uploadProgress = 0
      this.uploadCancelToken = source

      const formData = new FormData()
      formData.append('file', file)
      formData.append('namespace_id', this.selectedSpace.id)
      formData.append('parent_id', this.currentParentId)
      uploadFile(formData, {
        cancelToken: source.token,
        onUploadProgress: (progressEvent) => {
          const percentCompleted = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total
          )
          this.uploadProgress = percentCompleted
        }
      })
        .then((res) => {
          this.$Message.success('文件上传成功！')
          this.loadContent()
          this.isUploading = false
          this.uploadProgress = 0
          this.uploadCancelToken = null
        })
        .catch((error) => {
          this.isUploading = false
          this.uploadProgress = 0
          this.uploadCancelToken = null
          if (axios.isCancel(error)) {
            this.$Message.warning('上传已取消')
          } else {
            this.$Message.error('文件上传失败')
            console.error('Upload error:', error)
          }
        })
        .finally(() => {
          this.isUploading = false
          this.uploadProgress = 0
          this.uploadCancelToken = null
        })
    },
    cancelUpload () {
      if (this.uploadCancelToken) {
        this.uploadCancelToken.cancel('用户取消上传')
        this.isUploading = false
        this.uploadProgress = 0
        this.uploadCancelToken = null
      }
    },

    // 搜索文件
    searchFiles () {
      if (!this.searchKeyword.trim()) {
        this.refreshContent()
        return
      }
      this.loading = true
      setTimeout(() => {
        const filteredContent = this.contentList.filter(item =>
          item.name.toLowerCase().includes(this.searchKeyword.toLowerCase())
        )
        this.contentList = filteredContent
        this.loading = false
      }, 1500)
    },

    // 显示右键菜单
    showContextMenu (event, item) {
      this.selectedItem = item
      this.contextMenuLeft = event.clientX
      this.contextMenuTop = event.clientY
      this.contextMenuVisible = true
      this.activeName = ''
      event.stopPropagation()
    },

    // 处理右键菜单
    handleContextMenuSelect (name) {
      this.contextMenuVisible = false
      this.selectedItem = {}
      this.activeName = ''
      this.handleFileAction(name)
    },

    // 显示文件操作菜单
    showItemActions (event, item) {
      this.selectedItem = item
      this.itemActionsMenuLeft = event.clientX
      this.itemActionsMenuTop = event.clientY
      this.itemActionsMenuVisible = true
      event.stopPropagation()
    },

    // 处理文件操作菜单
    handleItemActionsMenuSelect (name) {
      this.itemActionsMenuVisible = false
      this.handleFileAction(name)
    },

    // 处理文件操作
    handleFileAction (action) {
      switch (action) {
        case 'open':
          this.openFolder(this.selectedItem)
          break
        case 'download':
          this.$Message.info(`开始下载: ${this.selectedItem.name}`)
          break
        case 'rename':
          this.$Modal.prompt({
            title: '重命名',
            content: '请输入新名称',
            value: this.selectedItem.name,
            onOk: (val) => {
              if (!val.trim()) {
                this.$Message.error('名称不能为空')
                return
              }
              setTimeout(() => {
                this.$Message.success(`已重命名为 "${val}"`)
                this.refreshContent()
              }, 500)
            }
          })
          break
        case 'delete':
          if (this.selectedItem.is_dir) {
            this.$Modal.confirm({
              title: '确认删除',
              content: `确定要删除 "${this.selectedItem.file_name}" 文件夹及下面的所有文件吗？`,
              onOk: () => {
                let ret = {
                  'folder_id': this.selectedItem.file_id,
                  'namespace_id': this.selectedSpace.id
                }
                deleteFolder(ret).then(res => {
                  this.$Message.success('删除成功!')
                  this.refreshContent()
                }).catch(() => {})
              }
            })
          } else {
            this.$Modal.confirm({
              title: '确认删除',
              content: `确定要删除 "${this.selectedItem.file_name}" 文件吗？`,
              onOk: () => {
                let ret = {
                  'file_id': this.selectedItem.file_id,
                  'namespace_id': this.selectedSpace.id
                }
                deleteFile(ret).then(res => {
                  this.$Message.success('删除成功!')
                  this.refreshContent()
                }).catch(() => {})
              }
            })
          }
          break
      }
    },

    // 处理选择变化
    handleSelectionChange (item) {
      console.log('Selected item:', item, 'Selected status:', this.selectedItems[item.id])
    },

    // 格式化文件大小
    formatSize (bytes) {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    },

    // 格式化日期
    formatDate (dateStr) {
      if (!dateStr) return ''
      const date = new Date(dateStr)
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hours = String(date.getHours()).padStart(2, '0')
      const minutes = String(date.getMinutes()).padStart(2, '0')
      return `${year}-${month}-${day} ${hours}:${minutes}`
    }
  },

  mounted () {
    // 点击页面其他地方关闭菜单
    document.addEventListener('click', () => {
      this.contextMenuVisible = false
      this.itemActionsMenuVisible = false
      this.spaceMenuVisible = false
    })
  },

  beforeDestroy () {
    document.removeEventListener('click', () => {
      this.contextMenuVisible = false
      this.itemActionsMenuVisible = false
      this.spaceMenuVisible = false
    })
  }
}
</script>
