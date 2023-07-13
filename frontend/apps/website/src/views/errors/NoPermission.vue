<template>
  <div class="text-center">
    <iconfont icon="ac-lock" :size="300" class="not-found-image" />
    <p class="not-found-tip">
      {{ $t('app.tips.noPermission') }}<router-link class="text-blue-600" :to="MAIN_PATH">{{ $t('app.common.goProjectList') }}</router-link>
    </p>
  </div>
</template>
<script setup>
import { MAIN_PATH } from '@/router/constant'
</script>
<style>
.not-found-image {
  height: 282px;
  margin-bottom: 6vh;
  margin-top: 20vh;
  display: inline-block;
}

.not-found-tip {
  font-size: 18px;
  font-family: PingFangSC-Regular, PingFang SC;
  font-weight: 400;
  color: rgba(0, 0, 0, 0.65);
  line-height: 25px;
  letter-spacing: 4px;
}
</style>
