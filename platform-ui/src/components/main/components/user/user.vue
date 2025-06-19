<template>
  <div class="user-avator-dropdown">
    <Dropdown @on-click="handleClick">
      <div class="user-name-display">
        <span>{{ userName }}</span>
        <Icon :size="18" type="md-arrow-dropdown"></Icon>
      </div>
      <DropdownMenu slot="list">
        <DropdownItem name="logout">退出登录</DropdownItem>
      </DropdownMenu>
    </Dropdown>
  </div>
</template>

<script>
import './user.less'
import { mapActions, mapState } from 'vuex'
export default {
  name: 'User',
  computed: {
    ...mapState({
      userNameFromStore: state => state.user.userName
    }),
    userName () {
      return this.userNameFromStore || '默认用户名'
    }
  },
  methods: {
    ...mapActions([
      'handleLogOut'
    ]),
    logout () {
      this.handleLogOut().then(() => {
        this.$router.push({
          name: 'login'
        })
      })
    },
    message () {
      this.$router.push({
        name: 'message_page'
      })
    },
    handleClick (name) {
      switch (name) {
        case 'logout': this.logout()
          break
      }
    }
  }
}
</script>
