import { NoMessageAjax } from './Ajax'
import { shortid } from '@ac/shared'
import { IMPORT_EXPORT_STATE } from '@/common/constant'

export const uploadInit = (data = {}) => NoMessageAjax.post('/upload_init', data)
export const uploadPath = (data = {}) => NoMessageAjax.post('/upload_path', data)

export const IMPORT_STATE = IMPORT_EXPORT_STATE

export const dataURLtoFile = (dataUrl, filename) => {
    let arr = dataUrl.split(','),
        mime = arr[0].match(/:(.*?);/)[1],
        str = atob(arr[1]),
        n = str.length,
        u8arr = new Uint8Array(n)
    let suffix = mime.split('/')[1]

    while (n--) {
        u8arr[n] = str.charCodeAt(n)
    }

    return new File([u8arr], (filename || shortid()) + '.' + suffix, { type: mime })
}

const IMAGE_CACHE = {}

export const ImageDisplay = {
    INLINE: 'inline',
    BREAK_TEXT: 'block',
    FLOAT_LEFT: 'left',
    FLOAT_RIGHT: 'right',
}

export const loadImage = (src) =>
    new Promise((resolve, reject) => {
        const result = {
            complete: false,
            width: 0,
            height: 0,
            src,
        }

        if (!src) {
            reject(result)
            return
        }

        if (IMAGE_CACHE[src]) {
            resolve({ ...IMAGE_CACHE[src] })
            return
        }

        const img = new Image()

        img.onload = () => {
            result.width = img.width
            result.height = img.height
            result.complete = true

            IMAGE_CACHE[src] = { ...result }
            resolve(result)
        }

        img.onerror = () => {
            reject(result)
        }

        img.src = src
    })
