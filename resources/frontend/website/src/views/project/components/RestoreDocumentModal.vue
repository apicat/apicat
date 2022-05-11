<template>
    <el-dialog
        v-model="isShow"
        :width="400"
        custom-class="show-footer-line vertical-center-modal"
        :close-on-click-modal="false"
        append-to-body
        title="原文档所在分类已被删除，请选择其他分类"
    >
        <el-form ref="teamForm" :model="form" :rules="rules" label-position="top" style="margin-bottom: -19px" @keyup.enter="handleSubmit('teamForm')">
            <el-form-item label="" prop="dir_id" class="hide_required">
                <el-cascader v-model="form.dir_id" class="w-full" :options="dir" :props="{ checkStrictly: true }" placeholder="请选择分类">
                    <template #empty>暂无数据</template>
                </el-cascader>
            </el-form-item>
        </el-form>

        <template #footer>
            <el-button @click="onCloseBtnClick()"> 取消 </el-button>
            <el-button :loading="isLoading" type="primary" @click="handleSubmit('teamForm')"> 确定 </el-button>
        </template>
    </el-dialog>
</template>
<script>
    import { getDirList } from '@/api/dir'
    import { restoreApiDocument } from '@/api/project'
    import { ElMessage as $Message } from 'element-plus'
    import { h } from 'vue'

    export default {
        emits: ['on-ok'],
        data() {
            return {
                project_id: this.$route.params.project_id || '',
                isShow: false,
                isLoading: false,
                dir: [],
                form: {
                    dir_id: [],
                },
                rules: {
                    dir_id: { required: true, min_len: 1, message: '请选择分类', trigger: 'change', type: 'array' },
                },
            }
        },
        watch: {
            isShow: function () {
                !this.isShow && this.reset()
            },

            'form.dir_id': function () {
                let doc_id = this.form.dir_id.slice(-1)[0]
                doc_id !== undefined && (this.document.node_id = doc_id)
            },
        },
        methods: {
            show(document) {
                if (!document) {
                    throw new Error('文档信息不能为空！')
                }
                this.document = document
                this.isShow = true
            },
            hide() {
                this.isShow = false
            },

            onCloseBtnClick() {
                this.isShow = false
                this.reset()
            },

            reset() {
                this.$refs['teamForm'].resetFields()
            },

            handleSubmit(name) {
                this.$refs[name].validate((valid) => valid && this.submit())
            },

            submit() {
                this.isLoading = true
                restoreApiDocument(this.document)
                    .then((res) => {
                        this.onCloseBtnClick()

                        $Message({
                            type: 'success',
                            closable: true,
                            message: h('span', null, [
                                '文档恢复成功，',
                                h('a', { class: 'text-blue-600', href: `/editor/${this.project_id}/doc/${this.document.doc_id}` }, '查看详情'),
                            ]),
                        })
                        this.$emit('on-ok')
                    })
                    .finally(() => {
                        this.isLoading = false
                    })
            },

            transferDir(dir) {
                let result = []

                let loop = (arr, container) => {
                    ;(arr || []).forEach((item) => {
                        let o = {
                            value: item.id,
                            label: item.title,
                            children: [],
                        }

                        container.push(o)

                        if (item.sub_nodes && item.sub_nodes.length) {
                            loop(item.sub_nodes, o.children)
                        }
                    })
                }

                loop(dir, result)

                return [{ value: 0, label: '根目录' }].concat(result)
            },

            async getDocumentDirList(project_id) {
                getDirList(project_id).then(({ data }) => {
                    this.dir = this.transferDir(data)
                })
            },
        },

        mounted() {
            this.getDocumentDirList(this.project_id)
        },
    }
</script>
