<script setup lang="tsx">
import { useNamespace } from '@/hooks'
import AcIconLogout from '~icons/mdi/logout'
import AcIconUserSetting from '~icons/ph/user-bold'
import { useUserStore } from '@/store/user'
import { storeToRefs } from 'pinia'
import UserSettingModal from '@/views/user/UserSettingModal.vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const ns = useNamespace('aside-nav')
const userStore = useUserStore()
const { userInfo } = storeToRefs(userStore)
const userSettingModalRef = ref<InstanceType<typeof UserSettingModal>>()

const handleLogout = () => userStore.logout()
const onUserSettingIconClick = () => userSettingModalRef.value?.show()

const userSettingMenus = [
  { elIcon: markRaw(AcIconUserSetting), text: t('app.user.nav.userSetting'), size: 18, onClick: onUserSettingIconClick },
  { elIcon: markRaw(AcIconLogout), text: t('app.user.nav.logout'), size: 18, onClick: handleLogout },
]

const navMenus: Array<{ name: string; path: string; icon: string }> = [
  { name: t('app.project.title'), path: '/projects', icon: 'ac-xiangmu' },
  { name: t('app.iteration.title'), path: '/iterations', icon: 'ac-diedai' },
  { name: t('app.member.title'), path: '/members', icon: 'ac-members' },
]
</script>

<template>
  <aside :class="ns.b()">
    <AcLogo pure />
    <main class="flex-1 w-full mt-30px px-5px">
      <router-link
        v-for="menu in navMenus"
        :key="menu.path"
        :to="menu.path"
        class="flex flex-col items-center rounded hover:bg-gray-110 pt-10px pb-11px mb-10px"
        active-class="font-500 bg-gray-110"
      >
        <Iconfont :size="28" :icon="menu.icon" />
        <p>{{ menu.name }}</p>
      </router-link>
    </main>

    <el-popover placement="right-start" trigger="hover" width="auto">
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

<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;
@include b(aside-nav) {
  @apply h-full w-80px  flex-col flex-y-center py-20px;
  box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.1);
}
</style>
