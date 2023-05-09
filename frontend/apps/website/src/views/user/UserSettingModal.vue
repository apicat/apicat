<template>
  <el-dialog
    v-model="dialogVisible"
    append-to-body
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    class="fullscree hide-header"
    destroy-on-close
    align-center
    width="50%"
  >
    <ModalLayout style="height: 40vh" :slide-width="200">
      <template #nav>
        <p class="text-16px text-gray-950 font-500">{{ $t('app.common.setting') }}</p>
        <ul class="mt-20px">
          <li
            v-for="menu in menus"
            :key="menu.title"
            @click="onMenuTabClick(menu)"
            class="cursor-pointer py-10px flex-y-center"
            :class="{ 'text-blue-primary': activeTab.title === menu.title }"
          >
            <el-icon :size="17" class="mr-4px">
              <component :is="menu.icon" />
            </el-icon>
            {{ menu.title }}
          </li>
        </ul>
      </template>
      <template #title>{{ activeTab.title }}</template>
      <component :is="activeTab.component" />
    </ModalLayout>
  </el-dialog>
</template>
<script setup lang="ts">
import ModalLayout from '@/layouts/ModalLayout.vue'
import UserProfiles from '@/views/user/UserSettingPages/UserProfiles.vue'
import ModifyPassword from '@/views/user/UserSettingPages/ModifyPassword.vue'
import AcIconLock from '~icons/material-symbols/lock-outline'
import AcIconUserSetting from '~icons/ph/user-bold'
import { useModal } from '@/hooks'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const menus = [
  { title: t('app.user.nav.userSetting'), icon: markRaw(AcIconUserSetting), component: UserProfiles },
  { title: t('app.user.nav.modifyPassword'), icon: markRaw(AcIconLock), component: ModifyPassword },
]

const activeTab = shallowRef(menus[0])

const { dialogVisible, showModel } = useModal()

const onMenuTabClick = (menu: any) => {
  activeTab.value = menu
}

const show = () => {
  showModel()
}

defineExpose({
  show,
})
</script>
