<template>
    <node-view-wrapper :class="className" tabindex="0">
        <div class="http-url--url" v-if="attrs.url">{{ attrs.url }}</div>
        <div class="http-url--path">{{ attrs.path }}</div>
    </node-view-wrapper>
</template>

<script>
    /**
     * 该组件与HttpUrl相比，更为广泛，覆盖一些非http协议的文档
     */
    import { NodeViewWrapper } from '../../../components/NodeViewWrapper'

    export default {
        name: 'ApiUrl',
        props: ['node', 'view', 'getPos', 'isSelected', 'isReadOnly', 'updateAttributes'],
        components: {
            NodeViewWrapper,
        },
        computed: {
            attrs: function () {
                return this.node.attrs
            },
            className() {
                return [
                    'http-url api-url',
                    {
                        'ProseMirror-selectednode': this.isSelected && !this.isReadOnly,
                    },
                ]
            },
        },
    }
</script>
