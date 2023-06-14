import elementPlusLocale from 'element-plus/lib/locale/lang/en'
export default {
  name: 'English',
  app: {
    common: {
      add: 'Add',
      confirm: 'OK',
      edit: 'Edit',
      save: 'Save',
      export: 'Export',
      import: 'Import',
      cancel: 'Cancel',
      restore: 'Restore',
      reanme: 'Rename',
      copy: 'Copy',
      create: 'Create New',
      generate: 'Generate',
      emptyDataTip: 'No data yet',
      delete: 'Delete',
      deleteTip: 'Delete tips',
      confirmDelete: 'Are you sure you want to delete{msg}?',
      copyAllPath: 'Copy Full URL',
      fetchMockData: 'Get Mock data',
      saving: 'Saving...',
      savedCloud: 'Saved cloud',
      preview: 'Preview',
      goHome: 'Back to homepage',
      goProjectList: 'Return to list of items',
      setting: 'Settings',
      register: 'Sign Up',
      registerAccount: 'Sign Up',
      login: 'Sign in',
    },
    tips: {
      notFound: 'Oops, page not found. We are looking for it...',
      noPermission: 'Unauthorized access…',
      permissionChangeTitle: 'Permissions Change Hint',
      copyed: 'Successfully copied',
    },
    table: {
      name: 'Name',
      operation: 'Operations',
      deleteAt: 'Deleted at',
      yes: 'Y',
      no: 'N',
      emptyText: 'No data yet',
    },
    form: {
      serverUrl: {
        desc: 'Description',
        url: 'Start with http:// or https://.',
      },
      user: {
        username: 'Username',
        email: 'Email',
        password: 'Password',
        oldPassword: 'Old password',
        newPassword: 'New password',
        confirmNewPassword: 'Confirm New Password',
      },
    },
    rules: {
      username: {
        required: 'Please enter a username',
      },
      email: {
        required: 'Please enter email',
        correct: 'Please enter the correct email address',
      },
      password: {
        required: 'Please enter a password',
        requiredOld: 'Please enter an old password',
        requiredNew: 'Please enter a new password',
        requiredConfirm: 'Please enter a confirmation password',
        noMatch: 'New passwords do not match',
        minLength: 'Password must be at least 8 bits',
      },
    },
    user: {
      tips: {
        permissionChange: 'User permissions changed, then reload the action.',
      },
      nav: {
        userSetting: 'Personal Settings',
        modifyPassword: 'Change Password',
        logout: 'Logout',
      },
    },
    project: {
      title: 'Project',
      tips: {
        quitProjectTitle: 'Exit Project',
        quitProject: 'Are you sure you want to exit this project?',
        transferProjectToMember: 'Are you sure you want to transfer the item to that member?',
        targetMemberPermissionError: 'This membership has changed, please try again.',
        permissionChange: 'You have changed your project permissions. You will need to refresh this action.',
        noData: 'No project information available',
      },
      list: {
        title: 'List of items',
        tabTitle: 'Project',
        auth: 'Project Permissions',
      },
      form: {
        cover: 'Project cover',
        coverColor: 'Cover color',
        coverIcon: 'Cover Icon',
        title: 'Project Name',
        desc: 'Project Description',
      },
      rules: {
        title: 'Please enter project name',
        titleMinLength: 'Project name must not be less than two characters',
        desc: 'Please enter project description',
        chooseMember: 'Please select a member',
      },
      createModal: {
        title: 'Create Project',
        dividerLine: 'Create from',
        blackProject: 'Blank items',
        importProject: 'Import JSON or YAML data files',
        importProjectTip: 'Support OpenAPI 2.0 and 3.0',
      },
      setting: {
        title: 'Project management',
        baseInfo: 'Project Settings',
        quitProject: 'Exit Project',
        serverUrl: 'URL Settings',
        globalParam: 'Global Parameter Settings',
        responseParam: 'Public Response Settings',
        member: 'Project members',
        export: 'Export items',
        deleteProject: 'Delete this item',
        deleteProjectTip:
          'Are you sure you want to delete this item?<br/> <span style="color:var(--el-color-danger)">Project deletes will not be able to operate this project.</span>',
        trash: 'Recycle Bin',
      },
      member: {
        title: 'Project members',
        addMember: 'Add Member',
        transferProject: 'Transfer of projects',
        deleteMember: 'Delete member',
        chooseMember: 'Select a member',
        chooseAuth: 'Select permissions',
      },
    },
    member: {
      title: 'Members',
      tips: {
        addMember: 'Add Member',
        editMember: 'Edit Member',
        removeMember: 'Delete member',
        deleteMemberTip: 'Are you sure you want to delete this member?',
      },
      form: {
        title: 'Member list',
        name: 'Name',
        email: 'Email',
        password: 'Password',
        role: 'Team Roles',
        accountStatus: 'Account Status',
        accountStatusNormal: 'Normal',
        accountStatusLock: 'Disabled',
      },
      rules: {},
    },
    serverUrl: {
      rules: {
        invalid: 'Please enter a valid URL',
      },
    },
    interface: {
      title: 'Interface',
      common: {
        aiGenerateInterface: 'AI Generate Interface',
      },
      tips: {
        unselectedInterface: 'Please select the interface to create',
        allInterfaceCreateFailure: 'All interfaces were not created, please try again',
      },
      form: {
        title: 'Please enter the interface title',
        modalTitle: 'Please enter the name of the interface you want to generate.',
      },
      table: {
        method: 'Methodology',
        path: 'Path',
        desc: 'Description',
      },
      popoverMenus: {
        aiGenerateInterface: 'AI Generate Interface',
        newInterface: 'New Interface',
        newGroup: 'New Category',
        confirmDeleteGroup: "Are you sure you want to delete the category `{0}'?",
        confirmDeleteInterface: 'Are you sure you want to delete the interface{0}?',
        unnamedInterface: 'Unnamed interface',
      },
    },
    schema: {
      title: 'Model',
      common: {
        aiGenerateSchema: 'AI Generate Model',
      },
      tips: {
        schemaInputTitle: 'Please enter the model name you want to generate',
        unselectedInterface: 'Please select the interface to create',
        allInterfaceCreateFailure: 'All interfaces were not created, please try again',
      },
      form: {
        title: 'Please enter model title',
        desc: 'Please enter model description',
      },
      popoverMenus: {
        aiGenerateSchema: 'AI Generate Model',
        newSchema: 'New Model',
        newGroup: 'New Category',
        confirmDeleteGroup: "Are you sure you want to delete the category `{0}'?",
        confirmDeleteSchema: "Are you sure you want to delete the model `{0}'?",
        unnamedSchema: 'Unnamed',
      },
    },
    definitionResponse: {
      title: 'Response',
      tips: {
        confirmDelete: 'Are you sure you want to delete this public response?',
        unref: 'Unquote the content of this response',
      },
      form: {
        title: 'Please enter response title',
        desc: 'Please enter response description',
      },
      popoverMenus: {
        newGroup: 'New Category',
        confirmDeleteDefinitionResponse: "Are you sure you want to delete the `{0}' response?",
        unnamedDefinitionResponse: 'Unnamed Response',
      },
    },
    response: {
      title: 'Response Parameters',
      fullname: 'Response Name',
      tips: {
        responseExample: 'Response Example',
      },
      table: {
        name: 'Name',
        code: 'Status Code',
        desc: 'Description',
      },
      model: {
        name: 'Response Name',
      },
      rules: {
        name: 'Response name cannot be empty',
      },
    },
    request: {
      title: 'Request Parameters',
      tips: {
        noRequestBody: 'NoBody for this request',
        selectFile: 'Please select file',
      },
    },
  },
  editor: {
    common: {
      error: {
        emptyParamName: 'Parameter name cannot be empty',
        paramNameDuplicate: 'Parameter {0} duplicated',
      },
      tips: {
        confirmDelete: 'Are you sure you want to delete{0}?',
        delete: 'Are you sure?',
      },
    },
    node: {
      httpMethod: {
        pathPlaceholder: 'Path started with "/"',
        pathError: 'Please enter a valid path',
      },
    },
    table: {
      paramName: 'Parameter Name',
      paramType: 'Parameter Type',
      required: 'Required',
      defaultValue: 'Default value',
      paramDesc: 'Parameter Description',
      paramMock: 'Mock',
      paramExample: 'Example value',
      yes: 'Y',
      no: 'N',
      removeBinding: 'Unbound',
      addNode: 'Add Child',
      addParam: 'Add parameter',
      refModel: 'Quote Model',
    },
  },
  elementPlusLocale,
}
