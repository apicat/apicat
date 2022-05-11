import Node from '../lib/Node'

export default class Text extends Node {
    get name() {
        return 'text'
    }

    get schema() {
        return {
            group: 'inline',
        }
    }
}
