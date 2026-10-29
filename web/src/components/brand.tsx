import lightWordmark from '@/assets/brand/apicat-wordmark-light-transparent-1920.png'
import darkWordmark from '@/assets/brand/apicat-wordmark-dark-transparent-1920.png'
import lightIcon from '@/assets/brand/apicat-icon-light-64.png'
import darkIcon from '@/assets/brand/apicat-icon-dark-64.png'

export function Brand({ compact = false }: { compact?: boolean }) {
  return (
    <span aria-hidden="true" className="block shrink-0">
      <img
        src={compact ? lightIcon : lightWordmark}
        width={compact ? 32 : 112}
        height={compact ? 32 : 36.75}
        className={compact ? 'brand-light size-8' : 'brand-light brand-wordmark'}
        alt=""
      />
      <img
        src={compact ? darkIcon : darkWordmark}
        width={compact ? 32 : 112}
        height={compact ? 32 : 36.75}
        className={compact ? 'brand-dark size-8' : 'brand-dark brand-wordmark'}
        alt=""
      />
    </span>
  )
}
