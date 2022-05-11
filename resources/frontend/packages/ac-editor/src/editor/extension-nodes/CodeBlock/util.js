import { javascript } from "@codemirror/lang-javascript";
import { java } from "@codemirror/lang-java";
import { rust } from "@codemirror/lang-rust";
import { sql } from "@codemirror/lang-sql";
import { json } from "@codemirror/lang-json";
import { python } from "@codemirror/lang-python";
import { html } from "@codemirror/lang-html";
import { css } from "@codemirror/lang-css";
import { cpp } from "@codemirror/lang-cpp";
import { markdown } from "@codemirror/lang-markdown";
import { xml } from "@codemirror/lang-xml";
import { StreamLanguage } from "@codemirror/stream-parser";
import { haskell } from "@codemirror/legacy-modes/mode/haskell";
import { clojure } from "@codemirror/legacy-modes/mode/clojure";
import { erlang } from "@codemirror/legacy-modes/mode/erlang";
import { groovy } from "@codemirror/legacy-modes/mode/groovy";
import { ruby } from "@codemirror/legacy-modes/mode/ruby";
import { shell } from "@codemirror/legacy-modes/mode/shell";
import { yaml } from "@codemirror/legacy-modes/mode/yaml";
import { go } from "@codemirror/legacy-modes/mode/go";

export const cleanLang = (lang) =>
  lang === "js"
    ? "javascript"
    : lang === "ts"
    ? "typescript"
    : lang === "cplusplus"
    ? "cpp"
    : lang === "scss"
    ? "css"
    : lang === "sass"
    ? "css"
    : lang === "less"
    ? "css"
    : lang === "c++"
    ? "cpp"
    : lang === "yml"
    ? "yaml"
    : lang === "shell"
    ? "bash"
    : lang;

export const getLangExtension = (lang) =>
  lang === "javascript"
    ? javascript()
    : lang === "typescript"
    ? javascript({ typescript: true })
    : lang === "java" || lang === "kotlin"
    ? java()
    : lang === "rust"
    ? rust()
    : lang === "sql"
    ? sql()
    : lang === "json"
    ? json()
    : lang === "python"
    ? python()
    : lang === "html"
    ? html()
    : lang === "css"
    ? css()
    : lang === "cpp"
    ? cpp()
    : lang === "markdown"
    ? markdown()
    : lang === "xml"
    ? xml()
    : lang === "haskell"
    ? StreamLanguage.define(haskell)
    : lang === "clojure"
    ? StreamLanguage.define(clojure)
    : lang === "erlang"
    ? StreamLanguage.define(erlang)
    : lang === "groovy"
    ? StreamLanguage.define(groovy)
    : lang === "ruby"
    ? StreamLanguage.define(ruby)
    : lang === "bash"
    ? StreamLanguage.define(shell)
    : lang === "yaml"
    ? StreamLanguage.define(yaml)
    : lang === "go"
    ? StreamLanguage.define(go)
    : markdown();

export function computeChange(oldVal, newVal) {
  if (oldVal === newVal) return null;
  let start = 0;
  let oldEnd = oldVal.length;
  let newEnd = newVal.length;

  while (
    start < oldEnd &&
    oldVal.charCodeAt(start) === newVal.charCodeAt(start)
  ) {
    ++start;
  }

  while (
    oldEnd > start &&
    newEnd > start &&
    oldVal.charCodeAt(oldEnd - 1) === newVal.charCodeAt(newEnd - 1)
  ) {
    oldEnd--;
    newEnd--;
  }

  return {
    from: start,
    to: oldEnd,
    text: newVal.slice(start, newEnd),
  };
}
