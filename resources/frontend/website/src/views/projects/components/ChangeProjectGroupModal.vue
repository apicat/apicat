<template>
    <el-dialog v-model="isShow" :width="340" :close-on-click-modal="false" title="移动项目分组" append-to-body class="show-footer-line">
        <el-form @submit.prevent="handleSubmit('teamForm')" ref="teamForm" :model="form" :rules="rules" label-position="top" style="margin-bottom: -19px">
            <el-form-item label="选择项目分组" prop="new_group_id" class="hide_required">
                <el-select
                    class="w-full"
                    v-model="form.new_group_id"
                    not-data-text="暂无项目分组"
                    no-match-text="未查找到分组"
                    placeholder="选择被移动的项目分组"
                    :teleported="false"
                    popper-class="w-full overflow-hidden"
                >
                    <el-option :value="0" label="不分组" />
                    <el-option v-for="item in projectGroupList" :value="item.id" :key="item.id" :label="item.name" />
                </el-select>
            </el-form-item>
        </el-form>

        <template #footer>
            <div>
                <el-button @click="onCloseBtnClick()">取消</el-button>
                <el-button :loading="isLoading" type="primary" @click="handleSubmit('teamForm')">确定</el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script>
    import { changeProjectGroup } from '@/api/project'
    import { useProjectsStore } from '@/stores/projects'
    import { storeToRefs } from 'pinia'

    export default {
        emits: ['on-ok'],
        data() {
            return {
                isShow: false,
                isLoading: false,
                form: {
                    new_group_id: '',
                    project_id: '',
                    old_group_id: '',
                },
                rules: {
                    new_group_id: {
                        required: true,
                        message: '请选择项目分组',
                        type: 'number',
                        trigger: 'blur',
                    },
                },
            }
        },

        watch: {
            isShow: function () {
                !this.isShow && this.reset()
            },
        },

        methods: {
            show(project) {
                this.form.project_id = project.id
                let groupId = project.group_id || parseInt(this.$route.params.cid || 0, 10)
                this.form.old_group_id = isNaN(groupId) ? 0 : groupId
                this.form.new_group_id = isNaN(groupId) ? 0 : groupId
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
                this.form.project_id = ''
                this.form.old_group_id = ''
                this.form.new_group_id = ''
            },

            handleSubmit(name) {
                this.$refs[name].validate((valid) => valid && this.submit())
            },

            submit() {
                this.isLoading = true
                changeProjectGroup(this.form)
                    .then((res) => {
                        this.$Message.success(res.msg || '更换分组成功!')
                        this.onCloseBtnClick()
                        this.$emit('on-ok')
                    })
                    .finally(() => {
                        this.isLoading = false
                    })
                    .catch((e) => {})
            },
        },

        setup() {
            const state = useProjectsStore()

            return {
                ...storeToRefs(state),
            }
        },
    }
</script>
