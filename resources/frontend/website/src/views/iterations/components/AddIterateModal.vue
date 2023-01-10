<template>
    <el-dialog v-model="isShow" :width="400" :close-on-click-modal="false" :title="title" append-to-body class="show-footer-line">
        <el-form
            ref="iterateForm"
            :model="form"
            :rules="rules"
            label-position="top"
            class="py-2 pl-2"
            @submit.prevent="handleSubmit(iterateForm)"
            v-loading="isLoadingForProject"
        >
            <el-form-item label="迭代名称" prop="title" class="hide_required">
                <el-input v-model="form.title" placeholder="迭代名称" maxlength="255" />
            </el-form-item>

            <el-form-item label="所属项目" prop="project_id" class="hide_required">
                <el-select
                    :disabled="isEdit"
                    v-model="form.project_id"
                    filterable
                    no-data-text="暂无数据"
                    no-match-text="暂无数据"
                    placeholder="请选择所属项目"
                    class="w-full"
                    popper-class="w-full"
                    :teleported="false"
                >
                    <el-option v-for="item in projects" :key="item.id" :label="item.name" :value="item.id" />
                </el-select>
            </el-form-item>

            <el-form-item label="迭代描述">
                <el-input v-model="form.description" type="textarea" :autosize="{ minRows: 4, maxRows: 4 }" maxlength="255" />
            </el-form-item>
        </el-form>

        <template v-slot:footer>
            <el-button @click="onCloseBtnClick()"> 取消 </el-button>
            <el-button :loading="isLoading" type="primary" @click="handleSubmit(iterateForm)"> 确 定 </el-button>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
    import { ref, toRaw, watch, reactive } from 'vue'
    import { createIteration, editIteration } from '@/api/iterate'
    import type { FormInstance } from 'element-plus'
    import { useIterateStore } from '@/stores/iterate'
    import { storeToRefs } from 'pinia'
    import useApi from '@/hooks/useApi'
    import { getEditableProjectListForCreateIterate } from '@/api/project'

    const iterateStore = useIterateStore()
    const { activeTab } = storeToRefs(iterateStore)

    const emit = defineEmits(['ok'])

    const [isLoadingForProject, fetchEditableProjectListForCreateIterate] = useApi(getEditableProjectListForCreateIterate)

    const iterateForm = ref<FormInstance>()
    const isShow = ref(false)
    const isLoading = ref(false)
    const isEdit = ref(false)
    const title = ref('新建迭代')
    const projects: any = ref([])

    let execute: any = null

    const initForm: any = {
        title: '',
        project_id: '',
        description: '',
    }

    let form = reactive({
        ...initForm,
    })

    const rules = {
        title: [
            { required: true, message: '请输入迭代名称', trigger: 'blur' },
            { min: 2, message: '迭代名称不能少于两个字', trigger: 'blur' },
        ],
        project_id: { required: true, message: '请选择所属项目', trigger: 'change' },
    }

    const show = (iterate: any) => {
        isShow.value = true
        isEdit.value = !!iterate
        title.value = isEdit.value ? '编辑迭代' : '新建迭代'
        execute = isEdit.value ? editIteration : createIteration
        iterate && (form = Object.assign(form, iterate))

        getUsableProjectList()
    }

    const hide = () => {
        isShow.value = false
    }

    const reset = () => {
        iterateForm.value?.resetFields()
        form = Object.assign(form, initForm)
        form.project_id = activeTab.value || ''
    }

    const handleSubmit = async (formEl: FormInstance | undefined) => {
        if (!formEl) return
        await formEl.validate(async (valid) => {
            if (valid) {
                try {
                    let data = { ...toRaw(form) }
                    if (isEdit.value) {
                        data = {
                            iteration_id: form.id,
                            title: form.title,
                            description: form.description,
                        }
                    }

                    isLoading.value = true
                    try {
                        await execute(data)
                        iterateStore.switchActiveCollectTab(form.project_id)
                        emit('ok')
                        hide()
                    } catch (e) {
                        //
                    } finally {
                        isLoading.value = false
                    }
                } catch (e) {
                    //
                }
            }
        })
    }

    const onCloseBtnClick = () => {
        hide()
    }

    const getUsableProjectList = async () => {
        const { data } = await fetchEditableProjectListForCreateIterate()
        projects.value = data || []

        const projectItem = projects.value.find((item: any) => item.id === activeTab.value)
        // 选中默认项目
        form.project_id = projectItem ? projectItem.id : ''
        setTimeout(() => {
            iterateForm.value?.clearValidate()
        }, 0)
    }

    watch(
        () => isShow.value,
        () => !isShow.value && reset()
    )

    // // 收藏项目 切换
    // watch(activeTab, async () => {
    //     form.project_id = activeTab.value || ''
    //     iterateForm.value?.clearValidate()
    // })

    defineExpose({
        show,
    })
</script>
