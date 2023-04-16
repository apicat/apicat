import elementPlusLocale from 'element-plus/lib/locale/lang/en'

export default {
  name: 'English',
  app: {
    common: {
      add: 'Add',
      confirm: 'Confirm',
      save: 'Save',
      export: 'Export',
      import: 'Import',
      cancel: 'Cancel',
      restore: 'Restore',
      emptyDataTip: 'No data available',
      delete: 'Delete',
      deleteTip: 'Delete confirmation message',
      confirmDelete: 'Are you sure you want to delete {message}?',
      copyAllPath: 'Copy full URL',
    },
    table: {
      paramName: 'Parameter name',
      paramType: 'Parameter type',
      required: 'Required',
      defaultValue: 'Default value',
      paramDesc: 'Parameter description',
      deleteResponseConfirm: 'Are you sure you want to delete this public response?',
    },
    form: {
      serverUrl: {
        desc: 'Description',
        url: 'Start with http:// or https://',
      },
    },
    project: {
      list: {
        title: 'Project list',
        tabTitle: 'Project',
      },
      form: {
        title: 'Project name',
        desc: 'Project description',
      },
      rules: {
        title: 'Enter project name',
        desc: 'Enter project description',
      },
      createModal: {
        title: 'Create project',
        dividerLine: 'From following ways',
        blackProject: 'Empty project',
        importProject: 'Import JSON data file',
        importProjectTip: 'Supports OpenAPI 2.0 and 3.0',
      },
      setting: {
        title: 'Project management',
        baseInfo: 'Project settings',
        serverUrl: 'URL settings',
        requestParam: 'Public parameter settings',
        responseParam: 'Public response settings',
        export: 'Export project',
        trash: 'Recycle bin',
      },
    },
  },
  editor: {},
  elementPlusLocale,
}
