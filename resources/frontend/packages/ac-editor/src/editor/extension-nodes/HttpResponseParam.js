import HttpResponseParamComp from './vue-node-view/HttpResponseParam.vue'
import ParamTable from './ParamTable'

export default class HttpResponseParam extends ParamTable {
    get name() {
        return 'http_api_response_parameter'
    }

    get schema() {
        return {
            content: 'inline*',
            attrs: {
                title: {
                    default: '返回参数',
                },
                response_header: {
                    default: {
                        params: [],
                        title: '返回头部',
                    },
                },
                response_body: {
                    default: {
                        params: [],
                        title: '返回参数',
                    },
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
        return HttpResponseParamComp
    }
}
