import { TextSelection } from 'prosemirror-state'
import isObject from "./isObject"

export default function isTextSelection(value){
    return isObject(value) && value instanceof TextSelection
}
