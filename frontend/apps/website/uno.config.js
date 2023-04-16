import { defineConfig, presetUno, presetAttributify } from 'unocss'
import transformerDirectives from '@unocss/transformer-directives'

export default defineConfig({
  presets: [presetUno(), presetAttributify()],
  transformers: [transformerDirectives()],
  rules: [['rounded', { 'border-radius': '5px' }]],
  shortcuts: {
    'wh-full': 'w-full h-full',
    'flex-center': 'flex justify-center items-center',
    'flex-x-center': 'flex justify-center',
    'flex-y-center': 'flex items-center',
  },
  theme: {
    colors: {
      blue: {
        primary: '#006bff',
      },
      gray: {
        '06': 'rgba(0,0,0,.06)',
        65: 'rgba(0,0,0,.65)',
        45: 'rgba(0,0,0,.45)',
        100: '#f5f7fa',
      },
    },
    height: {
      14: '48px',
    },
  },
})
