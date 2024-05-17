<script setup lang="ts">
import DownloadPane from "./download/DownloadPane.vue";
import ExportPane from "./export/ExportPane.vue";



</script>

<template>
  <n-flex justify="space-around" style="padding: 24px">
    <n-card title="下载" size="small" hoverable style="flex: 1">
      <download-pane/>
    </n-card>
    <n-card title="导出" size="small" hoverable style="flex: 1">
      <export-pane/>
    </n-card>
  </n-flex>
</template>

