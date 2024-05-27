<script setup lang="ts">
import DownloadTree from "./DownloadTree.vue";
import DownloadSearchBar from "./DownloadSearchBar.vue";
import {computed, ref, watch} from "vue";
import {TreeInst, TreeOption} from "naive-ui";
import DownloadList from "./DownloadList.vue";
import {DownloadStatus} from "../../constants/download-constant";
import {useDownloaderStore} from "../../stores/downloader";
import {search} from "../../../wailsjs/go/models";
import DownloadSearchList from "./DownloadSearchList.vue";
import ComicSearchResult = search.ComicSearchResult;

const store = useDownloaderStore()

const downloadTreeInst = ref<TreeInst | null>(null);
const downloadTreeOptions = ref<TreeOption[]>([]);
const downloadDefaultExpandKeys = ref<string[]>([]);
const downloadDefaultCheckedKeys = ref<string[]>([]);
watch(() => store.cacheDirectory, () => downloadTreeOptions.value = [])

const searchDisabled = ref<boolean>(false);
const searchResultType = ref<"empty" | "tree" | "list">("empty")

const downloadSearchBarRef = ref<InstanceType<typeof DownloadSearchBar> | null>(null);
const searchByKeyWordResult = ref<ComicSearchResult[]>([]);

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
    <div class="flex-1 flex flex-col p-2">
      <download-search-bar :disabled="searchDisabled"
                           v-model:download-tree-options="downloadTreeOptions"
                           v-model:download-default-expand-keys="downloadDefaultExpandKeys"
                           v-model:download-default-checked-keys="downloadDefaultCheckedKeys"
                           v-model:search-by-keyword-result="searchByKeyWordResult"
                            v-model:search-result-type="searchResultType"
                           ref="downloadSearchBarRef"
      />
      <n-result v-if="searchResultType == 'empty'" title="请搜索漫画"/>
      <download-tree v-else-if="searchResultType == 'tree'"
                     ref="downloadSearchBar"
                     v-model:download-tree-inst="downloadTreeInst"
                     v-model:download-tree-options="downloadTreeOptions"
                     v-model:download-default-expand-keys="downloadDefaultExpandKeys"
                     v-model:download-default-checked-keys="downloadDefaultCheckedKeys"
      />
      <download-search-list v-else-if="searchResultType == 'list'"
                            :search-by-keyword-result="searchByKeyWordResult"
                            :download-search-bar-ref="downloadSearchBarRef"
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

