<script lang="ts">
import { defineComponent, h, inject } from 'vue'

import type { ComponentInternalInstance } from 'vue'
import type { RootTreeType } from './tree.type'
import { useNamespace } from '@/hooks/useNamespace'

export default defineComponent({
  name: 'AcTreeNodeContent',
  props: {
    node: {
      type: Object,
      required: true,
    },
    renderContent: Function,
  },
  setup(props) {
    const namespaceRef = ref('el')
    const ns = useNamespace('tree', namespaceRef)
    const nodeInstance = inject<ComponentInternalInstance>('NodeInstance')
    const tree = inject<RootTreeType>('RootTree') as any
    return () => {
      const node = props.node
      const { data, store } = node
      return props.renderContent
        ? props.renderContent(h, { _self: nodeInstance, node, data, store })
        : tree.ctx.slots.default
          ? tree.ctx.slots.default({ node, data })
          : h('span', { class: ns.be('node', 'label') }, [node.label])
    }
  },
})
</script>
