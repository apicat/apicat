import { createApp } from 'vue'
import tippy from 'tippy.js'
import { $once } from '@ac/shared'

export default class NodeEditViewManager {
    constructor(extensions, editor) {
        this.$vm = null
        this.$app = null

        this.components = this.createComponents(extensions)

        this.editor = editor
        this.view = editor.view

        this.createPopper()
    }

    createComponents(extensions = []) {
        let comps = {}
        extensions
            .filter((extension) => extension.edit)
            .forEach((extension) => {
                comps[extension.name] = extension.edit
            })

        return comps
    }

    renderVm(componentKey, props = {}) {
        const sfc = this.components[componentKey]
        if (!sfc) {
            return null
        }

        const dom = document.createElement('div')

        this.$app = createApp(sfc, {
            ...props,
            editor: this.editor,
        })

        const vm = this.$app.mount(dom)

        $once(vm, 'on-update-attr', () => {
            this.tippy.hide()
        })

        $once(vm, 'on-create', () => {
            this.tippy.hide()
        })

        $once(vm, 'on-close', () => {
            this.tippy.hide()
        })

        return vm
    }

    createPopper() {
        this.tippy = tippy(document.body, {
            appendTo: () => document.body,
            duration: 0,
            getReferenceClientRect: null,
            content: '',
            theme: 'light',
            placement: 'bottom-start',
            trigger: 'manual',
            arrow: false,
            interactive: true,
            maxWidth: 'none',
            onHide: () => this.hide(),
        })

        this.tippy.hide()
    }

    createNode(node, attrs) {
        node = this.editor.schema.nodes[node.type.name]
        if (node) {
            const { dispatch, state } = this.view
            const tr = state.tr
            dispatch(tr.replaceSelectionWith(node.create(attrs)))
        }
    }

    updateNodeAttrs(node, attrs) {
        if (!this.view.editable) {
            return
        }

        const { state } = this.view
        const { node: selectNode, from } = state.selection

        if (selectNode.type.name !== node.type.name) {
            return
        }

        const transaction = state.tr.setNodeMarkup(from, null, attrs)
        this.view.dispatch(transaction)
    }

    hasEditView(node) {
        if (!node || !node.type || !node.type.name) {
            return false
        }

        return this.components[node.type.name] !== undefined
    }

    updateAttrs(node, nodeDom) {
        this.$vm = this.renderVm(node.type.name)

        if (node && nodeDom && this.$vm) {
            this.$vm.setNode && this.$vm.setNode({ node })
            this.tippy.setProps({
                getReferenceClientRect: () => nodeDom.getBoundingClientRect(),
                content: this.$vm.$el,
            })
            this.tippy.show()
            this.focus()
        }
    }

    create(schemaName) {
        this.$vm = this.renderVm(schemaName, { isCreate: true })
        const schema = this.editor.schema.nodes[schemaName]

        let { view } = this.editor
        let { selection } = this.editor.state
        const { node: element } = view.domAtPos(selection.from)

        if (schema && this.$vm) {
            let node = schema.create()

            this.$vm.setNode && this.$vm.setNode({ node })

            this.tippy.setProps({
                getReferenceClientRect: () => {
                    let rect = element.getBoundingClientRect().toJSON()
                    rect.width = 1
                    return rect
                },
                content: this.$vm.$el,
            })

            this.tippy.show()
            this.focus()
        }
    }

    focus() {
        // todo 为什么会影响滚动行为？
        setTimeout(() => {
            this.$vm && this.$vm.focus && this.$vm.focus()
        }, 100)
    }

    hide() {
        this.$vm.onHide && this.$vm.onHide((isCreate, node, attrs) => (isCreate ? this.createNode(node, attrs) : this.updateNodeAttrs(node, attrs)))
        this.$app && this.$app.unmount()
        this.$vm = null
        this.$app = null

        setTimeout(() => this.view.focus(), 0)
    }

    destroy() {
        this.tippy && this.tippy.destroy()
    }
}
