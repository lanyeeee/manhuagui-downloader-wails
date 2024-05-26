<script setup lang="ts">
import {defineProps, ref} from "vue"
import {TreeInst, TreeOption} from "naive-ui";
import DownloadProgress from "./DownloadProgress.vue";
import DownloadButton from "./DownloadButton.vue";

const props = defineProps<{
  downloadTreeInst: TreeInst | null,
  downloadTreeOptions: TreeOption[],
  optionsToDownload: (TreeOption | null)[],
  optionsDownloading: (TreeOption | null)[],
}>()

const downloadProgresses = ref<(InstanceType<typeof DownloadProgress>)[]>([])

const searchDisabled = defineModel<boolean>("searchDisabled", {required: true})

</script>

<template>
  <n-scrollbar>
    <div class="flex flex-col gap-y-2">
      <n-h3>下载队列</n-h3>
      <download-progress v-for="option in props.optionsDownloading"
                         :key="option?.key"
                         ref="downloadProgresses"
                         :tree-option="option"
      />
    </div>
  </n-scrollbar>

  <download-button
      :download-tree-inst="downloadTreeInst"
      :download-tree-options="downloadTreeOptions"
      :options-to-download="optionsToDownload"
      :download-progresses="downloadProgresses"
      v-model:search-disabled="searchDisabled"
  />
</template>
