<script setup lang="ts">
import {RefreshOutline as RefreshIcon} from "@vicons/ionicons5";
import {defineModel, ref, watch} from "vue";
import {useDownloaderStore} from "../../stores/downloader";
import {types} from "../../../wailsjs/go/models";
import {ScanCacheDir} from "../../../wailsjs/go/api/ExportApi";
import {TreeOption, useNotification} from "naive-ui";
import {ExportStatus} from "../../constants/export-constant";

const notification = useNotification()
const store = useDownloaderStore()

const props = defineProps<{
  refreshDisabled: boolean,
}>()

const exportTreeOptions = defineModel<TreeOption[]>("exportTreeOptions", {required: true});
const exportDefaultExpandKeys = defineModel<string[]>("exportDefaultExpandKeys", {required: true});
const exportDefaultCheckedKeys = defineModel<string[]>("exportDefaultCheckedKeys", {required: true});

const loading = ref<boolean>(false)
watch(() => store.cacheDirectory, onRefresh)


async function buildOptionTree(node: types.TreeNode) {
  const nodeOption: TreeOption = {
    key: node.key,
    label: node.label,
    isLeaf: node.isLeaf,
    disabled: node.disabled,
    children: []
  }
  if (node.defaultChecked) {
    exportDefaultCheckedKeys.value.push(node.key)
    nodeOption.suffix = () => ExportStatus.COMPLETED
  }
  if (node.defaultExpand) {
    exportDefaultExpandKeys.value.push(node.key)
  }

  for (const child of node.children) {
    const childOption = await buildOptionTree(child)
    nodeOption.children?.push(childOption);
  }

  return nodeOption
}

async function onRefresh() {
  try {
    loading.value = true

    const response = await ScanCacheDir(store.cacheDirectory, store.exportDirectory, store.exportTreeMaxDepth)
    if (response.code != 0) {
      notification.create({type: "error", title: "扫描缓存目录失败", content: response.msg})
      return
    }

    const roots: types.TreeNode[] = response.data
    // 清空原有的数据
    const options: TreeOption[] = []
    exportDefaultCheckedKeys.value.length = 0
    exportDefaultExpandKeys.value.length = 0
    for (const root of roots) {
      const rootOption = await buildOptionTree(root)
      options.push(rootOption)
    }

    exportTreeOptions.value = options
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
  <div class="flex flex-col gap-y-2">
    <n-button class="flex-1"
              @click="onRefresh"
              type="primary"
              secondary
              :loading="loading"
              :disabled="refreshDisabled">
      重新扫描
      <template #icon>
        <n-icon>
          <refresh-icon/>
        </n-icon>
      </template>
    </n-button>

  </div>
</template>
