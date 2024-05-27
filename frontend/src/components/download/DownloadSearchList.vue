<script setup lang="ts">
import {defineProps} from "vue"
import DownloadSearchBar from "./DownloadSearchBar.vue";
import {search} from "../../../wailsjs/go/models";

const props = defineProps<{
  downloadSearchBarRef: InstanceType<typeof DownloadSearchBar> | null
  searchByKeywordResult: search.ComicSearchResult[]
}>()

async function onSelectItem(comic: string) {
  if (props.downloadSearchBarRef === null) {
    return
  }

  await props.downloadSearchBarRef?.searchById(comic)
}

</script>

<template>
  <n-scrollbar>
    <div class="flex flex-col gap-y-2">
      <n-button v-for="result in searchByKeywordResult"
                :key="result.comicId"
                @click="onSelectItem(result.comicId)"
      >{{ result.title }}
      </n-button>
    </div>
  </n-scrollbar>
</template>
