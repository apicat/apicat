<template>
    <div class="ac-table">
        <!-- 表格 -->
        <el-table v-loading="loading" :data="tableData" :empty-text="emptyText" @selection-change="handleSelectionChange">
            <el-table-column v-if="isShowIndex" type="index" width="50" label="序号" />

            <el-table-column v-if="isShowSelection" type="selection" width="55" />
            <!-- 表头 -->
            <template v-for="column in columns" :key="column.prop">
                <el-table-column
                    v-if="column[dataKey.label] && !column.slot"
                    :prop="column[dataKey.prop]"
                    show-overflow-tooltip
                    :label="column[dataKey.label]"
                    :formatter="column.render || column.formatter"
                    :width="column.width"
                />
                <slot v-if="column.slot && column[dataKey.label]" :name="column.slot" />
            </template>
            <slot name="operation" />
        </el-table>

        <!-- 分页 -->
        <div v-if="isShowPager && pageTotal" class="mt-4 flex justify-end">
            <el-pagination
                :page-size="pageSize"
                layout="prev, pager, next"
                :total="pageTotal"
                v-model:current-page="currentPage"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
            />
        </div>
    </div>
</template>

<script>
    export default {
        name: 'AcTable',
        props: {
            isShowIndex: {
                type: Boolean,
                default: false,
            },
            isShowSelection: {
                type: Boolean,
                default: false,
            },
            columns: {
                type: Array,
                default: () => [],
            },
            tableData: {
                type: Array,
                default: () => [],
            },
            pageSize: {
                type: Number,
                default: () => 15,
            },
            pageTotal: {
                type: Number,
                default: () => 0, //总条数
            },
            pageCount: {
                type: Number,
                default: () => 0, //总页数
            },
            currentPage: {
                type: Number,
                default: () => 1, //当前页码
            },
            loading: {
                type: Boolean,
                default: false,
            },
            isShowPager: {
                type: Boolean,
                default: true,
            },
            emptyText: {
                type: String,
                default: '暂无数据',
            },
            dataKey: {
                type: Object,
                default: () => {
                    return { label: 'title', prop: 'key' }
                },
            },
        },

        setup(props, ctx) {
            const { emit } = ctx
            const handleSizeChange = (val) => emit('size-change', val)
            const handleCurrentChange = (val) => {
                emit('update:currentPage', props.currentPage)
                emit('current-change', val)
            }
            const handleSelectionChange = (val) => emit('selection-change', val)

            return {
                handleSizeChange,
                handleCurrentChange,
                handleSelectionChange,
            }
        },
    }
</script>
