export default class ueinfoActions {
    static SET_REG_UE = 'UEINFO/SET_REG_UE';
    static SET_UE_DETAIL = 'UEINFO/SET_UE_DETAIL';
    static SET_UE_DETAIL_AMF = 'UEINFO/SET_UE_DETAIL_AMF'
    static SET_UE_DETAIL_SMF = 'UEINFO/SET_UE_DETAIL_SMF'
    static SET_REG_UE_ERR = 'UEINFO/SET_REG_UE_ERR';
    static SET_UE_DETAIL_SM_CTX_REF = 'UEINFO/SET_UE_DETAIL_SM_CTX_REF'
  
    /**
     * @param users  {User}
     */
    static setRegisteredUE(users) {
      return {
        type: this.SET_REG_UE,
        registered_users: users,
      };
    }
  
    static setUEInfoDetail(ueInfoDetail) {
      return {
        type: this.SET_UE_DETAIL,
        ueInfoDetail: ueInfoDetail,
      };
    }

    static setUEInfoDetailAMF(AMFDetail) {
      return {
        type: this.SET_UE_DETAIL_AMF,
        amfInfo: AMFDetail
      };
    }

    static setUEInfoDetailSMF(SMFDetail) {
      return {
        type: this.SET_UE_DETAIL_SMF,
        smfInfo: SMFDetail
      };
    }

    static setUEInfoDetailSmContextRef(smContextRef) {
      return {
        type: this.SET_UE_DETAIL_SM_CTX_REF,
        smContextRef: smContextRef
      };

    }

    static setRegisteredUEError(errMsg) {
      return {
        type: this.SET_REG_UE_ERR,
        get_registered_ue_err: true,
        registered_ue_err_msg: errMsg
      };

    }
  }