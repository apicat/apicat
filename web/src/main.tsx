import { useI18n } from '@/lib/i18n'
import { LanguageSwitcher } from '@/components/language-switcher'
import { StrictMode, Component, type ErrorInfo, type ReactNode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import './index.css'

function RenderError() {
  const { t } = useI18n()
  return (
    <div className="mx-auto max-w-lg p-10">
      <LanguageSwitcher />
      <div className="mt-6">
        <h1 className="text-xl font-semibold">{t('页面渲染失败')}</h1>
        <p className="mt-3">{t('请刷新页面；如问题持续，请检查文档或向项目反馈。')}</p>
        <button className="mt-4 underline" onClick={() => window.location.reload()}>
          {t('重新加载')}
        </button>
      </div>
    </div>
  )
}

class ErrorBoundary extends Component<{ children: ReactNode }, { error: boolean }> {
  state = { error: false }
  static getDerivedStateFromError() {
    return { error: true }
  }
  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error(error, info.componentStack)
  }
  render() {
    return this.state.error ? <RenderError /> : this.props.children
  }
}
createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ErrorBoundary>
      <App />
    </ErrorBoundary>
  </StrictMode>,
)
