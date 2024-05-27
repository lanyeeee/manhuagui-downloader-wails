<script lang="ts" setup>
import {NTree, TreeInst, TreeOption, TreeOverrideNodeClickBehavior} from "naive-ui"
import {computed, defineModel, ref, watch} from "vue"


const exportTreeInst = defineModel<TreeInst | null>("exportTreeInst", {required: true});
const exportTreeOptions = defineModel<TreeOption[]>("exportTreeOptions", {required: true});
const exportDefaultExpandKeys = defineModel<string[]>("exportDefaultExpandKeys", {required: true});
const exportDefaultCheckedKeys = defineModel<string[]>("exportDefaultCheckedKeys", {required: true});

// TODO: 换成<"empty", "tree">的联合类型
const showInfo = computed<boolean>(() => exportTreeOptions.value.length === 0)
const treeRef = ref<TreeInst | null>(null)
watch(treeRef, () => {
  exportTreeInst.value = treeRef.value
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
        :data="exportTreeOptions"
        :default-expanded-keys="exportDefaultExpandKeys"
        :default-checked-keys="exportDefaultCheckedKeys"
        :override-default-node-click-behavior="treeNodeClickBehaviour"
        ref="treeRef"
    />
  </div>
</template>
