<template>
    <div class="ac-preview-catalog scroll-content">
        <div class="sidebar-header"></div>
        <document-tree :tree="tree" />
    </div>
</template>

<script>
    import DocumentTree from './DocumentTree.vue'
    import { Storage } from '@natosoft/shared'
    import Tree from './Tree'

    export default {
        name: 'DocumentCatalog',
        components: {
            DocumentTree,
        },
        props: {
            list: {
                type: Array,
                default: () => [],
            },
            project: {
                type: Object,
                default: () => ({}),
            },
        },
        watch: {
            '$route.params.node_id': function () {
                let node_id = parseInt(this.$route.params.node_id, 10)
                this.node_id = isNaN(node_id) ? null : node_id
                this.activeNode()
            },

            list: function () {
                this.initProjectCatalog()
            },
        },
        data() {
            let node_id = parseInt(this.$route.params.node_id, 10)

            return {
                project_id: this.$route.params.project_id,
                node_id: isNaN(node_id) ? null : node_id,
                token: Storage.get(Storage.KEYS.SECRET_PROJECT_TOKEN + this.$route.params.project_id || '', true),
                tree: [],
            }
        },
        methods: {
            reactiveNode() {
                let hasSelectedNode = false
                Tree.traverse((item) => {
                    if (item.selected) {
                        hasSelectedNode = item.selected
                        return false
                    }
                }, this.tree)

                if (!hasSelectedNode) {
                    let node = null
                    Tree.traverse((item) => {
                        if (item.isLeaf) {
                            node = item
                            return false
                        }
                    }, this.tree)

                    let params = { project_id: this.project_id }
                    // 还存在文档
                    if (node) {
                        params.node_id = node.id

                        if (this.node_id !== node.id) {
                            this.$router.push({ name: 'preview.project.document', params })
                        }
                    }
                }
            },

            activeNode() {
                let node = null
                Tree.traverse((item) => {
                    item.selected = false
                    if (item.id === parseInt(this.node_id, 10)) {
                        item.selected = true
                        node = item
                    }
                }, this.tree)

                while (node && node.parent) {
                    node = node.parent
                    if (!node.isRoot) {
                        node.expanded = true
                    }
                }
            },

            initProjectCatalog() {
                let { root } = new Tree(this.list || [])
                this.tree = root.children || []
                // this.$$store.tree = root.children || []
                // console.log("项目目录：", this.tree);
                if (!this.node_id) {
                    this.reactiveNode()
                } else {
                    this.activeNode()
                }
            },
        },
        created() {
            this.initProjectCatalog()
        },
    }
</script>
