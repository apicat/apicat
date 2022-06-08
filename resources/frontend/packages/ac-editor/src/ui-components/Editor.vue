<template>
    <div ref="editor" class="ac-editor"></div>
</template>

<script>
    import { markRaw } from 'vue'
    import AcEditor, { EDITOR_EVENTS } from '../editor'

    export default {
        name: 'AcEditor',
        emits: ['on-change'],
        props: {
            document: {
                type: Object,
                default: () => ({}),
            },

            readonly: {
                type: Boolean,
                default: false,
            },

            options: {
                type: Object,
                default: () => ({}),
            },
        },

        watch: {
            document: function () {
                this.editor && this.editor.setContent(this.document)
            },
        },

        data() {
            return {
                editor: null,
            }
        },

        methods: {
            focus() {
                this.editor.view.focus()
            },
        },

        mounted() {
            this.editor = markRaw(
                AcEditor.create(this.$refs.editor, {
                    ...this.options,
                    readonly: this.readonly,
                    content: this.document,
                })
            )

            this.editor.on(EDITOR_EVENTS.Update, () => {
                this.$emit('on-change', this.editor)
            })

            window.$editor = this.editor
            window.$commands = this.editor.commands
        },

        unmounted() {
            this.editor && this.editor.destroy()
        },
    }
</script>
