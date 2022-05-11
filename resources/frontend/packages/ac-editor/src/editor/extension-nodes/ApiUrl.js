import Node from '../lib/Node'
import ApiUrlComponent from '../extension-nodes/vue-node-view/ApiUrl.vue'
import UpdateApiUrlAttrs from '../../ui-components/editAttrs/UpdateApiUrlAttrs.vue'
import { capitalize } from 'lodash-es'

export default class ApiUrl extends Node {
    get name() {
        return 'api_url'
    }

    get schema() {
        return {
            content: 'inline*',
            attrs: {
                url: {
                    default: '',
                },
                path: {
                    default: '',
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
                    getAttrs: (dom) => {
                        let url = dom.getAttribute('data-url') || null
                        let func = dom.getAttribute('data-path') || null
                        return {
                            url,
                            func,
                        }
                    },
                },
            ],
            toDOM: () => ['div', 0],
        }
    }

    commands({ type }) {
        return {
            [`create${capitalize(this.name)}`]: () => () => {
                this.editor.nodeEditViewManager.create(type.name)
            },
        }
    }

    get component() {
        return ApiUrlComponent
    }

    get edit() {
        return UpdateApiUrlAttrs
    }
}
