<template>
    <div class="ac-hero">
        <div class="viewport ac-hero__main lg:flex lg:justify-between text-center lg:text-left">
            <div class="ac-hero__desc flex flex-col justify-center lg:pb-0 pb-4 text-3xl leading-10 pt-9">
                <h2>ApiCat是一款</h2>
                <h2>API协作开发提效软件</h2>
            </div>
            <img class="ac-hero__img" src="@/assets/image/banner.png" alt="" />
        </div>
    </div>

    <div class="ac-feature mt-10">
        <div class="viewport grid lg:grid-cols-4 lg:gap-10 px-5 gap-8">
            <div class="ac-feature__item" v-for="item in featureList" :key="item.title">
                <img class="w-10" :src="item.icon" alt="" />
                <p class="ac-feature__subtitle pt-7 lg:pb-4 pb-2">{{ item.title }}</p>
                <p class="ac-feature__desc">{{ item.desc }}</p>
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
    import { useBackground } from '@/hooks/useBackground'
    import iconApi from '@/assets/image/icon-api.png'
    import iconTeam from '@/assets/image/icon-team.png'
    import iconMock from '@/assets/image/icon-mock.png'
    import iconImport from '@/assets/image/icon-import.png'

    useBackground(document.getElementById('app'))

    const featureList = [
        {
            icon: iconApi,
            title: 'API文档',
            desc: '方便的API文档管理，无论是HTTP、Websocket还是TCP、UDP等协议，都可以轻松描述。Block化的编辑器，同时还支持Markdown语法。',
        },
        { icon: iconTeam, title: '团队协作', desc: '不同团队自由切换，项目相互隔离，一个团队一个身份，灵活的权限控制，高效管理团队中各个项目和成员。' },
        { icon: iconMock, title: '方便的Mock', desc: '便捷的Mock功能，简单的语法，能方便的模拟出你想要的数据。' },
        { icon: iconImport, title: '导入导出', desc: '支持Postman、Swagger、JMeter等导入导出，还可以将文档导出为PDF、Markdown、HTML等格式。' },
    ]
</script>
<style lang="scss">
    @use '@/assets/stylesheet/mixins/mixins' as *;

    @include b(hero) {
        background-image: url(@/assets/image/banner-bg.png);
        background-repeat: no-repeat;
        background-position: left bottom;
        background-size: cover;
        height: 550px;

        @include e(desc) {
            font-family: PingFangSC-Regular, PingFang SC;
            color: rgba(0, 0, 0, 0.85);
        }

        @include e(img) {
            width: 80%;
            margin: auto;
        }

        @screen sm {
            @include e(img) {
                width: 50%;
            }
        }

        @screen lg {
            @include e(main) {
                padding-top: 52px;
            }

            @include e(desc) {
                font-size: 44px;
                font-family: PingFangSC-Regular, PingFang SC;
                font-weight: 400;
                color: rgba(0, 0, 0, 0.85);
                line-height: 68px;
            }
            @include e(img) {
                width: 615px;
            }
        }
    }

    @include b(feature) {
        position: relative;

        @include e(title) {
            text-align: center;
            font-size: 40px;
            font-weight: 500;
            color: #ffffff;
            line-height: 68px;
        }

        @include e(list) {
            position: absolute;
            z-index: 10;
            margin-left: auto;
            margin-right: auto;
            left: 0;
            right: 0;
            background: #ffffff;
            border-radius: 8px;
            padding: 48px 40px;
        }

        @include e(item) {
            &:first-child {
                padding-left: 0;
            }
        }

        @include e(subtitle) {
            font-size: 24px;
            color: rgba(0, 0, 0, 0.85);
        }

        @include e(desc) {
            font-size: 16px;
            font-weight: 400;
            color: rgba(0, 0, 0, 0.65);
            line-height: 28px;
        }
    }
</style>
