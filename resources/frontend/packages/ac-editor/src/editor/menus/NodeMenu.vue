<template>
    <ul class="float-menu-wrapper" v-show="menus.length">
        <li v-for="item in menus" :key="item._id" @click="onMenuClick(item)">
            <i :class="iconClass(item)"></i>
            <span>{{ item.text }}</span>
        </li>
    </ul>
</template>

<script>
    import { $emit } from '@natosoft/shared'

    export default {
        name: 'NodeMenu',
        data() {
            return {
                menus: [],
            }
        },
        computed: {},
        methods: {
            iconClass(item) {
                return ['editor-font', item.icon]
            },
            onMenuClick(menu) {
                $emit(this, 'on-click-menu', menu)
            },
        },
    }
</script>

<style scoped lang="less">
    .float-menu-wrapper {
        list-style: none;
        margin: 0;
        padding: 0;

        li {
            cursor: pointer;
            display: flex;
            align-items: center;
            padding: 4px 8px;
            //border-bottom: 1px solid #d3d3d3;
            font-size: 14px;

            i {
                margin-right: 5px;
            }

            &:hover {
                //background: #eaeaea;
            }

            &:last-child {
                border-bottom: none;
            }
        }
    }
</style>
