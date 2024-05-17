import {defineStore} from "pinia"
import {TreeInst, TreeOption} from "naive-ui";
import * as path from "../../wailsjs/go/api/PathApi";
import {GetCpuNum} from "../../wailsjs/go/api/UtilsApi";

export const useDownloaderStore = defineStore('downloader', {
    state: () => ({
        downloadTreeInst: null as TreeInst | null,
        downloadTreeOptions: [] as TreeOption[],
        downloadDefaultExpandKeys: [] as string[],
        downloadDefaultCheckedKeys: [] as string[],

        exportTreeInst: null as TreeInst | null,
        exportTreeOptions: [] as TreeOption[],
        exportDefaultExpandKeys: [] as string[],
        exportDefaultCheckedKeys: [] as string[],

        searchDisabled: false,
        refreshDisabled: false,

        proxyUrl: "http://127.0.0.1:7890",
        downloadConcurrentCount: 3,
        exportConcurrentCount: 1,
        cacheDirectory: "",
        exportDirectory: "",
        downloadInterval: 10,
        exportTreeMaxDepth: 3,
    }),
    getters: {
        checkedDownloadTreeOptions(state): (TreeOption | null)[] | undefined {
            return state.downloadTreeInst?.getCheckedData().options
        },
        checkedExportTreeOptions(state): (TreeOption | null)[] | undefined {
            return state.exportTreeInst?.getCheckedData().options
        },
    },
    actions: {
        async init() {
            try {
                const userDownloadPath = await path.UserDownloadPath()
                const cacheDirectory = await path.Join([userDownloadPath, "漫画缓存"])
                if (!await path.PathExists(cacheDirectory)) {
                    await path.MkDirAll(cacheDirectory)
                }
                this.cacheDirectory = await path.Join([userDownloadPath, "漫画缓存"])

                const exportDirectory = await path.Join([userDownloadPath, "漫画导出"])
                if (!await path.PathExists(exportDirectory)) {
                    await path.MkDirAll(exportDirectory)
                }
                this.exportDirectory = await path.Join([userDownloadPath, "漫画导出"])

                this.exportConcurrentCount = await GetCpuNum() / 2
            } catch (e) {
                console.error(e)
            }
        },
    },
})
