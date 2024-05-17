<script setup lang="ts">
import {computed, ref} from "vue"
import {TreeOption, useNotification} from "naive-ui"
import {useDownloaderStore} from "../../stores/downloader"
import {SearchOutline as SearchIcon} from "@vicons/ionicons5"
import {DownloadStatus} from "../../constants/download-constant"
import * as path from "../../../wailsjs/go/api/PathApi"
import {SearchComicInfo} from "../../../wailsjs/go/api/DownloadApi"
import {search, types} from "../../../wailsjs/go/models"


const store = useDownloaderStore()
const notification = useNotification()

const searchInput = ref<string>("")
const loading = ref<boolean>(false)
const disabled = computed<boolean>(() => store.searchDisabled)

async function onSearch() {
  if (loading.value || disabled.value) {
    return
  }

  const comicId = extractComicIdFromInput()
  if (!comicId) {
    notification.create({type: "error", title: "搜索失败", content: "请输入漫画ID或漫画链接", duration: 2000,})
    return
  }

  try {
    loading.value = true
    const response: types.Response = await SearchComicInfo(comicId, store.proxyUrl)
    if (response.code != 0) {
      notification.create({type: "error", title: "搜索失败", meta: response.msg,})
      return
    }

    const comicInfo: search.ComicInfo = response.data as search.ComicInfo

    console.log("搜索结果", comicInfo)

    await handleSearchResult(comicInfo)
  } finally {
    loading.value = false
  }
}

async function handleSearchResult(comicInfo: search.ComicInfo) {
  const defaultExpandedKeys: string[] = []
  const defaultCheckedKeys: string[] = []

  const comicTreeOption: TreeOption = {
    label: comicInfo.title,
    key: comicInfo.title,
    children: [],
    isLeaf: false
  }
  defaultExpandedKeys.push(comicTreeOption.key as string)

  for (const chapterType of comicInfo.chapterTypes) {
    const chapterTypeTreeOption: TreeOption = {
      label: chapterType.title,
      key: await path.Join([comicInfo.title, chapterType.title]),
      children: [],
      isLeaf: false
    }
    defaultExpandedKeys.push(chapterTypeTreeOption.key as string)


    for (const chapterPager of chapterType.chapterPagers) {
      const chapterPagerTreeOption: TreeOption = {
        label: chapterPager.title,
        key: await path.Join([comicInfo.title, chapterType.title, chapterPager.title]),
        children: [],
        isLeaf: false
      }

      for (const chapter of chapterPager.chapters) {
        const saveDirectory = await path.Join([store.cacheDirectory, chapterPagerTreeOption.key as string, chapter.title])

        const key = JSON.stringify({
          href: chapter.href,
          saveDirectory: saveDirectory
        })

        const directoryExists = await path.PathExists(saveDirectory)
        const chapterTreeOption: TreeOption = {
          label: chapter.title,
          key: key,
          isLeaf: true,
          disabled: directoryExists,
        }
        if (directoryExists) {
          chapterTreeOption.suffix = () => DownloadStatus.COMPLETED
          defaultCheckedKeys.push(key)
        }

        chapterPagerTreeOption.children?.push(chapterTreeOption)
      }

      chapterTypeTreeOption.children?.push(chapterPagerTreeOption)
    }

    comicTreeOption.children?.push(chapterTypeTreeOption)
  }

  loading.value = false
  store.downloadTreeOptions = [comicTreeOption]
  store.downloadDefaultExpandKeys = defaultExpandedKeys
  store.downloadDefaultCheckedKeys = defaultCheckedKeys
}

function onKeyEnterDown() {
  onSearch()
}

function isNumeric(value: string) {
  return !isNaN(Number(value))
}

function extractComicIdFromInput(): string | null {
  const input = searchInput.value.trim()
  if (isNumeric(input)) {
    return input
  }

  const regex = /\/comic\/(\d+)\//
  const match = input.match(regex)
  if (match && match[1]) {
    return match[1]
  }
  return null
}

</script>

<template>
  <n-button text tag="a" href="https://www.manhuagui.com/" target="_blank" type="primary">
    漫画柜
  </n-button>
  <n-input class="search-input" v-model:value="searchInput" placeholder="漫画ID或漫画链接" clearable
           @keydown.enter="onKeyEnterDown"
  />
  <n-popover trigger="hover">
    <template #trigger>
      <n-button class="search-button" @click="onSearch" type="primary" :loading="loading" :disabled="disabled"
                secondary>搜索
        <template #icon>
          <n-icon>
            <search-icon/>
          </n-icon>
        </template>
      </n-button>
    </template>
    <span>直接使用[回车键]也能搜索</span>
  </n-popover>
</template>

<style scoped>
.search-input {
  flex: 4
}

.search-button {
  flex: 1
}
</style>
