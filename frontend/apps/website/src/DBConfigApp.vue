<script setup lang="ts">
import type { FormInstance } from 'element-plus'
import { saveDBConfig } from '@/api/systemConfig'
import useApi from '@/hooks/useApi'

interface DBConfig {
  [key: string]: {
    value: string,
    type: 'env' | 'value',
  }
}

const [isLoading, saveDBConfigApi] = useApi(saveDBConfig)
const ruleFormRef = ref<FormInstance>()
// 数据类型：值，环境变量
const DataTypes = [
  { label: '值', value: 'value' },
  { label: '环境变量', value: 'env' },
]

const config = ref<DBConfig>((window as any)['DB_CONFIG'] || {
  host: {
    value: '',
    type: 'value',
  },
  port: {
    value: '',
    type: 'value',
  },
  user: {
    value: '',
    type: 'value',
  },
  password: {
    value: '',
    type: 'value',
  },
  dbname: {
    value: '',
    type: 'value',
  },
})

const rules = reactive({
  'host.value': { required: true, message: '请输入Host', trigger: 'blur' },
  'port.value': { required: true, message: '请输入Port', trigger: 'blur' },
  'user.value': { required: true, message: '请输入数据库用户名称', trigger: 'blur' },
  'dbname.value': { required: true, message: '请输入数据库名称', trigger: 'blur' },
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return

  try {
    await formEl.validate()
    await saveDBConfigApi(toRaw(config.value))
  } catch (error) {
    //
  }

}

</script>

<template>
  <el-form :model="config" :rules="rules" ref="ruleFormRef" class="db-config" label-position="top" size="large"
    @keyup.enter="submitForm(ruleFormRef)" @submit.prevent="submitForm(ruleFormRef)">
    <h1 class="text-center mb-20px text-24px">MySQL 设置</h1>

    <el-form-item label="Host" prop="host.value">
      <el-input v-model="config.host.value" placeholder="Host">
        <template #append>
          <el-select v-model="config.host.type" class="w-110px">
            <el-option v-for="item in DataTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item label="Port" prop="port.value">
      <el-input v-model="config.port.value" placeholder="Port" maxlength="9">
        <template #append>
          <el-select v-model="config.port.type" class="w-110px">
            <el-option v-for="item in DataTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item label="User" prop="user.value">
      <el-input v-model="config.user.value" placeholder="User" maxlength="50">
        <template #append>
          <el-select v-model="config.user.type" class="w-110px">
            <el-option v-for="item in DataTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item label="Password" prop="password.value">
      <el-input v-model="config.password.value" placeholder="Password" maxlength="100">
        <template #append>
          <el-select v-model="config.password.type" class="w-110px">
            <el-option v-for="item in DataTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item label="Database" prop="dbname.value">
      <el-input v-model="config.dbname.value" placeholder="Database" maxlength="100">
        <template #append>
          <el-select v-model="config.dbname.type" class="w-110px">
            <el-option v-for="item in DataTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item class="mt-40px">
      <el-button type="primary" :loading="isLoading" class="w-full" @click="submitForm(ruleFormRef)">保存</el-button>
    </el-form-item>
  </el-form>
</template>

<style lang="scss" scoped>
.db-config {
  width: 400px;

  :deep(.el-form-item__label) {
    font-weight: bold;
  }
}
</style>
