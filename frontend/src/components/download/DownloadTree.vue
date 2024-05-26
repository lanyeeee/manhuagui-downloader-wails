<script lang="ts" setup>
import {NTree, TreeInst, TreeOption, TreeOverrideNodeClickBehavior} from "naive-ui"
import {defineModel} from "vue"

const downloadTreeInst = defineModel<TreeInst | null>("downloadTreeInst", {required: true});
const downloadTreeOptions = defineModel<TreeOption[]>("downloadTreeOptions", {required: true});
const downloadDefaultExpandKeys = defineModel<string[]>("downloadDefaultExpandKeys", {required: true});
const downloadDefaultCheckedKeys = defineModel<string[]>("downloadDefaultCheckedKeys", {required: true});

const treeNodeClickBehaviour: TreeOverrideNodeClickBehavior = ({option}) => {
  if (option.children?.length === 0) {
    return "toggleCheck"
  }
  return "toggleExpand"
}


</script>

<template>
  <n-tree
      class="text-align-left overflow-hidden"
      virtual-scroll
      block-line
      show-line
      cascade
      checkable
      :data="downloadTreeOptions"
      :default-expanded-keys="downloadDefaultExpandKeys"
      :default-checked-keys="downloadDefaultCheckedKeys"
      :override-default-node-click-behavior="treeNodeClickBehaviour"
      ref="downloadTreeInst"
  />
</template>
