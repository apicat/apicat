<template>
  <el-form :inline="true" :model="form" ref="fromRef" :rules="rules">
    <el-form-item :label="$t('app.project.member.chooseMember')" prop="user_ids" class="w-260px">
      <el-select
        ref="selectRef"
        v-model="form.user_ids"
        :placeholder="$t('app.project.member.chooseMember')"
        filterable
        multiple
        collapse-tags
        collapse-tags-tooltip
        class="w-full"
        @change="handleMemberChange"
      >
        <el-option v-for="member in members" :label="member.username" :value="member.user_id!" />
      </el-select>
    </el-form-item>
    <el-form-item :label="$t('app.project.member.chooseAuth')">
      <el-select v-model="form.authority">
        <el-option v-for="item in projectAuths" :label="item.text" :value="item.value" />
      </el-select>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" :loading="isLoading" @click="onSubmit(fromRef)">{{ $t('app.project.member.addMember') }}</el-button>
    </el-form-item>
  </el-form>
</template>
<script lang="ts" setup>
import { getMembersWithoutProject, addMemberToProject } from '@/api/project'
import useApi from '@/hooks/useApi'
import { useParams } from '@/hooks/useParams'
import useProjectStore from '@/store/project'
import { ProjectMember } from '@/typings/member'
import { MemberAuthorityInProject } from '@/typings/member'
import { FormInstance } from 'element-plus'
import { reactive } from 'vue'
import { useI18n } from 'vue-i18n'

const emits = defineEmits(['ok'])
const { t } = useI18n()
const { project_id } = useParams()
const { projectAuths } = useProjectStore()
const [isLoading, addMemberToProjectRequest] = useApi(addMemberToProject(project_id as string))
const fromRef = ref<FormInstance>()

const form = reactive({
  authority: MemberAuthorityInProject.READ,
  user_ids: [],
})
const selectRef = ref()

const rules = {
  user_ids: [{ required: true, message: t('app.project.member.chooseMember') }],
}

const members = shallowRef<ProjectMember[]>([])

const handleMemberChange = () => {
  selectRef.value.query = ''
}

const onSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  try {
    const valid = await formEl.validate()
    if (valid) {
      await addMemberToProjectRequest(toRaw(form))
      await refreshMemberList()
      formEl.resetFields()
      emits('ok')
    }
  } catch (error) {
    //
  }
}

const refreshMemberList = async () => {
  try {
    members.value = (await getMembersWithoutProject(project_id as string)) as any
  } catch (error) {
    members.value = []
  }
}

onMounted(async () => await refreshMemberList())

defineExpose({
  refreshMemberList,
})
</script>
