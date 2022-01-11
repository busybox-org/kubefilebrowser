<template>
  <div>
    <el-card shadow="never">
      <div>
        <el-select v-model="namespace" @click.native="getNamespace" @change="selectedNamespace" style="width: 100%" filterable :placeholder="$t('please_select_namespace')">
          <el-option
              v-for="item in namespaces"
              :label="item"
              :value="item"
              :key="item"
          ></el-option>
        </el-select>
      </div>
      <div style="margin-top: 15px" class="table">
        <div style="float: right; text-align: center;margin-bottom: 20px;">
          <el-pagination
              background
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="currentPage"
              :page-sizes="pageSizes"
              :page-size="pageSize"
              layout="total, sizes, prev, pager, next, jumper"
              :total="totalCount">
          </el-pagination>
        </div>
        <!-- div style="float: right;width: 300px" class="search-Box">
          <el-input  :placeholder="$t('keyword_search')" icon="search"  class="search"  v-model="searchKey"/>
        </div -->
        <el-table
            v-tableDrag
            class="app-table"
            size="medium"
            :data="tableData"
            style="width: 100%"
            :default-sort="{prop: 'PodName', order: 'ascending'}"
          >
          <el-table-column prop="PodName" :label="$t('pod')" sortable show-overflow-tooltip min-width="400" fixed="left" :sort-orders="['ascending', 'descending']"></el-table-column>
          <el-table-column :label="$t('state')" sortable show-overflow-tooltip width="100" :sort-orders="['ascending', 'descending']">
            <template slot-scope="scope">
              <el-button disabled v-if="scope.row.State === 'Running' || scope.row.State === 'Succeeded'" type="success"
                         size="mini" plain round>
                {{ scope.row.State }}
              </el-button>
              <el-button disabled v-if="scope.row.State !== 'Running' && scope.row.State !== 'Succeeded'" type="warning"
                         size="mini" plain round>
                {{ scope.row.State }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column label="Pod IP"  sortable width="150" :sort-orders="['ascending', 'descending']">
            <template slot-scope="scope">
              <span v-for="list in scope.row.PodIPs">{{list.ip}}<br></span>
            </template>
          </el-table-column>
          <el-table-column :label="$t('node')" sortable width="300" :sort-orders="['ascending', 'descending']">
            <template slot-scope="scope">
              {{scope.row.NodeName}}<br>
              {{scope.row.HostIP}}
            </template>
          </el-table-column>
          <el-table-column prop="CreateTime" :label="$t('create_time')" sortable fix width="200" :sort-orders="['ascending', 'descending']"></el-table-column>
          <el-table-column prop="RestartCount" :label="$t('restart_count')" sortable width="100" :sort-orders="['ascending', 'descending']"></el-table-column>
          <el-table-column :label="$t('operate')" min-width="150" fixed="right" style="text-align: right" :sort-orders="['ascending', 'descending']">
            <template v-slot:default="{row}">
              <el-dropdown style="margin-left: 10px" :hide-on-click="false" v-if="row.State === 'Running'">
                <el-dropdown-menu></el-dropdown-menu>
                <el-popover placement="left" trigger="hover">
                  <div v-for="c in row.Containers" :key="c.name">
                    <p style="margin: 0">
                      <el-button @click="openTerminal(row, c)" type="text">{{ c.name }}</el-button>
                    </p>
                  </div>
                  <el-dropdown-item slot="reference" icon="el-icon-s-fold">
                    {{ $t("terminal") }}
                    <i class="el-icon-arrow-right"/>
                  </el-dropdown-item>
                </el-popover>
                <el-popover placement="left" trigger="hover">
                  <div v-for="c in row.Containers" :key="c.name">
                    <p style="margin: 0">
                      <el-button @click="openFileBrowser(row, c, '/')" type="text">{{ c.name }}</el-button>
                    </p>
                  </div>
                  <el-dropdown-item slot="reference" icon="el-icon-files">
                    {{ $t("file_browser") }}
                    <i class="el-icon-arrow-right"/>
                  </el-dropdown-item>
                </el-popover>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    <el-dialog
        :visible.sync="dialogTerminalVisible"
        :title="$t('terminal')"
        center
        fullscreen
        :modal="false"
        :destroy-on-close="true"
        @opened="doOpened"
        @close="doClose"
    >
      <div style="margin-top: -25px;">
        <div ref="terminal" />
      </div>
    </el-dialog>
    <el-dialog
        center
        fullscreen
        :title="$t('file_browser')"
        :visible.sync="dialogFileBrowserVisible"
        @close="dialogFileBrowserVisible = false">
      <div style="margin-top: -25px;">
        <el-table-header store="">
          <el-dropdown  type="info" class="avatar-container" trigger="click" style="height: 36px; float: right; margin-bottom: 10px;">
            <div class="avatar-wrapper">
              <el-button style="width: 90px; height: 30px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="success" round class="el-icon-s-tools" size="medium">
                {{ $t('create') }}
                <i class="el-icon-caret-bottom" />
              </el-button>
            </div>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>
                <span style="display:block;" @click="openFileDialog(globalPath, 'create')">{{ $t('create_file') }}</span>
              </el-dropdown-item>
              <el-dropdown-item>
                <span style="display:block;" @click="createDir()">{{ $t('create_dir') }}</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
          <el-dropdown type="success" class="avatar-container" trigger="click" style="height: 36px;float: right;margin-bottom: 10px;">
            <div class="avatar-wrapper">
              <el-button style="width: 90px; height: 30px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="success" round class="el-icon-upload" size="medium">
                {{ $t('upload') }}
                <i class="el-icon-caret-bottom" />
              </el-button>
            </div>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>
                <span class="fake-file-btn">
                  {{ $t('upload_file') }}
                  <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event, globalPath)" name="files" multiple="true">
                </span>
              </el-dropdown-item>
              <el-dropdown-item divided>
                <span class="fake-file-btn">
                  {{ $t('upload_dir') }}
                  <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event, globalPath)" name="files" webkitdirectory mozdirectory accept="*/*">
                </span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
          <el-dropdown type="primary" class="el-upload avatar-container" trigger="click" style="height: 36px;float: right;margin-bottom: 10px;">
            <div class="avatar-wrapper">
              <el-button style="width: 120px; height: 30px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="primary" round class="el-icon-download" size="medium">
                {{ $t('bulk_download') }}
                <i class="el-icon-caret-bottom" />
              </el-button>
            </div>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>
                <span style="display:block;" @click="bulkDownload(bulkPath, 'tar')">TAR{{ $t('download') }}</span>
              </el-dropdown-item>
              <el-dropdown-item divided>
                <span style="display:block;" @click="bulkDownload(bulkPath, 'zip')">ZIP{{ $t('download') }}</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
          <ul>
            <li style="float: left; margin-top: 10px; list-style: none;" v-for="(item) in headerPaths">
              <a style="margin-right: 5px; font-size: 16px" class="el-icon-folder-opened" @click="openFileBrowser(null, null, item.path)">{{item.name}}</a>
            </li>
          </ul>
          &nbsp;&nbsp;&nbsp;&nbsp;
          <span style="float: left;">
            <el-button type="info" style="padding: 3px;margin-top: 8px;" icon="el-icon-refresh" circle @click="openFileBrowser(null,null, path)"></el-button>
          </span>
        </el-table-header>
        <el-table
            v-tableDrag
            id="tableData"
            class="app-table"
            border
            style="width: 100%"
            size="100%"
            :cell-style="{padding:'6px 0'}"
            :data="fileBrowserData"
            @selection-change="handleSelectionChange"
            :default-sort="{prop: 'Name', order: 'ascending'}">
          <el-table-column type="selection" fixed="left"></el-table-column>
          <el-table-column
              fixed="left"
              min-width="80px"
              prop="Name"
              :label="$t('name')"
              sortable
              :sort-orders="['ascending', 'descending']"
          >
            <template slot-scope="scope">
              <div class="el-icon-folder"  v-if="scope.row.IsDir" @click="openFileBrowser(null,null, scope.row.Path)">&nbsp;&nbsp;{{scope.row.Name}}</div>
              <div class="el-icon-files"  v-else @click="openFileDialog(scope.row.Path, 'open')">{{scope.row.Name}}</div>
            </template>
          </el-table-column>
          <el-table-column
              prop="Size"
              :label="$t('size')"
              sortable
              :sort-orders="['ascending', 'descending']"
          >
          </el-table-column>
          <el-table-column
              prop="Mode"
              :label="$t('mode')"
          >
          </el-table-column>
          <el-table-column
              prop="ModTime"
              :label="$t('mod_time')"
              sortable
              :sort-orders="['ascending', 'descending']"
          >
          </el-table-column>
          <el-table-column
              fixed="right"
              prop="Download"
              :label="$t('operate')" align="center"
          >
            <template slot-scope="scope">
              <el-dropdown type="info" class="avatar-container" trigger="click" style="height: 36px;font-size: 9px">
                <div class="avatar-wrapper">
                  <el-button style="width: 90px; height: 30px; margin-top: 4px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="success" round class="el-icon-s-tools" size="medium">
                    {{ $t('operate') }}
                    <i class="el-icon-caret-bottom" />
                  </el-button>
                </div>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item v-if="!scope.row.IsDir">
                    <span class="fake-file-btn" @click="openFileDialog(scope.row.Path, 'open')">{{ $t('change') }}</span>
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <span class="fake-file-btn" @click="openRenameDialog(scope.row.Name)">{{ $t('rename') }}</span>
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <span class="fake-file-btn" @click="removeFileOrDir(scope.row.Path)">{{ $t('remove') }}</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
              <span>
                &nbsp;&nbsp;
              </span>
              <el-dropdown v-if="scope.row.IsDir" type="success" class="avatar-container" trigger="click" style="height: 36px;font-size: 9px">
                <div class="avatar-wrapper">
                  <el-button style="width: 90px; height: 30px; margin-top: 4px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="success" round class="el-icon-upload" size="medium">
                    {{ $t('upload') }}
                    <i class="el-icon-caret-bottom" />
                  </el-button>
                </div>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item>
                    <span class="fake-file-btn">
                      {{ $t('upload_file') }}
                      <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event, scope.row.Path)" name="files" multiple="true">
                    </span>
                  </el-dropdown-item>
                  <el-dropdown-item divided>
                    <span class="fake-file-btn">
                      {{ $t('upload_dir') }}
                      <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event, scope.row.Path)" name="files" webkitdirectory mozdirectory accept="*/*">
                    </span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
              <span>
                &nbsp;&nbsp;
              </span>
              <el-dropdown type="primary" class="avatar-container" trigger="click" style="height: 36px;font-size: 9px">
                <div class="avatar-wrapper">
                  <el-button style="width: 90px; height: 30px; margin-top: 4px; margin-right: 6px; padding-top: 7px; padding-left: 14px;" type="primary" round class="el-icon-download" size="medium">
                    {{ $t('download') }}
                    <i class="el-icon-caret-bottom" />
                  </el-button>
                </div>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item>
                    <span style="display:block;" @click="download(scope.row.Path, 'tar')">TAR{{ $t('download') }}</span>
                  </el-dropdown-item>
                  <el-dropdown-item divided>
                    <span style="display:block;" @click="download(scope.row.Path, 'zip')">ZIP{{ $t('download') }}</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>
    <el-dialog
        :append-to-body="true"
        :with-header="false"
        width="80%"
        :title="$t('file')"
        :visible.sync="dialogFileVisible"
        @close="dialogFileVisible = false">
      <div style="margin-top: -45px">
        <span style="display:block; float: left;font-size: 25px;margin-top: 20px; margin-left: 12px">{{createForPath}}</span>
        <el-input v-if="isNewFile" v-model="createName" size="small" style="margin-top: 17px; margin-left: 6px;float: left; width: auto;" autocomplete="off" :placeholder="$t('please_input_name')"></el-input>
        <el-button @click.native="saveFile" style="float: right;margin-right: 12px;margin-top: 12px">{{ $t('enter') }}</el-button>
        <el-input v-model="fileContent" rows="15" type="textarea" style="margin-top: 15px" :placeholder="$t('please_input_content')"></el-input>
        <!--        style="margin-top: 15px;margin-bottom: 60px;padding-bottom: 70px;height: 600px;"-->
        <!--        <quill-editor-->
        <!--            style="margin-top: 15px;height: 100%"-->
        <!--            v-model="fileContent"-->
        <!--            :options="editorOption">-->
        <!--        </quill-editor>-->
      </div>
    </el-dialog>
  </div>
</template>


<style>
.fake-file-btn {
}
.fake-file-btn:active {
  box-shadow: 0 1px 5px 1px rgba(0, 255, 255, 0.3) inset;
}
.fake-file-btn input[type=file] {
  position: absolute;
  font-size: 100px;
  right: 0;
  top: 0;
  opacity: 0;
  filter: alpha(opacity=0);
  cursor: pointer
}
</style>

<script>
import {Download, GetNamespace, GetPods, Upload,} from '../api/kubeapiproxy'
import {
  FileBrowserCreateDir,
  FileBrowserCreateFile,
  FileBrowserList,
  FileBrowserOpen,
  FileBrowserRemove,
  FileBrowserRename,
} from "../api/filebrowser";
import {Terminal} from 'xterm'
import * as fit from 'xterm/lib/addons/fit/fit'
import {Base64} from 'js-base64'
import * as webLinks from 'xterm/lib/addons/webLinks/webLinks'
import * as search from 'xterm/lib/addons/search/search'
import 'xterm/lib/addons/fullscreen/fullscreen.css'
import 'xterm/dist/xterm.css'

const toolbarOptions = [
  ["bold", "italic", "underline", "strike"], // 加粗 斜体 下划线 删除线
  ["blockquote", "code-block"], // 引用  代码块
  [{ header: 1 }, { header: 2 }], // 1、2 级标题
  [{ list: "ordered" }, { list: "bullet" }], // 有序、无序列表
  [{ script: "sub" }, { script: "super" }], // 上标/下标
  [{ indent: "-1" }, { indent: "+1" }], // 缩进
  // [{'direction': 'rtl'}],                         // 文本方向
  [{ size: ["small", false, "large", "huge"] }], // 字体大小
  [{ header: [1, 2, 3, 4, 5, 6, false] }], // 标题
  [{ color: [] }, { background: [] }], // 字体颜色、字体背景颜色
  [{ font: [] }], // 字体种类
  [{ align: [] }], // 对齐方式
  ["clean"], // 清除文本格式
  ["link"] // 链接、图片
];

const defaultTheme = {
  foreground: '#ffffff', // 字体
  background: '#1b212f', // 背景色
  cursor: '#ffffff', // 设置光标
  selection: 'rgba(255, 255, 255, 0.3)',
  black: '#000000',
  brightBlack: '#808080',
  red: '#ce2f2b',
  brightRed: '#f44a47',
  green: '#00b976',
  brightGreen: '#05d289',
  yellow: '#e0d500',
  brightYellow: '#f4f628',
  magenta: '#bd37bc',
  brightMagenta: '#d86cd8',
  blue: '#1d6fca',
  brightBlue: '#358bed',
  cyan: '#00a8cf',
  brightCyan: '#19b8dd',
  white: '#e5e5e5',
  brightWhite: '#ffffff'
}
const bindTerminalResize = (term, websocket) => {
  const onTermResize = size => {
    websocket.send(
        Base64.encode(
            JSON.stringify({
              type: 'resize',
              rows: size.rows,
              cols: size.cols
            })
        )
    )
  }
  // register resize event.
  term.on('resize', onTermResize)
  // unregister resize event when WebSocket closed.
  websocket.addEventListener('close', function() {
    term.off('resize', onTermResize)
  })
}
const bindTerminal = (term, websocket, bidirectional, bufferedTime) => {
  term.socket = websocket
  let messageBuffer = null
  const handleWebSocketMessage = function(ev) {
    if (bufferedTime && bufferedTime > 0) {
      if (messageBuffer) {
        messageBuffer += Base64.decode(ev.data)
      } else {
        messageBuffer = Base64.decode(ev.data)
        setTimeout(function() {
          term.write(messageBuffer)
        }, bufferedTime)
      }
    } else {
      term.write(Base64.decode(ev.data))
    }
  }
  const handleTerminalData = function(data) {
    websocket.send(
        Base64.encode(
            JSON.stringify({
              type: 'input',
              input: data
            })
        )
    )
  }
  websocket.onmessage = handleWebSocketMessage
  if (bidirectional) {
    term.on('data', handleTerminalData)
  }
  // send heartbeat package to avoid closing webSocket connection in some proxy environmental such as nginx.
  const heartBeatTimer = setInterval(function() {
    websocket.send(
        Base64.encode(
            JSON.stringify({
              type: 'heartbeat',
              data: ''
            })
        )
    )
    // websocket.send('1')
  }, 20 * 1000)
  websocket.addEventListener('close', function() {
    websocket.removeEventListener('message', handleWebSocketMessage)
    term.off('data', handleTerminalData)
    delete term.socket
    clearInterval(heartBeatTimer)
  })
}

export default {
  data() {
    return {
      namespace: "",
      namespaces: [],
      tableData: [],
      tableDatas: {},
      pods: [],
      fileBrowserData: [],
      continue: "",
      totalCount: 0,
      pageSizes: [10, 30, 50, 100, 150, 200],
      pageSize: 10,
      currentPage: 1,
      podName: "",
      container:"",
      path: "",
      bulkPath: [],
      globalPath: "",
      headerPaths: [],
      dialogTerminalVisible: false,
      dialogFileBrowserVisible: false,
      dialogFileVisible: false,
      wsUrl: "",
      isFullScreen: false,
      searchKey: '',
      v: this.visible,
      ws: null,
      term: null,
      thisV: this.visible,
      createForPath: "",
      createName: "",
      fileContent: "",
      isNewFile: false,
      editorOption: {
        theme: "snow", // or 'bubble'
        modules: {
          toolbar: {
            container: toolbarOptions,
          }
        }
      },
    }
  },
  methods: {
    getNamespace() {
      GetNamespace().then(res =>{
        if (res) {
          this.namespaces = []
          const data = res.items
          for(const key in data){
            this.namespaces.push(data[key].metadata.name)
          }
        }
      })
    },
    // 分页
    // 每页显示的条数
    handleSizeChange(val) {
      // 改变每页显示的条数
      this.pageSize=val
      if (this.namespace !== "") {
        // 点击每页显示的条数时，显示第一页
        this.selectedNamespace(this.namespace)
        // 注意：在改变每页显示的条数时，要将页码显示到第一页
        this.currentPage=1
      }
    },
    // 显示第几页
    handleCurrentChange(val) {
      // 改变默认的页数
      this.currentPage=val
      this.tableData = this.tableDatas[val]
    },
    selectedNamespace(options) {
      this.pods = []
      this.tableData = []
      this.continue = ""
      this.namespace = options
      // 处理第一页
      GetPods({namespace: options, limit: this.pageSize}).then(async res => {
        this.totalCount = res.metadata.remainingItemCount
        this.tableDatas[1] = this.getTableData(res)
        if (this.totalCount === null || this.totalCount === undefined) {
          this.totalCount = this.tableDatas[1].length
        } else {
          this.totalCount += this.pageSize
        }
        // 处理剩余页面
        for (let i = 1; i < Math.ceil(this.totalCount / this.pageSize); i++) {
          await GetPods({namespace: options, limit: this.pageSize, continue: this.continue}).then(res => {
            this.tableDatas[i+1] = this.getTableData(res)
          })
        }
        this.tableData = this.tableDatas[1]
      }, err => {
        this.$message.error(err.message)
      })
    },
    getTableData (res) {
      this.continue = res.metadata.continue
      let tableData = []
      for (const i in res.items) {
        const pod = res.items[i]
        let tr = {
          Namespace: pod.metadata.namespace,
          PodName: pod.metadata.name,
          State: pod.status.phase,
          OS: pod.metadata.annotations.os,
          Arch: pod.metadata.annotations.arch,
          CreateTime: pod.metadata.creationTimestamp,
          PodIPs: pod.status.podIPs,
          NodeName: pod.spec.nodeName,
          HostIP: pod.status.hostIP,
          Containers: pod.spec.containers,
          RestartCount: this.getRestartTimes(pod),
        }
        tableData.push(tr)
      }
      return tableData
    },
    getRestartTimes (row) {
      if (row.status.containerStatuses) {
        let restartCount = 0
        for (const c of row.status.containerStatuses) {
          restartCount += c.restartCount
        }
        return restartCount
      }
      return 0
    },
    openTerminal(options, container) {
      let shell = "bash"
      if (options.OS === "windows") {
        shell = "cmd"
      }
      this.dialogTerminalVisible = true
      this.wsUrl = "ws://"+window.location.host+"/api/kubeapiproxy/terminal?namespace="+this.namespace+"&pod="+options.PodName+"&container="+container.name+"&shell="+shell;
    },
    openFileBrowser(options, container, path) {
      if (path === undefined) {
        path = "/"
      }
      if (path === "/" && options !== null) {
        this.podName = options.PodName
      }
      if (container != null) {
        this.container = container.name
      }
      this.headerPaths = []
      this.globalPath=path
      this.headerPaths.push(path)
      if (path !== undefined) {
        let _p = path.split('/')
        let _pa = ""
        this.headerPaths = []
        _p.forEach((item,index) => {
          if (index === 0) {
            _pa = "/"
            item = "/"
            this.headerPaths.push({
              name: item,
              path: _pa,
            })
          }
          if (index !== 0 && item !== "") {
            _pa += item + "/"
            this.headerPaths.push({
              name: item,
              path: _pa,
            })
          }
        })
      }
      this.path = path
      this.fileBrowserData = []
      FileBrowserList({
        namespace: this.namespace,
        pod: this.podName,
        container: this.container,
        path: path,
      }).then(res => {
        this.dialogFileBrowserVisible = true
        this.fileBrowserData = []
        if (res !== undefined) {
          this.fileBrowserData = res
        }
      }, err => {
        this.$message.error(err.message)
      })
    },
    openFileDialog(path, type) {
      this.fileContent = ""
      this.dialogFileVisible=true
      this.createForPath = path
      this.isNewFile = true
      if (type === "open") {
        this.isNewFile = false
        this.createName = ""
        FileBrowserOpen({
          namespace: this.namespace,
          pod: this.podName,
          container: this.container,
          path: path,
        }).then(res => {
          console.log(res)
          if (res !== undefined) {
            this.fileContent = res
          }
        }, err => {
          this.dialogFileVisible=false
          this.$message.error(err.message)
        })
      }
    },
    saveFile() {
      let path = this.createForPath+"/"+this.createName
      if (this.createName === "") {
        path = this.createForPath
      }
      FileBrowserCreateFile(this.fileContent, {
        namespace: this.namespace,
        pod: this.podName,
        container: this.container,
        path: path,
      }).then(res => {
        console.log(res)
        if (res !== undefined) {
          this.dialogFileVisible = false
          this.fileContent = ""
          this.createName = ""
          this.$message.success(res)
          this.openFileBrowser(null, null, this.path)
        }
      }, err => {
        this.$message.error(err.message)
      })
    },
    createDir() {
      this.$prompt(this.$t('please_input_name'), this.$t('tips'), {
        confirmButtonText: this.$t('enter'),
        cancelButtonText: this.$t('cancel'),
        type: 'warning'
      }).then(({value}) => {
        if(!value) {//对输入内容校验
          return this.$t('please_input_name');
        }
        FileBrowserCreateDir({
          namespace: this.namespace,
          pod: this.podName,
          container: this.container,
          path: this.path+"/"+value,
        }).then(res => {
          if (res !== undefined) {
            this.$message.success(res)
            this.openFileBrowser(null, null, this.path)
          }
        }, err => {
          this.$message.error(err.message)
        })
      }).catch(() => {
        this.$message.info(this.$t('cancel'))
      });
    },
    openRenameDialog(oldName) {
      this.$prompt(this.$t('please_input_name')+"\n"+oldName, this.$t('tips'), {
        confirmButtonText: this.$t('enter'),
        cancelButtonText: this.$t('cancel'),
        type: 'warning'
      }).then(({value}) => {
        if(!value) {//对输入内容校验
          return this.$t('please_input_name');
        }
        FileBrowserRename({
          namespace: this.namespace,
          pod: this.podName,
          container: this.container,
          old_path: this.path+"/"+oldName,
          path: this.path+"/"+value,
        }).then(res => {
          console.log(res)
          if (res !== undefined) {
            this.$message.success(res)
            this.openFileBrowser(null, null, this.path)
          }
        }, err => {
          this.$message.error(err.message)
        })
      }).catch(() => {
        this.$message.info(this.$t('cancel'))
      });
    },
    removeFileOrDir(path) {
      this.$confirm(this.$t('tips_msg')+"\n"+path, this.$t('tips'), {
        confirmButtonText: this.$t('enter'),
        cancelButtonText: this.$t('cancel'),
        type: 'warning'
      }).then(() => {
        FileBrowserRemove({
          namespace: this.namespace,
          pod: this.podName,
          container: this.container,
          path: path,
        }).then(res => {
          console.log(res)
          if (res !== undefined) {
            this.$message.success(res)
            this.openFileBrowser(null, null, this.path)
          }
        }, err => {
          this.$message.error(err.message)
        })
      }).catch(() => {
        this.$message.info(this.$t('cancel'))
      });
    },
    handleSelectionChange(val) {
      this.bulkPath = []
      val.forEach((item) => {
        this.bulkPath.push(item.Path)
      })
    },
    download(path, style) {
      Download({namespace: this.namespace, pod: this.podName, container: this.container, dest_paths: path, style: style}).then(res=>{})
    },
    bulkDownload(paths, style) {
      if (paths.length === 0) {
        this.$message.error(this.$t('cannot_empty'))
        return
      }
      let path = ""
      paths.forEach(item => {
        path += "&dest_paths="+item
      })
      Download({namespace: this.namespace, pod: this.podName, container: this.container, dest_paths: path, style: style}).then(res=>{})
    },
    uploadFileOrDir(e, path) {
      const files = e.target.files;
      if (files.length === 0 ) {
        e.target.value = ""
        return
      }
      const formData = new FormData();
      //追加文件数据
      for (let i = 0; i < files.length; i++) {
        formData.append("files", files[i]);
      }
      Upload(formData, {
        namespace:this.namespace,
        pod:this.podName,
        containers:this.container,
        dest_path:path},{"Content-Type":"multipart/form-data"}).then((res) => {
          this.openFileBrowser(null, null, path)
          this.$prompt(res, this.$t('tips'), {
            confirmButtonText: this.$t('enter'),
            cancelButtonText: this.$t('cancel'),
            showInput: false,
            type: 'success'
          })
      }, (err) => {
        this.$message.error(err.message)
      })
      e.target.value = ""
    },

    onWindowResize() {
      // console.log("resize")
      // this.term.fit() // it will make terminal resized.
      // this.term.scrollToBottom();
      let height = document.body.clientHeight;
      let rows = height/23;
      this.term.fit();
      this.term.resize(this.term.cols,parseInt(rows))//终端窗口重新设置大小 并触发term.on("resize"
      this.term.scrollToBottom();
    },
    doLink(ev, url) {
      if (ev.type === 'click') {
        window.open(url)
      }
    },
    doClose() {
      window.removeEventListener('resize', this.onWindowResize)
      // term.off("resize", this.onTerminalResize);
      if (this.ws) {
        this.ws.close()
      }
      if (this.term) {
        this.term.dispose()
      }
      this.$emit('pclose', false)// 子组件对openStatus修改后向父组件发送事件通知
    },
    doOpened() {
      Terminal.applyAddon(fit)
      Terminal.applyAddon(webLinks)
      Terminal.applyAddon(search)
      this.term = new Terminal({
        rendererType: 'canvas', // 渲染类型
        rows: parseInt(document.body.clientHeight/23),
        cols: parseInt(document.body.clientWidth),
        convertEol: true, // 启用时，光标将设置为下一行的开头
        // scrollback: 10, // 终端中的回滚量
        disableStdin: false, // 是否应禁用输入
        fontSize: 18,
        cursorBlink: true, // 光标闪烁
        cursorStyle: 'bar', // 光标样式 underline
        bellStyle: 'sound',
        theme: defaultTheme
      })
      this.term._initialized = true
      this.term.prompt = () => {
        this.term.write('\r\n')
      }
      this.term.prompt()
      this.term.on('key', function(key, ev) {
        console.log(key, ev, ev.keyCode)
      })
      this.term.open(this.$refs.terminal)
      this.term.webLinksInit(this.doLink)
      // term.on("resize", this.onTerminalResize);
      this.term.on('resize', this.onWindowResize)
      window.addEventListener('resize', this.onWindowResize)
      this.term.fit() // first resizing
      this.ws = new WebSocket(this.wsUrl)
      this.ws.onerror = () => {
        this.$message.error(this.$t('web_socker_connection_failed'))
      }
      this.ws.onclose = () => {
        this.term.setOption('cursorBlink', false)
        this.$message(this.$t('web_socket_disconnect'))
      }
      bindTerminal(this.term, this.ws, true, -1)
      bindTerminalResize(this.term, this.ws)
    }
  },
  watch:{
    visible(val) {
      this.v = val// 新增result的watch，监听变更并同步到myResult上
    }
  },
  created:function (){
    this.getNamespace()
  },
}
</script>