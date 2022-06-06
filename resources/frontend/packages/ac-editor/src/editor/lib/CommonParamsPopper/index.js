import tippy from 'tippy.js'
import { h, render } from 'vue'
import ParamList from './ParamList.vue'
import noop from 'lodash/noop'
import { $on } from '@ac/shared'

export default class CommonParamsPopper {
    constructor(editor) {
        const { getAllCommonParams, addCommonParam, deleteCommonParam } = editor.options
        this.editorDom = editor.view.dom

        this.getAllCommonParams = getAllCommonParams
        this.addCommonParam = addCommonParam
        this.deleteCommonParam = deleteCommonParam
        this.onParamItemClick = noop

        this.refDom = null
        this.node = null
        this.$vm = null
        this.$app = null
        this.popper = null
        this.searchKeyword = ''

        this.params = {
            list: [],
            map: {},
        }

        this._event = []

        this.initialize()
    }

    initialize() {
        this.renderVm()
        this.bindEvent()
        this.createPopper()
        this.getAllParams()
    }

    renderVm() {
        const dom = document.createElement('div')

        const ParamListComponent = h(ParamList)
        render(ParamListComponent, dom)

        this.$vm = ParamListComponent.component.proxy
    }

    bindEvent() {
        $on(this.$vm, 'on-close', () => this.hide())

        $on(this.$vm, 'on-delete', (val) => {
            this.deleteParam(val)
        })

        $on(this.$vm, 'on-ok', (val) => {
            this.onParamItemClick && this.onParamItemClick(this.node, val)
            this.hide()
        })
    }

    createPopper() {
        this.popper = tippy(document.body, {
            appendTo: () => this.editorDom.parentElement,
            content: this.$vm.$el,
            theme: 'params',
            placement: 'bottom',
            maxWidth: 'none',
            trigger: 'manual',
            arrow: true,
            interactive: true,
            popperOptions: {
                strategy: 'fixed',
                modifiers: [
                    {
                        name: 'arrow',
                        options: {
                            left: 20,
                        },
                    },
                ],
            },
            onShow: () => {
                this.$vm.isActive = true
            },

            onHide: () => {
                this.$vm.isActive = false
            },
        })

        this.popper.hide()
    }

    show({ inputDom, node, onParamItemClick }) {
        this.hide()

        if (!inputDom || inputDom.tagName.toLowerCase() !== 'input') {
            return
        }

        this.refDom = inputDom
        this.node = node
        this.onParamItemClick = onParamItemClick

        if (this.popper && this.refDom) {
            this.popper.setProps({
                getReferenceClientRect: () => {
                    if (this.refDom) {
                        return this.refDom.getBoundingClientRect()
                    }
                    return { top: -999, left: -999, width: 0, height: 0 }
                },
            })
            this.popper.show()
        }

        if (this.refDom) {
            this.queryParams(this.refDom.value)
        }
    }

    queryParams(searchKeyword) {
        if (!searchKeyword && this.refDom) {
            searchKeyword = this.refDom.value
        }

        this.searchKeyword = searchKeyword
        this.renderParamList()
    }

    getAllParams() {
        if (this.getAllCommonParams) {
            this.getAllCommonParams().then((params) => {
                const { list = [], map = {} } = params
                this.params.list = list
                this.params.map = map
                this.renderParamList()
                this.broadcast()
            })
        }
    }

    renderParamList() {
        let list = this.params.list.map((item) => ({ value: item }))

        if (!this.searchKeyword) {
            this.$vm.list = list.slice(0, 5)
        }

        let result = list.filter((item) => (item.value || '').toLowerCase().indexOf(this.searchKeyword.toLowerCase()) !== -1)

        if (!result.length) {
            this.$vm.list = []
        }

        this.$vm.list = result.slice(0, 5)

        if (!this.$vm.list || !this.$vm.list.length) {
            this.popper.hide()
        } else {
            this.refDom && this.popper.show()
        }
    }

    addParam(param) {
        if (!param || !param.name || this.hasParam(param.name)) {
            return
        }

        this.addCommonParam(param).then((res) => {
            updateApiParam(this.params, res)
            this.renderParamList()
            this.broadcast()
        })
    }

    deleteParam(key) {
        let param = this.getParamByKey(key)
        param &&
            this.deleteCommonParam(param).then(() => {
                updateApiParam(this.params, param, true)
                this.renderParamList()
                this.broadcast()
            })
    }

    getParamByKey(key) {
        return this.params.map[key]
    }

    hasParam(name) {
        return this.params.list.indexOf(name) !== -1
    }

    hide() {
        this.popper && this.popper.hide()
        this.searchKeyword = ''
        this.node = null
    }
    on(cb) {
        this._event.push(cb)
    }

    broadcast() {
        this._event.forEach((cb) => cb())
    }

    destroy() {
        this.hide()
        this.popper && this.popper.destroy()
        this.$vm && this.$vm.$destroy()
        this.$vm = null
        this.refDom = null
        this.popper = null
        this.node = null
    }
}

const updateApiParam = (state, param, isRemove = false) => {
    let { name } = param
    let { list, map } = state
    // 删除
    if (isRemove) {
        let pos = list.indexOf(name)
        if (pos !== -1) {
            list.splice(pos, 1)
            delete map[name]
        }

        return
    }
    // 添加
    if (list.indexOf(name) === -1) {
        list.unshift(name)
    }
    map[name] = param
}
