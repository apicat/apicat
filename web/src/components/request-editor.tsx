import { useEffect, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import { Send, Square, RotateCcw, Clock3, ArrowDown, CircleCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { PairEditor } from './pair-editor'
import { Code, Empty, formatJSON } from './docs'
import { api } from '@/lib/api'
import { useWorkbench } from '@/lib/workbench'
import { executionDraft, initialDraft } from '@/lib/request'
import type { Endpoint, RequestDraft, RunResponse } from '@/lib/types'

export function RequestEditor({ endpoint }: { endpoint: Endpoint }) {
  const { workspace, drafts, setDraft } = useWorkbench()
  const key = endpoint.Method + ' ' + endpoint.Path
  const draft = drafts[key] || initialDraft(endpoint)
  const update = (data: Partial<RequestDraft>) => setDraft(key, { ...draft, ...data })
  const [response, setResponse] = useState<RunResponse | null>(null)
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  const authorization = draft.headers.find(
    (h) => h.key.toLowerCase() === 'authorization' && h.enabled,
  )
  const [auth, setAuth] = useState(authorization?.value.startsWith('Bearer ') ? 'bearer' : 'none')
  const token = authorization?.value.replace(/^Bearer /, '') || ''
  const setToken = (value: string) =>
    update({
      headers: [
        ...draft.headers.filter((h) => h.key.toLowerCase() !== 'authorization'),
        { key: 'Authorization', value: `Bearer ${value}`, enabled: true },
      ],
    })
  const controller = useRef<AbortController | null>(null)
  useEffect(() => () => controller.current?.abort(), [])
  const environment = workspace.environments.find((e) => e.id === workspace.activeEnvironment)
  const variables = Object.fromEntries(
    [...workspace.variables, ...(environment?.variables || [])]
      .filter((v) => v.enabled)
      .map((v) => [v.key, v.value]),
  )
  async function send() {
    setError('')
    setResponse(null)
    setBusy(true)
    const abort = new AbortController()
    controller.current = abort
    try {
      const input = executionDraft(
        { ...draft, environmentId: workspace.activeEnvironment },
        endpoint,
        variables,
      )
      if (auth === 'bearer')
        input.headers = [
          ...input.headers.filter((h) => h.key.toLowerCase() !== 'authorization'),
          { key: 'Authorization', value: `Bearer ${token}`, enabled: true },
        ]
      const contentType =
        input.headers.find((h) => h.enabled && h.key.toLowerCase() === 'content-type')?.value || ''
      if (contentType.includes('multipart/form-data'))
        throw new Error('此版本暂不支持 multipart 文件上传；可使用 JSON、文本或 URL 编码表单。')
      setResponse(
        await api<RunResponse>('/requests/execute', {
          method: 'POST',
          body: JSON.stringify(input),
          signal: abort.signal,
        }),
      )
    } catch (e) {
      setError((e as Error).name === 'AbortError' ? '请求已取消' : (e as Error).message)
    } finally {
      if (controller.current === abort) {
        setBusy(false)
        controller.current = null
      }
    }
  }
  return (
    <div className="grid items-start gap-5 min-[1440px]:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)]">
      <div className="min-w-0 rounded-xl border bg-card p-5 shadow-panel">
        <div className="mb-4 flex flex-wrap items-center justify-between gap-2 text-xs text-muted-foreground">
          <span className="min-w-0 break-all">
            {environment
              ? `当前环境：${environment.name} · ${environment.baseUrl || '尚未设置地址'}`
              : '未选择环境，请输入完整地址或配置环境'}
          </span>
          <Link to="/environments" className="text-brand-text hover:underline">
            管理环境 →
          </Link>
        </div>
        <div className="flex flex-wrap gap-2">
          <select
            aria-label="请求方法"
            className="h-10 rounded-md border bg-background px-2 font-mono text-xs font-semibold"
            value={draft.method}
            onChange={(e) => update({ method: e.target.value })}
          >
            {['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS', 'TRACE'].map((m) => (
              <option key={m}>{m}</option>
            ))}
          </select>
          <Input
            aria-label="请求地址"
            value={draft.url}
            onChange={(e) => update({ url: e.target.value })}
            className="h-10 min-w-48 flex-1 font-mono text-xs"
            placeholder="/users 或 https://api.example.com/users"
          />
          <Button className="h-10" onClick={busy ? () => controller.current?.abort() : send}>
            {busy ? <Square size={15} /> : <Send size={15} />}
            {busy ? '取消' : '发送请求'}
          </Button>
        </div>
        {!environment && !!endpoint.Servers?.length && (
          <div className="mt-3 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
            <span>文档服务器：</span>
            {endpoint.Servers.map((s) => (
              <button
                key={s.URL}
                className="rounded border px-2 py-1 text-brand-text"
                onClick={() => update({ url: s.URL.replace(/\/$/, '') + endpoint.Path })}
              >
                {s.Description || s.URL}
              </button>
            ))}
          </div>
        )}
        <Tabs defaultValue="query" className="mt-5">
          <TabsList variant="line" className="h-auto! max-w-full flex-wrap">
            <TabsTrigger value="query">
              查询参数 <span className="ml-1 text-[10px]">{draft.query.length}</span>
            </TabsTrigger>
            <TabsTrigger value="path">路径参数</TabsTrigger>
            <TabsTrigger value="headers">请求头</TabsTrigger>
            <TabsTrigger value="body">请求体</TabsTrigger>
            <TabsTrigger value="auth">认证</TabsTrigger>
            <TabsTrigger value="settings">选项</TabsTrigger>
          </TabsList>
          <TabsContent value="query" className="pt-4">
            <PairEditor
              rows={draft.query}
              onChange={(query) => update({ query })}
              label="查询参数"
            />
            <p className="mt-3 text-xs text-muted-foreground">
              数组和对象参数使用 JSON 值，按文档中的 style / explode 编码。重复参数可添加多行。
            </p>
          </TabsContent>
          <TabsContent value="path" className="pt-4">
            <PairEditor rows={draft.path} onChange={(path) => update({ path })} label="路径参数" />
          </TabsContent>
          <TabsContent value="headers" className="pt-4">
            <PairEditor
              rows={draft.headers}
              onChange={(headers) => update({ headers })}
              label="请求头"
            />
          </TabsContent>
          <TabsContent value="body" className="pt-4">
            <div className="mb-2 flex justify-between text-xs text-muted-foreground">
              <span>原始请求体 · Content-Type 在请求头中设置</span>
              <Button
                size="sm"
                variant="ghost"
                onClick={() => update({ body: formatJSON(draft.body) })}
              >
                格式化 JSON
              </Button>
            </div>
            <Textarea
              spellCheck={false}
              aria-label="请求体"
              className="min-h-64 font-mono text-xs leading-6"
              value={draft.body}
              onChange={(e) => update({ body: e.target.value })}
              placeholder="输入 JSON、文本或 key=value&key2=value2"
            />
          </TabsContent>
          <TabsContent value="auth" className="space-y-4 pt-4">
            <select
              aria-label="认证方式"
              value={auth}
              onChange={(e) => setAuth(e.target.value)}
              className="h-9 rounded border bg-background px-3 text-sm"
            >
              <option value="none">使用请求头中的认证</option>
              <option value="bearer">Bearer Token</option>
            </select>
            {auth === 'bearer' && (
              <Input
                type="password"
                aria-label="Bearer Token"
                autoComplete="off"
                placeholder="Token 或 {{token}}"
                value={token}
                onChange={(e) => setToken(e.target.value)}
              />
            )}
            <p className="text-xs text-muted-foreground">
              其他认证可直接编辑 Authorization、Cookie 或 API Key 请求头。
            </p>
          </TabsContent>
          <TabsContent value="settings" className="space-y-3 pt-4">
            <label className="flex items-center gap-3 text-sm">
              超时（毫秒）
              <Input
                type="number"
                min={100}
                max={120000}
                className="w-36"
                value={draft.timeoutMs}
                onChange={(e) => update({ timeoutMs: Number(e.target.value) })}
              />
            </label>
            <p className="text-xs text-muted-foreground">
              请求不自动重试，不跟随重定向；Cookie 不跨请求保存。响应展示上限 2 MiB。
            </p>
          </TabsContent>
        </Tabs>
        <div className="mt-4 flex justify-end">
          <Button
            variant="ghost"
            size="sm"
            disabled={busy}
            onClick={() => {
              if (window.confirm('重置为文档中的请求示例？')) setDraft(key, initialDraft(endpoint))
            }}
          >
            <RotateCcw size={13} />
            重置请求
          </Button>
        </div>
      </div>
      <div className="min-w-0 space-y-5">
        {error && (
          <p
            role="alert"
            className="rounded-lg border border-destructive/30 bg-destructive/5 p-4 text-sm text-destructive"
          >
            {error}
          </p>
        )}
        {response ? (
          <div className="min-w-0 rounded-xl border bg-card shadow-panel">
            <div className="flex flex-wrap items-center gap-5 border-b px-5 py-4 text-xs">
              <span
                className={`flex items-center gap-2 font-semibold ${response.status < 300 ? 'text-success' : response.status < 400 ? 'text-info' : response.status < 500 ? 'text-warning' : 'text-danger'}`}
              >
                <CircleCheck size={15} />
                {response.status} {response.statusText}
              </span>
              <span className="flex items-center gap-1.5 text-muted-foreground">
                <Clock3 size={14} />
                {response.durationMs} ms
              </span>
              <span className="flex items-center gap-1.5 text-muted-foreground">
                <ArrowDown size={14} />
                {(response.size / 1024).toFixed(2)} KiB
              </span>
              {response.truncated && <Badge variant="destructive">响应已截断</Badge>}
            </div>
            <Tabs defaultValue="body" className="p-5">
              <TabsList>
                <TabsTrigger value="body">响应体</TabsTrigger>
                <TabsTrigger value="headers">响应头</TabsTrigger>
              </TabsList>
              <TabsContent value="body" className="pt-3">
                <Code
                  label={response.encoding === 'base64' ? '二进制响应 · Base64' : '响应正文'}
                  value={response.encoding === 'base64' ? response.body : formatJSON(response.body)}
                />
              </TabsContent>
              <TabsContent value="headers" className="pt-3">
                <Code
                  label="响应头"
                  value={Object.entries(response.headers)
                    .flatMap(([key, values]) => values.map((v) => `${key}: ${v}`))
                    .join('\n')}
                />
              </TabsContent>
            </Tabs>
          </div>
        ) : (
          <Empty
            title={busy ? '正在等待响应…' : '准备好后，发送第一个请求'}
            description={
              busy
                ? '请求由本地服务执行，可随时取消。'
                : '响应状态、耗时、请求头和正文会显示在这里。'
            }
          />
        )}
      </div>
    </div>
  )
}
