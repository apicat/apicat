<script setup lang="ts">
import type { CheckboxValueType } from 'element-plus'
import { ElTree } from 'element-plus'
import type { TreeNode } from 'element-plus/es/components/tree-v2/src/types'
import { useI18n } from 'vue-i18n'

const props = withDefaults(
  defineProps<{
    title: string
    nodeKey: string
    data: any[]
    defaultCheckedKeys?: any[]
    width?: string
    height?: string
    filter?: boolean
    placeholder?: string
    defaultProps?: {
      children: string
      label: string
      parentKey?: string
      rootValue?: any
    }
    parentId?: string
  }>(),
  {
    width: '100%',
    height: '320px',
    nodeKey: 'id',
    parentId: 'parent_id',
    filter: true,
    defaultProps: () => ({
      children: 'children',
      label: 'title',
      parentKey: 'parent_id',
      rootValue: 0,
    }),
  },
)

const emits = defineEmits<{
  (event: 'change', value: any[]): void
}>()

const { t } = useI18n()

const checkedKeys = computed(() => props.defaultCheckedKeys?.map(val => val[props.nodeKey]))

const treeIns = ref<InstanceType<typeof ElTree>>()
const checkAll = ref(false)
const isIndeterminate = ref(false)
const filterKeyword = ref('')
const checkedNodes = ref<any[]>([])
const {
  rootValue = 0,
  parentKey = 'parent_id',
  ...treePorps
} = props.defaultProps || { children: 'children', label: 'title' }

function notifyChange() {
  emits('change', treeIns.value?.getCheckedNodes() || [])
}

function handleCheckAll(isChecked: CheckboxValueType) {
  checkedNodes.value = isChecked ? props.data : []
  treeIns.value?.setCheckedNodes(checkedNodes.value)
  notifyChange()
}

// 切换选中节点
function onNodeCheckChange() {
  checkedNodes.value = treeIns.value?.getCheckedNodes() || []
  notifyChange()
}

// 点击非目录节点，切换选中
function onNodeClick(data: any, node: TreeNode) {
  if (!node.isLeaf)
    return

  const keys = treeIns.value?.getCheckedKeys()
  keys?.includes(data.id) ? treeIns.value?.setChecked(node, false, false) : treeIns.value?.setChecked(node, true, false)
}

// 节点过滤
function filterNode(searchKeyword: string, data: any) {
  if (!searchKeyword)
    return true
  return (data[treePorps.label] || '').includes(searchKeyword)
}

// 初始化默认选中项
watch(
  () => props.defaultCheckedKeys,
  (val) => {
    checkedNodes.value = val || []
  },
  { immediate: true },
)

// checkall model
watch(checkedNodes, (_checkedNodes) => {
  // 全不选
  if (!_checkedNodes.length) {
    checkAll.value = false
    isIndeterminate.value = false
  }
  else {
    // 全选
    if (_checkedNodes.filter(item => item[parentKey] === rootValue).length === props.data.length) {
      checkAll.value = true
      isIndeterminate.value = false
    }
    else {
      checkAll.value = false
      isIndeterminate.value = true
    }
  }
})

// filter tree
watchDebounced(
  filterKeyword,
  (val) => {
    treeIns.value?.filter(val)
  },
  { debounce: 300 },
)
</script>

<template>
  <div class="wl-transfer transfer" :style="{ width, height }">
    <div class="transfer-left">
      <h3 class="transfer-title">
        <el-checkbox
          v-model="checkAll"
          :disabled="!data || data.length === 0"
          :indeterminate="isIndeterminate"
          @change="handleCheckAll"
        />
        <span>{{ title }}</span>
        <slot name="title-left" />
      </h3>
      <div class="transfer-main">
        <el-input
          v-if="filter"
          v-model="filterKeyword"
          clearable
          size="small"
          :placeholder="placeholder"
          class="filter-tree"
          maxlength="255"
        />
        <ElTree
          ref="treeIns"
          :filter-node-method="filterNode"
          show-checkbox
          :empty-text="t('app.common.emptyDataTip')"
          default-expand-all
          :data="data"
          :props="treePorps"
          :node-key="nodeKey"
          :default-checked-keys="checkedKeys"
          @node-click="onNodeClick"
          @check-change="onNodeCheckChange"
        >
          <template #default="{ node, data: item }">
            <slot :node="node" :data="item">
              {{ node.label }}
            </slot>
            <!-- </div> -->
          </template>
        </ElTree>
        <slot name="left-footer" />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import '@/styles/ac-transfer.scss';
</style>
