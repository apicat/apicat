<template>
    <div class="search-popper" @click="onStop">
        <div class="search-popper-input">
            <input type="text" v-model="form.keywords" placeholder="搜索" maxlength="255" @keyup.enter="onSearchDocument" />
            <i class="icon iconfont iconsousuo"></i>
        </div>
        <div class="search-popper-result scroll-content" @click="onStop" v-show="isShowSearchResult">
            <div v-show="isLoading" style="padding: 10px; text-align: center">搜索中...</div>
            <ul v-if="data.length">
                <li v-for="item in data" :key="item.node_id">
                    <a href="javascript:void(0);" @click="onClickSearchResult(item)" :title="item.title">{{ item.title }}</a>
                </li>
            </ul>
            <div class="result p-2" v-if="isNoData">
                <div class="result-icon">
                    <img src="@/assets/image/icon-no-result.png" alt="" />
                </div>
                <span>未找到搜索结果！</span>
            </div>
        </div>
    </div>
</template>

<script>
    import { isEmpty } from 'lodash-es'
    import { trim, debounce } from 'lodash-es'
    import { searchProjectInfo } from '@/api/preview'
    import { Storage, traverseTree } from '@ac/shared'

    export default {
        name: 'ProjectSearch',

        data() {
            return {
                data: [],
                isLoading: false,
                isNoData: false,
                isShowSearchResult: false,
                form: {
                    token: Storage.get(Storage.KEYS.SECRET_PROJECT_TOKEN + this.$route.params.project_id || '', true),
                    project_id: this.$route.params.project_id,
                    keywords: '',
                },
            }
        },

        methods: {
            onClickSearchResult(node) {
                let hasNode = false
                traverseTree((item) => {
                    if (item.id === (node.node_id || node.doc_id)) {
                        hasNode = true
                        return false
                    }
                }, [])

                if (hasNode) {
                    this.$router.push(node.link)
                } else {
                    location.href = node.href
                }
            },

            onStop(e) {
                e.stopPropagation()
            },
            onSearchDocument: debounce(function () {
                if (!isEmpty(trim(this.form.keywords)) && this.form.keywords.length >= 2) {
                    this.isLoading = true

                    this.data = []

                    this.isShowSearchResult = true

                    searchProjectInfo(this.form.token, this.form.project_id, this.form.keywords)
                        .then((res) => {
                            this.data = (res.data || []).map((item) => {
                                item.link = {
                                    name: 'preview.project.document',
                                    params: {
                                        project_id: this.form.project_id,
                                        node_id: item.node_id || item.doc_id,
                                    },
                                }
                                item.href = `/app/${this.form.project_id}/${item.node_id || item.doc_id}`
                                return item
                            })
                            this.isNoData = !this.data.length
                        })
                        .finally(() => {
                            this.isLoading = false
                        })
                }
            }, 300),
        },
        created() {
            document.body.addEventListener('click', (e) => {
                this.data = []
                this.form.keywords = ''
                this.isShowSearchResult = false
            })
        },
    }
</script>
