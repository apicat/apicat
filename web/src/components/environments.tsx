import { useEffect, useState } from 'react'
import { Globe2, Plus, Save, Trash2, FolderOpen, CheckCircle2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogHeader,
} from '@/components/ui/dialog'
import { PairEditor } from './pair-editor'
import { Panel, Empty } from './docs'
import { useWorkbench } from '@/lib/workbench'
import type { Environment } from '@/lib/types'

export function Environments() {
  const { workspace, saveWorkspace, session } = useWorkbench()
  const [draft, setDraft] = useState(() => structuredClone(workspace))
  const [selected, setSelected] = useState(
    workspace.activeEnvironment || workspace.environments?.[0]?.id || '',
  )
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)
  const [remove, setRemove] = useState(false)
  const dirty = JSON.stringify(draft) !== JSON.stringify(workspace)
  useEffect(() => {
    const handler = (e: BeforeUnloadEvent) => {
      if (dirty) e.preventDefault()
    }
    window.addEventListener('beforeunload', handler)
    return () => window.removeEventListener('beforeunload', handler)
  }, [dirty])
  // Native navigation and SPA links both ask before abandoning unsaved settings.
  useEffect(() => {
    const guard = (e: MouseEvent) => {
      const link = (e.target as HTMLElement).closest('a')
      if (
        dirty &&
        link &&
        new URL(link.href).pathname !== '/environments' &&
        !window.confirm('环境配置尚未保存，确定离开？')
      ) {
        e.preventDefault()
        e.stopPropagation()
      }
    }
    document.addEventListener('click', guard, true)
    return () => document.removeEventListener('click', guard, true)
  }, [dirty])
  function add(name = '新环境', baseUrl = '') {
    const env: Environment = { id: crypto.randomUUID(), name, baseUrl, variables: [] }
    setDraft((d) => ({
      ...d,
      environments: [...d.environments, env],
      activeEnvironment: d.activeEnvironment || env.id,
    }))
    setSelected(env.id)
    setSaved(false)
  }
  const env = draft.environments.find((e) => e.id === selected)
  function update(update: Partial<Environment>) {
    setDraft((d) => ({
      ...d,
      environments: d.environments.map((e) => (e.id === selected ? { ...e, ...update } : e)),
    }))
    setSaved(false)
  }
  async function save() {
    setSaving(true)
    setError('')
    setSaved(false)
    try {
      await saveWorkspace(draft)
      setDraft((d) => ({ ...d, revision: d.revision + 1 }))
      setSaved(true)
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setSaving(false)
    }
  }
  return (
    <div className="mx-auto max-w-6xl space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <p className="mb-2 text-xs font-medium uppercase tracking-[.12em] text-muted-foreground">
            WORKSPACE / ENVIRONMENTS
          </p>
          <h1 className="text-2xl font-semibold tracking-tight">环境与变量</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            切换部署地址，在请求中使用 {'{{变量名}}'} 引用变量。
          </p>
        </div>
        <Button onClick={save} disabled={saving || !dirty}>
          <Save size={15} />
          {saving ? '保存中…' : '保存配置'}
        </Button>
      </div>
      {error && (
        <div
          role="alert"
          className="rounded-lg border border-destructive/30 bg-destructive/5 p-3 text-destructive"
        >
          {error}
          <Button variant="link" onClick={() => window.location.reload()}>
            重新加载
          </Button>
        </div>
      )}
      {saved && (
        <p role="status" className="flex items-center gap-2 text-success">
          <CheckCircle2 size={16} />
          配置已保存
        </p>
      )}
      <div className="grid gap-5 lg:grid-cols-[220px_minmax(0,1fr)]">
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <h2 className="text-xs font-semibold text-muted-foreground">部署环境</h2>
            <Button size="icon" variant="ghost" aria-label="添加环境" onClick={() => add()}>
              <Plus size={16} />
            </Button>
          </div>
          {draft.environments.map((e) => (
            <button
              key={e.id}
              onClick={() => setSelected(e.id)}
              className={`flex w-full items-center gap-2 rounded-lg border p-3 text-left ${selected === e.id ? 'border-primary/40 bg-selected text-brand-text' : 'border-transparent hover:bg-muted'}`}
            >
              <Globe2 size={15} className="shrink-0" />
              <span className="flex-1 truncate text-sm">{e.name}</span>
              {draft.activeEnvironment === e.id && (
                <span className="size-1.5 rounded-full bg-success" />
              )}
            </button>
          ))}
          {draft.environments.length === 0 && (
            <div className="space-y-2">
              <p className="text-xs text-muted-foreground">快速创建常用环境</p>
              {['开发环境', '测试环境', '生产环境'].map((name) => (
                <Button
                  key={name}
                  variant="outline"
                  className="w-full justify-start"
                  onClick={() => add(name)}
                >
                  <Plus size={14} />
                  {name}
                </Button>
              ))}
            </div>
          )}
        </div>
        <div className="min-w-0">
          {env ? (
            <Panel
              title={env.name}
              action={
                <Button
                  size="icon"
                  variant="ghost"
                  aria-label="删除环境"
                  onClick={() => setRemove(true)}
                >
                  <Trash2 size={15} />
                </Button>
              }
            >
              <div className="space-y-6">
                <div className="grid gap-4 sm:grid-cols-2">
                  <label className="space-y-2 text-xs font-medium">
                    环境名称
                    <Input
                      aria-label="环境名称"
                      className="mt-2"
                      value={env.name}
                      onChange={(e) => update({ name: e.target.value })}
                    />
                  </label>
                  <label className="space-y-2 text-xs font-medium">
                    基础地址
                    <Input
                      aria-label="环境基础地址"
                      className="mt-2 font-mono"
                      placeholder="https://api.example.com/v1"
                      value={env.baseUrl}
                      onChange={(e) => update({ baseUrl: e.target.value })}
                    />
                  </label>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={draft.activeEnvironment === env.id}
                  onClick={() => setDraft((d) => ({ ...d, activeEnvironment: env.id }))}
                >
                  {draft.activeEnvironment === env.id ? '当前默认环境' : '设为默认环境'}
                </Button>
                <div>
                  <h3 className="mb-3 font-medium">
                    环境变量{' '}
                    <Badge variant="secondary" className="ml-2">
                      {env.variables.length}
                    </Badge>
                  </h3>
                  <PairEditor
                    rows={env.variables}
                    onChange={(variables) => update({ variables })}
                    secrets
                    label="环境变量"
                  />
                </div>
              </div>
            </Panel>
          ) : (
            <Empty
              title="为接口配置运行环境"
              description="创建开发、测试或生产环境，设置基础地址和对应变量。"
            />
          )}
        </div>
      </div>
      <Panel title="公共变量">
        <p className="mb-4 text-xs text-muted-foreground">
          所有环境共享；同名的环境变量会覆盖公共变量。变量替换为单层替换，不递归求值。
        </p>
        <PairEditor
          rows={draft.variables}
          onChange={(variables) => {
            setDraft((d) => ({ ...d, variables }))
            setSaved(false)
          }}
          secrets
          label="公共变量"
        />
      </Panel>
      <div className="flex items-start gap-3 rounded-lg bg-muted/60 p-4 text-xs text-muted-foreground">
        <FolderOpen size={17} className="shrink-0" />
        <div>
          <p className="break-all">配置保存在 {session.workspacePath}</p>
          <p className="mt-1">
            隐藏值仅在界面遮罩，文件仍以明文保存在本机。不要将含密钥的环境文件加入版本控制。
          </p>
        </div>
      </div>
      <Dialog open={remove} onOpenChange={setRemove}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>删除环境「{env?.name}」？</DialogTitle>
            <DialogDescription>删除后点击“保存配置”才会写入磁盘。</DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setRemove(false)}>
              取消
            </Button>
            <Button
              variant="destructive"
              onClick={() => {
                setDraft((d) => ({
                  ...d,
                  environments: d.environments.filter((e) => e.id !== selected),
                  activeEnvironment: d.activeEnvironment === selected ? '' : d.activeEnvironment,
                }))
                setSelected('')
                setRemove(false)
              }}
            >
              删除
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
