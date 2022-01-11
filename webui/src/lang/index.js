import Vue from 'vue'
import VueI18n from 'vue-i18n'

Vue.use(VueI18n)

import zh from './zh_cn.js'
import en from './en.js'
const messages = {
    'zh-cn': zh,
    'en': zh,
}

let localeLang
if (global.navigator.language) {
    localeLang = global.navigator.language
    localeLang = localeLang.toLowerCase()
}
if (localeLang.indexOf('en') !== 0) {
    localeLang = 'zh-cn'
}

const i18n = new VueI18n({
    locale: localeLang,
    fallbackLocale: 'zh-cn',
    silentFallbackWarn: true,
    messages,
})

export default i18n