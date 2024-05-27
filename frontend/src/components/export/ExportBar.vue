<script setup lang="ts">

import {NButton, NSpin, TreeOption, useNotification} from "naive-ui"
import {computed, h, ref} from "vue"
import {useDownloaderStore} from "../../stores/downloader"
import {BookOutline as ExportIcon} from "@vicons/ionicons5"
import {ExportStatus} from "../../constants/export-constant"
import {export_pdf} from "../../../wailsjs/go/models"
import {BrowserOpenURL, EventsOff, EventsOn} from "../../../wailsjs/runtime"
import * as path from "../../../wailsjs/go/api/PathApi"
import {CreatePdfs, MergePdfs} from "../../../wailsjs/go/api/ExportApi"
import ExportDirectoryInput from "../settings/ExportDirectoryInput.vue";

const store = useDownloaderStore()
const notification = useNotification()

const disabled = computed<boolean>(() => store.exportTreeOptions.length === 0)
const buttonLoading = ref<boolean>(false)
const createProgressIndicator = ref<string>("")
const createProgressPercentage = ref<number>(0)
const mergeProgressIndicator = ref<string>("")
const mergeProgressPercentage = ref<number>(0)

// TODO: 修改获取要被导出的option的逻辑，目前的实现没法配合列表动态调整
async function onExport() {
  const leafOptionsToExport = store.checkedExportTreeOptions
      ?.filter(option =>
          option !== null &&
          !option.disabled &&
          option.isLeaf
      ) as TreeOption[] ?? []

  const nonLeafOptionsToExport = (store.checkedExportTreeOptions
      ?.filter(option =>
          option !== null &&
          !option.disabled &&
          !option.isLeaf
      ) as TreeOption[] ?? [])
      .sort((a: TreeOption, b: TreeOption) => {
        // 按照路径深度排序，深度深的排在前面
        const aSlashCount = (a.key as string).split("/").length
        const bSlashCount = (b.key as string).split("/").length
        return bSlashCount - aSlashCount
      })

  // 如果leafOptionsToExport和nonLeafOptionsToExport都为空，报错
  if (leafOptionsToExport.length === 0 && nonLeafOptionsToExport.length === 0) {
    notification.create({type: "error", title: "导出失败", content: "请选择要导出的章节", duration: 2000,})
    return
  }

  try {
    buttonLoading.value = true
    store.refreshDisabled = true
    createProgressPercentage.value = 0
    createProgressIndicator.value = `(${0}/${leafOptionsToExport.length})`
    mergeProgressPercentage.value = 0
    mergeProgressIndicator.value = `(${0}/${nonLeafOptionsToExport.length})`

    // 先把所有的option都disable掉，防止用户在导出过程中进行其他操作
    leafOptionsToExport.forEach(option => {
      option.disabled = true
      option.suffix = () => ExportStatus.WAITING
    })
    nonLeafOptionsToExport.forEach(option => {
      option.disabled = true
      option.suffix = () => ExportStatus.WAITING
    })

    await exportLeafOptions(leafOptionsToExport)
    await exportNonLeafOptions(nonLeafOptionsToExport)
  } finally {
    buttonLoading.value = false
    store.refreshDisabled = false
  }
}

async function exportLeafOptions(optionsToExport: TreeOption []) {
  const request = new export_pdf.CreatePdfsRequest({
    tasks: [],
    concurrentCount: store.exportConcurrentCount
  })

  // 填充request的tasks，并且给每个option添加prefix和suffix
  for (const option of optionsToExport) {
    option.prefix = () => h(NSpin, {size: 15})
    option.suffix = () => ExportStatus.CREATING

    const imageDirectory = option.key as string
    // 获取相对路径
    const relativePath = await path.GetRelPath(store.cacheDirectory, imageDirectory)
    const outputPath = await path.Join([store.exportDirectory, relativePath + ".pdf"])
    const task: export_pdf.CreatePdfTask = {
      imgDir: option.key as string,
      outputPath: outputPath,
      optionKey: option.key as string,
    }
    request.tasks.push(task)
  }

  EventsOn("create_pdf", (completedOptionKey: string, msg: string, percentage: number) => {
    const completedOption = optionsToExport.find(option => option.key === completedOptionKey) as TreeOption | undefined
    if (completedOption !== undefined) {
      completedOption.suffix = () => ExportStatus.COMPLETED
      completedOption.prefix = undefined
    }
    createProgressIndicator.value = msg
    createProgressPercentage.value = percentage
  })
  console.log(request)
  const response = await CreatePdfs(request)
  if (response.code != 0) {
    notification.create({type: "error", title: "导出失败", content: "创建PDF文件失败", meta: response.msg,})
  }
  EventsOff("create_pdf")
}

async function exportNonLeafOptions(optionsToExport: TreeOption[]) {
  for (const [i, option] of optionsToExport.entries()) {
    try {
      option.prefix = () => h(NSpin, {size: 15})
      option.suffix = () => ExportStatus.MERGING

      const relativePath = await path.GetRelPath(store.cacheDirectory, option.key as string)
      const inputDirectory = await path.Join([store.exportDirectory, relativePath])
      const outputPath = await path.Join([store.exportDirectory, relativePath + ".pdf"])

      const response = await MergePdfs(inputDirectory, outputPath)
      if (response.code != 0) {
        notification.create({type: "error", title: "导出失败", content: "合并PDF文件失败", meta: response.msg,})
        option.prefix = undefined
        option.suffix = () => ExportStatus.FAILED
        continue
      }

      option.prefix = undefined
      option.suffix = () => ExportStatus.COMPLETED
    } finally {
      option.disabled = false
      mergeProgressPercentage.value = ((i + 1) / optionsToExport.length) * 100
      mergeProgressIndicator.value = `(${i + 1}/${optionsToExport.length})`
    }
  }
}

async function onOpenExportDirectory() {
  BrowserOpenURL(store.exportDirectory)
}

</script>

<template>
  <div class="flex flex-col gap-y-2">
    <div class="flex">
      <n-text>生成进度：</n-text>
      <n-progress class="flex-1"
                  type="line"
                  :percentage="createProgressPercentage"
                  :height="25"
                  indicator-placement="inside"
                  indicator-text-color="black"
      >{{ createProgressIndicator }}
      </n-progress>
    </div>
    <div class="flex">
      <n-text>合并进度：</n-text>
      <n-progress class="n-progress flex-1"
                  type="line"
                  :percentage="mergeProgressPercentage"
                  :height="25"
                  indicator-placement="inside"
                  indicator-text-color="black"
      >{{ mergeProgressIndicator }}
      </n-progress>
    </div>

    <export-directory-input/>
    <n-button class="n-progress"
              @click="onExport"
              type="primary"
              :disabled="disabled"
              :loading="buttonLoading"
    >开始导出
      <template #icon>
        <n-icon>
          <export-icon/>
        </n-icon>
      </template>
    </n-button>
    <n-button @click="onOpenExportDirectory">打开导出目录</n-button>
  </div>
</template>
