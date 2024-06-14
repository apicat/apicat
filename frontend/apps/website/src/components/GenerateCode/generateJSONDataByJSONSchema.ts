import type { JSONSchema } from '@apicat/editor'

type Ctx = Record<string, boolean>

function _t(ctx: Ctx, o: JSONSchema, schemas: any) {
  const type = o.type
  const ps = o.properties

  if (o['x-apicat-temp-prop'])
    return undefined

  let data: any
  switch (type) {
    case 'string':
      data = 'string'
      break
    case 'number':
    case 'integer':
      data = 0
      break
    case 'boolean':
      data = false
      break
    case 'null':
      data = null
      break
    case 'any':
      data = null
      break
    case 'array':
      data = [_t(ctx, o.items as JSONSchema, schemas)]
      break
    case 'object':
      data = {}
      if (ps) {
        for (const key in ps)
          data[key] = _t(ctx, ps[key], schemas)
      }
      break
    default:
      if (Object.prototype.hasOwnProperty.call(o, 'anyOf')) {
        data = {}
        for (const subSchema of o.anyOf!) {
          const subData = _t(ctx, subSchema, schemas)
          Object.assign(data, subData)
        }
      }
      else if (Object.prototype.hasOwnProperty.call(o, 'allOf')) {
        data = {}
        for (const subSchema of o.allOf!) {
          const subData = _t(ctx, subSchema, schemas)
          Object.assign(data, subData)
        }
      }
      else if (Object.prototype.hasOwnProperty.call(o, 'oneOf')) {
        const validSubSchemas = o.oneOf!.filter((subSchema) => {
          try {
            const subData = _t(ctx, subSchema, schemas)
            return subData !== undefined
          }
          catch (e) {
            return false
          }
        })
        if (validSubSchemas.length > 0)
          data = _t(ctx, validSubSchemas[0], schemas)
        else
          data = undefined
      }
      else if (Object.prototype.hasOwnProperty.call(o, '$ref')) {
        const uri = o.$ref!
        const id = Number.parseInt(uri.split('/').pop()!)
        if (ctx[id]) {
          data = {}
        }
        else {
          const ref = (schemas || []).find((d: any) => d.id === id)
          if (ref) {
            ctx = JSON.parse(JSON.stringify(ctx))
            ctx[ref.id] = true
            data = _t(ctx, ref.schema, schemas)
          }
          else {
            data = null
          }
        }
      }
      else {
        data = null
      }
  }
  return data
}

export function generateJSONDataByJSONSchema(o: JSONSchema, schemas: any, id: number) {
  const ctx: Ctx = {}
  ctx[id!] = true
  const data = _t(ctx, o, schemas)
  return data
}
