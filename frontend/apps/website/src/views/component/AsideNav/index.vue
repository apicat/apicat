<script setup lang="tsx">
import { useI18n } from 'vue-i18n'
import { useNamespace } from '@apicat/hooks'
import { storeToRefs } from 'pinia'
import UserPopover from '@/views/component/AsideNav/UserPopover.vue'
import { useTeamStore } from '@/store/team'
import { SYSTEM_SETTING_PATH, TEAM_PATH } from '@/router'
import { useUserStore } from '@/store/user'

const { t } = useI18n()
const { hasTeam } = storeToRefs(useTeamStore())
const ns = useNamespace('aside-nav')
const { isAdmin } = storeToRefs(useUserStore())
const router = useRouter()

const navMenus: Array<{ name: string; path: string; icon: string }> = [
  { name: 'app.project.title', path: '/projects', icon: 'ac-xiangmu' },
  { name: 'app.iteration.title', path: '/iterations', icon: 'ac-diedai' },
  { name: 'app.team.title', path: TEAM_PATH, icon: 'ac-members' },
]

const sysActiveClass = computed(() => {
  if ((router.currentRoute.value.name as string || '').includes('system'))
    return 'font-500 bg-gray-110'

  return ''
})
</script>

<template>
  <aside :class="ns.b()">
    <AcLogo pure />
    <main class="flex-1 w-full mt-30px px-5px text-gray-title">
      <template v-if="hasTeam">
        <router-link
          v-for="menu in navMenus" :key="menu.path" :to="menu.path"
          class="flex flex-col items-center rounded hover:bg-gray-110 pt-10px pb-11px mb-10px"
          active-class="font-500 bg-gray-110"
        >
          <Iconfont :size="28" :icon="menu.icon" />
          <p>{{ t(menu.name) }}</p>
        </router-link>
      </template>

      <router-link
        v-if="isAdmin"
        :key="SYSTEM_SETTING_PATH"
        :to="SYSTEM_SETTING_PATH"
        class="flex flex-col items-center rounded hover:bg-gray-110 pt-10px pb-11px mb-10px"
        :class="sysActiveClass"
      >
        <Iconfont :size="28" icon="ac-setting" />
        <p>{{ t('app.system.title') }}</p>
      </router-link>
    </main>

    <UserPopover />
  </aside>
</template>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;

@include b(aside-nav) {
  @apply h-full w-80px flex-col flex-y-center py-20px;
  box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.1);
}
</style>
