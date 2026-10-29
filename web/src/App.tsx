import { lintSummary } from '@/lib/lint-messages'
import { LanguageSwitcher } from '@/components/language-switcher'
import { useI18n, asError } from '@/lib/i18n'
import { useEffect, useState } from 'react'
import {
  BrowserRouter,
  Link,
  NavLink,
  Route,
  Routes,
  useLocation,
  useSearchParams,
  useParams,
} from 'react-router-dom'
import {
  ArrowUpRight,
  BookOpen,
  Boxes,
  ChevronDown,
  ChevronRight,
  CircleDot,
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
import { Brand } from '@/components/brand'
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
  return `/endpoints/${encodeURIComponent(e.RouteID)}`
}
function uniqueEndpoints(doc: APIDoc) {
  return [
    ...new Map(
      (doc.TagGroups || []).flatMap((g) => g.Endpoints || []).map((e) => [e.Method + e.Path, e]),
    ).values(),
  ]
}
function useThemePreference() {
  const [theme, setTheme] = useState(() => {
    try {
      return localStorage.getItem('apicat-theme') || 'system'
    } catch {
      return 'system'
    }
  })
  useEffect(() => {
    const system = matchMedia('(prefers-color-scheme: dark)')
    const apply = () => {
      document.documentElement.classList.toggle(
        'dark',
        theme === 'dark' || (theme === 'system' && system.matches),
      )
      const active = document.documentElement.classList.contains('dark') ? 'dark' : 'light'
      document.querySelectorAll<HTMLLinkElement>('[data-icon-theme]').forEach((icon) => {
        icon.media = icon.dataset.iconTheme === active ? 'all' : 'not all'
      })
    }
    apply()
    system.addEventListener('change', apply)
    try {
      localStorage.setItem('apicat-theme', theme)
    } catch {
      /* unavailable storage */
    }
    return () => system.removeEventListener('change', apply)
  }, [theme])
  return { theme, setTheme }
}
function Theme({ theme, setTheme }: ReturnType<typeof useThemePreference>) {
  const { t } = useI18n()
  return (
    <div className="flex items-center gap-1 rounded-lg border bg-card p-1" aria-label={t('主题')}>
      {[
        { id: 'light', icon: Sun, label: t('浅色主题') },
        { id: 'system', icon: Monitor, label: t('跟随系统') },
        { id: 'dark', icon: Moon, label: t('深色主题') },
      ].map(({ id, icon: Icon, label }) => (
        <button
          key={id}
          title={label}
          aria-label={label}
          aria-pressed={theme === id}
          onClick={() => setTheme(id)}
          className={cn(
            'rounded p-1.5 text-muted-foreground hover:text-foreground',
            theme === id && 'bg-selected text-brand-text',
          )}
        >
          <Icon size={14} />
        </button>
      ))}
    </div>
  )
}
function Sidebar({
  close,
  theme,
}: {
  close: () => void
  theme: ReturnType<typeof useThemePreference>
}) {
  const { t, number } = useI18n()
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
      <div className="flex h-16 shrink-0 items-center gap-3 border-b px-5">
        <Link to="/" onClick={close} aria-label={t('apicat 首页')} className="rounded-sm">
          <Brand />
        </Link>
        <span className="rounded border px-1.5 py-0.5 text-[10px] font-medium tracking-wider text-muted-foreground">
          LOCAL
        </span>
        <Button
          size="icon"
          variant="ghost"
          className="ml-auto lg:hidden"
          aria-label={t('关闭导航')}
          onClick={close}
        >
          <X size={17} />
        </Button>
      </div>
      <div className="p-4">
        <div className="relative">
          <Search
            className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
            size={15}
          />
          <Input
            aria-label={t('搜索接口和组件')}
            className="pl-9 text-xs"
            placeholder={t('搜索接口、路径或组件…')}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>
      <nav className="min-h-0 flex-1 overflow-y-auto px-3 pb-5">
        <div className="space-y-1 pb-5">
          {[
            { to: '/', icon: BookOpen, text: t('项目概览') },
            { to: '/endpoints', icon: Layers3, text: t('接口文档') },
            { to: '/components', icon: Boxes, text: t('公共组件') },
            { to: '/environments', icon: Settings2, text: t('环境与变量') },
          ].map(({ to, icon: Icon, text }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              onClick={close}
              className={({ isActive }) =>
                cn(
                  'nav-item flex items-center gap-3 rounded-md px-3 py-2.5 text-[13px] text-secondary-text transition-colors hover:bg-hover',
                  isActive && 'bg-selected font-medium text-brand-text',
                )
              }
            >
              <Icon size={16} />
              {text}
            </NavLink>
          ))}
        </div>
        <div className="px-3 pb-3 text-[11px] font-semibold tracking-[.12em] text-muted-foreground">
          {t('接口集合')}
        </div>
        {groups.map((g) => {
          const key = `tag:${g.Tag.Name}`
          const open = q !== '' || expanded[key]
          return (
            <div key={key} className="mb-1">
              <button
                className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-[13px] text-secondary-text hover:bg-hover"
                aria-expanded={!!open}
                onClick={() => toggle(key)}
              >
                {open ? <ChevronDown size={13} /> : <ChevronRight size={13} />}
                <span className="truncate">{g.Tag.Name}</span>
                <span className="ml-auto text-[10px] text-muted-foreground">
                  {number(g.Endpoints.length)}
                </span>
              </button>
              {open && (
                <div className="space-y-0.5 py-1">
                  {g.Endpoints.map((e) => (
                    <Link
                      key={e.Method + e.Path}
                      to={endpointURL(e)}
                      aria-current={location.pathname === endpointURL(e) ? 'page' : undefined}
                      onClick={close}
                      title={`${e.Method} ${e.Path}${e.Summary?.trim() ? ` · ${e.Summary.trim()}` : ''}`}
                      className={cn(
                        'nav-item ml-2 flex items-center gap-2 rounded-md px-2 py-2 text-secondary-text hover:bg-hover',
                        location.pathname === endpointURL(e) && 'bg-selected',
                      )}
                    >
                      <Method method={e.Method} />
                      <span className="min-w-0 truncate text-xs">
                        {e.Summary?.trim() || e.Path}
                      </span>
                    </Link>
                  ))}
                </div>
              )}
            </div>
          )
        })}
        <div className="px-3 pb-3 pt-6 text-[11px] font-semibold tracking-[.12em] text-muted-foreground">
          {t('组件')}
        </div>
        {components.map((g) => {
          const key = `component:${g.Kind}`
          const open = q !== '' || expanded[key]
          return (
            <div key={key}>
              <button
                className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-[13px] text-secondary-text hover:bg-hover"
                aria-expanded={!!open}
                onClick={() => toggle(key)}
              >
                {open ? <ChevronDown size={13} /> : <ChevronRight size={13} />}
                <span>{g.Title}</span>
                <span className="ml-auto text-muted-foreground">
                  {number(g.Definitions.length)}
                </span>
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
          <p className="p-3 text-xs text-muted-foreground">{t('未找到匹配的接口或组件')}</p>
        )}
      </nav>
      <div className="space-y-3 border-t px-4 py-3">
        <LanguageSwitcher />
        <div className="flex items-center justify-between">
          <span className="flex items-center gap-2 text-[11px] text-muted-foreground">
            <span className="size-1.5 rounded-full bg-success" />
            {t('本地工作台')}
          </span>
          <Theme {...theme} />
        </div>
      </div>
    </>
  )
}
function Overview() {
  const { t, number, locale } = useI18n()
  const { doc, workspace } = useWorkbench()
  const endpoints = uniqueEndpoints(doc)
  return (
    <div className="mx-auto max-w-6xl space-y-7">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <p className="mb-3 text-xs font-medium tracking-[.12em] text-muted-foreground">
            {t('API 概览')}
          </p>
          <h1 className="text-3xl font-semibold tracking-tight">{doc.Title}</h1>
          <div className="mt-3 flex items-center gap-3">
            <Badge variant="outline" className="font-mono">
              v{doc.Version}
            </Badge>
            <span className="text-xs text-muted-foreground">
              {t('文档、环境与调试，一个工作台。')}
            </span>
          </div>
        </div>
        <Button asChild>
          <Link to="/endpoints">
            {t('浏览接口')}
            <ArrowUpRight size={15} />
          </Link>
        </Button>
      </div>
      <div className="grid gap-4 sm:grid-cols-3">
        {[
          {
            label: t('接口数量'),
            value: endpoints.length,
            icon: Layers3,
            sub: t('来自当前 OpenAPI 文档'),
          },
          {
            label: t('运行环境'),
            value: workspace.environments.length,
            icon: Globe2,
            sub: t('在环境设置中管理地址与变量'),
          },
          {
            label: t('文档质量提示'),
            value: doc.LintTotal,
            icon: ShieldCheck,
            sub: doc.LintTotal ? t('查看需要完善的文档信息') : t('已通过当前文档质量检查'),
          },
        ].map(({ label, value, icon: Icon, sub }) => (
          <div key={label} className="rounded-xl border bg-card p-5 shadow-panel">
            <div className="flex items-center justify-between text-xs text-muted-foreground">
              {label}
              <Icon size={17} />
            </div>
            <div className="mb-2 mt-4 text-3xl font-semibold tabular-nums">{number(value)}</div>
            <p className="text-xs text-muted-foreground">{sub}</p>
          </div>
        ))}
      </div>
      {doc.Description && (
        <Panel title={t('关于此 API')}>
          <Markdown>{doc.Description}</Markdown>
        </Panel>
      )}
      <div className="grid items-start gap-5 xl:grid-cols-[minmax(0,1.4fr)_minmax(0,1fr)]">
        <Panel
          title={t('文档服务器')}
          action={<Globe2 size={15} className="text-muted-foreground" />}
        >
          {doc.Servers?.length ? (
            <div className="space-y-4">
              {doc.Servers.map((s, i) => (
                <div key={i} className="rounded-lg border bg-muted/30 p-3">
                  <code className="break-all text-xs text-secondary-text">{s.URL}</code>
                  {s.Description && (
                    <p className="mt-1 text-xs text-muted-foreground">{s.Description}</p>
                  )}
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">
              {t('文档未声明服务器地址。可在环境设置中添加。')}
            </p>
          )}
        </Panel>
        <Panel title={t('认证方案')}>
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
              <p className="text-sm text-muted-foreground">{t('文档未声明认证方案')}</p>
            )}
          </div>
        </Panel>
      </div>
      <Panel
        title={t('文档质量')}
        action={<Sparkles size={15} className="text-muted-foreground" />}
      >
        {doc.LintTotal ? (
          <ul className="divide-y">
            {doc.LintCounts?.map((c) => (
              <li key={c.Rule} className="flex items-center justify-between py-3 text-sm">
                <span>{lintSummary(locale, c.Rule, c.Message)}</span>
                <Badge variant="outline">{number(c.Count)}</Badge>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-sm text-success">
            {t('没有质量提示。接口语义与实际行为仍需通过调试验证。')}
          </p>
        )}
      </Panel>
      {doc.Tags?.length ? (
        <Panel title={t('接口分组')}>
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
  const { t } = useI18n()
  const { doc } = useWorkbench()
  const { operationId: requested } = useParams()
  const endpoints = uniqueEndpoints(doc)
  const endpoint = endpoints.find((e) => e.RouteID === requested)
  if (!endpoint)
    return (
      <div className="mx-auto max-w-6xl">
        <p className="mb-2 text-xs tracking-[.12em] text-muted-foreground">{t('接口集合')}</p>
        <h1 className="mb-2 text-2xl font-semibold">{t('接口文档')}</h1>
        <p className="mb-6 text-sm text-muted-foreground">
          {t('选择接口查看定义、编辑参数并发起调试请求。')}
        </p>
        {requested ? (
          <Empty
            title={t('接口不存在')}
            description={t('文档可能已更新，请从列表中重新选择接口。')}
          />
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
              {!!e.LintHints?.length && (
                <Badge variant="outline">
                  {t('{count} 条提示', { count: e.LintHints.length })}
                </Badge>
              )}
              <ChevronRight size={15} />
            </Link>
          ))}
        </div>
        {!endpoints.length && <Empty title={t('文档中没有接口')} />}
      </div>
    )
  return (
    <div className="mx-auto max-w-6xl" key={endpoint.Method + endpoint.Path}>
      <div className="mb-6">
        <div className="mb-3 flex flex-wrap items-center gap-3">
          <Method method={endpoint.Method} />
          <code className="break-all text-sm text-muted-foreground">{endpoint.Path}</code>
          {endpoint.Deprecated && <Badge variant="destructive">{t('已弃用')}</Badge>}
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
            {t('文档')}
          </TabsTrigger>
          <TabsTrigger value="debug">
            <CircleDot size={14} />
            {t('接口调试')}
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
  const { t, number } = useI18n()
  const { doc } = useWorkbench()
  const [params] = useSearchParams()
  const kind = params.get('type'),
    name = params.get('name')
  const group = doc.ComponentGroups?.find((g) => g.Kind === kind)
  const component = group?.Definitions.find((d) => d.Name === name)
  if (kind && name && !component)
    return (
      <Empty
        title={t('公共组件不存在')}
        description={t('此定义可能已被移除，请从侧边栏重新选择。')}
      />
    )
  if (component)
    return (
      <div className="mx-auto max-w-6xl space-y-5">
        <p className="text-xs uppercase tracking-widest text-muted-foreground">
          {t('组件')} / {group?.Title}
        </p>
        <h1 className="break-all text-2xl font-semibold">{component.Name}</h1>
        <Markdown>{component.Description}</Markdown>
        {component.Schema && (
          <Panel title="Schema">
            <SchemaTree schema={component.Schema} />
          </Panel>
        )}
        {component.Parameters && (
          <Panel title={t('参数')}>
            <Parameters parameters={component.Parameters} />
          </Panel>
        )}
        {component.RequestBody && (
          <Panel title={t('请求体')}>
            <BodyDoc body={component.RequestBody} />
          </Panel>
        )}
        {component.Response && (
          <Panel title={t('响应')}>
            <Markdown>{component.Response.Description}</Markdown>
            <MediaContent content={component.Response.Content} />
          </Panel>
        )}
        <Code label={t('完整定义 · YAML')} value={component.Definition} />
      </div>
    )
  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <div>
        <p className="mb-2 text-xs tracking-[.12em] text-muted-foreground">{t('可复用定义')}</p>
        <h1 className="text-2xl font-semibold">{t('公共组件')}</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          {t('浏览共享的 Schema、参数、响应以及其他可复用定义。')}
        </p>
      </div>
      {doc.ComponentGroups?.length ? (
        doc.ComponentGroups.map((g) => (
          <Panel
            key={g.Kind}
            title={g.Title}
            action={<Badge variant="secondary">{number(g.Definitions.length)}</Badge>}
          >
            <div className="divide-y">
              {g.Definitions.map((d) => (
                <Link
                  key={d.Name}
                  to={`/components?type=${g.Kind}&name=${encodeURIComponent(d.Name)}`}
                  className="flex items-start justify-between gap-4 rounded px-2 py-3 hover:bg-muted"
                >
                  <div>
                    <code className="text-sm font-medium text-foreground">{d.Name}</code>
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
        <Empty title={t('文档中没有公共组件')} />
      )}
    </div>
  )
}
function Shell({ pollError }: { pollError: string }) {
  const { t, errorText } = useI18n()
  const theme = useThemePreference()
  const { doc, workspace, saveWorkspace } = useWorkbench()
  const location = useLocation()
  const [mobile, setMobile] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  const [switching, setSwitching] = useState(false)
  const title =
    (
      {
        '/': t('项目概览'),
        '/endpoints': t('接口文档'),
        '/components': t('公共组件'),
        '/environments': t('环境与变量'),
      } as Record<string, string>
    )[location.pathname] ||
    (location.pathname.startsWith('/endpoints/') ? t('接口文档') : t('工作台'))
  useEffect(() => {
    document.title = `${title} · ${doc.Title} — apicat`
  }, [title, doc.Title])
  async function choose(value: string) {
    setSwitching(true)
    setError(null)
    try {
      await saveWorkspace({ ...workspace, activeEnvironment: value === 'none' ? '' : value })
    } catch (e) {
      setError(asError(e))
    } finally {
      setSwitching(false)
    }
  }
  return (
    <div className="min-h-screen">
      <aside className="fixed inset-y-0 left-0 z-20 hidden w-[240px] flex-col border-r bg-app lg:flex">
        <Sidebar close={() => setMobile(false)} theme={theme} />
      </aside>
      {mobile && (
        <div className="fixed inset-0 z-50 bg-black/40 lg:hidden" onClick={() => setMobile(false)}>
          <aside
            className="flex h-full w-[min(88vw,300px)] flex-col bg-app"
            onClick={(e) => e.stopPropagation()}
          >
            <Sidebar close={() => setMobile(false)} theme={theme} />
          </aside>
        </div>
      )}
      <div className="lg:pl-[240px]">
        <header className="sticky top-0 z-10 flex h-16 items-center justify-between gap-4 border-b bg-app px-5 sm:px-8">
          <div className="flex min-w-0 items-center gap-3">
            <Button
              variant="ghost"
              size="icon"
              className="lg:hidden"
              aria-label={t('打开导航')}
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
              <SelectTrigger aria-label={t('当前运行环境')} className="w-44 bg-card text-xs">
                <Globe2 size={14} />
                <SelectValue placeholder={t('选择环境')} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">{t('不使用环境')}</SelectItem>
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
          <div role="alert" className="border-b bg-warning/10 px-8 py-3 text-xs text-warning">
            {error ? errorText(error) : t('与本地服务的连接已中断，当前显示上次成功加载的文档。')}
            <button className="ml-3 underline" onClick={() => window.location.reload()}>
              {t('重新加载')}
            </button>
          </div>
        )}
        <main className="p-5 pb-12 sm:p-8 xl:p-10">
          <Routes>
            <Route path="/" element={<Overview />} />
            <Route path="/endpoints" element={<Endpoints />} />
            <Route path="/endpoints/:operationId" element={<Endpoints />} />
            <Route path="/components" element={<Components />} />
            <Route path="/environments" element={<Environments />} />
            <Route path="*" element={<Empty title={t('页面不存在')} />} />
          </Routes>
        </main>
      </div>
    </div>
  )
}
export default function App() {
  const { t, errorText } = useI18n()
  const [data, setData] = useState<{ doc: APIDoc; workspace: Workspace; session: Session } | null>(
    null,
  )
  const [error, setError] = useState<Error | null>(null)
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
        if (e.name !== 'AbortError') setError(asError(e))
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
        if ((e as Error).name !== 'AbortError') setPollError('disconnected')
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
        <LanguageSwitcher />
        <Empty title={t('无法加载工作台')} description={errorText(error)}>
          <Button className="mt-5" onClick={() => window.location.reload()}>
            <RefreshCw size={15} />
            {t('重试')}
          </Button>
        </Empty>
      </div>
    )
  if (!data)
    return (
      <div className="flex min-h-screen items-center justify-center gap-3 text-muted-foreground">
        <Brand compact />
        {t('正在加载工作台…')}
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
