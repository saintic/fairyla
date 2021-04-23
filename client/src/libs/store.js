import { reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { http, getStorage, isValidMap, isObject } from './util.js'

console.log('init state')

export const state = reactive(getStorage() || {})
console.log({ ...state })

export const mutations = {
    isLogged() {
        return state.isLogin === true
    },
    setLogin(user) {
        this.commit('isLogin', true)
        this.commit('user', user)
    },
    clearLogin() {
        this.commit('isLogin', false)
        this.commit('user', '')
    },
    commit(key, value) {
        state[key] = value
    }
}

export const actions = {
    fetchConfig() {
        //get public config
        http.get('/config')
            .then(function (res) {
                console.log(res.data)
                /*
                let d = res.data
                mutations.commit('isLogin', d.status.isLogin)
                mutations.commit('username', d.status.user)
                mutations.commit('sapicAPI', d.sapic.api)
                mutations.commit('sapicSDK', d.sapic.sdk)
                mutations.commit('sapicToken', d.sapic.token)
                */
                for (let k in res.data) {
                    mutations.commit(k, res.data[k])
                }
            })
            .catch(function (e) {
                console.log(e)
            })
    }
}

/**
 * 将特定格式的Array|Object转化为Array
 * @param {Array} map 状态字段，可以嵌套Object
 */
function normalizeMap(map) {
    if (isObject(map)) map = [map]
    const ret = []
    for (let item of map) {
        if (isObject(item)) {
            ret.push(
                ...Object.keys(item).map((key) => ({ key, val: item[key] }))
            )
        } else {
            ret.push({ key: item, val: item })
        }
    }
    return ret
}

/**
 * 获取需要的状态数据对象
 * @param {String | Array | Object} sts 状态字段
 * @returns {Object}
 */
export const mapState = (sts) => {
    if (typeof sts === 'string') sts = [sts]
    if (!isValidMap(sts)) throw Error('Invalid type')
    const res = {}
    normalizeMap(sts).forEach(({ key, val }) => {
        res[key] = function mappedState() {
            return typeof val === 'function'
                ? val.call(this, state)
                : state[val]
        }
    })
    return res
}

export default { state, actions, mutations, mapState }
