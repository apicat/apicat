<template>
  <el-dialog title="Mock 规则" v-model="visible" :close-on-click-modal="false" :close-on-press-escape="false" class="ac-mock" append-to-body>
    <div class="ac-mock-container">
      <el-form :model="form" :rules="rules" ref="mockForm" @submit.prevent="onOkBtnClick" @keyup.enter="onOkBtnClick">
        <el-form-item label="" prop="mock_rule">
          <el-input placeholder="Mock Rules" v-model="form.mock_rule" clearable />
        </el-form-item>
      </el-form>

      <div class="ac-mock-layout">
        <div class="ac-mock-layout__left">
          <el-input
            :prefix-icon="search"
            size="small"
            class="mb-10px"
            v-model="mockQuery"
            @keydown.down.stop.prevent="onSelectMockRuleByKeyDown('next')"
            @keydown.up.stop.prevent="onSelectMockRuleByKeyDown('prev')"
            @keydown.enter.prevent="onSelectedMockRuleByEnter"
            placeholder="关键字"
            clearable
          />
          <ul class="ac-mock-rule__list">
            <li
              v-for="(rule, idx) in mockRules"
              :key="rule.key"
              @mouseenter="onHoverMockRuleItem($event, idx)"
              :id="'ac-mock-rule__item' + rule.key"
              :class="ruleClass(rule)"
              @click="onRuleItemClick(rule)"
            >
              {{ rule.cnName }}({{ rule.name }})
            </li>
          </ul>
        </div>
        <div class="ac-mock-layout__right" v-if="currentRule">
          <div class="ac-mock-rule__syntax">
            <p class="ac-mock-rule__label">语法：</p>
            <div class="ac-mock-rule__desc" v-html="currentRule.syntax" />
          </div>

          <div class="ac-mock-rule__eg">
            <p class="ac-mock-rule__label">举例：</p>
            <div class="ac-mock-rule__desc" v-html="currentRule.example" />
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <el-button @click="visible = false">取 消</el-button>
      <el-button type="primary" @click="onOkBtnClick">确 定</el-button>
    </template>
  </el-dialog>
</template>

<script>
import { ElButton, ElDialog, ElForm, ElFormItem, ElInput } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getMockRules } from './constants'
import MockRuleParser from './parser'
import MockValidator from './valid'
import { debounce } from 'lodash-es'
import scrollIntoView from 'smooth-scroll-into-view-if-needed'
import { markRaw } from 'vue'
import { $emit } from '@apicat/shared'

const MOCK_RULES = getMockRules()

let lastCursorPos = {
  x: 0,
  y: 0,
}

export default defineComponent({
  name: 'MockRules',
  components: {
    ElButton,
    ElDialog,
    ElForm,
    ElFormItem,
    ElInput,
  },
  emits: ['on-ok'],
  data() {
    return {
      search: markRaw(Search),
      visible: false,
      paramType: null,
      currentRule: null,
      hoverRule: null,
      hoverIndex: -1,
      mockQuery: '',
      mockRules: [],

      allMockRules: [],

      form: {
        mock_rule: '',
      },

      rules: {
        mock_rule: {
          validator: debounce((rule, value, callback) => {
            const val = value.trim()

            if (!val) {
              return callback(new Error('请输入Mock规则'))
            }
            const result = MockValidator.validate(value, this.paramType)
            return result.length ? callback(new Error(result[0])) : callback()
          }, 200),
          trigger: 'change',
        },
      },
    }
  },
  watch: {
    mockQuery: debounce(function () {
      const query = this.mockQuery.replace(/\s+/g, '')
      this.mockRules = !query ? this.allMockRules.concat([]) : this.allMockRules.filter((rule) => rule.searchKey.indexOf(query) !== -1)
      this.hoverIndex = -1
    }, 300),

    visible: function () {
      !this.visible && this.reset()
    },
    'form.mock_rule': debounce(function () {
      if (this.form.mock_rule) {
        const result = MockValidator.validate(this.form.mock_rule, this.paramType)
        !result.length && this.getCurrentMockRule(this.form.mock_rule)
      }
    }, 300),
  },

  methods: {
    onHoverMockRuleItem(e, idx) {
      // cursor didn't move
      if (e.screenX === lastCursorPos.x && e.screenY === lastCursorPos.y) {
        return
      }

      lastCursorPos = {
        x: e.screenX,
        y: e.screenY,
      }

      this.hoverIndex = idx

      this.hoverRule = this.mockRules[this.hoverIndex]
    },

    onSelectMockRuleByKeyDown(direction) {
      const len = this.mockRules.length

      if (direction === 'next') {
        this.hoverIndex++
      } else {
        this.hoverIndex--
      }

      if (this.hoverIndex >= len) {
        this.hoverIndex = 0
      }

      if (this.hoverIndex <= -1) {
        this.hoverIndex = len - 1
      }

      this.hoverRule = this.mockRules[this.hoverIndex]

      this.hoverRule &&
        this.$nextTick(() => {
          let node = document.getElementById('ac-mock-rule__item' + this.hoverRule.key)
          if (node) {
            node.scrollIntoView ? node.scrollIntoView() : scrollIntoView(node)
          }
        })
    },

    onSelectedMockRuleByEnter() {
      if (this.hoverIndex !== -1 && this.mockRules[this.hoverIndex]) {
        this.onRuleItemClick(this.mockRules[this.hoverIndex])
      }
    },

    ruleClass: function (rule) {
      return [
        'ac-mock-rule__item',
        {
          current: this.currentRule && this.currentRule.key === rule.key,
          hover: this.currentRule && this.hoverRule && this.hoverRule.key === rule.key && this.currentRule.key !== this.hoverRule.key,
        },
      ]
    },

    onRuleItemClick(rule) {
      this.currentRule = rule
      this.form.mock_rule = rule.defaultValue || rule.name
    },

    show(node) {
      this.visible = true
      this.form.mock_rule = node.mockRule || ''
      this.generateMockRulesByParamType(node.mockType)
      this.getCurrentMockRule(node.mockRule)
    },

    generateMockRulesByParamType(mockType) {
      const defaultMockRule = (MOCK_RULES[mockType] || { rules: [] }).rules
      let rules = []

      if (defaultMockRule && defaultMockRule.length) {
        rules = defaultMockRule
        this.paramType = mockType
      } else {
        Object.keys(MOCK_RULES).forEach((key) => {
          let _rules = MOCK_RULES[key].rules || []
          rules = rules.concat(_rules)
        })
      }

      this.allMockRules = rules
      this.mockRules = this.allMockRules.concat([])
    },

    getCurrentMockRule(rule) {
      const { type: ruleName } = MockRuleParser.getRuleName(rule)
      const key = `${this.paramType}-${ruleName}`
      this.currentRule = this.allMockRules.find((item) => item.key === key)

      if (!this.currentRule) {
        this.onRuleItemClick(this.allMockRules[0])
      }

      this.$nextTick(() => {
        let node = document.getElementById('ac-mock-rule__item' + this.currentRule.key)
        if (node) {
          node.scrollIntoView ? node.scrollIntoView({ behavior: 'smooth', block: 'nearest' }) : scrollIntoView(node, { behavior: 'smooth', block: 'nearest' })
        }
      })
    },

    onOkBtnClick() {
      this.$refs['mockForm'].validate((valid) => {
        if (valid) {
          $emit(this, 'on-ok', this.form.mock_rule)
          this.visible = false
        }
      })
    },

    reset() {
      this.$refs['mockForm'].resetFields()
      this.currentRule = null
      this.form.mock_rule = ''
      this.mockQuery = ''
    },
  },
})
</script>

<style lang="scss">
@use './style.scss';
</style>
