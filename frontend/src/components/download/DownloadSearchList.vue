<script setup lang="ts">
import {defineProps, ref} from "vue"
import DownloadSearchBar from "./DownloadSearchBar.vue";
import {search} from "../../../wailsjs/go/models";
import ComicSearchResult = search.ComicSearchResult;

const paginationDisabled = ref<boolean>(false)

const props = defineProps<{
  downloadSearchBarRef: InstanceType<typeof DownloadSearchBar> | null
  searchByKeywordResult: ComicSearchResult,
}>()

async function onSelectItem(comic: string) {
  if (props.downloadSearchBarRef === null || props.searchByKeywordResult === undefined) {
    return
  }

  await props.downloadSearchBarRef.searchById(comic)
}

async function onPageChange(pageNumber: number) {
  if (props.downloadSearchBarRef === null) {
    return
  }

  paginationDisabled.value = true
  await props.downloadSearchBarRef.searchByKeyword(props.downloadSearchBarRef.searchByKeywordInput, pageNumber)
  paginationDisabled.value = false
}

</script>

<template>
  <div>
    <n-scrollbar>
      <div class="flex flex-col gap-y-2">
        <n-button v-for="info in searchByKeywordResult?.infos"
                  :key="info.comicId"
                  size="large"
                  @click="onSelectItem(info.comicId)">
          《{{ info.title }}》<br/> 作者：{{ info.authors }}
        </n-button>
      </div>
    </n-scrollbar>
    <n-pagination v-model:page="searchByKeywordResult.currentPage"
                  :page-count="searchByKeywordResult.totalPage"
                  :disabled=paginationDisabled
                  @update:page="onPageChange"
                  show-quick-jumper
    >
      <template #goto>
        跳转至
      </template>
    </n-pagination>
  </div>
</template>
