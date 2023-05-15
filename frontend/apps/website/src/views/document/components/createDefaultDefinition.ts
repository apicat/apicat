import { DefinitionTypeEnum } from '@/commons/constant'
import { DefinitionSchema, JSONSchema } from '@/components/APIEditor/types'
import { DefinitionResponse } from '@/typings'

export const createDefaultSchema = (overwrite?: Partial<JSONSchema>) => ({
  type: 'object',
  properties: {},
  required: [],
  'x-apicat-orders': [],
  example: '',
  ...overwrite,
})

export const createDefaultSchemaDefinition = (overwrite?: Partial<DefinitionSchema>) => ({
  name: '',
  description: '',
  parent_id: 0,
  type: DefinitionTypeEnum.SCHEMA,
  schema: createDefaultSchema(),
  ...overwrite,
})

export const createDefaultResponseContent = () => ({
  'application/json': {
    schema: createDefaultSchema(),
  },
})

export const createDefaultResponseDefinition = (overwrite?: Partial<DefinitionResponse>) => ({
  name: '',
  description: '',
  parent_id: 0,
  type: DefinitionTypeEnum.RESPONSE,
  header: [],
  content: createDefaultResponseContent(),
  ...overwrite,
})

export default createDefaultSchemaDefinition
