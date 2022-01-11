import {get} from '../lib/fetch.js'
import {post} from '../lib/fetch.js'

export function FileBrowserList(params) {
    return get('/filebrowser/list', params)
}

export function FileBrowserOpen(params) {
    return get('/filebrowser/open', params)
}

export function FileBrowserCreateFile(data, params) {
    return post('/filebrowser/createfile', data, params, {"Content-Type":"multipart/form-data"})
}

export function FileBrowserCreateDir(params) {
    return post('/filebrowser/createdir', null, params)
}

export function FileBrowserRename(params) {
    return post('/filebrowser/rename', null, params)
}

export function FileBrowserRemove(params) {
    return post('/filebrowser/remove', null, params)
}