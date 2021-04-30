import axios from 'axios'
import { ElMessage } from 'element-plus'
import { STORAGE_KEY, ErrMsgMap } from './vars.js'

export const http = axios.create({
    baseURL: '/api',
    timeout: 5000,
    withCredentials: false,
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    transformRequest: [
        function (data) {
            let ret = ''
            for (let it in data) {
                ret +=
                    encodeURIComponent(it) +
                    '=' +
                    encodeURIComponent(data[it]) +
                    '&'
            }
            ret = ret.substring(0, ret.lastIndexOf('&'))
            return ret
        }
    ]
})

// request拦截器 ==> 对请求参数进行处理
http.interceptors.request.use(
    (config) => {
        // 可以在发送请求之前做些事情
        // 比如请求参数的处理、在headers中携带token等等
        console.log('send ajax with config')
        console.log(config)
        let s = getStorage()
        console.log(s)
        if (s) {
            let token = s.token
            if (token) {
                config.headers.Authorization = `Bearer ${token}`
            }
        }
        console.log(config)
        return config
    },
    (error) => {
        // 处理请求错误
        console.error(error)
        return Promise.reject(error)
    }
)

// respone拦截器 ==> 对响应做处理
http.interceptors.response.use(
    (response) => {
        // 2xx 范围内的状态码都会触发该函数。
        // 接口返回success字段不为true表示请求错误
        console.log(response)
        let data = response.data
        if (!data.success) {
            let prefix = getErrMsgPrefix(response.config.url),
                text = data.message || 'Error'
            ElMessage.error(prefix + text)
            return Promise.reject(new Error(text))
        } else {
            return data
        }
    },
    (error) => {
        let prefix = getErrMsgPrefix(error.config.url),
            text = ''

        if (error.response) {
            // 请求成功发出且服务器也响应了状态码，但状态代码超出了 2xx 的范围
            console.log(error.response.data)
            console.log(error.response.status)
            console.log(error.response.headers)
            switch (error.response.status) {
                case 400:
                    text = '参数错误'
                    break
                case 401:
                    text = '未登录'
                    return this.$router.replace('/login')
                case 403:
                    text = '拒绝访问'
                    break
                case 404:
                    text = '无效接口'
                    break
                case 500:
                    text = '服务器错误'
                    break
                case 501:
                    text = '服务未实现'
                    break
                case 502:
                    text = '网关错误'
                    break
                case 503:
                    text = '服务不可用'
                    break
                case 504:
                    text = '网关超时'
                    break
                default:
                    text = '网络连接出错'
            }
            let data = error.response.data
            if (isObject(data) && data.message !== '') {
                text = text + `（ ${data.message} ）`
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
        console.log(error.config)
        ElMessage.error(prefix + text)
        return Promise.reject(error)
    }
)

function getErrMsgPrefix(url) {
    for (let key of Object.keys(ErrMsgMap)) {
        if (url.startsWith(key)) {
            return ErrMsgMap[key]
        }
    }
    return ''
}

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

export function formatUnixTimestamp(unixtimestamp) {
    var unixtimestamp = new Date(unixtimestamp * 1000)
    var year = 1900 + unixtimestamp.getYear()
    var month = '0' + (unixtimestamp.getMonth() + 1)
    var date = '0' + unixtimestamp.getDate()
    var hour = '0' + unixtimestamp.getHours()
    var minute = '0' + unixtimestamp.getMinutes()
    //var second = '0' + unixtimestamp.getSeconds()
    return (
        year +
        '-' +
        month.substring(month.length - 2, month.length) +
        '-' +
        date.substring(date.length - 2, date.length) +
        ' ' +
        hour.substring(hour.length - 2, hour.length) +
        ':' +
        minute.substring(minute.length - 2, minute.length)
    )
}
