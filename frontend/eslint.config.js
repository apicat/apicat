import antfu from '@antfu/eslint-config'

export default antfu(
  {
    vue: true,
    typescript: true,
  },
  {
    ignores: ['dist', 'iconfont/*'],
    rules: {
      'no-console': ['warn'],
      'vue/custom-event-name-casing': ['off'],
      'vue/prop-name-casing': ['off'],
      'ts/no-use-before-define': ['off'],
    },
  },
)
