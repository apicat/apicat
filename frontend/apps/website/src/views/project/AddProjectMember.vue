<template>
  <el-form :inline="true" :model="formInline" class="demo-form-inline">
    <el-form-item label="Approved by">
      <el-input v-model="formInline.user" placeholder="Approved by" />
    </el-form-item>
    <el-form-item label="Activity zone">
      <el-select v-model="formInline.region" placeholder="Activity zone">
        <el-option label="Zone one" value="shanghai" />
        <el-option label="Zone two" value="beijing" />
      </el-select>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="onSubmit">{{ $t('app.project.member.addMember') }}</el-button>
    </el-form-item>
  </el-form>
</template>
<script lang="ts" setup>
import { getMembersWithoutProject } from '@/api/project'
import { useParams } from '@/hooks/useParams'
import { reactive } from 'vue'

const { project_id } = useParams()
const formInline = reactive({
  user: '',
  region: '',
})

const onSubmit = () => {
  console.log('submit!')
}

const onRefreshMemberList = async () => {
  await getMembersWithoutProject(project_id as string)
}

defineExpose({
  onRefreshMemberList,
})
</script>
