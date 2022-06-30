<template>
    <ul class="ac-popper-menu ac-popper-menu--thin">
        <li v-for="menu in menus" class="ac-popper-menu__item" :class="{ 'border-t': menu.divided }" :key="menu.value" @click.stop="onPopperItemClick(menu)">
            {{ menu.text }}
        </li>
    </ul>
</template>

<script>
    import { DIR_OPERATE_MENUS } from './menus'
    import { $emit } from '@natosoft/shared'

    export default {
        name: 'DirOperateMenus',
        data() {
            return {
                menus: DIR_OPERATE_MENUS,
            }
        },
        methods: {
            onPopperItemClick(menu) {
                $emit(this, 'on-menu-item-click', menu)
            },
        },
    }
</script>
