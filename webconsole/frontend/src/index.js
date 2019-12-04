import React from 'react';
import ReactDOM from 'react-dom';
import registerServiceWorker from './registerServiceWorker';
import {HashRouter} from 'react-router-dom';
import './assets/styles/base.scss';
import App from './pages/App';
import configureStore from './config/configureStore';
import {Provider} from 'react-redux';

export const store = configureStore();
const rootElement = document.getElementById('root');

const renderApp = Component => {
  ReactDOM.render(
    <Provider store={store}>
      <HashRouter>
        <Component />
      </HashRouter>
    </Provider>,
    rootElement
  );
};

renderApp(App);

if (module.hot) {
  module.hot.accept('./pages/App', () => {
    const NextApp = require('./pages/App').default;
    renderApp(NextApp);
  });
}

registerServiceWorker();
