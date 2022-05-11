import ParamTableComp from './vue-node-view/ParamTable.vue'
import Node from '../lib/Node'

export default class ParamTable extends Node {
    get name() {
        return 'api_parameter'
    }

    get schema() {
        return {
            content: 'inline*',
            attrs: {
                params: {
                    default: [],
                },
            },
            group: 'block',
            atom: true,
            isolating: true,
            marks: '',
            selectable: true,
            draggable: false,
            parseDOM: [
                {
                    tag: `div[data-node-type='${this.name}']`,
                },
            ],
            toDOM: () => [
                'div',
                {
                    'data-node-type': this.name,
                },
            ],
        }
    }

    get component() {
        return ParamTableComp
    }

    stopEvent(event) {
        let node = event.target

        // 拖拽句柄时，不阻止事件。
        if (node.classList.contains('handle')) {
            return false
        }

        const isClickEvent = event.type === 'mousedown'
        const isDragEvent = event.type.startsWith('drag') || event.type === 'drop'
        const isDragEnd = event.type === 'dragend'

        if (isClickEvent && node.classList.contains('drag_btn')) {
            this.isDragSortHandler = true
        }

        if (isDragEnd) {
            this.isDragSortHandler = false
        }

        if (this.isDragSortHandler && isDragEvent) {
            return true
        }

        if (isDragEvent) {
            return false
        }

        let isDragTableHeader = false

        while (node && node.parentNode) {
            if (node && node.classList.contains('SortableTreeTableHeader')) {
                isDragTableHeader = true
                break
            }
            node = node.parentNode
        }

        return !isDragTableHeader
    }
}
