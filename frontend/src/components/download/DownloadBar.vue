<script setup lang="ts">
import {NButton, NSpin, TreeOption, useNotification} from "naive-ui"
import {computed, h, onMounted, ref} from "vue"
import {EventsOff, EventsOn} from "../../../wailsjs/runtime"
import {useDownloaderStore} from "../../stores/downloader"
import {DownloadOutline as DownloadIcon} from "@vicons/ionicons5"
import {DownloadStatus} from "../../constants/download-constant"
import {DownloadChapter} from "../../../wailsjs/go/api/DownloadApi"

const store = useDownloaderStore()
const notification = useNotification()

const downloadButtonDisabled = computed<boolean>(() => store.downloadTreeOptions.length === 0)
const downloadButtonLoading = ref<boolean>(false)
const chapterProgressIndicator = ref<string>("")
const chapterProgressPercentage = ref<number>(0)
const overallProgressIndicator = ref<string>("")
const overallProgressPercentage = ref<number>(0)

onMounted(() => {
  EventsOn("download", (speed: string, percentage: number) => {
    chapterProgressIndicator.value = speed
    chapterProgressPercentage.value = percentage
  })
})

async function onDownload() {
  // 满足不为null，是叶子节点，没有disabled的才是要下载的章节
  const optionsToDownload = store.checkedDownloadTreeOptions?.filter(option => option !== null && option.isLeaf && !option.disabled)
  if (optionsToDownload === undefined || optionsToDownload.length === 0) {
    notification.create({type: "error", title: "下载失败", content: "请选择要下载的章节", duration: 2000,})
    return
  }
  try {
    store.searchDisabled = true
    downloadButtonLoading.value = true
    await downloadOptions(optionsToDownload)
  } finally {
    store.searchDisabled = false
    downloadButtonLoading.value = false
  }
}

async function downloadOptions(optionsToDownload: (TreeOption | null)[]) {
  // 先把所有的option都disable掉，防止用户在下载过程中进行其他操作
  optionsToDownload.forEach(option => {
    if (option !== null) {
      option.disabled = true
      option.suffix = () => DownloadStatus.WAITING
    }
  })
  overallProgressIndicator.value = `(0/${optionsToDownload.length})`

  for (const [i, option] of optionsToDownload.entries()) {
    if (option === null) {
      continue
    }

    chapterProgressPercentage.value = 0
    chapterProgressIndicator.value = `即将下载 ${option.label}`

    option.prefix = () => h(NSpin, {size: 15})
    option.suffix = () => DownloadStatus.DOWNLOADING

    const key = JSON.parse(option.key as string) as { href: string, saveDir: string }
    const chapterUrl = "https://www.manhuagui.com" + key.href
    const saveDirectory = key.saveDir
    const concurrentCount = store.downloadConcurrentCount
    const proxyUrl = store.proxyUrl
    EventsOn("download", (msg: string, percentage: number) => {
      chapterProgressIndicator.value = msg
      chapterProgressPercentage.value = percentage
    })
    try {
      const response = await DownloadChapter(chapterUrl, saveDirectory, concurrentCount, proxyUrl)
      // 如果下载失败
      if (response.code != 0) {
        console.error(response.msg)
        notification.create({type: "error", title: "下载失败", meta: response.msg,})
        option.disabled = false
        option.prefix = undefined
        option.suffix = () => DownloadStatus.FAILED
        return
      }
      // 下载成功
      option.prefix = undefined
      option.suffix = () => DownloadStatus.COMPLETED
    } finally {
      EventsOff("download")
      // 更新总体进度
      overallProgressPercentage.value = ((i + 1) / optionsToDownload.length) * 100
      // 更新总体进度提示
      overallProgressIndicator.value = `(${i + 1}/${optionsToDownload.length})`
      // 如果不是最后一个章节，在下载完成后等待store.downloadInterval秒
      if (i !== optionsToDownload.length - 1) {
        //每秒更新一次剩余时间
        for (let i = store.downloadInterval; i > 0; i--) {
          chapterProgressIndicator.value = `${option.label} 下载完成，将在${i}秒后继续下载(防止被ban IP)`
          await new Promise(resolve => setTimeout(resolve, 1000))
        }
      }
    }
  }
}

</script>

<template>
  <div class="flex flex-col gap-y-3">
    <div class="flex">
      <n-text>章节进度：</n-text>
      <n-progress type="line"
                  :percentage="chapterProgressPercentage"
                  :height="25"
                  indicator-placement="inside"
                  indicator-text-color="black"
                  class="flex-1"
      >{{ chapterProgressIndicator }}
      </n-progress>
    </div>
    <div class="flex">
      <n-text>总体进度：</n-text>
      <n-progress type="line"
                  :percentage="overallProgressPercentage"
                  :height="25"
                  indicator-placement="inside"
                  indicator-text-color="black"
                  class="flex-1"
      >{{ overallProgressIndicator }}
      </n-progress>
    </div>
    <template v-if="downloadButtonDisabled">
      <n-popover trigger="hover">
        <template #trigger>
          <n-button @click="onDownload" type="primary"
                    :disabled="downloadButtonDisabled"
                    :loading="downloadButtonLoading">
            下载
            <template #icon>
              <n-icon>
                <download-icon/>
              </n-icon>
            </template>
          </n-button>
        </template>
        <span>请先搜索漫画</span>
      </n-popover>
    </template>
    <template v-else>
      <n-button @click="onDownload" type="primary"
                :disabled="downloadButtonDisabled"
                :loading="downloadButtonLoading">
        下载
        <template #icon>
          <n-icon>
            <download-icon/>
          </n-icon>
        </template>
      </n-button>
    </template>
  </div>
</template>

