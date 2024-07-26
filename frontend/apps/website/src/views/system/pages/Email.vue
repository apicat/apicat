<script setup lang="ts">
import { apiGetEmail } from '@/api/system'
import { SysEmail } from '@/commons'
import { useCollapse } from '@/components/collapse/useCollapse'
import SMTP from '@/views/system/pages/Email/SMTP.vue'
import SendCloud from '@/views/system/pages/Email/SendCloud.vue'

const tBase = 'app.system.email'
const collapse = useCollapse({})

interface A {
  [SysEmail.SMTP]: SystemAPI.EmailSMTP
  [SysEmail.SendCloud]: SystemAPI.EmailSendCloud
}

const data = ref<A>({
  [SysEmail.SMTP]: {
    host: '',
    user: '',
    address: '',
    password: '',
  },
  [SysEmail.SendCloud]: {
    apiUser: '',
    apiKey: '',
    fromAddress: '',
    fromName: '',
  },
})

apiGetEmail().then((res) => {
  for (let i = 0; i < res.length; i++) {
    const v = res[i]
    data.value[v.driver as keyof A] = v.config as any
    if (v.use)
      collapse.ctx.open(v.driver)
  }
})
</script>

<template>
  <div class="bg-white w-85%">
    <h1>{{ $t(`${tBase}.title`) }}</h1>

    <div class="flex flex-col">
      <SMTP v-model:config="data[SysEmail.SMTP]" class="collapse-box" :name="SysEmail.SMTP" :collapse="collapse" />
      <SendCloud
        v-model:config="data[SysEmail.SendCloud]"
        class="collapse-box mt-30px"
        :name="SysEmail.SendCloud"
        :collapse="collapse"
      />
    </div>
  </div>
</template>

<style scoped>
.row {
  margin-top: 1em;
  margin-bottom: 1em;
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.left,
.right {
  display: flex;
  align-items: center;
}

.left {
  justify-content: flex-start;
  /* flex-grow: 1; */
}

.right {
  /* justify-content: flex-end; */
  flex-grow: 1;
}
</style>
