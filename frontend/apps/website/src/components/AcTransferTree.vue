<template>
  <div class="wl-transfer transfer" :style="{ width, height }">
    <!-- 左侧穿梭框 原料框 -->
    <div class="transfer-left">
      <h3 class="transfer-title">
        <el-checkbox :indeterminate="from_is_indeterminate" v-model="from_check_all" @change="fromAllBoxChange"></el-checkbox>
        <span>{{ fromTitle }}</span>
        <slot name="title-left"></slot>
      </h3>
      <!-- 内容区 -->
      <div class="transfer-main">
        <el-input v-if="filter" clearable size="small" :placeholder="placeholder" v-model="filterFrom" class="filter-tree" />
        <el-tree
          ref="from-tree"
          show-checkbox
          empty-text="暂无数据"
          :lazy="lazy"
          :indent="indent"
          :node-key="node_key"
          :load="leftloadNode"
          :data="self_from_data"
          :accordion="accordion"
          :props="selfDefaultProps"
          :default-expand-all="openAll"
          :highlight-current="highLight"
          :check-strictly="checkStrictly"
          :filter-node-method="filterNodeFrom"
          :check-on-click-node="checkOnClickNode"
          :render-after-expand="renderAfterExpand"
          :expand-on-click-node="expandOnClickNode"
          :default-checked-keys="defaultCheckedKeys"
          :default-expanded-keys="from_expanded_keys"
          @node-expand="(data, node) => onNodeExpand(data, node, 'from')"
          @node-collapse="(data, node) => onNodeCollapse(data, node, 'from')"
          @check="fromTreeChecked"
        >
          <template #default="{ node, data }">
            <div class="block truncate" :title="node.label">
              <slot name="content-left" :node="node" :data="data">
                {{ node.label }}
              </slot>
            </div>
          </template>
        </el-tree>
        <slot name="left-footer"></slot>
      </div>
    </div>
    <!-- 穿梭区 按钮框 -->
    <div class="transfer-center">
      <template v-if="button_text">
        <p class="transfer-center-item">
          <el-button type="primary" @click="addToAims(true)" :disabled="from_disabled" :icon="DArrowRight" />
        </p>
        <p class="transfer-center-item">
          <el-button type="primary" @click="removeToSource" :disabled="to_disabled" :icon="DArrowLeft" />
        </p>
      </template>
      <template v-else>
        <p class="transfer-center-item">
          <el-button type="primary" @click="addToAims(true)" :disabled="from_disabled" circle />
        </p>
        <p class="transfer-center-item">
          <el-button type="primary" @click="removeToSource" :disabled="to_disabled" icon="el-icon-arrow-left" circle />
        </p>
      </template>
    </div>
    <!-- 右侧穿梭框 目标框 -->
    <div class="transfer-right">
      <h3 class="transfer-title">
        <el-checkbox :indeterminate="to_is_indeterminate" v-model="to_check_all" @change="toAllBoxChange"></el-checkbox>
        <span>{{ toTitle }}</span>
        <slot name="title-right"></slot>
      </h3>
      <!-- 内容区 -->
      <div class="transfer-main">
        <el-input v-if="filter" clearable size="small" v-model="filterTo" :placeholder="placeholder" class="filter-tree" />
        <el-tree
          ref="to-tree"
          show-checkbox
          empty-text="暂无数据"
          :indent="indent"
          :lazy="lazyRight"
          :data="self_to_data"
          :node-key="node_key"
          :load="rightloadNode"
          :props="selfDefaultProps"
          :default-expand-all="openAll"
          :highlight-current="highLight"
          :check-strictly="checkStrictly"
          :filter-node-method="filterNodeTo"
          :check-on-click-node="checkOnClickNode"
          :render-after-expand="renderAfterExpand"
          :expand-on-click-node="expandOnClickNode"
          :default-expanded-keys="to_expanded_keys"
          @node-expand="(data, node) => onNodeExpand(data, node, 'to')"
          @node-collapse="(data, node) => onNodeCollapse(data, node, 'to')"
          @check="toTreeChecked"
        >
          <template #default="{ node, data }">
            <div class="block truncate" :title="node.label">
              <slot name="content-left" :node="node" :data="data">
                {{ node.label }}
              </slot>
            </div>
          </template>
        </el-tree>
        <slot name="right-footer"></slot>
      </div>
    </div>
  </div>
</template>

<script>
import { differenceBy } from 'lodash-es'
import DArrowRight from '~icons/ep/d-arrow-right'
import DArrowLeft from '~icons/ep/d-arrow-left'
import { findParents, flattenDeep } from '@/commons/helper'
import { markRaw, toRaw } from 'vue'

export default {
  name: 'AcTransferTree',
  props: {
    // 宽度
    width: {
      type: String,
      default: '100%',
    },
    // 高度
    height: {
      type: String,
      default: '320px',
    },
    // 标题
    title: {
      type: Array,
      default: () => ['源列表', '目标列表'],
    },
    // 穿梭按钮名字
    button_text: {
      type: Array,
      default: () => ['添加', '移除'],
    },
    // 源数据
    from_data: {
      type: Array,
      default: () => [],
    },
    // 选中数据
    to_data: {
      type: Array,
      default: () => [],
    },
    // el-tree 配置项
    defaultProps: Object,
    // el-tree node-key 必须唯一
    node_key: {
      type: String,
      default: 'id',
    },
    // 自定义 pid参数名
    pid: {
      type: String,
      default: 'pid',
    },
    // 自定义根节点pid的值，用于结束递归
    rootPidValue: {
      type: [String, Number],
      default: 0,
    },
    // 是否启用筛选
    filter: {
      type: Boolean,
      default: false,
    },
    // 是否展开所有节点
    openAll: {
      type: Boolean,
      default: false,
    },
    // 穿梭后是否展开节点
    transferOpenNode: {
      type: Boolean,
      default: true,
    },
    // 源数据 默认选中节点
    defaultCheckedKeys: {
      type: Array,
      default: () => [],
    },
    // 源数据 默认展开节点
    defaultExpandedKeys: {
      type: Array,
      default: () => [],
    },
    // 筛选placeholder
    placeholder: {
      type: String,
      default: '输入关键字进行过滤',
    },
    // 自定义筛选函数
    filterNode: Function,
    // 默认穿梭一次默认选中数据
    defaultTransfer: {
      type: Boolean,
      default: false,
    },
    // 是否启用懒加载
    lazy: {
      type: Boolean,
      default: false,
    },
    // 是否右侧树也启用懒加载
    lazyRight: {
      type: Boolean,
      default: false,
    },
    // 懒加载的回调函数
    lazyFn: Function,
    // 是否高亮当前选中节点，默认值是 false。
    highLight: {
      type: Boolean,
      default: false,
    },
    // 是否遵循父子不关联
    checkStrictly: {
      type: Boolean,
      default: false,
    },
    // 父子不关联模式
    checkStrictlyType: {
      type: String,
      default: 'authorization',
      validator: function (value) {
        /**
         * @name 父子不关联的三种模式，第一种适合业务授权场景，后两种不存在快速选中需要手选
         * @param authorization授权模式：左侧选择子节点自动带着父节点；右侧选择父节点自动带着子节点；此模式两侧可能存在相同的非叶子节点
         * @param puppet木偶模式：纯父子不关联穿梭，但要保持完整的树形结构，只自动带上穿梭到对面拼接所需的骨架结构；此模式两侧可能存在相同的非叶子节点
         * @param modular积木模式：纯父子不关联穿梭，也不保持完整的树形结构，像积木一样右侧要形成树形则需要把左侧拆除，左侧拆的越多右侧形成的树结构越完整；此模式左右两侧保证严格的唯一性
         */
        return ['authorization', 'puppet', 'modular'].indexOf(value) !== -1
      },
    },
    // 是否每次只打开一个同级树节点
    accordion: {
      type: Boolean,
      default: false,
    },
    // 是否在第一次展开某个树节点后才渲染其子节点
    renderAfterExpand: {
      type: Boolean,
      default: true,
    },
    // 是否在点击节点的时候展开或者收缩节点
    expandOnClickNode: {
      type: Boolean,
      default: true,
    },
    // 是否在点击节点的时候选中节点
    checkOnClickNode: {
      type: Boolean,
      default: false,
    },
    // 相邻级节点间的水平缩进，单位为像素
    indent: {
      type: Number,
      default: 16,
    },
  },
  data() {
    return {
      allExpandedKeys: [],
      DArrowRight: markRaw(DArrowRight),
      DArrowLeft: markRaw(DArrowLeft),
      from_is_indeterminate: false, // 源数据是否半选
      from_check_all: false, // 源数据是否全选
      from_expanded_keys: [], // 源数据展开节点
      from_disabled: true, // 添加按钮是否禁用
      from_check_keys: [], // 源数据选中key数组 以此属性关联穿梭按钮，总全选、半选状态
      from_array_clone: [], // 左侧数据一维化后存储为json格式
      to_check_all: false, // 目标数据是否全选
      to_is_indeterminate: false, // 目标数据是否半选
      to_expanded_keys: [], // 目标数据展开节点
      to_disabled: true, // 移除按钮是否禁用
      to_check_keys: [], // 目标数据选中key数组 以此属性关联穿梭按钮，总全选、半选状态
      to_array_clone: [], // 右侧数据一维化后存储为json格式
      filterFrom: '', // 源数据筛选
      filterTo: '', // 目标数据筛选
      strictly_parents: [], // 当使用父子不关联时，将左侧数据向右侧移动时，为了保证在右侧能形成树结构，必须将父节点也移动
      strictly_transferred: [], // 父子不关联时已经穿梭过的节点记录，用于第一次拼接父节点穿梭后，其他子节点不再拼接父节点
    }
  },
  computed: {
    // 左侧数据
    self_from_data() {
      return this.from_data
    },
    // 右侧数据
    self_to_data() {
      return this.to_data
    },
    // 左侧菜单名
    fromTitle() {
      let [text] = this.title
      return text
    },
    // 右侧菜单名
    toTitle() {
      let [, text] = this.title
      return text
    },
    // 配置项
    selfDefaultProps() {
      return {
        label: 'label',
        children: 'children',
        ...this.defaultProps,
      }
    },
  },
  watch: {
    // 左侧 状态监测
    from_check_keys(val) {
      if (val.length > 0) {
        // 穿梭按钮是否禁用
        this.from_disabled = false
        // 总半选是否开启
        this.from_is_indeterminate = true
        // 总全选是否开启 - 根据选中节点中为根节点的数量是否和源数据长度相等
        let allCheck = false
        if (!this.checkStrictly) {
          const roots = val.filter((item) => item[this.pid] === this.rootPidValue)
          allCheck = roots.length === this.self_from_data.length
        } else {
          allCheck = val.length === this.from_array_clone.length
        }
        if (allCheck) {
          // 关闭半选 开启全选
          this.from_is_indeterminate = false
          this.from_check_all = true
        } else {
          this.from_is_indeterminate = true
          this.from_check_all = false
        }
      } else {
        this.from_disabled = true
        this.from_is_indeterminate = false
        this.from_check_all = false
      }
    },
    // 右侧 状态监测
    to_check_keys(val) {
      if (val.length > 0) {
        // 穿梭按钮是否禁用
        this.to_disabled = false
        // 总半选是否开启
        this.to_is_indeterminate = true
        // 总全选是否开启 - 根据选中节点中为根节点的数量是否和源数据长度相等
        let allCheck = false
        if (!this.checkStrictly) {
          const roots = val.filter((item) => item[this.pid] === this.rootPidValue)
          allCheck = roots.length === this.self_to_data.length
        } else {
          allCheck = val.length === this.to_array_clone.length
        }
        if (allCheck) {
          // 关闭半选 开启全选
          this.to_is_indeterminate = false
          this.to_check_all = true
        } else {
          this.to_is_indeterminate = true
          this.to_check_all = false
        }
      } else {
        this.to_disabled = true
        this.to_is_indeterminate = false
        this.to_check_all = false
      }
    },
    // 左侧 数据筛选
    filterFrom(val) {
      this.$refs['from-tree'].filter(val)
    },
    // 右侧 数据筛选
    filterTo(val) {
      this.$refs['to-tree'].filter(val)
    },
    // 监视默认选中
    defaultCheckedKeys: {
      handler(val) {
        this.from_check_keys = val || []
        if (this.defaultTransfer && this.from_check_keys.length) {
          this.$nextTick(() => {
            this.addToAims(false)
          })
        }
      },
      immediate: true,
    },
    // 监视默认展开
    defaultExpandedKeys: {
      handler(val) {
        let _form = new Set(this.from_expanded_keys.concat(val))
        this.from_expanded_keys = [..._form]
        let _to = new Set(this.to_expanded_keys.concat(val))
        this.to_expanded_keys = [..._to]
      },
      immediate: true,
    },
  },
  methods: {
    onNodeExpand(data, node, dist) {
      let source = dist === 'from' ? this.from_expanded_keys : this.to_expanded_keys
      const expandedKeys = new Set(source.concat(data.id))
      this[dist === 'from' ? 'from_expanded_keys' : 'to_expanded_keys'] = [...expandedKeys]
    },

    onNodeCollapse(data, node, dist) {
      const store = node.store
      const removeChildId = [data.id]

      // 点击父级时，移除子集展开项
      let source = dist === 'from' ? this.from_expanded_keys : this.to_expanded_keys
      ;(source || []).forEach((id) => {
        let cNode = store.getNode(id)
        if (node.contains(cNode)) {
          removeChildId.push(id)
        }
      })
      const expandedKeys = source.filter((id) => removeChildId.indexOf(id) === -1)
      this[dist === 'from' ? 'from_expanded_keys' : 'to_expanded_keys'] = [...expandedKeys]
    },

    // 添加按钮
    addToAims(emit) {
      // 获取选中通过穿梭框的keys - 仅用于传送纯净的id数组到父组件同后台通信
      let keys = this.$refs['from-tree'].getCheckedKeys()
      // 获取半选通过穿梭框的keys - 仅用于传送纯净的id数组到父组件同后台通信
      let harfKeys = this.$refs['from-tree'].getHalfCheckedKeys()
      // 选中节点数据
      let arrayCheckedNodes = this.$refs['from-tree'].getCheckedNodes()
      // 半选中节点数据
      let arrayHalfCheckedNodes = this.$refs['from-tree'].getHalfCheckedNodes()
      // 自定义参数读取设置
      let children__ = this.selfDefaultProps.children || 'children'
      let pid__ = this.pid || 'pid'
      let id__ = this['node_key'] || 'id'
      let root__ = this.rootPidValue || 0
      // 将目标侧数据拉平为一维数组，查询速度比每个节点去目标侧递归查询快
      this.to_array_clone = flattenDeep(this.self_to_data, this.selfDefaultProps.children)

      // 第一步：排除在对面已经存在的半选节点，然后将需穿梭半选节点的children设置为[]并穿梭;
      arrayHalfCheckedNodes.forEach((i) => {
        let harfInTarget = this.to_array_clone.some((t) => t[id__] === i[id__])
        if (harfInTarget) return
        let _parent = root__ !== i[pid__] ? i[pid__] : null
        this.$refs['to-tree'].append(
          Object.assign({}, i, {
            [children__]: [],
          }),
          _parent
        )
      })
      // 第二步：先将对面存在的节点抛弃
      let notInTargetNodes = differenceBy(arrayCheckedNodes, this.to_array_clone, id__)
      // 第三步：若a节点的父节点也在选中节点中，则将a节点也抛弃，最后将剩余的节点穿梭
      notInTargetNodes.forEach((i) => {
        let parentInHere = notInTargetNodes.some((t) => t[id__] === i[pid__])
        if (parentInHere) return
        let _parent = root__ !== i[pid__] ? i[pid__] : null
        this.$refs['to-tree'].append(i, _parent)
      })

      // 左侧删掉选中数据
      arrayCheckedNodes.map((item) => this.$refs['from-tree'].remove(item))

      // 处理完毕按钮恢复禁用状态
      this.from_check_keys = []
      // 清空对面选中
      this.$refs['to-tree'].setCheckedKeys([])
      this.to_check_all = false
      this.to_is_indeterminate = false

      // 传递信息给父组件
      const all_move_nodes = [...arrayCheckedNodes, ...this.strictly_parents]

      emit &&
        this.$emit('add-btn', this.self_from_data, this.self_to_data, {
          keys,
          nodes: all_move_nodes,
          harfKeys,
          halfNodes: arrayHalfCheckedNodes,
        })

      // 处理完毕取消选中
      this.$refs['from-tree'].setCheckedKeys([])
    },
    // 移除按钮
    removeToSource() {
      // 获取选中通过穿梭框的keys - 仅用于传送纯净的id数组到父组件同后台通信
      let keys = this.$refs['to-tree'].getCheckedKeys()
      // 获取半选通过穿梭框的keys - 仅用于传送纯净的id数组到父组件同后台通信
      let harfKeys = this.$refs['to-tree'].getHalfCheckedKeys()
      // 获取选中通过穿梭框的nodes 选中节点数据
      let arrayCheckedNodes = this.$refs['to-tree'].getCheckedNodes()
      // 半选中节点数据
      let arrayHalfCheckedNodes = this.$refs['to-tree'].getHalfCheckedNodes()
      // 自定义参数读取设置
      let children__ = this.selfDefaultProps.children || 'children'
      let pid__ = this.pid || 'pid'
      let id__ = this['node_key'] || 'id'
      let root__ = this.rootPidValue || 0
      // 将目标侧数据拉平为一维数组，查询速度比每个节点去目标侧递归查询快
      this.from_array_clone = flattenDeep(this.self_from_data, this.selfDefaultProps.children)

      // 第一步：排除在对面已经存在的半选节点，然后将需穿梭半选节点的children设置为[]并穿梭;
      arrayHalfCheckedNodes.forEach((i) => {
        let harfInTarget = this.from_array_clone.some((t) => t[id__] === i[id__])
        if (harfInTarget) return
        let _parent = root__ !== i[pid__] ? i[pid__] : null
        this.$refs['from-tree'].append(
          Object.assign({}, i, {
            [children__]: [],
          }),
          _parent
        )
      })
      // 第二步：先将对面存在的节点抛弃
      let notInTargetNodes = differenceBy(arrayCheckedNodes, this.from_array_clone, id__)
      // 第三步：若a节点的父节点也在选中节点中，则将a节点也抛弃，最后将剩余的节点穿梭
      notInTargetNodes.forEach((i) => {
        let parentInHere = notInTargetNodes.some((t) => t[id__] === i[pid__])
        if (parentInHere) return
        let _parent = root__ !== i[pid__] ? i[pid__] : null
        this.$refs['from-tree'].append(i, _parent)
      })

      // 右侧删掉选中数据
      arrayCheckedNodes.map((item) => this.$refs['to-tree'].remove(item))

      // 处理完毕按钮恢复禁用状态
      this.to_check_keys = []
      // 清空对面选中
      this.$refs['from-tree'].setCheckedKeys([])
      this.from_check_all = false
      this.from_is_indeterminate = false

      // 传递信息给父组件
      this.$emit('remove-btn', this.self_from_data, this.self_to_data, {
        keys,
        nodes: arrayCheckedNodes,
        harfKeys,
        halfNodes: arrayHalfCheckedNodes,
      })
      // 处理完毕取消选中
      this.$refs['to-tree'].setCheckedKeys([])
    },

    // 异步加载左侧
    leftloadNode(node, resolve) {
      this.lazyFn && this.lazyFn(node, resolve, 'left')
    },
    // 异步加载右侧
    rightloadNode(node, resolve) {
      this.lazyFn && this.lazyFn(node, resolve, 'right')
    },
    // 源树选中事件 - 是否禁用穿梭按钮
    fromTreeChecked(nodeObj, treeObj) {
      this.from_check_keys = treeObj.checkedNodes
      // 父子不关联-授权模式
      if (this.checkStrictly && this.checkStrictlyType === 'authorization' && this.from_check_keys.some((i) => i[this.node_key] === nodeObj[this.node_key])) {
        this.authorizationAutoCheckLeft(nodeObj)
      }
      this.$nextTick(() => {
        this.$emit('left-check-change', nodeObj, treeObj, this.from_check_all)
      })
    },
    // 父子不关联-授权模式-左侧选中子节点自动选中父节点
    authorizationAutoCheckLeft(node) {
      // 查询所有父节点
      const parents = findParents(node, this.from_array_clone, {
        id: this.node_key,
        parentId: this.pid,
        // children: this.selfDefaultProps.children,
        root: this.rootPidValue,
      })
      if (!parents.length) return
      // 过滤掉已经选中过的父节点
      const autoAddParents = differenceBy(parents, this.from_check_keys, this.node_key)
      autoAddParents.forEach((i) => {
        this.$refs['from-tree'].setChecked(i, true)
        this.from_check_keys.push(i)
      })
    },
    // 目标树选中事件 - 是否禁用穿梭按钮
    toTreeChecked(nodeObj, treeObj) {
      this.to_check_keys = treeObj.checkedNodes
      // 父子不关联-授权模式
      if (this.checkStrictly && this.checkStrictlyType === 'authorization' && this.to_check_keys.some((i) => i[this.node_key] === nodeObj[this.node_key])) {
        this.authorizationAutoCheckRight(nodeObj)
      }
      this.$nextTick(() => {
        this.$emit('right-check-change', nodeObj, treeObj, this.to_check_all)
      })
    },
    // 父子不关联-授权模式-右侧选中父节点自动选中子节点
    authorizationAutoCheckRight(nodeObj) {
      // 查询所有子节点
      const children = flattenDeep([nodeObj], this.selfDefaultProps.children)
      if (!children.length) return
      // 过滤掉已经选中过的子节点
      const autoAddChildren = differenceBy(children, this.to_check_keys, this.node_key)
      autoAddChildren.forEach((i) => {
        this.to_check_keys.push(i)
        this.$refs['to-tree'].setChecked(i, true)
      })
    },
    // 源数据 总全选checkbox
    fromAllBoxChange(val) {
      if (!this.self_from_data.length) return
      if (val) {
        this.from_check_keys = this.checkStrictly ? flattenDeep(this.self_from_data, this.selfDefaultProps.children) : this.self_from_data
        this.$refs['from-tree'].setCheckedNodes(this.from_check_keys)
      } else {
        this.$refs['from-tree'].setCheckedNodes([])
        this.from_check_keys = []
      }
      this.$emit('left-check-change', null, null, this.from_check_all)
    },
    // 目标数据 总全选checkbox
    toAllBoxChange(val) {
      if (!this.self_to_data.length) return
      if (val) {
        this.to_check_keys = this.checkStrictly ? flattenDeep(this.self_to_data, this.selfDefaultProps.children) : this.self_to_data
        this.$refs['to-tree'].setCheckedNodes(this.to_check_keys)
      } else {
        this.$refs['to-tree'].setCheckedNodes([])
        this.to_check_keys = []
      }
      this.$emit('right-check-change', null, null, this.to_check_all)
    },
    // 源数据 筛选
    filterNodeFrom(value, data) {
      if (this.filterNode) {
        return this.filterNode(value, data, 'form')
      }
      if (!value) return true
      return data[this.selfDefaultProps.label].indexOf(value) !== -1
    },
    // 目标数据筛选
    filterNodeTo(value, data) {
      if (this.filterNode) {
        return this.filterNode(value, data, 'to')
      }
      if (!value) return true
      return data[this.selfDefaultProps.label].indexOf(value) !== -1
    },
    // 以下为提供方法 ----------------------------------------------------------------方法--------------------------------------
    /**
     * @name 清空选中节点
     * @param {String} type left左边 right右边 all全部 默认all
     */
    clearChecked(type = 'all') {
      if (type === 'left') {
        this.$refs['from-tree'].setCheckedKeys([])
        this.from_is_indeterminate = false
        this.from_check_all = false
      } else if (type === 'right') {
        this.$refs['to-tree'].setCheckedKeys([])
        this.to_is_indeterminate = false
        this.to_check_all = false
      } else {
        this.$refs['from-tree'].setCheckedKeys([])
        this.$refs['to-tree'].setCheckedKeys([])
        this.from_is_indeterminate = false
        this.from_check_all = false
        this.to_is_indeterminate = false
        this.to_check_all = false
      }
    },

    /**
     * 清空展开项
     */
    clearExpand() {
      this.from_expanded_keys = []
      this.to_expanded_keys = []
    },

    /**
     * @name 获取选中数据
     */
    getChecked() {
      // 左侧选中信息
      let leftKeys = this.$refs['from-tree'].getCheckedKeys()
      let leftHarfKeys = this.$refs['from-tree'].getHalfCheckedKeys()
      let leftNodes = this.$refs['from-tree'].getCheckedNodes()
      let leftHalfNodes = this.$refs['from-tree'].getHalfCheckedNodes()
      // 右侧选中信息
      let rightKeys = this.$refs['to-tree'].getCheckedKeys()
      let rightHarfKeys = this.$refs['to-tree'].getHalfCheckedKeys()
      let rightNodes = this.$refs['to-tree'].getCheckedNodes()
      let rightHalfNodes = this.$refs['to-tree'].getHalfCheckedNodes()
      return {
        leftKeys,
        leftHarfKeys,
        leftNodes,
        leftHalfNodes,
        rightKeys,
        rightHarfKeys,
        rightNodes,
        rightHalfNodes,
      }
    },
    /**
     * @name 设置选中数据
     * @param {Array} leftKeys 左侧ids
     * @param {Array} rightKeys 右侧ids
     */
    setChecked(leftKeys = [], rightKeys = []) {
      this.$refs['from-tree'].setCheckedKeys(leftKeys)
      this.$refs['to-tree'].setCheckedKeys(rightKeys)
    },
    /**
     * @name 清除搜索条件
     * @param {String} type left左边 right右边 all全部 默认all
     */
    clearFilter(type = 'all') {
      if (type === 'left') {
        this.filterFrom = ''
      } else if (type === 'right') {
        this.filterTo = ''
      } else {
        this.filterFrom = ''
        this.filterTo = ''
      }
    },

    /**
     * 返回一维数据，便于获取值
     * @returns {*}
     */
    getFlattenValues() {
      return flattenDeep(this.self_to_data, this.selfDefaultProps.children).map((item) => toRaw(item))
    },
  },
}
</script>
<style lang="scss">
@import '@/styles/ac-transfer.scss';
</style>
