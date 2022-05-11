import 'codemirror/mode/clike/clike'
import 'codemirror/mode/python/python'
import 'codemirror/mode/dockerfile/dockerfile'
import 'codemirror/mode/go/go'
import 'codemirror/mode/erlang/erlang'
import 'codemirror/mode/fortran/fortran'
import 'codemirror/mode/shell/shell'
import 'codemirror/mode/javascript/javascript'
import 'codemirror/mode/php/php'
import 'codemirror/mode/xml/xml'
import 'codemirror/mode/htmlmixed/htmlmixed'
import 'codemirror/mode/htmlembedded/htmlembedded'
import 'codemirror/mode/css/css'
import 'codemirror/mode/sass/sass'
import 'codemirror/mode/http/http'
import 'codemirror/mode/lua/lua'
import 'codemirror/mode/powershell/powershell'
import 'codemirror/mode/protobuf/protobuf'
import 'codemirror/mode/ruby/ruby'
import 'codemirror/mode/rust/rust'
import 'codemirror/mode/sql/sql'
import 'codemirror/mode/swift/swift'
import 'codemirror/mode/yaml/yaml'

export const languages = {
    none: null,
    json: { name: 'javascript', json: true },

    C: 'clike',
    'C++': 'clike',
    'C#': 'clike',

    css: 'css',
    scss: 'sass',
    sass: 'sass',

    Cython: 'python',
    Dockerfile: 'dockerfile',
    Erlang: 'erlang',
    Fortran: 'fortran',
    Go: 'go',
    'HTML embedded': 'htmlembedded',
    'HTML mixed-mode': 'htmlmixed',
    HTTP: 'http',
    Java: 'clike',
    javascript: 'javascript',
    Kotlin: 'clike',
    Lua: 'lua',
    php: 'php',
    PowerShell: 'powershell',
    ProtoBuf: 'protobuf',
    Python: 'python',
    Ruby: 'ruby',
    Rust: 'rust',
    Scala: 'clike',
    Shell: 'shell',
    SQL: 'sql',
    Swift: 'swift',
    'XML/HTML': 'xml',
    YAML: 'yaml',
}

export default Object.keys(languages)
