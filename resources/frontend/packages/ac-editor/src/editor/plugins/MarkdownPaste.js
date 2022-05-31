import { Plugin } from "prosemirror-state";
import { toggleMark } from "prosemirror-commands";
import Extension from "@/editor/lib/Extension";
import { isUrl } from "@/editor/utils";

export default class MarkdownPaste extends Extension {
  get name() {
    return "markdown-paste";
  }

  get plugins() {
    return [
      new Plugin({
        props: {
          handlePaste: (view, event) => {
            if (view.props.editable && !view.props.editable(view.state)) {
              return false;
            }
            if (!event.clipboardData) return false;

            const text = event.clipboardData.getData("text/plain");
            const html = event.clipboardData.getData("text/html");
            const { state, dispatch } = view;

            // first check if the clipboard contents can be parsed as a single
            // url, this is mainly for allowing pasted urls to become embeds
            if (isUrl(text)) {
              // just paste the link mark directly onto the selected text
              if (!state.selection.empty) {
                toggleMark(this.editor.schema.marks.link, { href: text })(
                  state,
                  dispatch
                );
                return true;
              }

              // well, it's not an embed and there is no text selected – so just
              // go ahead and insert the link directly
              const transaction = view.state.tr
                .insertText(text, state.selection.from, state.selection.to)
                .addMark(
                  state.selection.from,
                  state.selection.to + text.length,
                  state.schema.marks.link.create({ href: text })
                );
              view.dispatch(transaction);
              return true;
            }

            // otherwise, if we have html on the clipboard that looks like it
            // came from Prosemirror then use the default HTML parser behavior
            if (text.length === 0 || (html && html.includes("data-pm-slice"))) {
              return false;
            }

            event.preventDefault();

            // If the users selection is currently in a code block then paste
            // as plain text, ignore all formatting.
              view.dispatch(view.state.tr.insertText(text));
            return true;
          },
        },
      }),
    ];
  }
}
