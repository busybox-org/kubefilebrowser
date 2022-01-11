import i18n from '../lang'
import {Loading} from "element-ui";
import _ from 'lodash'

let needLoadingRequestCount = 0
let loading


export function startLoading() {
    loading = Loading.service({
        lock: true,
        text: i18n.t('loading'),
        background: 'rgba(0, 0, 0, 0)'
    })
}

export function endLoading() {
    loading.close()
}

const tryCloseLoading = () => {
    if (needLoadingRequestCount === 0) {
        endLoading()
    }
}

export function showFullScreenLoading() {
    if (needLoadingRequestCount === 0) {
        startLoading()
    }
    needLoadingRequestCount++
}

export function tryHideFullScreenLoading() {
    if (needLoadingRequestCount <= 0) return
    needLoadingRequestCount--
    if (needLoadingRequestCount === 0) {
        _.debounce(tryCloseLoading, 100)()  //延迟100ms
    }
}