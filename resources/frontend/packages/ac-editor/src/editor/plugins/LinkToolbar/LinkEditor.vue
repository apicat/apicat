<template>
    <div class="link-bar-wrapper">
        <input type="text" v-model="attrs.href" class="link-input" ref="input" placeholder="链接地址" @keydown.enter="onKeydown" />
        <button class="link-btn" v-if="!isCreate" @click="onRemoveLinkClick">
            <i class="editor-font editor-trash" title="移除链接"></i>
        </button>

        <button class="link-btn" v-if="isCreate" @click="onCloseClick">
            <i class="editor-font el-icon-close" title="关闭"></i>
        </button>
    </div>
</template>

<script>
    import { $emit } from '@natosoft/shared'

    export default {
        props: {
            mark: {
                type: Object,
                required: true,
                default: () => ({ attrs: {} }),
            },

            isCreate: {
                type: Boolean,
                default: false,
            },
        },
        name: 'LinkEditor',
        data() {
            return {
                attrs: { ...this.mark.attrs },
            }
        },
        watch: {
            mark: function () {
                this.attrs = { ...this.mark.attrs }
            },
        },
        computed: {
            isOpenClass: function () {
                return [
                    'link-btn',
                    {
                        active: this.mark.attrs.openInNewTab,
                    },
                ]
            },
        },
        methods: {
            onCloseClick(e) {
                e.stopPropagation()
                this.$refs.input.value = ''
                $emit(this, 'on-close')
            },

            toggleIsOpen(e) {
                e.stopPropagation()
                $emit(this, 'toggle-blank', { ...this.attrs, openInNewTab: !this.mark.attrs.openInNewTab })
                this.$refs.input.value = ''
            },

            onRemoveLinkClick(e) {
                e.stopPropagation()
                $emit(this, 'on-remove')
            },

            onKeydown(event) {
                event.preventDefault()
                $emit(this, 'on-create', { ...this.attrs, href: this.$refs.input.value })
                this.$refs.input.value = ''
            },

            focus() {
                this.$nextTick(() => {
                    if (this.$refs.input || this.mark.attrs.href === '') {
                        this.$refs.input.focus()
                    }
                })
            },
        },
    }
</script>
