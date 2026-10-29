import { describe, it, expect } from 'vitest'
import { executionDraft, initialDraft, sample, serializeQuery } from './request'
import type { Endpoint, Parameter, Schema } from './types'
const schema = (data: Partial<Schema>) =>
  ({ Type: 'string', Default: '', Example: '', JSONExample: '', ...data }) as Schema
const parameter = (data: Partial<Parameter>) =>
  ({
    Name: 'tags',
    In: 'query',
    Style: '',
    Explode: null,
    Schema: schema({ Type: 'array' }),
    ...data,
  }) as Parameter
const endpoint = (data: Partial<Endpoint>) =>
  ({ Method: 'GET', Path: '/pets/{id}', Parameters: [], ...data }) as Endpoint
describe('request construction', () => {
  it('uses typed JSON examples and request media examples', () => {
    expect(sample(schema({ Type: 'object', JSONExample: '{"n":3}' }))).toEqual({ n: 3 })
    const draft = initialDraft(
      endpoint({
        Parameters: [
          parameter({ Name: 'id', In: 'path', Required: true, Example: '"a/b"' }),
          parameter({ Name: 'optional', Required: false, Schema: schema({}) }),
        ],
        RequestBody: {
          Required: true,
          Description: '',
          Content: [{ MediaType: 'application/json', Schema: null, Example: '{"name":"cat"}' }],
        },
      }),
    )
    expect(draft.path).toEqual([{ key: 'id', value: 'a/b', enabled: true }])
    expect(draft.query[0].enabled).toBe(false)
    expect(JSON.parse(draft.body)).toEqual({ name: 'cat' })
  })
  it('serializes repeated and delimited arrays', () => {
    const row = { key: 'tags', value: '["one","two"]', enabled: true }
    expect(serializeQuery(row, parameter({})).map((p) => p.value)).toEqual(['one', 'two'])
    expect(serializeQuery(row, parameter({ Explode: false }))[0].value).toBe('one,two')
    expect(serializeQuery(row, parameter({ Style: 'pipeDelimited' }))[0].value).toBe('one|two')
  })
  it('serializes objects and prevents invalid nested query serialization', () => {
    const p = parameter({ Style: 'deepObject', Schema: schema({ Type: 'object' }) })
    expect(serializeQuery({ key: 'filter', value: '{"name":"cat"}', enabled: true }, p)).toEqual([
      { key: 'filter[name]', value: 'cat', enabled: true },
    ])
    expect(() =>
      serializeQuery({ key: 'filter', value: '{"nested":{"x":1}}', enabled: true }, p),
    ).toThrow('嵌套')
  })
  it('expands variables before array serialization without changing draft', () => {
    const ep = endpoint({ Parameters: [parameter({})] })
    const draft = initialDraft(ep)
    draft.query = [{ key: 'tags', value: '{{tags}}', enabled: true }]
    expect(executionDraft(draft, ep, { tags: '[1,2]' }).query.map((p) => p.value)).toEqual([
      '1',
      '2',
    ])
    expect(draft.query[0].value).toBe('{{tags}}')
  })
  it('does not invent complex path serialization', () => {
    const ep = endpoint({ Parameters: [parameter({ Name: 'id', In: 'path', Required: true })] })
    expect(() => executionDraft(initialDraft(ep), ep)).toThrow('复杂序列化')
  })
})
