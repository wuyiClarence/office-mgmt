<template>
    <Modal
            v-model="localValue"
            :title="title"
            :closable="false"
            :mask-closable="false"
            @on-cancel="cancel"
            :styles="{ backgroundColor: 'transparent' }"
            class-name="vertical-center-modal">
          <slot name="main-content"></slot>
          <slot></slot>
        <div v-if="showFooter" slot="footer">
            <Button type="primary" @click.native="ok">{{okText}}</Button>
            <Button @click.native="cancel">{{cancelText}}</Button>
        </div>
        <div v-else slot="footer"></div>
    </Modal>
</template>

<script>
import './index.less'
export default {
  name: 'TipModal',
  props: {
    value: {
      required: true,
      type: Boolean
    },
    title: {
      type: String,
      default: '提示信息'
    },
    showFooter: {
      type: Boolean,
      default: true
    },
    okText: {
      type: String,
      default: '确定'
    },
    cancelText: {
      type: String,
      default: '取消'
    }
  },
  data () {
    return {
      localValue: this.value
    }
  },
  watch: {
    value (newVal) {
      this.localValue = newVal
    },
    localValue (newVal) {
      this.$emit('input', newVal)
    }
  },
  methods: {
    ok () {
      this.$emit('click-ok')
      // this.localValue = false
    },
    cancel () {
      this.$emit('click-cancel')
      // this.localValue = false
    }
  }
}
</script>
