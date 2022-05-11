import { h, render } from 'vue'
import tippy from 'tippy.js'
import { $on } from '@ac/shared'
import { deleteSelection } from 'prosemirror-commands'
import NodeMenu from '../menus/NodeMenu.vue'

const BaseCommands = [
    {
        icon: 'editor-trash',
        text: '删除',
        _id: 'delete_node',
        command: (state, dispatch) => deleteSelection(state, dispatch),
    },
]

export default class HandleMenu {
    constructor(editor) {
        this.tippy = null
        this.ref = null
        this.editor = editor

        this.node = null
        this.nodeDom = null

        this.renderMenuWrapper()
    }

    renderMenuWrapper() {
        let VNodeMenu = h(NodeMenu)
        let dom = document.createElement('div')
        render(VNodeMenu, dom)
        this.$vm = VNodeMenu.component.proxy

        $on(this.$vm, 'on-click-menu', (menu) => menu.command && this.execCommand(menu))
        this.menuWrapper = this.$vm.$el
    }

    renderMenu() {
        if (!this.node) {
            // console.log("未选中node！");
            return
        }
        const menus = []
        if (this.editor.nodeEditViewManager.hasEditView(this.node)) {
            menus.push({
                icon: 'editor-bianji',
                text: '编辑',
                _id: 'edit_node',
                command: (state, dispatch, node, nodeDom) => this.editor.nodeEditViewManager.updateAttrs(node, nodeDom),
            })
        }
        this.$vm.menus = menus.concat(BaseCommands)
    }

    show(ref, node, nodeDom) {
        if (!ref || !node || !nodeDom) {
            return
        }

        this.destroy()
        this.ref = ref
        this.node = node
        this.nodeDom = nodeDom

        this.renderMenu()

        const { view } = this.editor

        this.tippy = tippy(view.dom.parentNode || document.body, {
            content: this.menuWrapper,
            getReferenceClientRect: () => {
                if (document.body.contains(ref)) {
                    return ref.getBoundingClientRect()
                }

                let rect = nodeDom.getBoundingClientRect().toJSON()
                rect.left -= 30
                return rect
            },
            theme: 'light',
            placement: 'left-start',
            trigger: 'manual',
            arrow: false,
            interactive: true,
        })
        this.tippy.show()
    }

    execCommand(menu) {
        const { view } = this.editor
        if (!view || !menu || !menu.command) {
            return
        }
        let { state, dispatch } = view

        menu.command(state, dispatch, this.node, this.nodeDom)

        this.destroy()

        menu._id === 'delete_node' && view && view.focus()
    }

    destroy() {
        this.tippy && this.tippy.destroy()
        this.tippy = null
        this.ref = null
        this.node = null
        this.nodeDom = null
    }
}
