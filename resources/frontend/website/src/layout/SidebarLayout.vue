<template>
    <div class="flex justify-between flex-col lg:flex-row">
        <!-- 左侧导航 -->
        <aside class="sidebar-nav flex-none lg:w-64 mb-4 lg:mb-0">
            <slot><SidebarMenu v-bind="$attrs" /></slot>
        </aside>

        <!-- 右侧内容 -->
        <div class="flex-auto lg:ml-4 overflow-hidden">
            <RouterView />
        </div>
    </div>
</template>
