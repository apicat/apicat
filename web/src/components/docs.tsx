import { lintDetail } from '@/lib/lint-messages'
import { useI18n } from '@/lib/i18n'
import { Fragment, useState, type ReactNode } from 'react'
import { Link } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { ChevronDown, ChevronRight, Copy, Check, FileCode2 } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import type { Schema, Parameter, Media, RequestBody, Endpoint } from '@/lib/types'

export function Markdown({ children }: { children?: string }) {
  return children ? (
    <div className="markdown">
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{children}</ReactMarkdown>
    </div>
  ) : null
}
export function Panel({
  title,
  children,
  action,
  className,
}: {
  title?: string
  children: ReactNode
  action?: ReactNode
  className?: string
}) {
  return (
    <section className={cn('rounded-xl border bg-card shadow-panel', className)}>
      {title && (
        <div className="flex items-center justify-between border-b px-5 py-3.5">
          <h2 className="text-base font-semibold">{title}</h2>
          {action}
        </div>
      )}
      <div className="p-5">{children}</div>
    </section>
  )
}
export function Method({ method }: { method: string }) {
  return (
    <span
      data-method={method}
      className="http-method inline-flex min-w-12 shrink-0 justify-center rounded-md px-1.5 py-0.5 font-mono text-[10px] font-semibold tracking-wide"
    >
      {method}
    </span>
  )
}
export function Code({ value, label }: { value: string; label?: string }) {
  const { t } = useI18n()
  label ??= t('代码')
  const [copied, setCopied] = useState(false)
  const [error, setError] = useState(false)
  async function copy() {
    try {
      await navigator.clipboard.writeText(value)
      setCopied(true)
      setError(false)
      setTimeout(() => setCopied(false), 1500)
    } catch {
      setError(true)
    }
  }
  return (
    <div className="relative overflow-hidden rounded-lg border bg-muted/60">
      <div className="flex items-center justify-between border-b px-3 py-1 text-xs text-muted-foreground">
        <span>{label}</span>
        <Button
          size="sm"
          variant="ghost"
          onClick={copy}
          aria-label={t('复制 {label}', { label })}
          className="h-7"
        >
          {copied ? <Check size={13} /> : <Copy size={13} />}
          {copied ? t('已复制') : t('复制')}
        </Button>
      </div>
      <pre className="max-h-[520px] overflow-auto p-4 text-xs leading-6">
        <code>{value}</code>
      </pre>
      {error && (
        <p role="status" className="p-2 text-xs text-destructive">
          {t('无法访问剪贴板，请手动选择文本复制。')}
        </p>
      )}
    </div>
  )
}
export function Empty({
  title,
  description,
  children,
}: {
  title: string
  description?: string
  children?: ReactNode
}) {
  return (
    <div className="flex flex-col items-center justify-center rounded-xl border border-dashed px-8 py-16 text-center">
      <FileCode2 className="mb-4 text-muted-foreground" size={30} />
      <h2 className="text-base font-semibold">{title}</h2>
      <p className="mt-2 max-w-md text-sm text-muted-foreground">{description}</p>
      {children}
    </div>
  )
}
function TypeLabel({ schema }: { schema: Schema }) {
  const { t } = useI18n()
  return (
    <div className="space-y-1">
      <code className="text-xs text-info">
        {schema.Type || (schema.Properties?.length ? 'object' : '—')}
        {schema.Format && `<${schema.Format}>`}
      </code>
      {schema.Nullable && <div className="text-[10px] text-muted-foreground">nullable</div>}
      {schema.Deprecated && <Badge variant="outline">{t('已弃用')}</Badge>}
      {schema.IsRef && (
        <Link
          className="block max-w-44 truncate text-xs text-info underline"
          to={`/components?type=schemas&name=${encodeURIComponent(schema.RefName)}`}
        >
          {schema.RefName}
        </Link>
      )}
    </div>
  )
}
export function SchemaTree({ schema }: { schema: Schema | null }) {
  const { t } = useI18n()
  const [collapsed, setCollapsed] = useState<Set<string>>(new Set())
  if (!schema) return null
  function rows(s: Schema, name: string, key: string, depth: number): ReactNode {
    const children: { s: Schema; name: string; key: string }[] = []
    s.Properties?.forEach((p, i) => children.push({ s: p, name: p.Name, key: `${key}.p${i}` }))
    if (s.Items) children.push({ s: s.Items, name: 'items[]', key: `${key}.items` })
    for (const kind of ['AllOf', 'OneOf', 'AnyOf'] as const)
      s[kind]?.forEach((p, i) =>
        children.push({
          s: p,
          name: `${kind[0].toLowerCase()}${kind.slice(1)} · ${i + 1}`,
          key: `${key}.${kind}${i}`,
        }),
      )
    const open = !collapsed.has(key)
    const constraints = [
      ['minLength', s.MinLength],
      ['maxLength', s.MaxLength],
      ['minimum', s.Minimum],
      ['maximum', s.Maximum],
      ['minItems', s.MinItems],
      ['maxItems', s.MaxItems],
      ['pattern', s.Pattern],
    ].filter(([, value]) => value !== null && value !== undefined && value !== '')
    return (
      <Fragment key={key}>
        <tr className="border-b last:border-0">
          <td className="py-3 pr-3 align-top">
            <div className="flex items-start gap-1.5" style={{ paddingLeft: depth * 16 }}>
              {children.length > 0 ? (
                <button
                  aria-label={t('展开或折叠 {name}', { name })}
                  aria-expanded={open}
                  className="mt-0.5 shrink-0 rounded hover:bg-accent"
                  onClick={() =>
                    setCollapsed((prev) => {
                      const next = new Set(prev)
                      next.has(key) ? next.delete(key) : next.add(key)
                      return next
                    })
                  }
                >
                  {open ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
                </button>
              ) : (
                <span className="w-3.5 shrink-0" />
              )}
              <code className="break-all text-xs font-medium">{name || 'root'}</code>
            </div>
          </td>
          <td className="p-3 align-top">
            <TypeLabel schema={s} />
          </td>
          <td className="p-3 align-top text-xs">
            {s.Required ? <span className="text-secondary-text">{t('必填')}</span> : '—'}
          </td>
          <td className="min-w-52 p-3 align-top text-xs">
            <Markdown>{s.Description}</Markdown>
            {s.Name && !s.Description && <span className="text-warning">{t('缺少字段描述')}</span>}
            {s.Recursive && <p className="text-muted-foreground">{t('递归引用')}</p>}
            {s.Truncated && <p className="text-muted-foreground">{t('已达到展示深度上限')}</p>}
            {s.Enum?.length ? (
              <p className="mt-1 break-all text-muted-foreground">enum: {s.Enum.join(' · ')}</p>
            ) : null}
            {constraints.map(([k, v]) => (
              <span key={String(k)} className="mr-2 text-muted-foreground">
                {k}: {v}
              </span>
            ))}
          </td>
          <td className="max-w-52 p-3 align-top">
            <pre className="whitespace-pre-wrap break-all text-xs">{s.Example || '—'}</pre>
            {s.Default !== '' && (
              <p className="mt-2 text-xs text-muted-foreground">default: {s.Default}</p>
            )}
          </td>
        </tr>
        {open && children.map((c) => rows(c.s, c.name, c.key, depth + 1))}
      </Fragment>
    )
  }
  return (
    <div className="overflow-x-auto">
      <table className="w-full min-w-[650px] text-left">
        <thead className="border-b text-[11px] text-muted-foreground">
          <tr>
            <th className="pb-3 font-medium">{t('字段')}</th>
            <th className="px-3 pb-3 font-medium">{t('类型')}</th>
            <th className="px-3 pb-3 font-medium">{t('必填')}</th>
            <th className="px-3 pb-3 font-medium">{t('说明 / 约束')}</th>
            <th className="px-3 pb-3 font-medium">{t('示例 / 默认值')}</th>
          </tr>
        </thead>
        <tbody>{rows(schema, schema.Name || 'root', 'root', 0)}</tbody>
      </table>
    </div>
  )
}
export function Parameters({ parameters }: { parameters: Parameter[] | null }) {
  const { t } = useI18n()
  if (!parameters?.length) return null
  return (
    <div className="space-y-5">
      {parameters.map((p) => (
        <div key={`${p.In}:${p.Name}`} className="rounded-lg border p-4">
          <div className="mb-2 flex items-center gap-2">
            <code className="font-semibold">{p.Name}</code>
            <Badge variant="secondary">{p.In}</Badge>
            {p.Required && <Badge variant="outline">{t('必填')}</Badge>}
            {p.Deprecated && <Badge variant="outline">{t('已弃用')}</Badge>}
          </div>
          <div className="mb-3 text-sm text-muted-foreground">
            <Markdown>{p.Description}</Markdown>
            {!p.Description && <span className="text-warning">{t('缺少参数描述')}</span>}
          </div>
          <SchemaTree schema={p.Schema} />
        </div>
      ))}
    </div>
  )
}
export function MediaContent({ content }: { content: Media[] | null }) {
  const { t } = useI18n()
  return (
    <div className="space-y-4">
      {content?.map((mt) => (
        <div key={mt.MediaType}>
          <Badge variant="secondary" className="mb-3 font-mono">
            {mt.MediaType}
          </Badge>
          <SchemaTree schema={mt.Schema} />
          {mt.Example && (
            <div className="mt-3">
              <Code value={formatJSON(mt.Example)} label={t('示例')} />
            </div>
          )}
        </div>
      ))}
    </div>
  )
}
export function BodyDoc({ body }: { body: RequestBody | null }) {
  const { t } = useI18n()
  return body ? (
    <div className="space-y-3">
      <div className="flex gap-2">
        {body.Required && <Badge variant="outline">{t('必填')}</Badge>}
        <Markdown>{body.Description}</Markdown>
      </div>
      <MediaContent content={body.Content} />
    </div>
  ) : null
}
export function EndpointDoc({ endpoint }: { endpoint: Endpoint }) {
  const { t, locale } = useI18n()
  return (
    <div className="space-y-5">
      {endpoint.LintHints?.length ? (
        <div className="rounded-lg border border-warning/25 bg-warning/5 px-4 py-3 text-xs text-warning">
          <strong>{t('文档质量提示')}</strong>
          <ul className="mt-2 list-inside list-disc space-y-1">
            {endpoint.LintDetails?.map((h, i) => (
              <li key={i}>{lintDetail(locale, h)}</li>
            ))}
          </ul>
        </div>
      ) : null}
      <Markdown>{endpoint.Description}</Markdown>
      {endpoint.Parameters?.length ? (
        <Panel title={t('请求参数')}>
          <Parameters parameters={endpoint.Parameters} />
        </Panel>
      ) : null}
      {endpoint.RequestBody && (
        <Panel title={t('请求体')}>
          <BodyDoc body={endpoint.RequestBody} />
        </Panel>
      )}
      {endpoint.Responses?.map((r) => (
        <Panel key={r.StatusCode} title={t('响应 {status}', { status: r.StatusCode })}>
          <div className="mb-4 text-muted-foreground">
            <Markdown>{r.Description}</Markdown>
          </div>
          <MediaContent content={r.Content} />
        </Panel>
      ))}
    </div>
  )
}
export function formatJSON(value: string) {
  try {
    return JSON.stringify(JSON.parse(value), null, 2)
  } catch {
    return value
  }
}
