/*
   Copyright 2021 Hiroshi.tao

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import { reactive } from 'vue'
import {
    http,
    getStorage,
    setStorage,
    clearStorage,
    isValidMap,
    isObject
} from './util.js'

export const state = reactive(getStorage() || {})

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
        http.get('/config').then(function (res) {
            for (let k in res.data) {
                mutations.commit(k, res.data[k])
            }
            actions.saveConfig2Local()
        })
    },
    saveConfig2Local() {
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
