import Vue from 'vue'
import ElementUI, {MessageBox} from 'element-ui'
import Plugin from 'v-fit-columns'
import 'element-ui/lib/theme-chalk/index.css'
import VueQuillEditor from 'vue-quill-editor'
import 'quill/dist/quill.core.css'
import 'quill/dist/quill.snow.css'
import 'quill/dist/quill.bubble.css'
import moment from 'moment'
import App from './App.vue'
import router from './router'
import i18n from './lang'
import util from './lib/util.js'
import data from './lib/data.js'
import 'xterm/dist/xterm.css'
import 'xterm/dist/xterm.js'
import './scss/app.scss'
import './lib/directives.js'

let localeLang
if (global.navigator.language) {
    localeLang = global.navigator.language
    localeLang = localeLang.toLowerCase()
}
if (localeLang.indexOf('en') !== 0) {
    localeLang = 'zh-cn'
}
moment.locale(localeLang)
Vue.config.debug = true;
Vue.config.productionTip = false
Vue.use(ElementUI)
Vue.use(Plugin)
Vue.use(VueQuillEditor)
Vue.prototype.$msgbox = MessageBox;
Vue.prototype.$alert = MessageBox.alert;
Vue.prototype.$confirm = MessageBox.confirm;
Vue.prototype.$prompt = MessageBox.prompt;

new Vue({
    i18n,
    router,
    methods: util,
    data: data,
    render: h => h(App)
}).$mount('#app')