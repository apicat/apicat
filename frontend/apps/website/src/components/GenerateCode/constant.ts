import type { OptionDefinition, TargetLanguage } from 'quicktype-core'

import {
  CPlusPlusTargetLanguage,
  CSharpTargetLanguage,
  DartTargetLanguage,
  GoTargetLanguage,
  HaskellTargetLanguage,
  JavaScriptPropTypesTargetLanguage,
  JavaScriptTargetLanguage,
  JavaTargetLanguage,
  KotlinTargetLanguage,
  ObjectiveCTargetLanguage,
  PythonTargetLanguage,
  RubyTargetLanguage,
  RustTargetLanguage,
  SwiftTargetLanguage,
  TypeScriptTargetLanguage,
} from 'quicktype-core'

import { DEFAULT_LANGUAGE } from '@/commons'

export interface TargetLanguageOption {
  name: string
  description: string
  defaultValue?: string | boolean
  legalValues?: string[]
  isEnumOption?: boolean
  isStringOption?: boolean
  isBooleanOption?: boolean
  isPrimaryOption?: boolean
}

export interface CodeGenerateLanguage {
  sort: number
  logo: string
  name: string
  options: TargetLanguageOption[]
  primaryOptions?: TargetLanguageOption[]
  secondaryOptions?: TargetLanguageOption[]
  targetLanguage: TargetLanguage | null
  label: string
}

export interface CodeGenerateLanguageWithI18n {
  lang: string
  alias: string[]
  codeGenerateLanguage: CodeGenerateLanguage[]
}

// zh-CN: 生成代码支持的语言
const CodeGenerateSupportedLanguagesWithZhCN: CodeGenerateLanguage[] = [
  {
    sort: 1,
    logo: 'logo.svg',
    name: 'JSON5',
    options: [],
    label: 'JSON',
    targetLanguage: null,
  },

  {
    sort: 2,
    logo: 'logo.svg',
    name: 'JSON',
    options: [],
    label: 'JSONSchema',
    targetLanguage: null,
  },
  {
    sort: 3,
    logo: 'logo.svg',
    name: 'typescript',
    options: [
      { name: 'just-types', description: '只生成类型定义' },
      { name: 'nice-property-names', description: '属性名转为 JavaScripty 风格' },
      // { name: 'explicit-unions', description: 'Explicitly name unions' },
      // { name: 'runtime-typecheck', description: '运行时校验 JSON.parse 结果' },
      // { name: 'runtime-typecheck-ignore-unknown-properties', description: '运行时忽略未定义的属性校验' },
      { name: 'acronym-style', description: '字段命名风格' },
      // { name: 'converters', description: '使用哪种 converters 来生成 (默认为 top-level)' },
      // { name: 'raw-type', description: '原始输入类型 (默认为 json)' },
      // { name: 'prefer-unions', description: '使用 union type 替代枚举' },
    ],
    label: 'TypeScript',
    targetLanguage: new TypeScriptTargetLanguage(),
  },
  {
    sort: 5,
    logo: 'logo.svg',
    name: 'go',
    options: [
      { name: 'just-types', description: '只生成 Types' },
      { name: 'package', description: '生成包名' },
      { name: 'multi-file-output', description: '将每个顶级对象呈现为其自己的Go文件', defaultValue: true },
      { name: 'just-types-and-package', description: '只要 types 和 package' },
    ],
    label: 'Go',
    targetLanguage: new GoTargetLanguage(),
  },
  {
    sort: 10,
    logo: 'logo.svg',
    name: 'c#',
    options: [
      { name: 'namespace', description: '生成 namespace' },
      { name: 'csharp-version', description: 'C# 版本' },
      { name: 'density', description: '代码密度' },
      // { name: 'array-type', description: '使用 T[] 或 List<T>' },
      // { name: 'number-type', description: 'number 数据类型' },
      { name: 'features', description: 'Output features' },
      // { name: 'check-required', description: '校验 required 属性' },
      // { name: 'any-type', description: 'any 使用类型' },
      // { name: 'base-class', description: 'Base class' },
      // { name: 'virtual', description: '生成 virtual 属性' },
    ],
    label: 'C#',
    targetLanguage: new CSharpTargetLanguage(),
  },
  {
    sort: 30,
    logo: 'logo.svg',
    name: 'rust',
    options: [
      { name: 'density', description: '代码密度' },
      { name: 'visibility', description: '字段可见性' },
      { name: 'derive-debug', description: 'Derive Debug impl' },
    ],
    label: 'Rust',
    targetLanguage: new RustTargetLanguage(),
  },
  {
    sort: 50,
    logo: 'logo.svg',
    name: 'c++',
    options: [
      { name: 'just-types', description: '只生成 Types' },
      { name: 'namespace', description: '生成的 namespace' },
      { name: 'code-format', description: '生成带 getters/setters 的 class, 而不是 structs' },
      { name: 'wstring', description: '使用 Utf-16 std::wstring 存储 strings, 而不是 Utf-8 std::string' },
      { name: 'const-style', description: '将 const 放置在左侧/西侧 (const T) 还是右侧/东侧 (T const)' },
      { name: 'type-style', description: 'types 命名风格' },
      { name: 'member-style', description: 'members 命名风格' },
      { name: 'enumerator-style', description: 'enumerators 命名风格' },
      { name: 'boost', description: '需要依赖 boost。如果没有 boost，需要 C++17 支持' },
    ],
    label: 'C++',
    targetLanguage: new CPlusPlusTargetLanguage(),
  },
  {
    sort: 60,
    logo: 'logo.svg',
    name: 'objective-c',
    options: [
      { name: 'just-types', description: '只生成 Types' },
      { name: 'class-prefix', description: 'Class prefix' },
      { name: 'features', description: 'Interface 和 implementation' },
      { name: 'extra-comments', description: 'Extra comments' },
      { name: 'functions', description: 'C语言风格 functions' },
    ],
    label: 'Objective-C',
    targetLanguage: new ObjectiveCTargetLanguage(),
  },
  {
    sort: 70,
    logo: 'logo.svg',
    name: 'java',
    options: [
      { name: 'just-types', description: '只生成 Types' },
      { name: 'acronym-style', description: '字段命名风格' },
      { name: 'package', description: '生成 package name' },
      { name: 'lombok', description: '使用 lombok' },
    ],
    label: 'Java',
    targetLanguage: new JavaTargetLanguage(),
  },

  {
    sort: 90,
    logo: 'logo.svg',
    name: 'javascript',
    options: [
      { name: 'runtime-typecheck', description: '运行时校验 JSON.parse 结果', defaultValue: false },
      { name: 'acronym-style', description: '字段命名风格' },
      { name: 'converters', description: '使用哪种 converters 来生成 (默认为 top-level)' },
      { name: 'raw-type', description: '原始输入类型 (默认为 json)' },
    ],
    label: 'JavaScript',
    targetLanguage: new JavaScriptTargetLanguage(),
  },
  {
    sort: 100,
    logo: 'logo.svg',
    name: 'javascript',
    options: [
      { name: 'acronym-style', description: '字段命名风格' },
      { name: 'converters', description: '使用哪种 converters 来生成 (默认为 top-level)' },
    ],
    label: 'JavaScript PropTypes',
    targetLanguage: new JavaScriptPropTypesTargetLanguage(),
  },
  {
    sort: 120,
    logo: 'logo.svg',
    name: 'swift',
    options: [
      { name: 'just-types', description: '只生成 Types' },
      { name: 'struct-or-class', description: 'Structs 或 classes' },
      { name: 'access-level', description: 'Access level' },
      { name: 'type-prefix', description: 'type names 的前缀' },
      { name: 'protocol', description: 'Make types implement protocol' },
      { name: 'acronym-style', description: '字段命名风格' },
      { name: 'multi-file-output', description: '每个顶级对象在其自己的 Swift 文件中呈现', defaultValue: true },
      { name: 'mutable-properties', description: 'object 的属性使用 var 替代 let' },
    ],
    label: 'Swift',
    targetLanguage: new SwiftTargetLanguage(),
  },
  {
    sort: 130,
    logo: 'logo.svg',
    name: 'kotlin',
    options: [
      { name: 'framework', description: 'Serialization framework' },
      { name: 'package', description: 'Package' },
    ],
    label: 'Kotlin',
    targetLanguage: new KotlinTargetLanguage(),
  },
  {
    sort: 150,
    logo: 'logo.svg',
    name: 'ruby',
    options: [
      { name: 'just-types', description: '只生成 Types' },
      { name: 'strictness', description: '严格模式 Type' },
    ],
    label: 'Ruby',
    targetLanguage: new RubyTargetLanguage(),
  },
  {
    sort: 160,
    logo: 'logo.svg',
    name: 'dart',
    options: [
      { name: 'just-types', description: '只生成 Types' },
      { name: 'from-map', description: '使用 fromMap() & toMap() 方法名' },
      { name: 'required-props', description: '所有属性设置为 required' },
      { name: 'final-props', description: '所有属性设置为 final' },
      { name: 'copy-with', description: '生成 CopyWith 方法' },
      { name: 'use-freezed', description: '生成 @freezed 兼容的 class 定义' },
      { name: 'use-hive', description: '生成 Hive type adapters 注解' },
      { name: 'part-name', description: '在 `part` 指令中使用此名称' },
    ],
    label: 'Dart',
    targetLanguage: new DartTargetLanguage(),
  },
  {
    sort: 170,
    logo: 'logo.svg',
    name: 'python',
    options: [
      { name: 'python-version', description: 'Python 版本' },
      { name: 'just-types', description: '只生成 Class' },
      { name: 'nice-property-names', description: '属性名转成 Pythonic 风格' },
    ],
    label: 'Python',
    targetLanguage: new PythonTargetLanguage('Python', ['python', 'py'], 'py'),
  },
  {
    sort: 190,
    logo: 'logo.svg',
    name: 'haskell',
    options: [
      { name: 'just-types', description: '只生成 Types' },
      { name: 'module', description: '生成 module name' },
      { name: 'array-type', description: '使用 Array 或 List' },
    ],
    label: 'Haskell',
    targetLanguage: new HaskellTargetLanguage(),
  },
]

// en-US: Supported languages for code generation
const CodeGenerateSupportedLanguagesWithEnUS: CodeGenerateLanguage[] = [
  {
    sort: 1,
    logo: 'logo.svg',
    name: 'JSON5',
    options: [],
    label: 'JSON',
    targetLanguage: null,
  },

  {
    sort: 2,
    logo: 'logo.svg',
    name: 'JSON',
    options: [],
    label: 'JSONSchema',
    targetLanguage: null,
  },
  {
    sort: 3,
    logo: 'logo.svg',
    name: 'typescript',
    options: [
      { name: 'just-types', description: 'Only generate type definitions' },
      { name: 'nice-property-names', description: 'Convert property names to JavaScripty style' },
      // { name: 'explicit-unions', description: 'Explicitly name unions' },
      // { name: 'runtime-typecheck', description: 'Verify JSON.parse results at runtime' },
      // { name: 'runtime-typecheck-ignore-unknown-properties', description: 'Ignore undefined property checks at runtime' },
      { name: 'acronym-style', description: 'Field naming style' },
      // { name: 'converters', description: 'Which converters to use for generation (default is top-level)' },
      // { name: 'raw-type', description: 'Raw input type (default is json)' },
      // { name: 'prefer-unions', description: 'Use union type instead of enumeration' },
    ],
    label: 'TypeScript',
    targetLanguage: new TypeScriptTargetLanguage(),
  },
  {
    sort: 5,
    logo: 'logo.svg',
    name: 'go',
    options: [
      { name: 'just-types', description: 'Only generate Types' },
      { name: 'package', description: 'Generate package name' },
      { name: 'multi-file-output', description: 'Render each top-level object as its own Go file', defaultValue: true },
      { name: 'just-types-and-package', description: 'Just types and package' },
    ],
    label: 'Go',
    targetLanguage: new GoTargetLanguage(),
  },
  {
    sort: 10,
    logo: 'logo.svg',
    name: 'c#',
    options: [
      { name: 'namespace', description: 'Generate namespace' },
      { name: 'csharp-version', description: 'C# version' },
      { name: 'density', description: 'code density' },
      // { name: 'array-type', description: 'Use T[] or List<T>' },
      // { name: 'number-type', description: 'number data type' },
      { name: 'features', description: 'Output features' },
      // { name: 'check-required', description: 'Check required attribute' },
      // { name: 'any-type', description: 'any usage type' },
      // { name: 'base-class', description: 'Base class' },
      // { name: 'virtual', description: 'Generate virtual properties' },
    ],
    label: 'C#',
    targetLanguage: new CSharpTargetLanguage(),
  },
  {
    sort: 30,
    logo: 'logo.svg',
    name: 'rust',
    options: [
      { name: 'density', description: 'code density' },
      { name: 'visibility', description: 'Field visibility' },
      { name: 'derive-debug', description: 'Derive Debug impl' },
    ],
    label: 'Rust',
    targetLanguage: new RustTargetLanguage(),
  },
  {
    sort: 50,
    logo: 'logo.svg',
    name: 'c++',
    options: [
      { name: 'just-types', description: 'Only generate Types' },
      { name: 'namespace', description: 'Generated namespace' },
      { name: 'code-format', description: 'Generate classes with getters/setters instead of structs' },
      { name: 'wstring', description: 'Use Utf-16 std::wstring to store strings instead of Utf-8 std::string' },
      // { name: 'msbuildPermissive', description: 'Move to_json and from_json types into the nlohmann::details namespace so that msbuild can build it when consistency mode is disabled' },
      { name: 'const-style', description: 'Whether const is placed on the left/west side (const T) or on the right/east side (T const)' },
      // { name: 'source-style', description: 'Source code generation type, whether it is a single file or multiple files' },
      // { name: 'include-location', description: 'Locate json.hpp as a global or local file' },
      { name: 'type-style', description: 'types naming style' },
      { name: 'member-style', description: 'members naming style' },
      { name: 'enumerator-style', description: 'enumerators naming style' },
      // { name: 'enum-type', description: 'enum class type' },
      { name: 'boost', description: 'Needs to depend on boost. If there is no boost, C++17 support is required' },
      // { name: 'hide-null-optional', description: 'Hide the null value of an optional field' },
    ],
    label: 'C++',
    targetLanguage: new CPlusPlusTargetLanguage(),
  },
  {
    sort: 60,
    logo: 'logo.svg',
    name: 'objective-c',
    options: [
      { name: 'just-types', description: 'Only generate Types' },
      { name: 'class-prefix', description: 'Class prefix' },
      { name: 'features', description: 'Interface and implementation' },
      { name: 'extra-comments', description: 'Extra comments' },
      { name: 'functions', description: 'C language style functions' },
    ],
    label: 'Objective-C',
    targetLanguage: new ObjectiveCTargetLanguage(),
  },
  {
    sort: 70,
    logo: 'logo.svg',
    name: 'java',
    options: [
      // { name: 'array-type', description: 'Use T[] or List<T>' },
      { name: 'just-types', description: 'Only generate Types' },
      // { name: 'datetime-provider', description: 'Date time provider type' },
      { name: 'acronym-style', description: 'Field naming style' },
      { name: 'package', description: 'Generate package name' },
      { name: 'lombok', description: 'Use lombok' },
      // { name: 'lombok-copy-annotations', description: 'Copy accessor annotations' },
    ],
    label: 'Java',
    targetLanguage: new JavaTargetLanguage(),
  },

  {
    sort: 90,
    logo: 'logo.svg',
    name: 'javascript',
    options: [
      { name: 'runtime-typecheck', description: 'Verify JSON.parse results at runtime', defaultValue: false },
      // { name: 'runtime-typecheck-ignore-unknown-properties', description: 'Ignore undefined property checks at runtime' },
      { name: 'acronym-style', description: 'Field naming style' },
      { name: 'converters', description: 'Which converters to use for generation (default is top-level)' },
      { name: 'raw-type', description: 'Raw input type (default is json)' },
    ],
    label: 'JavaScript',
    targetLanguage: new JavaScriptTargetLanguage(),
  },
  {
    sort: 100,
    logo: 'logo.svg',
    name: 'javascript',
    options: [
      { name: 'acronym-style', description: 'Field naming style' },
      { name: 'converters', description: 'Which converters to use for generation (default is top-level)' },
    ],
    label: 'JavaScript PropTypes',
    targetLanguage: new JavaScriptPropTypesTargetLanguage(),
  },
  // {
  // sort: 110,
  // logo: 'logo.svg',
  // name: 'flow',
  // options: [
  // { name: 'just-types', description: 'Only generate type definitions' },
  // { name: 'nice-property-names', description: 'Convert property names to JavaScripty style' },
  // { name: 'explicit-unions', description: 'Explicitly name unions' },
  // { name: 'runtime-typecheck', description: 'Verify JSON.parse results at runtime' },
  // { name: 'runtime-typecheck-ignore-unknown-properties', description: 'Ignore undefined property checks at runtime' },
  // { name: 'acronym-style', description: 'Field naming style' },
  // { name: 'converters', description: 'Which converters to use for generation (default is top-level)' },
  // { name: 'raw-type', description: 'Raw input type (default is json)' },
  // { name: 'prefer-unions', description: 'Use union type instead of enumeration' },
  // ],
  // label: 'Flow',
  // },
  {
    sort: 120,
    logo: 'logo.svg',
    name: 'swift',
    options: [
      { name: 'just-types', description: 'Only generate Types' },
      { name: 'struct-or-class', description: 'Structs or classes' },
      // { name: 'density', description: 'code density' },
      // { name: 'initializers', description: 'Generate initializers and mutators' },
      // { name: 'coding-keys', description: 'Explicit CodingKey value in Codable type' },
      { name: 'access-level', description: 'Access level' },
      // { name: 'url-session', description: 'URLSession task extension' },
      // { name: 'alamofire', description: 'Alamofire extension' },
      // { name: 'support-linux', description: 'Support Linux' },
      { name: 'type-prefix', description: 'prefix of type names' },
      { name: 'protocol', description: 'Make types implement protocol' },
      { name: 'acronym-style', description: 'Field naming style' },
      // { name: 'objective-c-support', description: 'Object inherits from NSObject and adds @objcMembers to the class' },
      // { name: 'swift-5-support', description: 'Renders output uses Swift 5 compatibility mode' },
      { name: 'multi-file-output', description: 'Each top-level object is rendered in its own Swift file', defaultValue: true },
      { name: 'mutable-properties', description: 'Use var instead of let\' for object properties' },
    ],
    label: 'Swift',
    targetLanguage: new SwiftTargetLanguage(),
  },
  {
    sort: 130,
    logo: 'logo.svg',
    name: 'kotlin',
    options: [
      { name: 'framework', description: 'Serialization framework' },
      { name: 'package', description: 'Package' },
    ],
    label: 'Kotlin',
    targetLanguage: new KotlinTargetLanguage(),
  },
  {
    sort: 150,
    logo: 'logo.svg',
    name: 'ruby',
    options: [
      { name: 'just-types', description: 'Only generate Types' },
      { name: 'strictness', description: 'Strict Mode Type' },
    ],
    label: 'Ruby',
    targetLanguage: new RubyTargetLanguage(),
  },
  {
    sort: 160,
    logo: 'logo.svg',
    name: 'dart',
    options: [
      { name: 'just-types', description: 'Only generate Types' },
      { name: 'from-map', description: 'Use fromMap() & toMap() method name' },
      { name: 'required-props', description: 'Set all properties to required' },
      { name: 'final-props', description: 'Set all properties to final' },
      { name: 'copy-with', description: 'Generate CopyWith method' },
      { name: 'use-freezed', description: 'Generate @freezed compatible class definition' },
      { name: 'use-hive', description: 'Generate Hive type adapters annotations' },
      { name: 'part-name', description: 'Use this name in the `part` directive' },
    ],
    label: 'Dart',
    targetLanguage: new DartTargetLanguage(),
  },
  {
    sort: 170,
    logo: 'logo.svg',
    name: 'python',
    options: [
      { name: 'python-version', description: 'Python version' },
      { name: 'just-types', description: 'Only generate Class' },
      { name: 'nice-property-names', description: 'Convert property names to Pythonic style' },
    ],
    label: 'Python',
    targetLanguage: new PythonTargetLanguage('Python', ['python', 'py'], 'py'),
  },
  {
    sort: 190,
    logo: 'logo.svg',
    name: 'haskell',
    options: [
      { name: 'just-types', description: 'Only generate Types' },
      { name: 'module', description: 'Generate module name' },
      { name: 'array-type', description: 'Use Array or List' },
    ],
    label: 'Haskell',
    targetLanguage: new HaskellTargetLanguage(),
  },
]

// zh-CN: 生成代码支持的语言
const ZhCNCodeGenerateConfig: CodeGenerateLanguageWithI18n = {
  lang: 'zh-CN',
  alias: ['zh', 'zh-CN', 'zh-Hans', 'zh-Hans-CN'],
  codeGenerateLanguage: CodeGenerateSupportedLanguagesWithZhCN,
}

// en-US: Supported languages for code generation
const EnUSCodeGenerateConfig: CodeGenerateLanguageWithI18n = {
  lang: 'en-US',
  alias: ['en', 'en-US'],
  codeGenerateLanguage: CodeGenerateSupportedLanguagesWithEnUS,
}

const allCodeGenerateConfig: CodeGenerateLanguageWithI18n[] = [ZhCNCodeGenerateConfig, EnUSCodeGenerateConfig]

const cacheMap: Map<string, CodeGenerateLanguage[]> = new Map()

export function getAllCodeGenerateSupportedLanguages(lang: string = DEFAULT_LANGUAGE): CodeGenerateLanguage[] {
  if (!cacheMap.has(lang)) {
    const config = allCodeGenerateConfig.find((item: CodeGenerateLanguageWithI18n) => item.lang === lang || item.alias.includes(lang))

    if (config) {
      config.codeGenerateLanguage.map((item: CodeGenerateLanguage): CodeGenerateLanguage => {
        const options = item.targetLanguage ? item.options.map((opt: TargetLanguageOption) => initOptionDefaultValueMapper(opt, item.targetLanguage!)) : []
        const primaryOpt: TargetLanguageOption[] = []
        const secondaryOpt: TargetLanguageOption[] = []

        // 分组options
        options.forEach((i: TargetLanguageOption) => {
          if (i.isPrimaryOption)
            primaryOpt.push(i)
          else
            secondaryOpt.push(i)
        })

        // 排序options
        const compare = (pre: TargetLanguageOption) => (pre.isBooleanOption ? 1 : -1)
        primaryOpt.sort(compare)
        secondaryOpt.sort(compare)

        // 重新赋值
        item.primaryOptions = primaryOpt
        item.secondaryOptions = secondaryOpt
        item.options = [...primaryOpt, ...secondaryOpt]
        return item
      })

      cacheMap.set(lang, config.codeGenerateLanguage)
    }
  }

  return cacheMap.get(lang) || CodeGenerateSupportedLanguagesWithEnUS
}

/**
 * option 默认值处理
 * @param option
 * @param targetLanguage
 */
function initOptionDefaultValueMapper(option: TargetLanguageOption, targetLanguage: TargetLanguage): TargetLanguageOption {
  // just-types
  if (option.name === 'just-types')
    option.defaultValue = option.defaultValue !== undefined ? option.defaultValue : true

  // package
  if (option.name === 'package')
    option.defaultValue = option.defaultValue !== undefined ? option.defaultValue : 'com.apicat'

  // namespace
  if (option.name === 'namespace')
    option.defaultValue = option.defaultValue !== undefined ? option.defaultValue : 'com.apicat'

  const optionDefinition = targetLanguage.optionDefinitions.find((i: OptionDefinition) => i.name === option.name)

  // not support option
  if (!optionDefinition) {
    // console.warn(`not support option: ${option.name}`)
    return option
  }

  option.defaultValue = option.defaultValue !== undefined ? option.defaultValue : optionDefinition.defaultValue

  option.isPrimaryOption = optionDefinition.kind === 'primary'

  if (optionDefinition.type === String && optionDefinition.legalValues) {
    option.isEnumOption = true
    option.legalValues = optionDefinition.legalValues
  }
  else {
    if (optionDefinition.type === String)
      option.isStringOption = true

    if (optionDefinition.type === Boolean)
      option.isBooleanOption = true
  }

  return option
}
