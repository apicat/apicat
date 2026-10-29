import { useEffect, useState } from 'react'
import {
  BrowserRouter,
  Link,
  NavLink,
  Route,
  Routes,
  useLocation,
  useSearchParams,
} from 'react-router-dom'
import {
  ArrowUpRight,
  BookOpen,
  Boxes,
  ChevronDown,
  ChevronRight,
  CircleDot,
  Command,
  Globe2,
  Layers3,
  Menu,
  Search,
  Settings2,
  ShieldCheck,
  Sparkles,
  Sun,
  Moon,
  Monitor,
  X,
  RefreshCw,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  BodyDoc,
  Code,
  Empty,
  EndpointDoc,
  Markdown,
  MediaContent,
  Method,
  Panel,
  Parameters,
  SchemaTree,
} from '@/components/docs'
import { Environments } from '@/components/environments'
import { RequestEditor } from '@/components/request-editor'
import { WorkbenchContext, useWorkbench } from '@/lib/workbench'
import { api, setToken } from '@/lib/api'
import { cn } from '@/lib/utils'
import type { APIDoc, Endpoint, RequestDraft, Session, Workspace } from '@/lib/types'

function endpointURL(e: Endpoint) {
  return `/endpoints?operation=${encodeURIComponent(e.Method + ' ' + e.Path)}`
}
function uniqueEndpoints(doc: APIDoc) {
  return [
    ...new Map(
      (doc.TagGroups || []).flatMap((g) => g.Endpoints || []).map((e) => [e.Method + e.Path, e]),
    ).values(),
  ]
}
function Theme() {
  const [theme, setTheme] = useState(() => {
    try {
      return localStorage.getItem('apicat-theme') || 'system'
    } catch {
      return 'system'
    }
  })
  useEffect(() => {
    const system = matchMedia('(prefers-color-scheme: dark)')
    const apply = () =>
      document.documentElement.classList.toggle(
        'dark',
        theme === 'dark' || (theme === 'system' && system.matches),
      )
    apply()
    system.addEventListener('change', apply)
    try {
      localStorage.setItem('apicat-theme', theme)
    } catch {
      /* unavailable storage */
    }
    return () => system.removeEventListener('change', apply)
  }, [theme])
  return (
    <div className="flex items-center gap-1 rounded-lg border bg-card p-1" aria-label="主题">
      {[
        { id: 'light', icon: Sun, label: '浅色主题' },
        { id: 'system', icon: Monitor, label: '跟随系统' },
        { id: 'dark', icon: Moon, label: '深色主题' },
      ].map(({ id, icon: Icon, label }) => (
        <button
          key={id}
          title={label}
          aria-label={label}
          aria-pressed={theme === id}
          onClick={() => setTheme(id)}
          className={cn(
            'rounded p-1.5 text-muted-foreground hover:text-foreground',
            theme === id && 'bg-accent text-primary',
          )}
        >
          <Icon size={14} />
        </button>
      ))}
    </div>
  )
}
function Sidebar({ close }: { close: () => void }) {
  const { doc } = useWorkbench()
  const location = useLocation()
  const [search, setSearch] = useState('')
  const [expanded, setExpanded] = useState<Record<string, boolean>>({})
  const q = search.toLowerCase().trim()
  const matches = (text: string) => text.toLowerCase().includes(q)
  const groups = (doc.TagGroups || [])
    .map((g) => ({
      ...g,
      Endpoints: g.Endpoints.filter((e) =>
        matches([g.Tag.Name, e.Method, e.Path, e.Summary, e.Description, e.OperationID].join(' ')),
      ),
    }))
    .filter((g) => g.Endpoints.length)
  const components = (doc.ComponentGroups || [])
    .map((g) => ({
      ...g,
      Definitions: g.Definitions.filter((d) => matches(`${g.Title} ${d.Name} ${d.Description}`)),
    }))
    .filter((g) => g.Definitions.length)
  function toggle(key: string) {
    setExpanded((v) => ({ ...v, [key]: !v[key] }))
  }
  return (
    <>
      <div className="flex h-[72px] items-center gap-3 border-b px-6">
        <div className="flex size-9 items-center justify-center rounded-xl bg-primary text-primary-foreground">
          <Command size={21} />
        </div>
        <Link to="/" onClick={close} className="text-xl font-bold tracking-tight">
          apicat
          <span className="ml-2 text-[10px] font-medium uppercase tracking-widest text-muted-foreground">
            local
          </span>
        </Link>
        <Button
          size="icon"
          variant="ghost"
          className="ml-auto lg:hidden"
          aria-label="关闭导航"
          onClick={close}
        >
          <X size={17} />
        </Button>
      </div>
      <div className="p-4">
        <div className="relative">
          <Search className="absolute left-3 top-2.5 text-muted-foreground" size={15} />
          <Input
            aria-label="搜索接口和组件"
            className="pl-9 text-xs"
            placeholder="搜索接口、路径或组件…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>
      <nav className="min-h-0 flex-1 overflow-y-auto px-3 pb-5">
        <div className="space-y-1 pb-5">
          {[
            { to: '/', icon: BookOpen, text: '项目概览' },
            { to: '/endpoints', icon: Layers3, text: '接口文档' },
            { to: '/components', icon: Boxes, text: '公共组件' },
            { to: '/environments', icon: Settings2, text: '环境与变量' },
          ].map(({ to, icon: Icon, text }) => (
            <NavLink
              key={to}
              to={to}
              end
              onClick={close}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-3 rounded-lg px-3 py-2.5 text-[13px] text-muted-foreground transition-colors hover:bg-muted',
                  isActive && 'bg-accent font-medium text-primary',
                )
              }
            >
              <Icon size={16} />
              {text}
            </NavLink>
          ))}
        </div>
        <div className="px-3 pb-3 text-[10px] font-semibold tracking-[.18em] text-muted-foreground">
          API COLLECTION
        </div>
        {groups.map((g) => {
          const key = `tag:${g.Tag.Name}`
          const open = q !== '' || expanded[key]
          return (
            <div key={key} className="mb-1">
              <button
                className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-xs hover:bg-muted"
                aria-expanded={!!open}
                onClick={() => toggle(key)}
              >
                {open ? <ChevronDown size={13} /> : <ChevronRight size={13} />}
                <span className="truncate">{g.Tag.Name}</span>
                <span className="ml-auto text-[10px] text-muted-foreground">
                  {g.Endpoints.length}
                </span>
              </button>
              {open && (
                <div className="space-y-0.5 py-1">
                  {g.Endpoints.map((e) => (
                    <Link
                      key={e.Method + e.Path}
                      to={endpointURL(e)}
                      onClick={close}
                      title={`${e.Method} ${e.Path} · ${e.Summary}`}
                      className={cn(
                        'ml-2 flex items-center gap-2 rounded-md px-2 py-2 hover:bg-muted',
                        location.pathname + location.search === endpointURL(e) && 'bg-accent',
                      )}
                    >
                      <Method method={e.Method} />
                      <span className="truncate font-mono text-[11px]">{e.Path}</span>
                    </Link>
                  ))}
                </div>
              )}
            </div>
          )
        })}
        <div className="px-3 pb-3 pt-6 text-[10px] font-semibold tracking-[.18em] text-muted-foreground">
          COMPONENTS
        </div>
        {components.map((g) => {
          const key = `component:${g.Kind}`
          const open = q !== '' || expanded[key]
          return (
            <div key={key}>
              <button
                className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-xs hover:bg-muted"
                aria-expanded={!!open}
                onClick={() => toggle(key)}
              >
                {open ? <ChevronDown size={13} /> : <ChevronRight size={13} />}
                <span>{g.Title}</span>
                <span className="ml-auto text-muted-foreground">{g.Definitions.length}</span>
              </button>
              {open &&
                g.Definitions.map((d) => (
                  <Link
                    key={d.Name}
                    onClick={close}
                    className="ml-5 block truncate rounded px-3 py-2 font-mono text-xs text-muted-foreground hover:bg-muted hover:text-foreground"
                    to={`/components?type=${g.Kind}&name=${encodeURIComponent(d.Name)}`}
                  >
                    {d.Name}
                  </Link>
                ))}
            </div>
          )
        })}
        {q && !groups.length && !components.length && (
          <p className="p-3 text-xs text-muted-foreground">未找到匹配的接口或组件</p>
        )}
      </nav>
      <div className="flex items-center justify-between border-t px-5 py-4">
        <span className="flex items-center gap-2 text-[11px] text-muted-foreground">
          <span className="size-1.5 rounded-full bg-primary" />
          本地工作台
        </span>
        <Theme />
      </div>
    </>
  )
}
function Overview() {
  const { doc, workspace } = useWorkbench()
  const endpoints = uniqueEndpoints(doc)
  return (
    <div className="mx-auto max-w-6xl space-y-7">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <p className="mb-3 text-xs font-medium tracking-[.18em] text-primary">
            YOUR API, AT A GLANCE
          </p>
          <h1 className="text-3xl font-semibold tracking-tight">{doc.Title}</h1>
          <div className="mt-3 flex items-center gap-3">
            <Badge variant="outline" className="font-mono">
              v{doc.Version}
            </Badge>
            <span className="text-xs text-muted-foreground">文档、环境与调试，一个工作台。</span>
          </div>
        </div>
        <Button asChild>
          <Link to="/endpoints">
            浏览接口 <ArrowUpRight size={15} />
          </Link>
        </Button>
      </div>
      <div className="grid gap-4 sm:grid-cols-3">
        {[
          {
            label: '接口数量',
            value: endpoints.length,
            icon: Layers3,
            sub: '来自当前 OpenAPI 文档',
          },
          {
            label: '运行环境',
            value: workspace.environments.length,
            icon: Globe2,
            sub: '在环境设置中管理地址与变量',
          },
          {
            label: '文档质量提示',
            value: doc.LintTotal,
            icon: ShieldCheck,
            sub: doc.LintTotal ? '查看需要完善的文档信息' : '已通过当前文档质量检查',
          },
        ].map(({ label, value, icon: Icon, sub }) => (
          <div key={label} className="rounded-xl border bg-card p-5">
            <div className="flex items-center justify-between text-xs text-muted-foreground">
              {label}
              <Icon size={17} />
            </div>
            <div className="mb-2 mt-4 text-3xl font-semibold tabular-nums">{value}</div>
            <p className="text-[11px] text-muted-foreground">{sub}</p>
          </div>
        ))}
      </div>
      {doc.Description && (
        <Panel title="关于此 API">
          <Markdown>{doc.Description}</Markdown>
        </Panel>
      )}
      <div className="grid items-start gap-5 xl:grid-cols-[minmax(0,1.4fr)_minmax(0,1fr)]">
        <Panel title="文档服务器" action={<Globe2 size={15} className="text-muted-foreground" />}>
          {doc.Servers?.length ? (
            <div className="space-y-4">
              {doc.Servers.map((s, i) => (
                <div key={i} className="rounded-lg border bg-muted/30 p-3">
                  <code className="break-all text-xs text-primary">{s.URL}</code>
                  {s.Description && (
                    <p className="mt-1 text-xs text-muted-foreground">{s.Description}</p>
                  )}
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">
              文档未声明服务器地址。可在环境设置中添加。
            </p>
          )}
        </Panel>
        <Panel title="认证方案">
          <div className="space-y-4">
            {doc.Security?.length ? (
              doc.Security.map((s) => (
                <div key={s.Name}>
                  <div className="mb-1 text-sm font-medium">{s.Name}</div>
                  <Badge variant="secondary">
                    {s.Type}
                    {s.Scheme ? ` · ${s.Scheme}` : ''}
                  </Badge>
                </div>
              ))
            ) : (
              <p className="text-sm text-muted-foreground">文档未声明认证方案</p>
            )}
          </div>
        </Panel>
      </div>
      <Panel title="文档质量" action={<Sparkles size={15} className="text-primary" />}>
        {doc.LintTotal ? (
          <ul className="divide-y">
            {doc.LintCounts?.map((c) => (
              <li key={c.Rule} className="flex items-center justify-between py-3 text-sm">
                <span>{c.Message}</span>
                <Badge variant="outline">{c.Count}</Badge>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-sm text-primary">没有质量提示。接口语义与实际行为仍需通过调试验证。</p>
        )}
      </Panel>
      {doc.Tags?.length ? (
        <Panel title="接口分组">
          <div className="divide-y">
            {doc.Tags.map((t) => (
              <div key={t.Name} className="grid gap-2 py-3 sm:grid-cols-[150px_1fr]">
                <span className="font-medium">{t.Name}</span>
                <div className="text-sm text-muted-foreground">
                  <Markdown>{t.Description}</Markdown>
                </div>
              </div>
            ))}
          </div>
        </Panel>
      ) : null}
    </div>
  )
}
function Endpoints() {
  const { doc } = useWorkbench()
  const [params] = useSearchParams()
  const location = useLocation()
  const endpoints = uniqueEndpoints(doc)
  const requested = params.get('operation')
  const endpoint = endpoints.find((e) =>
    requested ? e.Method + ' ' + e.Path === requested : e.Anchor === location.hash.slice(1),
  )
  if (!endpoint)
    return (
      <div className="mx-auto max-w-6xl">
        <p className="mb-2 text-xs tracking-[.18em] text-primary">API COLLECTION</p>
        <h1 className="mb-2 text-2xl font-semibold">接口文档</h1>
        <p className="mb-6 text-sm text-muted-foreground">
          选择接口查看定义、编辑参数并发起调试请求。
        </p>
        {requested ? (
          <Empty title="接口不存在" description="文档可能已更新，请从列表中重新选择接口。" />
        ) : null}
        <div className="mt-4 divide-y rounded-xl border bg-card">
          {endpoints.map((e) => (
            <Link
              key={e.Method + e.Path}
              to={endpointURL(e)}
              className="flex flex-wrap items-center gap-4 p-4 transition-colors hover:bg-muted"
            >
              <Method method={e.Method} />
              <code className="min-w-44 text-xs">{e.Path}</code>
              <span className="flex-1 text-sm text-muted-foreground">
                {e.Summary || e.OperationID}
              </span>
              {!!e.LintHints?.length && <Badge variant="outline">{e.LintHints.length} 提示</Badge>}
              <ChevronRight size={15} />
            </Link>
          ))}
        </div>
        {!endpoints.length && <Empty title="文档中没有接口" />}
      </div>
    )
  return (
    <div className="mx-auto max-w-6xl" key={endpoint.Method + endpoint.Path}>
      <div className="mb-6">
        <div className="mb-3 flex flex-wrap items-center gap-3">
          <Method method={endpoint.Method} />
          <code className="break-all text-sm text-muted-foreground">{endpoint.Path}</code>
          {endpoint.Deprecated && <Badge variant="destructive">已弃用</Badge>}
        </div>
        <h1 className="text-2xl font-semibold tracking-tight">
          {endpoint.Summary || endpoint.OperationID || endpoint.Path}
        </h1>
        {endpoint.OperationID && (
          <p className="mt-2 font-mono text-xs text-muted-foreground">{endpoint.OperationID}</p>
        )}
      </div>
      <Tabs defaultValue="documentation">
        <TabsList className="mb-5">
          <TabsTrigger value="documentation">
            <BookOpen size={14} />
            文档
          </TabsTrigger>
          <TabsTrigger value="debug">
            <CircleDot size={14} />
            接口调试
          </TabsTrigger>
        </TabsList>
        <TabsContent value="documentation">
          <EndpointDoc endpoint={endpoint} />
        </TabsContent>
        <TabsContent value="debug">
          <RequestEditor endpoint={endpoint} />
        </TabsContent>
      </Tabs>
    </div>
  )
}
function Components() {
  const { doc } = useWorkbench()
  const [params] = useSearchParams()
  const kind = params.get('type'),
    name = params.get('name')
  const group = doc.ComponentGroups?.find((g) => g.Kind === kind)
  const component = group?.Definitions.find((d) => d.Name === name)
  if (kind && name && !component)
    return <Empty title="公共组件不存在" description="此定义可能已被移除，请从侧边栏重新选择。" />
  if (component)
    return (
      <div className="mx-auto max-w-6xl space-y-5">
        <p className="text-xs uppercase tracking-widest text-primary">
          COMPONENTS / {group?.Title}
        </p>
        <h1 className="break-all text-2xl font-semibold">{component.Name}</h1>
        <Markdown>{component.Description}</Markdown>
        {component.Schema && (
          <Panel title="Schema">
            <SchemaTree schema={component.Schema} />
          </Panel>
        )}
        {component.Parameters && (
          <Panel title="参数">
            <Parameters parameters={component.Parameters} />
          </Panel>
        )}
        {component.RequestBody && (
          <Panel title="请求体">
            <BodyDoc body={component.RequestBody} />
          </Panel>
        )}
        {component.Response && (
          <Panel title="响应">
            <Markdown>{component.Response.Description}</Markdown>
            <MediaContent content={component.Response.Content} />
          </Panel>
        )}
        <Code label="完整定义 · YAML" value={component.Definition} />
      </div>
    )
  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <div>
        <p className="mb-2 text-xs tracking-[.18em] text-primary">REUSABLE DEFINITIONS</p>
        <h1 className="text-2xl font-semibold">公共组件</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          浏览共享的 Schema、参数、响应以及其他可复用定义。
        </p>
      </div>
      {doc.ComponentGroups?.length ? (
        doc.ComponentGroups.map((g) => (
          <Panel
            key={g.Kind}
            title={g.Title}
            action={<Badge variant="secondary">{g.Definitions.length}</Badge>}
          >
            <div className="divide-y">
              {g.Definitions.map((d) => (
                <Link
                  key={d.Name}
                  to={`/components?type=${g.Kind}&name=${encodeURIComponent(d.Name)}`}
                  className="flex items-start justify-between gap-4 rounded px-2 py-3 hover:bg-muted"
                >
                  <div>
                    <code className="text-sm font-medium text-primary">{d.Name}</code>
                    {d.Description && (
                      <p className="mt-1 line-clamp-2 text-xs text-muted-foreground">
                        {d.Description}
                      </p>
                    )}
                  </div>
                  <ChevronRight size={16} className="shrink-0 text-muted-foreground" />
                </Link>
              ))}
            </div>
          </Panel>
        ))
      ) : (
        <Empty title="文档中没有公共组件" />
      )}
    </div>
  )
}
function Shell({ pollError }: { pollError: string }) {
  const { doc, workspace, saveWorkspace } = useWorkbench()
  const location = useLocation()
  const [mobile, setMobile] = useState(false)
  const [error, setError] = useState('')
  const [switching, setSwitching] = useState(false)
  const title =
    (
      {
        '/': '项目概览',
        '/endpoints': '接口文档',
        '/components': '公共组件',
        '/environments': '环境与变量',
      } as Record<string, string>
    )[location.pathname] || '工作台'
  useEffect(() => {
    document.title = `${title} · ${doc.Title} — apicat`
  }, [title, doc.Title])
  async function choose(value: string) {
    setSwitching(true)
    setError('')
    try {
      await saveWorkspace({ ...workspace, activeEnvironment: value === 'none' ? '' : value })
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setSwitching(false)
    }
  }
  return (
    <div className="min-h-screen">
      <aside className="fixed inset-y-0 left-0 z-20 hidden w-[272px] flex-col border-r bg-card lg:flex">
        <Sidebar close={() => setMobile(false)} />
      </aside>
      {mobile && (
        <div className="fixed inset-0 z-50 bg-black/40 lg:hidden" onClick={() => setMobile(false)}>
          <aside
            className="flex h-full w-[min(88vw,300px)] flex-col bg-card"
            onClick={(e) => e.stopPropagation()}
          >
            <Sidebar close={() => setMobile(false)} />
          </aside>
        </div>
      )}
      <div className="lg:pl-[272px]">
        <header className="sticky top-0 z-10 flex h-[72px] items-center justify-between gap-4 border-b bg-background/95 px-5 backdrop-blur sm:px-8">
          <div className="flex min-w-0 items-center gap-3">
            <Button
              variant="ghost"
              size="icon"
              className="lg:hidden"
              aria-label="打开导航"
              onClick={() => setMobile(true)}
            >
              <Menu size={18} />
            </Button>
            <span className="hidden truncate text-xs text-muted-foreground sm:block">
              {doc.Title}
            </span>
            <ChevronRight size={12} className="hidden text-muted-foreground sm:block" />
            <span className="shrink-0 text-xs font-medium">{title}</span>
          </div>
          {location.pathname !== '/environments' && (
            <Select
              disabled={switching}
              value={workspace.activeEnvironment || 'none'}
              onValueChange={choose}
            >
              <SelectTrigger aria-label="当前运行环境" className="w-44 bg-card text-xs">
                <Globe2 size={14} />
                <SelectValue placeholder="选择环境" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">不使用环境</SelectItem>
                {workspace.environments.map((e) => (
                  <SelectItem key={e.id} value={e.id}>
                    {e.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        </header>
        {(error || pollError) && (
          <div
            role="alert"
            className="border-b bg-amber-100/30 px-8 py-3 text-xs text-amber-700 dark:text-amber-300"
          >
            {error || pollError}
            <button className="ml-3 underline" onClick={() => window.location.reload()}>
              重新加载
            </button>
          </div>
        )}
        <main className="p-5 pb-12 sm:p-8 xl:p-10">
          <Routes>
            <Route path="/" element={<Overview />} />
            <Route path="/endpoints" element={<Endpoints />} />
            <Route path="/components" element={<Components />} />
            <Route path="/environments" element={<Environments />} />
            <Route path="*" element={<Empty title="页面不存在" />} />
          </Routes>
        </main>
      </div>
    </div>
  )
}
export default function App() {
  const [data, setData] = useState<{ doc: APIDoc; workspace: Workspace; session: Session } | null>(
    null,
  )
  const [error, setError] = useState('')
  const [pollError, setPollError] = useState('')
  const [drafts, setDrafts] = useState<Record<string, RequestDraft>>({})
  useEffect(() => {
    const controller = new AbortController()
    let disposed = false
    Promise.all([
      api<APIDoc>('/document', { signal: controller.signal }),
      api<Workspace>('/workspace', { signal: controller.signal }),
      api<Session>('/session', { signal: controller.signal }),
    ])
      .then(([doc, workspace, session]) => {
        if (disposed) return
        setToken(session.token)
        setData({
          doc,
          workspace: {
            ...workspace,
            environments: workspace.environments || [],
            variables: workspace.variables || [],
          },
          session,
        })
      })
      .catch((e) => {
        if (e.name !== 'AbortError') setError(e.message)
      })
    return () => {
      disposed = true
      controller.abort()
    }
  }, [])
  useEffect(() => {
    if (!data) return
    let pending = false
    const controller = new AbortController()
    const timer = setInterval(async () => {
      if (pending || document.hidden) return
      pending = true
      try {
        const doc = await api<APIDoc>('/document', { signal: controller.signal })
        setData((current) =>
          current && JSON.stringify(current.doc) !== JSON.stringify(doc)
            ? { ...current, doc }
            : current,
        )
        setPollError('')
      } catch (e) {
        if ((e as Error).name !== 'AbortError')
          setPollError('与本地服务的连接已中断，当前显示上次成功加载的文档。')
      } finally {
        pending = false
      }
    }, 3000)
    return () => {
      clearInterval(timer)
      controller.abort()
    }
  }, [!!data])
  async function saveWorkspace(state: Workspace) {
    const workspace = await api<Workspace>('/workspace', {
      method: 'PUT',
      body: JSON.stringify(state),
    })
    setData((d) => (d ? { ...d, workspace } : d))
  }
  if (error)
    return (
      <div className="mx-auto max-w-lg p-10">
        <Empty title="无法加载工作台" description={error}>
          <Button className="mt-5" onClick={() => window.location.reload()}>
            <RefreshCw size={15} />
            重试
          </Button>
        </Empty>
      </div>
    )
  if (!data)
    return (
      <div className="flex min-h-screen items-center justify-center gap-3 text-muted-foreground">
        <Command className="animate-pulse text-primary" />
        正在加载工作台…
      </div>
    )
  return (
    <WorkbenchContext.Provider
      value={{
        ...data,
        saveWorkspace,
        drafts,
        setDraft: (key, draft) => setDrafts((d) => ({ ...d, [key]: draft })),
      }}
    >
      <BrowserRouter>
        <Shell pollError={pollError} />
      </BrowserRouter>
    </WorkbenchContext.Provider>
  )
}
