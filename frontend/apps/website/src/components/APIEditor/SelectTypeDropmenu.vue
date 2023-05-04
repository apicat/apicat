<template>
  <div class="ac-sce-selecttype">
    <div v-if="showRef" style="padding: 12px; width: 240px">
      <el-space direction="vertical" fill warp style="width: 100%">
        <el-text>{{ $t('editor.table.refModel') }}</el-text>
        <el-input v-model="searchRef">
          <template #prefix>
            <el-icon>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
                <path
                  d="M456.69 421.39L362.6 327.3a173.81 173.81 0 0 0 34.84-104.58C397.44 126.38 319.06 48 222.72 48S48 126.38 48 222.72s78.38 174.72 174.72 174.72A173.81 173.81 0 0 0 327.3 362.6l94.09 94.09a25 25 0 0 0 35.3-35.3zM97.92 222.72a124.8 124.8 0 1 1 124.8 124.8a124.95 124.95 0 0 1-124.8-124.8z"
                  fill="currentColor"
                ></path>
              </svg>
            </el-icon>
          </template>
        </el-input>
        <el-radio-group v-model="refName" @change="changeSchemaTypeRef">
          <el-tree :data="treeData" check-on-click-node class="typeSelect" style="width: 100%; max-height: 300px; overflow-y: scroll">
            <template #default="{ data }">
              <span v-if="data.isDir" class="-ml-4px">
                <el-space align-items="center" :size="4">
                  <el-icon :size="20">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16">
                      <g fill="none">
                        <path
                          d="M2 5v6a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2H7.175l-1.113-.89A.5.5 0 0 0 5.75 3H4a2 2 0 0 0-2 2zm1 0a1 1 0 0 1 1-1h1.575l.868.694l-.886.806H3V5zm4.593 0H12a1 1 0 0 1 1 1v5a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V6.5h2.75a.5.5 0 0 0 .336-.13L7.593 5z"
                          fill="currentColor"
                        ></path>
                      </g>
                    </svg>
                  </el-icon>
                  <span>{{ data.label }}</span>
                </el-space>
              </span>
              <el-radio class="flex-1" :label="data.key" v-else style="margin-left: -1px; --el-radio-font-weight: 400">
                {{ data.label }}
              </el-radio>
            </template>
          </el-tree>
        </el-radio-group>
      </el-space>
    </div>
    <el-dropdown-menu v-else>
      <el-dropdown-item :class="{ 'ac-sec-selected': data.refObj }" @click.prevent="openRefMode">{{ $t('editor.table.refModel') }}</el-dropdown-item>
      <el-dropdown-item
        v-for="(item, i) in basicTypes"
        :divided="i == 0"
        :key="item"
        :class="{
          'ac-sec-selected': !data.refObj && data.type == item,
        }"
        @click="changeSchemaType(item)"
        >{{ item }}</el-dropdown-item
      >
    </el-dropdown-menu>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, ref, watchEffect } from 'vue'
import type { Tree, Definition } from './types'
import { basicTypes } from './types'
import { RefPrefixKeys } from '@/commons'
const props = defineProps<{
  data: Tree
  showRef: boolean
}>()

const emits = defineEmits(['showRef', 'change'])

interface treeNode {
  key: number
  label: string
  isDir: boolean
  children?: treeNode[]
}

const DefSchemas = inject('definitions') as () => Definition[]

function listToTree(parentId: number): treeNode[] {
  const tree: treeNode[] = []
  DefSchemas()
    ?.filter((ele) => ele.parent_id === parentId)
    .forEach((ele) => {
      const item: treeNode = {
        key: ele.id as any,
        label: ele.name,
        isDir: ele.type === 'category',
      }
      if (item.isDir) {
        item.children = listToTree(ele.id as any)
      }
      tree.push(item)
    })
  return tree.sort((v) => (v.isDir ? -1 : 1))
}
const searchRef = ref('')
const treeData = computed(() => {
  if (searchRef.value == '') {
    return listToTree(0)
  }
  return DefSchemas()
    .filter((v) => v.type === 'schema' && v.name.toLowerCase().indexOf(searchRef.value.toLowerCase()) != -1)
    .map((v) => {
      return {
        key: v.id,
        label: v.name,
      }
    })
})

function resetObject(v: Object) {
  for (let k of Object.keys(v)) {
    if (k != 'description') {
      delete (v as any)[k]
    }
  }
}

const refName = ref()
watchEffect(() => {
  refName.value = props.data.refObj?.id
})
const openRefMode = () => {
  emits('showRef', true)
}
const changeSchemaTypeRef = (r: any) => {
  if (props.data.refObj?.name == r) {
    return
  }
  const sc = props.data.schema
  resetObject(sc)
  sc.$ref = `${RefPrefixKeys.DefinitionsSchema.key}${r}`
  emits('change')
}

const changeSchemaType = (vtype: string) => {
  if (vtype == props.data.type) {
    return
  }
  const sc = props.data.schema
  resetObject(sc)
  sc.type = vtype
  if (sc.type == 'array') {
    sc.items = {
      type: 'string',
    }
  } else if (sc.type == 'object') {
    sc.properties = {}
  }
  emits('change')
}
</script>

<style>
.ac-sce-selecttype .ac-sec-selected {
  color: var(--el-color-primary);
}

.ac-sce-selecttype {
  user-select: none;
  -webkit-user-select: none;
}
</style>
