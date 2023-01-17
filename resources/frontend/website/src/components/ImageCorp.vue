<template>
    <el-dialog :width="600" :close-on-click-modal="false" v-model="state.visible" :title="title || '选择裁剪区域'" class="show-footer-line">
        <div v-loading v-show="!state.imgLoading" ref="corpImgContainer" class="img-corp" :style="state.corpImgContainerStyle"></div>
        <template #footer>
            <el-button @click="onCloseBtnClick()">取消</el-button>
            <el-button :loading="state.isLoading" type="primary" @click="handleSubmit">确 定</el-button>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
    import 'cropperjs/dist/cropper.css'
    import Cropper from 'cropperjs'
    import { settingAvatar } from '@/api/user'
    import { nextTick, onMounted, reactive, ref, watch } from 'vue'
    import { ElMessage } from 'element-plus'
    //定义最大范围，以此作为缩放比例
    const IMG_CORP = {
        MIX_WIDTH: 100,
        MIX_HEIGHT: 100,
        MAX_WIDTH: 400,
        MAX_HEIGHT: 500,
        aspectRatio: 1,
    } as any

    interface IImageCorp {
        imgUrl: null | string
        file: null | File
        title?: null | string
        handleUpload?: (data: any) => Promise<any>
        cropperConfig?: any
    }

    const props = withDefaults(defineProps<IImageCorp>(), {
        imgUrl: null,
        file: null,
    })

    const emit = defineEmits(['on-ok'])

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

    const initCropper = (image: HTMLImageElement) => {
        cropper = new Cropper(image, {
            initialAspectRatio: 1,
            aspectRatio: 1,
            viewMode: 3,
            rotatable: false,
            scalable: false,
            minCropBoxWidth: IMG_CORP.MIX_WIDTH,
            minCropBoxHeight: IMG_CORP.MIX_HEIGHT,
            minContainerWidth: IMG_CORP.MIX_WIDTH,
            minContainerHeight: IMG_CORP.MIX_HEIGHT,
            ...(props.cropperConfig || {}),
        })
    }

    const loadImage = (url: string) => {
        let img = new Image()
        img.style.display = 'none'

        img.onerror = (e) => {
            reset()

            state.corpImgContainerStyle.width = 'auto'
            state.corpImgContainerStyle.height = 'auto'

            if (corpImgContainer.value?.innerHTML) {
                corpImgContainer.value.innerHTML = '<p class="text-center">图片加载失败，请重试。</p>'
            }
        }

        img.onload = function () {
            let { width, height } = adjustImage(this as HTMLImageElement)

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

            if (corpImgContainer.value?.innerHTML) {
                corpImgContainer.value.innerHTML = ''
            }
            corpImgContainer.value.appendChild(img)
        })
    }

    const adjustImage = (img: HTMLImageElement) => {
        let ratio = img.width / img.height
        let width = img.width < IMG_CORP.MIX_WIDTH ? IMG_CORP.MIX_WIDTH : Math.min(img.width, IMG_CORP.MAX_WIDTH)
        let height = img.height < IMG_CORP.MIX_HEIGHT ? IMG_CORP.MIX_HEIGHT : Math.min(img.height, IMG_CORP.MAX_HEIGHT)

        //宽 > 高  横向图片
        if (ratio > 1) {
            height = width / ratio
        }

        //高 > 宽 竖向图片
        else if (ratio < 1) {
            width = height * ratio
        }

        //高 = 宽 等比
        else {
            // height = width = IMG_CORP.WIDTH;
        }

        width = width < IMG_CORP.MIX_WIDTH ? IMG_CORP.MIX_WIDTH : width
        height = height < IMG_CORP.MIX_HEIGHT ? IMG_CORP.MIX_HEIGHT : height

        return {
            width: width + 'px',
            height: height + 'px',
        }
    }

    const reset = () => {
        cropper && cropper.destroy()
        cropper = null
        state.imgLoading = false
    }

    const onCloseBtnClick = () => {
        state.visible = false
    }

    const handleSubmit = async () => {
        if (cropper) {
            var cropperData = cropper.getData()
            var formData = new FormData()
            formData.append('cropped_x', cropperData.x.toFixed())
            formData.append('cropped_y', cropperData.y.toFixed())
            formData.append('cropped_width', cropperData.width.toFixed())
            formData.append('cropped_height', cropperData.height.toFixed())

            state.isLoading = true

            if (props.handleUpload) {
                try {
                    const { data } = (await props.handleUpload(formData)) || {}
                    emit('on-ok', data)
                    state.visible = false
                } catch (error) {
                    //
                } finally {
                    state.isLoading = false
                }
                return
            }

            formData.append('avatar', props.file as any)

            try {
                const { data } = (await settingAvatar(formData)) || {}
                ElMessage.success('头像更新成功!')
                emit('on-ok', data)
                state.visible = false
            } catch (error) {
                //
            } finally {
                state.isLoading = false
            }
        } else {
            onCloseBtnClick()
        }
    }

    const init = (url: string | null) => {
        if (!url) return
        state.url = url
        state.visible = !!url
        state.imgLoading = true
        state.visible && loadImage(url)
    }

    watch(
        () => props.imgUrl,
        () => init(props.imgUrl)
    )

    watch(
        () => state.visible,
        () => {
            !state.visible && reset()
        }
    )

    onMounted(() => {
        init(props.imgUrl)
    })
</script>
