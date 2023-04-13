import { DefinitionTypeEnum } from '@/commons/constant'
import { Definition, JSONSchema } from '@/components/APIEditor/types'

export const createDefaultSchema = (overwrite?: Partial<JSONSchema>) => ({
  type: 'object',
  properties: {},
  required: [],
  'x-apicat-orders': [],
  example: '',
  ...overwrite,
})

export const createDefaultDefinition = (overwrite?: Partial<Definition>) => ({
  name: '',
  description: '',
  parent_id: 0,
  type: DefinitionTypeEnum.SCHEMA,
  schema: createDefaultSchema(),
  ...overwrite,
})

export default createDefaultDefinition
