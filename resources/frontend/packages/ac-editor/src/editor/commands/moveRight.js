/*
Copyright 2020 Atlassian Pty Ltd

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 */

// This file is based on the implementation found here:
// https://bitbucket.org/atlassian/design-system-mirror/src/master/editor/editor-core/src/plugins/text-formatting/commands/text-formatting.ts

import isMarkActive from '../utils/isMarkActive'

export default function moveRight() {
    return (state, dispatch) => {
        const { code } = state.schema.marks
        const { empty, $cursor } = state.selection
        if (!empty || !$cursor) {
            return false
        }

        const { storedMarks } = state.tr
        if (code) {
            const insideCode = isMarkActive(code)(state)
            const currentPosHasCode = state.doc.rangeHasMark($cursor.pos, $cursor.pos, code)
            const nextPosHasCode = state.doc.rangeHasMark($cursor.pos, $cursor.pos + 1, code)

            const exitingCode = !currentPosHasCode && !nextPosHasCode && (!storedMarks || !!storedMarks.length)
            const enteringCode = !currentPosHasCode && nextPosHasCode && (!storedMarks || !storedMarks.length)

            // entering code mark (from the left edge): don't move the cursor, just add the mark
            if (!insideCode && enteringCode) {
                dispatch(state.tr.addStoredMark(code.create()))

                return true
            }

            // exiting code mark: don't move the cursor, just remove the mark
            if (insideCode && exitingCode) {
                dispatch(state.tr.removeStoredMark(code))

                return true
            }
        }

        return false
    }
}
