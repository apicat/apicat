import { getMockRules } from './constants'
import MockRuleParser from './parser'

export default {
  validate: (inputRule, type) => Diff.valid(inputRule, type),
}

var Diff = {
  allMockRules: getMockRules(),
  mockType: null,
  mockRuleName: null,
  currentTypeMock: null,
  currentMockRule: null,
  currentTypeMockRules: [],

  initCurrentTypeMockRules(type) {
    this.mockType = type
    this.currentTypeMock = this.allMockRules[type]
    this.currentTypeMockRules = this.currentTypeMock ? this.currentTypeMock.rules : []
  },

  valid: function valid(inputRule, type) {
    const result = []

    this.initCurrentTypeMockRules(type)

    // 先检测是否存在为同一类型并且该类型下存在该规则 如果匹配，才有必要继续检测
    if (this.type(inputRule, result) && this.name(inputRule, result)) {
      this.validate(inputRule, result)
    }

    return result
  },

  type: function (inputRule, result) {
    var length = result.length

    if (!this.currentTypeMock) {
      result.push(`${inputRule}，不存在的Mock类型`)
    }

    return result.length === length
  },

  name: function (inputRule, result) {
    var length = result.length
    var name = MockRuleParser.getRuleName(inputRule).type
    var filter = (rule) => {
      var isExist = rule.name === name
      this.currentMockRule = isExist ? rule : null
      return isExist
    }

    Assert.oneOf(this.currentTypeMockRules, filter, result, '语法有误')

    var isValid = result.length === length
    this.mockRuleName = isValid ? name : null
    return isValid
  },

  validate: function (inputRule, result) {
    var length = result.length
    if (this.currentMockRule) {
      const { allow, alias } = this.currentMockRule

      // 精准匹配的类型
      if (allow === undefined) {
        Assert.equal(this.mockRuleName, inputRule, result, `${this.currentMockRule.name} 语法有误`)
      } else {
        const regexp = allow.regexp
        const range = allow.range
        const validate = this[`validate${(alias || this.mockType).replace(/^\S/, (s) => s.toUpperCase())}`]
        // 有正则 || 有范围
        if (regexp && regexp.test(inputRule)) {
          validate && validate.call(this, inputRule, allow, result)
        } else if (!regexp && range !== null) {
          validate && validate.call(this, inputRule, allow, result)
        } else {
          result.push(`${this.currentMockRule.name} 语法有误`)
        }
      }
    }
    return result.length === length
  },

  validateRange(inputRule, allow, result) {
    const parse = allow.parse || MockRuleParser.parseString
    const rule = parse(inputRule, allow.regexp)
    console.log(rule)

    // 范围校验
    if (allow.range !== undefined) {
      let { min, max, minActionText, maxActionText } = allow.range
      let actionText = allow.actionText || '长度'

      minActionText = minActionText || actionText
      maxActionText = maxActionText || actionText

      // |min-max
      if (rule.min !== undefined && rule.max !== undefined) {
        // isSwap 允许最大值,最小值位置交替
        if (allow.isSwap) {
          if (rule.min < min) {
            result.push(`${this.currentMockRule.name} ${minActionText}不能小于${min}`)
            return
          }

          if (rule.min > max) {
            result.push(`${this.currentMockRule.name} ${minActionText}不能大于${max}`)
            return
          }

          if (rule.max < min) {
            result.push(`${this.currentMockRule.name} ${maxActionText}不能小于${min}`)
            return
          }

          if (rule.max > max) {
            result.push(`${this.currentMockRule.name} ${maxActionText}不能大于${max}`)
            return
          }
        } else {
          if (rule.min > rule.max) {
            result.push(`${this.currentMockRule.name} 最小${actionText}不能大于最大${actionText}`)
            return
          }

          // 小于最小值
          if (rule.min < min) {
            result.push(`${this.currentMockRule.name} ${minActionText || actionText}不能小于${min}`)
            return
          }

          if (rule.max > max) {
            result.push(`${this.currentMockRule.name} ${maxActionText || actionText}不能大于${max}`)
            return
          }
        }
      }

      // count
      if (rule.min === undefined && rule.max === undefined && rule.count !== undefined) {
        if (rule.count < min) {
          result.push(`${this.currentMockRule.name} ${actionText}不能小于${min}`)
          return
        }
        if (rule.count > max) {
          result.push(`${this.currentMockRule.name} ${actionText}不能大于${max}`)
          return
        }
      }
    }

    // 格式校验
    if (allow.oneOfTypes && rule.oneOfType && allow.oneOfTypes.indexOf(rule.oneOfType) === -1) {
      let typeText = allow.typeText || '图片'
      result.push(`${this.currentMockRule.name} 不支持的${typeText}类型`)
      return
    }
  },

  validateString: function (inputRule, allow, result) {
    this.validateRange(inputRule, allow, result)
  },

  validateArrayObject: function (inputRule, allow, result) {
    this.validateRange(inputRule, allow, result)
  },

  validateArray: function (inputRule, allow, result) {
    this.validateRange(inputRule, allow, result)
  },

  validateBoolean: function (inputRule, allow, result) {
    const parse = allow.parse || MockRuleParser.parseBoolean
    const rule = parse(inputRule, allow.regexp)

    // 范围校验
    if (allow.range !== undefined && rule.probability !== undefined) {
      const { min, max } = allow.range

      // 数字范围
      if (rule.probability < min) {
        result.push(`${this.currentMockRule.name} 概率不能小于${min}`)
        return
      }

      if (rule.probability > max) {
        result.push(`${this.currentMockRule.name} 概率不能大于${max}`)
      }
    }
  },

  validateInteger: function (inputRule, allow, result) {
    const rule = MockRuleParser.parseNumber(inputRule)

    // 范围校验
    if (allow.range !== undefined) {
      let { min, max } = allow.range
      let actionText = allow.actionText || '范围'
      let prefix = allow.prefix || ''

      // |min-max
      if (rule.min !== undefined && rule.max !== undefined) {
        if (rule.min > rule.max) {
          result.push(`${this.currentMockRule.name} ${prefix}最小范围不能大于最大范围`)
          return
        }

        // 小于最小值
        if (rule.min < min) {
          result.push(`${this.currentMockRule.name} ${actionText}不能小于${min}`)
          return
        }

        if (rule.max > max) {
          result.push(`${this.currentMockRule.name} ${actionText}不能大于${max}`)
          return
        }
      }

      // count
      if (rule.min !== undefined && rule.max === undefined) {
        if (rule.count < min) {
          result.push(`${this.currentMockRule.name} ${actionText}不能小于${min}`)
          return
        }
        if (rule.count > max) {
          result.push(`${this.currentMockRule.name} ${actionText}不能大于${max}`)
          return
        }
      }

      // .decimal
      if (rule.decimal && allow.decimal) {
        const { min: dmin, max: dmax } = allow.decimal

        // dmin-dmax
        if (rule.dmin !== undefined && rule.dmax !== undefined) {
          if (rule.dmin > rule.dmax) {
            result.push(`${this.currentMockRule.name} 小数位最小位数不能大于最大位数`)
            return
          }

          // 小于最小值
          if (rule.dmin < dmin) {
            result.push(`${this.currentMockRule.name} 小数位不能小于${dmin}位`)
            return
          }

          if (rule.dmax > dmax) {
            result.push(`${this.currentMockRule.name} 小数位不能大于${dmax}位`)
            return
          }
        }

        // dcount
        if (rule.dmin !== undefined && rule.dmax === undefined) {
          // 小于最小值
          if (rule.dmin < dmin) {
            result.push(`${this.currentMockRule.name} 小数位不能小于${dmin}位`)
            return
          }

          if (rule.dmin > dmax) {
            result.push(`${this.currentMockRule.name} 小数位不能大于${dmax}位`)
            return
          }
        }
      }
    }
  },

  validateNumber: function (inputRule, allow, result) {
    this.validateInteger(inputRule, allow, result)
  },

  validateFile: function (inputRule, allow, result) {
    this.validateRange(inputRule, allow, result)
  },
}

const Assert = {
  message: function (item) {
    return (item.message || '[{utype}] Expect {action} {expected}, but is {actual}')
      .replace('{utype}', item.type.toUpperCase())
      .replace('{ltype}', item.type.toLowerCase())
      .replace('{path}', item.path)
      .replace('{action}', item.action)
      .replace('{expected}', item.expected)
      .replace('{actual}', item.actual)
  },
  equal: function (actual, expected, result, message) {
    if (actual === expected) return true
    result.push(message)
    return false
  },

  match: function (type, path, actual, expected, result, message) {
    if (expected.test(actual)) return true
    result.push(
      Assert.message({
        path: path,
        type: type,
        actual: actual,
        expected: expected,
        action: 'matches',
        message: message,
      })
    )
    return false
  },

  notEqual: function (type, path, actual, expected, result, message) {
    if (actual !== expected) return true
    var item = {
      path: path,
      type: type,
      actual: actual,
      expected: expected,
      action: 'is not equal to',
      message: message,
    }
    item.message = Assert.message(item)
    result.push(item)
    return false
  },

  greaterThan: function (type, path, actual, expected, result, message) {
    if (actual > expected) return true
    var item = {
      path: path,
      type: type,
      actual: actual,
      expected: expected,
      action: 'is greater than',
      message: message,
    }
    item.message = Assert.message(item)
    result.push(item)
    return false
  },

  lessThan: function (type, path, actual, expected, result, message) {
    if (actual < expected) return true
    var item = {
      path: path,
      type: type,
      actual: actual,
      expected: expected,
      action: 'is less to',
      message: message,
    }
    item.message = Assert.message(item)
    result.push(item)
    return false
  },

  greaterThanOrEqualTo: function (type, path, actual, expected, result, message) {
    if (actual >= expected) return true
    var item = {
      path: path,
      type: type,
      actual: actual,
      expected: expected,
      action: 'is greater than or equal to',
      message: message,
    }
    item.message = Assert.message(item)
    result.push(item)
    return false
  },

  lessThanOrEqualTo: function (type, path, actual, expected, result, message) {
    if (actual <= expected) return true
    var item = {
      path: path,
      type: type,
      actual: actual,
      expected: expected,
      action: 'is less than or equal to',
      message: message,
    }
    item.message = Assert.message(item)
    result.push(item)
    return false
  },

  oneOf: function (source, filter, result, message) {
    var res = (source || []).find(filter)
    !res && message && result.push(message)
  },
}
