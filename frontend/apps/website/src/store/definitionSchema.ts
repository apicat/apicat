import { defineStore } from 'pinia'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import { traverseTree } from '@apicat/shared'
import {
  apiAICreateSchema,
  apiCopySchema,
  apiCreateSchema,
  apiDeleteSchema,
  apiEditSchema,
  apiGetSchemaHistories,
  apiGetSchemaInfo,
  apiGetSchemaTree,
  apiMoveSchema,
} from '@/api/project/definition/schema'
import { RefPrefixKeys, SchemaTypeEnum } from '@/commons'
import { useLoading } from '@/hooks/useLoading'

interface SchemaState {
  t: any
  schemaDetail: Definition.SchemaNode | null
  schemas: Definition.SchemaNode[]
  histories: HistoryRecord.SchemaHistory[]
}

const { loadingForGetter, startLoading, endLoading } = useLoading()

export const useDefinitionSchemaStore = defineStore('project.definitionSchema', {
  state: (): SchemaState => {
    const { t } = useI18n()
    return {
      t,
      schemaDetail: null,
      schemas: [] as Definition.SchemaNode[],
      histories: [] as HistoryRecord.SchemaHistory[],
    }
  },
  getters: {
    flatSchemas: (state): Definition.SchemaNode[] => {
      const definitions: Definition.SchemaNode[] = []
      traverseTree<Definition.SchemaNode>(
        (item) => {
          if (item.type === SchemaTypeEnum.Schema) {
            ;(item.schema as Definition.SchemaNode & { title: string }).title = item.name
            definitions.push(item)
          }
          return true
        },
        state.schemas,
        { subKey: 'items' },
      )
      return definitions
    },
    definitionsForCodeGenerate: (state): Record<string, Definition.SchemaNode> => {
      const definitions: Record<string, Definition.SchemaNode> = {}
      traverseTree<Definition.SchemaNode>(
        (item) => {
          if (item.type === SchemaTypeEnum.Schema) {
            ;(item.schema as Definition.SchemaNode & { title: string }).title = item.name
            definitions[`${RefPrefixKeys.DefinitionSchema.refForCodeGeneratePrefix}${item.id}`]
              = item.schema as Definition.SchemaNode
          }
          return true
        },
        state.schemas,
        { subKey: 'items' },
      )

      return definitions
    },
    historyRecord: (state): HistoryRecord.SchemaHistoryNode[] => {
      const historyMap: Record<string, HistoryRecord.SchemaHistoryNode[]> = {}
      ;(state.histories || []).forEach((item: HistoryRecord.SchemaHistory) => {
        const categoryTitle = dayjs(item.createdAt * 1000).format('l')
        const items = (historyMap[categoryTitle] = historyMap[categoryTitle] || [])
        const node: HistoryRecord.SchemaHistoryNode = {
          id: item.id,
          title: dayjs(item.createdAt * 1000).format('LLL LT') + (item.createdBy ? `(${item.createdBy})` : ''),
        }
        items.push(node)
      })

      const tree = Object.keys(historyMap).map((key) => {
        const node: HistoryRecord.SchemaHistoryNode = {
          id: key,
          title: key,
          items: historyMap[key] || [],
        }

        return node
      })
      return tree
    },
    historyRecordForOptions(state): HistoryRecord.SchemaInfoForOptions[] {
      const options = this.historyRecord
        .reduce((result: any[], item: HistoryRecord.SchemaHistoryNode) => result.concat(item.items || []), [])
        .map(
          (item: HistoryRecord.SchemaHistoryNode) =>
            ({ id: item.id, title: item.title }) as HistoryRecord.SchemaInfoForOptions,
        )
      return [{ id: 0, title: state.t('app.historyLayout.current') }].concat(options)
    },
    isLoading: loadingForGetter,
  },
  actions: {
    async getSchemaDetail(projectID: string, schemaID: number) {
      try {
        startLoading()
        await nextTick()
        this.schemaDetail = await apiGetSchemaInfo(projectID, schemaID)
      }
      finally {
        endLoading()
      }
    },

    updateSchemaDetail(schema: Partial<Definition.SchemaNode>) {
      this.schemaDetail = {
        ...this.schemaDetail,
        ...(schema || {}),
      } as any
    },

    async getSchemas(projectID: string) {
      this.schemas = await apiGetSchemaTree(projectID)
      return this.schemas
    },

    async renameSchemaCategory(projectID: string, schema: Definition.SchemaNode) {
      await apiEditSchema(projectID, schema)
    },

    async updateSchema(projectID: string, schema: Definition.SchemaNode) {
      await apiEditSchema(projectID, schema)
      this.updateSchemaToStore(schema)
    },

    async createSchema(projectID: string, data: Omit<Definition.Schema, 'id'>) {
      return await apiCreateSchema(projectID, data)
    },

    async createSchemaWithAI(projectID: string, data: Definition.RequestCreateSchemaWithAI) {
      return await apiAICreateSchema(projectID, data)
    },

    // copy schema
    async copySchema(projectID: string, schema: Definition.SchemaNode) {
      return await apiCopySchema(projectID, schema.id)
    },

    // delete schema
    async deleteSchema(projectID: string, schema: Definition.SchemaNode, deref: boolean = false) {
      await apiDeleteSchema(projectID, schema.id, deref)
    },

    // move schema
    async moveSchema(projectID: string, data: Definition.RequestSortParams) {
      return await apiMoveSchema(projectID, data)
    },

    // 获取schema历史记录
    async getHistories(projectID: string, schemaID: number) {
      this.histories = await apiGetSchemaHistories(projectID, schemaID)
      return this.histories
    },

    updateSchemaToStore(newScehma: Partial<Definition.Schema>) {
      traverseTree<Definition.SchemaNode>(
        (schema) => {
          if (schema.id === newScehma.id) {
            Object.assign(schema, newScehma)
            return false
          }
          return true
        },
        this.schemas,
        { subKey: 'items' },
      )
    },
  },
})

export default useDefinitionSchemaStore
