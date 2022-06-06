<template>
    <div :class="wrapperClass" ref="wrapper">
        <ol ref="menus">
            <template v-for="(menu, index) in menus" :key="menu.uid">
                <li
                    v-if="menu.name !== 'separator'"
                    :class="menuItemClass(index)"
                    @click="onMenuItemClick(menu)"
                    @mousemove="onMouseOverMenuItem($event, index)"
                >
                    <img :src="menu.img" alt="" />
                    <div class="menu-item">
                        <p>{{ menu.title }}</p>
                        <span>{{ menu.desc }}</span>
                    </div>
                </li>

                <li v-else class="hr" :data-idx="index">
                    {{ menu.title || '' }}
                </li>
            </template>
        </ol>
        <input ref="file" style="display: none" type="file" @change="handleImagePicked" accept="image/*" />
    </div>
</template>
<script>
    import scrollIntoView from 'smooth-scroll-into-view-if-needed'
    import { findParentNode } from 'prosemirror-utils'
    import { capitalize, noop } from 'lodash-es'
    import uuid from 'shortid'
    import { getDataTransferFiles } from '../../utils'
    import { insertFiles } from '../../commands'
    import getMenuItems from '../../menus/block'
    import { $emit } from '@ac/shared'
    import { markRaw } from 'vue'

    export default {
        props: {
            editor: {
                type: Object,
                default: () => ({}),
            },
            dictionary: {
                type: Object,
                default: () => ({}),
            },
            onImageUploadStart: {
                type: Function,
                default: () => {},
            },
            uploadImage: {
                type: Function,
                default: () => {},
            },
            onImageUploadStop: {
                type: Function,
                default: () => {},
            },
            isActive: {
                type: Boolean,
                default: false,
            },
            searchKeyword: {
                type: String,
                default: '',
            },
        },
        data() {
            return {
                mouseHoverInfo: { x: -1, y: -1 },
                allMenus: this.getMenus().concat([]),
                view: null,
                commands: null,
                selectedIndex: -1,
            }
        },
        computed: {
            wrapperClass() {
                return ['menu-wrapper scroll-content']
            },
            menus: function () {
                if (!(this.searchKeyword || '').trim()) {
                    return this.allMenus
                }

                let menus = this.allMenus.filter((menu) => {
                    if (menu.name === 'separator' || !menu.keywords) {
                        return true
                    }
                    return menu.keywords.includes(this.searchKeyword)
                })

                if (!menus.length) {
                    return this.allMenus
                }
                let len = menus.length
                let removeIndex = []

                for (let i = 0; i < len; i++) {
                    let next = menus[i + 1]
                    // 相邻
                    if (menus[i].name === 'separator' && next && next.name === 'separator') {
                        removeIndex.push(i)
                    }

                    // 末尾
                    if (menus[i].name === 'separator' && !next) {
                        removeIndex.push(i)
                    }
                }

                // 移除分割,倒序移除，防止数组下表错乱
                for (let j = removeIndex.length - 1; j >= 0; j--) {
                    menus.splice(removeIndex[j], 1)
                }

                if (!menus.length) {
                    return this.allMenus
                }

                return menus
            },
        },

        watch: {
            searchKeyword: function () {
                if (this.isActive) {
                    this.selectedIndex = 0
                    this.changeSelectIndex(false)
                }
            },

            isActive: function () {
                if (this.isActive) {
                    this.selectedIndex = 0
                    this.changeSelectIndex(false)
                }
            },

            selectedIndex: function () {
                this.$refs.menus && this.goScrollIntoView(this.$refs.menus.children[this.selectedIndex])
            },
        },

        methods: {
            onMouseOverMenuItem(e, index) {
                const { x, y } = this.mouseHoverInfo

                if (e.clientX !== x && e.clientY !== y) {
                    this.selectedIndex = index

                    this.mouseHoverInfo.x = e.clientX
                    this.mouseHoverInfo.y = e.clientY
                }
            },

            getMenus() {
                let menus = getMenuItems(this.dictionary)
                menus.map((item) => {
                    item.uid = uuid()
                    return item
                })
                return menus
            },

            onMenuItemClick(menuItem) {
                this.close()

                switch (menuItem.name) {
                    case 'image':
                        this.triggerImagePick()
                        break
                    case 'link':
                        this.triggerCreateLink()
                        break
                    default:
                        this.insertBlock(menuItem)
                }
            },

            triggerCreateLink() {
                this.clearSearch()
                $emit(this, 'onCreateLinkTrigger')
            },

            triggerImagePick() {
                if (this.$refs.file) {
                    this.$refs.file.click()
                }
            },

            insertBlock(menuItem) {
                this.clearSearch()

                let command = this.commands[menuItem.name]

                if (!command) {
                    if (menuItem.name.startsWith('trigger_')) {
                        try {
                            const fn = this.editor.options[menuItem.name.split('_')[1]] || noop
                            fn && fn()
                        } catch (e) {
                            //
                        }
                        return
                    }
                    command = this.commands[`create${capitalize(menuItem.name)}`]
                }

                command && command(menuItem.attrs)
            },

            clearSearch() {
                const { state, dispatch } = this.view
                const parent = findParentNode((node) => !!node)(state.selection)

                if (parent) {
                    dispatch(state.tr.insertText('', parent.start, parent.start + parent.node.textContent.length))
                }
            },

            goScrollIntoView(node) {
                node &&
                    scrollIntoView(node, {
                        behavior: 'smooth',
                        block: 'nearest',
                        inline: 'nearest',
                    })
            },

            menuItemClass(index) {
                return {
                    selected: this.selectedIndex === index,
                }
            },

            close() {
                $emit(this, 'close')
            },

            handleImagePicked(event) {
                const files = getDataTransferFiles(event)
                const { uploadImage, onImageUploadStart, onImageUploadStop, dictionary } = this
                const { view } = this.editor
                const { state, dispatch } = view
                const parent = findParentNode((node) => !!node)(state.selection)

                if (parent) {
                    dispatch(state.tr.insertText('', parent.pos, parent.pos + parent.node.textContent.length + 1))

                    insertFiles(view, event, parent.pos, files, {
                        uploadImage,
                        onImageUploadStart,
                        onImageUploadStop,
                        dictionary,
                    })
                }

                if (this.$refs.file) {
                    this.$refs.file.value = ''
                }

                this.close()
            },

            changeSelectIndex(isPrev) {
                if (this.menus.length) {
                    const total = this.menus.length - 1
                    let index = this.selectedIndex - (isPrev ? 1 : -1)

                    if (isPrev && index <= 0) {
                        index = total
                    }

                    if (!isPrev && index > total) {
                        index = 0
                    }

                    const menuItem = this.menus[index]

                    if (menuItem && menuItem.name === 'separator') {
                        index = index - (isPrev ? 1 : -1)
                    }

                    this.selectedIndex = isPrev ? Math.max(0, index) : Math.min(index, total)
                } else {
                    this.close()
                }
            },

            handleKeyDown(event) {
                if (!this.isActive) return

                if (event.key === 'Enter') {
                    const item = this.menus[this.selectedIndex]
                    item && this.onMenuItemClick(item)
                }

                if (event.key === 'ArrowUp') {
                    event.preventDefault()
                    event.stopPropagation()
                    this.changeSelectIndex(true)
                }

                if (event.key === 'ArrowDown' || event.key === 'Tab') {
                    event.preventDefault()
                    event.stopPropagation()
                    this.changeSelectIndex(false)
                }

                if (event.key === 'Escape') {
                    this.close()
                }
            },
        },

        mounted() {
            const { view, commands } = this.editor
            this.view = markRaw(view)
            this.commands = commands

            window.addEventListener('keydown', this.handleKeyDown)
        },

        unmounted() {
            window.removeEventListener('keydown', this.handleKeyDown)
        },
    }
</script>
