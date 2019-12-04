import React from 'react';
import {Route} from 'react-router-dom';
import TasksOverview from "./TasksOverview";

const Tasks = ({match}) => (
  <div className="content">
    <Route exact path={`${match.url}/`} component={TasksOverview} />
    {/*<Route path={`${match.url}/:uuid`} component={} />*/}
  </div>
);

export default Tasks;
