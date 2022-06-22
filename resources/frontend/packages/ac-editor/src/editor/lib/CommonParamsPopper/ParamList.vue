<template>
    <ul class="param-list">
        <li v-for="(item, index) in list" :key="item.value" :class="paramItemClass(index)">
            <p class="param-list__text" @click="onParamItemClick(item.value)">{{ item.value }}</p>
            <el-icon @click="onDeleteParamBtnClick(item.value)"><Delete /></el-icon>
        </li>
    </ul>
</template>

<script>
    import { ElIcon } from 'element-plus'
    import { Delete } from '@element-plus/icons-vue'
    import { $emit } from '@natosoft/shared'

    export default {
        name: 'ParamList',
        components: {
            ElIcon,
            Delete,
        },
        data() {
            return {
                list: [],
                isActive: false,
                selectedIndex: -1,
            }
        },
        watch: {
            isActive: function () {
                if (!this.isActive) {
                    this.selectedIndex = -1
                }
            },
        },
        methods: {
            paramItemClass(index) {
                return [
                    'param-list__item',
                    {
                        active: this.selectedIndex === index,
                    },
                ]
            },

            close() {
                $emit(this, 'on-close')
            },

            onDeleteParamBtnClick(val) {
                $emit(this, 'on-delete', val)
            },

            onParamItemClick(val) {
                $emit(this, 'on-ok', val)
            },

            handleKeyDown(event) {
                if (!this.isActive) return

                if (event.key === 'Enter') {
                    event.preventDefault()
                    event.stopPropagation()
                    const item = this.list[this.selectedIndex]
                    item && this.onParamItemClick(item.value)
                }

                if (event.key === 'ArrowUp') {
                    event.preventDefault()
                    event.stopPropagation()

                    if (this.list.length) {
                        const prevIndex = this.selectedIndex - 1
                        const total = this.list.length - 1

                        if (this.selectedIndex === 0) {
                            this.selectedIndex = total
                        } else {
                            this.selectedIndex = Math.max(0, prevIndex)
                        }
                    }
                }

                if (event.key === 'ArrowDown' || event.key === 'Tab') {
                    event.preventDefault()
                    event.stopPropagation()

                    if (this.list.length) {
                        const total = this.list.length - 1
                        const nextIndex = this.selectedIndex + 1

                        if (this.selectedIndex === total) {
                            this.selectedIndex = 0
                        } else {
                            this.selectedIndex = Math.min(nextIndex, total)
                        }
                    }
                }

                if (event.key === 'Escape') {
                    this.close()
                }
            },
        },
        mounted() {
            window.addEventListener('keydown', this.handleKeyDown)
        },

        destroyed() {
            window.removeEventListener('keydown', this.handleKeyDown)
        },
    }
</script>
