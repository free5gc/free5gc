
const config = {
  API_URL: '',

  USERNAME: 'admin',
  PASSWORD: 'free5gc',
};

// [Note]
//  React's process.env is passed (configured) from npm's process.env by /config/env.js,
//  which only works when running on npm server (dev/testing purpose).
if (process.env.NODE_ENV === 'test') {
  config.API_URL = `http://localhost:${process.env.PORT}/api`;
} else {
  config.API_URL = process.env.REACT_APP_HTTP_API_URL ? process.env.REACT_APP_HTTP_API_URL : "/api";
}

export default config;
