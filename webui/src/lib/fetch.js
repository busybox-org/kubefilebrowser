import axios from 'axios'
import qs from 'qs'
import jsd from 'js-file-download'
import Vue from 'vue'
import i18n from '../lang'
import Code from './code.js'
import {showFullScreenLoading, tryHideFullScreenLoading} from "./loading";

let API_URL = '/api'
let CancelToken = axios.CancelToken

Vue.prototype.$CancelAjaxRequet = function() {}
Vue.prototype.$IsCancel = axios.isCancel

const service = axios.create({
    baseURL: API_URL + '/',
    timeout: 600000,
    withCredentials: true,
})

service.interceptors.request.use(config => {
    config.headers.Accept = "application/json, text/plain, */*";
    showFullScreenLoading()
    return config
}, error => {
    Promise.reject(error)
})

service.interceptors.response.use(response => {
    tryHideFullScreenLoading()
    let res = response.data
    if (!res) {
        res = {
            code: -1,
            message: i18n.t("network_error"),
        }
    }
    // 文件下载
    if (response.headers["content-type"].indexOf('application/octet-stream') !== -1) {
        let fileName = response.headers['content-disposition'].split('=')[1]
        jsd(res, fileName)
        return
    }
    // json数据处理
    if (res.code !== 0) {
        switch (res.code) {
            case Code.CODE_ERR_NETWORK:
            case Code.CODE_ERR_APP:
                Vue.prototype.$message.error(res.message)
                break
        }
        return Promise.reject(res)
    }
    return res.data
}, error => {
    tryHideFullScreenLoading()
    if (!axios.isCancel(error)) {
        let res = {
            code: -1,
            message: error.message ? error.message : i18n.t("unknown_error"),
        }
        Vue.prototype.$message.error(res.message)
        return Promise.reject(res)
    }
    return Promise.reject(error)
})

export function post(url, data, params, headers) {
    if (!params) {
        params = {}
    }
    params._t = new Date().getTime()
    let config = {
        method: 'post',
        url: url,
        params,
        paramsSerializer: params => {
            return qs.stringify(params, { indices: false})
        }
    }
    if (data) {
        if (headers && headers['Content-Type'] === 'multipart/form-data') {
            config.data = data
        } else {
            config.data = qs.stringify(data, { indices: false })
        }
    }
    if (headers) {
        config.headers = headers
    }

    config.cancelToken = new CancelToken(function(cancel) {
        Vue.prototype.$CancelAjaxRequet = function() {
            cancel()
        }
    })

    return service(config)
}

export function get(url, params, headers, responseType) {
    if (!params) {
        params = {}
    }
    params._t = new Date().getTime()
    // params = qs.stringify(params, {arrayFormat: 'repeat'})
    let config = {
        method: 'get',
        url: url,
        params,
        paramsSerializer: params => {
            return qs.stringify(params, { indices: false})
        }
    }
    if (responseType) {
        config.responseType = responseType
    }
    if (headers) {
        config.headers = headers
    }

    config.cancelToken = new CancelToken(function(cancel) {
        Vue.prototype.$CancelAjaxRequet = function() {
            cancel()
        }
    })

    return service(config)
}

export default service;