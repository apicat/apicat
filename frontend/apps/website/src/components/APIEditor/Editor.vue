<template>
  <div class="ac-sce">
    <EditorRow :level="1" :data="root" :readonly="readonly"></EditorRow>
  </div>
</template>

<script setup lang="ts">
import { computed, provide, ref, unref, watch } from 'vue'
import EditorRow from './EditorRow.vue'
import type { JSONSchema, Definition, Tree } from './types'
import { constNodeType, typename } from './types'

const props = withDefaults(
  defineProps<{
    // 根schema
    modelValue: JSONSchema
    // 引用模型的集合
    definitions?: Definition[]
    readonly?: boolean
  }>(),
  {
    definitions: () => [],
  }
)

const emits = defineEmits(['update:modelValue'])
const expandKeys = ref<Set<string>>(new Set([constNodeType.root]))
const localSchema = ref(props.modelValue)
watch(
  () => props.modelValue,
  () => {
    localSchema.value = props.modelValue
  }
)

const root = computed(() => {
  return convertTreeData(undefined, constNodeType.root, constNodeType.root, localSchema.value)
})

provide('expandKeys', expandKeys.value)
provide('definitions', () => props.definitions)
provide('change', changeEvent)

function changeEvent(root?: JSONSchema) {
  if (root) {
    emits('update:modelValue', root)
  } else {
    emits('update:modelValue', localSchema.value)
  }
}

function convertTreeData(parent: Tree | undefined, key: string, label: string, schema: JSONSchema): Tree {
  const item: Tree = {
    key,
    label,
    schema,
    parent,
    type: '',
  }
  if (schema.$ref != undefined) {
    const name = schema.$ref.match(/#\/definitions\/(.*)/)?.[1]
    const refschema = props.definitions?.find((v) => v.name === name)
    if (refschema && refschema.schema) {
      item.refObj = refschema
      schema = refschema.schema
    }
  }
  item.type = typename(schema.type)
  switch (item.type) {
    case 'object':
      let children: Tree[] = []
      const ps = schema.properties
      if (ps) {
        const orders = schema['x-apicat-orders'] || Object.keys(ps)
        for (let k of orders) {
          if (k != 'additionalMetadata') {
            const p = key + '.' + k
            children.push(convertTreeData(item, p, k, ps[k]))
          }
        }
        schema['x-apicat-orders'] = orders
      } else {
        schema.properties = {}
      }
      if (!schema.required) schema.required = []
      if (!schema['x-apicat-orders']) schema['x-apicat-orders'] = []
      item.children = children
      break
    case 'array':
      if (schema.items) {
        item.children = [convertTreeData(item, `${key}.${constNodeType.items}`, constNodeType.items, schema.items as JSONSchema)]
      }
  }

  // default expand children
  if (item.children && item.children.length) {
    expandKeys.value.add(key)
  }

  return item
}

// 处理拖拽
provide('drop', dropHandler)
function dropHandler(offset: number, to: Tree, source: string) {
  const p = to.parent
  if (!p) {
    return
  }
  const container = to.parent?.schema
  if (container && checkValidObject(container)) {
    const from = findTreeFromKey(root.value, source)
    if (from) {
      const orders = container['x-apicat-orders'] || []
      if (from.parent?.key === to.parent?.key) {
        // 同一目录下 只需要变换位置
        const s = orders.filter((v) => v !== from.label)
        let i = s.indexOf(to.label)
        if (offset > 0) i += 1
        s.splice(i < 0 ? 0 : i, 0, from.label)
        container['x-apicat-orders'] = s
      } else {
        // 添加新的同时删除老的
        if (container.properties) {
          container.properties[from.label] = getSchemaSource(from)
          let i = orders.indexOf(to.label)
          if (offset > 0) i += 1
          orders.splice(i < 0 ? 0 : i, 0, from.label)
          const p = from.parent?.schema
          if (p) {
            p['x-apicat-orders'] = p['x-apicat-orders']?.filter((v) => v != from.label)
            delete p.properties?.[from.label]
          }
        }
      }
    }
    changeEvent()
  }
}

function findTreeFromKey(n: Tree, k: string): Tree | undefined {
  if (n.children) {
    for (let v of n.children) {
      if (v.key === k) {
        return v
      }
      const n = findTreeFromKey(v, k)
      if (n) {
        return n
      }
    }
  }
  return undefined
}

function getSchemaSource(t: Tree) {
  if (t.refObj) {
    return {
      $ref: `#/definitions/${t.refObj.name}`,
      description: t.schema.description,
    }
  }
  return t.schema
}

function checkValidObject(schema: JSONSchema) {
  return schema.type === 'object' && schema.properties && schema.required && schema['x-apicat-orders']
}
</script>

<style>
.ac-sce {
  background: var(--el-fill-color-blank);
  color: var(--el-tree-text-color);
  font-size: var(--el-font-size-base);
  border: 1px var(--el-border-color-lighter) solid;
  padding: 4px 0px;
  border-radius: var(--el-border-radius-base);
}

.ac-sce:focus-within {
  box-shadow: rgba(0, 0, 0, 0.1) 0 2px 4px;
}
</style>
