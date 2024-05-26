<script setup lang="ts">

import DownloadTree from "./DownloadTree.vue";
import DownloadSearchBar from "./DownloadSearchBar.vue";
import DownloadBar from "./DownloadBar.vue";
</script>

<template>
  <div class="flex flex-col h-full">
    <download-search-bar/>
    <download-tree/>
    <download-bar/>
  </div>
</template>

