<template>
    <node-view-wrapper :class="className">
        <ParamTreeTable :editor="editor" :readonly="isReadOnly" :data="params" @on-update="updateAttrs" @on-sort="onSortEnd" />
    </node-view-wrapper>
</template>

<script>
    import ParamTreeTable from '../../../components/ParamTreeTable.vue'
    import { NodeViewWrapper } from '../../../components/NodeViewWrapper'
    import { ref, watch } from 'vue'

    export default {
        name: 'ParamTable',
        props: ['editor', 'node', 'view', 'getPos', 'isSelected', 'isReadOnly', 'updateAttributes'],
        components: {
            ParamTreeTable,
            NodeViewWrapper,
        },
        computed: {
            className() {
                return [
                    'param-table',
                    {
                        'ProseMirror-selectednode': this.isSelected && !this.isReadOnly,
                    },
                ]
            },
        },
        methods: {
            updateAttrs(params) {
                this.updateAttributes({ params })
            },
            onSortEnd() {
                this.view && this.view.focus()
            },
        },

        setup(props) {
            const params = ref(props.node.attrs.params)

            watch(
                params.value,
                () => {
                    props.updateAttributes({ params: params.value })
                },
                { deep: true }
            )

            return {
                params,
            }
        },
    }
</script>
