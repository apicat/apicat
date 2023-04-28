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
      emptyDataTip: '暂无数据',
      delete: '删除',
      deleteTip: '删除提示',
      confirmDelete: '确认删除{msg}吗？',
      copyAllPath: '复制完整URL',
    },
    table: {
      paramName: '参数名称',
      paramType: '参数类型',
      required: '必须',
      defaultValue: '默认值',
      paramDesc: '参数说明',
      deleteResponseConfirm: '确认删除该公共响应吗？',
    },
    form: {
      serverUrl: {
        desc: '描述',
        url: '以http://或https://开始',
      },
    },
    project: {
      list: {
        title: '项目列表',
        tabTitle: '项目',
      },
      form: {
        title: '项目名称',
        desc: '项目描述',
      },
      rules: {
        title: '请输入项目名称',
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
  },
  // 编辑器
  editor: {
    common: {
      error: {
        paramNameDuplicate: '参数{0}名称重复',
      },
      tips: {
        confirmDelete: '确认删除{0}吗？',
      },
    },
    node: {
      httpMethod: {
        pathPlaceholder: 'Path, 以"/"开始',
        pathError: '请输入有效的路径',
      },
    },
  },
  elementPlusLocale,
}
