<template>
  <div class="flex min-h-70vh">
    <div class="w-315px bg-[var(--bg-color)]">
      <slot name="nav"></slot>
    </div>
    <div class="flex-1">
      <slot></slot>
    </div>
  </div>
</template>
