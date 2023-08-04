import { getDefinitionSchemaList, updateDefinitionSchema, createDefinitionSchema, copyDefinitionSchema, deleteDefinitionSchema } from '@/api/definitionSchema'
import { DefinitionTypeEnum, RefPrefixKeys, markDataWithKey, removeJsonSchemaTempProperty } from '@/commons'
import { DefinitionSchema, JSONSchema } from '@/components/APIEditor/types'
import { traverseTree } from '@apicat/shared'
import { cloneDeep } from 'lodash-es'
import { defineStore } from 'pinia'

export const extendDocTreeFeild = (node = {} as any) => {
  node = node || {}
  node._extend = {
    isLeaf: node.type !== DefinitionTypeEnum.DIR,
    isEditable: false,
    isCurrent: false,
  }
  markDataWithKey(node.schema)
  return node
}

export const useDefinitionStore = defineStore('definitionSchema', {
  state: () => ({
    definitions: [] as DefinitionSchema[],
    tempCreateSchemaParentId: undefined as number | undefined,
  }),
  getters: {
    definitionsForCodeGenerate: (state): Record<string, JSONSchema> => {
      const definitions: Record<string, JSONSchema> = {}

      state.definitions.forEach((item: DefinitionSchema) => {
        ;(item.schema as JSONSchema & { title: string }).title = item.name
        definitions[`${RefPrefixKeys.DefinitionSchema.refForCodeGeneratePrefix}${item.id}`] = item.schema
      })

      return definitions
    },
  },
  actions: {
    transformSchemaForCodeGenerate(definitionSchema: DefinitionSchema): string {
      try {
        const schema = definitionSchema.schema || {}
        let json = JSON.stringify({ ...schema, description: definitionSchema.description || schema.description, definitions: this.definitionsForCodeGenerate })
        json = json.replaceAll(RefPrefixKeys.DefinitionSchema.key, RefPrefixKeys.DefinitionSchema.replaceForCodeGenerate)
        return json
      } catch (error) {
        return JSON.stringify({ type: 'object' })
      }
    },

    async getDefinitions(project_id: string) {
      const tree = await getDefinitionSchemaList(project_id)
      this.definitions = traverseTree((item: any) => extendDocTreeFeild(item), tree) as any
      return this.definitions
    },

    async updateDefinition(data: any) {
      const copy = cloneDeep(data)
      const jsonSchema = copy.schema as JSONSchema
      removeJsonSchemaTempProperty(jsonSchema)
      await updateDefinitionSchema(copy)
      data.id = data.def_id
      this.updateDefinitionStore(data)
    },

    async createDefinition(data: any) {
      try {
        const definition: any = await createDefinitionSchema(data)
        this.definitions.unshift(extendDocTreeFeild(definition))
        return definition
      } catch (error) {
        // error
      }
    },

    async copyDefinition(project_id: string, def_id: string | number) {
      const definition: any = await copyDefinitionSchema(project_id, def_id)
      const index = this.definitions.findIndex((item: any) => item.id === def_id)
      if (index !== -1) {
        this.definitions.splice(index + 1, 0, extendDocTreeFeild(definition))
      }
      return definition
    },

    updateDefinitionStore(definition: any) {
      const { id, name, description, schema } = definition
      const target = this.definitions.find((item: any) => item.id === id)
      if (target) {
        target.name = name
        target.description = description
        target.schema = { ...target.schema, ...schema }
      }
    },

    async deleteDefinition(project_id: string, def_id: string | number, is_unref = 1) {
      await deleteDefinitionSchema(project_id as string, def_id, is_unref)
    },
  },
})

export const useDefinitionSchemaStore = useDefinitionStore

export default useDefinitionStore
