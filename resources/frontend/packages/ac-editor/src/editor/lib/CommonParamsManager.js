export default class CommonParamsManager {
    constructor(editor) {
        const { getAllCommonParams, addCommonParam, deleteCommonParam } = editor.options

        this.getAllCommonParams = getAllCommonParams
        this.addCommonParam = addCommonParam
        this.deleteCommonParam = deleteCommonParam

        this.params = {
            list: [],
            map: {},
        }

        this._event = []

        !editor.options.readonly && this.getAllParams()
    }

    getAllParams() {
        if (this.getAllCommonParams) {
            this.getAllCommonParams().then((params) => {
                const { map = {} } = params
                Object.keys(map).forEach((key) => {
                    updateApiParam(this.params, map[key])
                })
                this.broadcast()
            })
        }
    }

    queryParams(queryString) {
        let list = this.params.list.map((item) => ({ value: item }))

        if (!queryString) {
            return list.slice(0, 5)
        }

        let result = list.filter((item) => (item.value || '').toLowerCase().indexOf(queryString.toLowerCase()) !== -1)

        if (!result.length) {
            return []
        }

        return result.slice(0, 5)
    }

    addParam(param) {
        if (!param || !param.name || this.hasParam(param.name)) {
            return
        }

        this.addCommonParam &&
            this.addCommonParam(param).then((newParam) => {
                updateApiParam(this.params, newParam)
                this.broadcast(newParam)
            })
    }

    deleteParam(key) {
        let param = this.getParamByKey(key)
        if (param) {
            updateApiParam(this.params, param, true)
            this.deleteCommonParam && this.deleteCommonParam(param)
            this.broadcast(param)
        }
    }

    getParamByKey(key) {
        return this.params.map[key]
    }

    hasParam(name) {
        return this.params.list.indexOf(name) !== -1
    }

    on(cb) {
        this._event.push(cb)
    }

    broadcast(param) {
        this._event.forEach((cb) => cb(param))
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
