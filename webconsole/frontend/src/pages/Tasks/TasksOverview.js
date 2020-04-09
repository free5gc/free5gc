import React, {Component} from 'react';
import {withRouter} from "react-router-dom";
import {connect} from "react-redux";
import queryString from 'query-string';
import AppUtils from "../../util/AppUtils";

class TasksOverview extends Component {

  constructor(props) {
    super(props);
    this.refreshEnabled = false;
    this.targetTaskUuid = null;
    this.targetDownloaded = false;
  }

  componentWillMount() {
    this.refreshEnabled = true;
    this.updateTasksTable().then();
  }

  componentWillUnmount() {
    this.refreshEnabled = false;
  }

  componentWillReceiveProps(nextProps) {
    let urlParams = queryString.parse(this.props.location.search);

    if (urlParams['target'] !== undefined) {
      let dashedUuid = AppUtils.dashUuid(urlParams['target']);
      if (this.targetTaskUuid !== dashedUuid) {
        this.targetTaskUuid = dashedUuid;
        this.targetDownloaded = false;
      }
    }

    // Track the target task
    /*
    if (this.targetTaskUuid !== null && !this.targetDownloaded && nextProps.tasksMap[this.targetTaskUuid] !== undefined) {
      let task = nextProps.tasksMap[this.targetTaskUuid];
      if (task['queue_status'] === Constants.QueueStatus.FINISHED) {
        // TODO
        this.targetDownloaded = true;
      }
    } */
  }

  async updateTasksTable() {
    if (this.refreshEnabled) {
      // await ApiHelper.fetchTasks();
      await AppUtils.wait(1000);
      this.updateTasksTable().then();
    }
  }

  render() {
    return (
      <div className="container-fluid">

        <div className="row">
          <div className="col-md-12">
            {/*<TasksTable data={this.props.tasks} targetTask={this.targetTaskUuid}/>*/}
          </div>
        </div>

      </div>
    );
  }
}

const mapStateToProps = state => ({
  // tasks: state.tasks.tasks,
  // tasksMap: state.tasks.tasksMap
});

export default withRouter(connect(mapStateToProps)(TasksOverview));
