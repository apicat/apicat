import { h, defineComponent } from 'vue'

export const NodeViewWrapper = defineComponent({
    props: {
        as: {
            type: String,
            default: 'div',
        },
    },

    render() {
        return h(
            this.as,
            {
                style: {
                    whiteSpace: 'normal',
                },
                'data-node-view-wrapper': '',
            },
            this.$slots.default?.()
        )
    },
})
