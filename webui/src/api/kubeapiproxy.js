import {get, post} from "../lib/fetch";

export function GetNamespace(params) {
    return get('/kubeapiproxy/namespace', params)
}

export function GetPods(params) {
    return get('/kubeapiproxy/pods', params)
}

export function GetPvcs(params) {
    return get('/kubeapiproxy/pvcs', params)
}

export function Upload(data, params, headers) {
    return post('/kubeapiproxy/upload', data, params, headers)
}

export function MultiUpload(data, urlParams, headers) {
    return post(`/kubeapiproxy/multiupload?${urlParams}`, data, null, headers)
}

export function UploadPVC(data, params, headers) {
    return post('/kubeapiproxy/uploadpvc', data, params, headers)
}

export function Download(params) {
    return get('/kubeapiproxy/download', params, null, 'blob')
}

export function BulkDownload(urlParams) {
    return get(`/kubeapiproxy/download?${urlParams}`, null, null, 'blob')
}

export function CTerminal(params) {
    return get('/kubeapiproxy/terminal', params)
}

export function Exec(params) {
    return get('/kubeapiproxy/exec', params)
}