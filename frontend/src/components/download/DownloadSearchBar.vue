<script setup lang="ts">
import {computed, ref} from "vue"
import {TreeOption, useNotification} from "naive-ui"
import {useDownloaderStore} from "../../stores/downloader"
import {SearchOutline as SearchIcon} from "@vicons/ionicons5"
import {SearchComicInfo} from "../../../wailsjs/go/api/DownloadApi"
import {types} from "../../../wailsjs/go/models"
import {DownloadStatus} from "../../constants/download-constant";


const store = useDownloaderStore()
const notification = useNotification()

const searchInput = ref<string>("")
const loading = ref<boolean>(false)
const disabled = computed<boolean>(() => store.searchDisabled)

function buildOptionTree(node: types.TreeNode): TreeOption {
  const nodeOption: TreeOption = {
    key: node.key,
    label: node.label,
    isLeaf: node.isLeaf,
    disabled: node.disabled,
    children: []
  }

  if (node.defaultChecked) {
    store.downloadDefaultCheckedKeys.push(node.key)
    nodeOption.suffix = () => DownloadStatus.COMPLETED
  }
  if (node.defaultExpand) {
    store.downloadDefaultExpandKeys.push(node.key)
  }

  for (const child of node.children) {
    const childOption = buildOptionTree(child)
    nodeOption.children?.push(childOption);
  }

  return nodeOption
}

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
    const response: types.Response = await SearchComicInfo(comicId, store.proxyUrl, store.cacheDirectory)
    if (response.code != 0) {
      notification.create({type: "error", title: "搜索失败", meta: response.msg,})
      return
    }

    const root: types.TreeNode = response.data
    console.log("搜索结果", root)
    const rootOption: TreeOption = buildOptionTree(root)

    store.downloadTreeOptions = [rootOption]

  } finally {
    loading.value = false
  }
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
