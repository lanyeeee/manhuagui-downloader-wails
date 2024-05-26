<script setup lang="ts">
import DownloadTree from "./DownloadTree.vue";
import DownloadSearchBar from "./DownloadSearchBar.vue";
import {computed, ref, watch} from "vue";
import {TreeInst, TreeOption} from "naive-ui";
import DownloadList from "./DownloadList.vue";
import {DownloadStatus} from "../../constants/download-constant";
import {useDownloaderStore} from "../../stores/downloader";

const store = useDownloaderStore()


const downloadTreeInst = ref<TreeInst | null>(null);
const downloadTreeOptions = ref<TreeOption[]>([]);
const downloadDefaultExpandKeys = ref<string[]>([]);
const downloadDefaultCheckedKeys = ref<string[]>([]);
const searchDisabled = ref<boolean>(false);

watch(() => store.cacheDirectory, () => downloadTreeOptions.value = [])

const showInfo = computed<boolean>(() => downloadTreeOptions.value.length === 0)
const optionToDownload = computed<(TreeOption | null)[]>(() => downloadTreeInst
        .value
        ?.getCheckedData()
        .options
        // 所有状态不是已完成的叶子节点就是要下载的
        .filter(option => option !== null && option.isLeaf && (option.suffix === undefined || option.suffix() !== DownloadStatus.COMPLETED))
    ?? [])
const optionDownloading = computed<(TreeOption | null)[]>(() => downloadTreeInst
        .value
        ?.getCheckedData()
        .options
        .filter(option => option !== null
            && option.isLeaf
            && option.suffix !== undefined
            && (option.suffix() === DownloadStatus.WAITING || option.suffix() === DownloadStatus.DOWNLOADING))
    ?? [])

</script>

<template>
  <div class="flex h-full">
    <div class="flex-col flex-1 p-2">
      <download-search-bar :disabled="searchDisabled"
                           v-model:download-tree-options="downloadTreeOptions"
                           v-model:download-default-expand-keys="downloadDefaultExpandKeys"
                           v-model:download-default-checked-keys="downloadDefaultCheckedKeys"
      />
      <n-result v-if="showInfo" title="在搜索框中输入漫画链接"/>
      <download-tree v-else
                     v-model:download-tree-inst="downloadTreeInst"
                     v-model:download-tree-options="downloadTreeOptions"
                     v-model:download-default-expand-keys="downloadDefaultExpandKeys"
                     v-model:download-default-checked-keys="downloadDefaultCheckedKeys"
      />
    </div>
    <div class="flex-1">
      <download-list :download-tree-inst="downloadTreeInst"
                     :download-tree-options="downloadTreeOptions"
                     :options-to-download="optionToDownload"
                     :options-downloading="optionDownloading"
                     v-model:search-disabled="searchDisabled"
      />
    </div>
  </div>
</template>

