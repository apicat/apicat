import { h, render } from 'vue'
import Component from './MockRules.vue'
import { $on } from '@ac/shared'

export default class MockRules {
    constructor(editor) {
        this.editor = editor
        this.$vm = null
        this.paramNode = null

        this.element = this.renderVm()
        this.bindEvent()
    }

    renderVm() {
        const dom = document.createElement('div')
        const vNode = h(Component)
        render(vNode, dom)
        this.$vm = vNode.component.proxy
        return this.$vm.$el
    }

    bindEvent() {
        this.editor.on('destroy', () => this.destroy())
        $on(this.$vm, 'on-ok', (rule) => this.onConfirmBtnClick(rule))
    }

    show($responseParamTableVm, model) {
        this.paramNode = model.node
        this.$vm.show(model.node)
    }

    onConfirmBtnClick(rule) {
        if (this.paramNode) {
            this.paramNode['mock_rule'] = rule
            // Vue.set(this.paramNode, 'mock_rule', rule)
        }
    }

    destroy() {
        this.$vm = null
        this.paramNode = null
    }
}
