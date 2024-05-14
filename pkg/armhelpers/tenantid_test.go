// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package armhelpers

// func TestGetTenantID(t *testing.T) {
// 	fakeServer := fake.Server{
// 		Get: func(ctx context.Context, subscriptionID string, options *azarm.ClientGetOptions) (resp azfake.Responder[azarm.ClientGetResponse], errResp azfake.ErrorResponder) {
// 			errResp.SetError(&azcore.ResponseError{
// 				StatusCode: http.StatusUnauthorized,
// 				RawResponse: &http.Response{
// 					StatusCode: http.StatusUnauthorized,
// 				},
// 			})
// 			return
// 		},
// 	}

// 	tenantID, err := GetTenantID(&azfake.TokenCredential{}, "foobarsubscription", &arm.ClientOptions{
// 		ClientOptions: azcore.ClientOptions{
// 			Transport: fake.NewServerTransport(&fakeServer),
// 			PerCallPolicies: []policy.Policy{
// 				runtime.NewBearerTokenPolicy(&azfake.TokenCredential{}, []string{"https://management.azure.com/"}, &policy.BearerTokenOptions{
// 					AuthorizationHandler: policy.AuthorizationHandler{
// 						OnRequest: func(req *policy.Request, _ func(policy.TokenRequestOptions) error) error {
// 							req.Raw().Header.Set("Authorization", "Bearer fake_token")
// 							return nil
// 						},
// 						// 			OnChallenge: func(_ *policy.Request, resp *http.Response, authNZ func(policy.TokenRequestOptions) error) error {
// 						// 				resp.Header.Add("WWW-Authenticate", `authorization_uri="https://login.windows.net/faketenantid"`)
// 						// 				return runtime.NewResponseError(resp)
// 						// 			},
// 					},
// 				}),
// 			},
// 		},
// 	})

// 	if err != nil {
// 		t.Error("Did not expect error")
// 	}
// 	if tenantID != "faketenantid" {
// 		t.Errorf("expected tenant Id : %s, but got %s", "faketenantid", tenantID)
// 	}
// }

// func TestGetTenantID_UnexpectedResponse(t *testing.T) {
// 	mux := http.NewServeMux()
// 	server := httptest.NewServer(mux)
// 	defer server.Close()

// 	mux.HandleFunc("/subscriptions/foobarsubscription", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusBadRequest)
// 	})

// 	_, err := GetTenantID(&azfake.TokenCredential{}, "foobarsubscription", nil)

// 	expectedMsg := "Unexpected response from Get Subscription: 400"

// 	if err == nil || err.Error() != expectedMsg {
// 		t.Errorf("expected error with msg : %s to be thrown", expectedMsg)
// 	}
// }

// func TestGetTenantID_InvalidHeader(t *testing.T) {
// 	mux := http.NewServeMux()
// 	server := httptest.NewServer(mux)
// 	defer server.Close()

// 	mux.HandleFunc("/subscriptions/foobarsubscription", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Header().Set("fookey", "bazvalue")
// 	})

// 	_, err := GetTenantID(&azfake.TokenCredential{}, "foobarsubscription", nil)

// 	expectedMsg := "Header WWW-Authenticate not found in Get Subscription response"

// 	if err == nil || err.Error() != expectedMsg {
// 		t.Errorf("expected error with msg : %s to be thrown", expectedMsg)
// 	}
// }

// func TestGetTenantID_InvalidHeaderValue(t *testing.T) {
// 	mux := http.NewServeMux()
// 	server := httptest.NewServer(mux)
// 	defer server.Close()

// 	mux.HandleFunc("/subscriptions/foobarsubscription", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("WWW-Authenticate", `sample_invalid_auth_uri`)
// 		w.WriteHeader(http.StatusUnauthorized)
// 		_, _ = w.Write([]byte("Unauthorized"))
// 	})

// 	_, err := GetTenantID(&azfake.TokenCredential{}, "foobarsubscription", nil)

// 	expectedMsg := "Could not find the tenant ID in header: WWW-Authenticate \"sample_invalid_auth_uri\""

// 	if err == nil || err.Error() != expectedMsg {
// 		t.Errorf("expected error with msg : %s to be thrown", expectedMsg)
// 	}
// }
