import { reactive } from 'vue'
import {
    http,
    getStorage,
    setStorage,
    clearStorage,
    isValidMap,
    isObject
} from './util.js'

console.log('init state')

export const state = reactive(getStorage() || {})
console.log({ ...state })

export const mutations = {
    isLogged() {
        return state.isLogin === true && state.token !== '' && state.user !== ''
    },
    setLogin(user, jwt) {
        this.commit('isLogin', true)
        this.commit('user', user)
        this.commit('token', jwt)
    },
    clearLogin() {
        this.commit('token', '')
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
                actions.saveConfig2Local()
            })
            .catch(function (e) {
                console.log(e)
            })
    },
    saveConfig2Local() {
        console.log('hit scl')
        setStorage({ ...state })
    },
    removeConfig() {
        clearStorage()
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
 * sts 可以是Array(String包装成Array)，获取其中的字段，适用于data
 * sts 可以是Object，key是键，val是字段，可以是function(state)，适用于computed
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
