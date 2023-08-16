package common

import "errors"

var ErrNotFoundTopic = errors.New("TopicId: NotFound")
var ErrValidateBytes = errors.New("byte size is not valid")
