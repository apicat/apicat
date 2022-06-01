<template>
    <div>
        <FloatMenu v-if="!isLink && state && menus.length" :editor="editor" :menus="menus" :state="state" @menu-click="onMenuClick" />

        <LinkEditor v-show="isLink && range" ref="linkEditor" :mark="mark" @on-create="onCreateLink" @on-remove="onRemoveLink" @toggle-blank="onToggleBlank" />
    </div>
</template>

<script>
    import { setTextSelection } from 'prosemirror-utils'

    import FloatMenu from './FloatMenu.vue'
    import LinkEditor from '../LinkToolbar/LinkEditor.vue'

    import getTableColMenuItems from '../../menus/tableCol'
    import getTableRowMenuItems from '../../menus/tableRow'
    import getImageMenuItems from '../../menus/image'
    import getTableMenuItems from '../../menus/table'
    import getFormattingMenuItems from '../../menus/formatting'
    import { getColumnIndex, getMarkRange, getRowIndex, isMarkActive, isNodeActive, getMarkAttrs } from '../../utils'

    export default {
        props: ['editor', 'dictionary'],
        name: 'SelectionMenus',
        components: {
            FloatMenu,
            LinkEditor,
        },
        data() {
            return {
                menus: [],
                state: null,
                isActive: false,
                isLink: false,
                range: null,
                mark: { attrs: {} },
            }
        },
        methods: {
            onMenuClick(menu) {
                const { commands } = this.editor
                const command = commands[menu.name]
                command && command(menu.attrs)
            },

            onToggleBlank(attrs) {
                this.editor.commands.link(attrs)
                this.view.focus()
            },

            onRemoveLink() {
                this.editor.commands.link()
                this.view.focus()
            },

            onCreateLink(attrs) {
                this.editor.commands.link(attrs)
                this.range && this.view.dispatch(setTextSelection(this.range.to)(this.view.state.tr))
                this.$emit('on-close')
                this.view.focus()
            },

            updateFloatMenus() {
                const { dictionary, editor } = this
                const { view } = editor
                const { state } = view

                const { selection } = state
                const isCodeSelection = isNodeActive(state.schema.nodes.code_block)(state)

                this.range = null
                this.mark = { attrs: {} }
                this.isLink = false
                if (isCodeSelection) {
                    return []
                }

                const colIndex = getColumnIndex(state.selection)
                const rowIndex = getRowIndex(state.selection)

                const isTableSelection = colIndex !== undefined && rowIndex !== undefined
                const isLink = isMarkActive(state.schema.marks.link)(state)
                const range = getMarkRange(selection.$from, state.schema.marks.link)
                const isImageSelection = selection.node && selection.node.type.name === 'image'

                if (isLink && range) {
                    this.range = range
                    this.mark = this.range.mark
                    this.isLink = isLink

                    return []
                }

                let items = []
                if (isTableSelection) {
                    items = getTableMenuItems(dictionary)
                } else if (colIndex !== undefined) {
                    items = getTableColMenuItems(state, colIndex, dictionary)
                } else if (rowIndex !== undefined) {
                    items = getTableRowMenuItems(state, rowIndex, dictionary)
                } else if (isImageSelection) {
                    items = getImageMenuItems(state, dictionary)
                } else {
                    items = getFormattingMenuItems(state, dictionary)
                }

                this.menus = items.map((item) => {
                    item.isActive = item.active ? item.active(state) : false
                    item.isDisable = item.disable ? item.disable(state) : false

                    if (item.markType) {
                        const attrs = getMarkAttrs(state, item.markType)
                        const range = getMarkRange(selection.$from, item.markType)
                        item.attrs = attrs ? attrs : item.attrs
                        item.mark = range ? range.mark : null
                    }

                    if (item.getStyle) {
                        item.style = item.getStyle(item.attrs)
                    }
                    return item
                })
            },
        },
        created() {
            this.state = this.editor.state
            this.view = this.editor.view
        },
    }
</script>
