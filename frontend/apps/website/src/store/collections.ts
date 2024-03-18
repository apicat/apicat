import { defineStore } from 'pinia'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import { traverseTree } from '@apicat/shared'
import {
  apiAICreateCollection,
  apiCopyCollection,
  apiCreateCollection,
  apiDeleteCollection,
  apiEditCollectionDetail,
  apiGetCollectDetail,
  apiGetCollections,
  apiMoveCollection,
  apiRenameCollection,
} from '@/api/project/collection'
import { getCollectionHistoryRecords } from '@/api/project/collectionHistoryRecord'
import { useParams } from '@/hooks/useParams'
import { getCollectionSharedToken, setCollectionSharedToken } from '@/api/shareToken'

interface CollectionsStore {
  // i18n
  t: any
  // loading: boolean
  loadingCounter: number
  collectionDetail: CollectionAPI.ResponseCollectionDetail | null
  collections: CollectionAPI.ResponseCollection[]
  histories: HistoryRecord.ResponseCollectionRecordList
  shareToken: string
}

export const useCollectionsStore = defineStore('project.collections', {
  state: (): CollectionsStore => {
    const { t } = useI18n()
    return {
      t,
      loadingCounter: 0,
      collectionDetail: null,
      collections: [],
      histories: [],
      shareToken: '',
    }
  },

  getters: {
    historyRecord: (state): HistoryRecord.CollectionTreeNode[] => {
      const historyMap: Record<string, HistoryRecord.CollectionTreeNode[]> = {}
      ;(state.histories || []).forEach((item: HistoryRecord.CollectionInfo) => {
        const categoryTitle = dayjs(item.createdAt * 1000).format('l')
        const items = (historyMap[categoryTitle] = historyMap[categoryTitle] || [])
        const node: HistoryRecord.CollectionTreeNode = {
          id: item.id,
          title: dayjs(item.createdAt * 1000).format('LLL LT') + (item.createdBy ? `(${item.createdBy})` : ''),
          type: item.type,
        }
        items.push(node)
      })

      const tree = Object.keys(historyMap).map((key) => {
        const node: HistoryRecord.CollectionTreeNode = {
          id: key,
          title: key,
          items: historyMap[key] || [],
        }

        return node
      })
      return tree
    },
    historyRecordForOptions(state): HistoryRecord.CollectionInfoForOptions[] {
      const options = this.historyRecord
        .reduce((result: any[], item: HistoryRecord.CollectionTreeNode) => result.concat(item.items || []), [])
        .map(
          (item: HistoryRecord.CollectionTreeNode) =>
            ({ id: item.id, title: item.title }) as HistoryRecord.CollectionInfoForOptions,
        )
      return [{ id: 0, title: state.t('app.historyLayout.current') }].concat(options)
    },
    loading: state => state.loadingCounter > 0,
  },
  actions: {
    // get collection detail
    async getCollectionDetail(projectID: string, collectionID: number) {
      try {
        this.loadingCounter++
        this.collectionDetail = await apiGetCollectDetail(projectID, collectionID)
      }
      catch (error) {
        //
      }
      finally {
        this.loadingCounter--
      }
      // NOTE: 鼠标点击频繁切换时，会导致loading状态不正确，渲染的数据不正确，因此改为计数方式来作为loading状态
      // this.loading = true
      // this.collectionDetail = await apiGetCollectDetail(projectID, collectionID)
      // this.loading = false
    },
    // get collection list
    async getCollections(projectID: string) {
      const { iterationID } = useParams()
      this.collections = await apiGetCollections(projectID, iterationID.value)
    },

    // rename collection
    async renameCollection(projectID: string, collection: CollectionAPI.ResponseCollection) {
      await apiRenameCollection(projectID, collection)
    },

    // update collection
    async updateCollection(projectID: string, collection: CollectionAPI.ResponseCollectionDetail) {
      await apiEditCollectionDetail(projectID, collection)
      this.updateCollectionToStore(collection)
    },

    // create collection
    async createCollection(
      projectID: string,
      data: Omit<CollectionAPI.ResponseCollectionDetail, 'id'> & { iterationID?: string },
    ) {
      const { iterationID } = useParams()
      data.iterationID = iterationID.value
      return await apiCreateCollection(projectID, data)
    },
    async createCollectionWithAI(projectID: string, data: ProjectAPI.RequestCreateCollectionWithAI) {
      const { iterationID } = useParams()
      data.iterationID = iterationID.value
      return await apiAICreateCollection(projectID, data)
    },

    // copy collection
    async copyCollection(projectID: string, collectionID: number) {
      const { iterationID } = useParams()
      return await apiCopyCollection(projectID, collectionID, { iterationID: iterationID.value })
    },

    // delete collection
    async deleteCollection(projectID: string, collectionID: number) {
      const { iterationID } = useParams()
      await apiDeleteCollection(projectID, collectionID, { iterationID: iterationID.value })
    },

    // sort collection
    async sortCollections(projectID: string, data: CollectionAPI.RequestMoveCollection) {
      await apiMoveCollection(projectID, data)
    },

    // update collection to store
    updateCollectionToStore(newCollection: Partial<CollectionAPI.ResponseCollectionDetail>) {
      traverseTree<CollectionAPI.ResponseCollection>(
        (collection) => {
          if (collection.id === newCollection.id) {
            Object.assign(collection, newCollection)
            return false
          }
          return true
        },
        this.collections,
        {
          subKey: 'items',
        },
      )
    },

    // 获取集合历史记录
    async getHistories(projectID: string, id: number) {
      this.histories = await getCollectionHistoryRecords(projectID, id)
      return this.histories
    },

    initShareToken(publicID: string) {
      this.shareToken = getCollectionSharedToken(publicID)
    },

    setShareToken(publicID: string, token: string) {
      this.shareToken = token
      setCollectionSharedToken(publicID, token)
    },
  },
})
