import axios from 'axios';
import config from '../config/config';

axios.defaults.baseURL = config.API_URL;
axios.defaults.headers.common.Accept = 'application/json';
axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
axios.defaults.crossdomain = true;

// Request interceptor
axios.interceptors.request.use(config => {
  // Before request is sent
  // config.headers.common['Authorization'] = generateAuthHeader();

  return config;
}, error => Promise.reject(error));

// Response interceptor
axios.interceptors.response.use(
  response => response,
  async error => Promise.reject(error)
);

// function generateAuthHeader() {
//   let apiTokens = store ? store.getState().Auth.apiTokens : null;
//   return apiTokens ? `Bearer ${apiTokens.accessToken}` : null;
// }

export default axios;
