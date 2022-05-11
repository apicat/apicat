export default class Extension {
  constructor(options = {}) {
    this.options = {
      ...this.defaultOptions,
      ...options,
    };
  }

  get allowCreate() {
    return true;
  }

  bindEditor(editor = null) {
    this.editor = editor;
  }

  get type() {
    return "extension";
  }

  get name() {
    return "";
  }

  get plugins() {
    return [];
  }

  keys() {
    return {};
  }

  inputRules() {
    return [];
  }

  commands() {
    return () => () => false;
  }

  get defaultOptions() {
    return {};
  }
}
