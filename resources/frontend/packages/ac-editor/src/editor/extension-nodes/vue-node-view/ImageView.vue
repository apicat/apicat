<template>
    <node-view-wrapper :class="imageViewClass">
        <div :class="imageViewBodyClass">
            <img
                :alt="node.attrs.alt"
                :src="src"
                :title="node.attrs.title"
                :width="width"
                :height="height"
                @click="selectImage"
                :data-zoomable="isReadOnly"
                class="image-view__body__image"
            />
            <div class="image-resizer" v-if="!isReadOnly" v-show="isSelected || resizing">
                <span
                    :class="`image-resizer__handler--${direction}`"
                    :key="direction"
                    @mousedown="onMouseDown($event, direction)"
                    class="image-resizer__handler"
                    v-for="direction in resizeDirections"
                />
            </div>
        </div>
    </node-view-wrapper>
</template>

<script>
    import { NodeSelection } from 'prosemirror-state'
    import { ResizeObserver } from '@juggle/resize-observer'
    import { clamp, loadImage } from '../../utils'
    import { NodeViewWrapper } from '../../../components/NodeViewWrapper'
    import { ImageDisplay } from '../../../common/constants'
    import { markRaw } from 'vue'

    const ResizeDirection = {
        TOP_LEFT: 'tl',
        TOP_RIGHT: 'tr',
        BOTTOM_LEFT: 'bl',
        BOTTOM_RIGHT: 'br',
    }

    const MIN_SIZE = 20
    const MAX_SIZE = 100000

    export default {
        components: {
            NodeViewWrapper,
        },
        props: ['node', 'view', 'getPos', 'isSelected', 'isReadOnly', 'updateAttributes'],
        name: 'ImageView',
        computed: {
            src() {
                return this.node.attrs.src
            },

            width() {
                return this.node.attrs.width
            },

            height() {
                return this.node.attrs.height
            },

            imageViewBodyClass() {
                return [
                    'image-view__body',
                    {
                        'image-view__body--focused': this.isSelected,
                        'image-view__body--resizing': this.resizing,
                    },
                ]
            },

            imageViewClass() {
                return ['image-view', `image-view--${this.alignment}`]
            },

            alignment() {
                return this.node.attrs.alignment || ImageDisplay.CENTER
            },
        },
        data() {
            return {
                maxSize: {
                    width: MAX_SIZE,
                    height: MAX_SIZE,
                },

                originalSize: {
                    width: 0,
                    height: 0,
                },

                resizeOb: markRaw(
                    new ResizeObserver(() => {
                        this.getMaxSize()
                    })
                ),

                resizeDirections: [ResizeDirection.TOP_LEFT, ResizeDirection.TOP_RIGHT, ResizeDirection.BOTTOM_LEFT, ResizeDirection.BOTTOM_RIGHT],

                resizing: false,

                resizerState: {
                    x: 0,
                    y: 0,
                    w: 0,
                    h: 0,
                    dir: '',
                },
            }
        },
        methods: {
            getMaxSize() {
                const { width } = getComputedStyle(this.view.dom)
                this.maxSize.width = parseInt(width, 10)
            },

            selectImage() {
                const { state } = this.view
                let { tr } = state
                const selection = NodeSelection.create(state.doc, this.getPos())
                tr = tr.setSelection(selection)
                this.view.dispatch(tr)
            },

            onMouseDown(e, dir) {
                e.preventDefault()
                e.stopPropagation()

                this.resizerState.x = e.clientX
                this.resizerState.y = e.clientY

                const originalWidth = this.originalSize.width
                const aspectRatio = this.originalSize.aspectRatio

                let { width, height } = this.node.attrs
                const maxWidth = this.maxSize.width

                if (width && !height) {
                    width = width > maxWidth ? maxWidth : width
                    height = Math.round(width / aspectRatio)
                } else if (height && !width) {
                    width = Math.round(height * aspectRatio)
                    width = width > maxWidth ? maxWidth : width
                } else if (!width && !height) {
                    width = originalWidth > maxWidth ? maxWidth : originalWidth
                    height = Math.round(width / aspectRatio)
                } else {
                    width = width > maxWidth ? maxWidth : width
                }

                this.resizerState.w = width
                this.resizerState.h = height
                this.resizerState.dir = dir

                this.resizing = true

                this.onEvents()
            },

            onMouseMove(e) {
                e.preventDefault()
                e.stopPropagation()
                if (!this.resizing) return

                const { x, w, dir } = this.resizerState
                const dx = (e.clientX - x) * (/l/.test(dir) ? -1 : 1)
                const width = clamp(w + dx, MIN_SIZE, this.maxSize.width)
                this.updateAttributes({
                    width,
                    height: width / this.originalSize.aspectRatio,
                })
            },

            onMouseUp(e) {
                e.preventDefault()
                e.stopPropagation()
                if (!this.resizing) return

                this.resizing = false

                this.resizerState = {
                    x: 0,
                    y: 0,
                    w: 0,
                    h: 0,
                    dir: '',
                }

                this.offEvents()
                this.selectImage()
            },

            onEvents() {
                document.addEventListener('mousemove', this.onMouseMove, true)
                document.addEventListener('mouseup', this.onMouseUp, true)
            },

            offEvents() {
                document.removeEventListener('mousemove', this.onMouseMove, true)
                document.removeEventListener('mouseup', this.onMouseUp, true)
            },
        },
        created() {
            loadImage(this.src).then((result) => {
                if (!result.complete) {
                    result.width = MIN_SIZE
                    result.height = MIN_SIZE
                }

                this.originalSize = {
                    width: result.width,
                    height: result.height,
                    aspectRatio: result.width / result.height,
                }
            })
        },

        mounted() {
            this.$nextTick(function () {
                this.resizeOb.observe(this.view.dom)
            })
        },

        beforeDestroy() {
            this.resizeOb.disconnect()
        },
    }
</script>
