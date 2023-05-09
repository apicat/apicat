<template>
  <aside :class="ns.b()">
    <AcLogo pure />
    <main class="flex-1 w-full mt-30px">
      <router-link to="/projects" class="flex flex-col items-center hover:bg-gray-110 pt-13px pb-18px" active-class="font-500 bg-gray-110">
        <Iconfont :size="28" icon="ac-xiangmu" />
        <p>项目</p>
      </router-link>
    </main>

    <el-popover placement="right-start" trigger="hover">
      <template #reference>
        <footer>
          <AcIcon :size="35">
            <ac-icon-ph:user-circle-light />
          </AcIcon>
        </footer>
      </template>
      <h1 :title="userInfo.username" class="font-bold truncate px-8px pb-16px">{{ userInfo.username }}</h1>
      <PopperMenu :menus="userSettingMenus" size="small" class="clear-popover-space" />
    </el-popover>
  </aside>

  <UserSettingModal ref="userSettingModalRef" />
</template>
<script setup lang="tsx">
import { useNamespace } from '@/hooks'
import AcIconLogout from '~icons/mdi/logout'
import AcIconUserSetting from '~icons/ph/user-bold'
import { useUserStore } from '@/store/user'
import { storeToRefs } from 'pinia'
import UserSettingModal from '@/views/user/UserSettingModal.vue'

const ns = useNamespace('aside-nav')
const userStore = useUserStore()
const { userInfo } = storeToRefs(userStore)
const userSettingModalRef = ref<InstanceType<typeof UserSettingModal>>()

const handleLogout = () => userStore.logout()
const onUserSettingIconClick = () => userSettingModalRef.value?.show()

const userSettingMenus = [
  { elIcon: markRaw(AcIconUserSetting), text: '个人设置', size: 18, onClick: onUserSettingIconClick },
  { elIcon: markRaw(AcIconLogout), text: '退出登录', size: 18, onClick: handleLogout },
]
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;
@include b(aside-nav) {
  @apply h-full w-80px  flex-col flex-y-center py-20px;
  box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.1);
}
</style>
