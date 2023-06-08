<template>
  <div class="container flex flex-col justify-center">
    <AddProjectMember v-if="isManager" ref="addProjectMemberRef" @ok="getTableData" />

    <AcSimpleTable isShowPager v-model:page="currentPage" v-model:page-size="pageSize" :columns="columns" :table-data="data" :loading="isLoading" :total="total">
      <template #auth>
        <el-table-column :label="$t('app.project.list.auth')">
          <template #default="{ row }">
            <div v-if="isManager && !row.isSelf" :ref="(el) => setButtonRef(el, row)" class="inline-flex items-center cursor-pointer" @click="showRoleDropdownMenu(row)">
              <span>{{ (MemberAuthorityMap as any)[row.authority] }}</span>
              <el-icon :class="['m-4px']">
                <ac-icon-ep-arrow-down />
              </el-icon>
            </div>
            <div v-else>
              <span>{{ (MemberAuthorityMap as any)[row.authority] }}</span>
            </div>
          </template>
        </el-table-column>
      </template>

      <template #operation v-if="isManager">
        <el-table-column :label="$t('app.table.operation')">
          <template #default="{ row }">
            <template v-if="!row.isSelf">
              <el-button link type="danger" size="small" @click="handlerRemoveMember(row)">{{ $t('app.project.member.deleteMember') }}</el-button>
              <el-button link size="small" v-if="row.authority === MemberAuthorityInProject.WRITE" @click="handlerTransferProject(row)">{{
                $t('app.project.member.transferProject')
              }}</el-button>
            </template>
          </template>
        </el-table-column>
      </template>
    </AcSimpleTable>
  </div>

  <el-popover :visible="isShowRoleDropdownMenu" :virtual-ref="popoverRefEl" trigger="click" virtual-triggering>
    <PopperMenu
      :active-menu-key="currentChangeUser?.authority"
      row-key="value"
      :menus="projectAuths"
      size="small"
      class="clear-popover-space"
      @menu-click="handlerChangeUserRole"
    />
  </el-popover>
</template>
<script setup lang="ts">
import { useTable } from '@/hooks/useTable'
import { getMembersInProject, deleteMemberFromProject, updateMemberAuthorityInProject, transferProject } from '@/api/project'
import { useI18n } from 'vue-i18n'
import { MemberAuthorityMap, ProjectMember, MemberAuthorityInProject } from '@/typings/member'
import { usePopover } from '@/hooks/usePopover'
import { useUserStore } from '@/store/user'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import NProgress from 'nprogress'
import { useParams } from '@/hooks/useParams'
import uesProjectStore from '@/store/project'
import AddProjectMember from '../AddProjectMember.vue'

const { t } = useI18n()
const { project_id } = useParams()
const buttonRefMap: Record<number, any> = {}
const { userInfo } = useUserStore()
const { projectAuths, isManager } = uesProjectStore()
const currentChangeUser = ref<ProjectMember | null>()
const addProjectMemberRef = ref<InstanceType<typeof AddProjectMember>>()

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

const { currentPage, pageSize, total, data, isLoading, getTableData } = useTable(getMembersInProject(project_id as string), {
  isLoaded: true,
  transform: (member: ProjectMember): ProjectMember => {
    member.isSelf = member.user_id === userInfo.id
    return member
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
    slot: 'auth',
  },
]

const showRoleDropdownMenu = (member: ProjectMember) => {
  showPopover(buttonRefMap[member.id!])
  currentChangeUser.value = member
}

const setButtonRef = (el: any, member: ProjectMember) => {
  buttonRefMap[member.id!] = el
}

// 移除成员
const handlerRemoveMember = (member: ProjectMember) => {
  AsyncMsgBox({
    title: t('app.common.deleteTip'),
    content: t('app.member.tips.deleteMemberTip'),
    onOk: async () => {
      await deleteMemberFromProject(project_id as string, member.user_id!)
      await getTableData()
      addProjectMemberRef.value?.refreshMemberList()
    },
  })
}

// 修改成员权限
const handlerChangeUserRole = async (role: any) => {
  if (!currentChangeUser.value) {
    return
  }

  const { user_id } = currentChangeUser.value
  NProgress.start()
  try {
    await updateMemberAuthorityInProject(project_id as string, user_id!, role.value)
    hidePopover()
    await getTableData()
  } catch (error) {
    //
  } finally {
    NProgress.done()
  }
}

// 移交项目
const handlerTransferProject = async (member: ProjectMember) => {
  AsyncMsgBox({
    title: t('app.common.deleteTip'),
    content: '确定移交项目给该成员?',
    onOk: async () => {
      await transferProject(project_id as string, member.id!)
      await getTableData()
    },
  })
}
</script>
