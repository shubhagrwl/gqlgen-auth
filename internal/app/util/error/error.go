package utils

// // HandleError is the method for returning the user facing error
// func HandleError(ctx context.Context, errType string, err error) error {
// 	log := logger.Logger(ctx)

// 	req, _ := utilContext.RequestFromContext(ctx)
// 	log.WithField("request", req).WithError(err).Error(errType)

// 	if config.BuildEnv() != "dev" && (errType == constants.InternalServerError) {
// 		notice := instance.Airbrake().Notice(err, nil, 0)
// 		notice.Context["body"] = req
// 		userID, err := utilContext	.UserIDFromContext(ctx)
// 		if err == nil {
// 			notice.Context["userID"] = userID
// 		}
// 		instance.Airbrake().Notify(notice, nil)
// 	}

// 	if errType != constants.InvalidRequestData {
// 		return errors.New(constants.ErrorString[errType])
// 	}
// 	if err == orm.ErrNoRows {
// 		err = errors.New(constants.NotFound)
// 	}

// 	errToReturn := constants.ErrorString[err.Error()]

// 	if errToReturn == "" {
// 		errToReturn = constants.ErrorString[errType]
// 	}

// 	return errors.New(errToReturn)

// }
