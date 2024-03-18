<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import MemberInfo from './MemberInfo.vue'
import { useTablev2 } from '@/hooks/useTable'
import { apiEditMember, apiGetTeamMembers, apiRemoveTeamMember } from '@/api/team'
import { usePopover } from '@/hooks/usePopover'
import { useUserStore } from '@/store/user'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { useTeamStore } from '@/store/team'
import { EditableRole, Role, Status } from '@/commons/constant'
import Page30pLayout from '@/layouts/Page30pLayout.vue'
import EmptyAvatar from '@/components/EmptyAvatar.vue'

const teamStore = useTeamStore()
const { currentID, isOwner, isAdmin } = storeToRefs(teamStore)

const { t } = useI18n()
const buttonRefMap: Record<number, any> = {}
const { userInfo } = useUserStore()
const currentChangeUser = ref<TeamAPI.TeamMember | null>()

const { isShow: isShowRoleDropdownMenu, popoverRefEl, showPopover, hidePopover } = usePopover({})

const { total, currentPage, pageSize, data, isLoading, refreshData } = useTablev2(apiGetTeamMembers, {
  addonArgs: [currentID.value],
})

const columns: any = [
  {
    slot: 'name',
  },
  {
    slot: 'email',
  },
  {
    slot: 'role',
  },
  {
    slot: 'status',
  },
]

function showRoleDropdownMenu(user: TeamAPI.TeamMember) {
  showPopover(buttonRefMap[user.id])
  currentChangeUser.value = user
}

function setButtonRef(el: any, user: TeamAPI.TeamMember) {
  buttonRefMap[user.id] = el
}

function handlerRemove(user: TeamAPI.TeamMember) {
  AsyncMsgBox({
    confirmButtonClass: 'red',
    confirmButtonText: t('app.team.member.remove.btn'),
    cancelButtonText: t('app.common.cancel'),
    title: t('app.team.member.remove.poptitle'),
    content: t('app.team.member.remove.poptip'),
    onOk: async () => {
      await apiRemoveTeamMember(user.teamID, user.id)
      await refreshData()
    },
  })
}

interface RoleItem {
  text: string
  status: boolean
  onClick: () => void
}
const roleMenu = computed(() => {
  if (!currentChangeUser.value)
    return []
  // click event
  const clickFunc = (role: TeamAPI.EditableRole): (() => void) => {
    return async () => {
      await apiEditMember(currentID.value, currentChangeUser.value!.id, {
        role,
      })
      hidePopover()
      refreshData()
    }
  }

  // list
  const list: RoleItem[] = []
  for (const key in EditableRole) {
    const role: any | EditableRole = (EditableRole as any)[key as any]
    const val: RoleItem = {
      text: t(`app.team.member.roles.${role}`),
      onClick: clickFunc(role),
      status: currentChangeUser.value!.role === role,
    }
    list.push(val)
  }
  return list
})
function isSelf(row: TeamAPI.TeamMember) {
  return row.user.email === userInfo.email
}
async function handleDisableOrEnable(row: TeamAPI.TeamMember) {
  const base = `app.team.member.${row.status === Status.Active ? 'disable' : 'enable'}`
  AsyncMsgBox({
    confirmButtonText: t(`${base}.btn`),
    cancelButtonText: t(`${base}.cancel`),
    title: t(`${base}.poptitle`),
    content: t(`${base}.poptip`),
    onOk: async () => {
      await apiEditMember(currentID.value, row.id, {
        status: row.status === Status.Active ? Status.Deactive : Status.Active,
      })
      refreshData()
    },
  })
}

function roleAvailable(row: TeamAPI.TeamMember): boolean {
  if (isSelf(row))
    return false
  else if (isOwner.value)
    return true
  else return false
}
function operationAvailable(row: TeamAPI.TeamMember): boolean {
  if (isSelf(row))
    return false
  else if (isOwner.value)
    return true
  else if (isAdmin.value)
    return row.role === Role.Member
  else return false
}
</script>

<template>
  <Page30pLayout>
    <h1>{{ $t('app.member.form.title') }}</h1>
    <AcSimpleTable
      v-model:page="currentPage"
      class="mt-40px"
      is-show-pager
      round-border
      row-class-name="memberRow"
      :border="false"
      :page-size="pageSize"
      :columns="columns"
      :table-data="data"
      :loading="isLoading"
      :total="total"
    >
      <template #name>
        <el-table-column show-overflow-tooltip width="200" :label="$t('app.team.member.table.name')">
          <template #default="{ row }">
            <div class="row">
              <div class="w-full left">
                <MemberInfo v-if="row.user.avatar" :user="row.user">
                  <img class="avatar" :src="row.user.avatar">
                </MemberInfo>
                <EmptyAvatar v-else />
                <span class="ml-2 truncate">
                  {{ isSelf(row) ? `${row.user.name}(${t('app.common.self')})` : row.user.name }}
                </span>
              </div>
            </div>
          </template>
        </el-table-column>
      </template>

      <template #email>
        <el-table-column show-overflow-tooltip :label="$t('app.team.member.table.email')">
          <template #default="{ row }">
            <span>{{ row.user.email }}</span>
          </template>
        </el-table-column>
      </template>

      <template #role>
        <el-table-column width="100" :label="$t('app.team.member.table.role')">
          <template #default="{ row }">
            <!-- v-if="true" -->
            <div
              v-if="roleAvailable(row)"
              :ref="(el) => setButtonRef(el, row)"
              class="inline-flex items-center w-full cursor-pointer row role-item"
              @click="showRoleDropdownMenu(row)"
            >
              <div class="left">
                <span>
                  <!-- {{ (Role as any)[row.role] }} -->
                  {{ $t(`app.team.member.roles.${[row.role]}`) }}
                </span>
              </div>
              <div class="right">
                <div class="arrow-icon">
                  <el-icon class="m-4px" size="12px">
                    <ac-icon-ep-arrow-down />
                  </el-icon>
                </div>
              </div>
            </div>
            <div v-else>
              <span>
                {{ $t(`app.team.member.roles.${[row.role]}`) }}
              </span>
            </div>
          </template>
        </el-table-column>
      </template>

      <template #status>
        <el-table-column :label="$t('app.team.member.table.status')" width="85" align="center">
          <template #default="{ row }">
            <p v-if="row.status" :color="row.status === Status.Active ? undefined : 'orange'">
              <!-- {{ (Status as any)[row.status as any] }} -->
              {{ $t(`app.team.member.status.${[row.status]}`) }}
            </p>
          </template>
        </el-table-column>
      </template>

      <template v-if="isOwner || isAdmin" #operation>
        <el-table-column width="140" align="center">
          <template #default="{ row }">
            <div v-if="operationAvailable(row)" class="opr-item">
              <!-- disable or enable -->
              <el-button link type="default" @click="handleDisableOrEnable(row)">
                {{ $t(`app.team.member.${row.status === Status.Active ? 'disable' : 'enable'}.btn`) }}
              </el-button>

              <!-- remove -->
              <el-button link type="danger" @click="handlerRemove(row)">
                {{ $t('app.team.member.remove.btn') }}
              </el-button>
            </div>
          </template>
        </el-table-column>
      </template>
    </AcSimpleTable>

    <el-popover
      :visible="isShowRoleDropdownMenu"
      :virtual-ref="popoverRefEl"
      :show-arrow="false"
      transition="fade-fast"
      trigger="click"
      destroy-on-close
      virtual-triggering
    >
      <PopperMenu :active-menu-key="true" row-key="status" :menus="roleMenu" size="small" class="clear-popover-space" />
    </el-popover>
  </Page30pLayout>
</template>

<style lang="scss" scoped>
.avatar {
  border-radius: 50%;
  width: 30px;
}

.row {
  margin: 0;
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
  flex-grow: 1;
}
.right {
  justify-content: flex-end;
}

.memberRow .role-item .arrow-icon {
  opacity: 0;
  transition: all 0.03s;
  transform: translateY(2px);
}
.memberRow:hover .role-item .arrow-icon {
  opacity: 1;
}

.memberRow .opr-item {
  transition: all 0.03s;
  opacity: 0;
}
.memberRow:hover .opr-item {
  opacity: 1;
}
</style>
