<template>
  <div class="flex w-full ac-top-nav shadow-black h-13">
    <AcLogo class="ac-top-nav-logo ml-7 lg:ml-15vw" href="/" :size="90" />
  </div>
  <slot />
</template>

<style lang="scss" scoped>
.ac-top-nav {
  box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.2);
}
* {
  transition: all 0.2s;
}
</style>
