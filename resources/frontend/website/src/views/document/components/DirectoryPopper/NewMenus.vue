<template>
    <ul class="ac-popper-menu">
        <template v-for="menu in menus" :key="menu.value">
            <li
                v-if="node.level < 5 || (node.level === 5 && menu.value !== 'newCatetory')"
                class="ac-popper-menu__item"
                :class="{ 'border-t': menu.divided }"
                @click.stop="onPopperItemClick(menu)"
            >
                <span v-if="menu.icon" class="icon iconfont ac-popper-menu__icon" :class="menu.icon" />
                <img v-if="menu.img" class="ac-popper-menu__icon" :src="menu.img" alt="" />{{ menu.text }}
            </li>
        </template>
    </ul>
</template>

<script>
    import { NEW_MENUS } from './menus'
    import { $emit } from '@ac/shared'

    export default {
        emits: ['on-menu-item-click'],
        name: 'NewMenus',
        props: {
            node: {
                type: Object,
                default: () => ({}),
            },
        },
        setup() {
            return {
                menus: NEW_MENUS,
            }
        },
        methods: {
            onPopperItemClick(menu) {
                $emit(this, 'on-menu-item-click', menu)
            },
        },
    }
</script>
