import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: '/api',
  timeout: 30000, // 30秒超时
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    console.log('发送请求:', config.method?.toUpperCase(), config.url)
    return config
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    console.log('收到响应:', response.status, response.data)
    return response
  },
  (error) => {
    console.error('响应错误:', error)
    
    // 处理不同类型的错误
    if (error.response) {
      // 服务器返回了错误状态码
      const { status, data } = error.response
      let message = '请求失败'
      
      switch (status) {
        case 400:
          message = data.error || '请求参数错误'
          break
        case 500:
          message = data.error || '服务器内部错误'
          break
        case 502:
          message = '服务器网关错误'
          break
        case 503:
          message = '服务暂时不可用'
          break
        case 504:
          message = '请求超时'
          break
        default:
          message = data.error || `请求失败 (${status})`
      }
      
      const customError = new Error(message)
      customError.status = status
      customError.data = data
      return Promise.reject(customError)
    } else if (error.request) {
      // 请求已发出但没有收到响应
      const customError = new Error('网络连接失败，请检查网络或服务器状态')
      customError.isNetworkError = true
      return Promise.reject(customError)
    } else {
      // 其他错误
      return Promise.reject(error)
    }
  }
)

// API服务类
class ApiService {
  /**
   * 上传文档
   * @param {string} title - 文档标题
   * @param {string} content - 文档内容
   * @returns {Promise<Object>} 上传结果
   */
  async uploadDocument(title, content) {
    try {
      const response = await api.post('/documents', {
        title,
        content
      })
      return response.data
    } catch (error) {
      throw new Error(error.message || '文档上传失败')
    }
  }

  /**
   * 发送问题
   * @param {string} question - 问题内容
   * @returns {Promise<Object>} 问答结果
   */
  async askQuestion(question) {
    try {
      const response = await api.post('/chat', {
        question
      })
      return response.data
    } catch (error) {
      throw new Error(error.message || '问答请求失败')
    }
  }

  /**
   * 检查服务器连接状态
   * @returns {Promise<boolean>} 连接状态
   */
  async checkConnection() {
    try {
      // 尝试发送一个简单的问答请求来检查连接
      await api.post('/chat', { question: 'test' })
      return true
    } catch (error) {
      return false
    }
  }
}

// 创建API服务实例
export const apiService = new ApiService()

// 导出API实例（如果需要直接使用axios）
export { api }
