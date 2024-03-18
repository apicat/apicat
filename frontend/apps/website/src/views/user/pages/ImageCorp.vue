<script setup lang="ts">
import 'cropperjs/dist/cropper.css'
import Cropper from 'cropperjs'
import { nextTick, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'

const props = withDefaults(defineProps<IImageCorp>(), {
  handleUpload: () => Promise.resolve(),
})

const { t } = useI18n()

export interface IImageCorp {
  title?: null | string
  handleUpload?: (data: any) => Promise<any>
  cropperConfig?: any
}

// 定义最大范围，以此作为缩放比例
const IMG_CORP = {
  MIN_WIDTH: 32,
  MIN_HEIGHT: 32,
  MIX_WIDTH: 100,
  MIX_HEIGHT: 100,
  MAX_WIDTH: 400,
  MAX_HEIGHT: 500,
  aspectRatio: 1,
} as any

const corpImgContainer = ref()

const state = reactive({
  url: '',
  visible: false,
  corpData: null,
  isLoading: false,
  imgLoading: true,
  corpImgContainerStyle: {
    width: 'auto',
    margin: 'auto',
    height: 'auto',
  },
  imgStyle: {
    width: 'auto',
  },
})

let cropper: any = null

function initCropper(image: HTMLImageElement) {
  cropper = new Cropper(image, {
    initialAspectRatio: 1,
    aspectRatio: 1,
    viewMode: 2,
    rotatable: false,
    scalable: false,
    minCropBoxWidth: IMG_CORP.MIX_WIDTH,
    minCropBoxHeight: IMG_CORP.MIX_HEIGHT,
    minContainerWidth: IMG_CORP.MIX_WIDTH,
    minContainerHeight: IMG_CORP.MIX_HEIGHT,
    ...(props.cropperConfig || {}),
  })
}

function loadImage(url: string) {
  const img = new Image()
  img.style.display = 'none'

  img.onerror = () => {
    reset()

    state.corpImgContainerStyle.width = 'auto'
    state.corpImgContainerStyle.height = 'auto'

    if (corpImgContainer.value?.innerHTML) {
      corpImgContainer.value.innerHTML = `<p class="text-center">${t('app.user.general.imgLoadFail')}</p>`
    }
  }

  img.onload = function () {
    const { width: originWidth, height: originHeight } = this as HTMLImageElement
    if (originWidth < IMG_CORP.MIN_WIDTH || originHeight < IMG_CORP.MIN_HEIGHT) {
      ElMessage.error(t('app.user.general.imgTooSmall', [IMG_CORP.MIN_WIDTH, IMG_CORP.MIN_HEIGHT]))
      return onCloseBtnClick()
    }
    const { width, height } = adjustImage(this as HTMLImageElement)

    state.corpImgContainerStyle.width = width
    state.corpImgContainerStyle.height = height

    img.style.width = width
    img.style.height = height
    img.style.display = 'block'

    state.imgLoading = false

    initCropper(this as HTMLImageElement)
  }

  nextTick().then(() => {
    img.src = url

    if (corpImgContainer.value?.innerHTML) corpImgContainer.value.innerHTML = ''

    corpImgContainer.value.appendChild(img)
  })
}

function adjustImage(img: HTMLImageElement) {
  const ratio = img.width / img.height
  let width = img.width < IMG_CORP.MIX_WIDTH ? IMG_CORP.MIX_WIDTH : Math.min(img.width, IMG_CORP.MAX_WIDTH)
  let height = img.height < IMG_CORP.MIX_HEIGHT ? IMG_CORP.MIX_HEIGHT : Math.min(img.height, IMG_CORP.MAX_HEIGHT)

  // 宽 > 高  横向图片
  if (ratio > 1) height = width / ratio
  // 高 > 宽 竖向图片
  else if (ratio < 1) width = height * ratio
  // 高 = 宽 等比
  else height = width = IMG_CORP.MAX_WIDTH

  width = width < IMG_CORP.MIX_WIDTH ? IMG_CORP.MIX_WIDTH : width
  height = height < IMG_CORP.MIX_HEIGHT ? IMG_CORP.MIX_HEIGHT : height

  return {
    width: `${width}px`,
    height: `${height}px`,
  }
}

function reset() {
  cropper && cropper.destroy()
  cropper = null
  state.imgLoading = false
  state.isLoading = false
}

function onCloseBtnClick() {
  state.isLoading = false
  state.visible = false
}

async function handleSubmit() {
  if (cropper) {
    const cropperData = cropper.getData()
    const data: UserAPI.RequestChangeAvatar = {
      croppedX: cropperData.x.toFixed(),
      croppedY: cropperData.y.toFixed(),
      croppedWidth: cropperData.width.toFixed(),
      croppedHeight: cropperData.height.toFixed(),
    }

    state.isLoading = true
    try {
      await props.handleUpload(data)
    } catch (error) {
      console.error(error)
    }
    // state.isLoading = true
    onCloseBtnClick()
    return
  }

  onCloseBtnClick()
}

function init(url: string | null) {
  if (!url) return
  state.url = url
  state.visible = !!url
  state.imgLoading = true
  state.visible && loadImage(url)
}

watch(
  () => state.visible,
  () => {
    !state.visible && reset()
  },
)

defineExpose({
  show: (url: string) => {
    state.visible = true
    init(url)
  },
})
</script>

<template>
  <el-dialog
    v-model="state.visible"
    :width="600"
    :close-on-click-modal="false"
    :title="title || $t('app.user.general.imgcut')"
    class="show-footer-line"
    destroy-on-close>
    <div
      v-show="!state.imgLoading"
      ref="corpImgContainer"
      v-loading
      class="img-corp"
      :style="state.corpImgContainerStyle" />
    <template #footer>
      <el-button @click="onCloseBtnClick()">
        {{ $t('app.common.cancel') }}
      </el-button>
      <el-button :loading="state.isLoading" type="primary" @click="handleSubmit">
        {{ $t('app.common.confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>
