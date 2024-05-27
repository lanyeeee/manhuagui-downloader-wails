<script lang="ts" setup>
import {NTree, TreeInst, TreeOverrideNodeClickBehavior} from "naive-ui"
import {computed, ref, watch} from "vue"
import {useDownloaderStore} from "../../stores/downloader"


const store = useDownloaderStore()

const showInfo = computed<boolean>(() => store.exportTreeOptions.length === 0)
const treeRef = ref<TreeInst | null>(null)
watch(treeRef, () => {
  store.exportTreeInst = treeRef.value
})

const treeNodeClickBehaviour: TreeOverrideNodeClickBehavior = ({option}) => {
  if (option.children?.length === 0) {
    return "toggleCheck"
  }
  return "toggleExpand"
}

</script>

<template>
  <div class="h-full overflow-hidden">
    <n-result v-if="showInfo" title="缓存目录为空"/>
    <n-tree
        class="text-align-left"
        v-if="!showInfo"
        virtual-scroll
        block-line
        show-line
        checkable
        cascade
        :data="store.exportTreeOptions"
        :default-expanded-keys="store.exportDefaultExpandKeys"
        :default-checked-keys="store.exportDefaultCheckedKeys"
        :override-default-node-click-behavior="treeNodeClickBehaviour"
        ref="treeRef"
    />
  </div>
</template>
