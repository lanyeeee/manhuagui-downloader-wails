<script setup lang="ts">
import ExportTree from "./ExportTree.vue";
import ExportBar from "./ExportBar.vue";
import ExportRefreshButton from "./ExportRefreshButton.vue";
import CacheDirectoryInput from "../settings/CacheDirectoryInput.vue";
</script>

<template>
  <div class="flex h-full">
    <div class="flex-1 flex flex-col p-2">
      <div class="flex gap-x-2">
        <cache-directory-input/>
        <export-refresh-button/>
      </div>
      <export-tree/>
    </div>
    <div class="flex-1 flex flex-col p-2">
      <div class="flex-1"></div>
      <export-bar/>
    </div>
  </div>
</template>
