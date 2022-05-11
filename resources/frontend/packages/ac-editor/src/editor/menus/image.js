import { isNodeActive } from '../utils'
import { ImageDisplay } from '../../common/constants'

export default function imageMenuItems(state, dictionary) {
    const { schema } = state
    const isLeftAligned = isNodeActive(schema.nodes.image, {
        alignment: ImageDisplay.LEFT,
    })
    const isRightAligned = isNodeActive(schema.nodes.image, {
        alignment: ImageDisplay.RIGHT,
    })

    return [
        {
            name: 'alignLeft',
            tooltip: dictionary.alignLeft,
            icon: 'editor-text-align-left',
            visible: true,
            active: isLeftAligned,
        },
        {
            name: 'alignCenter',
            tooltip: dictionary.alignCenter,
            icon: 'editor-text-align-center',
            visible: true,
            active: (state) => isNodeActive(schema.nodes.image)(state) && !isLeftAligned(state) && !isRightAligned(state),
        },
        {
            name: 'alignRight',
            tooltip: dictionary.alignRight,
            icon: 'editor-text-align-right',
            visible: true,
            active: isRightAligned,
        },
        {
            name: 'separator',
            visible: true,
        },
        {
            name: 'deleteImage',
            tooltip: dictionary.deleteImage,
            icon: 'editor-trash',
            visible: true,
            active: () => false,
        },
    ]
}
