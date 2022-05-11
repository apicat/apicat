<template>
    <ul>
        <li v-for="item in tree" @click.stop="onItemClick(item)" :key="item.id" :class="getItemClass(item)">
            <div class="ac-preview-catalog_highlight" @mouseenter.stop="onItemMouseEnter(item, true)" @mouseleave.stop="onItemMouseEnter(item, false)" />

            <div class="flex items-center" @mouseenter.stop="onItemMouseEnter(item, true)" @mouseleave.stop="onItemMouseEnter(item, false)">
                <template v-if="!item.isLeaf">
                    <span class="text">{{ item.name }}</span>
                    <i class="icon iconfont iconarrow-right"></i>
                </template>

                <template v-else>
                    <a :id="'tree_' + item.id" href="javascript:void(0);" class="text dis_hover">{{ item.name }}</a>
                </template>
            </div>

            <document-tree v-if="item.children && item.children.length" :tree="item.children" />
        </li>
    </ul>
</template>

<script>
    export default {
        name: 'DocumentTree',
        props: {
            tree: {
                type: Array,
                default: () => [],
            },
        },

        data() {
            return {
                project_id: this.$route.params.project_id || '',
                node_id: this.$route.params.node_id || '',
            }
        },

        methods: {
            onItemClick(item) {
                if (!item.isLeaf) {
                    item.expanded = !item.expanded
                    return
                }
                this.$router.push({
                    name: 'preview.project.document',
                    params: { project_id: this.project_id, node_id: item.id },
                })
            },

            onItemMouseEnter(item, flag) {
                item.isHover = flag
            },

            getItemClass(item) {
                return {
                    open: item.expanded,
                    hover: item.isHover,
                    active: item.selected,
                }
            },
        },
    }
</script>
