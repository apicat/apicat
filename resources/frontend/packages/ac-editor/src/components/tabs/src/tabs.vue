<template>
    <div class="tabs">
        <div ref="navWrap" class="tabs-nav-wrap" :class="[scrollable ? 'tabs-nav-scrollable' : '']">
            <div ref="navScroll" class="tabs-nav-scroll">
                <div ref="nav" class="tabs-nav" :style="navStyle">
                    <div class="tabs-inv-bar" :style="barStyle"></div>
                    <div class="tabs-tab" v-for="(item, index) in navList" :key="item.id" @click="handleChange(index)">
                        <input
                            type="text"
                            :id="item.id"
                            :value="item.label"
                            v-show="item.isEdit"
                            style="width: 80px"
                            @keydown.enter="onEnter"
                            @keydown.esc="onCancel($event, item)"
                            @blur="(e) => onChangeTabName(e, item)"
                        />
                        <span v-show="!item.isEdit">{{ item.label }}</span>
                        <i class="el-icon-edit" v-show="!item.isEdit" @click="onEditTabHeader($event, item)"></i>
                    </div>
                </div>
            </div>
        </div>
        <div class="pane-content">
            <slot></slot>
        </div>
    </div>
</template>
<script>
    import { ResizeObserver } from '@juggle/resize-observer'
    import shortid from 'shortid'
    import { getCurrentInstance } from 'vue'

    export default {
        name: 'Tabs',
        provide() {
            return { tabsInstance: this }
        },
        props: {
            value: {
                type: [String, Number],
            },
        },
        data() {
            return {
                observer: new ResizeObserver(() => this.initTabs()),
                navList: [],
                activeKey: this.value,
                barWidth: 0,
                barOffset: 0,
                scrollable: false,
                navStyle: {
                    transform: '',
                },
            }
        },
        computed: {
            barStyle() {
                return {
                    width: `${this.barWidth}px`,
                    transform: `translate3d(${this.barOffset}px,0px,0px)`,
                }
            },
        },
        methods: {
            onCancel(e, item) {
                e.target.value = item.label
                e.target.blur && e.target.blur()
            },

            onEditTabHeader(e, item) {
                e.preventDefault()
                e.stopPropagation()

                item.isEdit = !item.isEdit
                this.$nextTick(() => {
                    let input = document.querySelector('#' + item.id)
                    input && input.focus()
                })
            },

            onEnter(e) {
                e.target.blur && e.target.blur()
            },

            onChangeTabName(e, item) {
                item.isEdit = false
                if (!e.target.value.trim()) {
                    return
                }
                this.$emit('update-tab-title', {
                    key: item.name,
                    oldVal: item.label,
                    newVal: e.target.value,
                })
            },

            getTabs() {
                return this.$slots.default().filter((item) => item.type.name === 'TabPane')
            },

            initTabs() {
                this.updateNav()
                this.updateStatus()
                this.updateBar()
            },

            updateNav() {
                let navList = []

                this.getTabs().forEach((pane, index) => {
                    let oldNavItem = this.navList.find((item) => item.name === pane.type.name) || { id: `tab_header_input_${shortid()}` }
                    const { label, name } = pane.props
                    navList.push({
                        id: oldNavItem.id,
                        label: label,
                        name: name || index,
                        isEdit: oldNavItem.isEdit || false,
                    })
                    if (index === 0 && !this.activeKey) {
                        this.activeKey = pane.type.name
                    }
                })

                this.navList = navList
            },

            updateBar() {
                this.$nextTick(() => {
                    const index = this.navList.findIndex((nav) => nav.name === this.activeKey)
                    const elemTabs = this.$refs.navWrap.querySelectorAll('.tabs-tab')
                    const elemTab = elemTabs[index]
                    this.barWidth = elemTab ? elemTab.offsetWidth : 0
                    if (index > 0) {
                        let offset = 0
                        for (let i = 0; i < index; i++) {
                            offset += elemTabs[i].offsetWidth + 16
                        }
                        this.barOffset = offset
                    } else {
                        this.barOffset = 0
                    }
                })
            },

            updateStatus() {
                const tabs = this.getTabs()
                tabs.forEach((tab) => {
                    tab.show = tab.name === this.activeKey
                })
            },

            handleChange(index) {
                const nav = this.navList[index]
                this.activeKey = nav.name
            },
        },
        watch: {
            value(val) {
                this.activeKey = val
            },

            activeKey() {
                this.updateStatus()
                this.updateBar()
            },
        },
        setup() {
            // const instance = getCurrentInstance()
            // console.log(instance)
        },
    }
</script>
<style lang="less" scoped>
    .tabs {
        .tabs-nav-wrap {
            position: relative;
            border-bottom: 1px solid #dcdee2;
            margin-bottom: 16px;
        }
        .tabs-tab {
            position: relative;
            display: inline-block;
            margin-right: 16px;
            padding: 8px 16px;
            cursor: pointer;
        }
        .tabs-inv-bar {
            position: absolute;
            left: 0;
            bottom: 0;
            background-color: #2d8cf0;
            height: 2px;
            transition: transform 300ms ease-in-out;
        }

        .tabs-nav-scroll {
            overflow: hidden;
            white-space: nowrap;
        }
        .tabs-nav {
            position: relative;
            float: left;
            transition: transform 0.5s ease-in-out;
        }

        .tabs-nav-prev,
        .tabs-nav-next {
            position: absolute;
            width: 32px;
            line-height: 32px;
            text-align: center;
            cursor: pointer;
        }
        .tabs-nav-prev {
            left: 0;
        }
        .tabs-nav-next {
            right: 0;
        }
        .tabs-nav-scrollable {
            padding: 0 32px;
        }
        .tabs-nav-scroll-disabled {
            display: none;
        }
    }
</style>
