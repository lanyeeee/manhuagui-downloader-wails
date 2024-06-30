import {defineStore} from "pinia";
import * as path from "../../wailsjs/go/api/PathApi";
import {GetCpuNum, GetUserDownloadPath, GetUserProxy} from "../../wailsjs/go/api/UtilsApi";

export const useDownloaderStore = defineStore("downloader", {
    state: () => ({
        proxyUrl: "",
        downloadConcurrentCount: 3,
        exportConcurrentCount: 1,
        cacheDirectory: "",
        exportDirectory: "",
        downloadInterval: 10,
        exportTreeMaxDepth: 3,
    }),
    getters: {},
    actions: {
        async init() {
            try {
                this.proxyUrl = await GetUserProxy();

                const userDownloadPath = await GetUserDownloadPath();
                const exportDirectory = await path.Join([userDownloadPath, "漫画导出"]);
                if (!await path.PathExists(exportDirectory)) {
                    await path.MkDirAll(exportDirectory);
                }
                const cacheDirectory = await path.Join([userDownloadPath, "漫画缓存"]);
                if (!await path.PathExists(cacheDirectory)) {
                    await path.MkDirAll(cacheDirectory);
                }
                this.exportDirectory = await path.Join([userDownloadPath, "漫画导出"]);
                this.cacheDirectory = await path.Join([userDownloadPath, "漫画缓存"]);

                this.exportConcurrentCount = await GetCpuNum() / 2;
            } catch (e) {
                console.error(e);
            }
        },
    },
});
