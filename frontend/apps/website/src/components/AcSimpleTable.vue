<script>
export default defineComponent({
  name: 'AcSimpleTable',
  props: {
    border: {
      type: Boolean,
      default: true,
    },
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
    roundBorder: {
      type: Boolean,
      default: false,
    },
    rowClassName: {
      type: String,
      default: undefined,
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
    const handleSizeChange = val => emit('update:page-size', val)
    const handleCurrentChange = (val) => {
      emit('update:page', val)
    }
    const handleSelectionChange = val => emit('selection-change', val)

    return {
      handleSizeChange,
      handleCurrentChange,
      handleSelectionChange,
    }
  },
})
</script>

<template>
  <div class="SimpleTable" :class="roundBorder ? 'roundBorder' : undefined">
    <slot name="search-form" />

    <!-- 表格 -->
    <el-table
      v-loading="loading"
      class="w-full"
      :row-class-name="rowClassName"
      :data="tableData"
      :border="border"
      @selection-change="handleSelectionChange"
    >
      <el-table-column v-if="isShowIndex" type="index" width="50" label="" align="center" />

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
        <slot v-if="column.slot" :name="column.slot" />
      </template>
      <slot name="operation" />
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

<style lang="scss" scoped>
.hideBackground {
  :deep(.el-table--enable-row-hover .el-table__body tr:hover > td.el-table__cell) {
    background-color: transparent;
  }
}

.roundBorder {
  :deep(.el-table__body) {
    border-collapse: separate;
    border-spacing: 0 5px;
  }
  :deep(.el-table--enable-row-hover .el-table__body tr td, .el-table__header thead tr th) {
    margin-top: 5px;
    border-bottom: none !important;
  }
  :deep(.el-table__inner-wrapper::before) {
    background-color: transparent !important;
  }
  :deep(.el-table--enable-row-hover .el-table__body tr td:first-child) {
    border-top-left-radius: 5px;
    border-bottom-left-radius: 5px;
  }
  :deep(.el-table--enable-row-hover .el-table__body tr td:last-child) {
    border-bottom-right-radius: 5px;
    border-top-right-radius: 5px;
  }
  :deep(.el-table__header thead tr th:first-child) {
    border-top-left-radius: 5px;
    border-bottom-left-radius: 5px;
  }
  :deep(.el-table__header thead tr th:last-child) {
    border-bottom-right-radius: 5px;
    border-top-right-radius: 5px;
  }
}
</style>
