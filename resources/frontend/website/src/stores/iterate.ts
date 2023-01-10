import { defineStore } from 'pinia'
import { getCollectionList, sortCollectionList, addCollection, cancelCollection, getIterationDetail } from '@/api/iterate'
import { ElMessage as $Message } from 'element-plus'

interface IterateState {
    activeTab: any
    iterateInfo: any
    collection: Array<any>
    isIterateRoute: boolean
}

export const useIterateStore = defineStore({
    id: 'iterations',

    state: (): IterateState => ({
        activeTab: undefined,
        collection: [],
        iterateInfo: null,
        isIterateRoute: false,
    }),

    getters: {
        isActiveIterateList: (state) => state.activeTab === '',
    },

    actions: {
        async getIterateInfo(iteration_public_id: any) {
            if (this.iterateInfo && this.iterateInfo.id_public === iteration_public_id) {
                return this.iterateInfo
            }

            const { data } = await getIterationDetail({ iteration_id: iteration_public_id })
            this.iterateInfo = data

            return this.iterateInfo
        },

        switchActiveCollectTab(tab: any) {
            const info = this.collection.find((item: any) => item.project_id === tab)
            this.activeTab = info ? tab : ''
        },

        async getIterateCollectionList() {
            try {
                const { data } = await getCollectionList()
                this.collection = data || []
                this.switchActiveCollectTab(this.collection.length ? this.collection[0].project_id : '')
            } catch (e) {
                //
            }
        },

        async addProjectToCollect(iterate: any) {
            const { project_id, project_title: project_name } = iterate
            await addCollection({ project_id })
            $Message.success('收藏成功')
            const item = this.collection.find((item) => item.project_id === project_id)
            !item && this.collection.push({ project_id, project_name })
        },

        async removeProjectFromCollect(iterate: any) {
            const { project_id } = iterate
            await cancelCollection({ project_id })
            $Message.success('取消收藏成功')
            this.collection = this.collection.filter((item) => item.project_id !== project_id)
        },

        async sortCollectionList(oldItemIndex: number, newItemIndex: number) {
            const changeItem = this.collection.splice(oldItemIndex, 1)[0]
            this.collection.splice(newItemIndex, 0, changeItem)
            await sortCollectionList({ project_ids: this.collection.map((item) => item.project_id) })
        },
    },
})
