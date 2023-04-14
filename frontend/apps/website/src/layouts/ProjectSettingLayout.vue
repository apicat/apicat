<template>
  <div class="flex h-screen max-h-screen">
    <div class="bg-gray-100 w-315px pt-20px pl-30px">
      <slot name="nav"></slot>
    </div>
    <div class="flex-1 pt-20px px-30px">
      <slot></slot>
    </div>
  </div>
</template>
