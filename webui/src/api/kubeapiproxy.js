import {get, post} from "../lib/fetch";

export function GetNamespace(params) {
    return get('/kubeapiproxy/namespace', params)
}

export function GetPods(params) {
    return get('/kubeapiproxy/pods', params)
}

export function Upload(data, params, headers) {
    return post('/kubeapiproxy/upload', data, params, headers)
}

export function MultiUpload(data, params, headers) {
    return post('/kubeapiproxy/multiupload', data, params, headers)
}

export function Download(params) {
    return get('/kubeapiproxy/download', params, null, 'blob')
}

export function CTerminal(params) {
    return get('/kubeapiproxy/terminal', params)
}

export function Exec(params) {
    return get('/kubeapiproxy/exec', params)
}