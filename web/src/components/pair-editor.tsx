import { useI18n } from '@/lib/i18n'
import { Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import type { Pair, Variable } from '@/lib/types'

export function PairEditor<T extends Pair>({
  rows,
  onChange,
  secrets = false,
  label,
}: {
  rows: T[]
  onChange: (rows: T[]) => void
  secrets?: boolean
  label?: string
}) {
  const { t } = useI18n()
  label ??= t('参数')
  function update(index: number, data: Partial<T>) {
    onChange(rows.map((r, i) => (i === index ? { ...r, ...data } : r)))
  }
  return (
    <div className="space-y-3">
      <div className="relative overflow-x-auto rounded-lg border">
        <table className="w-full min-w-[480px] text-left">
          <thead className="border-b bg-muted/50 text-xs text-muted-foreground">
            <tr>
              <th className="w-14 whitespace-nowrap p-3">{t('启用')}</th>
              <th className="p-3">{t('名称')}</th>
              <th className="p-3">{t('值')}</th>
              {secrets && <th className="w-14 p-3">{t('隐藏')}</th>}
              <th className="w-10">
                <span className="sr-only">{t('操作')}</span>
              </th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row, i) => (
              <tr key={i} className="border-b last:border-0">
                <td className="p-3">
                  <input
                    type="checkbox"
                    className="size-4 accent-[var(--accent-text)]"
                    checked={row.enabled}
                    onChange={(e) => update(i, { enabled: e.target.checked } as Partial<T>)}
                    aria-label={t('启用 {label} {index}', { label, index: i + 1 })}
                  />
                </td>
                <td className="p-2">
                  <Input
                    value={row.key}
                    aria-label={t('{label}名称 {index}', { label, index: i + 1 })}
                    placeholder="key"
                    className="font-mono text-xs"
                    onChange={(e) => update(i, { key: e.target.value } as Partial<T>)}
                  />
                </td>
                <td className="p-2">
                  <Input
                    type={secrets && (row as unknown as Variable).secret ? 'password' : 'text'}
                    value={row.value}
                    aria-label={t('{label}值 {index}', { label, index: i + 1 })}
                    placeholder={t('value 或 {{variable}}')}
                    className="font-mono text-xs"
                    autoComplete="off"
                    onChange={(e) => update(i, { value: e.target.value } as Partial<T>)}
                  />
                </td>
                {secrets && (
                  <td className="p-3">
                    <input
                      type="checkbox"
                      className="size-4 accent-[var(--accent-text)]"
                      checked={(row as unknown as Variable).secret || false}
                      onChange={(e) =>
                        update(i, { secret: e.target.checked } as unknown as Partial<T>)
                      }
                      aria-label={t('隐藏 {label} {index}', { label, index: i + 1 })}
                    />
                  </td>
                )}
                <td className="pr-2">
                  <Button
                    size="icon"
                    variant="ghost"
                    aria-label={t('删除 {label} {index}', { label, index: i + 1 })}
                    onClick={() => onChange(rows.filter((_, index) => i !== index))}
                  >
                    <Trash2 size={14} />
                  </Button>
                </td>
              </tr>
            ))}
            {rows.length === 0 && (
              <tr>
                <td
                  colSpan={secrets ? 5 : 4}
                  className="p-7 text-center text-xs text-muted-foreground"
                >
                  {t('尚未添加{label}', { label })}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      <Button
        variant="outline"
        size="sm"
        onClick={() =>
          onChange([
            ...rows,
            { key: '', value: '', enabled: true, ...(secrets ? { secret: false } : {}) } as T,
          ])
        }
      >
        <Plus size={14} />
        {t('添加{label}', { label })}
      </Button>
    </div>
  )
}
