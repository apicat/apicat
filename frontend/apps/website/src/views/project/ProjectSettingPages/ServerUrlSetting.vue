<script setup lang="ts">
import isURL from 'validator/lib/isURL'
import { useI18n } from 'vue-i18n'
import { isEmpty } from 'lodash-es'
import { storeToRefs } from 'pinia'
import { useNamespace } from '@apicat/hooks'
import { VueDraggableNext } from 'vue-draggable-next'
import useProjectStore from '@/store/project'
import { useServerUrlCompare } from '@/views/project/logic/useServerUrlCompare'

const ns = useNamespace('url-list')
const { t } = useI18n()
const projectStore = useProjectStore()
const { isReader } = storeToRefs(projectStore)

const { addURL, deleteURL, urlList, isLoadingUrls, draglist, saveURL, navFormRef } = useServerUrlCompare()

function isV4Localhost(url: string) {
  try {
    return new URL(url).host.includes('localhost')
  }
  catch (e) {
    return false
  }
}

const rules = {
  urlItem: {
    validator: (_: any, value: any, callback: any) => {
      const url = value
      if (
        isEmpty(url)
        || !(
          isURL(url, {
            protocols: ['http', 'https'],
            require_protocol: true,
          }) || isV4Localhost(url)
        )
      )
        return callback(new Error(t('app.serverUrl.rules.invalid')))

      let count = 2
      for (const key in urlList.value) {
        if (key in urlList.value) {
          const v = urlList.value[key]
          if (v.data.url === value)
            count--
          if (!count)
            return callback(new Error(t('app.serverUrl.rules.duplicated')))
        }
      }

      callback()
    },
    trigger: 'change',
  },
}
</script>

<template>
  <el-form
    ref="navFormRef"
    v-loading="isLoadingUrls"
    :model="urlList"
    :rules="rules"
    label-position="top"
    :label-width="0"
  >
    <VueDraggableNext handle=".icon-reorder" :list="draglist">
      <el-form-item v-for="rowID in draglist" :key="rowID" class="hide_required drag-node" :class="ns.e('item')">
        <div v-if="urlList[rowID]" class="row">
          <div class="left" style="margin-right: 10px">
            <div class="row">
              <div class="left" style="margin-right: 10px">
                <el-form-item style="width: 100%" :prop="`${rowID}.data.url`" :rules="rules.urlItem">
                  <el-input
                    v-model="urlList[rowID].data.url"
                    clearable
                    maxlength="255"
                    :disabled="isReader"
                    :placeholder="$t('app.form.serverUrl.url')"
                  />
                </el-form-item>
              </div>
              <div class="right">
                <el-form-item style="width: 250px">
                  <el-input
                    v-model="urlList[rowID].data.description"
                    clearable
                    maxlength="255"
                    :disabled="isReader"
                    :placeholder="$t('app.form.serverUrl.desc')"
                  />
                </el-form-item>
              </div>
            </div>
          </div>
          <div class="right">
            <div v-if="!isReader" class="operate">
              <el-icon style="cursor: pointer" class="icon-reorder mr-10px">
                <ac-icon-ep-sort />
              </el-icon>
              <el-icon style="cursor: pointer" @click="deleteURL(rowID)">
                <ac-icon-ep-delete />
              </el-icon>
            </div>
          </div>
        </div>
      </el-form-item>
    </VueDraggableNext>
    <el-form-item v-if="!isReader" label="">
      <el-col :span="24">
        <el-button class="w-full add-child-btn" @click="addURL">
          {{ $t('app.common.add') }}
        </el-button>
      </el-col>
    </el-form-item>

    <el-form-item v-if="!isReader" style="margin-top: 20px">
      <el-button type="primary" @click="saveURL">
        {{ $t('app.project.setting.basic.update') }}
      </el-button>
    </el-form-item>
  </el-form>
</template>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins.scss' as *;

@include b(url-list) {
  @include e(item) {
    .operate {
      // display: none;
      visibility: hidden;
    }

    &:hover .operate {
      // display: block;
      visibility: unset;
    }
  }

  .icon-reorder {
    cursor: pointer;
  }

  .operate .icon-reorder {
    cursor: grabbing;
  }
}

.add-child-btn {
  border-radius: var(--el-border-radius-base) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
  color: rgba(0, 0, 0, 0.6) !important;
  border-style: dashed !important;
}

.edit-node {
  padding: 0.3em;
  border-radius: 5px;
  background-color: #eef5ff;
}

.el-form-item {
  margin: 0;
}

.drag-node {
  margin-bottom: 6px;
}

.row {
  margin-bottom: 5px;
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
</style>
