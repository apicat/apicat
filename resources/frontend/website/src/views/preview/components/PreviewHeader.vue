<template>
    <header class="header header-fixed">
        <a class="logo" href="/">
            <img src="@/assets/image/logo.svg" alt="ApiCat" />
            <span class="logo"><span class="logo-text logo-apicat">ApiCat</span></span>
        </a>
        <ul class="header-nav"></ul>
        <project-search v-if="isShowSearch" />
    </header>
</template>

<script>
    import ProjectSearch from './ProjectSearch.vue'

    export default {
        name: 'PreviewHeader',
        props: {
            isShowSearch: {
                type: Boolean,
                default: true,
            },
        },
        components: {
            ProjectSearch,
        },
    }
</script>
