import { wrapIn, lift } from 'prosemirror-commands'
import { nodeIsActive } from '../utils'

export default function (type, attrs = {}) {
    return (state, dispatch, view) => {
        const isActive = nodeIsActive(state, type, attrs)

        if (isActive) {
            return lift(state, dispatch)
        }

        return wrapIn(type, attrs)(state, dispatch, view)
    }
}
