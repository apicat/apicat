import Node from '../lib/Node'
import HttpUrlComp from './vue-node-view/HttpUrl.vue'
import UpdateHttpUrlAttrs from '../../ui-components/editAttrs/UpdateHttpUrlAttrs.vue'
import { capitalize } from 'lodash-es'

export default class HttpUrl extends Node {
    get name() {
        return 'http_api_url'
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
                method: {
                    default: 1,
                },
                bodyDataType: {
                    default: 0,
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
                    getAttrs(dom) {
                        let url = dom.getAttribute('data-url') || ''
                        let path = dom.getAttribute('data-path') || ''
                        let method = parseInt(dom.getAttribute('data-method'), 10)
                        let bodyDataType = parseInt(dom.getAttribute('data-body-type'), 10)

                        method = isNaN(method) ? 1 : method
                        bodyDataType = isNaN(bodyDataType) ? 1 : bodyDataType

                        return {
                            url,
                            path,
                            method,
                            bodyDataType,
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
        return HttpUrlComp
    }

    get edit() {
        return UpdateHttpUrlAttrs
    }
}
