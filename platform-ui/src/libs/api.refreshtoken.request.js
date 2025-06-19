import config from '@/config'
import axios from 'axios'
const baseUrl = process.env.NODE_ENV === 'development' ? config.baseUrl.dev : config.baseUrl.pro

const refreshTokenAxios = axios.create({
  baseURL: baseUrl
})
export default refreshTokenAxios
