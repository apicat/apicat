/**
 * 解析 mock 规则
 */

const RE_NORMAL_NUMBER = /^(\d|[1-9]\d+)$/

const parseName = (name, isRemoveSpace = true) => {
  name = (name === undefined ? '' : name + '').replace(isRemoveSpace ? /\s*/g : '', '')
  var rules = name.split('|')
  var type = undefined,
    lang = undefined

  if (rules[0] && rules.length <= 2 && !/\|$/.test(name)) {
    const result = rules[0].match(new RegExp(`^([\\w&]+)(\\((\\w*)\\))?$`, 'i'))
    type = result ? result[1] : undefined
    lang = result ? result[3] : undefined
  }

  return {
    rules,
    type,
    lang,
  }
}

/**
 * 图片规则解析
 */
export const parseImage = (name, regexp) => {
  const match = (name || '').match(regexp)
  let min,
    max,
    oneOfType = undefined

  if (match) {
    min = +(match[3] || match[6])
    max = +(match[4] || match[6])
  }

  return {
    min,
    max,
    oneOfType,
  }
}

/**
 * bool规则解析
 */
export const parseBoolean = (name, reg) => {
  let { type } = parseName(name)
  let probability = undefined

  if (reg) {
    let matched = name.match(reg) || []
    probability = matched[3] ? +matched[3] : undefined

    return {
      type,
      probability,
    }
  }

  return name
}

/**
 * string|{type},min,max解析
 * @param {*} name
 * @param {*} reg
 * @returns
 */
export const parseStringWithType = (name, reg) => {
  var range, min, max, count, oneOfType

  if (reg) {
    var matched = name.match(reg) || []

    oneOfType = matched[3] ? matched[3] : undefined
    range = matched[4] && matched[4].split(',')
    min = matched[5] ? +matched[5] : undefined
    max = matched[6] ? +matched[6] : undefined
    count = matched[7] ? +matched[7] : undefined
  }

  return {
    oneOfType,
    range,
    min,
    max,
    count,
  }
}

export default {
  getRuleName(name) {
    return parseName(name)
  },

  /**
   * 解析字符串类型mock规则
   */
  parseString: function (name, reg) {
    var range, min, max, count

    if (reg) {
      var matched = name.match(reg) || []

      range = matched[4] && matched[4].split(',')
      min = matched[5] ? +matched[5] : undefined
      max = matched[6] ? +matched[6] : undefined
      count = matched[7] ? +matched[7] : undefined

      return {
        lang: matched[2],
        range,
        min,
        max,
        count,
      }
    }

    var { rules, type } = parseName(name)

    var ruleRight = rules[1]

    range = ruleRight && ruleRight.split(',')

    // 正常规则
    if (range && range.length <= 2) {
      min = range[0] && RE_NORMAL_NUMBER.test(range[0]) ? +range[0] : undefined
      max = range[1] && RE_NORMAL_NUMBER.test(range[0]) ? +range[1] : undefined
      count = !range[1] && RE_NORMAL_NUMBER.test(range[0]) ? +range[0] : undefined
    }

    return {
      type,
      range,
      min,
      max,
      count,
    }
  },

  /**
   * 解析图片规则类型mock规则
   */
  parseImage: (name) => parseImage(name),

  /**
   * 解析数字类型mock规则
   */
  parseNumber: function (name) {
    var { rules, type } = parseName(name)

    var range, min, max, count, decimal, dmin, dmax, dcount

    var ruleRight = rules[1]

    var numberRules = ruleRight && ruleRight.split(',')

    if (numberRules && numberRules.length > 1) {
      range = [numberRules[0], numberRules[1]]
      min = range[0] && +range[0]
      max = range[1] && +range[1]
      count = !range[1] ? +range[0] : undefined
    }

    if (numberRules && numberRules.length >= 3) {
      dmin = numberRules[2] && +numberRules[2]
      dmax = numberRules[2] && +numberRules[2]
      dcount = !numberRules[2] ? +numberRules[1] : undefined
    }

    return {
      type,
      // 取值范围
      range,
      min,
      max,
      // min-max
      count,
      // 是否有 decimal
      decimal,
      dmin,
      dmax,
      // dmin-dimax
      dcount,
    }
  },

  parseBoolean: (name, reg) => parseBoolean(name, reg),
}
