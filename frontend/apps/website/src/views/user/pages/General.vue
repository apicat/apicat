<script setup lang="ts">
import { Picture as IconPicture } from '@element-plus/icons-vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { FileUploaderWrapper } from '@apicat/components'
import { apiUpdateGeneral, apiUploadAvatar } from '@/api/user'
import { useUserStore } from '@/store/user'
import { useLocaleStore } from '@/store/locale'
import ImageCorp from '@/views/user/pages/ImageCorp.vue'

type GeneralUserInfo = Pick<UserAPI.ResponseUserInfo, 'name' | 'avatar' | 'language'>

const fileUploaderWrapperRef = shallowRef()
const { t } = useI18n()
const store = useUserStore()
const localeStore = useLocaleStore()
const { userInfo } = storeToRefs(store)
const { languagesForSelect } = storeToRefs(localeStore)
const form = ref<GeneralUserInfo>({
  name: '',
  avatar: '',
  language: '',
})

const userinfo = computed((): GeneralUserInfo => {
  const keys = Object.keys(form.value)
  for (const key in userInfo.value) {
    if (keys.includes(key))
      form.value[key as keyof typeof form.value] = userInfo.value[key as keyof UserAPI.ResponseUserInfo]
  }
  return form.value
})

const formRef = ref<FormInstance>()
const rules: FormRules<typeof userinfo.value> = {
  name: [
    {
      required: true,
      type: 'string',
      message: t('app.rules.username.required'),
      trigger: 'blur',
    },
    {
      validator(_: any, value: string, callback: any) {
        if (value.length < 2 || value.length > 64)
          callback(new Error(t('app.rules.username.wrongLength')))
        else callback()
      },
      trigger: 'blur',
    },
  ],
  language: [
    {
      required: true,
      type: 'string',
      message: t('app.rules.lang.required'),
      trigger: 'blur',
    },
  ],
  avatar: [
    {
      type: 'string',
      message: t('app.rules.ava.required'),
      trigger: 'blur',
    },
  ],
}
const avatarURL = ref<string>()
let avatarFile: File
const imgCorpRef = ref<typeof ImageCorp>()

function handleChange(file: File) {
  const src: string = URL.createObjectURL(file)
  avatarURL.value = src
  avatarFile = file
  imgCorpRef.value!.show(src)
}

async function submit() {
  try {
    await formRef.value?.validate()
    await apiUpdateGeneral(form.value as UserAPI.RequestGeneral)
    await store.updateUserInfo(userinfo.value as UserAPI.ResponseUserInfo)
    ElMessage.success(t('app.user.general.success'))
  }
  catch (e) {
    //
  }
}

async function submitAvatar(data: UserAPI.RequestChangeAvatar) {
  data.avatar = avatarFile
  const res = await apiUploadAvatar(data)
  userinfo.value.avatar = res.avatar
  store.updateUserInfo(res as UserAPI.ResponseUserInfo)
}
</script>

<template>
  <ImageCorp ref="imgCorpRef" :handle-upload="submitAvatar" />

  <div class="flex flex-col justify-center mx-auto px-36px" style="align-items: center">
    <div class="items-start text-start w-40vw">
      <div class="bg-white w-450px">
        <h1 class="font-500 text-24px text-gray-title">
          {{ $t('app.user.general.title') }}
        </h1>
        <ElForm ref="formRef" :rules="rules" :model="form" class="content" label-position="top">
          <!-- head img -->
          <ElFormItem prop="avatar">
            <div class="row">
              <div class="left" style="margin-right: 50px">
                <FileUploaderWrapper ref="fileUploaderWrapperRef" accept=".png, .jpg, .jpeg" @change="handleChange">
                  <el-image style="cursor: pointer" :src="userinfo.avatar">
                    <template #error>
                      <div class="image-slot">
                        <el-icon><IconPicture /></el-icon>
                      </div>
                    </template>
                  </el-image>
                </FileUploaderWrapper>
              </div>
            </div>
          </ElFormItem>

          <div style="margin-top: 40px">
            <!-- username -->
            <ElFormItem prop="name" :label="$t('app.user.general.username')">
              <ElInput v-model="userinfo.name" maxlength="64" class="h-40px" />
            </ElFormItem>

            <!-- lang -->
            <ElFormItem prop="language" :label="$t('app.user.general.lang')">
              <ElSelect v-model="userinfo.language" class="w-full">
                <el-option
                  v-for="item in languagesForSelect"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </ElSelect>
            </ElFormItem>
          </div>

          <!-- submit -->
          <ElButton class="w-full" type="primary" @click="submit">
            {{ $t('app.user.general.save') }}
          </ElButton>
        </ElForm>
      </div>
    </div>
  </div>
</template>

<style scoped>
:deep(.el-select .el-input) {
  height: 40px;
}

:deep(.el-button) {
  height: 40px;
}

.row {
  margin-top: 1em;
  margin-bottom: 1em;
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
  justify-content: flex-start;
  /* flex-grow: 1; */
}

.right {
  /* justify-content: flex-end; */
  flex-grow: 1;
}

.content {
  margin-top: 40px;
}

/* el-upload */
:deep(.content .el-upload) {
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

/* el-image */
.content .block {
  padding: 30px 0;
  text-align: center;
  border-right: solid 1px var(--el-border-color);
  display: inline-block;
  width: 49%;
  box-sizing: border-box;
  vertical-align: top;
}
.content .demonstration {
  display: block;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  margin-bottom: 20px;
}
.content .el-image {
  width: 200px;
  height: 200px;
  border-radius: 50%;
}

.content .image-slot {
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
