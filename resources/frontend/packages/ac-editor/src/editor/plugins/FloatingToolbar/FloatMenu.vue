<template>
    <div v-if="menus.length" class="float-menu">
        <template v-for="(item, index) in menus">
            <button
                v-if="item.name !== 'separator' && item.visible !== false && item.icon"
                :key="index"
                @click="onMenuItemClick(item)"
                :class="getFloatMenuClass(item)"
            >
                <i :title="item.tooltip" :data-key="item.name" :ref="item.popper ? 'popper' : null" :class="'editor-font ' + item.icon" :style="item.style" />
            </button>
        </template>
    </div>
</template>

<script>
    import FloatMenuPopper from './FloatMenuPopper'

    export default {
        props: {
            menus: {
                type: Array,
                default: () => [],
            },
            state: {
                type: Object,
                default: null,
            },
            editor: {
                type: Object,
                default: null,
            },
        },
        name: 'FloatMenu',
        data() {
            return {
                floatMenuPopper: null,
            }
        },
        watch: {
            menus: function () {
                this.$nextTick(() => {
                    if (this.floatMenuPopper) {
                        this.floatMenuPopper.updateMenus(this.menus)
                        return
                    }

                    this.initPopper()
                })
            },
        },
        methods: {
            getFloatMenuClass(item) {
                return [
                    'float-menu-item',
                    {
                        active: item.isActive,
                        disable: item.isDisable,
                    },
                ]
            },

            onMenuItemClick(item) {
                if (item.isDisable) {
                    return
                }

                item.name && this.$emit('menu-click', item)
            },

            initPopper() {
                if (this.$refs['popper'] && this.editor) {
                    this.floatMenuPopper = new FloatMenuPopper(this.$refs['popper'], this.menus, this.editor)
                }
            },
        },
        mounted() {
            this.initPopper()
        },
        unmounted() {
            this.floatMenuPopper && this.floatMenuPopper.destroy()
            this.floatMenuPopper = null
            // console.log("FloatMenu.vue destroyed")
        },
    }
</script>
