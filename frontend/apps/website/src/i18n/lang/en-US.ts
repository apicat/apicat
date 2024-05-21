import elementPlusLocale from 'element-plus/es/locale/lang/en'
import { localeEN as components } from '@apicat/components'
import { localeEN as editor } from '@apicat/editor'
import dayjsLocale from '../dayjs/en-US'

const en: any = {
  name: 'English',
  app: {
    description:
      'An efficient API document management tool that fully complies with the OpenAPI specification and incorporates advanced LLM technology. Automatically generate API documents, data models and test cases, greatly improving development efficiency and quality.',
    pageTitles: {
      iterationList: 'Iteration List',
      projectList: 'Project List',
      myStarProjectList: 'Star Project List',
      myProjectList: 'My Project List',
      createProject: 'Create Project',
      team: 'Team',
      createTeam: 'Create Team',
      joinTeam: 'Join Team',
      connectOAuth: 'Connect OAuth',
      schemaHistory: 'Model History',
      collectionHistory: 'Collection History',
      login: 'Login',
      register: 'Register',
      resetPass: 'Reset Password',
      forgetPass: 'Forget Password',
      verificationEmail: 'Verification Email',
      completeInfo: 'Complete Information',
      userSetting: {
        general: 'General',
        email: 'Email',
        github: 'Github',
        password: 'Password',
      },
      systemSetting: {
        service: 'Service',
        oauth: 'OAuth',
        storage: 'Storage',
        email: 'Email',
        model: 'Model',
        users: 'User List',
      },
    },
    common: {
      private: 'Private',
      public: 'Public',
      add: 'Add',
      confirm: 'OK',
      edit: 'Edit',
      save: 'Save',
      export: 'Export',
      import: 'Import',
      cancel: 'Cancel',
      restore: 'Restore',
      rename: 'Rename',
      copy: 'Copy',
      create: 'Create',
      generate: 'Generate',
      emptyDataTip: 'No Datas',
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
      login: 'Sign in',
      loginDivider: 'OR',
      loginGithub: 'Sign in with GitHub',
      register: 'Create an account',
      registerAccount: 'Create an account',
      loginTip: 'Already have an account?',
      generateModelCode: 'Generate Model Code',
      generateCode: 'Generate Code',
      code: 'Code',
      aiPlaceholder: 'Please enter the interface name you want to generate',
      change: 'Change',
      self: 'You',
      resourceLoadErrorTip: 'Resource loading failed, please try again',
      update: 'Update',
      createNew: 'Create New',
      search: 'Search',
    },
    codeGen: {
      model: {
        name: 'Model Name',
      },
      rules: {
        name: 'Model name cannot be empty',
      },
      tips: {
        chooseLanguage: 'Please choose language',
      },
    },
    sign: {
      forgotPass: 'Forgot your password?',
      forgotPassEmail: 'Enter your email',
      forgotPassSend: 'Send',
      resetPass: 'Reset your password',
      resetPassNew: 'New password',
      resetPassRepeat: 'Confirm your new password',
      resetPassSend: 'Reset',
      nocode: 'Invalid code',
      completeInfo: 'Complete your information',
      completeInfoSend: 'Confirm',
    },
    tips: {
      notFound: 'Oops, page not found. We are looking for it...',
      noPermission: 'Unauthorized access…',
      permissionChangeTitle: 'Permissions Change Hint',
      copyed: 'Successfully copied',
      opeartionSuceess: 'Operation Successfully',
      autoSave: 'The system has been automatically saved.',
    },
    table: {
      name: 'Name',
      operation: 'Operations',
      deleteAt: 'Deleted at',
      yes: 'Y',
      no: 'N',
    },
    form: {
      serverUrl: {
        desc: 'Description',
        url: 'Start with http:// or https://',
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
      lang: {
        required: 'Please select a language',
      },
      ava: {
        required: 'Please set an avatar',
      },
      username: {
        required: 'Please enter your username',
        wrongLength: 'username must be 2-64 characters',
      },
      email: {
        required: 'Please enter your email address',
        correct: 'Please enter a valid email address',
      },
      password: {
        required: 'Please enter your password',
        requiredOld: 'Please enter your old password',
        requiredNew: 'Please enter your new password',
        requiredConfirm: 'Please confirm your new password',
        noMatch: 'New passwords do not match',
        minLength: 'Password must be at least 8 characters',
      },
    },
    user: {
      tips: {
        permissionChange: 'User permissions changed, then reload the action.',
      },
      nav: {
        teams: 'Your teams',
        create: 'Create team',
        userSetting: 'Settings',
        modifyPassword: 'Change Password',
        logout: 'Logout',
      },
      general: {
        left_title: 'General',
        title: 'General setting',
        edit: 'Edit',
        username: 'Username',
        lang: 'Language',
        save: 'Save',
        success: 'Update success',
        imgLoadFail: 'Image fail to load, plase try again.',
        imgcut: 'Crop region selection',
        imgTooLarge: 'Image too large. It should be less than 1MB.',
        imgTooSmall: 'The image dimensions should be greater than \'{0}x{1}\'.',
      },
      email: {
        left_title: 'Email',
        title: 'Change email',
        email: 'Email',
        send: 'Send',
        success: 'Change email success.',
      },
      github: {
        left_title: 'GitHub',
        title: 'GitHub connection',
        tip: 'Connect a GitHub account to your ApiCat account to login with GitHub',
        conn: 'Connect',
        disconn: 'Disconnect',
      },
      password: {
        left_title: 'Password',
        title: 'Change password',
        old: 'Old password',
        new: 'New password',
        confirm: 'Confirm new password',
        update: 'Update',
        forgot: 'I forgot my password',
        success: 'Email has been sent.',
        resetSuccess: 'Reset success',
      },
    },
    team: {
      title: 'Team',
      no_current: 'No team selected',
      create: {
        title0: 'You don\'t have a team yet',
        title1: 'Create team',
        name: 'Team name',
        submit: 'Create a new team',
        rule: 'Please enter a team name',
      },
      member: {
        left_title: 'Member',
        disable: {
          btn: 'Disable',
          poptitle: 'Disable team member',
          poptip: 'Do you want to disable this member?',
          ok: 'Disable',
        },
        enable: {
          btn: 'Enable',
          poptitle: 'Enable team member',
          poptip: 'Do you want to enable this member?',
          ok: 'Enable',
        },
        remove: {
          btn: 'Remove',
          poptitle: 'Remove team member',
          poptip: 'Do you want to remove this member from team?',
          ok: 'Remove',
        },
        roles: {
          admin: 'Admin',
          member: 'Member',
          owner: 'Owner',
          Admin: 'Admin',
          Member: 'Member',
          Owner: 'Owner',
        },
        table: {
          name: 'Name',
          email: 'Email',
          role: 'Role',
          status: 'Status',
        },
        status: {
          active: 'Active',
          deactive: 'Deactive',
        },
      },
      invite: {
        left_title: 'Invite',
        title: 'Invite members',
        copy: 'Copy',
        tip1: 'Invite members to join quickly via a public link.',
        tip2: {
          text: 'Anyone who sees the invitation link can join the team. If you want the invitation link to become invalid, please ',
          link: 'reset it',
        },
      },
      setting: {
        left_title: 'Setting',
        title: 'Team setting',
        name: {
          title: 'Team name',
          btn: 'Change',
        },
        transfer: {
          title: 'Transfer ownership',
          holder: 'Select an admin for the team',
          tip: 'You can transfer this team to another administrator, and once that\'s done, you\'ll no longer be the owner of the team, but become a team member.',
          btn: 'Transfer',
          pop: {
            tip: 'After the transfer, you will no longer be the team owner. Please confirm.',
          },
          selectAdmin: 'Please select a team admin',
        },
        rm: {
          title: 'Delete this team',
          tip: 'After deleting the team, all project data within the team will be deleted simultaneously.',
          btn: 'Delete',
          pop: {
            tip: 'After deletion, it will not be possible to access or recover the data through any means. Please confirm.',
          },
        },
        quit: {
          title: 'Leave the team',
          tip: 'After leaving the team, you will no longer have access to view all project data within the team.',
          btn: 'Leave',
          pop: {
            tip: 'Are you leaving the team? please confirm.',
          },
        },
      },
      join: {
        title: 'Welcome to ApiCat',
        tip1: ' invites you to join the ',
        tip2: ' team',
        btn: 'Join the team',
        expire: 'Link expired',
      },
    },
    iter: {
      table: {
        title: 'Iterations',
        project: 'Project: ',
        apis: 'APIs: ',
        created_at: 'Created at: ',
        delete: {
          title: 'Delete iteration',
          tip: 'Are you sure you want to delete this iteration?',
        },
      },
      create: {
        title: 'Create iteration',
        edit_title: 'Edit iteration',
        name: 'Name',
        name_hold: 'iteration name',
        project: 'Project',
        project_hold: 'Select project',
        desc: 'Description',
        plan: {
          title: 'Planning',
          tip: 'Select the APIs involved in this iteration.',
          table: {
            title: 'APIs',
            hold: 'Filter by keywords',
          },
        },
      },
      star: {
        title: 'Stars',
      },
    },
    project: {
      title: 'Project',
      projects: {
        title: 'Projects',
        star: 'Star',
        unstar: 'Unstar',
      },
      stars: {
        title: 'Stars',
        star: 'Star',
        unstar: 'Unstar',
      },
      mypro: {
        title: 'My projects',
      },
      groups: {
        title: 'Project groups',
        grouping: 'Project grouping',
        editGroup: 'Rename group',
        noGroup: 'Not grouped',
        createGroup: 'Create group',
        inputGroupName: 'Enter group name',
        inputGroupNameTip: 'Group name',
        delete: 'Delete',
        rename: 'Rename',
        pop: {
          title: 'Delete this group',
          content: 'Are you sure you want to delete this group?',
        },
      },
      tips: {
        quitProjectTitle: 'Exit Project',
        quitProject: 'Are you sure you want to exit this project?',
        transferProjectToMember: 'Are you sure you want to transfer this project to that member?',
        targetMemberPermissionError: 'The member\'s role has been changed. Please try again.',
        noData: 'No project information available',
      },
      list: {
        title: 'List of Projects',
        tabTitle: 'Project',
        auth: 'Project Permissions',
        emptyDataTip: 'No projects',
      },
      form: {
        visibility: 'Visibility',
        cover: 'Cover',
        coverColor: 'Background Color',
        coverIcon: 'Icon',
        group: 'Group',
        title: 'Name',
        desc: 'Description',
        create: 'Create',
      },
      rules: {
        title: 'Please enter the project name.',
        titleMinLength: 'Project name must not be less than two characters',
        desc: 'Please enter project description',
        chooseMember: 'Please select a member',
      },
      createModal: {
        title: 'Create Project',
        dividerLine: 'Create from',
        blackProject: 'Blank Project',
        importProject: 'Import JSON or YAML data files',
        importProjectTip: 'Support OpenAPI 2.0 and 3.0',
        noInput: 'Please select the file to import.',
        noInputTip: 'Please select a file or change to type "Blank items"',
      },
      setting: {
        title: 'Project settings',
        baseInfo: 'Project Settings',
        quitProject: 'Exit Project',
        deleteProject: 'Delete this item',
        deleteProjectTip:
          'Are you sure you want to delete this item?<br/> <span style="color:var(--el-color-danger)">Project deletes will not be able to operate this project.</span>',
        basic: {
          title: 'General',
          alias: 'Setting',
          update: 'Update',
          detailTitle: 'General',
        },
        member: {
          detailTitle: 'Members',
          title: 'Members',
          smem: 'Select member',
          sper: 'Select permissions',
          add: 'Add member',
          transfer: {
            btn: 'Transfer',
            pop: {
              title: 'Transfer project',
              tip: 'Do you want to transfer the current project to this member?',
            },
          },
          form: {
            memberIDs: 'Please select at least one member',
          },
          poprm: {
            title: 'Remove project member',
            tip: 'Do you want to remove this member from the project?',
            ok: 'Remove',
            trigger: 'Remove',
          },
          table: {
            name: 'Name',
            email: 'Email',
            permissions: 'Permissions',
            status: 'Status',
            rm: 'Remove',
            active: 'Active',
            inactive: 'Inactive',
          },
          auth: {
            manage: 'Manage',
            none: 'None',
            read: 'Read',
            write: 'Write',
          },
        },
        urls: {
          detailTitle: 'URL setting',
          title: 'URLs',
        },
        globalParam: {
          detailTitle: 'Global parameters setting',
          title: 'Global parameters',
          tips: {
            delete: 'Are you sure you want to delete this global parameter?',
            unref: 'Change this parameter to a separate parameter for each API',
          },
        },
        share: {
          detailTitle: 'Share this project',
          title: 'Share',
        },
        export: {
          detailTitle: 'Export project data',
          title: 'Export',
        },
        trash: {
          detailTitle: 'Trash',
          title: 'Trash',
          table: {
            title: 'Title',
            deletedName: 'Deleted by',
            deletedAt: 'Deleted at',
            restore: 'Restore',
          },
        },
        delete: {
          detailTitle: 'Delete project',
          btitle: 'Delete project',
          tip: 'Once you delete the project, there is no going back. Please be certain.',
          btn: 'Delete this project',
          pop: {
            tip: 'Do you want to delete this project?',
            btn: 'Delete',
          },
        },
        quit: {
          detailTitle: 'Quit project',
          btitle: 'Quit project',
          tip: 'Once you quit the project, there is no going back. Please be certain.',
          btn: 'Quit this project',
          success: 'Quit project success',
          pop: {
            tip: 'Do you want to quit this project?',
            btn: 'Quit',
          },
        },
      },
      member: {
        addMember: 'Add Member',
        deleteMember: 'Delete member',
        chooseMember: 'Select a member',
        chooseAuth: 'Select permissions',
      },
      share: {
        title: 'Share this project',
        tip: 'Once you start sharing, anyone with this link and password can view the API content。',
        link: 'Link',
        password: 'Password',
        reset: 'Reset',
        cl: 'Copy link',
        clap: 'Copy link and password',
        copied: 'Copied!',
        noAuthInfo: 'Project authorization unset',
        copylink: 'Document link: ',
        copypass: 'Password: ',
      },
      collection: {
        copy: {
          copied: 'Copied!',
        },
        share: {
          title: 'Share this API',
          tip: 'Once you start sharing, anyone with this link and password can view the API content。',
          link: 'Link',
          password: 'Password',
          reset: 'Reset',
          cl: 'Copy link',
          clap: 'Copy link and password',
          copied: 'Copied!',
          copylink: 'Document link: ',
          copypass: 'Password: ',
        },
        page: {
          sharedoc: 'Share document',
          exportdoc: 'Export document',
          history: 'History',
          edit: {
            titleNull: 'Title can not be none',
          },
        },
        history: {
          diff: {
            title: 'Diff',
            holder: 'Compare documents',
          },
          restore: {
            title: 'Rollback this version',
            poptitle: 'Rollback this version',
            popcontent: 'Are you sure you want to rollback this version of content',
            success: 'Rollback complete',
          },
        },
        export: {
          title: 'Export document',
        },
        ai: {
          title: 'Create API with AI',
          aiPromptPlaceholder: 'Input the description information of the API you want.',
        },
        test: {
          title: 'Test cases',
          table: {
            name: 'Test case',
            time: 'Created at',
            empty: 'No test case',
          },
          btn: {
            generate: 'Generate test cases',
            continueGen: 'Continue generating',
            ai: 'Ask AI to generate',
            regen: 'Regenerate all test cases',
            loading: 'Generating',
          },
          holder: 'Tell the AI what test case you want. For example: Edge cases that may not have been anticipated',
          detail: {
            regenerate: 'Regenerate',
            fail: 'Regenerate fatal',
            del: {
              title: 'Delete test case',
              tip: 'Do you want to delete this test case',
              confirm: 'Delete',
            },
            holder: 'Optional. Tell the AI which parts of the content need improvement',
          },
        },
      },
      infoHeader: {
        private: 'Private',
      },
      history: {
        diff: {
          current: 'Current',
        },
      },
    },
    iteration: {
      title: 'Iteration',
      left: {
        allIter: 'Iterations',
        createIter: 'Create new',
        star: 'Star',
      },
      table: {
        project: 'Project',
        apiCount: 'APIs',
        createAt: 'Create at',
      },
      form: {
        create: 'Create',
        edit: 'Edit iteration',
        name: 'Name',
        inpIterName: 'Iteration name',
        inpIterNameTip: 'Please enter iteration name',
        project: 'Project',
        selectProject: 'Select the project',
        selectProjectTip: 'Please select project',
        desc: 'Description',
        inpDesc: 'Iteration name',
        descTip: 'Please enter iteration description',
        plan: 'Planning',
        planTip: 'Select the APIs involved in this iteration.',
      },
      list: {
        emptyDataTip: 'No iterations',
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
        title: 'Team members',
        name: 'Name',
        email: 'Email',
        password: 'Password',
        accountStatus: 'Account Status',
        accountStatusNormal: 'Normal',
        accountStatusLock: 'Disabled',
      },
    },
    serverUrl: {
      rules: {
        invalid: 'Please enter a valid URL',
        duplicated: 'URL duplicated',
      },
    },
    interface: {
      title: 'APIs',
      common: {
        aiGenerateInterface: 'AI Generate API',
      },
      tips: {
        allInterfaceCreateFailure: 'All interfaces were not created, please try again',
      },
      form: {
        title: 'Please enter the interface title',
        modalTitle: 'Please enter the name of the interface you want to generate.',
      },
      table: {
        desc: 'Description',
      },
      popoverMenus: {
        aiGenerateInterface: 'Create an API with AI',
        newInterface: 'Create an API',
        newGroup: 'Create a category',
        confirmDeleteGroup: 'Are you sure you want to delete the category \'{0}\'?',
        confirmDeleteInterface: 'Are you sure you want to delete the API \'{0}\'?',
      },
      unnamedInterface: 'Unnamed',
      unnamedCategory: 'UnnamedCategory',
    },
    schema: {
      title: 'Models',
      common: {
        aiGenerateSchema: 'Create a model with AI',
      },
      tips: {
        schemaInputTitle: 'Please enter the model name you want to generate',
        unref: 'Unlink content referencing this model',
      },
      form: {
        title: 'Please enter model title',
        desc: 'Please enter model description',
      },
      popoverMenus: {
        aiGenerateSchema: 'Create a model with AI',
        newSchema: 'Create a model',
        newGroup: 'Create a category',
        confirmDeleteGroup: 'Are you sure you want to delete category \'{0}\'?',
        confirmDeleteSchema: 'Are you sure you want to delete model \'{0}\'?',
      },
      history: {
        title: 'History',
        diff: {
          title: 'Diff',
          holder: 'Compare documents',
        },
        restore: {
          title: 'Rollback this version',
          poptitle: 'Rollback this version',
          popcontent: 'Are you sure you want to rollback this version of content',
          success: 'Rollback complete',
        },
      },
      ai: {
        title: 'Create Model with AI',
        aiPromptPlaceholder: 'Input the description information of the model you want.',
      },
      page: {
        edit: {
          titleNull: 'Title can not be none',
        },
      },
      unnamedSchema: 'UnnamedModel',
      unnamedCategory: 'UnnamedCategory',
    },
    definitionResponse: {
      title: 'Responses',
      tips: {
        confirmDelete: 'Are you sure you want to delete this public response?',
      },
      form: {
        title: 'Please enter response title',
        desc: 'Please enter response description',
      },
      popoverMenus: {
        newDefinitionResponse: 'Create a response',
        newGroup: 'Create a category',
        confirmDeleteDefinitionResponse: 'Are you sure you want to delete the \'{0}\' response?',
        unnamedDefinitionResponse: 'Unnamed Response',
      },
      page: {
        edit: {
          titleNull: 'Title can not be none',
        },
      },
      unnamedResponse: 'UnnamedResponse',
      unnamedCategory: 'UnnamedCategory',
    },
    response: {
      title: 'Response Parameters',
      fullname: 'Response Name',
      tips: {
        responseExample: 'Response Example',
        unref: 'Unquote the content of this response',
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
    acSelect: {
      holder: 'Please choose one',
    },
    historyLayout: {
      record: 'History',
      current: 'Current',
    },
    memberAuth: {
      manage: 'Manage',
      write: 'Write',
      read: 'Read',
      none: 'None',
    },
    save: {
      saving: 'Saving...',
      error: 'Save error!',
      saved: 'Saved in the cloud',
    },
    verifyShare: {
      holder: 'Access password',
      btn: 'Verify',
    },
    verifyEmail: {
      p1: 'Your email has been verified, redirect to main page in ',
      p2: ' seconds, if nothing happens, ',
      p3: 'click here',
    },
    system: {
      title: 'Setting',
      service: {
        title: 'Service',
        left_title: 'Service',
        appname: 'App name',
        appurl: 'App url',
        mockurl: 'Mock url',
        rules: {
          appname: 'App name is required',
          appurl: 'Please enter a valid url',
          mockurl: 'Please enter a valid url',
        },
      },
      oauth: {
        title: 'OAuth',
        left_title: 'OAuth',
        github: {
          id: 'Client ID',
          secret: 'Client Secret',
        },
        rules: {
          clientID: 'Client ID is required',
          clientSecret: 'Client Secret is required',
        },
      },
      storage: {
        title: 'Storage',
        left_title: 'Storage',
        local: {
          title: 'Local disk',
          path: 'Path',
          rules: {
            path: 'Path is required',
          },
        },
        cf: {
          title: 'CLOUDFLARE',
          accountID: 'Account ID',
          accessKeyID: 'Access key ID',
          accessKeySecret: 'Access key secret',
          bucketName: 'Bucket name',
          bucketUrl: 'Bucket url',
          rules: {
            accountID: 'Account ID is required',
            accessKeyID: 'Access key ID is required',
            accessKeySecret: 'Access key secret is required',
            bucketName: 'Bucket name is required',
            bucketUrl: 'Please enter a valid url',
          },
        },
        qiniu: {
          title: 'QiNiuYun',
          accessKey: 'Access key',
          secretKey: 'Secret key',
          bucketName: 'Bucket name',
          bucketUrl: 'Bucket url',
          rules: {
            accessKey: 'Access key is required',
            secretKey: 'Secret key is required',
            bucketName: 'Bucket name is required',
            bucketUrl: 'Please enter a valid url',
          },
        },
      },
      email: {
        title: 'Email',
        left_title: 'Email',
        smtp: {
          title: 'SMTP',
          host: 'Host',
          user: 'User',
          address: 'Address',
          pw: 'Password',
          rules: {
            host: 'Host is required',
            address: 'Address is required',
            pw: 'Password is required',
          },
        },
        sendcloud: {
          title: 'SendCloud',
          apiUser: 'API user',
          apiKey: 'API key',
          fromEmail: 'From email',
          fromName: 'From name',
          rules: {
            apiUser: 'API user is required',
            apiKey: 'API Key is required',
            fromName: 'From Name is required',
            fromEmail: 'From Email is required',
          },
        },
      },
      model: {
        title: 'Model',
        left_title: 'Model',
        openai: {
          title: 'OpenAI',
          apiKey: 'API key',
          organizationID: 'Organization ID',
          apiBase: 'API base',
          llmName: 'Model name',
          rules: {
            apiKey: 'API Key is required',
            apiBase: 'API Base is required',
            organizationID: 'OrganizationID is required',
            llmName: 'Model name is required',
          },
        },
        azure: {
          title: 'Azure',
          apiKey: 'API key',
          endpoint: 'Endpoint',
          llmName: 'Model name',
          rules: {
            apiKey: 'API Key is required',
            endpoint: 'Endpoint is required',
            llmName: 'Model name is required',
          },
        },
      },
      users: {
        title: 'Users',
        left_title: 'Users',
        searchKeywordPlaceholder: 'Search by username',
        updatePasswordTitle: 'Update password',
        removeUserTitle: 'Remove user',
        removeUserTip: 'Do you want to remove this user?',
      },
    },
    app: {
      system: {
        model: {
          openai: {
            ruls: {
              apiKey: 'API Key is required',
              apiBase: 'API Base is required',
              organizationID: 'OrganizationID is required',
            },
          },
          azure: {
            ruls: {
              apiKey: 'API Key is required',
              endpoint: 'Endpoint is required',
            },
          },
        },
        email: {
          sendcloud: {
            rules: {
              path: 'Path is required',
              apiUser: 'API user is required',
              apiKey: 'API Key is required',
              fromName: 'From Name is required',
              fromEmail: 'From Email is required',
            },
          },
          smtp: {
            title: 'SMTP',
            rules: {
              host: 'Host is required',
            },
            host: {
              label: 'Host',
            },
          },
        },
      },
    },
  },
  editor,
  components,
  elementPlusLocale,
  dayjsLocale,
}

export { en as default }
