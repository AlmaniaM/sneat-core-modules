package api4memberus

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/memberus/facade4memberus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var createMember = facade4memberus.CreateMember

// httpPostCreateMember is an API endpoint that adds a members to a team
func httpPostCreateMember(w http.ResponseWriter, r *http.Request) {
	var request dal4contactus.CreateMemberRequest
	handler := func(ctx context.Context, userCtx facade.User) (interface{}, error) {
		return createMember(ctx, userCtx, request)
	}
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, handler, http.StatusCreated, verify.DefaultJsonWithAuthRequired)
}
