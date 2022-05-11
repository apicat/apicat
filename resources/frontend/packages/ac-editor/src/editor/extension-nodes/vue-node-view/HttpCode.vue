<template>
    <node-view-wrapper :class="className">
        <span class="intro">{{ attrs.intro }}</span>
        <span class="code" ref="code" :style="codeBgColor">{{ attrs.code }}</span>
    </node-view-wrapper>
</template>

<script>
    import { NodeViewWrapper } from '../../../components/NodeViewWrapper'
    import HttpCodeMap, { ColorMap } from '../../../common/HttpCodeMap'
    import tippy from 'tippy.js'

    export default {
        name: 'HttpCode',
        props: ['node', 'view', 'getPos', 'isSelected', 'isReadOnly', 'updateAttributes'],
        components: {
            NodeViewWrapper,
        },
        watch: {
            node: function () {
                this.tippy && this.tippy.setContent(this.attrs.codeDesc || HttpCodeMap[this.attrs.code])
            },
        },
        computed: {
            attrs() {
                return this.node.attrs
            },

            className() {
                return [
                    'http-code',
                    {
                        'ProseMirror-selectednode': this.isSelected && !this.isReadOnly,
                    },
                ]
            },

            codeBgColor() {
                return {
                    background: ColorMap[(this.attrs.code + '')[0]],
                }
            },
        },

        mounted() {
            let dest = HttpCodeMap.find((item) => item.code === this.attrs.code) || null
            if (dest) {
                this.tippy = tippy(this.$refs.code, {
                    // appendTo: () => document.body,
                    interactive: true,
                    theme: 'light',
                    // getReferenceClientRect:()=> this.$refs.code.getBoundingClientRect(),
                    allowHTML: true,
                    content: this.attrs.codeDesc || (dest || { desc: '' }).desc,
                })
            }
        },

        unmounted() {
            this.tippy && this.tippy.destroy()
            this.tippy = null
        },
    }
</script>
