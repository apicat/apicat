<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { storeToRefs } from 'pinia'
import { useNamespace } from '@apicat/hooks'
import LeftRightLayout from '@/layouts/LeftRightLayout.vue'
import { useTeamStore } from '@/store/team'
import { Role } from '@/commons/constant'

const ns = useNamespace('group-list')
const teamStore = useTeamStore()
const menus = (() => {
  const c: Record<string, any> = {
    member: {
      icon: 'uil:user',
      component: defineAsyncComponent(() => import('./pages/member/MemberListPage.vue')),
    },
  }
  if (teamStore.currentRole !== Role.Member) {
    c.invite = {
      icon: 'ion:add',
      component: defineAsyncComponent(() => import('./pages/Invite.vue')),
    }
  }
  c.setting = {
    icon: 'uil:setting',
    component: defineAsyncComponent(() => import('./pages/setting/Setting.vue')),
  }
  return c
})()
const currentPage = ref<string>(Object.keys(menus)[0])
const { currentTeam } = storeToRefs(teamStore)

function activeClass(key: string) {
  return currentPage.value === key ? 'active' : ''
}
</script>

<template>
  <LeftRightLayout main-width="auto">
    <template #left>
      <div class="flex flex-col h-full">
        <div :class="[ns.b(), ns.m('header')]">
          <div
            v-for="(menu, name) in menus"
            :key="name"
            :class="[ns.e('item'), activeClass(name)]"
            @click="currentPage = name"
          >
            <Icon color="rgba(72,148,255)" class="mr-2" :icon="menu.icon" width="18" />
            <span>{{ $t(`app.team.${name}.left_title`) }}</span>
          </div>
        </div>
      </div>
    </template>

    <component :is="menus[currentPage].component" v-if="currentTeam" />
    <div v-else>
      <h1 style="text-align: center; font-size: 20px; font-weight: 700">
        {{ $t('app.team.no_current') }}
      </h1>
    </div>
  </LeftRightLayout>
</template>
