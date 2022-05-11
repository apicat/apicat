<template>
    <div class="color-panel">
        <p class="group-title"><span class="group-title-name">字体颜色</span></p>
        <div class="panel-menu-group">
            <button v-for="item in fontColors" :key="item" :class="fontColorClass(item)" @click="onFontColorClick($event, item)">
                <svg width="22" height="22" viewBox="0 0 22 22" class="panel-btn-icon" :style="{ fill: item }">
                    <path d="M10.28 5.5L6 16.5h1.34l1.16-3.08h4.99l1.16 3.08H16l-4.28-11h-1.44zm-1.39 6.86l2.09-5.5h.06l2.05 5.5H8.9z"></path>
                </svg>
            </button>
        </div>
        <p class="group-title"><span class="group-title-name">背景颜色</span></p>

        <div class="panel-menu-group">
            <button v-for="item in highlightColors" :key="item" :class="bgColorClass(item)" @click="onBgColorClick($event, item)">
                <svg width="22" height="22" viewBox="0 0 22 22" class="panel-btn-icon" :style="{ 'background-color': item }">
                    <path d="M10.28 5.5L6 16.5h1.34l1.16-3.08h4.99l1.16 3.08H16l-4.28-11h-1.44zm-1.39 6.86l2.09-5.5h.06l2.05 5.5H8.9z"></path>
                </svg>
            </button>
        </div>

        <button type="button" class="clear-btn" @click="onClearBtnClick">
            <span><span>清除</span></span>
        </button>
    </div>
</template>

<script>
    import { FONT_COLOR, HIGHLIGHT_COLOR } from '../../../common/constants'
    import { $emit } from '@ac/shared'

    export default {
        name: 'ColorPanel',
        data() {
            return {
                fontColors: FONT_COLOR,
                highlightColors: HIGHLIGHT_COLOR,
                type: null,
                attrs: {},
            }
        },

        methods: {
            bgColorClass(item) {
                return [
                    'panel-btn',
                    {
                        selected: item === this.attrs.bgColor,
                    },
                ]
            },

            fontColorClass(item) {
                return [
                    'panel-btn text',
                    {
                        selected: item === this.attrs.fontColor,
                    },
                ]
            },

            onFontColorClick(e, item) {
                e.stopPropagation()
                let fontColor = item
                if (fontColor === this.attrs.fontColor) {
                    fontColor = ''
                }
                this.updateAttr({ ...this.attrs, fontColor })
            },

            onBgColorClick(e, item) {
                e.stopPropagation()
                let bgColor = item

                if (bgColor === this.attrs.bgColor) {
                    bgColor = ''
                }
                this.updateAttr({ ...this.attrs, bgColor })
            },

            onClearBtnClick(e) {
                e.stopPropagation()
                this.updateAttr({})
            },

            updateAttr(attrs) {
                $emit(this, 'on-update', this.type, attrs)
            },
        },
    }
</script>
