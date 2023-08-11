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
    'transition-base': 'transition-all duration-300 ease-in-out',
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
        lighter: '#eceeef',
        100: '#fafafa',
        110: '#f2f2f2',
        title: '#101010',
        helper: '#9C9898',
      },
    },
    height: {
      14: '48px',
    },
  },
})
