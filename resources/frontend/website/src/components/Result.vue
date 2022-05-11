<template>
    <div class="result">
        <div class="result-icon" :style="styles">
            <slot name="icon"> </slot>
        </div>
        <slot name="title">
            <p class="result-title">
                {{ title }}
            </p>
        </slot>

        <p class="mt-1 result-subtitle">
            <slot name="subTitle">
                {{ subTitle }}
            </slot>
        </p>
    </div>
</template>

<script>
    export default {
        props: {
            status: {
                type: String,
                default: () => 'success',
                validate(value) {
                    return ['success', 'error'].indexOf(value) !== -1
                },
            },
            title: {
                type: String,
                default: () => '',
            },
            subTitle: {
                type: String,
                default: () => '',
            },
            styles: {
                type: Object,
                default: () => ({}),
            },
        },
    }
</script>
<style lang="scss" scoped>
    .result {
        padding: 48px 32px;
        text-align: center;

        &-icon {
            width: 95px;
            height: 95px;
            margin: auto;
            img {
                width: 100%;
            }
        }

        &-title {
            font-size: 24px;
            font-weight: 500;
            color: theme('colors.zinc.600');
            margin-top: 6px;
        }

        &-subtitle {
            max-width: 502px;
            margin: auto;
            color: theme('colors.zinc.500');
        }
    }
</style>
