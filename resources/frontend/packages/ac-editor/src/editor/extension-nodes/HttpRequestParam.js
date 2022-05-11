import HttpRequestParamComp from './vue-node-view/HttpRequestParam.vue'
import ParamTable from './ParamTable'

export default class HttpRequestParam extends ParamTable {
    get name() {
        return 'http_api_request_parameter'
    }

    get schema() {
        return {
            content: 'inline*',
            attrs: {
                title: {
                    default: '请求参数',
                },
                request_header: {
                    default: {
                        params: [],
                        title: 'Header 请求参数',
                    },
                },
                request_body: {
                    default: {
                        params: [],
                        title: 'Body 请求参数',
                    },
                },
                request_query: {
                    default: {
                        params: [],
                        title: 'Query 请求参数',
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
        return HttpRequestParamComp
    }
}
