import axios from 'axios'
import { ElMessage } from 'element-plus'
import { STORAGE_KEY, ErrMsgMap } from './vars.js'

export const http = axios.create({
    baseURL: '/api',
    timeout: 5000,
    withCredentials: false
})

// request拦截器 ==> 对请求参数进行处理
http.interceptors.request.use(
    (config) => {
        // 可以在发送请求之前做些事情
        // 比如请求参数的处理、在headers中携带token等等
        let s = getStorage()
        console.log(s)
        if (s) {
            let token = s.token
            if (token) {
                config.headers.Authorization = `Bearer ${token}`
            }
        }
        return config
    },
    (error) => {
        // 处理请求错误
        console.log(error) // for debug
        return Promise.reject(error)
    }
)

// respone拦截器 ==> 对响应做处理
http.interceptors.response.use(
    (response) => {
        // 2xx 范围内的状态码都会触发该函数。
        // 接口返回success字段不为true表示请求错误
        const res = response.data
        if (!res.success) {
            let msg = res.message || 'Error'
            let err = new Error(msg)
            err.text = msg
            return Promise.reject(err)
        } else {
            return res
        }
    },
    (error) => {
        let text = ''

        if (error.response) {
            // 请求成功发出且服务器也响应了状态码，但状态代码超出了 2xx 的范围
            console.log(error.response.data)
            console.log(error.response.status)
            console.log(error.response.headers)

            switch (error.response.status) {
                case 400:
                    text = '请求错误(400)，请重新申请'
                    break
                case 401:
                    text = '登录错误(401)，请重新登录'
                    return this.$router.replace('/login')
                case 403:
                    text = '拒绝访问(403)'
                    break
                case 404:
                    text = '请求出错(404)'
                    break
                case 408:
                    text = '请求超时(408)'
                    break
                case 500:
                    text = '服务器错误(500)'
                    break
                case 501:
                    text = '服务未实现(501)'
                    break
                case 502:
                    text = '网络错误(502)'
                    break
                case 503:
                    text = '服务不可用(503)'
                    break
                case 504:
                    text = '网络超时(504)'
                    break
                case 505:
                    text = 'HTTP版本不受支持(505)'
                    break
                default:
                    text = '网络连接出错'
            }
        } else if (error.request) {
            // 请求已经成功发起，但没有收到响应
            // `error.request` 在浏览器中是 XMLHttpRequest 的实例，
            // 而在node.js中是 http.ClientRequest 的实例
            text = '连接服务器失败，请稍后重试！'
        } else {
            // 请求已经成功发起，但没有收到响应
            text = error.message
        }
        // TODO
        let prefix = ErrMsgMap['xxx']
        ElMessage.error(prefix + text)
        return Promise.reject(error)
    }
)

/**
 * 设置过期存储
 * @param {*} data 数据
 * @param {String} key 索引
 * @param {Number} expire 过期毫秒数，默认0不过期
 */
export function setStorage(data, key = STORAGE_KEY, expire = 0) {
    if (key && data) {
        return localStorage.setItem(
            key,
            JSON.stringify({
                exp: expire || 0,
                time: Date.now(),
                data: data
            })
        )
    } else {
        return false
    }
}

export function getStorage(key = STORAGE_KEY) {
    try {
        let obj = JSON.parse(localStorage.getItem(key))
        if (obj && isObject(obj)) {
            let exp = parseInt(obj.exp) || 0,
                time = obj.time,
                data = obj.data
            //如果exp为真则说明有过期时间，否则永久有效直接返回数据
            let now = Date.now()
            if (exp > 0 && now > time + exp) {
                clearStorage(key)
            } else {
                return data
            }
        }
    } catch (e) {
        console.error(e)
    }
}

export function clearStorage(key = STORAGE_KEY) {
    return localStorage.removeItem(key)
}

export function isObject(o) {
    return o !== null && typeof o === 'object' && Array.isArray(o) === false
}

export function isValidMap(map) {
    return Array.isArray(map) || isObject(map)
}
