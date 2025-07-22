package processor

import (
	"net/http"
	"os"
	"testing"

	nef_context "github.com/free5gc/nef/internal/context"
	"github.com/free5gc/nef/internal/sbi/consumer"
	"github.com/free5gc/nef/internal/sbi/notifier"
	"github.com/free5gc/nef/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

type nefTestApp struct {
	cfg      *factory.Config
	nefCtx   *nef_context.NefContext
	consumer *consumer.Consumer
	notifier *notifier.Notifier
	proc     *Processor
}

func newTestApp(cfg *factory.Config, tlsKeyLogPath string) (*nefTestApp, error) {
	var err error
	nef := &nefTestApp{cfg: cfg}

	if nef.nefCtx, err = nef_context.NewContext(nef); err != nil {
		return nil, err
	}
	if nef.consumer, err = consumer.NewConsumer(nef); err != nil {
		return nil, err
	}
	if nef.notifier, err = notifier.NewNotifier(); err != nil {
		return nil, err
	}
	if nef.proc, err = NewProcessor(nef); err != nil {
		return nil, err
	}
	return nef, nil
}

func (a *nefTestApp) Config() *factory.Config {
	return a.cfg
}

func (a *nefTestApp) Context() *nef_context.NefContext {
	return a.nefCtx
}

func (a *nefTestApp) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *nefTestApp) Notifier() *notifier.Notifier {
	return a.notifier
}

func (a *nefTestApp) Processor() *Processor {
	return a.proc
}

var (
	nefApp *nefTestApp

	pfd1 = models.Pfd{
		PfdId: "pfd1",
		FlowDescriptions: []string{
			"permit in ip from 10.68.28.39 80 to any",
			"permit out ip from any to 10.68.28.39 80",
		},
	}
	pfd2 = models.Pfd{
		PfdId: "pfd2",
		Urls: []string{
			"^http://test.example.com(/\\S*)?$",
		},
	}
	pfd3 = models.Pfd{
		PfdId: "pfd3",
		Urls: []string{
			"^http://test.example2.net(/\\S*)?$",
		},
	}

	pfdDataForApp1 = models.PfdDataForApp{
		ApplicationId: "app1",
		Pfds: []models.PfdContent{
			{
				PfdId: "pfd1",
				FlowDescriptions: []string{
					"permit in ip from 10.68.28.39 80 to any",
					"permit out ip from any to 10.68.28.39 80",
				},
			},
			{
				PfdId: "pfd2",
				Urls: []string{
					"^http://test.example.com(/\\S*)?$",
				},
			},
		},
	}
	pfdDataForApp2 = models.PfdDataForApp{
		ApplicationId: "app2",
		Pfds: []models.PfdContent{
			{
				PfdId: "pfd3",
				Urls: []string{
					"^http://test.example2.net(/\\S*)?$",
				},
			},
		},
	}
)

func TestMain(m *testing.M) {
	var err error
	openapi.InterceptH2CClient()
	initNRFNfmStub()
	initNRFDiscUDRStub()

	cfg := &factory.Config{
		Info: &factory.Info{
			Version: "1.0.0",
		},
		Configuration: &factory.Configuration{
			Sbi: &factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.5",
				BindingIPv4:  "127.0.0.5",
				Port:         8000,
			},
			NrfUri: "http://127.0.0.10:8000",
			ServiceList: []factory.Service{
				{
					ServiceName: factory.ServiceNefPfd,
				},
			},
		},
	}
	nefApp, err = newTestApp(cfg, "")
	if err != nil {
		panic(err)
	}

	exitVal := m.Run()
	openapi.RestoreH2CClient()
	os.Exit(exitVal)
}

func TestGetPFDManagementTransactions(t *testing.T) {
	initUDRDrGetPfdDatasStub()
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Valid input",
			afID:        "af1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body: &[]models.PfdManagement{
					{
						Self: nefApp.Processor().genPfdManagementURI("af1", "1"),
						PfdDatas: map[string]models.PfdData{
							"app1": {
								ExternalAppId: "app1",
								Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app1"),
								Pfds: map[string]models.Pfd{
									"pfd1": pfd1,
									"pfd2": pfd2,
								},
							},
							"app2": {
								ExternalAppId: "app2",
								Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app2"),
								Pfds: map[string]models.Pfd{
									"pfd3": pfd3,
								},
							},
						},
					},
				},
			},
		},
		{
			description: "TC2: Invalid AF ID, should return ProblemDetails",
			afID:        "af2",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound(DetailNoAF),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			af.Mu.Lock()
			afPfdTr := af.NewPfdTrans()
			af.PfdTrans[afPfdTr.TransID] = afPfdTr
			afPfdTr.AddExtAppID("app1")
			afPfdTr.AddExtAppID("app2")
			af.Mu.Unlock()

			rsp := nefApp.Processor().GetPFDManagementTransactions(tc.afID)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestDeletePFDManagementTransactions(t *testing.T) {
	initUDRDrDeletePfdDataStub()
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Valid input",
			afID:        "af1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNoContent,
			},
		},
		{
			description: "TC2: Invalid AF ID, should return ProblemDetails",
			afID:        "af2",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound(DetailNoAF),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			af.Mu.Lock()
			afPfdTr := af.NewPfdTrans()
			af.PfdTrans[afPfdTr.TransID] = afPfdTr
			af.Mu.Unlock()

			rsp := nefApp.Processor().DeletePFDManagementTransactions(tc.afID)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestPostPFDManagementTransactions(t *testing.T) {
	initUDRDrPutPfdDataStub(http.StatusCreated)
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		pfdManagement    *models.PfdManagement
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Valid input",
			afID:        "af1",
			pfdManagement: &models.PfdManagement{
				PfdDatas: map[string]models.PfdData{
					"app1": {
						ExternalAppId: "app1",
						Pfds: map[string]models.Pfd{
							"pfd1": pfd1,
							"pfd2": pfd2,
						},
					},
					"app2": {
						ExternalAppId: "app2",
						Pfds: map[string]models.Pfd{
							"pfd3": pfd3,
						},
					},
				},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusCreated,
				Body: &models.PfdManagement{
					Self: nefApp.Processor().genPfdManagementURI("af1", "1"),
					PfdDatas: map[string]models.PfdData{
						"app1": {
							ExternalAppId: "app1",
							Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app1"),
							Pfds: map[string]models.Pfd{
								"pfd1": pfd1,
								"pfd2": pfd2,
							},
						},
						"app2": {
							ExternalAppId: "app2",
							Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app2"),
							Pfds: map[string]models.Pfd{
								"pfd3": pfd3,
							},
						},
					},
					PfdReports: map[string]models.PfdReport{},
				},
			},
		},
		{
			description: "TC2: Invalid AF ID, should return ProblemDetails",
			afID:        "af2",
			pfdManagement: &models.PfdManagement{
				PfdDatas: map[string]models.PfdData{
					"app1": {
						ExternalAppId: "app1",
						Pfds: map[string]models.Pfd{
							"pfd1": pfd1,
							"pfd2": pfd2,
						},
					},
					"app2": {
						ExternalAppId: "app2",
						Pfds: map[string]models.Pfd{
							"pfd3": pfd3,
						},
					},
				},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound(DetailNoAF),
			},
		},
		{
			description: "Invalid PfdManagement, should return ProblemDetails",
			afID:        "af1",
			pfdManagement: &models.PfdManagement{
				PfdDatas: map[string]models.PfdData{},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound(DetailNoPfdData),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			rsp := nefApp.Processor().PostPFDManagementTransactions(tc.afID, tc.pfdManagement)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestGetIndividualPFDManagementTransaction(t *testing.T) {
	initUDRDrGetPfdDatasStub()
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		transID          string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Valid input",
			afID:        "af1",
			transID:     "1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body: &models.PfdManagement{
					Self: nefApp.Processor().genPfdManagementURI("af1", "1"),
					PfdDatas: map[string]models.PfdData{
						"app1": {
							ExternalAppId: "app1",
							Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app1"),
							Pfds: map[string]models.Pfd{
								"pfd1": pfd1,
								"pfd2": pfd2,
							},
						},
						"app2": {
							ExternalAppId: "app2",
							Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app2"),
							Pfds: map[string]models.Pfd{
								"pfd3": pfd3,
							},
						},
					},
				},
			},
		},
		{
			description: "TC2: Invalid transaction ID, should return ProblemDetails",
			afID:        "af1",
			transID:     "-1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound("PFD transaction not found"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			af.Mu.Lock()
			afPfdTr := af.NewPfdTrans()
			af.PfdTrans[afPfdTr.TransID] = afPfdTr
			afPfdTr.AddExtAppID("app1")
			afPfdTr.AddExtAppID("app2")
			af.Mu.Unlock()

			rsp := nefApp.Processor().GetIndividualPFDManagementTransaction(tc.afID, tc.transID)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestDeleteIndividualPFDManagementTransaction(t *testing.T) {
	initUDRDrDeletePfdDataStub()
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		transID          string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Valid input",
			afID:        "af1",
			transID:     "1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNoContent,
			},
		},
		{
			description: "TC2: Invalid transaction ID, should return ProblemDetails",
			afID:        "af2",
			transID:     "1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound("AF not found"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			af.Mu.Lock()
			afPfdTr := af.NewPfdTrans()
			af.PfdTrans[afPfdTr.TransID] = afPfdTr
			af.Mu.Unlock()

			rsp := nefApp.Processor().DeleteIndividualPFDManagementTransaction(tc.afID, tc.transID)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestPutIndividualPFDManagementTransaction(t *testing.T) {
	initUDRDrPutPfdDataStub(http.StatusOK)
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		transID          string
		pfdManagement    *models.PfdManagement
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Valid input",
			afID:        "af1",
			transID:     "1",
			pfdManagement: &models.PfdManagement{
				PfdDatas: map[string]models.PfdData{
					"app1": {
						ExternalAppId: "app1",
						Pfds: map[string]models.Pfd{
							"pfd1": pfd1,
							"pfd2": pfd2,
						},
					},
					"app2": {
						ExternalAppId: "app2",
						Pfds: map[string]models.Pfd{
							"pfd3": pfd3,
						},
					},
				},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body: &models.PfdManagement{
					Self: nefApp.Processor().genPfdManagementURI("af1", "1"),
					PfdDatas: map[string]models.PfdData{
						"app1": {
							ExternalAppId: "app1",
							Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app1"),
							Pfds: map[string]models.Pfd{
								"pfd1": pfd1,
								"pfd2": pfd2,
							},
						},
						"app2": {
							ExternalAppId: "app2",
							Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app2"),
							Pfds: map[string]models.Pfd{
								"pfd3": pfd3,
							},
						},
					},
					PfdReports: map[string]models.PfdReport{},
				},
			},
		},
		{
			description: "TC2: Invalid transaction ID, should return ProblemDetails",
			afID:        "af1",
			transID:     "-1",
			pfdManagement: &models.PfdManagement{
				PfdDatas: map[string]models.PfdData{
					"app1": {
						ExternalAppId: "app1",
						Pfds: map[string]models.Pfd{
							"pfd1": pfd1,
							"pfd2": pfd2,
						},
					},
					"app2": {
						ExternalAppId: "app2",
						Pfds: map[string]models.Pfd{
							"pfd3": pfd3,
						},
					},
				},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound("PFD transaction not found"),
			},
		},
		{
			description: "TC3: Invalid PfdManagement, should return ProblemDetails",
			afID:        "af1",
			transID:     "1",
			pfdManagement: &models.PfdManagement{
				PfdDatas: map[string]models.PfdData{},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound(DetailNoPfdData),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			af.Mu.Lock()
			afPfdTr := af.NewPfdTrans()
			af.PfdTrans[afPfdTr.TransID] = afPfdTr
			af.Mu.Unlock()

			rsp := nefApp.Processor().PutIndividualPFDManagementTransaction(tc.afID, tc.transID, tc.pfdManagement)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestGetIndividualApplicationPFDManagement(t *testing.T) {
	initUDRDrGetPfdDataStub()
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		transID          string
		appID            string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Valid input",
			afID:        "af1",
			transID:     "1",
			appID:       "app1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body: &models.PfdData{
					ExternalAppId: "app1",
					Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app1"),
					Pfds: map[string]models.Pfd{
						"pfd1": pfd1,
						"pfd2": pfd2,
					},
				},
			},
		},
		{
			description: "TC2: Invalid App ID, should return ProblemDetails",
			afID:        "af1",
			transID:     "1",
			appID:       "app2",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound("Application ID not found"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			af.Mu.Lock()
			afPfdTr := af.NewPfdTrans()
			af.PfdTrans[afPfdTr.TransID] = afPfdTr
			afPfdTr.AddExtAppID("app1")
			af.Mu.Unlock()

			rsp := nefApp.Processor().GetIndividualApplicationPFDManagement(tc.afID, tc.transID, tc.appID)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestDeleteIndividualApplicationPFDManagement(t *testing.T) {
	initUDRDrDeletePfdDataStub()
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		transID          string
		appID            string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Valid input",
			afID:        "af1",
			transID:     "1",
			appID:       "app1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNoContent,
			},
		},
		{
			description: "TC2: Invalid App ID, should return ProblemDetails",
			afID:        "af1",
			transID:     "1",
			appID:       "app2",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound("Application ID not found"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			af.Mu.Lock()
			afPfdTr := af.NewPfdTrans()
			af.PfdTrans[afPfdTr.TransID] = afPfdTr
			afPfdTr.AddExtAppID("app1")
			af.Mu.Unlock()

			rsp := nefApp.Processor().DeleteIndividualApplicationPFDManagement(tc.afID, tc.transID, tc.appID)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestPutIndividualApplicationPFDManagement(t *testing.T) {
	initUDRDrPutPfdDataStub(http.StatusOK)
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		transID          string
		appID            string
		pfdData          *models.PfdData
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Valid input",
			afID:        "af1",
			transID:     "1",
			appID:       "app1",
			pfdData: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": pfd1,
					"pfd2": pfd2,
				},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body: &models.PfdData{
					ExternalAppId: "app1",
					Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app1"),
					Pfds: map[string]models.Pfd{
						"pfd1": pfd1,
						"pfd2": pfd2,
					},
				},
			},
		},
		{
			description: "TC2: Invalid App ID, should return ProblemDetails",
			afID:        "af1",
			transID:     "1",
			appID:       "app2",
			pfdData: &models.PfdData{
				ExternalAppId: "app2",
				Pfds: map[string]models.Pfd{
					"pfd1": pfd1,
					"pfd2": pfd2,
				},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound("Application ID not found"),
			},
		},
		{
			description: "TC3: Invalid PfdData, should return ProblemDetails",
			afID:        "af1",
			transID:     "1",
			appID:       "app1",
			pfdData: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": {
						PfdId: "pfd1",
					},
				},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound(DetailNoPfdInfo),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			af.Mu.Lock()
			afPfdTr := af.NewPfdTrans()
			af.PfdTrans[afPfdTr.TransID] = afPfdTr
			afPfdTr.AddExtAppID("app1")
			af.Mu.Unlock()

			rsp := nefApp.Processor().PutIndividualApplicationPFDManagement(tc.afID, tc.transID, tc.appID, tc.pfdData)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestPatchIndividualApplicationPFDManagement(t *testing.T) {
	initUDRDrGetPfdDataStub()
	initUDRDrPutPfdDataStub(http.StatusOK)
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		transID          string
		appID            string
		pfdData          *models.PfdData
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Valid input",
			afID:        "af1",
			transID:     "1",
			appID:       "app1",
			pfdData: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": {
						PfdId: "pfd1",
					},
				},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body: &models.PfdData{
					ExternalAppId: "app1",
					Self:          nefApp.Processor().genPfdDataURI("af1", "1", "app1"),
					Pfds: map[string]models.Pfd{
						"pfd2": pfd2,
					},
				},
			},
		},
		{
			description: "TC2: Invalid App ID, should return ProblemDetails",
			afID:        "af1",
			transID:     "1",
			appID:       "app2",
			pfdData: &models.PfdData{
				ExternalAppId: "app2",
				Pfds: map[string]models.Pfd{
					"pfd1": {
						PfdId: "pfd1",
					},
				},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound("Application ID not found"),
			},
		},
		{
			description: "TC3: Invalid PfdData, should return ProblemDetails",
			afID:        "af1",
			transID:     "1",
			appID:       "app1",
			pfdData: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd3": {
						PfdId: "pfd3",
					},
				},
			},
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ProblemDetailsDataNotFound(DetailNoPfdInfo),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			af.Mu.Lock()
			afPfdTr := af.NewPfdTrans()
			af.PfdTrans[afPfdTr.TransID] = afPfdTr
			afPfdTr.AddExtAppID("app1")
			af.Mu.Unlock()

			rsp := nefApp.Processor().PatchIndividualApplicationPFDManagement(
				tc.afID, tc.transID, tc.appID, tc.pfdData)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestValidatePfdManagement(t *testing.T) {
	testCases := []struct {
		description     string
		pfdManagement   *models.PfdManagement
		expectedProblem *models.ProblemDetails
		expectedReports map[string]models.PfdReport
	}{
		{
			description: "TC1: Valid",
			pfdManagement: &models.PfdManagement{
				PfdDatas: map[string]models.PfdData{
					"app1": {
						ExternalAppId: "app1",
						Pfds: map[string]models.Pfd{
							"pfd1": pfd1,
							"pfd2": pfd2,
						},
					},
					"app2": {
						ExternalAppId: "app2",
						Pfds: map[string]models.Pfd{
							"pfd3": pfd3,
						},
					},
				},
			},
			expectedProblem: nil,
			expectedReports: map[string]models.PfdReport{},
		},
		{
			description: "TC2: Empty PfdDatas, should return ProblemDetails",
			pfdManagement: &models.PfdManagement{
				PfdDatas: map[string]models.PfdData{},
			},
			expectedProblem: openapi.ProblemDetailsDataNotFound(DetailNoPfdData),
			expectedReports: map[string]models.PfdReport{},
		},
		{
			description: "TC3: An appID is already provisioned, should mark in PfdReports",
			pfdManagement: &models.PfdManagement{
				PfdDatas: map[string]models.PfdData{
					"app100": {
						ExternalAppId: "app100",
						Pfds: map[string]models.Pfd{
							"pfd1": pfd1,
						},
					},
					"app101": {
						ExternalAppId: "app101",
						Pfds: map[string]models.Pfd{
							"pfd1": pfd1,
						},
					},
				},
			},
			expectedProblem: nil,
			expectedReports: map[string]models.PfdReport{
				string(models.FailureCode_APP_ID_DUPLICATED): {
					ExternalAppIds: []string{"app100"},
					FailureCode:    models.FailureCode_APP_ID_DUPLICATED,
				},
			},
		},
		{
			description: "TC4: None of the PFDs were created, should return ProblemDetails and mark in PfdReports",
			pfdManagement: &models.PfdManagement{
				PfdDatas: map[string]models.PfdData{
					"app100": {
						ExternalAppId: "app100",
						Pfds: map[string]models.Pfd{
							"pfd1": pfd1,
						},
					},
				},
			},
			expectedProblem: openapi.ProblemDetailsSystemFailure("None of the PFDs were created"),
			expectedReports: map[string]models.PfdReport{
				string(models.FailureCode_APP_ID_DUPLICATED): {
					ExternalAppIds: []string{"app100"},
					FailureCode:    models.FailureCode_APP_ID_DUPLICATED,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			af := nefApp.Context().NewAf("af1")
			nefApp.Context().AddAf(af)
			defer nefApp.Context().DeleteAf("af1")

			af.Mu.Lock()
			afPfdTr := af.NewPfdTrans()
			af.PfdTrans[afPfdTr.TransID] = afPfdTr
			afPfdTr.AddExtAppID("app100")
			af.Mu.Unlock()

			rst := validatePfdManagement("af2", "1", tc.pfdManagement, nefApp.Context())
			require.Equal(t, tc.expectedProblem, rst)
			require.Equal(t, tc.expectedReports, tc.pfdManagement.PfdReports)
		})
	}
}

func TestValidatePfdData(t *testing.T) {
	testCases := []struct {
		description    string
		pfdData        *models.PfdData
		expectedResult *models.ProblemDetails
	}{
		{
			description: "TC1: Valid",
			pfdData: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": pfd1,
				},
			},
			expectedResult: nil,
		},
		{
			description: "TC2: Without ExternalAppId, should return ProblemDetails",
			pfdData: &models.PfdData{
				Pfds: map[string]models.Pfd{
					"pfd1": pfd1,
				},
			},
			expectedResult: openapi.ProblemDetailsDataNotFound(DetailNoExtAppID),
		},
		{
			description: "TC3: Empty Pfds, should return ProblemDetails",
			pfdData: &models.PfdData{
				ExternalAppId: "app1",
			},
			expectedResult: openapi.ProblemDetailsDataNotFound(DetailNoPfd),
		},
		{
			description: "TC4: Without PfdID, should return ProblemDetails",
			pfdData: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": {
						FlowDescriptions: []string{
							"permit in ip from 10.68.28.39 80 to any",
							"permit out ip from any to 10.68.28.39 80",
						},
					},
				},
			},
			expectedResult: openapi.ProblemDetailsDataNotFound(DetailNoPfdID),
		},
		{
			description: "TC5: FlowDescriptions, Urls and DomainNames are all empty, should return ProblemDetails",
			pfdData: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": {
						PfdId: "pfd1",
					},
				},
			},
			expectedResult: openapi.ProblemDetailsDataNotFound(DetailNoPfdInfo),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rst := validatePfdData(tc.pfdData, nefApp.Context(), false)
			require.Equal(t, tc.expectedResult, rst)
		})
	}
}

func TestPatchModifyPfdData(t *testing.T) {
	testCases := []struct {
		description     string
		old             *models.PfdData
		new             *models.PfdData
		expectedProblem *models.ProblemDetails
		expectedResult  *models.PfdData
	}{
		{
			description: "TC1: Given a PfdData with non-existing appID, should append the Pfds to the PfdData",
			old: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": pfd1,
				},
			},
			new: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd2": pfd2,
				},
			},
			expectedProblem: nil,
			expectedResult: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": pfd1,
					"pfd2": pfd2,
				},
			},
		},
		{
			description: "TC2: Given a PfdData with existing appID, should update the PfdData",
			old: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": pfd1,
				},
			},
			new: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": {
						PfdId: "pfd1",
						Urls: []string{
							"^http://test.example.com(/\\S*)?$",
						},
					},
				},
			},
			expectedProblem: nil,
			expectedResult: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": {
						PfdId: "pfd1",
						Urls: []string{
							"^http://test.example.com(/\\S*)?$",
						},
					},
				},
			},
		},
		{
			description: "TC3: Given a PfdData with existing appID and empty content, should delete the PfdData",
			old: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": pfd1,
					"pfd2": pfd2,
				},
			},
			new: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": {
						PfdId: "pfd1",
					},
				},
			},
			expectedProblem: nil,
			expectedResult: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd2": pfd2,
				},
			},
		},
		{
			description: "TC4: Given an invalid PfdData, should return ProblemDetails",
			old: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": pfd1,
				},
			},
			new: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd2": {
						PfdId: "pfd2",
					},
				},
			},
			expectedProblem: openapi.ProblemDetailsDataNotFound(DetailNoPfdInfo),
			expectedResult: &models.PfdData{
				ExternalAppId: "app1",
				Pfds: map[string]models.Pfd{
					"pfd1": pfd1,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			problemDetail := patchModifyPfdData(tc.old, tc.new)
			require.Equal(t, tc.expectedProblem, problemDetail)
			require.Equal(t, tc.expectedResult, tc.old)
		})
	}
}

func initNRFNfmStub() {
	nrfRegisterInstanceRsp := models.NfProfile{
		NfInstanceId: "nef-pfd-unit-testing",
	}
	gock.New("http://127.0.0.10:8000/nnrf-nfm/v1").
		Put("/nf-instances/.*").
		MatchType("json").
		JSON(".*").
		Reply(http.StatusCreated).
		SetHeader("Location", "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances/12345").
		JSON(nrfRegisterInstanceRsp)
}

func initNRFDiscUDRStub() {
	searchResult := &models.SearchResult{
		ValidityPeriod: 100,
		NfInstances: []models.NfProfile{
			{
				NfInstanceId: "nef-unit-testing",
				NfType:       "UDR",
				NfStatus:     "REGISTERED",
				UdrInfo: &models.UdrInfo{
					SupportedDataSets: []models.DataSetId{
						"SUBSCRIPTION",
					},
				},
				NfServices: &[]models.NfService{
					{
						ServiceInstanceId: "datarepository",
						ServiceName:       "nudr-dr",
						Versions: &[]models.NfServiceVersion{
							{
								ApiVersionInUri: "v1",
								ApiFullVersion:  "1.0.0",
							},
						},
						Scheme:          "http",
						NfServiceStatus: "REGISTERED",
						IpEndPoints: &[]models.IpEndPoint{
							{
								Ipv4Address: "127.0.0.4",
								Transport:   "TCP",
								Port:        8000,
							},
						},
						ApiPrefix: "http://127.0.0.4:8000",
					},
				},
			},
		},
	}

	gock.New("http://127.0.0.10:8000/nnrf-disc/v1").
		Get("/nf-instances").
		MatchParam("target-nf-type", "UDR").
		MatchParam("requester-nf-type", "NEF").
		MatchParam("service-names", "nudr-dr").
		Reply(http.StatusOK).
		JSON(searchResult)
}

func initNRFDiscPCFStub() {
	searchResult := &models.SearchResult{
		ValidityPeriod: 100,
		NfInstances: []models.NfProfile{
			{
				NfInstanceId: "nef-unit-testing",
				NfType:       "PCF",
				NfStatus:     "REGISTERED",
				Ipv4Addresses: []string{
					"127.0.0.7",
				},
				PcfInfo: &models.PcfInfo{
					DnnList: []string{
						"free5gc",
						"internet",
					},
				},
				NfServices: &[]models.NfService{
					{
						ServiceInstanceId: "1",
						ServiceName:       "npcf-policyauthorization",
						Versions: &[]models.NfServiceVersion{
							{
								ApiVersionInUri: "v1",
								ApiFullVersion:  "1.0.0",
							},
						},
						Scheme:          "http",
						NfServiceStatus: "REGISTERED",
						IpEndPoints: &[]models.IpEndPoint{
							{
								Ipv4Address: "127.0.0.7",
								Transport:   "TCP",
								Port:        8000,
							},
						},
						ApiPrefix: "http://127.0.0.7:8000",
					},
				},
			},
		},
	}

	gock.New("http://127.0.0.10:8000/nnrf-disc/v1").
		Get("/nf-instances").
		MatchParam("target-nf-type", "PCF").
		MatchParam("requester-nf-type", "NEF").
		MatchParam("service-names", "npcf-policyauthorization").
		Reply(http.StatusOK).
		JSON(searchResult)
}

func initUDRDrGetPfdDatasStub() {
	gock.New("http://127.0.0.4:8000/nudr-dr/v1").
		Get("/application-data/pfds").
		// To Matching the request for both app1 and app2.
		// Should be clarified if there is a way to exact match multiple parameters with the same key.
		MatchParam("appId", "app1").
		Persist().
		Reply(http.StatusOK).
		JSON([]models.PfdDataForApp{pfdDataForApp1, pfdDataForApp2})

	gock.New("http://127.0.0.4:8000/nudr-dr/v1").
		Get("/application-data/pfds").
		MatchParam("appId", "app3").
		Persist().
		Reply(http.StatusNotFound).
		JSON(models.ProblemDetails{Status: http.StatusNotFound})
}

func initUDRDrGetPfdDataStub() {
	gock.New("http://127.0.0.4:8000/nudr-dr/v1").
		Get("/application-data/pfds/app1").
		Persist().
		Reply(http.StatusOK).
		JSON(pfdDataForApp1)

	gock.New("http://127.0.0.4:8000/nudr-dr/v1").
		Get("/application-data/pfds/app3").
		Persist().
		Reply(http.StatusNotFound).
		JSON(models.ProblemDetails{Status: http.StatusNotFound})
}

func initUDRDrDeletePfdDataStub() {
	gock.New("http://127.0.0.4:8000/nudr-dr/v1").
		Delete("/application-data/pfds/.*").
		Persist().
		Reply(http.StatusNoContent)
}

func initUDRDrPutPfdDataStub(statusCode int) {
	gock.New("http://127.0.0.4:8000/nudr-dr/v1").
		Put("/application-data/pfds/.*").
		Persist().
		Reply(statusCode).
		JSON(pfdDataForApp1)
}
