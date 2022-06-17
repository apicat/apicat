import tippy from 'tippy.js'
import MENU_TYPES from './menus'
import { noop } from 'lodash-es'
import { $once } from '@ac/shared'

export { NEW_MENUS } from './menus'

export default {
    onPopperItemClick: noop,
    onPopperHide: noop,

    popper: tippy(document.body, {
        appendTo: () => document.body,
        trigger: 'manual',
        placement: 'bottom',
        interactive: true,
        theme: 'light',
    }),

    show(type, el, props, appendToEl) {
        const menuConfig = MENU_TYPES[type]

        if (!menuConfig) {
            return
        }

        this.$vm = menuConfig.render(props)

        $once(this.$vm.component.proxy, 'on-menu-item-click', (menu) => {
            this.onPopperItemClick({ menu, ...props })
            this.hide()
        })

        this.popper.setProps({
            appendTo: () => appendToEl || document.body,
            getReferenceClientRect: () => el.getBoundingClientRect(),
            content: this.$vm.el,
            onHide: () => {
                this.onPopperHide()
            },
        })

        this.popper.show()
    },

    hide() {
        this.popper.hide()
    },
}
