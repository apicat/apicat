import { h, render } from 'vue'
import Component from './MockRules.vue'
import { $off, $on } from '@apicat/shared'

export default class MockRules {
  constructor() {
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
    $on(this.$vm, 'on-ok', (rule) => this.onConfirmBtnClick(rule))
  }

  show(model) {
    this.paramNode = model.node
    this.$vm.show(model.node)
  }

  onConfirmBtnClick(rule) {
    if (this.paramNode) {
      this.paramNode['x-apicat-mock'] = rule
    }
  }

  destroy() {
    $off(this.$vm, 'on-ok')
    this.$vm = null
    this.paramNode = null
  }
}
