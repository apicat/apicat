import { keydownHandler } from 'prosemirror-keymap'
import { TextSelection, NodeSelection, Plugin } from 'prosemirror-state'
import { Decoration, DecorationSet } from 'prosemirror-view'
import { findDomRefAtPos, findPositionOfNodeBefore } from 'prosemirror-utils'
import { GapCursorSelection, JSON_ID, Side } from './GapCursorSelection'
import { toDOM } from './utils/place-gap-cursor'
import { Direction, isBackward, isForward } from './direction'
import { isTextBlockNearPos } from './utils/common'
import { isValidTargetNode } from './utils/is-valid-target-node'
import { atTheBeginningOfDoc, atTheEndOfDoc } from './position'

// :: () → Plugin
// Create a gap cursor plugin. When enabled, this will capture clicks
// near and arrow-key-motion past places that don't have a normally
// selectable position nearby, and create a gap cursor selection for
// them. The cursor is drawn as an element with class
// `ProseMirror-gapcursor`. You can either include
// `style/gapcursor.css` from the package's directory or add your own
// styles to make it visible.
export const gapCursor = function () {
    return new Plugin({
        props: {
            decorations: ({ doc, selection }) => {
                if (selection instanceof GapCursorSelection) {
                    const { $from, side } = selection

                    // render decoration DOM node always to the left of the target node even if selection points to the right
                    // otherwise positioning of the right gap cursor is a nightmare when the target node has a nodeView with vertical margins
                    let position = selection.head
                    const isRightCursor = side === Side.RIGHT
                    if (isRightCursor && $from.nodeBefore) {
                        const nodeBeforeStart = findPositionOfNodeBefore(selection)
                        if (typeof nodeBeforeStart === 'number') {
                            position = nodeBeforeStart
                        }
                    }

                    // const node = isRightCursor ? $from.nodeBefore : $from.nodeAfter;
                    const breakoutMode = ''
                    return DecorationSet.create(doc, [
                        Decoration.widget(position, toDOM, {
                            key: `${JSON_ID}-${side}-${breakoutMode}`,
                            side: breakoutMode ? -1 : 0,
                        }),
                    ])
                }

                return null
            },

            createSelectionBetween(view, $anchor, $head) {
                if ($anchor.pos === $head.pos && GapCursorSelection.valid($head)) {
                    return new GapCursorSelection($head)
                }
                return
            },

            handleKeyDown,
        },
    })
}

export { GapCursorSelection }

const arrow = (dir) => (state, dispatch, view) => {
    const endOfTextblock = view ? view.endOfTextblock.bind(view) : undefined
    const { doc, schema, selection, tr } = state

    let $pos = isBackward(dir) ? selection.$from : selection.$to
    let mustMove = selection.empty

    // start from text selection
    if (selection instanceof TextSelection) {
        // if cursor is in the middle of a text node, do nothing
        if (!endOfTextblock || !endOfTextblock(dir.toString())) {
            return false
        }

        // UP/DOWN jumps to the nearest texblock skipping gapcursor whenever possible
        if (
            (dir === Direction.UP && !atTheBeginningOfDoc(state) && isTextBlockNearPos(doc, schema, $pos, -1)) ||
            (dir === Direction.DOWN && (atTheEndOfDoc(state) || isTextBlockNearPos(doc, schema, $pos, 1)))
        ) {
            return false
        }
        // otherwise resolve previous/next position
        $pos = doc.resolve(isBackward(dir) ? $pos.before() : $pos.after())
        mustMove = false
    }

    if (selection instanceof NodeSelection) {
        if (selection.node.isInline) {
            return false
        }
        if (dir === Direction.UP || dir === Direction.DOWN) {
            // We dont add gap cursor on node selections going up and down
            return false
        }
    }

    // if (!shouldHandleMediaGapCursor(dir, state)) {
    //     return false;
    // }

    // when jumping between block nodes at the same depth, we need to reverse cursor without changing ProseMirror position
    if (
        selection instanceof GapCursorSelection &&
        // next node allow gap cursor position
        isValidTargetNode(isBackward(dir) ? $pos.nodeBefore : $pos.nodeAfter) &&
        // gap cursor changes block node
        ((isBackward(dir) && selection.side === Side.LEFT) || (isForward(dir) && selection.side === Side.RIGHT))
    ) {
        // reverse cursor position
        if (dispatch) {
            dispatch(tr.setSelection(new GapCursorSelection($pos, selection.side === Side.RIGHT ? Side.LEFT : Side.RIGHT)).scrollIntoView())
        }
        return true
    }

    if (view) {
        const domAtPos = view.domAtPos.bind(view)
        const target = findDomRefAtPos($pos.pos, domAtPos)

        if (target && target.textContent === 0) {
            return false
        }
    }

    const nextSelection = GapCursorSelection.findFrom($pos, isBackward(dir) ? -1 : 1, mustMove)

    if (!nextSelection) {
        return false
    }

    if (!isValidTargetNode(isForward(dir) ? nextSelection.$from.nodeBefore : nextSelection.$from.nodeAfter)) {
        // reverse cursor position
        if (dispatch) {
            dispatch(tr.setSelection(new GapCursorSelection(nextSelection.$from, isForward(dir) ? Side.LEFT : Side.RIGHT)).scrollIntoView())
        }
        return true
    }

    if (dispatch) {
        dispatch(tr.setSelection(nextSelection).scrollIntoView())
    }
    return true
}

const handleKeyDown = keydownHandler({
    ArrowLeft: arrow(Direction.LEFT),
    ArrowRight: arrow(Direction.RIGHT),
    ArrowUp: arrow(Direction.UP),
    ArrowDown: arrow(Direction.DOWN),
})
