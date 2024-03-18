<script lang="ts">
import { defineComponent, getCurrentInstance, inject, nextTick, provide, ref, watch } from 'vue'
import { CaretRight, Loading } from '@element-plus/icons-vue'
import type { ComponentInternalInstance, PropType } from 'vue'
import { isFunction } from '@apicat/shared'
import { isString } from 'lodash-es'
import NodeContent from './tree-node-content.vue'
import { getNodeKey as getNodeKeyUtil } from './model/util'
import { useNodeExpandEventBroadcast } from './model/useNodeExpandEventBroadcast'
import { dragEventsKey } from './model/useDragNode'
import Node from './model/node'

import type { Nullable } from './utils'
import type { RootTreeType, TreeNodeData, TreeOptionProps } from './tree.type'
import { useNamespace } from '@/hooks/useNamespace'

export default defineComponent({
  name: 'AcTreeNode',
  components: {
    NodeContent,
    Loading,
  },
  props: {
    node: {
      type: Node,
      default: () => ({}),
    },
    props: {
      type: Object as PropType<TreeOptionProps>,
      default: () => ({}),
    },
    accordion: Boolean,
    renderContent: Function,
    renderAfterExpand: Boolean,
    showCheckbox: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['node-expand'],
  setup(props: any, ctx) {
    const namespaceRef = ref('el')
    const ns = useNamespace('tree', namespaceRef)
    const { broadcastExpanded } = useNodeExpandEventBroadcast(props)
    const tree = inject<RootTreeType>('RootTree') as any
    const expanded = ref(false)
    const childNodeRendered = ref(false)
    const oldChecked = ref<boolean>(false)
    const oldIndeterminate = ref<boolean>(false)
    const node$ = ref<Nullable<HTMLElement>>(null)
    const dragEvents = inject(dragEventsKey) as any
    const instance = getCurrentInstance()

    provide('NodeInstance', instance)
    if (!tree) {
      //
    }

    if (props.node.expanded) {
      expanded.value = true
      childNodeRendered.value = true
    }

    const childrenKey = tree.props.children || 'children'

    const getNodeKey = (node: any): any => {
      return getNodeKeyUtil(tree.props.nodeKey, node.data)
    }

    const getNodeClass = (node: Node) => {
      const nodeClassFunc = props.props.class
      if (!nodeClassFunc)
        return {}

      let className
      if (isFunction(nodeClassFunc)) {
        const { data } = node
        className = nodeClassFunc(data, node)
      }
      else {
        className = nodeClassFunc
      }

      if (isString(className))
        return { [className]: true }
      else return className
    }

    const handleSelectChange = (checked: boolean, indeterminate: boolean) => {
      if (
        oldChecked.value !== checked
        || oldIndeterminate.value !== indeterminate
      )
        tree.ctx.emit('check-change', props.node.data, checked, indeterminate)

      oldChecked.value = checked
      oldIndeterminate.value = indeterminate
    }

    const handleClick = (e: MouseEvent) => {
      const store = tree.store.value
      store.setCurrentNode(props.node)
      tree.ctx.emit(
        'current-change',
        store.currentNode ? store.currentNode.data : null,
        store.currentNode,
      )
      tree.currentNode.value = props.node

      if (tree.props.expandOnClickNode)
        handleExpandIconClick()

      if (tree.props.checkOnClickNode && !props.node.disabled) {
        handleCheckChange({
          target: { checked: !props.node.checked },
        })
      }
      tree.ctx.emit('node-click', props.node.data, props.node, instance, e)
    }

    const handleContextMenu = (event: Event) => {
      if (tree.instance.vnode.props.onNodeContextmenu) {
        event.stopPropagation()
        event.preventDefault()
      }
      tree.ctx.emit(
        'node-contextmenu',
        event,
        props.node.data,
        props.node,
        instance,
      )
    }

    const handleExpandIconClick = () => {
      if (props.node.isLeaf)
        return
      if (expanded.value) {
        tree.ctx.emit('node-collapse', props.node.data, props.node, instance)
        props.node.collapse()
      }
      else {
        props.node.expand()
        ctx.emit('node-expand', props.node.data, props.node, instance)
      }
    }

    const handleCheckChange = (ev: any) => {
      props.node.setChecked(ev.target.checked, !tree.props.checkStrictly)
      nextTick(() => {
        const store = tree.store.value
        tree.ctx.emit('check', props.node.data, {
          checkedNodes: store.getCheckedNodes(),
          checkedKeys: store.getCheckedKeys(),
          halfCheckedNodes: store.getHalfCheckedNodes(),
          halfCheckedKeys: store.getHalfCheckedKeys(),
        })
      })
    }

    const handleChildNodeExpand = (
      nodeData: TreeNodeData,
      node: Node,
      instance: ComponentInternalInstance,
    ) => {
      broadcastExpanded(node)
      tree.ctx.emit('node-expand', nodeData, node, instance)
    }

    const handleDragStart = (event: DragEvent) => {
      if (!tree.props.draggable)
        return
      dragEvents.treeNodeDragStart({ event, treeNode: props })
    }

    const handleDragOver = (event: DragEvent) => {
      event.preventDefault()
      if (!tree.props.draggable)
        return
      dragEvents.treeNodeDragOver({
        event,
        treeNode: { $el: node$.value, node: props.node },
      })
    }

    const handleDrop = (event: DragEvent) => {
      event.preventDefault()
    }

    const handleDragEnd = (event: DragEvent) => {
      if (!tree.props.draggable)
        return
      dragEvents.treeNodeDragEnd(event)
    }

    watch(
      () => {
        const children = props.node.data[childrenKey]
        return children && [...children]
      },
      () => {
        props.node.updateChildren()
      },
    )

    watch(
      () => props.node.indeterminate,
      (val) => {
        handleSelectChange(props.node.checked, val)
      },
    )

    watch(
      () => props.node.checked,
      (val) => {
        handleSelectChange(val, props.node.indeterminate)
      },
    )

    watch(
      () => props.node.expanded,
      (val) => {
        nextTick(() => (expanded.value = val))
        if (val)
          childNodeRendered.value = true
      },
    )

    return {
      ns,
      node$,
      tree,
      expanded,
      childNodeRendered,
      oldChecked,
      oldIndeterminate,
      getNodeKey,
      getNodeClass,
      handleSelectChange,
      handleClick,
      handleContextMenu,
      handleExpandIconClick,
      handleCheckChange,
      handleChildNodeExpand,
      handleDragStart,
      handleDragOver,
      handleDrop,
      handleDragEnd,
      CaretRight,
    }
  },
})
</script>

<template>
  <div
    v-show="node.visible"
    ref="node$"
    :class="[
      ns.b('node'),
      ns.is('expanded', expanded),
      ns.is('current', node.isCurrent),
      ns.is('hidden', !node.visible),
      ns.is('focusable', !node.disabled),
      ns.is('checked', !node.disabled && node.checked), getNodeClass(node),
    ]"
    role="treeitem"
    tabindex="-1"
    :aria-expanded="expanded"
    :aria-disabled="node.disabled"
    :aria-checked="node.checked"
    :draggable="tree.props.draggable"
    :data-key="getNodeKey(node)"
    @click.stop="handleClick"
    @contextmenu="handleContextMenu"
    @dragstart.stop="handleDragStart"
    @dragover.stop="handleDragOver"
    @dragend.stop="handleDragEnd"
    @drop.stop="handleDrop"
  >
    <div
      :class="ns.be('node', 'content')"
      :style="{ paddingLeft: `${(node.level - 1) * tree.props.indent}px` }"
    >
      <el-icon
        v-if="tree.props.icon || CaretRight"
        :class="[
          ns.be('node', 'expand-icon'),
          ns.is('leaf', node.isLeaf),
          {
            expanded: !node.isLeaf && expanded,
          },
        ]"
        @click.stop="handleExpandIconClick"
      >
        <component :is="tree.props.icon || CaretRight" />
      </el-icon>
      <el-checkbox
        v-if="showCheckbox"
        :model-value="node.checked"
        :indeterminate="node.indeterminate"
        :disabled="!!node.disabled"
        @click.stop
        @change="handleCheckChange($event)"
      />
      <el-icon
        v-if="node.loading"
        :class="[ns.be('node', 'loading-icon'), ns.is('loading')]"
      >
        <Loading />
      </el-icon>
      <NodeContent :node="node" :render-content="renderContent" />
    </div>

    <div
      v-if="!renderAfterExpand || childNodeRendered"
      v-show="expanded"
      :class="ns.be('node', 'children')"
      role="group"
      :aria-expanded="expanded"
    >
      <ac-tree-node
        v-for="child in node.childNodes"
        :key="getNodeKey(child)"
        :render-content="renderContent"
        :render-after-expand="renderAfterExpand"
        :show-checkbox="showCheckbox"
        :node="child"
        :accordion="accordion"
        :props="props"
        @node-expand="handleChildNodeExpand"
      />
    </div>
  </div>
</template>
