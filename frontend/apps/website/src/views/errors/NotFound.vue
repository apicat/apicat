<template>
  <div class="text-center">
    <img class="not-found-image" src="@/assets/images/404@2x.png" alt="404">
    <p class="not-found-tip">
      {{ $t('app.tips.notFound') }}<router-link class="text-blue-600" to="/">
        {{ $t('app.common.goHome') }}
      </router-link>
    </p>
  </div>
</template>

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
