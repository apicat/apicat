<template>
    <div class="ac-document is-detail">
        <h1 class="ac-document__title">
            {{ doc.title }}
        </h1>

        <div v-if="doc.content" class="ProseMirror readonly" v-html="doc.content" />

        <div v-html="zoomTemplate"></div>
    </div>
</template>

<script>
    import tippy from 'tippy.js'
    import mediumZoom from 'medium-zoom'
    import { toggleClass, getAttr, hasClass, showOrHide } from '@ac/shared'
    import { useHighlight } from '@/hooks/useHighlight'

    function expand(pid, isExpand) {
        document.querySelectorAll('[data-pid="' + pid + '"]').forEach(function (el) {
            let arrow = el.querySelector('i.editor-arrow-right')
            if (arrow && !hasClass(arrow, 'expand')) {
                toggleClass(arrow, 'expand')
            }
            let id = getAttr(el, 'data-id')
            el.style.display = isExpand ? null : 'none'
            id && expand(id, isExpand)
        })
    }

    export default {
        name: 'AcEditorDocument',
        props: ['doc'],
        watch: {
            doc: function () {
                this.init()
            },
        },
        setup() {
            const { initHighlight } = useHighlight()

            return {
                initHighlight,
            }
        },
        data() {
            return {
                zoomTemplate: `<template id="template-zoom-image">
                            <div class="zoom-image-wrapper">
                                <div class="zoom-image-container" data-zoom-container></div>
                            </div>
                          </template>`,
                zoomImageOption: {
                    template: '#template-zoom-image',
                    container: '[data-zoom-container]',
                },
            }
        },
        methods: {
            initTableToggle() {
                document.querySelectorAll('.ac-param-table .editor-arrow-right').forEach(function (el) {
                    el.onclick = function () {
                        expand(getAttr(this, 'data-id'), !hasClass(this, 'expand'))
                        toggleClass(this, 'expand')
                    }
                })

                document.querySelectorAll('div.collapse-title .response_body_title').forEach(function (el) {
                    el.onclick = function () {
                        let h3 = this.parentElement
                        let parent = h3.parentElement
                        let isShow = hasClass(parent, 'close')
                        showOrHide(h3.nextElementSibling, isShow)
                        showOrHide(parent.nextElementSibling, isShow)
                        toggleClass(parent, 'close')
                    }
                })

                document.querySelectorAll('h3.collapse-title >span').forEach(function (el) {
                    el.onclick = function () {
                        let parent = this.parentElement
                        let isShow = hasClass(parent, 'close')
                        showOrHide(parent.nextElementSibling, isShow)
                        toggleClass(parent, 'close')
                    }
                })
            },

            initMediumZoom() {
                mediumZoom('.ProseMirror .image-view img', this.zoomImageOption)
            },

            initTippy() {
                tippy('[data-tippy-content]', {
                    theme: 'light',
                    appendTo: document.querySelector('.ProseMirror'),
                })
            },

            initCodeBlockToClipboard() {
                document.querySelectorAll('.code-block button').forEach((el) => {
                    el.setAttribute('data-text', el.parentElement.querySelector('code').innerText)
                })
            },

            init() {
                this.$nextTick(() => {
                    this.initTableToggle()
                    this.initTippy()
                    this.initMediumZoom()
                    this.initCodeBlockToClipboard()
                    this.initHighlight(document.querySelectorAll('pre code'))
                })
            },
        },
        mounted() {
            this.init()
        },
    }
</script>
