package consumer

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/free5gc/openapi"
	udm_context "github.com/free5gc/udm/internal/context"
	"github.com/free5gc/udm/pkg/app"
)

func TestSendRegisterNFInstance(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	openapi.InterceptH2CClient()
	defer openapi.RestoreH2CClient()

	gock.New("http://127.0.0.10:8000").
		Put("/nnrf-nfm/v1/nf-instances/1").
		Reply(200).
		JSON(map[string]string{})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApp := app.NewMockApp(ctrl)
	consumer, err := NewConsumer(mockApp)
	require.NoError(t, err)

	mockApp.EXPECT().Context().Times(1).Return(
		&udm_context.UDMContext{
			NrfUri: "http://127.0.0.10:8000",
			NfId:   "1",
		},
	)

	_, _, err = consumer.RegisterNFInstance(context.TODO())
	require.NoError(t, err)
}
