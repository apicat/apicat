<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import ResetPasswordDialog from '../ResetPasswordDialog.vue'
import { useUserStore } from '@/store/user'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import EmptyAvatar from '@/components/EmptyAvatar.vue'
import { apiDeleteSystemUser, apiGetSystemUserList } from '@/api/user'
import useTable from '@/hooks/useTable'

const { t } = useI18n()
const { userInfo } = useUserStore()
const keywordRef = ref('')
const resetPasswordDialogRef = ref<InstanceType<typeof ResetPasswordDialog>>()
const { currentPage, pageSize, data, isLoading, total, getTableData } = useTable(apiGetSystemUserList, {
  totalKey: 'totalCount',
  searchParam: { keywords: keywordRef },
})

const columns: any = [{ label: 'ID', prop: 'id', width: 30 }, { slot: 'name' }, { slot: 'email' }]

function handlerRemove(user: UserAPI.ResponseUserInfo) {
  AsyncMsgBox({
    confirmButtonClass: 'red',
    confirmButtonText: t('app.team.member.remove.btn'),
    cancelButtonText: t('app.common.cancel'),
    title: t('app.system.users.removeUserTitle'),
    content: t('app.system.users.removeUserTip'),
    onOk: async () => {
      await apiDeleteSystemUser(user.id)
      await getTableData()
    },
  })
}

function showChangePasswordDialog(user: UserAPI.ResponseUserInfo) {
  resetPasswordDialogRef.value?.show(user)
}

function isSelf(row: UserAPI.ResponseUserInfo) {
  return row.email === userInfo.email
}
</script>

<template>
  <h1>{{ $t('app.system.users.title') }}</h1>
  <AcSimpleTable
    v-model:page="currentPage"
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
    <template #search-form>
      <el-form @submit.prevent="getTableData">
        <div class="grid grid-cols-2 gap-3">
          <el-form-item label="">
            <el-input v-model="keywordRef" :placeholder="$t('app.system.users.searchKeywordPlaceholder')" clearable />
          </el-form-item>

          <el-form-item>
            <el-button type="" @click="getTableData">
              {{ $t('app.common.search') }}
            </el-button>
          </el-form-item>
        </div>
      </el-form>
    </template>

    <template #name>
      <el-table-column show-overflow-tooltip :label="$t('app.team.member.table.name')">
        <template #default="{ row }">
          <div class="row">
            <div class="w-full left">
              <img v-if="row.avatar" class="avatar" :src="row.avatar">
              <EmptyAvatar v-else />
              <span class="ml-2 truncate">
                {{ isSelf(row) ? `${row.name}(${t('app.common.self')})` : row.name }}
              </span>
            </div>
          </div>
        </template>
      </el-table-column>
    </template>

    <template #email>
      <el-table-column show-overflow-tooltip :label="$t('app.team.member.table.email')">
        <template #default="{ row }">
          <span>{{ row.email }}</span>
        </template>
      </el-table-column>
    </template>

    <template #operation>
      <el-table-column width="200" align="center">
        <template #default="{ row }">
          <div v-if="!isSelf(row)">
            <!-- password -->
            <el-button link type="default" @click="showChangePasswordDialog(row)">
              {{ $t('app.system.users.updatePassword') }}
            </el-button>
            <!-- remove -->
            <el-button link type="default" @click="handlerRemove(row)">
              {{ $t('app.team.member.remove.btn') }}
            </el-button>
          </div>
        </template>
      </el-table-column>
    </template>
  </AcSimpleTable>

  <ResetPasswordDialog ref="resetPasswordDialogRef" />
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
