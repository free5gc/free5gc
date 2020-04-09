import React from 'react';
import Constants from "../../../config/Constants";

const TasksTable = ({data, targetTask: targetTaskUuid}) => (
  <div className="card card-plain">
    <div className="header">
      <h4 className="title">Analysis Tasks</h4>
      {/*<p className="category">Here is a subtitle for this table</p>*/}
    </div>
    <div className="content table-responsive table-full-width">
      <table className="table table-hover tasks-table">
        <thead>
          <tr>
            <th>Status</th>
            <th className='col-md-2'>Type</th>
            <th className='col-md-2'>Created By</th>
            <th className='col-md-2'>Created At</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
        {data.map(item => (
          <tr key={item['uuid']} className={targetTaskUuid === item['uuid'] ? 'highlight' : ''}>
            <td width="1px">
              <StatusIndicator status={item['queue_status']}/>
            </td>
            <td>{taskName2Title(item.name)}</td>
            <td>{item['created_by'] === null ? '-' : item['created_by']}</td>
            <td>{item['created_at']}</td>
            <td>
              &nbsp;
            </td>
          </tr>
        ))}
        </tbody>
      </table>
    </div>
  </div>
);

const StatusIndicator = ({status}) => {
  if (status === Constants.QueueStatus.CREATED) {
    return (
      <div className="fa fa-clock icon-inactive"/>
    );
  } else if (status === Constants.QueueStatus.RUNNING) {
    return (
      <div className="spinner5"/>
    );
  } else if (status === Constants.QueueStatus.FINISHED) {
    return (
      <div className="fa fa-check-circle icon-active"/>
    );
  } else if (status === Constants.QueueStatus.ABORTED) {
    return (
      <div className="fa fa-exclamation-triangle icon-warning"/>
    );
  } else {
    return (
      <div/>
    );
  }
};

function taskName2Title(taskName) {
  if (taskName === Constants.TASK_CSV_EXPORT['name']) {
    return Constants.TASK_CSV_EXPORT['title'];
  } else {
    return '-';
  }
}

export default TasksTable;
