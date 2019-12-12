import ApiHelper from "./ApiHelper";

class TaskUtils {

  static async launchSpotsCsvDownload(taskItem) {
    let taskUuid = taskItem['uuid'];
    let t = taskItem['created_at'];

    // e.g. '2018-07-04 08:51:50' -> '0704_085150.csv'
    let downloadName = t.substr(5, 2) + t.substr(8, 2) + '_' + t.substr(11, 2) + t.substr(14, 2) + t.substr(17, 2) + '.csv';

    window.location.href = await ApiHelper.getDownloadUrl('spots-csv/' + taskUuid + '/' + downloadName);
  }

}

export default TaskUtils;
