<template>
  <el-form ref="navFormRef" :model="form" :rules="rules" label-position="top" :label-width="0" v-loading="isLoadingUrls">
    <div ref="sortableList" :class="ns.b()">
      <el-form-item class="hide_required" :class="ns.e('item')" v-for="(item, index) in form.urls" :key="item._id">
        <el-row class="flex-1">
          <el-col :span="14" class="pr-2">
            <el-form-item :prop="'urls.' + index + '.url'" :rules="rules.urlItem">
              <el-input v-model.trim="item.url" :placeholder="$t('app.form.serverUrl.url')" maxlength="255" clearable />
            </el-form-item>
          </el-col>

          <el-col :span="6">
            <el-form-item :prop="'urls.' + index + '.description'">
              <el-input v-model.trim="item.description" :placeholder="$t('app.form.serverUrl.desc')" maxlength="255" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="4" class="flex items-center pl-4 operate">
            <el-icon class="icon-reorder mr-10px"><ac-icon-ep-sort /></el-icon>
            <el-icon @click="onRemoveNavItemIconClick(index)"><ac-icon-ep-delete /></el-icon>
          </el-col>
        </el-row>
      </el-form-item>
    </div>

    <el-form-item label="">
      <el-col :span="20">
        <el-button class="add-child-btn mt-2.5 w-full" @click="onAddNavItemBtnClick">
          {{ $t('app.common.add') }}
        </el-button>
      </el-col>
    </el-form-item>

    <el-form-item class="mt-5 mb-0">
      <el-button type="primary" :loading="isLoading" @click="handleSubmit(navFormRef)"> {{ $t('app.common.save') }} </el-button>
    </el-form-item>
  </el-form>
</template>
<script setup lang="ts">
import { useSortable } from '@/hooks/useSortable'
import { useNamespace } from '@/hooks'
import { FormInstance } from 'element-plus'
import isEmpty from 'lodash/isEmpty'
import isURL from 'validator/lib/isURL'
import { useProjectId } from '@/hooks/useProjectId'
import uesProjectStore from '@/store/project'
import useApi from '@/hooks/useApi'

const ns = useNamespace('url-list')
const project_id = useProjectId()
const projectStore = uesProjectStore()
const navFormRef = shallowRef()
const [isLoading, saveProjectServerUrlListApi] = useApi(projectStore.saveProjectServerUrlListApi)()
const [isLoadingUrls, getProjectServerUrlListApi] = useApi(projectStore.getUrlServers)()

type Url = {
  url: string
  description: string
  _id?: string | number
}

const form = reactive({
  urls: [] as Url[],
})

const rules = {
  urlItem: {
    validator: (rule: any, value: any, callback: any) => {
      let index = parseInt(rule.field.split('.')[1])
      let field = rule.field.split('.')[2]
      let { url } = (form.urls[index] || {}) as Url
      if ((isEmpty(url) || !isURL(url, { protocols: ['http', 'https'], require_protocol: true })) && field === 'url') {
        return callback(new Error('请输入有效的链接地址'))
      }
      callback()
    },
    trigger: 'blur',
  },
}

const handleSubmit = async (fromIns: FormInstance) => {
  try {
    const valid = await fromIns.validate()
    valid && (await saveProjectServerUrlListApi({ project_id, urls: toRaw(form.urls) }))
  } catch (error) {
    //
  }
}

const onRemoveNavItemIconClick = (index: number) => form.urls.splice(index, 1)

const onAddNavItemBtnClick = () => {
  form.urls.push({ description: '', url: '', _id: Date.now() })
}

const sortableList = ref()
const { initSortable } = useSortable(sortableList, {
  handle: '.icon-reorder',
  onEnd(e: any) {
    const changeItem = form.urls.splice(e.oldIndex, 1)[0]
    form.urls.splice(e.newIndex, 0, changeItem)
  },
})

onMounted(async () => {
  const urls = await getProjectServerUrlListApi(project_id)
  form.urls = urls.concat([])
  initSortable()
})
</script>
<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;

@include b(url-list) {
  @include e(item) {
    .operate {
      display: none;
    }

    &:hover .operate {
      display: block;
    }
  }

  .operate .icon-reorder {
    cursor: grabbing;
  }
}

.add-child-btn {
  border-radius: 2px !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
  color: rgba(0, 0, 0, 0.6) !important;
  border-style: dashed !important;
}
</style>
