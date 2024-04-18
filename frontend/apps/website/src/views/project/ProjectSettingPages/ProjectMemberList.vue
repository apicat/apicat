<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import NProgress from 'nprogress'
import { storeToRefs } from 'pinia'
import type { FormInstance } from 'element-plus'
import { useTablev2 } from '@/hooks/useTable'
import { usePopover } from '@/hooks/usePopover'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import useProjectStore from '@/store/project'
import { useApi } from '@/hooks/useApi'
import {
  apiCreateProjectMember,
  apiEditProjectMember,
  apiGetExcludedMembers,
  apiGetProjectMembers,
  apiRemoveProjectMember,
} from '@/api/project/setting/member'
import { Authority, Status } from '@/commons/constant'
import { useUserStore } from '@/store/user'
import MemberInfo from '@/views/team/pages/member/MemberInfo.vue'
import EmptyAvatar from '@/components/EmptyAvatar.vue'

const { t } = useI18n()
const buttonRefMap: Record<number, any> = {}
const projectStore = useProjectStore()
const { isManager } = storeToRefs(projectStore)
const currentChangeUser = ref<ProjectAPI.Member | null>()

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

const { currentPage, pageSize, total, data, isLoading, refreshData } = useTablev2(apiGetProjectMembers, {
  addonArgs: [projectStore.project?.id],
})

const columns: any = [
  {
    slot: 'name',
  },
  {
    slot: 'email',
  },
  {
    slot: 'authority',
  },
  {
    slot: 'status',
  },
]

function showRoleDropdownMenu(member: ProjectAPI.Member) {
  showPopover(buttonRefMap[member.id!])
  currentChangeUser.value = member
}

function setButtonRef(el: any, member: ProjectAPI.Member) {
  buttonRefMap[member.id!] = el
}

// 修改成员权限
async function handlerChangeUserAuth(auth: { text: string; value: Authority }) {
  if (!currentChangeUser.value)
    return

  NProgress.start()
  try {
    await apiEditProjectMember(projectStore.project!.id, currentChangeUser.value.id!, {
      permission: auth.value,
    })
    hidePopover()
    await refreshData()
  }
  catch (error) {
    //
  }
  finally {
    NProgress.done()
  }
}

// 移除成员
function handlerRemoveMember(member: ProjectAPI.Member) {
  AsyncMsgBox({
    title: t('app.project.setting.member.poprm.title'),
    content: t('app.member.tips.deleteMemberTip'),
    cancelButtonText: t('app.common.cancel'),
    confirmButtonClass: 'red',
    confirmButtonText: t('app.project.setting.member.table.rm'),
    onOk: async () => {
      try {
        await apiRemoveProjectMember(projectStore.project!.id, member.id!)
      }
      catch (error) {
        //
      }
      finally {
        await refreshData()
      }
    },
  })
}

// 移交项目
async function handlerTransferProject(member: ProjectAPI.Member) {
  AsyncMsgBox({
    title: t('app.project.setting.member.transfer.pop.title'),
    content: t('app.project.setting.member.transfer.pop.tip'),
    cancelButtonText: t('app.common.cancel'),
    confirmButtonText: t('app.project.setting.member.transfer.btn'),
    onOk: async () => {
      try {
        await projectStore.transferProject(member.id!)
      }
      catch (error) {
        //
      }
      finally {
        await refreshData()
      }
    },
  })
}

const excludedMembers = ref<TeamAPI.TeamMember[]>([])
const permissions = (() => {
  const map: Record<string, string> = {}
  for (const key in Authority) {
    const val = (Authority as any)[key as any]
    map[val] = t(`app.project.setting.member.auth.${val}`)
  }
  return map
})()
const editablePermissions = computed(() => {
  const map: Record<string, string> = {}
  for (const key in permissions) {
    if (key === Authority.None || key === Authority.Manage)
      continue
    map[key] = permissions[key]
  }
  return map
})
const addMemberForm = ref<{
  memberIDs: number[]
  permission: Authority
}>({
  memberIDs: [],
  permission: Authority.Read,
})
const rules = {
  memberIDs: [
    {
      trigger: 'blur',
      validator: (_: any, value: number[], callback: any) => {
        if (value.length === 0)
          callback(t('app.project.setting.member.form.memberIDs'))

        callback()
      },
    },
    {
      trigger: 'blur',
      required: true,
    },
  ],
  permission: [
    {
      trigger: 'blur',
      required: true,
    },
  ],
}
const formRef = ref<FormInstance>()
const [excludedLoading, _getExcludedMembers] = useApi(apiGetExcludedMembers)
async function getExcludedMembers(v: boolean) {
  if (v) {
    const res = await _getExcludedMembers(projectStore.project!.id)
    excludedMembers.value = (res || []).filter(val => val)!
  }
}

async function addMember() {
  try {
    await formRef.value?.validate()
    await apiCreateProjectMember({
      id: projectStore.project!.id,
      ...addMemberForm.value,
    })
    addMemberForm.value.memberIDs = []
    refreshData()
  }
  catch (e) {
    console.error(e)
  }
}
const userStore = useUserStore()
function isSelf(row: TeamAPI.TeamMember): boolean {
  return userStore.userInfo.email === row.user.email
}
</script>

<template>
  <div class="container flex flex-col justify-center">
    <div v-if="isManager" class="row mb-20px">
      <div class="left">
        <ElForm ref="formRef" :rules="rules" class="flex w-full" :model="addMemberForm">
          <ElFormItem prop="memberIDs" class="w-1/2">
            <ElSelect
              v-model="addMemberForm.memberIDs"
              :placeholder="$t('app.project.setting.member.smem')"
              class="w-full mr-3"
              multiple
              collapse-tags
              collapse-tags-tooltip
              :max-collapse-tags="2"
              :loading="excludedLoading"
              @visible-change="getExcludedMembers"
            >
              <ElOption v-for="item in excludedMembers" :key="item.id" :label="item.user.name" :value="item.id" />
            </ElSelect>
          </ElFormItem>
          <ElFormItem required>
            <ElSelect v-model="addMemberForm.permission" :placeholder="$t('app.project.setting.member.sper')">
              <ElOption v-for="(val, key) in editablePermissions" :key="key" :label="val" :value="key" />
            </ElSelect>
          </ElFormItem>
          <ElFormItem>
            <ElButton class="ml-3" type="primary" @click="addMember">
              {{ $t('app.project.setting.member.add') }}
            </ElButton>
          </ElFormItem>
        </ElForm>
      </div>
    </div>

    <AcSimpleTable
      v-model:page="currentPage"
      v-model:page-size="pageSize"
      is-show-pager
      round-border
      row-class-name="memberRow"
      :border="false"
      :columns="columns"
      :table-data="data"
      :loading="isLoading"
      :total="total"
    >
      <template #name>
        <el-table-column show-overflow-tooltip width="200" :label="$t('app.project.setting.member.table.name')">
          <template #default="{ row }">
            <div class="row" style="margin: 0">
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
        <el-table-column show-overflow-tooltip :label="$t('app.project.setting.member.table.email')">
          <template #default="{ row }">
            {{ row.user.email }}
          </template>
        </el-table-column>
      </template>

      <template #authority>
        <el-table-column width="105" :label="$t('app.project.setting.member.table.permissions')">
          <template #default="{ row }">
            <div
              v-if="isManager && !isSelf(row)"
              :ref="(el) => setButtonRef(el, row)"
              class="inline-flex items-center w-full cursor-pointer row auth-item"
              @click="showRoleDropdownMenu(row)"
            >
              <div class="left">
                <span>{{ permissions[row.permission] }}</span>
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
              <span>{{ permissions[row.permission] }}</span>
            </div>
          </template>
        </el-table-column>
      </template>

      <template #status>
        <el-table-column :label="$t('app.project.setting.member.table.status')" width="90" align="center">
          <template #default="{ row }">
            <p v-if="row.status" :color="row.status === Status.Active ? undefined : 'orange'">
              {{
                row.status === Status.Active
                  ? $t(`app.project.setting.member.table.active`)
                  : $t(`app.project.setting.member.table.inactive`)
              }}
            </p>
          </template>
        </el-table-column>
      </template>

      <template v-if="isManager" #operation>
        <el-table-column width="150" align="center">
          <template #default="{ row }">
            <div v-if="!isSelf(row)" class="opr-item">
              <el-button
                v-if="row.permission === Authority.Write && row.status === Status.Active"
                link
                size="small"
                @click="handlerTransferProject(row)"
              >
                {{ $t('app.project.setting.member.transfer.btn') }}
              </el-button>
              <el-button link type="danger" size="small" @click="handlerRemoveMember(row)">
                {{ $t('app.project.setting.member.table.rm') }}
              </el-button>
            </div>
          </template>
        </el-table-column>
      </template>
    </AcSimpleTable>
  </div>

  <el-popover
    :visible="isShowRoleDropdownMenu"
    :virtual-ref="popoverRefEl"
    :show-arrow="false"
    transition="fade-fast"
    width="80"
    trigger="hover"
    virtual-triggering
  >
    <PopperMenu
      row-key="value"
      center
      :active-menu-key="currentChangeUser?.permission"
      :menus="
        Object.keys(editablePermissions).map((val) => {
          return { text: permissions[val], value: val }
        })
      "
      size="small"
      class="clear-popover-space"
      @menu-click="handlerChangeUserAuth"
    />
  </el-popover>
</template>

<style lang="scss" scoped>
:deep(.el-form-item) {
  margin: 0;
}

.avatar {
  border-radius: 50%;
  width: 30px;
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
  // justify-content: flex-start;
  flex-grow: 1;
}
.right {
  justify-content: flex-end;
  // flex-grow: 1;
}

.memberRow .auth-item .arrow-icon {
  opacity: 0;
  transition: all 0.03s;
  transform: translateY(1px);
}
.memberRow:hover .auth-item .arrow-icon {
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
