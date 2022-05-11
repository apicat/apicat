<template>
    <node-view-wrapper :class="className" data-language="json">
        <div contenteditable="false">
            <select @change="handleLanguageChange"></select>
        </div>
        <pre>
      <code spellcheck="false"></code>
    </pre>
    </node-view-wrapper>
</template>

<script>
    import { NodeViewWrapper } from '../../../components/NodeViewWrapper'

    export default {
        name: 'CodeBlock',
        props: ['node', 'view', 'getPos', 'isSelected', 'updateAttributes', 'options'],
        components: {
            NodeViewWrapper,
        },
        computed: {
            attrs() {
                return this.node.attrs
            },

            className() {
                return [
                    'code-block',
                    {
                        'ProseMirror-selectednode': this.isSelected,
                    },
                ]
            },
        },

        methods: {
            handleLanguageChange(event) {
                const { state } = this.view
                const { tr } = state
                const element = event.target
                const { top, left } = element.getBoundingClientRect()
                const result = this.view.posAtCoords({ top, left })

                if (result) {
                    const transaction = tr.setNodeMarkup(result.inside, undefined, {
                        language: element.value,
                    })
                    this.view.dispatch(transaction)
                }
            },
        },

        created() {
            console.log(this.options)
        },
        unmounted() {},
    }
</script>
