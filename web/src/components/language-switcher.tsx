import { Languages } from 'lucide-react'
import { useI18n, type LanguagePreference } from '@/lib/i18n'

export function LanguageSwitcher() {
  const { t, language, setLanguage } = useI18n()
  return (
    <label className="flex items-center gap-2 rounded-md border bg-card px-2.5 text-xs text-secondary-text">
      <Languages size={15} className="shrink-0" aria-hidden="true" />
      <select
        aria-label={t('界面语言')}
        className="h-9 min-w-0 flex-1 bg-transparent"
        value={language}
        onChange={(event) => setLanguage(event.target.value as LanguagePreference)}
      >
        <option value="system">{t('跟随系统')}</option>
        <option value="zh-CN">简体中文</option>
        <option value="en">English</option>
      </select>
    </label>
  )
}
