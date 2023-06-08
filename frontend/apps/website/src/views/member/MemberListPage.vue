<template>
  <div class="container flex flex-col justify-center mx-auto px-20px">
    <div class="flex justify-between border-b border-solid border-gray-lighter py-20px mt-20px text-20px">
      <p>{{ $t('app.member.form.title') }}</p>
      <el-button type="primary" @click="addMemberModal?.show()">{{ $t('app.member.tips.addMember') }}</el-button>
    </div>
    <AcSimpleTable class="mt-40px" isShowPager v-model:page="currentPage" v-model:page-size="pageSize" :columns="columns" :table-data="data" :loading="isLoading" :total="total">
      <template #accountStatus>
        <el-table-column :label="$t('app.member.form.accountStatus')" width="100" align="center">
          <template #default="{ row }">
            <el-tag disable-transitions :type="row.accountStatusType">{{ row.accountStatus }}</el-tag>
          </template>
        </el-table-column>
      </template>

      <template #auth>
        <el-table-column :label="$t('app.member.form.role')">
          <template #default="{ row }">
            <div v-if="isSuperAdmin && !row.isSelf" :ref="(el) => setButtonRef(el, row)" class="inline-flex items-center cursor-pointer" @click="showRoleDropdownMenu(row)">
              <span>{{ (UserRoleInTeamMap as any)[row.role] }}</span>
              <el-icon :class="['m-4px']">
                <ac-icon-ep-arrow-down />
              </el-icon>
            </div>
            <div v-else>
              <span>{{ (UserRoleInTeamMap as any)[row.role] }}</span>
            </div>
          </template>
        </el-table-column>
      </template>

      <template #operation v-if="isSuperAdmin">
        <el-table-column :label="$t('app.table.operation')">
          <template #default="{ row }">
            <template v-if="!row.isSelf">
              <el-button link size="small" @click="handelEditUser(row)">{{ $t('app.member.tips.editMember') }}</el-button>
              <el-button link type="danger" size="small" @click="handleRemove(row)">{{ $t('app.member.tips.removeMember') }}</el-button>
            </template>
          </template>
        </el-table-column>
      </template>
    </AcSimpleTable>
  </div>

  <AddMemberModal ref="addMemberModal" @ok="getTableData" />
  <UpdateMemberModal ref="updateMemberModal" @ok="getTableData" />

  <el-popover :visible="isShowRoleDropdownMenu" :virtual-ref="popoverRefEl" trigger="click" virtual-triggering>
    <PopperMenu :active-menu-key="currentChangeUser?.role" row-key="value" :menus="userRoles" size="small" class="clear-popover-space" @menu-click="handelChangeUserRole" />
  </el-popover>
</template>
<script setup lang="ts">
import { useTable } from '@/hooks/useTable'
import { getMembers, deleteMember, updateMember } from '@/api/member'
import { useI18n } from 'vue-i18n'
import AddMemberModal from './AddMemberModal.vue'
import { UserInfo, UserRoleInTeamMap } from '@/typings/user'
import { usePopover } from '@/hooks/usePopover'
import { useUserStore } from '@/store/user'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import NProgress from 'nprogress'
import UpdateMemberModal from './UpdateMemberModal.vue'

const { t } = useI18n()
const buttonRefMap: Record<number, any> = {}
const { userInfo, isSuperAdmin, userRoles } = useUserStore()
const currentChangeUser = ref<UserInfo | null>()

const addMemberModal = ref<InstanceType<typeof AddMemberModal>>()
const updateMemberModal = ref<InstanceType<typeof UpdateMemberModal>>()

const {
  isShow: isShowRoleDropdownMenu,
  popoverRefEl,
  showPopover,
  hidePopover,
} = usePopover({
  onHide: () => {
    currentChangeUser.value = null
  },
})

const { currentPage, pageSize, total, data, isLoading, getTableData } = useTable(getMembers, {
  dataKey: 'members',
  totalKey: 'total_member',
  isLoaded: true,
  transform: (user: UserInfo): UserInfo => {
    user.accountStatus = user.is_enabled ? t('app.member.form.accountStatusNormal') : t('app.member.form.accountStatusLock')
    user.accountStatusType = user.is_enabled ? '' : 'info'
    user.isSelf = user.id === userInfo.id
    return user
  },
})

const columns: any = [
  {
    label: t('app.member.form.name'),
    prop: 'username',
  },
  {
    label: t('app.member.form.email'),
    prop: 'email',
  },
  {
    slot: 'accountStatus',
  },
  {
    slot: 'auth',
  },
]

const showRoleDropdownMenu = (user: UserInfo) => {
  showPopover(buttonRefMap[user.id!])
  currentChangeUser.value = user
}

const setButtonRef = (el: any, user: UserInfo) => {
  buttonRefMap[user.id!] = el
}

const handleRemove = (user: UserInfo) => {
  AsyncMsgBox({
    title: t('app.common.deleteTip'),
    content: t('app.member.tips.deleteMemberTip'),
    onOk: async () => {
      await deleteMember(user.id!)
      await getTableData()
    },
  })
}

const handelChangeUserRole = async (role: any) => {
  if (!currentChangeUser.value) {
    return
  }

  const { id, is_enabled } = currentChangeUser.value
  NProgress.start()
  try {
    await updateMember({ id, role: role.value, is_enabled })
    hidePopover()
    await getTableData()
  } catch (error) {
    //
  } finally {
    NProgress.done()
  }
}

const handelEditUser = async (user: UserInfo) => {
  updateMemberModal.value?.show(user)
}
</script>
