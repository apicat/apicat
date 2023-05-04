import elementPlusLocale from 'element-plus/lib/locale/lang/en'

export default {
  name: 'English',
  app: {
    common: {
      add: 'Add',
      confirm: 'Confirm',
      edit: 'Edit',
      save: 'Save',
      export: 'Export',
      import: 'Import',
      cancel: 'Cancel',
      restore: 'Restore',
      reanme: 'Rename',
      copy: 'Copy',
      create: 'Create',
      generate: 'Generate',
      emptyDataTip: 'No data available',
      delete: 'Delete',
      deleteTip: 'Delete tip',
      confirmDelete: 'Are you sure you want to delete {msg}?',
      copyAllPath: 'Copy full URL',
      saving: 'Saving...',
      savedCloud: 'Saved to Cloud',
      preview: 'Preview',
      goHome: 'Go back to Home',
    },
    tips: {
      notFound: 'Oops, page not found. We are looking for it...',
      copyed: 'Copied successfully',
    },
    table: {
      name: 'Name',
      operation: 'Operation',
      deleteAt: 'Delete Time',
      yes: 'Y',
      no: 'N',
    },
    form: {
      serverUrl: {
        desc: 'Description',
        url: 'Starts with http:// or https://',
      },
    },
    project: {
      list: {
        title: 'Project List',
        tabTitle: 'Project',
      },
      form: {
        title: 'Project Name',
        desc: 'Project Description',
      },
      rules: {
        title: 'Please enter project name',
        titleMinLength: 'Project name cannot be less than two characters',
        desc: 'Please enter project description information',
      },
      createModal: {
        title: 'Create Project',
        dividerLine: 'Create from the following methods',
        blackProject: 'Blank Project',
        importProject: 'Import JSON or YAML data file',
        importProjectTip: 'Supports OpenAPI 2.0, 3.0',
      },
      setting: {
        title: 'Project Management',
        baseInfo: 'Project Setting',
        serverUrl: 'URL Setting',
        globalParam: 'Global Parameter Setting',
        responseParam: 'Public Response Setting',
        export: 'Export Project',
        deleteProject: 'Delete Project',
        deleteProjectTip:
          'Are you sure you want to delete this project? <br/> <span style="color:var(--el-color-danger)">After deleting the project, the relevant personnel will not be able to operate this project.</span>',
        trash: 'Recycle Bin',
      },
    },
    serverUrl: {
      rules: {
        invalid: 'Please enter a valid link address',
      },
    },
    interface: {
      title: 'Interface',
      common: {
        aiGenerateInterface: 'AI Generate Interface',
      },
      tips: {
        unselectedInterface: 'Please select the interface to be created',
        allInterfaceCreateFailure: 'All interface creation failed, please try again',
      },
      form: {
        title: 'Please enter interface title',
        modalTitle: 'Please enter the name of the interface you want to generate',
      },
      table: {
        method: 'Method',
        path: 'Path',
        desc: 'Description',
      },
      popoverMenus: {
        aiGenerateInterface: 'AI Generate Interface',
        newInterface: 'New Interface',
        newGroup: 'New Group',
        confirmDeleteGroup: 'Are you sure you want to delete "{0}" group?',
        confirmDeleteInterface: 'Are you sure you want to delete the "{0}" interface?',
        unnamedInterface: 'Unnamed Interface',
      },
    },
    schema: {
      title: 'Schema',
      common: {
        aiGenerateSchema: 'AI Generate Schema',
      },
      tips: {
        schemaInputTitle: 'Please enter the name of the schema you want to generate',
        unselectedInterface: 'Please select the interface to be created',
        allInterfaceCreateFailure: 'All interface creation failed, please try again',
      },
      form: {
        title: 'Please enter schema title',
        desc: 'Please enter schema description',
      },
      popoverMenus: {
        aiGenerateSchema: 'AI Generate Schema',
        newSchema: 'New Schema',
        newGroup: 'New Group',
        confirmDeleteGroup: 'Are you sure you want to delete "{0}" group?',
        confirmDeleteSchema: 'Are you sure you want to delete the "{0}" schema?',
        unnamedSchema: 'Unnamed',
      },
    },
    publicResponse: {
      title: 'Public Response',
    },
    response: {
      title: 'Response Parameter',
      fullname: 'Response Name',
      tips: {
        confirmDelete: 'Are you sure you want to delete this public response?',
        responseExample: 'Response example',
      },
      table: {
        name: 'Name',
        code: 'Status Code',
        desc: 'Description',
      },
      model: {
        description: 'Public Response',
      },
      rules: {
        name: 'Response name cannot be empty',
      },
    },
    request: {
      title: 'Request Parameter',
      tips: {
        noRequestBody: 'No request body',
        selectFile: 'Please select a file',
      },
    },
  },
  editor: {
    common: {
      error: {
        emptyParamName: 'Parameter name cannot be empty',
        paramNameDuplicate: 'Duplicate parameter name for {0}',
      },
      tips: {
        confirmDelete: 'Are you sure you want to delete {0}?',
        delete: 'Delete it?',
      },
    },
    node: {
      httpMethod: {
        pathPlaceholder: 'Path, starts with "/"',
        pathError: 'Please enter a valid path',
      },
    },
    table: {
      paramName: 'Parameter Name',
      paramType: 'Parameter Type',
      required: 'Required',
      defaultValue: 'Default Value',
      paramDesc: 'Parameter Description',
      paramExample: 'Example Value',
      yes: 'Y',
      no: 'N',
      removeBinding: 'Remove Binding',
      addNode: 'Add Child Node',
      addParam: 'Add Parameter',
      refModel: 'Reference Model',
    },
  },
  elementPlusLocale,
}
