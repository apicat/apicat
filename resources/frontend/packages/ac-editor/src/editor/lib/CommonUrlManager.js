import { ref } from 'vue'

export default class CommonUrlManager {
    constructor(editor) {
        const { getUrlList, deleteUrl } = editor.options

        this.getUrlList = getUrlList
        this.deleteUrl = deleteUrl

        this.urls = []
        !editor.options.readonly && this.initUrlList()

        this.urlList = ref([])
    }

    initUrlList() {
        if (this.getUrlList) {
            this.getUrlList()
                .then((urls) => {
                    this.urls = urls || []
                })
                .catch((e) => {
                    //
                })
        }
    }

    filterUrls(queryString) {
        if (!queryString) {
            this.urlList.value = []
            return
        }

        let list = (this.urls || []).map((item) => {
            item.value = item.url
            return item
        })

        this.urlList.value = list.filter((item) => (item.value || '').toLowerCase().indexOf(queryString.toLowerCase()) !== -1)
    }

    queryUrls(queryString) {
        if (!queryString) {
            return []
        }

        let list = (this.urls || []).map((item) => {
            item.value = item.url
            return item
        })

        let result = list.filter((item) => (item.value || '').toLowerCase().indexOf(queryString.toLowerCase()) !== -1)

        if (!result.length) {
            return []
        }

        return result
    }

    deleteUrlById(id) {
        let url = this.getUrlById(id)
        if (url) {
            this.urls = this.urls.filter((item) => item.id !== id)
            this.urlList.value = this.urls.concat([])
            this.deleteUrl && this.deleteUrl(id)
        }
    }

    getUrlById(id) {
        return this.urls.find((item) => item.id === id)
    }
}
