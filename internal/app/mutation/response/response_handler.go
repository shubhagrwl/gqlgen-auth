package response

// func GetWho(ctx *fasthttp.RequestCtx) (*models.Who, error) {
// 	who := &models.Who{}
// 	if err := json.Unmarshal(ctx.Request.Header.Peek("who"), who); err != nil {
// 		return nil, err
// 	}
// 	if who.Token == "" {
// 		if strings.Contains(ctx.URI().String(), "sign-in") {
// 			return who, nil
// 		}
// 		return nil, errors.New("Token is not present in the header")
// 	}
// 	return who, nil
// }

// func ExtractToken(ctx *fasthttp.RequestCtx) (string, error) {
// 	authHeader := string(ctx.Request.Header.Peek("Authorization"))
// 	authHeaderContent := strings.Split(authHeader, " ")
// 	if len(authHeaderContent) != 2 {
// 		return "", errors.New("Token not provided or malformed")
// 	}
// 	return authHeaderContent[1], nil
// }

// func ResponseMethod(ctx *fasthttp.RequestCtx, code int, body interface{}) {
// 	type response struct {
// 		Error interface{} `json:"error"`
// 	}
// 	log := logger.Logger(ctx)
// 	if code >= 300 {
// 		var b response
// 		b.Error = body
// 		ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
// 		ctx.Response.SetStatusCode(code)
// 		err := json.NewEncoder(ctx).Encode(b)
// 		if err != nil {
// 			log.Error("Getting Error while encoding the response (error) :", err)
// 			return
// 		}
// 	} else {
// 		ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
// 		ctx.Response.SetStatusCode(code)
// 		err := json.NewEncoder(ctx).Encode(body)
// 		if err != nil {
// 			log.Error("Getting Error while encoding the response (error) :", err)
// 			return
// 		}
// 	}

// 	// start := ctx.UserValue("start").(time.Time)
// 	// log.Info("Handled request in " + time.Since(start).String())
// }
