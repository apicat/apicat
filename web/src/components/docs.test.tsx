import { describe, it, expect } from 'vitest'
import { renderToStaticMarkup } from 'react-dom/server'
import { MemoryRouter } from 'react-router-dom'
import { Markdown, SchemaTree } from './docs'
import type { Schema } from '@/lib/types'
const schema = (data: Partial<Schema>) =>
  ({ Name: '', Type: 'string', Default: '', Example: '', ...data }) as Schema
describe('document rendering', () => {
  it('preserves union, nullable, constraints, examples and recursion markers', () => {
    const root = schema({
      Type: 'object',
      Properties: [
        schema({
          Name: 'value',
          Type: 'string | integer',
          Nullable: true,
          Minimum: 0,
          Example: '0',
          Required: true,
        }),
        schema({ Name: 'child', Recursive: true, IsRef: true, RefName: 'Node' }),
        schema({ Name: 'deep', Truncated: true }),
      ],
    })
    const html = renderToStaticMarkup(
      <MemoryRouter>
        <SchemaTree schema={root} />
      </MemoryRouter>,
    )
    for (const text of [
      'string | integer',
      'nullable',
      'minimum',
      '递归引用',
      '展示深度上限',
      '必填',
      '缺少字段描述',
    ])
      expect(html).toContain(text)
  })
  it('does not execute HTML or javascript links from spec descriptions', () => {
    const html = renderToStaticMarkup(
      <Markdown>{'<script>alert(1)</script>\n\n[bad](javascript:alert(1))'}</Markdown>,
    )
    expect(html).not.toContain('<script>')
    expect(html).not.toContain('href="javascript:')
  })
})
