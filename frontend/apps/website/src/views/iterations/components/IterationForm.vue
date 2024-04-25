<script setup lang="ts">
import { useNamespace } from '@apicat/hooks'
import { useIterationPlan } from '../logic/useIterationPlan'
import { useIterationForm } from '../logic/useIterationForm'
import { CollectionTypeEnum } from '@/commons'

export interface IterationFormProps {
  iterationID: string | null
  projects: ProjectAPI.ResponseProject[]
}
export interface IterationFormEmits {
  (event: 'success' | 'cancel'): void
}

const props = defineProps<IterationFormProps>()
const emits = defineEmits<IterationFormEmits>()

const ns = useNamespace('iteration-detail')

const {
  isEditMode,
  isLoadingForSubmit,
  isLoadedIteration,
  iterationFormRef,
  iterationRules,
  iterationInfo,
  handleCancel,
  handleSubmit,
} = useIterationForm(props, emits)

const { isLoadingForTree, defaultProps, collections, selectedCollectionKeys, onTreeChange } = useIterationPlan(
  props,
  iterationInfo,
)
</script>

<template>
  <div v-loading="isLoadedIteration" :class="ns.b()" class="min-w-500px md:w-full lg:max-w-800px">
    <p class="border-b border-solid border-gray-lighter pb-15px text-24px text-gray-title mb-30px">
      {{ isEditMode ? $t('app.iter.create.edit_title') : $t('app.iter.create.title') }}
    </p>
    <el-form ref="iterationFormRef" :model="iterationInfo" :rules="iterationRules" label-position="top">
      <el-form-item :label="$t('app.iter.create.name')" prop="title">
        <el-input v-model="iterationInfo.title" maxlength="255" :placeholder="$t('app.iter.create.name_hold')" />
      </el-form-item>
      <el-form-item :label="$t('app.iter.create.project')" prop="projectID">
        <el-select
          v-model="iterationInfo.projectID" :disabled="isEditMode" popper-class="ac-select-popper"
          :teleported="false" class="w-full" :placeholder="$t('app.iter.create.project_hold')"
        >
          <el-option v-for="item in projects" :key="item.id" :label="item.title" :value="item.id">
            <div class="truncate max-w-700px">
              {{ item.title }}
            </div>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('app.iter.create.desc')" prop="description">
        <el-input v-model="iterationInfo.description" type="textarea" :rows="4" />
      </el-form-item>

      <el-form-item :label="$t('app.iter.create.plan.title')" prop="collectionIDs">
        <p class="text-gray-helper -mt-10px">
          {{ $t('app.iter.create.plan.tip') }}
        </p>
      </el-form-item>

      <el-form-item v-loading="isLoadingForTree">
        <AcTransferTree
          height="500px" :data="collections" :default-checked-keys="selectedCollectionKeys"
          :placeholder="$t('app.iter.create.plan.table.hold')" :title="$t('app.iter.create.plan.table.title')"
          node-key="id" :default-props="defaultProps" @change="onTreeChange"
        >
          <template #default="{ node, data }">
            <div class="flex items-center flex-1 overflow-hidden cursor-pointer">
              <i v-if="data.type !== CollectionTypeEnum.Dir" alt="" class="ac-doc ac-iconfont mr-2px" />
              <ac-icon-ic:outline-folder v-else class="mr-2px" />
              <label :title="node.label" class="truncate">{{ node.label }}</label>
            </div>
          </template>
        </AcTransferTree>
      </el-form-item>

      <el-form-item>
        <div class="flex-1 text-right">
          <el-button @click="handleCancel">
            {{ $t('app.common.cancel') }}
          </el-button>
          <el-button type="primary" :loading="isLoadingForSubmit" @click="handleSubmit(iterationFormRef)">
            {{ isEditMode ? $t('app.common.edit') : $t('app.iteration.form.create') }}
          </el-button>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>
