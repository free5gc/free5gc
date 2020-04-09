export default class UEInfoDetail {
  
  ueInfoDetail = {
      amfInfo:{},
      smfInfo:{},
      pcfInfo:{}
  }

  constructor(info) {
     this.amfInfo = info.amfInfo;
     this.smfInfo = info.smfInfo;
     this.pcfInfo = info.pcfInfo;
  }
}