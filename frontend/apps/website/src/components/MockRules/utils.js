import { getMockRules } from './constants'
import isPlainObject from 'lodash-es/isPlainObject'

const getMockRuleDefaultValue = function (rules, ruleName) {
  const rule = (rules || []).find((rule) => rule.name === ruleName)
  return rule ? rule.defaultValue || rule.name : ''
}

export const guessMockRule = (param) => {
  if (!isPlainObject(param)) {
    throw new Error('响应参数类型有误！')
  }

  const MOCK_RULES = getMockRules()

  const paramType = param.mockType
  const mockInfo = MOCK_RULES[paramType]

  let defaultRule = ''

  if (mockInfo) {
    const defaultMockRule = mockInfo.rules

    // 空参数名称，默认规则
    if (!param.name) {
      return getMockRuleDefaultValue(defaultMockRule, defaultMockRule[0].name)
    }

    // 精准匹配规则
    if (mockInfo.ruleKeys.indexOf(param.name) !== -1) {
      return getMockRuleDefaultValue(defaultMockRule, param.name)
    }

    // 类型检索
    const len = defaultMockRule.length
    for (let i = 0; i < len; i++) {
      const rule = defaultMockRule[i]
      if (rule.searchKey.indexOf(paramType) !== -1) {
        defaultRule = rule.defaultValue || rule.name
        break
      }
    }

    // searchkeys 再次精准匹配
    for (let i = 0; i < len; i++) {
      const rule = defaultMockRule[i]
      if (rule.searchKeys.indexOf(param.name) !== -1) {
        defaultRule = rule.defaultValue || rule.name
        break
      }
    }
  }

  return defaultRule
}
