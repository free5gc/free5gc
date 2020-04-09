import actions from '../actions/ueinfoActions';

const initialState = {
    ueInfoDetail: {
        amfInfo: {
        },
    
        smfInfo: {
        }, 
    
        pcfInfo: {
            AmPolicyData:  [{
                PolicyAssociationID: "string",
                AccessType: "string",
                Triggers: ["string"],
                Rfsp: "int32",
                RestrictionType: "string",
                Areas: [{
                    // only contain Tacs or AreaCodes
                    Tacs: ["string"],
                    AreaCodes: "string"
                }]
            }]
        }
    }, 
    amfInfo: {
        AccessType: "3GPP",
        Supi: "imsi-2089300007487",
        Guti: "guti-2089300007487",
        mcc: "123",
        mnc: "456",
        tac: "1",
        PduSessions: [{
          PduSessionId: "int",
          smContextRef: "string",
          sst: "int",
          sd: "string",
          Dnn: "internet", 
        }],
        CmState: "string" // CONNECTED or IDLE
    },

    smfInfo: {
        smContext: {
            Supi: "string",
            LocalSEID:    "string",
            RemoteSEID:   "string",
            PDUSessionID: "string",
            PduAddress: "string",
            AnType: "string",
            PDUAddress: "string",
            SessionRule: ["models.SessionRule"],
            Tunnel: "smf_context.UPTunnel",
        }
    },

    registered_users: [],
    get_registered_ue_err: false,
    registered_ue_err_msg: '',
    smContextRef: '',
    
  };

export default function reducer(state = initialState, action) {
let nextState = {...state};

switch (action.type) {
    case actions.SET_REG_UE:
        nextState.registered_users = action.registered_users;
        return nextState;

    case actions.SET_UE_DETAIL:
        nextState.ueInfoDetail = action.ueInfoDetail;
        return nextState;

    case actions.SET_REG_UE_ERR:
        nextState.get_registered_ue_err = action.get_registered_ue_err
        nextState.registered_ue_err_msg = action.registered_ue_err_msg
        return nextState;

    case actions.SET_UE_DETAIL_AMF:
        nextState.amfInfo = action.amfInfo
        return nextState;

    case actions.SET_UE_DETAIL_SMF:
        nextState.smfInfo = action.smfInfo
        return nextState;

    case actions.SET_UE_DETAIL_SM_CTX_REF:
        nextState.smContextRef = action.smContextRef
        return nextState;
    default:
    return state;
}
};
