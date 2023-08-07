import { typename } from './utils'
import { exampleSchema } from '@/mock/schemas'

export default defineComponent({
  props: {
    modelValue: {},
    definitionSchemas: {
      type: Array,
      default: () => [],
    },
  },
  name: 'SchemaTable',
  setup() {
    const keysMap = new Map()

    const convertSchemaToTree = (parent: any, key: any, name: string, schema: any) => {
      const node: any = {
        key,
        name,
        parent,
        level: parent ? parent.level + 1 : 0,
        schema,
        type: typename(schema.type),
      }

      keysMap.set(key, node)

      if (node.type === 'object') {
        node.children = []
        schema.required = schema.required || []

        const properties = (schema.properties = schema.properties || {})
        const propertiesKeys = schema['x-apicat-order'] || Object.keys(properties)

        for (let k of propertiesKeys) {
          const key = [node.key, 'properties', k].join('.')
          node.children.push(convertSchemaToTree(node, key, k, properties[k]))
        }

        schema['x-apicat-order'] = propertiesKeys
      }

      if (node.type === 'array') {
        schema.items = schema.items
        const key = [node.key, 'items'].join('.')
        node.children = [convertSchemaToTree(node, key, 'items', schema.items)]
      }

      return node
    }

    const root = computed(() => {
      return convertSchemaToTree(undefined, 'root', 'root', exampleSchema)
    })

    return {
      root,
    }
  },

  render() {
    const renderRow = (node: any) => {
      return (
        <div style="text-align:left;">
          {node.key} - {node.level}
          <input v-model={node.name} onInput={(e: any) => (node.name = e.target?.value)} />
          {node.children &&
            node.children.map((child: any) => {
              return <div style="padding-left:20px">{renderRow(child)}</div>
            })}
        </div>
      )
    }

    return (this.root.children || []).map((node: any) => renderRow(node))
  },
})
