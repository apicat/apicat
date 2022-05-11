import { h, markRaw, render } from 'vue'
import { NodeSelection } from 'prosemirror-state'
import getMarkRange from '../utils/getMarkRange'

export default class ComponentView {
    constructor(component, { editor, extension, node, view, getPos, decorations, innerDecorations }) {
        this.component = component
        this.editor = editor
        this.extension = extension
        this.node = node
        this.view = view
        this.decorations = decorations
        this.innerDecos = innerDecorations

        this.isNode = !!this.node.marks
        this.isMark = !this.isNode
        this.getPos = this.isMark ? this.getMarkPos : getPos

        this.isSelected = false

        this.dom = this.createDOM()
    }

    createDOM() {
        let nodeViewWrapper = this.node.type.spec.inline ? document.createElement('span') : document.createElement('div')

        let vueComponentWrapper = this.node.type.spec.inline ? document.createElement('span') : document.createElement('div')

        nodeViewWrapper.className = this.node.type.name + '_node_view_wrapper'

        // todo tabindex 影响鼠标拖拽图片，不知道为何，暂且记录
        if (this.node.type.name !== 'image') {
            nodeViewWrapper.setAttribute('tabindex', 0)
            vueComponentWrapper.setAttribute('tabindex', 0)
        }

        // 避免共享数据引用地址
        this.node.attrs = JSON.parse(JSON.stringify(this.node.attrs || {}))

        const props = {
            editor: markRaw(this.editor),
            view: markRaw(this.view),
            node: markRaw(this.node),
            getPos: () => this.getPos(),
            decorations: this.decorations,
            isSelected: this.isSelected,
            isReadOnly: !this.view.editable,
            updateAttributes: (attrs) => this.updateAttributes(attrs),
            options: this.extension.options,
        }

        const vNode = h(this.component, props)
        render(vNode, vueComponentWrapper)

        if (typeof this.extension.setSelection === 'function') {
            this.setSelection = this.extension.setSelection
        }

        if (typeof this.extension.update === 'function') {
            this.update = this.extension.update
        }

        this.innerDecos &&
            this.innerDecos.find &&
            this.innerDecos.find().map((d) => {
                const elem = typeof d.type.toDOM === 'function' ? d.type.toDOM() : d.type.toDOM
                nodeViewWrapper.appendChild(elem)
            })

        this.$vm = vNode.component.proxy
        this.$props = vNode.component.props

        nodeViewWrapper.appendChild(this.$vm.$el)

        return nodeViewWrapper
    }

    update(node, decorations) {
        if (node.type !== this.node.type) {
            return false
        }

        if (node === this.node && this.decorations === decorations) {
            return true
        }

        this.node = node
        this.decorations = decorations

        this.updateComponentProps({
            node,
            decorations,
        })

        return true
    }

    updateComponentProps(props) {
        if (!this.$props) {
            return
        }

        Object.entries(props).forEach(([key, value]) => {
            this.$props[key] = value
        })
    }

    updateAttributes(attributes) {
        if (!this.editor.view.editable) {
            return
        }

        const { state } = this.editor.view
        const pos = this.getPos()
        const transaction = state.tr.setNodeMarkup(pos, undefined, {
            ...this.node.attrs,
            ...attributes,
        })

        this.editor.view.dispatch(transaction)
    }

    // disable (almost) all prosemirror event listener for node views
    stopEvent(event) {
        if (!this.dom) {
            return false
        }

        if (typeof this.extension.stopEvent === 'function') {
            return this.extension.stopEvent(event)
        }

        const target = event.target
        let contentDOM

        const isInElement =
            this.dom.contains(target) && !((contentDOM = this.contentDOM) === null || contentDOM === void 0 ? void 0 : contentDOM.contains(target))

        // any event from child nodes should be handled by ProseMirror
        if (!isInElement) {
            return false
        }

        const isInput = ['INPUT', 'BUTTON', 'SELECT', 'TEXTAREA', 'CODE'].includes(target.tagName) || target.isContentEditable

        // any input event within node views should be ignored by ProseMirror
        if (isInput) {
            return true
        }

        // const isDraggable = !!this.node.type.spec.draggable;
        const isDraggable = true
        const isSelectable = NodeSelection.isSelectable(this.node)
        const isCopyEvent = event.type === 'copy'
        const isPasteEvent = event.type === 'paste'
        const isCutEvent = event.type === 'cut'
        const isClickEvent = event.type === 'mousedown'
        const isDragEvent = event.type.startsWith('drag') || event.type === 'drop'

        // ProseMirror tries to drag selectable nodes
        // even if `draggable` is set to `false`
        // this fix prevents that
        if (!isDraggable && isSelectable && isDragEvent) {
            event.preventDefault()
        }

        if (isDraggable && isDragEvent) {
            // event.preventDefault();
            return false
        }

        // these events are handled by prosemirror
        if (isCopyEvent || isPasteEvent || isCutEvent || (isClickEvent && isSelectable)) {
            return false
        }

        return true
    }

    ignoreMutation() {
        return true
    }

    selectNode() {
        this.updateComponentProps({
            isSelected: true,
        })
    }

    deselectNode() {
        this.updateComponentProps({
            isSelected: false,
        })
    }

    getMarkPos() {
        const pos = this.view.posAtDOM(this.dom)
        const resolvedPos = this.view.state.doc.resolve(pos)
        return getMarkRange(resolvedPos, this.node.type)
    }

    destroy() {
        this.$vm = null
        this.dom = null
    }
}
