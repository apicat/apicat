import elementPlusLocale from 'element-plus/lib/locale/lang/zh-cn'

export default {
  name: '中文',
  app: {
    common: {
      add: '添加',
      confirm: '确定',
      edit: '编辑',
      save: '保存',
      export: '导出',
      import: '导入',
      cancel: '取消',
      restore: '恢复',
      reanme: '重命名',
      copy: '复制',
      create: '创建',
      generate: '生成',
      emptyDataTip: '暂无数据',
      delete: '删除',
      deleteTip: '删除提示',
      confirmDelete: '确认删除{msg}吗？',
      copyAllPath: '复制完整URL',
      fetchMockData: '获取Mock数据',
      saving: '保存中...',
      savedCloud: '已保存在云端',
      preview: '预览',
      goHome: '回到首页',
      setting: '设置',
      register: '注册',
      registerAccount: '注册账号',
      login: '登录',
    },
    tips: {
      notFound: '啊哦，网页走丢了，正在努力寻找中…',
      copyed: '复制成功',
    },
    table: {
      name: '名称',
      operation: '操作',
      deleteAt: '删除时间',
      yes: '是',
      no: '否',
      emptyText: '暂无数据',
    },
    form: {
      serverUrl: {
        desc: '描述',
        url: '以http://或https://开始',
      },
      user: {
        username: '用户名',
        email: '邮箱',
        password: '密码',
        oldPassword: '旧密码',
        newPassword: '新密码',
        confirmNewPassword: '确认新密码',
      },
    },
    rules: {
      username: {
        required: '请输入用户名',
      },
      email: {
        required: '请输入邮箱',
        correct: '请输入正确的邮箱地址',
      },
      password: {
        required: '请输入密码',
        requiredOld: '请输入旧密码',
        requiredNew: '请输入新密码',
        requiredConfirm: '请输入确认新密码',
        noMatch: '新密码不一致',
        minLength: '密码至少8位',
      },
    },
    user: {
      nav: {
        userSetting: '个人设置',
        modifyPassword: '修改密码',
        logout: '退出登录',
      },
    },
    project: {
      title: '项目',
      list: {
        title: '项目列表',
        tabTitle: '项目',
      },
      form: {
        cover: '项目封面',
        coverColor: '封面颜色',
        coverIcon: '封面图标',
        title: '项目名称',
        desc: '项目描述',
      },
      rules: {
        title: '请输入项目名称',
        titleMinLength: '项目名称不能少于两个字',
        desc: '请输入项目描述信息',
      },
      createModal: {
        title: '创建项目',
        dividerLine: '从以下方式创建',
        blackProject: '空白项目',
        importProject: '导入JSON或YAML数据文件',
        importProjectTip: '支持OpenAPI2.0、3.0',
      },
      setting: {
        title: '项目管理',
        baseInfo: '项目设置',
        serverUrl: 'URL设置',
        globalParam: '全局参数设置',
        responseParam: '公共响应设置',
        export: '导出项目',
        deleteProject: '删除该项目',
        deleteProjectTip: '确定删除该项目吗？<br/> <span style="color:var(--el-color-danger)">项目删除后，相关人员将无法操作该项目。</span>',
        trash: '回收站',
      },
    },
    serverUrl: {
      rules: {
        invalid: '请输入有效的链接地址',
      },
    },
    interface: {
      title: '接口',
      common: {
        aiGenerateInterface: 'AI生成接口',
      },
      tips: {
        unselectedInterface: '请选择要创建的接口',
        allInterfaceCreateFailure: '所有接口创建失败，请重试',
      },
      form: {
        title: '请输入接口标题',
        modalTitle: '请输入您想生成的接口名称',
      },
      table: {
        method: '方式',
        path: '路径',
        desc: '描述',
      },
      popoverMenus: {
        aiGenerateInterface: 'AI生成接口',
        newInterface: '新建接口',
        newGroup: '新建分类',
        confirmDeleteGroup: '确定删除「{0}」分类吗？',
        confirmDeleteInterface: '确定删除「{0}」接口吗？',
        unnamedInterface: '未命名接口',
      },
    },
    schema: {
      title: '模型',
      common: {
        aiGenerateSchema: 'AI生成模型',
      },
      tips: {
        schemaInputTitle: '请输入您想生成的模型名称',
        unselectedInterface: '请选择要创建的接口',
        allInterfaceCreateFailure: '所有接口创建失败，请重试',
      },
      form: {
        title: '请输入模型标题',
        desc: '请输入模型描述',
      },
      popoverMenus: {
        aiGenerateSchema: 'AI生成模型',
        newSchema: '新建模型',
        newGroup: '新建分类',
        confirmDeleteGroup: '确定删除「{0}」分类吗？',
        confirmDeleteSchema: '确定删除「{0}」模型吗？',
        unnamedSchema: 'Unnamed',
      },
    },
    definitionResponse: {
      title: '响应',
      tips: {
        confirmDelete: '确认删除该公共响应吗？',
        unref: '对引用此响应的内容解引用',
      },
      form: {
        title: '请输入响应标题',
        desc: '请输入响应描述',
      },
      popoverMenus: {
        newGroup: '新建分类',
        confirmDeleteDefinitionResponse: '确定删除「{0}」响应吗？',
        unnamedDefinitionResponse: '未命名响应',
      },
    },
    response: {
      title: '响应参数',
      fullname: '响应名称',
      tips: {
        responseExample: '响应示例',
      },
      table: {
        name: '名称',
        code: '状态码',
        desc: '描述',
      },
      model: {
        name: 'Response Name',
      },
      rules: {
        name: '响应名称不能为空',
      },
    },
    request: {
      title: '请求参数',
      tips: {
        noRequestBody: '该请求没有Body体',
        selectFile: '请选择文件',
      },
    },
  },
  editor: {
    common: {
      error: {
        emptyParamName: '参数名不能为空',
        paramNameDuplicate: '参数{0}名称重复',
      },
      tips: {
        confirmDelete: '确认删除{0}吗？',
        delete: '确认删除？',
      },
    },
    node: {
      httpMethod: {
        pathPlaceholder: 'Path, 以"/"开始',
        pathError: '请输入有效的路径',
      },
    },
    table: {
      paramName: '参数名',
      paramType: '参数类型',
      required: '必须',
      defaultValue: '默认值',
      paramDesc: '参数说明',
      paramMock: 'Mock',
      paramExample: '示例值',
      yes: '是',
      no: '否',
      removeBinding: '解除绑定',
      addNode: '添加子节点',
      addParam: '添加参数',
      refModel: '引用模型',
    },
  },
  elementPlusLocale,
}
