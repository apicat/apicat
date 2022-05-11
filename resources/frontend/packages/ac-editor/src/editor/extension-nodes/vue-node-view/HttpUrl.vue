<template>
    <node-view-wrapper :class="className" tabindex="0">
        <div class="http-url--method" :style="methodTagStyle">{{ attrs.methodConfig.text }}</div>
        <div class="http-url--type" v-if="attrs.bodyDataTypeText">{{ attrs.bodyDataTypeText }}</div>
        <div class="http-url--url" v-if="attrs.url">{{ attrs.url }}</div>
        <div class="http-url--path">{{ attrs.path }}</div>
    </node-view-wrapper>
</template>

<script>
    import { HTTP_METHODS, REQUEST_BODY_DATA_TYPES } from '../../../common/constants'
    import { NodeViewWrapper } from '../../../components/NodeViewWrapper'

    export default {
        name: 'HttpUrl',
        props: ['node', 'view', 'getPos', 'isSelected', 'isReadOnly', 'updateAttributes'],
        components: {
            NodeViewWrapper,
        },
        computed: {
            attrs: function () {
                let attrs = { ...this.node.attrs }
                attrs.methodConfig = HTTP_METHODS.valueOf(attrs.method)
                attrs.bodyDataTypeText = REQUEST_BODY_DATA_TYPES.valueOf(attrs.bodyDataType)
                return attrs
            },
            methodTagStyle() {
                return {
                    background: this.attrs.methodConfig.color,
                }
            },
            className() {
                return [
                    'http-url',
                    {
                        'ProseMirror-selectednode': this.isSelected && !this.isReadOnly,
                    },
                ]
            },
        },
        created() {},
    }
</script>
