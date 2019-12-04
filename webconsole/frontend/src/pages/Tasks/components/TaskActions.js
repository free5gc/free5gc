import React from 'react';
import Constants from '../../../config/Constants';
import TaskUtils from "../../../util/TaskUtils";

const TaskActions = ({taskItem}) => {
  if (taskItem['name'] === Constants.TASK_CSV_EXPORT['name'] &&
      taskItem['queue_status'] === Constants.QueueStatus.FINISHED) {
    let exported = JSON.parse(taskItem['output'])['exported'];

    if (exported) {
      return (
        <button type="button" className="btn btn-primary"
                onClick={() => TaskUtils.launchSpotsCsvDownload(taskItem)}>
          <i className="fa fa-cloud-download-alt"/>&nbsp;
          CSV
        </button>
      );
    } else {
      return (
        <div className="fa fa-exclamation-triangle icon-warning" style={{width: 'min-content'}}/>
      );
    }

  }

  // Else
  return (
    <div/>
  );
};

export default TaskActions;
