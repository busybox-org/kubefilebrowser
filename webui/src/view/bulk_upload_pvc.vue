<template>
  <div>
    <el-card shadow="never">
      <div>
        <el-select style="width: 100%" v-model="namespace" @click.native="getNamespace" @change="selectedNamespace" filterable :placeholder="$t('please_select_namespace')">
          <el-option
              v-for="item in namespaces"
              :label="item"
              :value="item"
              :key="item"
          ></el-option>
        </el-select>
      </div>
      <div style="margin-top: 15px">
        <el-select style="width: 100%" :placeholder="$t('please_select_pvc')" filterable v-model="pvc">
          <el-option v-for="item in pvcs" :label="item.label" :value="item.value" :key="item.value"></el-option>
        </el-select>
      </div>
      <div style="margin-top: 15px">
        <el-dropdown  type="success" class="avatar-container" trigger="click" style="height: 36px;float: right;margin-bottom: 10px;">
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
                  <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event)" name="files" multiple="true">
                </span>
            </el-dropdown-item>
            <el-dropdown-item divided>
                <span class="fake-file-btn">
                  {{ $t('upload_dir') }}
                  <input type="file" style="display:block;" v-on:change="uploadFileOrDir($event)" name="files" webkitdirectory mozdirectory accept="*/*">
                </span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </el-card>
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
import {
  GetNamespace,
  GetPvcs,
  Upload,
  UploadPVC,
  Download,
  CTerminal,
  Exec,
} from '../api/kubeapiproxy'
import {
  FileBrowserCreateDir,
  FileBrowserCreateFile,
  FileBrowserList,
  FileBrowserOpen,
  FileBrowserRemove,
  FileBrowserRename,
} from "../api/filebrowser";

export default {
  data() {
    return {
      namespace: "",
      pvc: [],
      namespaces:[],
      pvcs: [],
      destPath:"",
      fileList:[],
    }
  },
  methods: {
    getNamespace() {
      GetNamespace().then(res => {
        if (res) {
          this.namespaces = []
          this.pvc = []
          this.pvcs = []
          const data = res.items;
          for (const key in data) {
            this.namespaces.push(data[key].metadata.name)
          }
        }
      })
    },
    selectedNamespace() {
      GetPvcs({namespace: this.namespace}).then(res => {
        if (res) {
          const pvcs = [];
          this.pvc = []
          this.pvcs = []
          const data = res.items;
          for (const key in data) {
            const _d = {label: data[key].metadata.name, value: data[key].metadata.name}
            pvcs.push(_d)
          }
          this.pvcs = pvcs
        }
      })
    },
    uploadFileOrDir(e) {
      const files = e.target.files;
      if (files.length === 0 ) {
        e.target.value = ""
        return
      }
      let pvc = this.pvc
      const formData = new FormData();
      //追加文件数据
      for (let i = 0; i < files.length; i++) {
        formData.append("files", files[i]);
      }
      UploadPVC(formData, {
        namespace:this.namespace,
        pvc:pvc},{"Content-Type":"multipart/form-data"}).then((res) => {
          this.$message.success("upload success!")
      }, (err) => {
        console.log(err.message)
        this.$message.error(err.message)
      })
      e.target.value = ""
    },
  }
}
</script>