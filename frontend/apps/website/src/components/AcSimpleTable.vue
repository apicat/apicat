<template>
  <div class="SimpleTable">
    <slot name="search-form"></slot>

    <!-- 表格 -->
    <el-table v-loading="loading" :data="tableData" @selection-change="handleSelectionChange" :border="true" class="w-full">
      <el-table-column v-if="isShowIndex" type="index" width="50" label="" />

      <el-table-column v-if="isShowSelection" type="selection" width="55" />
      <!-- 表头 -->
      <template v-for="column in columns" :key="column.prop">
        <el-table-column
          v-if="column[dataKey.label] && !column.slot"
          :prop="column[dataKey.prop]"
          show-overflow-tooltip
          :label="column[dataKey.label]"
          :formatter="column.formatter"
          :min-width="column.width"
        />
        <slot v-if="column.slot" :name="column.slot"></slot>
      </template>
      <slot name="operation"></slot>
    </el-table>

    <div class="flex justify-end mt-4">
      <!-- 分页 -->
      <el-pagination
        v-if="isShowPager && total"
        :background="false"
        :page-size="pageSize"
        layout="prev, pager, next"
        :total="total"
        :current-page="page"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script>
export default defineComponent({
  name: 'AcSimpleTable',
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
      default: () => 10,
    },
    total: {
      type: Number,
      default: () => 0, // 总条目数
    },
    page: {
      type: Number,
      default: () => 1, // 当前页码
    },
    loading: {
      type: Boolean,
      default: false,
    },
    isShowPager: {
      type: Boolean,
      default: false,
    },

    dataKey: {
      type: Object,
      default: () => {
        return { label: 'label', prop: 'prop' }
      },
    },
  },

  setup(props, ctx) {
    const { emit } = ctx
    const handleSizeChange = (val) => emit('update:page-size', val)
    const handleCurrentChange = (val) => emit('update:page', val)
    const handleSelectionChange = (val) => emit('selection-change', val)

    return {
      handleSizeChange,
      handleCurrentChange,
      handleSelectionChange,
    }
  },
})
</script>
