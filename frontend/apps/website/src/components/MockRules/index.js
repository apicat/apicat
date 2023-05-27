import { h, render } from 'vue'
import Component from './MockRules.vue'
import { $off, $on } from '@apicat/shared'

let ins = null
export default class MockRules {
  constructor() {
    if (ins) {
      return ins
    }

    this.$vm = null
    this.model = null
    this.onOk = null

    this.renderVm()
    this.bindEvent()

    ins = this
    return this
  }

  renderVm() {
    const dom = document.createElement('div')
    const vNode = h(Component)
    render(vNode, dom)
    this.$vm = vNode.component.proxy
    return this.$vm.$el
  }

  bindEvent() {
    $on(this.$vm, 'on-ok', (rule) => this.onConfirmBtnClick(rule))
  }

  show({ model, onOk }) {
    if (!this.$vm) {
      this.renderVm()
    }
    this.onOk = onOk
    this.model = model
    this.$vm.show(model)
  }

  onConfirmBtnClick(rule) {
    this.onOk && this.onOk(rule)
  }

  destroy() {
    $off(this.$vm, 'on-ok')
    this.$vm = null
    this.model = null
    this.onOk = null

    ins = null
  }
}

export const mockRulesModal = new MockRules()
