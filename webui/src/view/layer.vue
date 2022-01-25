<template>
  <div class="layer-global">
    <header class="layer-header">
      <div class="header-left">
        <div class="logo">
          <h2>KubernetesFileBrowser</h2>
          <i class="iconfont icon-kubefilebrowser"></i>
        </div>
      </div>
      <div class="header-right">
        <span class="r-item">
          <el-dropdown trigger="click">
            <span class="item app-cursor">
              <i class="iconfont icon-question-circle-fill"></i>
              <i class="iconfont small icon-arrow-down">{{ $t('more') }}</i>
            </span>
            <el-dropdown-menu slot="dropdown" class="app-header-dropdown">
              <a class="app-dropdown-link" href="https://github.com/xmapst/kubefilebrowser/issues" target="_blank">
                <el-dropdown-item><i class="iconfont small left icon-help"></i>{{ $t('help') }}</el-dropdown-item>
              </a>
              <a class="app-dropdown-link" href="https://github.com/xmapst/kubefilebrowser" target="_blank">
                <el-dropdown-item><i class="iconfont small left icon-pull-request"></i>{{ $t('contribute_to_kube_file_browser') }}</el-dropdown-item>
              </a>
            </el-dropdown-menu>
          </el-dropdown>
        </span>
      </div>
    </header>
    <section class="layer-container">
      <aside class="layer-aside">
        <ScrollBar>
          <el-menu class="aside-menu" :default-active="activeMenu" :router="true" :unique-opened="true">
            <template v-for="menu in AppMenu">
              <el-submenu v-if="menu.children && (menu.children.length > 1 || (menu.children.length === 1 && !menu.children[0].meta.single))" :index="menu.name" :key="menu.name">
                <template slot="title">
                  <span v-if="menu.meta.icon" class="iconfont left" :class="menu.meta.icon"></span><span>{{ menu.meta.title }}</span>
                </template>
                <template v-for="childMenu in menu.children">
                  <el-menu-item v-if="!(childMenu.meta && childMenu.meta.hide)" :route="{name: childMenu.name}" :index="childMenu.name" :key="childMenu.name">
                    <i class="iconfont small left">
                      <svg viewBox="0 0 1024 1024" width="1em" height="1em" fill="currentColor" aria-hidden="true"><path d="M384.023552 384.083968l256.016384 0 0 256.016384-256.016384 0 0-256.016384Z"></path></svg>
                    </i>
                    <span>{{ childMenu.meta.title }}</span>
                  </el-menu-item>
                </template>
              </el-submenu>
              <el-menu-item :route="{name: menu.children[0].name}" v-else-if="menu.children && menu.children.length === 1" :index="menu.children[0].name" :key="menu.children[0].name">
                <i v-if="menu.children[0].meta.icon" class="iconfont left" :class="menu.children[0].meta.icon"></i>
                <span>{{ menu.children[0].meta.title }}</span>
              </el-menu-item>
            </template>
          </el-menu>
        </ScrollBar>
      </aside>
      <main class="layer-main">
        <div class="container">
          <router-view/>
          <div class="app-cpy">
            <p>
              ©️ {{ new Date().getFullYear() }} <a href="https://github.com/xmapst/kubefilebrowser" target="_blank">KubernetesFileBrowser</a>. All Rights Reserved. MIT License.
            </p>
          </div>
        </div>
      </main>
    </section>
  </div>
</template>

<script>
import ScrollBar from '../component/ScrollBar';
import {routes} from '../router'

export default {
  data() {
    return {
      breadcrumb: [],
      activeMenu: '',
      btnLoading: false,
      settingDialogVisible: false,
      settingForm: {},
      passwordDialogVisible: false,
      passwordForm: {},
    }
  },
  computed: {
    AppMenu() {
      let menu = []
      routes.forEach(first => {
        let newFirst = Object.assign({}, first)
        menu.push(newFirst)
      })
      return menu
    }
  },
  watch: {
    '$route.name'() {
      this.breadcrumbItems()
      this.initActiveMenu()
    },
    AppMenu() {
      this.breadcrumbItems()
    },
  },
  components: {
    ScrollBar,
  },
  methods: {
    initActiveMenu() {
      this.activeMenu = this.$route.name
    },
    breadcrumbItems() {
      let breadcrumb = []
      this.AppMenu.forEach(menu => {
        menu.children.forEach(sub => {
          if (sub.name !== this.$route.name) {
            return
          }
          if (menu.meta.title) {
            breadcrumb.push(menu.meta.title)
          }
          breadcrumb.push(sub.meta.title)
        })
      })
      this.breadcrumb = breadcrumb
    },
  },
  mounted() {
    this.initActiveMenu()
  },
}
</script>
