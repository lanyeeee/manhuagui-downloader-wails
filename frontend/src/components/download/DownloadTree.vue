<script lang="ts" setup>
import {NTree, TreeInst, TreeOverrideNodeClickBehavior} from "naive-ui"
import {computed, ref, watch} from "vue"
import {useDownloaderStore} from "../../stores/downloader"

const store = useDownloaderStore()

const showInfo = computed<boolean>(() => store.downloadTreeOptions.length === 0)
const treeRef = ref<TreeInst | null>(null)
watch(treeRef, () => {
  store.downloadTreeInst = treeRef.value
})

const treeNodeClickBehaviour: TreeOverrideNodeClickBehavior = ({option}) => {
  if (option.children) {
    return "toggleExpand"
  }
  return "toggleCheck"
}

</script>

<template>
  <n-scrollbar style="height: 60vh">
    <n-result v-if="showInfo" title="在搜索框中输入漫画链接"/>
    <n-tree
        v-if="!showInfo"
        block-line
        show-line
        cascade
        checkable
        :data="store.downloadTreeOptions"
        :default-expanded-keys="store.downloadDefaultExpandKeys"
        :default-checked-keys="store.downloadDefaultCheckedKeys"
        :override-default-node-click-behavior="treeNodeClickBehaviour"
        ref="treeRef"
        style="text-align: left;"
    />
  </n-scrollbar>
</template>
