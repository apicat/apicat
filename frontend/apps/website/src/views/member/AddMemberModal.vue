<template>
  <el-dialog v-model="dialogVisible" :title="$t('app.member.tips.addMember')" :width="300" align-center :close-on-click-modal="false">
    <!-- 内容 -->
    <el-form label-position="top" label-width="100px" :model="form" :rules="rules" ref="fromRef" @submit.prevent="handleSubmit(fromRef)">
      <el-form-item :label="$t('app.form.user.email')" prop="email">
        <el-input v-model="form.email" :placeholder="$t('app.rules.email.required')" clearable />
      </el-form-item>

      <el-form-item :label="$t('app.form.user.password')" prop="password">
        <el-input show-password type="password" v-model="form.password" :placeholder="$t('app.rules.password.required')" clearable />
      </el-form-item>

      <el-form-item :label="$t('app.member.form.role')" prop="role">
        <el-select v-model="form.role" class="w-full">
          <template v-for="(item, key) in UserRoleInTeamMap" :key="key">
            <el-option v-if="key !== UserRoleInTeam.SUPER_ADMIN" :label="item" :value="key"></el-option>
          </template>
        </el-select>
      </el-form-item>
    </el-form>
    <!-- 底部按钮 -->
    <div slot="footer" class="text-right mt-20px">
      <el-button @click="dialogVisible = false">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="isLoading" @click="handleSubmit(fromRef)">
        {{ $t('app.common.add') }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { useModal } from '@/hooks'
import { useI18n } from 'vue-i18n'
import type { FormInstance, FormRules } from 'element-plus'
import { addMember } from '@/api/member'
import { UserInfo, UserRoleInTeam, UserRoleInTeamMap } from '@/typings/user'
import { isEmail } from '@apicat/shared'
import useApi from '@/hooks/useApi'

const emits = defineEmits(['ok'])

const { t } = useI18n()

const fromRef = ref<FormInstance>()
const { dialogVisible, showModel, hideModel } = useModal(fromRef as any)
const [isLoading, addMemberRequest] = useApi(addMember)

const form = reactive<UserInfo>({
  email: '',
  password: '',
  role: UserRoleInTeam.USER,
})

const rules = reactive<FormRules>({
  email: [
    { required: true, message: t('app.rules.email.required'), trigger: 'blur' },
    {
      validator(rule: any, value: any, callback: any) {
        if (!isEmail(value)) {
          callback(new Error(t('app.rules.email.correct')))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
  password: [{ required: true, min: 8, message: t('app.rules.password.minLength'), trigger: 'blur' }],
})

const handleSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  try {
    const valid = await formEl.validate()
    if (valid) {
      const user: UserInfo = await addMemberRequest(toRaw(form))
      hideModel()
      emits('ok', user)
    }
  } catch (error) {
    //
  }
}

defineExpose({
  show: showModel,
})
</script>
