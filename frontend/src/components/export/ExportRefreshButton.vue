<script setup lang="ts">
import {RefreshOutline as RefreshIcon} from "@vicons/ionicons5";
import {ref, watch} from "vue";
import {useDownloaderStore} from "../../stores/downloader";
import {scan_cache, types} from "../../../wailsjs/go/models";
import {ScanCacheDir} from "../../../wailsjs/go/api/ExportApi";
import {TreeOption, useNotification} from "naive-ui";
import * as path from "../../../wailsjs/go/api/PathApi";
import {ExportStatus} from "../../constants/export-constant";

const notification = useNotification()

const loading = ref<boolean>(false)
const store = useDownloaderStore()

watch(() => store.cacheDirectory, onRefresh)


async function buildOptionTree(node: scan_cache.TreeNode): Promise<TreeOption> {
  const nodeOption: TreeOption = {key: node.key, label: node.label, isLeaf: node.isLeaf,}
  const relativePath: string = await path.GetRelPath(store.cacheDirectory, node.key)
  const pdfPath: string = await path.Join([store.exportDirectory, relativePath + ".pdf"])
  if (await path.PathExists(pdfPath)) {
    nodeOption.suffix = () => ExportStatus.COMPLETED
    if (node.isLeaf) {
      nodeOption.disabled = true
      store.exportDefaultCheckedKeys.push(node.key)
    }
  }

  for (const child of node.children) {
    nodeOption.children ??= []
    const childOption = await buildOptionTree(child)
    nodeOption.children.push(childOption);
  }

  return nodeOption
}

async function onRefresh() {
  try {
    loading.value = true

    const response: types.Response = await ScanCacheDir(store.cacheDirectory, store.exportTreeMaxDepth)
    if (response.code != 0) {
      notification.create({type: "error", title: "扫描缓存目录失败", content: response.msg})
      return
    }

    const roots: scan_cache.TreeNode[] = response.data
    console.log(roots)
    // 清空原有的数据
    store.exportTreeOptions.length = 0
    store.exportDefaultExpandKeys.length = 0
    for (const root of roots) {
      const rootOption: TreeOption = await buildOptionTree(root)
      store.exportTreeOptions.push(rootOption)
      // 设置默认展开的节点
      store.exportDefaultExpandKeys.push(rootOption.key as string)
      rootOption.children?.forEach(child => {
        store.exportDefaultExpandKeys.push(child.key as string)
      })
    }
  } catch (e) {
    console.error(e)
    if (typeof e === "string") {
      notification.create({type: "error", title: "扫描缓存目录失败", content: "异常", meta: e})
    } else if (e instanceof Error) {
      notification.create({type: "error", title: "扫描缓存目录失败", content: "异常", meta: e.message})
    } else {
      notification.create({type: "error", title: "扫描缓存目录失败", content: "异常", meta: "未知异常"})
    }
  } finally {
    loading.value = false
  }

}
</script>

<template>
  <n-flex>
    <div style="flex: 3"/>
    <n-button @click="onRefresh" type="primary" secondary :loading="loading" :disabled="store.refreshDisabled"
              style="flex: 1">
      重新扫描缓存目录
      <template #icon>
        <n-icon>
          <refresh-icon/>
        </n-icon>
      </template>
    </n-button>
  </n-flex>
</template>
