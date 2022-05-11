<script setup>
    import { onMounted, markRaw } from 'vue'
    import AcEditor from '@ac/editor'
    import { onlyTable as doc } from './doc'

    const document = markRaw(doc)

    const options = {
        getAllCommonParams: () => getAllCommonParams(),
        addCommonParam: (param) => addCommonParam(param),
        deleteCommonParam: (param) => deleteCommonParam(param),
        getUrlList: () => getUrlList(),
        deleteUrl: (id) => deleteUrl(id),
    }

    const getAllCommonParams = () =>
        new Promise((resolve) => {
            setTimeout(() => {
                resolve({
                    list: ['page', 'page_num'],
                    map: {
                        page: {
                            id: 1,
                            name: 'page',
                            type: 1,
                            type_name: 'Int',
                            is_must: false,
                            default_value: '1',
                            description: '页码',
                        },
                        page_num: {
                            id: 2,
                            name: 'page_num',
                            type: 1,
                            type_name: 'Int',
                            is_must: true,
                            default_value: '15',
                            description: '每页记录数',
                        },
                    },
                })
            }, 3000)
        })

    const addCommonParam = (param) => Promise.resolve(param)
    const deleteCommonParam = (param) => Promise.resolve(param)

    const getUrlList = () => {
        return new Promise((resolve) => {
            setTimeout(() => {
                resolve(
                    [
                        { id: 1, url: 'http://www.baidu.com' },
                        { id: 2, url: 'http://www.baidu.com2' },
                        { id: 3, url: 'http://www.baidu.com3' },
                        { id: 4, url: 'http://www.baidu.com4' },
                        { id: 5, url: 'http://www.baidu.com5' },
                        { id: 6, url: 'http://www.baidu.com6' },
                    ].map((item) => {
                        item.value = item.url
                        item.label = item.id
                        return item
                    })
                )
            }, 3000)
        })
    }

    const deleteUrl = (id) => {
        return Promise.resolve(id)
    }

    onMounted(async () => {
        // getDetail()
    })
</script>

<template>
    <button @click="getDetail">获取文档详情</button>
    <AcEditor :document="document" :options="options" />
</template>

<style>
    #app {
        font-family: Avenir, Helvetica, Arial, sans-serif;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
        margin: 40px 200px;
    }
</style>
