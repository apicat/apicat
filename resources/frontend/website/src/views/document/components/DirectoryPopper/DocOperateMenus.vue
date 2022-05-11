<template>
    <ul class="ac-popper-menu ac-popper-menu--thin">
        <template v-for="menu in menus">
            <li
                v-if="
                    (menu.value === 'delete' || menu.value === 'export' || menu.value === 'share' || menu.value === 'copy') && isDetailPage && source.isCurrent
                "
                class="ac-popper-menu__item"
                :class="{ 'border-t': menu.divided }"
                :key="menu.value"
                @click.stop="onPopperItemClick(menu)"
            >
                {{ menu.text }}
            </li>

            <li
                v-if="menu.value === 'delete' && !isDetailPage && source.isCurrent"
                class="ac-popper-menu__item"
                :class="{ 'border-t': menu.divided }"
                :key="menu.value"
                @click.stop="onPopperItemClick(menu)"
            >
                {{ menu.text }}
            </li>

            <li
                v-if="!source.isCurrent"
                class="ac-popper-menu__item"
                :class="{ 'border-t': menu.divided }"
                :key="menu.value"
                @click.stop="onPopperItemClick(menu)"
            >
                {{ menu.text }}
            </li>
        </template>
    </ul>
</template>

<script>
    import { DOC_OPERATE_MENUS } from './menus'
    import { $emit } from '@ac/shared'

    export default {
        emits: ['on-menu-item-click'],
        name: 'DocOperateMenus',
        props: {
            node: {
                type: Object,
                default: () => ({}),
            },
            source: {
                type: Object,
                default: () => ({}),
            },
            isDetailPage: {
                type: Boolean,
                default: false,
            },
        },
        data() {
            return {
                menus: DOC_OPERATE_MENUS,
            }
        },
        methods: {
            onPopperItemClick(menu) {
                $emit(this, 'on-menu-item-click', menu)
            },
        },
    }
</script>
