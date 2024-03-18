<script setup lang="ts">
import { ClickOutside as vClickOutside } from 'element-plus'
import { useNamespace } from '@apicat/hooks'
import { useI18n } from 'vue-i18n'

const props = withDefaults(
  defineProps<{
    modelValue: string
    options: Array<{ value: string; label: string }>
    align?: 'left' | 'center'
    placeholder?: string
  }>(),
  {
    modelValue: '',
    options: () => [],
    align: 'left',
  },
)

const { t } = useI18n()

const ns = useNamespace('select')
const selected: any = useVModel(props)
const selectTriggerRef = ref()
const popoverRef = ref()

const states = reactive({
  visible: false,
})

const placeholderRef = computed(() => props.placeholder || t('acSelect.holder'))

const dropMenuVisible = computed({
  get() {
    return states.visible
  },
  set(val: boolean) {
    states.visible = val
  },
})

const popperPaneRef = computed(() => {
  return popoverRef.value?.popperRef?.contentRef
})

function toggleMenu() {
  states.visible = !states.visible
}

function handleSelect(option: any) {
  selected.value = option.value
}

function handleClose() {
  states.visible = false
}
</script>

<template>
  <div v-click-outside:[popperPaneRef]="handleClose" :class="ns.b()" class="el-select" @click.stop="toggleMenu">
    <el-popover
      ref="popoverRef"
      transition="el-zoom-in-top"
      trigger="click"
      :teleported="false"
      width="100%"
      :visible="dropMenuVisible"
      :show-arrow="false"
    >
      <template #reference>
        <div class="el-input el-input--suffix" :class="{ 'is-focus': dropMenuVisible }">
          <div ref="selectTriggerRef" class="el-input__wrapper">
            <div class="el-input__inner">
              <slot v-if="$slots.default && selected" :selected="selected">
                <span>{{ selected }}</span>
              </slot>
              <p v-else class="text-gray-helper">
                {{ placeholderRef }}
              </p>
            </div>
            <span class="el-input__suffix">
              <span class="el-input__suffix-inner">
                <el-icon class="el-select__caret el-select__icon" :class="{ 'is-reverse': dropMenuVisible }"><ac-icon-ep:arrow-down /></el-icon>
              </span>
            </span>
          </div>
        </div>
      </template>
      <template #default>
        <div :class="ns.e('options')">
          <div
            v-for="option in options"
            :key="option.value"
            :class="[ns.e('item'), { [ns.is('selected')]: selected === option.value }]"
            @click="handleSelect(option)"
          >
            <slot name="option" :option="option" />
          </div>
        </div>
      </template>
    </el-popover>
  </div>
</template>

<style lang="scss" scoped>
@use '@/styles/mixins/mixins' as *;

@include b(select) {
  .el-input__suffix {
    @apply absolute right-10px;
  }

  @include e(options) {
    @apply grid grid-cols-3 gap-4;
  }

  @include e(item) {
    @apply p-4px rounded;
    font-size: var(--el-font-size-base);
    position: relative;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    color: var(--el-text-color-regular);
    box-sizing: border-box;
    cursor: pointer;

    &:hover {
      background-color: var(--el-fill-color-light);
    }

    @include when('selected') {
      color: var(--el-color-primary);
      font-weight: bold;
      border: 1px solid var(--el-color-primary);
    }
  }
}
</style>
