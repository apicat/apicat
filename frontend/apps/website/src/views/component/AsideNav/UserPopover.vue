<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { Icon } from '@iconify/vue'
import AcIconLogout from '~icons/mdi/logout'
import { useUserStore } from '@/store/user'
import { TEAM_CREATE_PATH, USER_SETTING_NAME } from '@/router'
import { useTeamStore } from '@/store/team'
import EmptyAvatar from '@/components/EmptyAvatar.vue'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const teamStore = useTeamStore()
const { userInfo } = storeToRefs(userStore)
const { teams, currentTeam } = storeToRefs(teamStore)

const userSettingMenus = computed(() => [
  {
    iconify: 'uil:setting',
    text: t('app.user.nav.userSetting'),
    size: 18,
    onClick: () => router.push({ name: USER_SETTING_NAME }),
  },
  {
    elIcon: markRaw(AcIconLogout),
    text: t('app.user.nav.logout'),
    size: 18,
    onClick: () => userStore.logout(),
  },
])

const userTeamListToMenu = computed(() => teams.value.map(({ id, name: text }) => ({ text, id })))

async function onSwitchTeam(team: TeamAPI.Team) {
  await teamStore.switchTeam(team.id)
}
</script>

<template>
  <el-popover :show-arrow="false" placement="right-start" trigger="hover" transition="fade-fast" :width="300">
    <template #reference>
      <footer class="text-gray-title">
        <AcIcon :size="35">
          <ac-icon-ph:user-circle-light />
        </AcIcon>
      </footer>
    </template>
    <!-- user&team -->
    <div class="row p-5px">
      <div class="w-full left text-gray-title">
        <img v-if="userInfo.avatar" width="20" :src="userInfo.avatar" style="border-radius: 50%" />
        <EmptyAvatar v-else />
        <p class="ml-2 truncate">
          {{ userInfo.name }}
        </p>
      </div>
    </div>
    <h1 class="mt-2 font-500 truncate px-8px pb-16px">
      {{ $t('app.user.nav.teams') }}
    </h1>
    <PopperMenu
      style="min-width: 250px; padding-bottom: 17px"
      :menus="userTeamListToMenu"
      row-key="id"
      :active-menu-key="currentTeam?.id"
      size="small"
      class="clear-popover-space"
      @menu-click="onSwitchTeam">
      <template #prefix>
        <i style="width: 18px">
          <Icon icon="ph:dot" :width="18" />
        </i>
      </template>
      <template #suffix>
        <Icon icon="ph:check" color="rgba(72,148,255)" class="mr-1" :width="18" />
      </template>
    </PopperMenu>
    <div class="mt-0">
      <ElButton link type="primary" @click="() => $router.push(TEAM_CREATE_PATH)">
        <div class="row">
          <div class="left">
            <Icon class="mr-1" icon="ion:add-outline" :width="18" />
            {{ $t('app.user.nav.create') }}
          </div>
        </div>
      </ElButton>
    </div>

    <ElDivider />
    <!-- menu -->
    <PopperMenu
      style="min-width: 200px; padding-bottom: 17px"
      :menus="userSettingMenus"
      size="small"
      class="clear-popover-space" />
  </el-popover>
</template>

<style lang="scss" scoped>
.el-divider--horizontal {
  margin: 16px 0;
}

.row {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.left,
.right {
  display: flex;
  align-items: center;
}

.left {
  /* justify-content: flex-start; */
  flex-grow: 1;
}

.right {
  justify-content: flex-end;
  /* flex-grow: 1; */
}

/* el-image */
.row .block {
  padding: 30px 0;
  text-align: center;
  border-right: solid 1px var(--el-border-color);
  display: inline-block;
  width: 49%;
  box-sizing: border-box;
  vertical-align: top;
}
.row .demonstration {
  display: block;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-bottom: 20px;
}
.row .el-image {
  width: 20px;
  height: 20px;
  border-radius: 50%;
}

.row .image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-secondary);
  font-size: 30px;
}
</style>
