package api4schedulus

import (
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/facade4schedulus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

func httpRemoveParticipantFromHappening(w http.ResponseWriter, r *http.Request) {
	var request dto4schedulus.HappeningContactRequest
	request.HappeningRequest = getHappeningRequestParamsFromURL(r)
	ctx, userContext, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = facade4schedulus.RemoveParticipantFromHappening(ctx, userContext.GetID(), request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, nil)
}
