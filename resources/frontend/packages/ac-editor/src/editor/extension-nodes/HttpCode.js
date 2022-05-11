import Node from '../lib/Node'
import HttpCodeComp from './vue-node-view/HttpCode.vue'
import UpdateHttpCodeAttrs from '../../ui-components/editAttrs/UpdateHttpCodeAttrs.vue'
import { capitalize } from 'lodash-es'

export default class HttpCode extends Node {
    get name() {
        return 'http_status_code'
    }

    get schema() {
        return {
            content: 'inline*',
            attrs: {
                intro: {
                    default: 'Response Status Code:',
                },
                code: {
                    default: 200,
                },
                codeDesc: {
                    default: 'OK',
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
                        let intro = dom.getAttribute('data-intro') || ''
                        let code = dom.getAttribute('data-code') || 200
                        let codeDesc = dom.getAttribute('data-code-desc') || ''
                        return {
                            intro,
                            code,
                            codeDesc,
                        }
                    },
                },
            ],
            toDOM: () => ['div'],
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
        return HttpCodeComp
    }

    get edit() {
        return UpdateHttpCodeAttrs
    }
}
