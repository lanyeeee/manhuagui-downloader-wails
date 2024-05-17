<script setup lang="ts">

import DownloadTree from "./DownloadTree.vue";
import DownloadSearchBar from "./DownloadSearchBar.vue";
import DownloadBar from "./DownloadBar.vue";
</script>

<template>
  <n-flex vertical>
    <n-flex>
      <download-search-bar/>
    </n-flex>
    <download-tree/>
    <download-bar/>
  </n-flex>
</template>

