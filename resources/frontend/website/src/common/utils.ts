import { NOT_FOUND } from '@/router/constant'

export const wrapperOrigin = (hasOrigin = true) => (hasOrigin ? window['origin'] : '')

export const goNotFound = () => {
    // location.href = NOT_FOUND.path
}
