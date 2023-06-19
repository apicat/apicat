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
  TargetLanguage,
  OptionDefinition,
} from 'quicktype-core'

export interface TargetLanguageOption {
  name: string
  description: string
  defaultValue?: string | boolean
  legalValues?: string[]
  isEnumOption?: boolean
  isStringOption?: boolean
  isBooleanOption?: boolean
}

export interface CodeGenerateLanguage {
  sort: number
  logo: string
  name: string
  options: TargetLanguageOption[]
  targetLanguage: TargetLanguage
  label: string
}

const CodeGenerateSupportedLanguages: CodeGenerateLanguage[] = [
  {
    sort: 10,
    logo: 'logo.svg',
    name: 'c#',
    options: [
      { name: 'namespace', description: '生成 namespace' },
      { name: 'csharp-version', description: 'C# 版本' },
      { name: 'density', description: '代码密度' },
      { name: 'array-type', description: '使用 T[] 或 List<T>' },
      { name: 'number-type', description: 'number 数据类型' },
      { name: 'features', description: 'Output features' },
      { name: 'check-required', description: '校验 required 属性' },
      { name: 'any-type', description: 'any 使用类型' },
      { name: 'base-class', description: 'Base class' },
      { name: 'virtual', description: '生成 virtual 属性' },
    ],
    label: 'C#',
    targetLanguage: new CSharpTargetLanguage(),
  },
  {
    sort: 20,
    logo: 'logo.svg',
    name: 'go',
    options: [
      { name: 'just-types', description: '只生成 Types' },
      { name: 'package', description: '生成包名' },
      { name: 'multi-file-output', description: '将每个顶级对象呈现为其自己的 Go 文件' },
      { name: 'just-types-and-package', description: '只要 types 和 package' },
    ],
    label: 'Go',
    targetLanguage: new GoTargetLanguage(),
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
      { name: 'msbuildPermissive', description: '将 to_json 和 from_json 类型移动到 nlohmann::details 命名空间中，以便 msbuild 可以在禁用一致性模式时构建它' },
      { name: 'const-style', description: '将 const 放置在左侧/西侧 (const T) 还是右侧/东侧 (T const)' },
      { name: 'source-style', description: '源代码生成类型，是单个文件还是多个文件' },
      { name: 'include-location', description: '将 json.hpp 定位为全局还是本地文件' },
      { name: 'type-style', description: 'types 命名风格' },
      { name: 'member-style', description: 'members 命名风格' },
      { name: 'enumerator-style', description: 'enumerators 命名风格' },
      { name: 'enum-type', description: 'enum class 的类型' },
      { name: 'boost', description: '需要依赖 boost。如果没有 boost，需要 C++17 支持' },
      { name: 'hide-null-optional', description: '隐藏可选字段的 null 值' },
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
      { name: 'array-type', description: '使用 T[] 或 List<T>' },
      { name: 'just-types', description: '只生成 Types' },
      { name: 'datetime-provider', description: 'Date time provider type' },
      { name: 'acronym-style', description: '字段命名风格' },
      { name: 'package', description: '生成 package name' },
      { name: 'lombok', description: '使用 lombok' },
      { name: 'lombok-copy-annotations', description: 'Copy accessor annotations' },
    ],
    label: 'Java',
    targetLanguage: new JavaTargetLanguage(),
  },
  {
    sort: 80,
    logo: 'logo.svg',
    name: 'typescript',
    options: [
      { name: 'just-types', description: '只生成类型定义' },
      { name: 'nice-property-names', description: '属性名转为 JavaScripty 风格' },
      { name: 'explicit-unions', description: 'Explicitly name unions' },
      { name: 'runtime-typecheck', description: '运行时校验 JSON.parse 结果' },
      { name: 'runtime-typecheck-ignore-unknown-properties', description: '运行时忽略未定义的属性校验' },
      { name: 'acronym-style', description: '字段命名风格' },
      { name: 'converters', description: '使用哪种 converters 来生成 (默认为 top-level)' },
      { name: 'raw-type', description: '原始输入类型 (默认为 json)' },
      { name: 'prefer-unions', description: '使用 union type 替代枚举' },
    ],
    label: 'TypeScript',
    targetLanguage: new TypeScriptTargetLanguage(),
  },
  {
    sort: 90,
    logo: 'logo.svg',
    name: 'javascript',
    options: [
      { name: 'runtime-typecheck', description: '运行时校验 JSON.parse 结果' },
      { name: 'runtime-typecheck-ignore-unknown-properties', description: '运行时忽略未定义的属性校验' },
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
    name: 'javascript-proptypes',
    options: [
      { name: 'acronym-style', description: '字段命名风格' },
      { name: 'converters', description: '使用哪种 converters 来生成 (默认为 top-level)' },
    ],
    label: 'JavaScript PropTypes',
    targetLanguage: new JavaScriptPropTypesTargetLanguage(),
  },
  // {
  //   sort: 110,
  //   logo: 'logo.svg',
  //   name: 'flow',
  //   options: [
  //     { name: 'just-types', description: '只生成类型定义' },
  //     { name: 'nice-property-names', description: '转换属性名称为 JavaScripty 风格' },
  //     { name: 'explicit-unions', description: 'Explicitly name unions' },
  //     { name: 'runtime-typecheck', description: '运行时校验 JSON.parse 结果' },
  //     { name: 'runtime-typecheck-ignore-unknown-properties', description: '运行时忽略未定义的属性校验' },
  //     { name: 'acronym-style', description: '字段命名风格' },
  //     { name: 'converters', description: '使用哪种 converters 来生成 (默认为 top-level)' },
  //     { name: 'raw-type', description: '原始输入类型 (默认为 json)' },
  //     { name: 'prefer-unions', description: '使用 union type 替代枚举' },
  //   ],
  //   label: 'Flow',
  // },
  {
    sort: 120,
    logo: 'logo.svg',
    name: 'swift',
    options: [
      { name: 'just-types', description: '只生成 Types' },
      { name: 'struct-or-class', description: 'Structs 或 classes' },
      { name: 'density', description: '代码密度' },
      { name: 'initializers', description: 'Generate initializers and mutators' },
      { name: 'coding-keys', description: '在 Codable 类型中明确 CodingKey 值' },
      { name: 'access-level', description: 'Access level' },
      { name: 'url-session', description: 'URLSession task 扩展' },
      { name: 'alamofire', description: 'Alamofire 扩展' },
      { name: 'support-linux', description: '支持 Linux' },
      { name: 'type-prefix', description: 'type names 的前缀' },
      { name: 'protocol', description: 'Make types implement protocol' },
      { name: 'acronym-style', description: '字段命名风格' },
      { name: 'objective-c-support', description: '对象继承自 NSObject，并在类中添加 @objcMembers' },
      { name: 'swift-5-support', description: 'Renders 输出使用 Swift 5 兼容模式' },
      { name: 'multi-file-output', description: '每个顶级对象在其自己的 Swift 文件中呈现' },
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
  // {
  //   sort: 140,
  //   logo: 'logo.svg',
  //   name: 'elm',
  //   options: [
  //     { name: 'just-types', description: '只生成 Types' },
  //     { name: 'module', description: '生成 module name' },
  //     { name: 'array-type', description: '使用 Array 或 List' },
  //   ],
  //   label: 'Elm',
  // },
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
      { name: 'coders-in-class', description: '将编码器和解码器放置在类中' },
      { name: 'from-map', description: '使用 fromMap() &amp; toMap() 方法名' },
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
  // { sort: 180, logo: 'logo.svg', name: 'pike', options: [], label: 'Pike' },
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

let _chche: CodeGenerateLanguage[] | null = null

export const getAllCodeGenerateSupportedLanguages = () => {
  return _chche
    ? _chche
    : (_chche = CodeGenerateSupportedLanguages.map((item: CodeGenerateLanguage): CodeGenerateLanguage => {
        item.options = item.options.map((opt: TargetLanguageOption) => initOptionDefaultValueMapper(opt, item.targetLanguage))
        return item
      }))
}

/**
 * option 默认值处理
 * @param option
 * @param targetLanguage
 * @returns
 */
const initOptionDefaultValueMapper = (option: TargetLanguageOption, targetLanguage: TargetLanguage): TargetLanguageOption => {
  // just-types
  if (option.name === 'just-types') {
    option.defaultValue = option.defaultValue !== undefined ? option.defaultValue : true
  }

  // package
  if (option.name === 'package') {
    option.defaultValue = option.defaultValue !== undefined ? option.defaultValue : 'com.apicat'
  }

  // namespace
  if (option.name === 'namespace') {
    option.defaultValue = option.defaultValue !== undefined ? option.defaultValue : 'com.apicat'
  }

  const optionDefinition = targetLanguage.optionDefinitions.find((i: OptionDefinition) => i.name === option.name)

  if (!optionDefinition) {
    return option
  }

  option.defaultValue = option.defaultValue !== undefined ? option.defaultValue : optionDefinition.defaultValue

  if (optionDefinition.type === String && optionDefinition.legalValues) {
    option.isEnumOption = true
    option.legalValues = optionDefinition.legalValues
  } else {
    if (optionDefinition.type === String) {
      option.isStringOption = true
    }

    if (optionDefinition.type === Boolean) {
      option.isBooleanOption = true
    }
  }

  return option
}
